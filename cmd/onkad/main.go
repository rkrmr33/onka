package main

import (
	"context"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/rkrmr33/onka/cmd/onkad/commands"
	"github.com/rkrmr33/onka/pkg/util"
)

func main() {
	ctx := context.Background()
	ctx = util.ContextWithCancelOnSignals(ctx, syscall.SIGINT, syscall.SIGTERM)

	cmd := commands.NewRootCmd()

	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}
