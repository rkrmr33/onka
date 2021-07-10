package command

import (
	"context"
	"io"
	"sync"

	"github.com/gogo/protobuf/types"
	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
	log "github.com/sirupsen/logrus"
)

type Logger interface {
	io.Writer
	Log(*v1alpha1.LogEntry)
}

type Task struct {
	BaseCommand
	Logger Logger

	*v1alpha1.Task
	ctx    context.Context
	handle *TaskHandle
}

type TaskHandle struct {
	BaseHandle

	taskMux *sync.RWMutex
	task    *Task

	handlersMux      *sync.Mutex
	onUpdateHandlers []func(*v1alpha1.TaskStatus)
}

func NewTask(ctx context.Context, task *v1alpha1.Task) (*Task, *TaskHandle) {
	t := &Task{
		BaseCommand: BaseCommand{ctx: ctx},
		Logger:      &basicLogger{},
		Task:        task,
		ctx:         ctx,
	}

	h := &TaskHandle{
		BaseHandle:       NewBaseHandle(),
		task:             t,
		taskMux:          &sync.RWMutex{},
		handlersMux:      &sync.Mutex{},
		onUpdateHandlers: []func(*v1alpha1.TaskStatus){},
	}

	t.handle = h
	t.Statuses = []*v1alpha1.TaskStatus{
		{
			State: v1alpha1.TaskState_TASK_STATE_UNSPECIFIED,
			Cause: "initial status",
			From:  types.TimestampNow(),
		},
	}

	return t, h
}

func (t *Task) WithLogger(l Logger) *Task {
	t.Logger = l
	return t
}

func ToCmdState(s v1alpha1.TaskState) State {
	switch s {
	case v1alpha1.TaskState_TASK_STATE_PENDING:
		return Pending
	case v1alpha1.TaskState_TASK_STATE_PREPARE:
		return Preparing
	case v1alpha1.TaskState_TASK_STATE_RUNNING:
		return Running
	case v1alpha1.TaskState_TASK_STATE_SUCCESS:
		return Success
	case v1alpha1.TaskState_TASK_STATE_FAILURE:
		return Failure
	case v1alpha1.TaskState_TASK_STATE_ERROR:
		return Error
	default:
		return Unkown
	}
}

//*****************************
//         Task Impl          *
//*****************************

func (j *Task) UpdateState(s v1alpha1.TaskState, cause string) {
	j.handle.addStatus(&v1alpha1.TaskStatus{
		State: s,
		Cause: cause,
		From:  types.TimestampNow(),
	})
}

//****************************
//        Handle Impl        *
//****************************
func (h *TaskHandle) CurrentStatus() *v1alpha1.TaskStatus {
	h.taskMux.RLock()
	defer h.taskMux.RUnlock()
	return h.task.Statuses[len(h.task.Statuses)-1]
}

func (h *TaskHandle) AddOnUpdateHandler(handle func(*v1alpha1.TaskStatus)) {
	h.handlersMux.Lock()
	h.onUpdateHandlers = append(h.onUpdateHandlers, handle)
	h.handlersMux.Unlock()
}

func (h *TaskHandle) addStatus(status *v1alpha1.TaskStatus) {
	h.taskMux.Lock()
	h.task.Statuses = append(h.task.Statuses, status)
	h.taskMux.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(len(h.onUpdateHandlers))

	h.handlersMux.Lock()
	for _, handler := range h.onUpdateHandlers {
		go func(handler func(*v1alpha1.TaskStatus)) {
			defer wg.Done()
			handler(status)
		}(handler)
	}
	h.handlersMux.Unlock()

	wg.Wait()

	h.updateState(ToCmdState(status.State))
}

type basicLogger struct{}

func (l *basicLogger) Log(e *v1alpha1.LogEntry) {
	log.Info(string(e.Data))
}

func (l *basicLogger) Write(data []byte) (int, error) {
	log.Info(string(data))
	return len(data), nil
}

type callbackLogger struct {
	cb func(*v1alpha1.LogEntry)
}

func NewCallbackLogger(cb func(*v1alpha1.LogEntry)) Logger {
	return callbackLogger{cb}
}

func (l callbackLogger) Log(e *v1alpha1.LogEntry) {
	l.cb(e)
}

func (l callbackLogger) Write(data []byte) (int, error) {
	l.cb(&v1alpha1.LogEntry{
		Data:      data,
		Stream:    v1alpha1.Stream_ERR,
		Timestamp: types.TimestampNow(),
	})
	return len(data), nil
}
