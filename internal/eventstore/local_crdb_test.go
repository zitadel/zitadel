package eventstore_test

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/cmd/initialise"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/cockroach"
)

var (
	testCRDBClient *sql.DB
	localClient    *sql.DB
)

func TestMain(m *testing.M) {
	os.Exit(execTests(m))
}

func execTests(m *testing.M) int {
	ts, err := startTestserver()
	if err != nil {
		logging.Fatalf("unable to start testserver %v", err)
	}
	if err = startLocalDB(); err != nil {
		logging.Debug("unable to connect to local db")
	}
	defer func() {
		testCRDBClient.Close()
		ts.Stop()
		if localClient != nil {
			localClient.Close()
		}
	}()

	if err = initDB(testCRDBClient); err != nil {
		logging.WithFields("error", err).Fatal("migrations failed")
	}

	return m.Run()
}

func startTestserver() (testserver.TestServer, error) {
	ts, err := testserver.NewTestServer()
	if err != nil {
		logging.WithFields("error", err).Fatal("unable to start db")
		return nil, err
	}

	testCRDBClient, err = sql.Open("postgres", ts.PGURL().String())
	if err != nil {
		logging.WithFields("error", err).Fatal("unable to connect to db")
		return nil, err
	}
	if err = testCRDBClient.Ping(); err != nil {
		logging.WithFields("error", err).Fatal("unable to ping db")
		return nil, err
	}
	return ts, nil
}

func startLocalDB() (err error) {
	c := database.Config{}
	c.SetConnector(&cockroach.Config{
		Host:            "localhost",
		Port:            26257,
		Database:        "defaultdb",
		MaxOpenConns:    20,
		MaxIdleConns:    10,
		MaxConnLifetime: 30 * time.Minute,
		MaxConnIdleTime: 1 * time.Minute,
		User: cockroach.User{
			Username: "root",
			SSL:      cockroach.SSL{Mode: "disable"},
		},
	})
	if localClient, err = database.Connect(c, false); err != nil {
		return err
	}

	if err = localClient.Ping(); err != nil {
		logging.WithError(err).Fatal("unable to ping db")
		return err
	}

	if err = initDB(localClient); err != nil {
		return err
	}

	return nil
}

func initDB(db *sql.DB) error {
	initialise.ReadStmts("cockroach")
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
