package postgres

import (
	"os"
	"sync"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/zitadel/logging"
)

var (
	port   uint16 = 5432
	portMu        = sync.Mutex{}
)

func StartEmbedded() (embeddedpostgres.Config, func()) {
	tempPath, err := os.MkdirTemp("", "db")
	logging.OnError(err).Fatal("unable to create temp dir")

	portMu.Lock()
	logging.WithFields("port", port).Debug("starting embedded postgres")
	config := embeddedpostgres.DefaultConfig().Version(embeddedpostgres.V16).RuntimePath(tempPath).Port(uint32(port))
	port++
	portMu.Unlock()

	psql := embeddedpostgres.NewDatabase(config)
	err = psql.Start()
	logging.OnError(err).Fatal("unable to start db")

	return config, func() {
		logging.OnError(psql.Stop()).Error("unable to stop db")
	}
}
