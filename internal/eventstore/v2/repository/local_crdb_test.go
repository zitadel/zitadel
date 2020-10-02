package repository

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
	"time"

	"github.com/caos/logging"
	"github.com/cockroachdb/cockroach-go/v2/testserver"
)

var (
	migrationsPath = os.ExpandEnv("${GOPATH}/src/github.com/caos/zitadel/migrations/cockroach")
	db             *sql.DB
)

func TestMain(m *testing.M) {
	ts, err := testserver.NewTestServer()
	if err != nil {
		logging.LogWithFields("REPOS-RvjLG", "error", err).Fatal("unable to start db")
	}

	db, err = sql.Open("postgres", ts.PGURL().String())
	if err != nil {
		logging.LogWithFields("REPOS-CF6dQ", "error", err).Fatal("unable to connect to db")
	}

	defer func() {
		db.Close()
		ts.Stop()
	}()

	if err = executeMigrations(); err != nil {
		logging.LogWithFields("REPOS-jehDD", "error", err).Fatal("migrations failed")
	}

	os.Exit(m.Run())
}

func TestInsert(t *testing.T) {
	tx, _ := db.Begin()

	var seq Sequence
	var d time.Time

	row := tx.QueryRow(crdbInsert, "event.type", "aggregate.type", "aggregate.id", Version("v1"), nil, Data(nil), "editor.user", "editor.service", "resource.owner", Sequence(0), false)
	err := row.Scan(&seq, &d)

	row = tx.QueryRow(crdbInsert, "event.type", "aggregate.type", "aggregate.id", Version("v1"), nil, Data(nil), "editor.user", "editor.service", "resource.owner", Sequence(1), true)
	err = row.Scan(&seq, &d)

	row = tx.QueryRow(crdbInsert, "event.type", "aggregate.type", "aggregate.id", Version("v1"), nil, Data(nil), "editor.user", "editor.service", "resource.owner", Sequence(0), false)
	err = row.Scan(&seq, &d)

	tx.Commit()

	rows, err := db.Query("select * from eventstore.events order by event_sequence")
	defer rows.Close()
	fmt.Println(err)

	fmt.Println(rows.Columns())
	for rows.Next() {
		i := make([]interface{}, 12)
		var id string
		rows.Scan(&id, &i[1], &i[2], &i[3], &i[4], &i[5], &i[6], &i[7], &i[8], &i[9], &i[10], &i[11])
		i[0] = id

		fmt.Println(i)
	}

	t.Fail()
}

func executeMigrations() error {
	files, err := migrationFilePaths()
	if err != nil {
		return err
	}
	sort.Sort(files)
	for _, file := range files {
		migration, err := ioutil.ReadFile(string(file))
		if err != nil {
			return err
		}
		transactionInMigration := strings.Contains(string(migration), "BEGIN;")
		exec := db.Exec
		var tx *sql.Tx
		if !transactionInMigration {
			tx, err = db.Begin()
			if err != nil {
				return fmt.Errorf("begin file: %v || err: %w", file, err)
			}
			exec = tx.Exec
		}
		if _, err = exec(string(migration)); err != nil {
			return fmt.Errorf("exec file: %v || err: %w", file, err)
		}
		if !transactionInMigration {
			if err = tx.Commit(); err != nil {
				return fmt.Errorf("commit file: %v || err: %w", file, err)
			}
		}
	}
	return nil
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
