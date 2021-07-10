package commands

import (
	"fmt"

	"github.com/rkrmr33/onka/common"
	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func NewVersionCmd(client *client) *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		PreRunE: requiresAuth(client),
		RunE: func(cmd *cobra.Command, _ []string) error {
			res, err := client.daemon.Info(cmd.Context(), &v1alpha1.InfoRequest{})
			if err != nil {
				return err
			}

			remote, err := yaml.Marshal(res)
			if err != nil {
				return err
			}

			fmt.Println("onka client:")
			fmt.Printf("Version: %s\n", common.Version)
			fmt.Println("")
			fmt.Println("onkad:")
			fmt.Printf("%v", string(remote))

			return nil
		},
	}
}
