package daemon

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/rkrmr33/onka/common"
	"github.com/rkrmr33/onka/daemon/collector"
	"github.com/rkrmr33/onka/daemon/command"
	"github.com/rkrmr33/onka/daemon/runtime"
	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
	"github.com/rkrmr33/onka/pkg/util"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultUpdateTickRate = time.Second * 10
)

// Errors
var (
	ErrRuntimeNotRunning = errors.New("cannot execute tasks when runtime is shut down")
)

type (
	DaemonConfig struct {
		Thresholds collector.Thresholds
		UpdateTick time.Duration
	}

	Daemon interface {
		// Start starts the daemon.
		Start(ctx context.Context) error

		// Wait waits until the daemon has finished, you must cancel the context given to Start()
		// before calling this method.
		Stop() error

		// Info returns information about the daemon and the runtime
		Info(ctx context.Context) (*v1alpha1.InfoResponse, error)

		// StartRuntime starts the runtime
		StartRuntime(ctx context.Context) error

		// StopRuntime stops the runtime
		StopRuntime(ctx context.Context) error
	}

	daemon struct {
		mux            *sync.Mutex
		conf           DaemonConfig
		collector      collector.Collector
		r              runtime.Runtime
		cmdsrc         command.CmdSrc
		wg             *sync.WaitGroup
		runtimeStopped util.AtomicBool

		cmdC chan command.Cmd
		errC chan error
		ctx  context.Context
	}
)

func NewDaemon(conf DaemonConfig, r runtime.Runtime, collector collector.Collector, tasksrc command.CmdSrc) (Daemon, error) {
	return &daemon{
		r:         r,
		conf:      conf,
		collector: collector,
		cmdsrc:    tasksrc,
		wg:        &sync.WaitGroup{},
		mux:       &sync.Mutex{},

		cmdC: make(chan command.Cmd),
		errC: make(chan error),
	}, nil
}

func (d *daemon) Start(ctx context.Context) error {
	d.ctx = ctx

	d.runAsync(d.updateLoop)
	d.runAsync(d.executionLoop)
	d.runAsync(d.errorHandler)

	return d.updateStatus(true)
}

func (d *daemon) Stop() error {
	d.wg.Wait()
	return nil
}

func (d *daemon) Info(ctx context.Context) (*v1alpha1.InfoResponse, error) {
	rInfo, err := d.r.Info(ctx)
	if err != nil {
		return nil, err
	}

	return &v1alpha1.InfoResponse{
		DaemonVersion: common.Version,
		Runtime:       rInfo,
	}, nil
}

func (d *daemon) StartRuntime(ctx context.Context) error {
	d.runtimeStopped.Unset()
	return d.r.Start(ctx)
}

func (d *daemon) StopRuntime(ctx context.Context) error {
	d.runtimeStopped.Set()
	return d.r.Stop(ctx)
}

func (d *daemon) runAsync(asyncFn func()) {
	d.wg.Add(1)
	go func() {
		asyncFn()
		d.wg.Done()
	}()
}

func (d *daemon) updateLoop() {
	t := time.NewTicker(d.conf.UpdateTick)

	for {
		select {
		case <-d.ctx.Done():
			return
		case <-t.C:
			d.errC <- d.updateStatus(false)
		}
	}
}

func (d *daemon) executionLoop() {
	for {
		t, err := d.cmdsrc.Recv(d.ctx)
		if err != nil {
			switch err {
			case command.ErrCmdSrcClosed, command.ErrRecvCanceled:
				log.Debugf("task source closed, stopping execution loop: %s", err.Error())
				return
			}

			d.errC <- err
			continue
		}

		go d.handleCmd(t)
	}
}

func (d *daemon) errorHandler() {
	for {
		select {
		case <-d.ctx.Done():
			return
		case err := <-d.errC:
			if err == nil {
				continue
			}
			log.Errorf("runtime error: %s", err)
		}
	}
}

func (d *daemon) updateStatus(print bool) error {
	running, err := d.r.Running(d.ctx)
	if err != nil {
		return err
	}

	s, err := d.collector.Collect(d.ctx)
	if err != nil {
		return err
	}

	upholds := d.conf.Thresholds.Upholds(s)

	if running {
		if upholds {
			if print {
				log.Debug("runtime is up and upholds set threasholds")
			}
			// runtime is up and upholds threasholds
			// report everyting is OK
			return nil
		}

		log.Warn("runtime does not uphold set threasholds, shuting down...")
		// report threasholds issue

		return d.r.Stop(d.ctx)
	}

	if upholds && !d.runtimeStopped.IsSet() {
		// runtime was NOT intentionally stopped and needs to be started
		log.Debug("detected runtime is down, trying to start runtime...")
		if err = d.r.Start(d.ctx); err != nil {
			// report cannot start runtime
			return fmt.Errorf("failed to start runtime: %w", err)
		}

		// report runtime started everything is OK
		return nil
	}

	// report runtime is supposed to be down
	return nil
}

func (d *daemon) handleCmd(cmd command.Cmd) {
	switch cmd := cmd.(type) {
	case *command.Task:
		go d.handleTaskCmd(cmd)
	case *command.GetInfoCmd:
		go d.handleGetInfoCmd(cmd)
	case *command.StartRuntimeCmd:
		go d.handleStartRuntimeCmd(cmd)
	case *command.StopRuntimeCmd:
		go d.handleStopRuntimeCmd(cmd)
	default:
		d.errC <- fmt.Errorf("unhandled command type: %s", reflect.TypeOf(cmd).Elem().Name())
	}
}

func (d *daemon) handleTaskCmd(t *command.Task) {
	if running, err := d.r.Running(d.ctx); err != nil {
		d.errC <- err
		return
	} else if !running {
		return
	}

	d.errC <- d.r.Execute(d.ctx, t)
}

func (d *daemon) handleGetInfoCmd(cmd *command.GetInfoCmd) {
	info, err := d.Info(cmd.Context())
	cmd.SetResult(info, err)
}

func (d *daemon) handleStartRuntimeCmd(cmd *command.StartRuntimeCmd) {
	err := d.StartRuntime(cmd.Context())
	cmd.SetResult(err)
}

func (d *daemon) handleStopRuntimeCmd(cmd *command.StopRuntimeCmd) {
	err := d.StopRuntime(cmd.Context())
	cmd.SetResult(err)
}
