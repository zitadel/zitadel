package postgres

import (
	"sync/atomic"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/zitadel/logging"
)

var (
	runningTests atomic.Int32
	embedded     *embeddedpostgres.EmbeddedPostgres
)

func StartEmbedded() (embeddedpostgres.Config, func()) {
	runningCount := runningTests.Add(1)
	config := embeddedpostgres.DefaultConfig().Version(embeddedpostgres.V16)

	// postgres is already started if runningCount > 1
	if runningCount > 1 {
		return config, cleanup
	}

	embedded = embeddedpostgres.NewDatabase(config)
	err := embedded.Start()
	logging.OnError(err).Fatal("unable to start db")

	return config, cleanup
}

func cleanup() {
	if runningTests.Add(-1) > 0 {
		return
	}
	logging.OnError(embedded.Stop()).Error("unable to stop db")
}
