package machine

import (
	"context"
	"fmt"
	"path/filepath"
	goruntime "runtime"

	"github.com/docker/machine/drivers/virtualbox"
	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/host"
	"github.com/docker/machine/libmachine/state"
	"github.com/rkrmr33/onka/common"
	"github.com/rkrmr33/onka/daemon/runtime"
	"github.com/rkrmr33/onka/daemon/runtime/docker"
	"github.com/rkrmr33/onka/pkg/config"
	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
	"github.com/rkrmr33/onka/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type DockerMachineConfig struct {
	MachineName string
	StorePath   string
	CertsDir    string
}

type dockerMachine struct {
	docker.DockerClient
	machine *host.Host
	info    *v1alpha1.DockerMachineInfo
	synced  util.AtomicBool
}

const (
	Name = "docker-machine"
)

var (
	DefaultMachineName      = "onka-runtime"
	DefaultMachinesPath     = filepath.Join(common.HomeDir, ".docker", "machine")
	DefaultMachineCertsPath = filepath.Join(DefaultMachinesPath, "certs")
)

func AddFlags() *DockerMachineConfig {
	c := &DockerMachineConfig{}

	flags := &pflag.FlagSet{}
	flags.StringVar(&c.MachineName, "docker-machine-name", DefaultMachineName, "machine name")
	flags.StringVar(&c.StorePath, "docker-machine-store-path", DefaultMachinesPath, "machines store path")
	flags.StringVar(&c.CertsDir, "docker-machine-certs-path", DefaultMachineCertsPath, "machines certs path")
	config.BindUnder(common.RuntimeConfigPath, flags)

	return c
}

func (cfg *DockerMachineConfig) Build(ctx context.Context) (runtime.Runtime, error) {
	api := libmachine.NewClient(cfg.StorePath, cfg.CertsDir)

	log.Info("using machine: ", cfg.MachineName)

	// We can load the machine only once becasue the dynamic config values (such as the machine's IP) are
	// stored under machine.Driver, which computes them when accessed.
	h, err := api.Load(cfg.MachineName)
	if err != nil {
		return nil, err
	}

	dm := &dockerMachine{
		DockerClient: nil,
		machine:      h,
		synced:       0,
	}

	return dm, nil
}

func (dm *dockerMachine) Running(ctx context.Context) (bool, error) {
	running, err := dm.machineRunning()
	if !running { // or err != nil
		return false, err
	}

	if !dm.synced.IsSet() {
		if err = dm.reloadDockerClient(ctx); err != nil {
			return false, err
		}
	}

	return dm.DockerClient.Running(ctx)
}

func (dm *dockerMachine) Start(ctx context.Context) error {
	log.Debug("Starting Docker Machine runtime")
	err := dm.machine.Start()
	if err != nil {
		return err
	}

	// In case machine IP changed
	return dm.reloadDockerClient(ctx)
}

func (dm *dockerMachine) Stop(ctx context.Context) error {
	log.Debug("Stopping Docker Machine runtime")
	dm.synced.Unset()
	return dm.machine.Stop()
}

func (d *dockerMachine) Info(ctx context.Context) (*v1alpha1.RuntimeInfo, error) {
	ret := &v1alpha1.RuntimeInfo{
		State:   v1alpha1.RuntimeState_RUNTIME_STATE_RUNNING,
		Cause:   "machine is running",
		Runtime: &v1alpha1.RuntimeInfo_DockerMachine{DockerMachine: d.info},
	}

	running, err := d.Running(ctx)
	if err != nil {
		ret.State = v1alpha1.RuntimeState_RUNTIME_STATE_ERROR
		ret.Cause = err.Error()
	}

	if !running {
		ret.State = v1alpha1.RuntimeState_RUNTIME_STATE_STOPPED
		ret.Cause = "machine is stopped"
	}

	return ret, nil
}

func (d *dockerMachine) getInfo(ctx context.Context) (*v1alpha1.DockerMachineInfo, error) {
	sv, err := d.DockerClient.ServerVersion(ctx)
	if err != nil {
		return nil, err
	}

	var (
		cpu    int32
		mem    int64
		driver interface{} = d.machine.Driver
	)

	if vbd, ok := driver.(*virtualbox.Driver); ok {
		cpu = int32(vbd.CPU)
		mem = int64(vbd.Memory)
	}

	r := &v1alpha1.DockerMachineInfo{
		Driver:        d.machine.DriverName,
		KernelVersion: sv.KernelVersion,
		Mem:           mem,
		Cpu:           cpu,
		EngineVersion: sv.Version,
		EngineOsArch:  fmt.Sprintf("%s/%s", sv.Os, sv.Arch),
		ClientVersion: d.DockerClient.ClientVersion(),
		ClientOsArch:  fmt.Sprintf("%s/%s", goruntime.GOOS, goruntime.GOARCH),
	}

	return r, nil
}

func (dm *dockerMachine) machineRunning() (bool, error) {
	stt, err := dm.machine.Driver.GetState()
	return stt == state.Running, err
}

func (dm *dockerMachine) reloadDockerClient(ctx context.Context) error {
	addr, err := dm.machine.Driver.GetURL()
	if err != nil {
		return err
	}

	auth := dm.machine.AuthOptions()
	if auth == nil {
		return fmt.Errorf("docker-machine nil auth options")
	}

	client, err := docker.OptsDockerClient(docker.DockerClientOpts{
		Host:       addr,
		CACertPath: auth.CaCertPath,
		CertPath:   auth.ClientCertPath,
		KeyPath:    auth.ClientKeyPath,
		Timeout:    docker.DefaultDockerTimeout,

		// TODO: Dynamic client version based on machine daemon version
		ClientVersion: "1.40",
	})
	if err != nil {
		return err
	}

	dm.DockerClient = client

	if info, err := dm.getInfo(ctx); err != nil {
		return err
	} else {
		dm.info = info
	}

	dm.synced.Set()
	return nil
}
