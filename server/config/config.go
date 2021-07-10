package config

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/rkrmr33/onka/common"
	"github.com/rkrmr33/onka/pkg/config"
	"github.com/rkrmr33/onka/server"
)

var (
	configPath string
	rootCmd    *cobra.Command

	defaultConfigName = "config"
	defaultConfigPath = filepath.Join(common.HomeDir, ".onkaci", "server")

	mappings = map[string]string{}
)

func Init(cmd *cobra.Command) {
	config.Init(cmd, defaultConfigPath, defaultConfigName)
}

type Config struct {
	Server server.ServerConfig
}

func (c *Config) AddFlags() {
	serverFlags := &pflag.FlagSet{}
	serverFlags.StringVarP(&c.Server.ListenAddr, "listen", "l", server.DefaultListenAddr, "grpc server listen address")

	config.BindUnder("server", serverFlags)
}

func (c *Config) Load() error {
	return config.Load()
}
