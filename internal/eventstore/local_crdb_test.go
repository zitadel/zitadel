package eventstore_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/initialise"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/cockroach"
)

var (
	testCRDBClient *sql.DB
)

func TestMain(m *testing.M) {
	ts, err := testserver.NewTestServer()
	if err != nil {
		logging.WithFields("error", err).Fatal("unable to start db")
	}

	testCRDBClient, err = sql.Open("postgres", ts.PGURL().String())
	if err != nil {
		logging.WithFields("error", err).Fatal("unable to connect to db")
	}
	if err != nil {
		logging.WithFields("error", err).Fatal("unable to connect to db")
	}
	if err = testCRDBClient.Ping(); err != nil {
		logging.WithFields("error", err).Fatal("unable to ping db")
	}

	defer func() {
		testCRDBClient.Close()
		ts.Stop()
	}()

	if err = initDB(testCRDBClient); err != nil {
		logging.WithFields("error", err).Fatal("migrations failed")
	}

	os.Exit(m.Run())
}

func initDB(db *sql.DB) error {
	config := new(database.Config)
	config.SetConnector(&cockroach.Config{
		User: cockroach.User{
			Username: "zitadel",
		},
		Database: "zitadel",
	})
	err := initialise.Init(db,
		initialise.VerifyUser(config.Username(), ""),
		initialise.VerifyDatabase(config.Database()),
		initialise.VerifyGrant(config.Database(), config.Username()))
	if err != nil {
		return err
	}
	return initialise.VerifyZitadel(db, *config)
}
