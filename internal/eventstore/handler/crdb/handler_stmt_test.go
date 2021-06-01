package crdb

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	es_repo_mock "github.com/caos/zitadel/internal/eventstore/repository/mock"
)

var (
	errFilter = errors.New("filter err")
	errReduce = errors.New("reduce err")
)

func TestProjectionHandler_SearchQuery(t *testing.T) {
	type want struct {
		SearchQueryBuilder *eventstore.SearchQueryBuilder
		limit              uint64
		isErr              func(error) bool
		expectations       []mockExpectation
	}
	type fields struct {
		sequenceTable  string
		projectionName string
		aggregates     []eventstore.AggregateType
		events         []eventstore.EventType
		bulkLimit      uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "error in current sequence",
			fields: fields{
				sequenceTable:  "my_sequences",
				projectionName: "my_projection",
				aggregates:     []eventstore.AggregateType{"testAgg"},
				bulkLimit:      5,
			},
			want: want{
				limit: 0,
				isErr: func(err error) bool {
					return errors.Is(err, sql.ErrTxDone)
				},
				expectations: []mockExpectation{
					expectCurrentSequenceErr("my_sequences", "my_projection", sql.ErrTxDone),
				},
				SearchQueryBuilder: nil,
			},
		},
		{
			name: "only aggregates",
			fields: fields{
				sequenceTable:  "my_sequences",
				projectionName: "my_projection",
				aggregates:     []eventstore.AggregateType{"testAgg"},
				bulkLimit:      5,
			},
			want: want{
				limit: 5,
				isErr: func(err error) bool {
					return err == nil
				},
				expectations: []mockExpectation{
					expectCurrentSequence("my_sequences", "my_projection", 5),
				},
				SearchQueryBuilder: eventstore.
					NewSearchQueryBuilder(eventstore.ColumnsEvent, "testAgg").
					SequenceGreater(5).
					Limit(5),
			},
		},
		//TODO: discuss about event types in handler first
		// {
		// 	name: "aggregates and events",
		// 	fields: fields{
		// 		sequenceTable:  "my_sequences",
		// 		projectionName: "my_projection",
		// 		aggregates:     []eventstore.AggregateType{"testAgg"},
		// 		events:         []eventstore.EventType{"testAgg.added"},
		// 		bulkLimit:      5,
		// 	},
		// 	want: want{
		// 		limit: 5,
		// 		isErr: func(err error) bool {
		// 			return err == nil
		// 		},
		// 		expectations: []mockExpectation{
		// 			expectCurrentSequence("my_sequences", "my_projection", 5),
		// 		},
		// 		SearchQueryBuilder: eventstore.
		// 			NewSearchQueryBuilder(eventstore.ColumnsEvent, "testAgg").
		// 			EventTypes("testAgg.added").
		// 			SequenceGreater(5).
		// 			Limit(5),
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()

			h := &StatementHandler{
				ProjectionHandler: &handler.ProjectionHandler{
					ProjectionName: tt.fields.projectionName,
				},
				sequenceTable: tt.fields.sequenceTable,
				bulkLimit:     tt.fields.bulkLimit,
				aggregates:    tt.fields.aggregates,
				eventTypes:    tt.fields.events,
				client:        client,
			}

			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}

			query, limit, err := h.SearchQuery()
			if !tt.want.isErr(err) {
				t.Errorf("ProjectionHandler.prepareBulkStmts() error = %v", err)
				return
			}
			if !reflect.DeepEqual(query, tt.want.SearchQueryBuilder) {
				t.Errorf("unexpected query: expected %v, got %v", tt.want.SearchQueryBuilder, query)
			}
			if tt.want.limit != limit {
				t.Errorf("unexpected limit: got: %d want %d", limit, tt.want.limit)
			}

			mock.MatchExpectationsInOrder(true)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations not met: %v", err)
			}
		})
	}
}

func TestStatementHandler_Update(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		aggregates []eventstore.AggregateType
	}
	type want struct {
		expectations []mockExpectation
		isErr        func(error) bool
	}
	type args struct {
		ctx    context.Context
		stmts  []handler.Statement
		reduce handler.Reduce
	}
	tests := []struct {
		name   string
		fields fields
		want   want
		args   args
	}{
		{
			name: "begin fails",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				expectations: []mockExpectation{
					expectBeginErr(sql.ErrConnDone),
				},
				isErr: func(err error) bool {
					return errors.Is(err, sql.ErrConnDone)
				},
			},
		},
		{
			name: "current sequence fails",
			args: args{
				ctx: context.Background(),
			},
			want: want{
				expectations: []mockExpectation{
					expectBegin(),
					expectCurrentSequenceErr("my_sequences", "my_projection", sql.ErrTxDone),
					expectRollback(),
				},
				isErr: func(err error) bool {
					return errors.Is(err, sql.ErrTxDone)
				},
			},
		},
		{
			name: "fetch previous fails",
			fields: fields{
				eventstore: eventstore.NewEventstore(
					es_repo_mock.NewRepo(t).
						ExpectFilterEventsError(errFilter),
				),
				aggregates: []eventstore.AggregateType{"testAgg"},
			},
			args: args{
				ctx: context.Background(),
				stmts: []handler.Statement{
					NewNoOpStatement(6, 0),
				},
			},
			want: want{
				expectations: []mockExpectation{
					expectBegin(),
					expectCurrentSequence("my_sequences", "my_projection", 5),
					expectRollback(),
				},
				isErr: func(err error) bool {
					return errors.Is(err, errFilter)
				},
			},
		},
		{
			name: "no successful stmts",
			fields: fields{
				eventstore: eventstore.NewEventstore(
					es_repo_mock.NewRepo(t),
				),
				aggregates: []eventstore.AggregateType{"testAgg"},
			},
			args: args{
				ctx: context.Background(),
				stmts: []handler.Statement{
					NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 7, 6),
				},
			},
			want: want{
				expectations: []mockExpectation{
					expectBegin(),
					expectCurrentSequence("my_sequences", "my_projection", 5),
					expectCommit(),
				},
				isErr: func(err error) bool {
					return errors.Is(err, nil)
				},
			},
		},
		{
			name: "update current sequence fails",
			fields: fields{
				eventstore: eventstore.NewEventstore(
					es_repo_mock.NewRepo(t),
				),
				aggregates: []eventstore.AggregateType{"testAgg"},
			},
			args: args{
				ctx: context.Background(),
				stmts: []handler.Statement{
					NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 7, 5),
				},
			},
			want: want{
				expectations: []mockExpectation{
					expectBegin(),
					expectCurrentSequence("my_sequences", "my_projection", 5),
					expectSavePoint(),
					expectCreate("my_projection", []string{"col"}, []string{"$1"}),
					expectSavePointRelease(),
					expectUpdateCurrentSequenceNoRows("my_sequences", "my_projection", 7),
					expectRollback(),
				},
				isErr: func(err error) bool {
					return errors.Is(err, errSeqNotUpdated)
				},
			},
		},
		{
			name: "commit fails",
			fields: fields{
				eventstore: eventstore.NewEventstore(
					es_repo_mock.NewRepo(t),
				),
				aggregates: []eventstore.AggregateType{"testAgg"},
			},
			args: args{
				ctx: context.Background(),
				stmts: []handler.Statement{
					NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 7, 5),
				},
			},
			want: want{
				expectations: []mockExpectation{
					expectBegin(),
					expectCurrentSequence("my_sequences", "my_projection", 5),
					expectSavePoint(),
					expectCreate("my_projection", []string{"col"}, []string{"$1"}),
					expectSavePointRelease(),
					expectUpdateCurrentSequence("my_sequences", "my_projection", 7),
					expectCommitErr(sql.ErrConnDone),
				},
				isErr: func(err error) bool {
					return errors.Is(err, sql.ErrConnDone)
				},
			},
		},
		{
			name: "correct",
			fields: fields{
				eventstore: eventstore.NewEventstore(
					es_repo_mock.NewRepo(t),
				),
				aggregates: []eventstore.AggregateType{"testAgg"},
			},
			args: args{
				ctx: context.Background(),
				stmts: []handler.Statement{
					NewNoOpStatement(7, 5),
				},
			},
			want: want{
				expectations: []mockExpectation{
					expectBegin(),
					expectCurrentSequence("my_sequences", "my_projection", 5),
					expectUpdateCurrentSequence("my_sequences", "my_projection", 7),
					expectCommit(),
				},
				isErr: func(err error) bool {
					return errors.Is(err, nil)
				},
			},
		},
		{
			name: "fetch previous stmts no additional stmts",
			fields: fields{
				eventstore: eventstore.NewEventstore(
					es_repo_mock.NewRepo(t).ExpectFilterEvents(),
				),
				aggregates: []eventstore.AggregateType{"testAgg"},
			},
			args: args{
				ctx: context.Background(),
				stmts: []handler.Statement{
					NewNoOpStatement(7, 0),
				},
			},
			want: want{
				expectations: []mockExpectation{
					expectBegin(),
					expectCurrentSequence("my_sequences", "my_projection", 5),
					expectUpdateCurrentSequence("my_sequences", "my_projection", 7),
					expectCommit(),
				},
				isErr: func(err error) bool {
					return errors.Is(err, nil)
				},
			},
		},
		{
			name: "fetch previous stmts additional events",
			fields: fields{
				eventstore: eventstore.NewEventstore(
					es_repo_mock.NewRepo(t).ExpectFilterEvents(
						&repository.Event{
							Sequence:         6,
							PreviousSequence: 5,
						},
					),
				),
				aggregates: []eventstore.AggregateType{"testAgg"},
			},
			args: args{
				ctx: context.Background(),
				stmts: []handler.Statement{
					NewNoOpStatement(7, 0),
				},
				reduce: testReduce(),
			},
			want: want{
				expectations: []mockExpectation{
					expectBegin(),
					expectCurrentSequence("my_sequences", "my_projection", 5),
					expectUpdateCurrentSequence("my_sequences", "my_projection", 7),
					expectCommit(),
				},
				isErr: func(err error) bool {
					return errors.Is(err, nil)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()

			h := &StatementHandler{
				sequenceTable: "my_sequences",
				client:        client,
				ProjectionHandler: handler.NewProjectionHandler(handler.ProjectionHandlerConfig{
					HandlerConfig: handler.HandlerConfig{
						Eventstore: tt.fields.eventstore,
					},
					ProjectionName: "my_projection",
					RequeueEvery:   0,
				}),
				aggregates: tt.fields.aggregates,
			}

			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}

			_, err = h.Update(tt.args.ctx, tt.args.stmts, tt.args.reduce)
			if !tt.want.isErr(err) {
				t.Errorf("StatementHandler.Update() error = %v", err)
			}

			mock.MatchExpectationsInOrder(true)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations not met: %v", err)
			}
		})
	}
}

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
					es_repo_mock.NewRepo(t).ExpectFilterEventsError(errFilter),
				),
				aggregates: []eventstore.AggregateType{"testAgg"},
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, errFilter)
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
				reduce: testReduceErr(errReduce),
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
					return errors.Is(err, errReduce)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StatementHandler{
				ProjectionHandler: handler.NewProjectionHandler(handler.ProjectionHandlerConfig{
					HandlerConfig: handler.HandlerConfig{
						Eventstore: tt.fields.eventstore,
					},
					ProjectionName: "my_projection",
					RequeueEvery:   0,
				}),
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
		projectionName    string
		maxFailureCount   uint
		failedEventsTable string
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
					NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val1",
						},
					}, 5, 2),
				},
				currentSeq: 5,
			},
			want: want{
				expectations: []mockExpectation{},
				idx:          -1,
			},
		},
		{
			name: "previous sequence higher than sequence",
			fields: fields{
				projectionName: "my_projection",
			},
			args: args{
				stmts: []handler.Statement{
					NewCreateStatement([]handler.Column{
						{
							Name:  "col1",
							Value: "val1",
						},
					}, 5, 0),
					NewCreateStatement([]handler.Column{
						{
							Name:  "col2",
							Value: "val2",
						},
					}, 8, 7),
				},
				currentSeq: 2,
			},
			want: want{
				expectations: []mockExpectation{
					expectSavePoint(),
					expectCreate("my_projection", []string{"col1"}, []string{"$1"}),
					expectSavePointRelease(),
				},
				idx: 0,
			},
		},
		{
			name: "execute fails not continue",
			fields: fields{
				projectionName:    "my_projection",
				maxFailureCount:   5,
				failedEventsTable: "failed_events",
			},
			args: args{
				stmts: []handler.Statement{
					NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 5, 0),
					NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 6, 5),
					NewCreateStatement([]handler.Column{
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
					expectFailureCount("failed_events", "my_projection", 6, 4),
				},
				idx: 0,
			},
		},
		{
			name: "execute fails continue",
			fields: fields{
				projectionName:    "my_projection",
				maxFailureCount:   5,
				failedEventsTable: "failed_events",
			},
			args: args{
				stmts: []handler.Statement{
					NewCreateStatement([]handler.Column{
						{
							Name:  "col1",
							Value: "val1",
						},
					}, 5, 0),
					NewCreateStatement([]handler.Column{
						{
							Name:  "col2",
							Value: "val2",
						},
					}, 6, 5),
					NewCreateStatement([]handler.Column{
						{
							Name:  "col3",
							Value: "val3",
						},
					}, 7, 6),
				},
				currentSeq: 2,
			},
			want: want{
				expectations: []mockExpectation{
					expectSavePoint(),
					expectCreate("my_projection", []string{"col1"}, []string{"$1"}),
					expectSavePointRelease(),
					expectSavePoint(),
					expectCreateErr("my_projection", []string{"col2"}, []string{"$1"}, sql.ErrConnDone),
					expectSavePointRollback(),
					expectFailureCount("failed_events", "my_projection", 6, 5),
					expectSavePoint(),
					expectCreate("my_projection", []string{"col3"}, []string{"$1"}),
					expectSavePointRelease(),
				},
				idx: 2,
			},
		},
		{
			name: "correct",
			fields: fields{
				projectionName: "my_projection",
			},
			args: args{
				stmts: []handler.Statement{
					NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 5, 0),
					NewCreateStatement([]handler.Column{
						{
							Name:  "col",
							Value: "val",
						},
					}, 6, 5),
					NewCreateStatement([]handler.Column{
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
			client, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()

			h := NewStatementHandler(
				context.Background(),
				StatementHandlerConfig{
					ProjectionHandlerConfig: handler.ProjectionHandlerConfig{
						HandlerConfig: handler.HandlerConfig{
							Eventstore: nil,
						},
						ProjectionName: tt.fields.projectionName,
						RequeueEvery:   0,
					},
					Client:            client,
					FailedEventsTable: tt.fields.failedEventsTable,
					MaxFailureCount:   tt.fields.maxFailureCount,
				},
			)

			mock.ExpectBegin()

			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}

			mock.ExpectCommit()

			tx, err := client.Begin()
			if err != nil {
				t.Fatalf("unexpected err in begin: %v", err)
			}

			idx := h.executeStmts(tx, tt.args.stmts, tt.args.currentSeq)
			if idx != tt.want.idx {
				t.Errorf("unexpected index want: %d got %d", tt.want.idx, idx)
			}

			if err := tx.Commit(); err != nil {
				t.Fatalf("unexpected err in commit: %v", err)
			}

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
				stmt: NewCreateStatement([]handler.Column{
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
				stmt: NewCreateStatement([]handler.Column{
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
				stmt: NewCreateStatement([]handler.Column{
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
				stmt: NewNoOpStatement(1, 0),
			},
			want: want{
				isErr: func(err error) bool {
					return err == nil
				},
				expectations: []mockExpectation{},
			},
		},
		{
			name: "with op",
			fields: fields{
				projectionName: "my_projection",
			},
			args: args{
				stmt: NewCreateStatement([]handler.Column{
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
				ProjectionHandler: &handler.ProjectionHandler{
					ProjectionName: tt.fields.projectionName,
				},
			}

			client, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()

			mock.ExpectBegin()

			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}

			mock.ExpectCommit()

			tx, err := client.Begin()
			if err != nil {
				t.Fatalf("unexpected err in begin: %v", err)
			}

			err = h.executeStmt(tx, tt.args.stmt)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}

			if err = tx.Commit(); err != nil {
				t.Fatalf("unexpected err in begin: %v", err)
			}

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
					expectCurrentSequenceErr("my_table", "my_projection", sql.ErrConnDone),
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
				ProjectionHandler: &handler.ProjectionHandler{
					ProjectionName: tt.fields.projectionName,
				},
				sequenceTable: tt.fields.sequenceTable,
			}

			client, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()

			mock.ExpectBegin()

			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}

			mock.ExpectCommit()

			tx, err := client.Begin()
			if err != nil {
				t.Fatalf("unexpected err in begin: %v", err)
			}

			seq, err := h.currentSequence(tx.QueryRow)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}

			if err = tx.Commit(); err != nil {
				t.Fatalf("unexpected err in commit: %v", err)
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
				ProjectionHandler: &handler.ProjectionHandler{
					ProjectionName: tt.fields.projectionName,
				},
				sequenceTable: tt.fields.sequenceTable,
			}

			client, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()

			mock.ExpectBegin()
			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}
			mock.ExpectCommit()

			tx, err := client.Begin()
			if err != nil {
				t.Fatalf("unexpected error in begin: %v", err)
			}

			err = h.updateCurrentSequence(tx, tt.args.stmt)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error: %v", err)
			}

			err = tx.Commit()
			if err != nil {
				t.Fatalf("unexpected error in commit: %v", err)
			}

			mock.MatchExpectationsInOrder(true)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations not met: %v", err)
			}
		})
	}
}

func testReduce(stmts ...handler.Statement) handler.Reduce {
	return func(event eventstore.EventReader) ([]handler.Statement, error) {
		return []handler.Statement{
			NewNoOpStatement(event.Sequence(), event.PreviousSequence()),
		}, nil
	}
}

func testReduceErr(err error) handler.Reduce {
	return func(event eventstore.EventReader) ([]handler.Statement, error) {
		return nil, err
	}
}
