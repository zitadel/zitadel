package http

import (
	"net"
	"strings"

	"github.com/caos/logging"
)

func CreateListener(endpoint string) net.Listener {
	l, err := net.Listen("tcp", listenerEndpoint(endpoint))
	logging.Log("SERVE-6vasef").OnError(err).Fatal("creating listener failed")
	return l
}

func listenerEndpoint(endpoint string) string {
	if strings.Contains(endpoint, ":") {
		return endpoint
	}
	return ":" + endpoint
}
