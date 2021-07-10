package common

import (
	"os"

	"github.com/docker/machine/libmachine/log"
	"github.com/sirupsen/logrus"
)

var (
	BinaryName = "dev"
	Version    = "v9.9.9"

	HomeDir string

	RuntimeConfigPath = "runtime.config"
	ServerConfigPath  = "server"
)

func init() {
	var err error
	HomeDir, err = os.UserHomeDir()
	if err != nil {
		logrus.WithError(err).Fatal("failed to get user homedir")
	}

	log.Debugf("user home-dir: %s", HomeDir)
}
