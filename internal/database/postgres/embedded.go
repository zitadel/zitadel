package postgres

import (
	"net"
	"os"
	"strconv"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/zitadel/logging"
)

func StartEmbedded() (embeddedpostgres.Config, func()) {
	tempPath, err := os.MkdirTemp("", "db")
	logging.OnError(err).Fatal("unable to create temp dir")

	port := uint16(5432)
	for isPortInUse(port) {
		logging.WithFields("port", port).Debug("port in use, trying next")
		port++
	}

	config := embeddedpostgres.DefaultConfig().Version(embeddedpostgres.V16).RuntimePath(tempPath).Port(uint32(port))
	psql := embeddedpostgres.NewDatabase(config)
	err = psql.Start()
	logging.OnError(err).Fatal("unable to start db")

	return config, func() {
		logging.OnError(psql.Stop()).Error("unable to stop db")
	}
}

func isPortInUse(port uint16) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", strconv.FormatUint(uint64(port), 10)), timeout)
	if err != nil {
		return false
	}
	logging.OnError(conn.Close()).Debug("unable to close connection")
	return true
}
