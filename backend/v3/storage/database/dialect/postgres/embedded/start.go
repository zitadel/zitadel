// embedded is used for testing purposes
package embedded

import (
	"net"
	"os"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
)

// StartEmbedded starts an embedded postgres v16 instance and returns a database connector and a stop function
// the database is started on a random port and data are stored in a temporary directory
// its used for testing purposes only
func StartEmbedded() (connector database.Connector, stop func(), err error) {
	path, err := os.MkdirTemp("", "zitadel-embedded-postgres-*")
	logging.OnError(err).Fatal("unable to create temp dir")

	port, close := getPort()

	config := embeddedpostgres.DefaultConfig().Version(embeddedpostgres.V16).Port(uint32(port)).RuntimePath(path)
	embedded := embeddedpostgres.NewDatabase(config)

	close()
	err = embedded.Start()
	logging.OnError(err).Fatal("unable to start db")

	connector, err = postgres.DecodeConfig(config.GetConnectionURL())
	if err != nil {
		return nil, nil, err
	}

	return connector, func() {
		logging.OnError(embedded.Stop()).Error("unable to stop db")
	}, nil
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
