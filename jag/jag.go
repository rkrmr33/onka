package jag

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/rkrmr33/onka/jag/compiler"
	"github.com/rkrmr33/onka/jag/compiler/bexpr"
	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
	"github.com/rkrmr33/onka/pkg/util"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type ExecuteOptions struct {
	timeout time.Duration
}

type Executor interface {
	RunTask(context.Context, *v1alpha1.RunTaskRequest, ...grpc.CallOption) (v1alpha1.DaemonService_RunTaskClient, error)
}

type executionContext struct {
	ctx  context.Context
	errC chan error
	wg   sync.WaitGroup

	statusesLock sync.RWMutex
	statuses     map[string]*v1alpha1.TaskStatus

	executor Executor
}

func Execute(ctx context.Context, pipeline []byte, e Executor, opts *ExecuteOptions) error {
	p, entrypoints, err := compiler.Compile(pipeline) // TODO: bake with parameters, throw err if missing values
	if err != nil {
		return err
	}

	if opts != nil {
		if opts.timeout != 0 {
			var cancel func()
			ctx, cancel = context.WithTimeout(ctx, opts.timeout)
			defer cancel()
		}
	}

	exectx := &executionContext{
		ctx:          ctx,
		errC:         make(chan error, 10),
		wg:           sync.WaitGroup{},
		statusesLock: sync.RWMutex{},
		statuses:     make(map[string]*v1alpha1.TaskStatus, len(p.Spec.Tasks)),
		executor:     e,
	}

	exectx.initializeTaskStatuses(p.Spec.Tasks)

	for i := range entrypoints {
		exectx.wg.Add(1)
		go exectx.execStep(&entrypoints[i])
	}

	exectx.wg.Wait()

	log.Info("finished all steps")
	select {
	case err = <-exectx.errC:
	default:
	}

	return err
}

func (e *executionContext) execStep(t *v1alpha1.PipelineTask) {
	defer e.wg.Done()
	log.Infof("executing step: %s", t.Name)

	client, err := e.executor.RunTask(e.ctx, &v1alpha1.RunTaskRequest{
		Task: &v1alpha1.Task{
			Metadata: &v1alpha1.Metadata{
				Id: t.Name,
			},
			Spec: t.Task,
		},
		Watch: true,
	})
	if err != nil {
		e.reportTaskError(t, fmt.Errorf("failed to run task: %s: %w", t.Name, err))
		return
	}

	dmux := util.NewTaskWatchDemuxer(client)
	handlers := sync.WaitGroup{}

	handlers.Add(3)

	go func() {
		defer handlers.Done()
		e.handleLogs(t.Name, dmux.LogsC())
	}()

	go func() {
		defer handlers.Done()
		for err := range dmux.ErrC() {
			e.reportTaskError(t, fmt.Errorf("failed to receive task: %s: %w", t.Name, err))
		}
	}()

	go func() {
		defer handlers.Done()
		for status := range dmux.StatusC() {
			e.updateStepStatus(t, status)
		}
	}()

	handlers.Wait() // wait for step handlers to complete
}

func (e *executionContext) updateStepStatus(t *v1alpha1.PipelineTask, status *v1alpha1.TaskStatus) {
	e.statusesLock.Lock()
	defer e.statusesLock.Unlock()

	e.statuses[t.Name] = status

	fmt.Printf("[%s::%s]: %s\n", t.Name, status.State.String(), status.Cause)

	for _, dt := range t.Dependants {
		if s := e.statuses[dt.Name]; s.State != v1alpha1.TaskState_TASK_STATE_UNSPECIFIED {
			continue // this task is already being handled
		}

		ready, err := bexpr.EvalNeedsExpr(dt, e.statuses)
		if err != nil {
			e.errC <- err
			continue
		}
		if ready {
			e.wg.Add(1)
			e.statuses[dt.Name] = &v1alpha1.TaskStatus{State: v1alpha1.TaskState_TASK_STATE_PENDING}
			go e.execStep(dt)
		}
	}
}

func (e *executionContext) handleLogs(node string, logsC <-chan *v1alpha1.LogEntry) {
	for l := range logsC {
		var stream io.Writer
		if l.Stream == v1alpha1.Stream_ERR {
			stream = os.Stderr
		} else {
			stream = os.Stdout
		}
		fmt.Fprintf(stream, "[%s] => %s", node, string(l.Data))
	}
}

func (e *executionContext) reportTaskError(t *v1alpha1.PipelineTask, err error) {
	e.errC <- err
	e.updateStepStatus(t, &v1alpha1.TaskStatus{
		State: v1alpha1.TaskState_TASK_STATE_ERROR,
		Cause: err.Error(),
		From:  types.TimestampNow(),
	})
}

func (e *executionContext) initializeTaskStatuses(tasks map[string]*v1alpha1.TaskSpec) {
	for name := range tasks {
		e.statuses[name] = &v1alpha1.TaskStatus{State: v1alpha1.TaskState_TASK_STATE_UNSPECIFIED}
	}
}
