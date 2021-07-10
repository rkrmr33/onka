package commands

import (
	"fmt"

	"github.com/rkrmr33/onka/common"
	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use: "version",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Version: %s\n", common.Version)
		},
	}
}
