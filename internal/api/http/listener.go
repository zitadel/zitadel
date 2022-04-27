package http

import (
	"net"
	"strings"

	"github.com/zitadel/logging"
)

func CreateListener(endpoint string) net.Listener {
	l, err := net.Listen("tcp", Endpoint(endpoint))
	logging.Log("SERVE-6vasef").OnError(err).Fatal("creating listener failed")
	return l
}

func Endpoint(endpoint string) string {
	if strings.Contains(endpoint, ":") {
		return endpoint
	}
	return ":" + endpoint
}
