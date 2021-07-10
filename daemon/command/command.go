package command

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

// Errors
var (
	ErrRecvCanceled = errors.New("receive command canceled")
	ErrCmdSrcClosed = errors.New("command source closed")
)

type Cmd interface {
	Context() context.Context
}

type State int32

// Command states
const (
	Unkown State = iota
	Pending
	Preparing
	Running
	Success
	Failure
	Error
)

type CmdSrc interface {
	Recv(context.Context) (Cmd, error)
}

type cmdErr struct {
	cmd Cmd
	err error
}

type MuxCmdSrc struct {
	cmdC   chan cmdErr
	cancel func()
	ctx    context.Context
	srcsN  int32
}

func (s State) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Preparing:
		return "Preparing"
	case Running:
		return "Running"
	case Success:
		return "Success"
	case Failure:
		return "Failure"
	case Error:
		return "Error"
	case Unkown:
		return "Unkown"
	default:
		return "Unkown"
	}
}

// Empty returns a new empty command source
func MuxSrc(srcs ...CmdSrc) *MuxCmdSrc {
	cmdC := make(chan cmdErr)
	ctx, cancel := context.WithCancel(context.Background())
	t := &MuxCmdSrc{cmdC, cancel, ctx, int32(len(srcs))}

	for _, src := range srcs {
		t.Tap(src)
	}

	return t
}

func (t *MuxCmdSrc) Tap(src CmdSrc) *MuxCmdSrc {
	go func() {
		for {
			cmd, err := src.Recv(t.ctx)
			t.cmdC <- cmdErr{cmd: cmd, err: err}
		}
	}()

	return t
}

func (t *MuxCmdSrc) Recv(ctx context.Context) (Cmd, error) {
	select {
	case terr := <-t.cmdC:
		if terr.err == ErrCmdSrcClosed && atomic.AddInt32(&t.srcsN, -1) == 0 {
			return nil, ErrCmdSrcClosed // only once all cmd sources are closed
		}
		return terr.cmd, terr.err
	case <-ctx.Done():
		t.cancel() // cancel all srcs recv
		return nil, ErrRecvCanceled
	}
}

type BaseCommand struct {
	ctx context.Context
	h   *BaseHandle
}

func NewBaseCommand(ctx context.Context, h *BaseHandle) BaseCommand {
	return BaseCommand{ctx, h}
}

func (bc *BaseCommand) Context() context.Context {
	return bc.ctx
}

type BaseHandle struct {
	stateC chan State
	doneC  chan struct{}

	curStateLock *sync.Mutex
	curState     State

	handlersLock     *sync.Mutex
	onUpdateHandlers []func(State)
}

func NewBaseHandle() BaseHandle {
	return BaseHandle{
		stateC: make(chan State, 64),
		doneC:  make(chan struct{}),

		curStateLock: &sync.Mutex{},
		curState:     Unkown,

		handlersLock:     &sync.Mutex{},
		onUpdateHandlers: []func(State){},
	}
}

func (h *BaseHandle) WaitFinished(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-h.doneC:
		return nil
	}
}

func (h *BaseHandle) CurrentState() State {
	return h.curState
}

func (h *BaseHandle) AddOnUpdateHandler(handle func(State)) {
	h.handlersLock.Lock()
	h.onUpdateHandlers = append(h.onUpdateHandlers, handle)
	h.handlersLock.Unlock()
}

func (h *BaseHandle) updateState(state State) {
	h.curStateLock.Lock()
	defer h.curStateLock.Unlock()
	if isFinishedState(h.CurrentState()) {
		panic(fmt.Errorf("attempted state updated after final state: %s", state))
	}

	h.curState = state

	wg := sync.WaitGroup{}
	wg.Add(len(h.onUpdateHandlers))

	h.handlersLock.Lock()
	for _, handler := range h.onUpdateHandlers {
		go func(handler func(State)) {
			defer wg.Done()
			handler(state)
		}(handler)
	}
	h.handlersLock.Unlock()

	if isFinishedState(state) {
		wg.Wait()
		close(h.doneC)
	}
}

func isFinishedState(state State) bool {
	switch state {
	case Error, Success, Failure:
		return true
	}
	return false
}
