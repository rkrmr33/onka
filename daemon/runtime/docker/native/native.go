// +build linux

package native

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	goruntime "runtime"
	"time"

	dtypes "github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"

	"github.com/rkrmr33/onka/daemon/collector"
	"github.com/rkrmr33/onka/daemon/runtime"
	"github.com/rkrmr33/onka/daemon/runtime/docker"
	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
)

const (
	systemdPath       = "/etc/systemd/system"
	systemdDockerDir  = systemdPath + "/docker.service.d"
	systemdDockerConf = systemdDockerDir + "/onkad-docker.conf"

	daemonConfDir  = "/etc/docker"
	daemonConfPath = daemonConfDir + "/daemon.json"

	Name = "docker-native"
)

type EngineOpts struct {
	Hosts []string `json:"hosts,omitempty"`
}

type dockerNative struct {
	docker.DockerClient
	cfg  *DockerNativeOpts
	info *v1alpha1.DockerNativeInfo
}

type DockerNativeOpts docker.DockerClientOpts

func AddFlags() *DockerNativeOpts {
	c := &DockerNativeOpts{}

	return c
}

func (cfg *DockerNativeOpts) Build(ctx context.Context) (runtime.Runtime, error) {
	d := &dockerNative{cfg: cfg}
	if cfg != nil {
		dc, err := docker.OptsDockerClient(docker.DockerClientOpts(*cfg))
		if err != nil {
			return nil, fmt.Errorf("failed to create docker client: %w", err)
		}
		d.DockerClient = dc
	} else {
		dc, err := docker.EnvDockerClient()
		if err != nil {
			return nil, fmt.Errorf("failed to create docker client: %w", err)
		}
		d.DockerClient = dc
	}

	return d, nil
}

func (d *dockerNative) CollectStats(ctx context.Context) (collector.RuntimeStats, error) {
	return collector.RuntimeStats{}, nil
}

func (d *dockerNative) Start(ctx context.Context) error {
	log.Debug("Starting Docker runtime")

	// Make cross-platform
	var opts *EngineOpts
	if d.cfg != nil {
		opts = &EngineOpts{Hosts: []string{d.cfg.Host}}
	}

	if err := startEngine(ctx, opts); err != nil {
		return err
	}

	info, err := d.getInfo(ctx)
	if err != nil {
		return err
	}
	d.info = info

	return nil
}

func (d *dockerNative) Stop(ctx context.Context) error {
	log.Debug("Stopping Docker runtime")
	// Make cross-platform
	return stopEngine(ctx)
}

func (d *dockerNative) Info(ctx context.Context) (*v1alpha1.RuntimeInfo, error) {
	ret := &v1alpha1.RuntimeInfo{
		State:   v1alpha1.RuntimeState_RUNTIME_STATE_RUNNING,
		Cause:   "daemon is running",
		Runtime: &v1alpha1.RuntimeInfo_DockerNative{DockerNative: d.info},
	}

	running, err := d.Running(ctx)
	if err != nil {
		ret.State = v1alpha1.RuntimeState_RUNTIME_STATE_ERROR
		ret.Cause = err.Error()
	}

	if !running {
		ret.State = v1alpha1.RuntimeState_RUNTIME_STATE_STOPPED
		ret.Cause = "daemon is stopped"
	}

	return ret, nil
}

func (d *dockerNative) getInfo(ctx context.Context) (*v1alpha1.DockerNativeInfo, error) {
	var (
		sv  dtypes.Version
		err error
	)

	for {
		sv, err = d.DockerClient.ServerVersion(ctx)
		if err == nil {
			break
		}
		<-time.After(time.Second * 5)
		log.WithError(err).Debug("waiting for daemon to start")
	}

	r := &v1alpha1.DockerNativeInfo{
		KernelVersion: sv.KernelVersion,
		EngineVersion: sv.Version,
		EngineOsArch:  fmt.Sprintf("%s/%s", sv.Os, sv.Arch),
		ClientVersion: d.DockerClient.ClientVersion(),
		ClientOsArch:  fmt.Sprintf("%s/%s", goruntime.GOOS, goruntime.GOARCH),
	}

	return r, nil
}

var startEngine = func(ctx context.Context, e *EngineOpts) error {
	if e != nil {
		err := setOpts(*e)
		if err != nil {
			return err
		}
	} else {
		err := nullifyOpts()
		if err != nil {
			return err
		}
	}

	return exec.CommandContext(ctx, "service", "docker", "start").Run()
}

var stopEngine = func(ctx context.Context) error {
	return exec.CommandContext(ctx, "service", "docker", "stop").Run()
}

func setOpts(e EngineOpts) error {
	systemd, err := hasSystemd()
	if err != nil {
		return err
	}

	if systemd {
		return setupSystemdConf(e)
	} else {
		return setupDaemonConf(e)
	}
}

func setupSystemdConf(e EngineOpts) error {
	err := os.MkdirAll(systemdDockerDir, 0x775)
	if err != nil && !os.IsExist(err) {
		return err
	}

	f, err := os.Create(systemdDockerConf)

	if len(e.Hosts) != 0 {
		if err != nil {
			return err
		}

		// TODO: Bash escaping
		_, err = f.Write([]byte("[Service]\nExecStart=\nExecStart=/usr/bin/dockerd -H '" + e.Hosts[0] + "' --containerd=/run/containerd/containerd.sock\n"))
		if err != nil {
			return err
		}
		err = exec.Command("systemctl", "daemon-reload").Run()
		if err != nil {
			return err
		}

	}
	return nil
}

func setupDaemonConf(e EngineOpts) error {
	err := os.MkdirAll(daemonConfDir, 0x775)
	if err != nil && !os.IsExist(err) {
		return err
	}

	f, err := os.Create(daemonConfPath)
	if err != nil {
		return err
	}

	return json.NewEncoder(f).Encode(e)
}

func nullifyOpts() error {
	systemd, err := hasSystemd()
	if err != nil {
		return err
	}

	if systemd {
		_, err := os.Stat(systemdDockerConf)
		if err == nil || os.IsNotExist(err) {
			return os.Remove(systemdDockerConf)
		}
		return err
	} else {
		_, err := os.Stat(daemonConfPath)
		if err == nil || os.IsNotExist(err) {
			return os.Remove(daemonConfPath)
		}
		return err
	}
}

func hasSystemd() (bool, error) {
	_, err := os.Stat(systemdPath)
	isNotExists := os.IsNotExist(err)
	if err != nil && !isNotExists {
		return false, err
	}

	return !isNotExists, nil
}
