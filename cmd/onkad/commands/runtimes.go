package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/rkrmr33/onka/daemon/runtime"
)

func NewRuntimesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "runtimes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newRuntimesListCmd())

	return cmd
}

func newRuntimesListCmd() *cobra.Command {
	return &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			for _, n := range runtime.ListRuntimes() {
				fmt.Println(n)
			}
		},
	}
}
