package main

import (
	"context"
	"syscall"

	"github.com/rkrmr33/onka/cmd/onka/commands"
	"github.com/rkrmr33/onka/pkg/util"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	ctx = util.ContextWithCancelOnSignals(ctx, syscall.SIGINT, syscall.SIGTERM)

	cmd := commands.NewRootCmd()

	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}
