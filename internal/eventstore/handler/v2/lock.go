package handler

import (
	"database/sql"
	"sync"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
)

type Locker interface {
	Lock() error
	Unlock()
}

type MutexLocker struct {
	mu sync.Mutex
}

func NewMutexLocker() Locker {
	return &MutexLocker{}
}

func (l *MutexLocker) Lock() error {
	l.mu.Lock()
	return nil
}

func (l *MutexLocker) Unlock() {
	l.mu.Unlock()
}

type SQLLocker struct {
	client         *sql.DB
	interval       time.Duration
	workerName     string
	projectionName string

	hasLocked bool
	lockedMu  sync.Mutex
}

const (
	lockStmt = "INSERT INTO zitadel.projections.locks l" +
		" (locker_id, locked_until, projection_name) VALUES ($1, now()+$2::INTERVAL, $3)" +
		" ON CONFLICT (projection_name)" +
		" DO UPDATE SET locker_id = $1, locked_until = now()+$2::INTERVAL" +
		" WHERE l.projection_name = $3 AND (l.locker_id = $1 OR l.locked_until < now())"
)

type SQLLockerConfig struct {
	Client         *sql.DB
	Interval       time.Duration
	WorkerName     string
	ProjectionName string
}

func NewSQLLocker(config SQLLockerConfig) Locker {
	// TODO: validate config
	return &SQLLocker{
		client:         config.Client,
		interval:       config.Interval,
		projectionName: config.ProjectionName,
		workerName:     config.WorkerName,
	}
}

func (l *SQLLocker) Lock() error {
	l.lockedMu.Lock()
	defer l.lockedMu.Unlock()

	if l.hasLocked {
		return nil
	}

	res, err := l.client.Exec(lockStmt, l.workerName, l.interval, l.projectionName)
	if err != nil {
		return errors.ThrowInternal(err, "CRDB-uaDoR", "unable to execute lock")
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.ThrowAlreadyExists(nil, "CRDB-mmi4J", "projection already locked")
	}

	l.hasLocked = true
	return nil
}

func (l *SQLLocker) Unlock() {
	l.lockedMu.Lock()
	defer l.lockedMu.Unlock()

	if !l.hasLocked {
		return
	}

	res, err := l.client.Exec(lockStmt, l.workerName, l.interval, l.projectionName)
	if err != nil {
		logging.LogWithFields("V2-kqfo1", "projection", l.projectionName).WithError(err).Warn("unable to unlock")
		return
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		logging.LogWithFields("V2-vtkcM", "projection", l.projectionName).Warn("lock already exists")
	}

	l.hasLocked = false
}
