package command

import (
	"context"

	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
)

// GetInfoCmd
type GetInfoCmd struct {
	BaseCommand

	result *v1alpha1.InfoResponse
	err    error
}

type GetInfoHandle struct {
	BaseHandle
	cmd *GetInfoCmd
}

func NewGetInfoCmd(ctx context.Context) (*GetInfoCmd, *GetInfoHandle) {
	h := &GetInfoHandle{
		BaseHandle: NewBaseHandle(),
	}

	cmd := &GetInfoCmd{BaseCommand: NewBaseCommand(ctx, &h.BaseHandle)}

	h.cmd = cmd

	return cmd, h
}

func (cmd *GetInfoCmd) SetResult(res *v1alpha1.InfoResponse, err error) {
	cmd.result = res
	cmd.err = err

	if err != nil {
		cmd.h.updateState(Failure)
	} else {
		cmd.h.updateState(Success)
	}
}

func (h *GetInfoHandle) AwaitResult(ctx context.Context) (*v1alpha1.InfoResponse, error) {
	if err := h.WaitFinished(ctx); err != nil {
		return nil, err
	}

	return h.cmd.result, h.cmd.err
}

// StartRuntimeCmd
type StartRuntimeCmd struct {
	BaseCommand

	err error
}

type StartRuntimeHandle struct {
	BaseHandle
	cmd *StartRuntimeCmd
}

func NewStartRuntimeCmd(ctx context.Context) (*StartRuntimeCmd, *StartRuntimeHandle) {
	h := &StartRuntimeHandle{
		BaseHandle: NewBaseHandle(),
	}

	cmd := &StartRuntimeCmd{BaseCommand: NewBaseCommand(ctx, &h.BaseHandle)}

	h.cmd = cmd

	return cmd, h
}

func (cmd *StartRuntimeCmd) SetResult(err error) {
	cmd.err = err

	if err != nil {
		cmd.h.updateState(Failure)
	} else {
		cmd.h.updateState(Success)
	}
}

func (h *StartRuntimeHandle) AwaitResult(ctx context.Context) error {
	if err := h.WaitFinished(ctx); err != nil {
		return err
	}

	return h.cmd.err
}

// StopRuntimeCmd
type StopRuntimeCmd struct {
	BaseCommand

	err error
}

type StopRuntimeHandle struct {
	BaseHandle
	cmd *StopRuntimeCmd
}

func NewStopRuntimeCmd(ctx context.Context) (*StopRuntimeCmd, *StopRuntimeHandle) {
	h := &StopRuntimeHandle{
		BaseHandle: NewBaseHandle(),
	}

	cmd := &StopRuntimeCmd{BaseCommand: NewBaseCommand(ctx, &h.BaseHandle)}

	h.cmd = cmd

	return cmd, h
}

func (cmd *StopRuntimeCmd) SetResult(err error) {
	cmd.err = err

	if err != nil {
		cmd.h.updateState(Failure)
	} else {
		cmd.h.updateState(Success)
	}
}

func (h *StopRuntimeHandle) AwaitResult(ctx context.Context) error {
	if err := h.WaitFinished(ctx); err != nil {
		return err
	}

	return h.cmd.err
}
