package postgres

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/zitadel/logging"
)

func StartEmbedded() (embeddedpostgres.Config, func()) {
	path, err := os.MkdirTemp("", "zitadel-embedded-postgres-*")
	logging.OnError(err).Fatal("unable to create temp dir")

	port := getPort()

	config := embeddedpostgres.DefaultConfig().Version(embeddedpostgres.V16).Port(uint32(port)).RuntimePath(path)
	embedded := embeddedpostgres.NewDatabase(config)

	err = embedded.Start()
	logging.OnError(err).Fatal("unable to start db")

	return config, func() {
		logging.OnError(embedded.Stop()).Error("unable to stop db")
	}
}

var (
	nextPort = uint16(5432)
	portMu   = sync.Mutex{}
)

func getPort() uint16 {
	portMu.Lock()
	defer portMu.Unlock()
	for {
		timeout := time.Second
		_, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", fmt.Sprintf("%d", nextPort)), timeout)
		if err != nil {
			return nextPort
		}
		nextPort++
	}
}
