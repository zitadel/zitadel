package postgres

import (
	"net"
	"os"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/zitadel/logging"
)

func StartEmbedded() (embeddedpostgres.Config, func()) {
	path, err := os.MkdirTemp("", "zitadel-embedded-postgres-*")
	logging.OnError(err).Fatal("unable to create temp dir")

	port, close := getPort()

	config := embeddedpostgres.DefaultConfig().Version(embeddedpostgres.V16).Port(uint32(port)).RuntimePath(path)
	embedded := embeddedpostgres.NewDatabase(config)

	close()
	err = embedded.Start()
	logging.OnError(err).Fatal("unable to start db")

	return config, func() {
		logging.OnError(embedded.Stop()).Error("unable to stop db")
	}
}

// getPort returns a free port and locks it until close is called
func getPort() (port uint16, close func()) {
	l, err := net.Listen("tcp", ":0")
	logging.OnError(err).Fatal("unable to get port")
	port = uint16(l.Addr().(*net.TCPAddr).Port)
	logging.WithFields("port", port).Info("Port is available")
	return port, func() {
		logging.OnError(l.Close()).Error("unable to close port listener")
	}
}
