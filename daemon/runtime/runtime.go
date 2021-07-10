package runtime

import (
	"context"
	"errors"
	"fmt"

	"github.com/rkrmr33/onka/daemon/collector"
	"github.com/rkrmr33/onka/daemon/command"
	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
	"github.com/spf13/cobra"
)

//go:generate mockery -name Runtime -filename runtime.go

// errors
var (
	ErrUnknownRuntimeType   = errors.New("unknown runtime type")
	ErrInvalidRuntimeConfig = errors.New("invalid runtime config")
)

type (
	Builder interface {
		Build(context.Context) (Runtime, error)
	}

	RuntimeFlagsAdder func(*cobra.Command) interface{}

	Runtime interface {
		// Execute a task on the runtime, blocks until the task is completed
		Execute(context.Context, *command.Task) error

		// Running checks whether the runtime is currently in "running" state
		Running(context.Context) (bool, error)

		// CollectStats collects runtime stats
		CollectStats(context.Context) (collector.RuntimeStats, error)

		// Start starts the runtime if it's not already running
		Start(context.Context) error

		// Stop stop the runtime, sends all currently running tasks a SIGTERM and waits for
		// them to finish before returning
		Stop(context.Context) error

		// Info the first return value should be the runtime name, the second is a pointer to
		// a struct containing extra informantion on this runtime. If getting the extra info
		// fails for some reason, return an error
		Info(ctx context.Context) (*v1alpha1.RuntimeInfo, error)
	}
)

var runtimesMap = map[string]Builder{}

func RegisterRuntime(name string, builder Builder) {
	runtimesMap[name] = builder
}

func ListRuntimes() []string {
	res := make([]string, 0, len(runtimesMap))
	for name := range runtimesMap {
		res = append(res, name)
	}
	return res
}

func NewRuntime(ctx context.Context, runtimeType string) (Runtime, error) {
	builder, exists := runtimesMap[runtimeType]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrUnknownRuntimeType, runtimeType)
	}

	return builder.Build(ctx)
}
