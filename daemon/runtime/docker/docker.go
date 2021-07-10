package docker

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/jsonmessage"
	prototypes "github.com/gogo/protobuf/types"
	moby "github.com/moby/moby/client"
	"github.com/rkrmr33/onka/daemon/collector"
	"github.com/rkrmr33/onka/daemon/command"
	"github.com/rkrmr33/onka/daemon/util"
	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
	log "github.com/sirupsen/logrus"
)

const DefaultDockerTimeout = time.Second * 5

type (
	DockerClientOpts struct {
		// Host the docker host to connect to
		Host string `json:"host"`

		// CACertPath to use for docker connection
		CACertPath string `json:"cacertPath"`

		// CertPath to use for docker connection
		CertPath string `json:"certPath"`

		// KeyPath to use for docker connection
		KeyPath string `json:"keyPath"`

		// Timeout for requests to docker
		Timeout time.Duration `json:"timeout"`

		// Client version to use for daemon communication
		ClientVersion string `json:"clientVersion"`
	}

	DockerClient interface {
		moby.APIClient

		Execute(context.Context, *command.Task) error

		CollectStats(ctx context.Context) (collector.RuntimeStats, error)

		Running(context.Context) (bool, error)
	}

	dockerClientImpl struct {
		moby.APIClient
	}
)

var (
	EnvDockerClient = func() (DockerClient, error) {
		c, err := moby.NewClientWithOpts(moby.FromEnv)
		if err != nil {
			return nil, err
		}

		return &dockerClientImpl{c}, err
	}

	OptsDockerClient = func(d DockerClientOpts) (DockerClient, error) {
		c, err := moby.NewClientWithOpts(moby.WithHost(d.Host), moby.WithTLSClientConfig(d.CACertPath, d.CertPath, d.KeyPath), moby.WithVersion(d.ClientVersion))
		if err != nil {
			return nil, err
		}
		return &dockerClientImpl{c}, err
	}
)

func (c *dockerClientImpl) createAndPullIfNotPresent(ctx context.Context, t *command.Task) (*container.ContainerCreateCreatedBody, error) {
	image := t.Spec.Image
	if !strings.Contains(image, "/") { // need to add library/ prefix
		image = fmt.Sprintf("library/%s", image)
	}

	conf := &container.Config{
		Image: image,
		Tty:   false,
		Env:   t.Spec.Env,
	}

	if t.Spec.Entrypoint != "" {
		conf.Entrypoint = []string{t.Spec.Entrypoint}
	} else {
		shell := t.Spec.Shell
		if shell == "" {
			shell = "sh"
		}

		conf.Entrypoint = []string{shell, "-c"}
	}

	conf.Cmd = []string{util.GetShellCmd(t.Spec.Commands)}

	t.UpdateState(v1alpha1.TaskState_TASK_STATE_PREPARE, "creating container")
	cont, err := c.ContainerCreate(ctx, conf, nil, nil, nil, t.Metadata.Id)
	if err == nil {
		log.WithField("id", t.Metadata.Id).Debug("container created for task")

		return &cont, nil
	}
	if !strings.Contains(err.Error(), "No such image") {
		t.UpdateState(v1alpha1.TaskState_TASK_STATE_ERROR, fmt.Sprintf("failed to create container: %s", err))
		return nil, err
	}

	t.UpdateState(v1alpha1.TaskState_TASK_STATE_PREPARE, "pulling image")

	log.WithField("image", image).Debugf("pulling image")
	progress, err := c.ImagePull(ctx, image, dtypes.ImagePullOptions{})
	if err != nil {
		t.UpdateState(v1alpha1.TaskState_TASK_STATE_ERROR, fmt.Sprintf("failed to pull image: %s", err))

		return nil, err
	}

	if err = jsonmessage.DisplayJSONMessagesStream(progress, t.Logger, os.Stdout.Fd(), true, nil); err != nil {
		return nil, err
	}

	if err = progress.Close(); err != nil {
		return nil, err
	}
	log.WithField("image", image).Debug("pulled image")
	t.UpdateState(v1alpha1.TaskState_TASK_STATE_PREPARE, "image pulled")

	return c.createAndPullIfNotPresent(ctx, t)
}

func (c *dockerClientImpl) Execute(ctx context.Context, t *command.Task) error {
	log.Debugf("starting task %+v", t)

	cont, err := c.createAndPullIfNotPresent(ctx, t)
	if err != nil {
		return err
	}

	err = c.attachAndReportLogs(ctx, t, false)
	if err != nil {
		t.UpdateState(v1alpha1.TaskState_TASK_STATE_ERROR, fmt.Sprintf("failed to attach to container: %s", err))

		return err
	}

	// reporting here to reach before logs
	// TODO: possibly solve this in a better way
	t.UpdateState(v1alpha1.TaskState_TASK_STATE_RUNNING, "container started")
	err = c.ContainerStart(ctx, cont.ID, dtypes.ContainerStartOptions{})
	if err != nil {
		t.UpdateState(v1alpha1.TaskState_TASK_STATE_ERROR, fmt.Sprintf("failed to start container: %s", err))

		return err
	}
	log.Debug("container ", cont.ID, " started")

	poll, errc := c.ContainerWait(ctx, cont.ID, container.WaitConditionNextExit)

	select {
	case <-poll:
		break
	case err := <-errc:
		t.UpdateState(v1alpha1.TaskState_TASK_STATE_ERROR, fmt.Sprintf("failed to wait for container: %s", err))
		return err
	case <-ctx.Done():
		t.UpdateState(v1alpha1.TaskState_TASK_STATE_ERROR, fmt.Sprintf("failed to wait for container: %s", ctx.Err()))
		return ctx.Err()
	}

	insp, err := c.ContainerInspect(ctx, cont.ID)
	if err != nil {
		t.UpdateState(v1alpha1.TaskState_TASK_STATE_ERROR, fmt.Sprintf("failed to inspect container after finish: %s", err))
		return err
	}

	finishState := v1alpha1.TaskState_TASK_STATE_SUCCESS
	if insp.State.ExitCode != 0 {
		finishState = v1alpha1.TaskState_TASK_STATE_FAILURE
	}
	t.UpdateState(finishState, fmt.Sprintf("container exited with exit code: %d", insp.State.ExitCode))

	log.Debugf("finished task: %s", t.Metadata.Id)
	// TODO: enhance cleanup options
	if err = c.ContainerRemove(ctx, cont.ID, dtypes.ContainerRemoveOptions{}); err != nil {
		return err
	}

	return nil
}

func (c *dockerClientImpl) attachAndReportLogs(ctx context.Context, t *command.Task, tty bool) error {
	attch, err := c.ContainerAttach(ctx, t.Metadata.Id, dtypes.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return err
	}

	logsC := readLogs(ioutil.NopCloser(attch.Reader), tty, false)

	go func() {
		for log := range logsC {
			t.Logger.Log(log)
		}
	}()

	return nil
}

func (c *dockerClientImpl) CollectStats(ctx context.Context) (collector.RuntimeStats, error) {
	return collector.RuntimeStats{}, nil
}

func (c *dockerClientImpl) Running(ctx context.Context) (bool, error) {
	_, err := c.Ping(ctx)
	if err != nil && util.IsCONNREFUSED(err) {
		return false, nil
	}

	return err == nil, err
}

// Demultiplex mux by using docker stdcopy util. Returns
// stdout and stderr in that oreder.
func readLogs(r io.ReadCloser, tty, timestamps bool) <-chan *v1alpha1.LogEntry {
	logC := make(chan *v1alpha1.LogEntry)

	if tty {
		go readTTYLogs(r, logC, timestamps)
	} else {
		go demuxLogs(r, logC, timestamps)
	}

	return logC
}

func readTTYLogs(r io.ReadCloser, logC chan *v1alpha1.LogEntry, timestamps bool) {
	scn := bufio.NewScanner(r)
	defer r.Close()
	defer close(logC)

	if timestamps {
		for scn.Scan() {
			line, ts := getLogTimestamp(scn.Text())
			logC <- &v1alpha1.LogEntry{
				Data:      []byte(line),
				Stream:    v1alpha1.Stream_OUT,
				Timestamp: ts,
			}
		}
	} else {
		for scn.Scan() {
			logC <- &v1alpha1.LogEntry{
				Data:   scn.Bytes(),
				Stream: v1alpha1.Stream_OUT,
			}
		}
	}

	if scn.Err() != nil {
		log.WithError(scn.Err()).Error("failed to read container logs")
	}
}

func demuxLogs(mux io.ReadCloser, logC chan *v1alpha1.LogEntry, timestamps bool) {
	hdr := make([]byte, 8)
	defer mux.Close()
	defer close(logC)

	for {
		_, err := mux.Read(hdr)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.WithError(err).Error("failed to read container logs")
		}

		l := &v1alpha1.LogEntry{}

		if hdr[0] == 1 {
			l.Stream = v1alpha1.Stream_OUT
		} else {
			l.Stream = v1alpha1.Stream_ERR
		}

		size := binary.BigEndian.Uint32(hdr[4:])
		chunk := make([]byte, size)
		n, err := mux.Read(chunk)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.WithError(err).Error("failed to read container logs")
		}

		l.Data = chunk[:n]
		if timestamps {
			newLog, ts := getLogTimestamp(string(l.Data))
			l.Data = []byte(newLog)
			l.Timestamp = ts
		}

		logC <- l
	}
}

func getLogTimestamp(log string) (cleanLog string, stamp *prototypes.Timestamp) {
	ind := strings.Index(log, " ")
	if ind == -1 {
		return log, nil
	}

	subs := strings.SplitN(log, " ", 2)
	stampStr := subs[0]
	cleanLog = subs[1]

	ts, err := time.Parse(time.RFC3339Nano, stampStr)
	if err != nil {
		return
	}
	seconds := ts.Unix()
	nanos := int32(ts.Sub(time.Unix(seconds, 0)))

	stamp = &prototypes.Timestamp{
		Seconds: seconds,
		Nanos:   nanos,
	}

	return
}
