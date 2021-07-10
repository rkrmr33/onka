package commands

import (
	"io/ioutil"

	"github.com/rkrmr33/onka/jag"
	"github.com/spf13/cobra"
)

func NewDaemonCmd(client *client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "interact with a local onka daemon",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(NewDaemonRunCmd(client))

	return cmd
}

func NewDaemonRunCmd(client *client) *cobra.Command {
	var (
		file string
	)

	cmd := &cobra.Command{
		Use:     "run",
		Short:   "execute a pipeline on a local daemon",
		PreRunE: requiresAuth(client),
		RunE: func(cmd *cobra.Command, _ []string) error {
			p, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}

			return jag.Execute(cmd.Context(), p, client.daemon, &jag.ExecuteOptions{})
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "file that contains the pipeline definition that needs to be executed")

	return cmd
}
