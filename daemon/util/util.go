package util

import (
	"fmt"
	"net"
	"strings"
	"syscall"

	"github.com/fatih/color"
)

var executionMsg = color.New(color.FgGreen, color.Bold).Sprint("executing command: ")

func IsCONNREFUSED(err error) bool {
	switch t := err.(type) {
	case *net.OpError:
		if t.Op == "read" {
			return true
		}

	case syscall.Errno:
		if t == syscall.ECONNREFUSED {
			return true
		}
	}
	return false
}

func GetShellCmd(cmds []string) string {
	commands := make([]string, 0, len(cmds)*2)

	for _, cmd := range cmds {
		commands = append(commands, fmt.Sprintf("echo \"%s%s\"", executionMsg, cmd))
		commands = append(commands, cmd)
	}

	return strings.Join(commands, " && ")
}
