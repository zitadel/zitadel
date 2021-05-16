package crdb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	es_repo_mock "github.com/caos/zitadel/internal/eventstore/repository/mock"
)

var (
	filterErr = errors.New("filter err")
	reduceErr = errors.New("reduce err")
)

type mockExpectation func(sqlmock.Sqlmock)

func TestProjectionHandler_fetchPreviousStmts(t *testing.T) {
	type args struct {
		ctx        context.Context
		stmtSeq    uint64
		currentSeq uint64
		reduce     handler.Reduce
	}
	type want struct {
		stmtCount int
		isErr     func(error) bool
	}
	type fields struct {
		eventstore *eventstore.Eventstore
		aggregates []eventstore.AggregateType
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   want
	}{
		{
			name: "eventstore returns err",
			args: args{
				ctx:    context.Background(),
				reduce: testReduce(),
			},
			fields: fields{
				eventstore: eventstore.NewEventstore(
					es_repo_mock.NewRepo(t).ExpectFilterEventsError(filterErr),
				),
				aggregates: []eventstore.AggregateType{"testAgg"},
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, filterErr)
				},
			},
		},
		{
			name: "no events found",
			args: args{
				ctx:    context.Background(),
				reduce: testReduce(),
			},
			fields: fields{
				eventstore: eventstore.NewEventstore(
					es_repo_mock.NewRepo(t).ExpectFilterEvents(),
				),
				aggregates: []eventstore.AggregateType{"testAgg"},
			},
			want: want{
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
		{
			name: "found events",
			args: args{
				ctx:    context.Background(),
				reduce: testReduce(),
			},
			fields: fields{
				eventstore: eventstore.NewEventstore(
					es_repo_mock.NewRepo(t).ExpectFilterEvents(
						&repository.Event{
							ID:               "id",
							Sequence:         1,
							PreviousSequence: 0,
							CreationDate:     time.Now(),
							Type:             "test.added",
							Version:          "v1",
							AggregateID:      "testid",
							AggregateType:    "testAgg",
						},
						&repository.Event{
							ID:               "id",
							Sequence:         2,
							PreviousSequence: 1,
							CreationDate:     time.Now(),
							Type:             "test.changed",
							Version:          "v1",
							AggregateID:      "testid",
							AggregateType:    "testAgg",
						},
					),
				),
				aggregates: []eventstore.AggregateType{"testAgg"},
			},
			want: want{
				stmtCount: 2,
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
		{
			name: "reduce fails",
			args: args{
				ctx:    context.Background(),
				reduce: testReduceErr(reduceErr),
			},
			fields: fields{
				eventstore: eventstore.NewEventstore(
					es_repo_mock.NewRepo(t).ExpectFilterEvents(
						&repository.Event{
							ID:               "id",
							Sequence:         1,
							PreviousSequence: 0,
							CreationDate:     time.Now(),
							Type:             "test.added",
							Version:          "v1",
							AggregateID:      "testid",
							AggregateType:    "testAgg",
						},
					),
				),
				aggregates: []eventstore.AggregateType{"testAgg"},
			},
			want: want{
				stmtCount: 0,
				isErr: func(err error) bool {
					return errors.Is(err, reduceErr)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StatementHandler{
				eventstore: tt.fields.eventstore,
				aggregates: tt.fields.aggregates,
			}
			stmts, err := h.fetchPreviousStmts(tt.args.ctx, tt.args.stmtSeq, tt.args.currentSeq, tt.args.reduce)
			if !tt.want.isErr(err) {
				t.Errorf("ProjectionHandler.prepareBulkStmts() error = %v", err)
				return
			}
			if tt.want.stmtCount != len(stmts) {
				t.Errorf("unexpected length of stmts: got: %d want %d", len(stmts), tt.want.stmtCount)
			}
		})
	}
}

func TestStatementHandler_executeStmts(t *testing.T) {
	type fields struct {
		projectionName string
	}
	type args struct {
		stmts      []handler.Statement
		currentSeq uint64
	}
	type want struct {
		expectations []mockExpectation
		idx          int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "already inserted",
			fields: fields{
				projectionName: "my_projection",
			},
			args: args{
				stmts: []handler.Statement{
					handler.NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 5, 0),
					handler.NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 4, 3),
					handler.NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 6, 5),
				},
				currentSeq: 2,
			},
			want: want{
				expectations: []mockExpectation{
					expectSavePoint(),
					expectCreate("my_projection", []string{"col"}, []string{"$1"}),
					expectSavePointRelease(),
				},
				idx: 0,
			},
		},
		{
			name: "previous sequence higher than sequence",
			fields: fields{
				projectionName: "my_projection",
			},
			args: args{
				stmts: []handler.Statement{
					handler.NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 5, 0),
					handler.NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 8, 7),
					handler.NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 9, 8),
				},
				currentSeq: 2,
			},
			want: want{
				expectations: []mockExpectation{
					expectSavePoint(),
					expectCreate("my_projection", []string{"col"}, []string{"$1"}),
					expectSavePointRelease(),
				},
				idx: 0,
			},
		},
		{
			name: "execute fails",
			fields: fields{
				projectionName: "my_projection",
			},
			args: args{
				stmts: []handler.Statement{
					handler.NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 5, 0),
					handler.NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 6, 5),
					handler.NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 7, 6),
				},
				currentSeq: 2,
			},
			want: want{
				expectations: []mockExpectation{
					expectSavePoint(),
					expectCreate("my_projection", []string{"col"}, []string{"$1"}),
					expectSavePointRelease(),
					expectSavePoint(),
					expectCreateErr("my_projection", []string{"col"}, []string{"$1"}, sql.ErrConnDone),
					expectSavePointRollback(),
				},
				idx: 0,
			},
		},
		{
			name: "correct",
			fields: fields{
				projectionName: "my_projection",
			},
			args: args{
				stmts: []handler.Statement{
					handler.NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 5, 0),
					handler.NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 6, 5),
					handler.NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 7, 6),
				},
				currentSeq: 2,
			},
			want: want{
				expectations: []mockExpectation{
					expectSavePoint(),
					expectCreate("my_projection", []string{"col"}, []string{"$1"}),
					expectSavePointRelease(),
					expectSavePoint(),
					expectCreate("my_projection", []string{"col"}, []string{"$1"}),
					expectSavePointRelease(),
					expectSavePoint(),
					expectCreate("my_projection", []string{"col"}, []string{"$1"}),
					expectSavePointRelease(),
				},
				idx: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StatementHandler{
				projectionName: tt.fields.projectionName,
			}

			client, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}

			mock.ExpectBegin()

			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}

			mock.ExpectCommit()

			tx, _ := client.Begin()

			idx := h.executeStmts(tx, tt.args.stmts, tt.args.currentSeq)
			if idx != tt.want.idx {
				t.Errorf("unexpected index want: %d got %d", tt.want.idx, idx)
			}

			tx.Commit()

			mock.MatchExpectationsInOrder(true)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations not met: %v", err)
			}
		})
	}
}

func TestStatementHandler_executeStmt(t *testing.T) {
	type fields struct {
		projectionName string
	}
	type args struct {
		stmt handler.Statement
	}
	type want struct {
		expectations []mockExpectation
		isErr        func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "create savepoint fails",
			fields: fields{
				projectionName: "my_projection",
			},
			args: args{
				stmt: handler.Statement{},
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, sql.ErrConnDone)
				},
				expectations: []mockExpectation{
					expectSavePointErr(sql.ErrConnDone),
				},
			},
		},
		{
			name: "execute fails",
			fields: fields{
				projectionName: "my_projection",
			},
			args: args{
				stmt: handler.NewCreateStatement([]handler.Column{
					{
						Name:  "col",
						Value: "val",
					},
				}, 1, 0),
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, sql.ErrNoRows)
				},
				expectations: []mockExpectation{
					expectSavePoint(),
					expectCreateErr("my_projection", []string{"col"}, []string{"$1"}, sql.ErrNoRows),
					expectSavePointRollback(),
				},
			},
		},
		{
			name: "rollback savepoint fails",
			fields: fields{
				projectionName: "my_projection",
			},
			args: args{
				stmt: handler.NewCreateStatement([]handler.Column{
					{
						Name:  "col",
						Value: "val",
					},
				}, 1, 0),
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, sql.ErrConnDone)
				},
				expectations: []mockExpectation{
					expectSavePoint(),
					expectCreateErr("my_projection", []string{"col"}, []string{"$1"}, sql.ErrNoRows),
					expectSavePointRollbackErr(sql.ErrConnDone),
				},
			},
		},
		{
			name: "no op",
			fields: fields{
				projectionName: "my_projection",
			},
			args: args{
				stmt: handler.NewNoOpStatement(1, 0),
			},
			want: want{
				isErr: func(err error) bool {
					return err == nil
				},
				expectations: []mockExpectation{
					expectSavePoint(),
					expectSavePointRelease(),
				},
			},
		},
		{
			name: "with op",
			fields: fields{
				projectionName: "my_projection",
			},
			args: args{
				stmt: handler.NewCreateStatement([]handler.Column{
					{
						Name:  "col",
						Value: "val",
					},
				}, 1, 0),
			},
			want: want{
				isErr: func(err error) bool {
					return err == nil
				},
				expectations: []mockExpectation{
					expectSavePoint(),
					expectCreate("my_projection", []string{"col"}, []string{"$1"}),
					expectSavePointRelease(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StatementHandler{
				projectionName: tt.fields.projectionName,
			}

			client, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}

			mock.ExpectBegin()

			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}

			mock.ExpectCommit()

			tx, _ := client.Begin()

			err = h.executeStmt(tx, tt.args.stmt)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}

			tx.Commit()

			mock.MatchExpectationsInOrder(true)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations not met: %v", err)
			}
		})
	}
}

func TestStatementHandler_currentSequence(t *testing.T) {
	type fields struct {
		sequenceTable  string
		projectionName string
	}
	type args struct {
		stmt handler.Statement
	}
	type want struct {
		expectations []mockExpectation
		isErr        func(error) bool
		seq          uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "error in query",
			fields: fields{
				sequenceTable:  "my_table",
				projectionName: "my_projection",
			},
			args: args{
				stmt: handler.Statement{},
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, sql.ErrConnDone)
				},
				expectations: []mockExpectation{
					expectCurrentSequenceWithErr("my_table", "my_projection", sql.ErrConnDone),
				},
			},
		},
		{
			name: "no rows",
			fields: fields{
				sequenceTable:  "my_table",
				projectionName: "my_projection",
			},
			args: args{
				stmt: handler.Statement{},
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, nil)
				},
				expectations: []mockExpectation{
					expectCurrentSequenceNotExists("my_table", "my_projection"),
				},
				seq: 0,
			},
		},
		{
			name: "scan err",
			fields: fields{
				sequenceTable:  "my_table",
				projectionName: "my_projection",
			},
			args: args{
				stmt: handler.Statement{},
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, sql.ErrTxDone)
				},
				expectations: []mockExpectation{
					expectCurrentSequenceScanErr("my_table", "my_projection"),
				},
				seq: 0,
			},
		},
		{
			name: "found",
			fields: fields{
				sequenceTable:  "my_table",
				projectionName: "my_projection",
			},
			args: args{
				stmt: handler.Statement{},
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, nil)
				},
				expectations: []mockExpectation{
					expectCurrentSequence("my_table", "my_projection", 5),
				},
				seq: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StatementHandler{
				sequenceTable:  tt.fields.sequenceTable,
				projectionName: tt.fields.projectionName,
			}

			client, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}

			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}

			seq, err := h.currentSequence(client.QueryRow)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}

			mock.MatchExpectationsInOrder(true)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations not met: %v", err)
			}

			if seq != tt.want.seq {
				t.Errorf("expected sequence %d got %d", tt.want.seq, seq)
			}
		})
	}
}

func TestStatementHandler_updateCurrentSequence(t *testing.T) {
	type fields struct {
		sequenceTable  string
		projectionName string
	}
	type args struct {
		stmt handler.Statement
	}
	type want struct {
		expectations []mockExpectation
		isErr        func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "update sequence fails",
			fields: fields{
				sequenceTable:  "my_table",
				projectionName: "my_projection",
			},
			args: args{
				stmt: handler.Statement{
					Sequence: 5,
				},
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, sql.ErrConnDone)
				},
				expectations: []mockExpectation{
					expectUpdateCurrentSequenceErr("my_table", "my_projection", 5, sql.ErrConnDone),
				},
			},
		},
		{
			name: "update sequence returns no rows",
			fields: fields{
				sequenceTable:  "my_table",
				projectionName: "my_projection",
			},
			args: args{
				stmt: handler.Statement{
					Sequence: 5,
				},
			},
			want: want{
				isErr: func(err error) bool {
					return errors.As(err, &errSeqNotUpdated)
				},
				expectations: []mockExpectation{
					expectUpdateCurrentSequenceNoRows("my_table", "my_projection", 5),
				},
			},
		},
		{
			name: "correct",
			fields: fields{
				sequenceTable:  "my_table",
				projectionName: "my_projection",
			},
			args: args{
				stmt: handler.Statement{
					Sequence: 5,
				},
			},
			want: want{
				isErr: func(err error) bool {
					return err == nil
				},
				expectations: []mockExpectation{
					expectUpdateCurrentSequence("my_table", "my_projection", 5),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StatementHandler{
				sequenceTable:  tt.fields.sequenceTable,
				projectionName: tt.fields.projectionName,
			}

			client, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}

			mock.ExpectBegin()
			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}
			mock.ExpectCommit()

			tx, _ := client.Begin()

			err = h.updateCurrentSequence(tx, tt.args.stmt)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}

			tx.Commit()

			mock.MatchExpectationsInOrder(true)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations not met: %v", err)
			}
		})
	}
}

func expectCreate(projectionName string, columnNames, placeholders []string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		args := make([]driver.Value, len(columnNames))
		for i := 0; i < len(columnNames); i++ {
			args[i] = sqlmock.AnyArg()
			placeholders[i] = `\` + placeholders[i]
		}
		m.ExpectExec("INSERT INTO " + projectionName + ` \(` + strings.Join(columnNames, ", ") + `\) VALUES \(` + strings.Join(placeholders, ", ") + `\)`).
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func expectCreateErr(projectionName string, columnNames, placeholders []string, err error) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		args := make([]driver.Value, len(columnNames))
		for i := 0; i < len(columnNames); i++ {
			args[i] = sqlmock.AnyArg()
			placeholders[i] = `\` + placeholders[i]
		}
		m.ExpectExec("INSERT INTO " + projectionName + ` \(` + strings.Join(columnNames, ", ") + `\) VALUES \(` + strings.Join(placeholders, ", ") + `\)`).
			WithArgs(args...).
			WillReturnError(err)
	}
}

func expectSavePoint() func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("SAVEPOINT push_stmt").
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func expectSavePointErr(err error) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("SAVEPOINT push_stmt").
			WillReturnError(err)
	}
}

func expectSavePointRollback() func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("ROLLBACK TO SAVEPOINT push_stmt").
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func expectSavePointRollbackErr(err error) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("ROLLBACK TO SAVEPOINT push_stmt").
			WillReturnError(err)
	}
}

func expectSavePointRelease() func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("RELEASE push_stmt").
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
}

func expectCurrentSequence(tableName, projection string, seq uint64) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(`WITH seq AS \(SELECT current_sequence FROM ` + tableName + ` WHERE view_name = \$1 FOR UPDATE\)
SELECT 
	IF\(
		COUNT\(current_sequence\) > 0, 
		\(SELECT current_sequence FROM seq\),
		0 AS current_sequence
	\) 
FROM seq`).
			WithArgs(
				projection,
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"current_sequence"}).
					AddRow(seq),
			)
	}
}

func expectCurrentSequenceWithErr(tableName, projection string, err error) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(`WITH seq AS \(SELECT current_sequence FROM ` + tableName + ` WHERE view_name = \$1 FOR UPDATE\)
SELECT 
	IF\(
		COUNT\(current_sequence\) > 0, 
		\(SELECT current_sequence FROM seq\),
		0 AS current_sequence
	\) 
FROM seq`).
			WithArgs(
				projection,
			).
			WillReturnError(err)
	}
}

func expectCurrentSequenceNotExists(tableName, projection string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(`WITH seq AS \(SELECT current_sequence FROM ` + tableName + ` WHERE view_name = \$1 FOR UPDATE\)
SELECT 
	IF\(
		COUNT\(current_sequence\) > 0, 
		\(SELECT current_sequence FROM seq\),
		0 AS current_sequence
	\) 
FROM seq`).
			WithArgs(
				projection,
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"current_sequence"}).
					AddRow(0),
			)
	}
}

func expectCurrentSequenceScanErr(tableName, projection string) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(`WITH seq AS \(SELECT current_sequence FROM ` + tableName + ` WHERE view_name = \$1 FOR UPDATE\)
SELECT 
	IF\(
		COUNT\(current_sequence\) > 0, 
		\(SELECT current_sequence FROM seq\),
		0 AS current_sequence
	\) 
FROM seq`).
			WithArgs(
				projection,
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"current_sequence"}).
					RowError(0, sql.ErrTxDone).
					AddRow(0),
			)
	}
}

func expectUpdateCurrentSequence(tableName, projection string, seq uint64) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("UPSERT INTO "+tableName+` \(view_name, current_sequence, timestamp\) VALUES \(\$1, \$2, NOW\(\)\)`).
			WithArgs(
				projection,
				seq,
			).
			WillReturnResult(
				sqlmock.NewResult(1, 1),
			)
	}
}

func expectUpdateCurrentSequenceErr(tableName, projection string, seq uint64, err error) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("UPSERT INTO "+tableName+` \(view_name, current_sequence, timestamp\) VALUES \(\$1, \$2, NOW\(\)\)`).
			WithArgs(
				projection,
				seq,
			).
			WillReturnError(err)
	}
}

func expectUpdateCurrentSequenceNoRows(tableName, projection string, seq uint64) func(sqlmock.Sqlmock) {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec("UPSERT INTO "+tableName+` \(view_name, current_sequence, timestamp\) VALUES \(\$1, \$2, NOW\(\)\)`).
			WithArgs(
				projection,
				seq,
			).
			WillReturnResult(
				sqlmock.NewResult(0, 0),
			)
	}
}

func testReduce(stmts ...handler.Statement) handler.Reduce {
	return func(event eventstore.EventReader) ([]handler.Statement, error) {
		return []handler.Statement{
			handler.NewNoOpStatement(event.Sequence(), event.PreviousSequence()),
		}, nil
	}
}

func testReduceErr(err error) handler.Reduce {
	return func(event eventstore.EventReader) ([]handler.Statement, error) {
		return nil, err
	}
}
