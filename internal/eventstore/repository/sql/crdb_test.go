package sql

import (
	"database/sql"
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

func TestCRDB_placeholder(t *testing.T) {
	type args struct {
		query string
	}
	type res struct {
		query string
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no placeholders",
			args: args{
				query: "SELECT * FROM eventstore.events2",
			},
			res: res{
				query: "SELECT * FROM eventstore.events2",
			},
		},
		{
			name: "one placeholder",
			args: args{
				query: "SELECT * FROM eventstore.events2 WHERE aggregate_type = ?",
			},
			res: res{
				query: "SELECT * FROM eventstore.events2 WHERE aggregate_type = $1",
			},
		},
		{
			name: "multiple placeholders",
			args: args{
				query: "SELECT * FROM eventstore.events2 WHERE aggregate_type = ? AND aggregate_id = ? LIMIT ?",
			},
			res: res{
				query: "SELECT * FROM eventstore.events2 WHERE aggregate_type = $1 AND aggregate_id = $2 LIMIT $3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{}
			if query := db.placeholder(tt.args.query); query != tt.res.query {
				t.Errorf("CRDB.placeholder() = %v, want %v", query, tt.res.query)
			}
		})
	}
}

func TestCRDB_operation(t *testing.T) {
	type res struct {
		op string
	}
	type args struct {
		operation repository.Operation
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "no op",
			args: args{
				operation: repository.Operation(-1),
			},
			res: res{
				op: "",
			},
		},
		{
			name: "greater",
			args: args{
				operation: repository.OperationGreater,
			},
			res: res{
				op: ">",
			},
		},
		{
			name: "less",
			args: args{
				operation: repository.OperationLess,
			},
			res: res{
				op: "<",
			},
		},
		{
			name: "equals",
			args: args{
				operation: repository.OperationEquals,
			},
			res: res{
				op: "=",
			},
		},
		{
			name: "in",
			args: args{
				operation: repository.OperationIn,
			},
			res: res{
				op: "=",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{}
			if got := db.operation(tt.args.operation); got != tt.res.op {
				t.Errorf("CRDB.operation() = %v, want %v", got, tt.res.op)
			}
		})
	}
}

func TestCRDB_conditionFormat(t *testing.T) {
	type res struct {
		format string
	}
	type args struct {
		operation repository.Operation
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "default",
			args: args{
				operation: repository.OperationEquals,
			},
			res: res{
				format: "%s %s ?",
			},
		},
		{
			name: "in",
			args: args{
				operation: repository.OperationIn,
			},
			res: res{
				format: "%s %s ANY(?)",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{}
			if got := db.conditionFormat(tt.args.operation); got != tt.res.format {
				t.Errorf("CRDB.conditionFormat() = %v, want %v", got, tt.res.format)
			}
		})
	}
}

func TestCRDB_columnName(t *testing.T) {
	type res struct {
		name string
	}
	type args struct {
		field repository.Field
		useV1 bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "invalid field",
			args: args{
				field: repository.Field(-1),
			},
			res: res{
				name: "",
			},
		},
		{
			name: "aggregate id",
			args: args{
				field: repository.FieldAggregateID,
			},
			res: res{
				name: "aggregate_id",
			},
		},
		{
			name: "aggregate type",
			args: args{
				field: repository.FieldAggregateType,
			},
			res: res{
				name: "aggregate_type",
			},
		},
		{
			name: "editor service",
			args: args{
				field: repository.FieldEditorService,
				useV1: true,
			},
			res: res{
				name: "editor_service",
			},
		},
		{
			name: "editor service v2",
			args: args{
				field: repository.FieldEditorService,
			},
			res: res{
				name: "",
			},
		},
		{
			name: "editor user",
			args: args{
				field: repository.FieldEditorUser,
				useV1: true,
			},
			res: res{
				name: "editor_user",
			},
		},
		{
			name: "editor user v2",
			args: args{
				field: repository.FieldEditorUser,
			},
			res: res{
				name: "creator",
			},
		},
		{
			name: "event type",
			args: args{
				field: repository.FieldEventType,
			},
			res: res{
				name: "event_type",
			},
		},
		{
			name: "latest sequence",
			args: args{
				field: repository.FieldSequence,
				useV1: true,
			},
			res: res{
				name: "event_sequence",
			},
		},
		{
			name: "latest sequence v2",
			args: args{
				field: repository.FieldSequence,
			},
			res: res{
				name: `"sequence"`,
			},
		},
		{
			name: "resource owner",
			args: args{
				field: repository.FieldResourceOwner,
				useV1: true,
			},
			res: res{
				name: "resource_owner",
			},
		},
		{
			name: "resource owner v2",
			args: args{
				field: repository.FieldResourceOwner,
			},
			res: res{
				name: `"owner"`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &CRDB{}
			if got := db.columnName(tt.args.field, tt.args.useV1); got != tt.res.name {
				t.Errorf("CRDB.operation() = %v, want %v", got, tt.res.name)
			}
		})
	}
}

func generateEvent(t *testing.T, aggregateID string, opts ...func(*repository.Event)) *repository.Event {
	t.Helper()
	e := &repository.Event{
		AggregateID:   aggregateID,
		AggregateType: eventstore.AggregateType(t.Name()),
		EditorUser:    "user",
		ResourceOwner: sql.NullString{String: "ro", Valid: true},
		Typ:           "test.created",
		Version:       "v1",
		Pos:           42,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}
