package util

import (
	"io"

	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
)

// TaskWatchDemuxer is used to deconstruct a DaemonService_RunTaskClient stream into
// a logs channel, status channel and error channel, to allow for easier watching of
// a task's progress.
type TaskWatchDemuxer struct {
	logsC   chan *v1alpha1.LogEntry
	statusC chan *v1alpha1.TaskStatus
	errC    chan error
}

func NewTaskWatchDemuxer(i v1alpha1.DaemonService_RunTaskClient) TaskWatchDemuxer {
	twd := TaskWatchDemuxer{
		logsC:   make(chan *v1alpha1.LogEntry),
		statusC: make(chan *v1alpha1.TaskStatus),
		errC:    make(chan error, 1),
	}

	go twd.handleEvents(i)

	return twd
}

func (twd TaskWatchDemuxer) StatusC() <-chan *v1alpha1.TaskStatus {
	return twd.statusC
}

func (twd TaskWatchDemuxer) LogsC() <-chan *v1alpha1.LogEntry {
	return twd.logsC
}

func (twd TaskWatchDemuxer) ErrC() <-chan error {
	return twd.errC
}

func (twd TaskWatchDemuxer) handleEvents(i v1alpha1.DaemonService_RunTaskClient) {
	defer close(twd.logsC)
	defer close(twd.statusC)
	defer close(twd.errC)

	for {
		ev, err := i.Recv()
		if err != nil {
			if err == io.EOF {
				return
			}
			twd.errC <- err
			return
		}

		switch ev := ev.Event.(type) {
		case *v1alpha1.RunTaskResponse_LogEvent:
			twd.logsC <- ev.LogEvent
		case *v1alpha1.RunTaskResponse_StatusEvent:
			twd.statusC <- ev.StatusEvent
		}
	}
}
