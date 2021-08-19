package handler

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
)

type executer struct {
	Handler
	projection string
}

type StatementExecuter interface {
	Execute(ctx context.Context, stmts []handler.Statement) ([]handler.Statement, error)
}

type PrintExecuter struct {
	executer
}

func NewPrintHandler(projection string) StatementExecuter {
	return &PrintExecuter{
		executer: executer{
			projection: projection,
		},
	}
}

func (e *PrintExecuter) Execute(ctx context.Context, stmts []handler.Statement) ([]handler.Statement, error) {
	for _, stmt := range stmts {
		select {
		case <-ctx.Done():
			logging.Log("V2-CroeS").Warn("stop processing")
			return nil, nil
		default:
			fmt.Println(stmt.AggregateType, stmt.Sequence)
		}
	}
	return nil, nil
}

type SQLExecuter struct {
	executer

	client *sql.DB

	sequences map[eventstore.AggregateType]uint64
}

const (
	currentSequenceStmt            = `SELECT current_sequence, aggregate_type FROM zitadel.projections.current_sequences WHERE projection_name = $1 FOR UPDATE`
	updateCurrentSequencesBaseStmt = `UPSERT INTO zitadel.projections.current_sequences (projection_name, aggregate_type, current_sequence, timestamp) VALUES `
)

type SQLExecuterConfig struct {
	Client         *sql.DB
	ProjectionName string
	PreSteps       []PreStep
	PostSteps      []PostStep
}

func NewSQLExecuter(config SQLExecuterConfig) StatementExecuter {
	return &SQLExecuter{
		client: config.Client,
		executer: executer{
			projection: config.ProjectionName,
			Handler: Handler{
				preSteps:  config.PreSteps,
				postSteps: config.PostSteps,
			},
		},
	}
}

func (e *SQLExecuter) Execute(ctx context.Context, stmts []handler.Statement) ([]handler.Statement, error) {
	if err := e.execPreSteps(); err != nil {
		logging.Log("V2-HwdJl").WithError(err).Warn("pre step failed")
		return stmts, err
	}
	defer e.execPostSteps()

	if err := e.currentSequences(); err != nil {
		return stmts, err
	}

	tx, err := e.client.Begin()
	if err != nil {
		//TODO: map err
		return stmts, err
	}

	for i, stmt := range stmts {
		select {
		case <-ctx.Done():
			return e.writeSequences(tx, stmts, i)
		default:
			if stmt.Sequence <= e.sequences[stmt.AggregateType] {
				// stmt already processed
				continue
			}
			if stmt.PreviousSequence > 0 && stmt.PreviousSequence != e.sequences[stmt.AggregateType] {
				logging.LogWithFields("CRDB-jJBJn", "projection", e.projection, "seq", stmt.Sequence, "prevSeq", stmt.PreviousSequence, "currentSeq", e.sequences[stmt.AggregateType]).Warn("sequences do not match")
				break
			}
			err := stmt.Execute(tx, e.projection)
			if err != nil {
				logging.LogWithFields("V2-9hUkD", "projection", e.projection, "seq", stmt.Sequence).Warn("unable to execute stmt")
				return e.writeSequences(tx, stmts, i)
			}
			e.sequences[stmt.AggregateType] = stmt.Sequence
		}
	}
	return e.writeSequences(tx, stmts, len(stmts))
}

func (e *SQLExecuter) currentSequences() error {
	rows, err := e.client.Query(currentSequenceStmt, e.projection)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			aggregateType eventstore.AggregateType
			sequence      uint64
		)

		err = rows.Scan(&sequence, &aggregateType)
		if err != nil {
			return errors.ThrowInternal(err, "CRDB-dbatK", "scan failed")
		}

		e.sequences[aggregateType] = sequence
	}

	if err = rows.Close(); err != nil {
		return errors.ThrowInternal(err, "CRDB-h5i5m", "close rows failed")
	}

	if err = rows.Err(); err != nil {
		return errors.ThrowInternal(err, "CRDB-O8zig", "errors in scanning rows")
	}

	return nil
}

func (e *SQLExecuter) writeSequences(tx *sql.Tx, stmts []handler.Statement, currentIdx int) ([]handler.Statement, error) {
	if currentIdx <= 0 {
		tx.Rollback()
		return stmts, nil
	}

	valueQueries := make([]string, 0, len(e.sequences))
	valueCounter := 0
	values := make([]interface{}, 0, len(e.sequences)*3)
	for aggregate, sequence := range e.sequences {
		valueQueries = append(valueQueries, "($"+strconv.Itoa(valueCounter+1)+", $"+strconv.Itoa(valueCounter+2)+", $"+strconv.Itoa(valueCounter+3)+", NOW())")
		valueCounter += 3
		values = append(values, e.projection, aggregate, sequence)
	}

	res, err := tx.Exec(updateCurrentSequencesBaseStmt+strings.Join(valueQueries, ", "), values...)
	if err != nil {
		tx.Rollback()
		return stmts, errors.ThrowInternal(err, "CRDB-TrH2Z", "unable to exec update sequence")
	}
	if rows, _ := res.RowsAffected(); rows < 1 {
		return stmts, errors.ThrowInternal(err, "V2-uP2zB", "unable to update sequences") //errSeqNotUpdated
	}

	return stmts[currentIdx+1:], nil
}
