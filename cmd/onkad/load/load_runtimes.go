package load

import (
	"github.com/rkrmr33/onka/daemon/runtime"
	dockermachine "github.com/rkrmr33/onka/daemon/runtime/docker/machine"
)

func LoadRuntimes() {
	loadPlatformRuntimes()

	runtime.RegisterRuntime(dockermachine.Name, dockermachine.AddFlags())
}
