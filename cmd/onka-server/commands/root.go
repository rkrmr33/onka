package commands

import (
	"context"

	"github.com/rkrmr33/onka/common"
	"github.com/rkrmr33/onka/server/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

	cmd.AddCommand(NewVersionCmd())

	return cmd
}

func RunRootCmd(ctx context.Context) error {
	log.Info("running server")
	return nil
}
