//go:build !integration

package start

import (
	"net"
)

func listenConfig() *net.ListenConfig {
	return &net.ListenConfig{}
}
