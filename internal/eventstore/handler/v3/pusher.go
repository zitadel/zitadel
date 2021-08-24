package handler

import (
	"context"
	"database/sql"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
)

type Pusher struct {
	client          *sql.DB
	projectionName  string
	maxFailureCount uint

	stmtsMu sync.Mutex
	stmts   handler.Statements

	timerSet bool
	t        *time.Timer
	interval time.Duration
}

type PusherConfig struct {
	Client          *sql.DB
	Interval        time.Duration
	MaxFailureCount uint

	projectionName string
}

func NewPusher(config PusherConfig) *Pusher {
	p := &Pusher{
		client:          config.Client,
		projectionName:  config.projectionName,
		interval:        config.Interval,
		maxFailureCount: config.MaxFailureCount,
		t:               time.NewTimer(0),
	}

	//unitialized timer
	//https://github.com/golang/go/issues/12721
	<-p.t.C

	return p
}

func (p *Pusher) Process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			err := p.push(context.Background())
			logging.LogWithFields("V3-lN8e1", "projection", p.projectionName).OnError(err).Info("push before shutdown failed")
			return
		case <-p.t.C:
			err := p.push(context.Background())
			if err != nil {
				logging.LogWithFields("V3-HZXHY", "projection", p.projectionName).WithError(err).Info("timed push failed")
				p.t.Reset(p.interval)
				p.timerSet = true
				continue
			}
			p.timerSet = false
		}
	}
}

const (
	lockCurrentSequenceStmt        = `SELECT current_sequence, aggregate_type FROM zitadel.projections.current_sequences WHERE projection_name = $1 FOR UPDATE`
	updateCurrentSequencesBaseStmt = `UPSERT INTO zitadel.projections.current_sequences (projection_name, aggregate_type, current_sequence, timestamp) VALUES `
)

func (p *Pusher) appendStmts(stmts ...handler.Statement) {
	p.stmtsMu.Lock()
	p.stmts = append(p.stmts, stmts...)
	if !p.timerSet {
		p.timerSet = true
		p.t.Reset(p.interval)
	}
	p.stmtsMu.Unlock()
}

func (p *Pusher) push(ctx context.Context) error {
	p.stmtsMu.Lock()
	defer p.stmtsMu.Unlock()

	tx, err := p.client.Begin()
	if err != nil {
		//TODO: map err
		return err
	}

	sequences, err := p.currentSequences(tx)
	if err != nil {
		return err
	}

	sort.Sort(p.stmts)
	for i, stmt := range p.stmts {
		select {
		case <-ctx.Done():
			return p.writeSequences(tx, sequences, i)
		default:
			if stmt.Sequence <= sequences[stmt.AggregateType] {
				// stmt already processed
				continue
			}
			if stmt.PreviousSequence > 0 && stmt.PreviousSequence != sequences[stmt.AggregateType] {
				logging.LogWithFields("CRDB-jJBJn", "projection", p.projectionName, "seq", stmt.Sequence, "prevSeq", stmt.PreviousSequence, "currentSeq", sequences[stmt.AggregateType]).Warn("sequences do not match")
				break
			}
			err := p.execStmt(tx, stmt)
			if err == nil {
				sequences[stmt.AggregateType] = stmt.Sequence
				continue
			}

			shouldContinue := p.handleFailedStmt(tx, stmt, err)
			if !shouldContinue {
				break
			}
			sequences[stmt.AggregateType] = stmt.Sequence
		}
	}
	return p.writeSequences(tx, sequences, len(p.stmts)-1)
}

func (p *Pusher) execStmt(tx *sql.Tx, stmt handler.Statement) error {
	if stmt.IsNoop() {
		return nil
	}
	_, err := tx.Exec("SAVEPOINT push_stmt")
	if err != nil {
		return errors.ThrowInternal(err, "CRDB-i1wp6", "unable to create savepoint")
	}
	err = stmt.Execute(tx, p.projectionName)
	if err != nil {
		_, rollbackErr := tx.Exec("ROLLBACK TO SAVEPOINT push_stmt")
		if rollbackErr != nil {
			return errors.ThrowInternal(rollbackErr, "CRDB-zzp3P", "rollback to savepoint failed")
		}
		return errors.ThrowInternal(err, "CRDB-oRkaN", "unable execute stmt")
	}
	_, err = tx.Exec("RELEASE push_stmt")
	if err != nil {
		return errors.ThrowInternal(err, "CRDB-qWgwT", "unable to release savepoint")
	}
	return nil
}

func (p *Pusher) currentSequences(tx *sql.Tx) (map[eventstore.AggregateType]uint64, error) {
	sequences := make(map[eventstore.AggregateType]uint64)
	rows, err := tx.Query(lockCurrentSequenceStmt, p.projectionName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			aggregateType eventstore.AggregateType
			sequence      uint64
		)

		err = rows.Scan(&sequence, &aggregateType)
		if err != nil {
			return nil, errors.ThrowInternal(err, "CRDB-dbatK", "scan failed")
		}

		sequences[aggregateType] = sequence
	}

	if err = rows.Close(); err != nil {
		return nil, errors.ThrowInternal(err, "CRDB-h5i5m", "close rows failed")
	}

	if err = rows.Err(); err != nil {
		return nil, errors.ThrowInternal(err, "CRDB-O8zig", "errors in scanning rows")
	}

	return sequences, nil
}

func (p *Pusher) writeSequences(tx *sql.Tx, sequences map[eventstore.AggregateType]uint64, currentIdx int) error {
	if currentIdx <= 0 {
		tx.Rollback()
		return nil
	}

	valueQueries := make([]string, 0, len(sequences))
	valueCounter := 0
	values := make([]interface{}, 0, len(sequences)*3)
	for aggregate, sequence := range sequences {
		valueQueries = append(valueQueries, "($"+strconv.Itoa(valueCounter+1)+", $"+strconv.Itoa(valueCounter+2)+", $"+strconv.Itoa(valueCounter+3)+", NOW())")
		valueCounter += 3
		values = append(values, p.projectionName, aggregate, sequence)
	}

	res, err := tx.Exec(updateCurrentSequencesBaseStmt+strings.Join(valueQueries, ", "), values...)
	if err != nil {
		tx.Rollback()
		return errors.ThrowInternal(err, "CRDB-TrH2Z", "unable to exec update sequence")
	}
	if rows, _ := res.RowsAffected(); rows < 1 {
		tx.Rollback()
		return errors.ThrowInternal(err, "V2-uP2zB", "unable to update sequences") //errSeqNotUpdated
	}
	if err := tx.Commit(); err != nil {
		return errors.ThrowInternal(err, "V3-8kLut", "unable to commit updates")
	}

	p.stmts = p.stmts[currentIdx+1:]
	if len(p.stmts) == 0 {
		if !p.t.Stop() {
			<-p.t.C
		}
	}

	return nil
}
