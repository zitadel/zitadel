package sql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/caos/logging"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
)

var (
	migrationsPath = os.ExpandEnv("${GOPATH}/src/github.com/caos/zitadel/migrations/cockroach")
	testCRDBClient *sql.DB
)

func TestMain(m *testing.M) {
	ts, err := testserver.NewTestServer()
	if err != nil {
		logging.LogWithFields("REPOS-RvjLG", "error", err).Fatal("unable to start db")
	}

	testCRDBClient, err = sql.Open("postgres", ts.PGURL().String())

	if err != nil {
		logging.LogWithFields("REPOS-CF6dQ", "error", err).Fatal("unable to connect to db")
	}
	if err = testCRDBClient.Ping(); err != nil {
		logging.LogWithFields("REPOS-CF6dQ", "error", err).Fatal("unable to ping db")
	}

	defer func() {
		testCRDBClient.Close()
		ts.Stop()
	}()

	if err = executeMigrations(); err != nil {
		logging.LogWithFields("REPOS-jehDD", "error", err).Fatal("migrations failed")
	}

	os.Exit(m.Run())
}

func executeMigrations() error {
	files, err := migrationFilePaths()
	if err != nil {
		return err
	}
	sort.Sort(files)
	if err = setPasswordNULL(); err != nil {
		return err
	}
	if err = createFlywayHistory(); err != nil {
		return err
	}
	for _, file := range files {
		migrationData, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		migration := os.ExpandEnv(string(migrationData))
		if _, err = testCRDBClient.Exec(migration); err != nil {
			return fmt.Errorf("exec file: %v || err: %w", file, err)
		}
	}
	return nil
}

func setPasswordNULL() error {
	passwordNames := []string{
		"eventstorepassword",
		"managementpassword",
		"adminapipassword",
		"authpassword",
		"notificationpassword",
		"authzpassword",
		"queriespassword",
	}
	for _, name := range passwordNames {
		if err := os.Setenv(name, "NULL"); err != nil {
			return err
		}
	}
	return nil
}

func createFlywayHistory() error {
	_, err := testCRDBClient.Exec("CREATE TABLE defaultdb.flyway_schema_history(id TEXT, PRIMARY KEY(id));")
	return err
}

func fillUniqueData(unique_type, field string) error {
	_, err := testCRDBClient.Exec("INSERT INTO eventstore.unique_constraints (unique_type, unique_field) VALUES ($1, $2)", unique_type, field)
	return err
}

func fillAssets(id string, asset []byte) error {
	_, err := testCRDBClient.Exec("INSERT INTO eventstore.assets (id, asset) VALUES ($1, $2)", id, asset)
	return err
}

type migrationPaths []string

type version struct {
	major int
	minor int
}

func versionFromPath(s string) version {
	v := s[strings.Index(s, "/V")+2 : strings.Index(s, "__")]
	splitted := strings.Split(v, ".")
	res := version{}
	var err error
	if len(splitted) >= 1 {
		res.major, err = strconv.Atoi(splitted[0])
		if err != nil {
			panic(err)
		}
	}

	if len(splitted) >= 2 {
		res.minor, err = strconv.Atoi(splitted[1])
		if err != nil {
			panic(err)
		}
	}

	return res
}

func (a migrationPaths) Len() int      { return len(a) }
func (a migrationPaths) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a migrationPaths) Less(i, j int) bool {
	versionI := versionFromPath(a[i])
	versionJ := versionFromPath(a[j])

	return versionI.major < versionJ.major ||
		(versionI.major == versionJ.major && versionI.minor < versionJ.minor)
}

func migrationFilePaths() (migrationPaths, error) {
	files := make(migrationPaths, 0)
	err := filepath.Walk(migrationsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(info.Name(), ".sql") {
			return err
		}
		files = append(files, path)
		return nil
	})
	return files, err
}
