package server

import (
	"context"
	"fmt"
	"net"
	"time"

	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/rkrmr33/onka/daemon"
	"github.com/rkrmr33/onka/daemon/command"
	"github.com/rkrmr33/onka/pkg/proto/v1alpha1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	DefaultListenAddr     = ":6543"
	DefaultConnTimeout    = 180 * time.Second
	DefaultMaxSend        = 1 << 20 // 1mb
	DefaultMaxRecv        = 1 << 20 // 1mb
	DefaultTaskBufferSize = 100
)

type ServerConfig struct {
	ListenAddr        string
	ConnectionTimeout time.Duration
	MaxRecv           int
	MaxSend           int
}

type Server interface {
	command.CmdSrc

	// Start starts the server, this is non blocking
	Start() error

	// Stop starts server graceful shutdown and blocks until all live connections
	// are terminated
	Stop()
}

type server struct {
	addr   string
	grpcs  *grpc.Server
	daemon daemon.Daemon
	cmdC   chan command.Cmd
}

func NewServer(conf *ServerConfig, d daemon.Daemon) Server {
	s := &server{
		daemon: d,
		addr:   conf.ListenAddr,
		cmdC:   make(chan command.Cmd, DefaultTaskBufferSize),
	}

	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(conf.MaxRecv),
		grpc.MaxSendMsgSize(conf.MaxSend),
		grpc.ConnectionTimeout(conf.ConnectionTimeout),

		grpc.ChainStreamInterceptor(
			grpclogrus.StreamServerInterceptor(
				log.NewEntry(log.StandardLogger()),
				grpclogrus.WithDurationField(grpclogrus.DurationToDurationField),
			),
		),

		grpc.ChainUnaryInterceptor(
			grpclogrus.UnaryServerInterceptor(
				log.NewEntry(log.StandardLogger()),
				grpclogrus.WithDurationField(grpclogrus.DurationToDurationField),
			),
		),
	}

	s.grpcs = grpc.NewServer(opts...)

	v1alpha1.RegisterDaemonServiceServer(s.grpcs, s)
	reflection.Register(s.grpcs)

	return s
}

func (s *server) Start() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on: %s, %w", s.addr, err)
	}
	log.Infof("listening on: %s", s.addr)

	go func() {
		if err = s.grpcs.Serve(l); err != nil {
			log.Fatalf("grpc server error: %s", err)
		}
	}()

	return nil
}

func (s *server) Stop() {
	log.Debug("grpc server graceful shutdown initiated")
	s.grpcs.GracefulStop()
}

func (s *server) Recv(ctx context.Context) (command.Cmd, error) {
	select {
	case <-ctx.Done():
		return nil, command.ErrRecvCanceled
	case t, more := <-s.cmdC:
		if !more {
			return nil, command.ErrCmdSrcClosed
		}
		return t, nil
	}
}

func (s *server) Info(ctx context.Context, _ *v1alpha1.InfoRequest) (*v1alpha1.InfoResponse, error) {
	cmd, handle := command.NewGetInfoCmd(ctx)
	s.cmdC <- cmd
	return handle.AwaitResult(ctx)
}

func (s *server) RunTask(req *v1alpha1.RunTaskRequest, res v1alpha1.DaemonService_RunTaskServer) error {
	task, handle := command.NewTask(res.Context(), req.Task)
	task.WithLogger(command.NewCallbackLogger(forwardLogs(res)))
	handle.AddOnUpdateHandler(forwardStatus(res))
	task.UpdateState(v1alpha1.TaskState_TASK_STATE_PENDING, "task has been received")

	if !req.Watch {
		return res.Send(&v1alpha1.RunTaskResponse{
			Event: &v1alpha1.RunTaskResponse_StatusEvent{
				StatusEvent: handle.CurrentStatus(),
			},
		})
	}
	s.cmdC <- task

	return handle.WaitFinished(res.Context())
}

func (s *server) StopRuntime(ctx context.Context, req *v1alpha1.StopRuntimeRequest) (*v1alpha1.StopRuntimeResponse, error) {
	cmd, handle := command.NewStopRuntimeCmd(ctx)
	s.cmdC <- cmd
	return &v1alpha1.StopRuntimeResponse{}, handle.AwaitResult(ctx)
}

func (s *server) StartRuntime(ctx context.Context, req *v1alpha1.StartRuntimeRequest) (*v1alpha1.StartRuntimeResponse, error) {
	cmd, handle := command.NewStartRuntimeCmd(ctx)
	s.cmdC <- cmd
	return &v1alpha1.StartRuntimeResponse{}, handle.AwaitResult(ctx)
}

func forwardLogs(res v1alpha1.DaemonService_RunTaskServer) func(*v1alpha1.LogEntry) {
	return func(l *v1alpha1.LogEntry) {
		err := res.Send(&v1alpha1.RunTaskResponse{
			Event: &v1alpha1.RunTaskResponse_LogEvent{
				LogEvent: l,
			},
		})
		if err != nil {
			log.WithError(err).Error("failed to send task logs")
		}
	}
}

func forwardStatus(res v1alpha1.DaemonService_RunTaskServer) func(*v1alpha1.TaskStatus) {
	return func(s *v1alpha1.TaskStatus) {
		err := res.Send(&v1alpha1.RunTaskResponse{
			Event: &v1alpha1.RunTaskResponse_StatusEvent{
				StatusEvent: s,
			},
		})
		if err != nil {
			log.WithError(err).Error("failed to send task status")
		}
	}
}
