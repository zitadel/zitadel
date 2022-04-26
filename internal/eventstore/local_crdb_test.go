package eventstore_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/caos/logging"
	"github.com/cockroachdb/cockroach-go/v2/testserver"

	"github.com/zitadel/zitadel/cmd/admin/initialise"
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
	username := "zitadel"
	database := "zitadel"
	err := initialise.Initialise(db, initialise.VerifyUser(username, ""),
		initialise.VerifyDatabase(database),
		initialise.VerifyGrant(database, username))
	if err != nil {
		return err
	}
	return initialise.VerifyZitadel(db)
}
