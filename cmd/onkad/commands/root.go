package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/rkrmr33/onka/cmd/onkad/load"
	"github.com/rkrmr33/onka/common"
	"github.com/rkrmr33/onka/daemon"
	"github.com/rkrmr33/onka/daemon/collector"
	"github.com/rkrmr33/onka/daemon/command"
	"github.com/rkrmr33/onka/daemon/config"
	"github.com/rkrmr33/onka/daemon/runtime"
	"github.com/rkrmr33/onka/daemon/server"
)

var (
	conf config.Config
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   common.BinaryName,
		Short: "starts " + common.BinaryName,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := conf.Load(); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunRootCmd(cmd.Context())
		},
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       common.Version,
	}

	config.Init(cmd)

	conf.AddFlags()

	load.LoadRuntimes()

	cmd.AddCommand(NewRuntimesCmd())
	cmd.AddCommand(NewVersionCmd())

	return cmd
}

func RunRootCmd(ctx context.Context) error {
	r, err := runtime.NewRuntime(ctx, conf.Runtime.Type)
	if err != nil {
		return fmt.Errorf("failed to create runtime: %w", err)
	}

	c := collector.NewCollector()

	muxTaskSrc := command.MuxSrc()

	d, err := daemon.NewDaemon(conf.Daemon, r, c, muxTaskSrc)
	if err != nil {
		return err
	}

	srv := server.NewServer(&conf.Server, d)

	muxTaskSrc.Tap(srv)

	if err = srv.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	if err = d.Start(ctx); err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	<-ctx.Done() // wait for signal

	srv.Stop()

	if err = d.Stop(); err != nil {
		return fmt.Errorf("failed to stop daemon: %w", err)
	}

	return nil
}
