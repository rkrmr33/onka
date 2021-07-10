package load

import (
	log "github.com/sirupsen/logrus"

	"github.com/rkrmr33/onka/daemon/runtime"
	dockernative "github.com/rkrmr33/onka/daemon/runtime/docker/native"
)

func loadPlatformRuntimes() {
	log.Debug("Using linux, loading runtime...")

	runtime.RegisterRuntime(dockernative.Name, dockernative.AddFlags())
}
