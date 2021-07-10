package config

import (
	"path/filepath"

	"github.com/rkrmr33/onka/common"
	"github.com/rkrmr33/onka/daemon"
	"github.com/rkrmr33/onka/daemon/server"
	"github.com/rkrmr33/onka/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	defaultConfigName  = "config"
	defaultRuntimeType = "docker-machine"
	defaultConfigPath  = filepath.Join(common.HomeDir, ".onkaci", "onkad")
)

func Init(cmd *cobra.Command) {
	config.Init(cmd, defaultConfigPath, defaultConfigName)
}

type Config struct {
	Daemon  daemon.DaemonConfig
	Server  server.ServerConfig
	Runtime struct {
		Type   string
		Config interface{} `mapstructure:"config"`
	}
}

func (c *Config) AddFlags() {
	rootFlags := &pflag.FlagSet{}
	rootFlags.StringVar(&c.Runtime.Type, "runtime", defaultRuntimeType, "runtime type")
	config.BindConfig("runtime.type", rootFlags.Lookup("runtime"))

	daemonFlags := &pflag.FlagSet{}
	daemonFlags.DurationVar(&c.Daemon.UpdateTick, "update-rate", daemon.DefaultUpdateTickRate, "the rate in which to send status updates")
	daemonFlags.Uint64Var(&c.Daemon.Thresholds.RAM, "max-ram", ^uint64(0), "RAM usage threshold to trigger a runtime shutdown")
	daemonFlags.Float64Var(&c.Daemon.Thresholds.Load, "max-cpu", 1.0, "CPU 5m load threshold to trigger a runtime shutdown")
	config.BindUnder("daemon", daemonFlags)

	serverFlags := &pflag.FlagSet{}
	serverFlags.StringVarP(&c.Server.ListenAddr, "listen", "l", server.DefaultListenAddr, "grpc server listen address")
	serverFlags.DurationVar(&c.Server.ConnectionTimeout, "connection-timeout", server.DefaultConnTimeout, "grpc server connection timeout")
	serverFlags.IntVar(&c.Server.MaxSend, "max-send", server.DefaultMaxSend, "grpc server max send size")
	serverFlags.IntVar(&c.Server.MaxRecv, "max-receive", server.DefaultMaxRecv, "grpc server max receive size")
	config.BindUnder("server", serverFlags)
}

func (c *Config) Load() error {
	return config.Load()
}
