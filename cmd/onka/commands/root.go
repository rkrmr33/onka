package commands

import (
	"fmt"

	"github.com/rkrmr33/onka/common"
	onkaclient "github.com/rkrmr33/onka/pkg/client"
	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type client struct {
	daemon v1alpha1.DaemonServiceClient
}

func NewRootCmd() *cobra.Command {
	var (
		logLvl     string
		client     client
		clientOpts onkaclient.ClientOptions
	)

	cmd := &cobra.Command{
		Use:   common.BinaryName,
		Short: common.BinaryName + " is the rkrmr33 cli tool",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			l, err := log.ParseLevel(logLvl)
			if err != nil {
				return err
			}

			log.SetLevel(l)

			client.daemon, err = clientOpts.New(cmd.Context())
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       common.Version,
	}

	cmd.PersistentFlags().StringVar(&logLvl, "log-level", log.GetLevel().String(), "log level")

	clientOpts = onkaclient.AddFlags(cmd)

	cmd.AddCommand(NewVersionCmd(&client))
	cmd.AddCommand(NewDaemonCmd(&client))

	return cmd
}

func requiresAuth(c *client) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		if _, err := c.daemon.Info(cmd.Context(), &v1alpha1.InfoRequest{}); err != nil {
			return fmt.Errorf("connection to daemon failed: %w", err)
		}
		return nil
	}
}
