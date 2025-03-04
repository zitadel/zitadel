package postgres

import (
	"os"
	"sync"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/zitadel/logging"
)

var (
	port   = uint16(5432)
	portMu = sync.Mutex{}
)

func StartEmbedded() (embeddedpostgres.Config, func()) {
	path, err := os.MkdirTemp("", "zitadel-embedded-postgres-*")
	logging.OnError(err).Fatal("unable to create temp dir")

	portMu.Lock()
	startPort := port
	port++
	portMu.Unlock()

	config := embeddedpostgres.DefaultConfig().Version(embeddedpostgres.V16).Port(uint32(startPort)).RuntimePath(path)
	embedded := embeddedpostgres.NewDatabase(config)

	err = embedded.Start()
	logging.OnError(err).Fatal("unable to start db")

	return config, func() {
		logging.OnError(embedded.Stop()).Error("unable to stop db")
	}
}
