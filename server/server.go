package server

const (
	DefaultListenAddr = ":6554"
)

type (
	ServerConfig struct {
		ListenAddr string
	}
)
