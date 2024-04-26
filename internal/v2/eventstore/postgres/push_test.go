package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/v2/database/mock"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func Test_uniqueConstraints(t *testing.T) {
	type args struct {
		commands     []*command
		expectations []mock.Expectation
	}
	execErr := errors.New("exec err")
	tests := []struct {
		name      string
		args      args
		assertErr func(t *testing.T, err error) bool
	}{
		{
			name: "no commands",
			args: args{
				commands:     []*command{},
				expectations: []mock.Expectation{},
			},
			assertErr: expectNoErr,
		},
		{
			name: "command without constraints",
			args: args{
				commands: []*command{
					{},
				},
				expectations: []mock.Expectation{},
			},
			assertErr: expectNoErr,
		},
		{
			name: "add 1 constraint 1 command",
			args: args{
				commands: []*command{
					{
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: eventstore.Aggregate{
								Instance: "instance",
							},
						},
						uniqueConstraints: []*eventstore.UniqueConstraint{
							eventstore.NewAddEventUniqueConstraint("test", "id", "error"),
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectExec(
						"INSERT INTO eventstore.unique_constraints (instance_id, unique_type, unique_field) VALUES ($1, $2, $3)",
						mock.WithExecArgs("instance", "test", "id"),
						mock.WithExecRowsAffected(1),
					),
				},
			},
			assertErr: expectNoErr,
		},
		{
			name: "add 1 global constraint 1 command",
			args: args{
				commands: []*command{
					{
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: eventstore.Aggregate{
								Instance: "instance",
							},
						},
						uniqueConstraints: []*eventstore.UniqueConstraint{
							eventstore.NewAddGlobalUniqueConstraint("test", "id", "error"),
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectExec(
						"INSERT INTO eventstore.unique_constraints (instance_id, unique_type, unique_field) VALUES ($1, $2, $3)",
						mock.WithExecArgs("", "test", "id"),
						mock.WithExecRowsAffected(1),
					),
				},
			},
			assertErr: expectNoErr,
		},
		{
			name: "add 2 constraint 1 command",
			args: args{
				commands: []*command{
					{
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: eventstore.Aggregate{
								Instance: "instance",
							},
						},
						uniqueConstraints: []*eventstore.UniqueConstraint{
							eventstore.NewAddEventUniqueConstraint("test", "id", "error"),
							eventstore.NewAddEventUniqueConstraint("test", "id2", "error"),
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectExec(
						"INSERT INTO eventstore.unique_constraints (instance_id, unique_type, unique_field) VALUES ($1, $2, $3)",
						mock.WithExecArgs("instance", "test", "id"),
						mock.WithExecRowsAffected(1),
					),
					mock.ExpectExec(
						"INSERT INTO eventstore.unique_constraints (instance_id, unique_type, unique_field) VALUES ($1, $2, $3)",
						mock.WithExecArgs("instance", "test", "id2"),
						mock.WithExecRowsAffected(1),
					),
				},
			},
			assertErr: expectNoErr,
		},
		{
			name: "add 1 constraint per command",
			args: args{
				commands: []*command{
					{
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: eventstore.Aggregate{
								Instance: "instance",
							},
						},
						uniqueConstraints: []*eventstore.UniqueConstraint{
							eventstore.NewAddEventUniqueConstraint("test", "id", "error"),
						},
					},
					{
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: eventstore.Aggregate{
								Instance: "instance",
							},
						},
						uniqueConstraints: []*eventstore.UniqueConstraint{
							eventstore.NewAddEventUniqueConstraint("test", "id2", "error"),
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectExec(
						"INSERT INTO eventstore.unique_constraints (instance_id, unique_type, unique_field) VALUES ($1, $2, $3)",
						mock.WithExecArgs("instance", "test", "id"),
						mock.WithExecRowsAffected(1),
					),
					mock.ExpectExec(
						"INSERT INTO eventstore.unique_constraints (instance_id, unique_type, unique_field) VALUES ($1, $2, $3)",
						mock.WithExecArgs("instance", "test", "id2"),
						mock.WithExecRowsAffected(1),
					),
				},
			},
			assertErr: expectNoErr,
		},
		{
			name: "remove instance constraints 1 command",
			args: args{
				commands: []*command{
					{
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: eventstore.Aggregate{
								Instance: "instance",
							},
						},
						uniqueConstraints: []*eventstore.UniqueConstraint{
							eventstore.NewRemoveInstanceUniqueConstraints(),
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectExec(
						"DELETE FROM eventstore.unique_constraints WHERE instance_id = $1",
						mock.WithExecArgs("instance"),
						mock.WithExecRowsAffected(10),
					),
				},
			},
			assertErr: expectNoErr,
		},
		{
			name: "remove instance constraints 2 commands",
			args: args{
				commands: []*command{
					{
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: eventstore.Aggregate{
								Instance: "instance",
							},
						},
						uniqueConstraints: []*eventstore.UniqueConstraint{
							eventstore.NewRemoveInstanceUniqueConstraints(),
						},
					},
					{
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: eventstore.Aggregate{
								Instance: "instance",
							},
						},
						uniqueConstraints: []*eventstore.UniqueConstraint{
							eventstore.NewRemoveInstanceUniqueConstraints(),
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectExec(
						"DELETE FROM eventstore.unique_constraints WHERE instance_id = $1",
						mock.WithExecArgs("instance"),
						mock.WithExecRowsAffected(10),
					),
					mock.ExpectExec(
						"DELETE FROM eventstore.unique_constraints WHERE instance_id = $1",
						mock.WithExecArgs("instance"),
						mock.WithExecRowsAffected(0),
					),
				},
			},
			assertErr: expectNoErr,
		},
		{
			name: "remove 1 constraint 1 command",
			args: args{
				commands: []*command{
					{
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: eventstore.Aggregate{
								Instance: "instance",
							},
						},
						uniqueConstraints: []*eventstore.UniqueConstraint{
							eventstore.NewRemoveUniqueConstraint("test", "id"),
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectExec(
						`DELETE FROM eventstore.unique_constraints WHERE (instance_id = $1 AND unique_type = $2 AND unique_field = ( SELECT unique_field from ( SELECT instance_id, unique_type, unique_field FROM eventstore.unique_constraints WHERE instance_id = $1 AND unique_type = $2 AND unique_field = $3 UNION ALL SELECT instance_id, unique_type, unique_field FROM eventstore.unique_constraints WHERE instance_id = $1 AND unique_type = $2 AND unique_field = LOWER($3) ) AS case_insensitive_constraints LIMIT 1) )`,
						mock.WithExecArgs("instance", "test", "id"),
						mock.WithExecRowsAffected(1),
					),
				},
			},
			assertErr: expectNoErr,
		},
		{
			name: "remove 1 global constraint 1 command",
			args: args{
				commands: []*command{
					{
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: eventstore.Aggregate{
								Instance: "instance",
							},
						},
						uniqueConstraints: []*eventstore.UniqueConstraint{
							eventstore.NewRemoveGlobalUniqueConstraint("test", "id"),
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectExec(
						`DELETE FROM eventstore.unique_constraints WHERE (instance_id = $1 AND unique_type = $2 AND unique_field = ( SELECT unique_field from ( SELECT instance_id, unique_type, unique_field FROM eventstore.unique_constraints WHERE instance_id = $1 AND unique_type = $2 AND unique_field = $3 UNION ALL SELECT instance_id, unique_type, unique_field FROM eventstore.unique_constraints WHERE instance_id = $1 AND unique_type = $2 AND unique_field = LOWER($3) ) AS case_insensitive_constraints LIMIT 1) )`,
						mock.WithExecArgs("", "test", "id"),
						mock.WithExecRowsAffected(1),
					),
				},
			},
			assertErr: expectNoErr,
		},
		{
			name: "remove 2 constraints 1 command",
			args: args{
				commands: []*command{
					{
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: eventstore.Aggregate{
								Instance: "instance",
							},
						},
						uniqueConstraints: []*eventstore.UniqueConstraint{
							eventstore.NewRemoveUniqueConstraint("test", "id"),
							eventstore.NewRemoveUniqueConstraint("test", "id2"),
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectExec(
						`DELETE FROM eventstore.unique_constraints WHERE (instance_id = $1 AND unique_type = $2 AND unique_field = ( SELECT unique_field from ( SELECT instance_id, unique_type, unique_field FROM eventstore.unique_constraints WHERE instance_id = $1 AND unique_type = $2 AND unique_field = $3 UNION ALL SELECT instance_id, unique_type, unique_field FROM eventstore.unique_constraints WHERE instance_id = $1 AND unique_type = $2 AND unique_field = LOWER($3) ) AS case_insensitive_constraints LIMIT 1) )`,
						mock.WithExecArgs("instance", "test", "id"),
						mock.WithExecRowsAffected(1),
					),
					mock.ExpectExec(
						`DELETE FROM eventstore.unique_constraints WHERE (instance_id = $1 AND unique_type = $2 AND unique_field = ( SELECT unique_field from ( SELECT instance_id, unique_type, unique_field FROM eventstore.unique_constraints WHERE instance_id = $1 AND unique_type = $2 AND unique_field = $3 UNION ALL SELECT instance_id, unique_type, unique_field FROM eventstore.unique_constraints WHERE instance_id = $1 AND unique_type = $2 AND unique_field = LOWER($3) ) AS case_insensitive_constraints LIMIT 1) )`,
						mock.WithExecArgs("instance", "test", "id2"),
						mock.WithExecRowsAffected(1),
					),
				},
			},
			assertErr: expectNoErr,
		},
		{
			name: "remove 1 constraints per command",
			args: args{
				commands: []*command{
					{
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: eventstore.Aggregate{
								Instance: "instance",
							},
						},
						uniqueConstraints: []*eventstore.UniqueConstraint{
							eventstore.NewRemoveUniqueConstraint("test", "id"),
						},
					},
					{
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: eventstore.Aggregate{
								Instance: "instance",
							},
						},
						uniqueConstraints: []*eventstore.UniqueConstraint{
							eventstore.NewRemoveUniqueConstraint("test", "id2"),
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectExec(
						`DELETE FROM eventstore.unique_constraints WHERE (instance_id = $1 AND unique_type = $2 AND unique_field = ( SELECT unique_field from ( SELECT instance_id, unique_type, unique_field FROM eventstore.unique_constraints WHERE instance_id = $1 AND unique_type = $2 AND unique_field = $3 UNION ALL SELECT instance_id, unique_type, unique_field FROM eventstore.unique_constraints WHERE instance_id = $1 AND unique_type = $2 AND unique_field = LOWER($3) ) AS case_insensitive_constraints LIMIT 1) )`,
						mock.WithExecArgs("instance", "test", "id"),
						mock.WithExecRowsAffected(1),
					),
					mock.ExpectExec(
						`DELETE FROM eventstore.unique_constraints WHERE (instance_id = $1 AND unique_type = $2 AND unique_field = ( SELECT unique_field from ( SELECT instance_id, unique_type, unique_field FROM eventstore.unique_constraints WHERE instance_id = $1 AND unique_type = $2 AND unique_field = $3 UNION ALL SELECT instance_id, unique_type, unique_field FROM eventstore.unique_constraints WHERE instance_id = $1 AND unique_type = $2 AND unique_field = LOWER($3) ) AS case_insensitive_constraints LIMIT 1) )`,
						mock.WithExecArgs("instance", "test", "id2"),
						mock.WithExecRowsAffected(1),
					),
				},
			},
			assertErr: expectNoErr,
		},
		{
			name: "exec fails no error specified",
			args: args{
				commands: []*command{
					{
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: eventstore.Aggregate{
								Instance: "instance",
							},
						},
						uniqueConstraints: []*eventstore.UniqueConstraint{
							eventstore.NewAddEventUniqueConstraint("test", "id", ""),
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectExec(
						"INSERT INTO eventstore.unique_constraints (instance_id, unique_type, unique_field) VALUES ($1, $2, $3)",
						mock.WithExecArgs("instance", "test", "id"),
						mock.WithExecErr(execErr),
					),
				},
			},
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, zerrors.ThrowAlreadyExists(execErr, "POSTG-QzjyP", "Errors.Internal"))
				if !is {
					t.Errorf("no error expected got: %v", err)
				}
				return is
			},
		},
		{
			name: "exec fails error specified",
			args: args{
				commands: []*command{
					{
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: eventstore.Aggregate{
								Instance: "instance",
							},
						},
						uniqueConstraints: []*eventstore.UniqueConstraint{
							eventstore.NewAddEventUniqueConstraint("test", "id", "My.Error"),
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectExec(
						"INSERT INTO eventstore.unique_constraints (instance_id, unique_type, unique_field) VALUES ($1, $2, $3)",
						mock.WithExecArgs("instance", "test", "id"),
						mock.WithExecErr(execErr),
					),
				},
			},
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, zerrors.ThrowAlreadyExists(execErr, "POSTG-QzjyP", "My.Error"))
				if !is {
					t.Errorf("no error expected got: %v", err)
				}
				return is
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMock := mock.NewSQLMock(t, append([]mock.Expectation{mock.ExpectBegin(nil)}, tt.args.expectations...)...)
			tx, err := dbMock.DB.Begin()
			if err != nil {
				t.Errorf("unexpected error in begin: %v", err)
				t.FailNow()
			}
			err = uniqueConstraints(context.Background(), tx, tt.args.commands)
			tt.assertErr(t, err)
			dbMock.Assert(t)
		})
	}
}

var errReduce = errors.New("reduce err")

func Test_lockAggregates(t *testing.T) {
	type args struct {
		pushIntent   *eventstore.PushIntent
		expectations []mock.Expectation
	}
	type want struct {
		intents   []*intent
		assertErr func(t *testing.T, err error) bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "1 intent",
			args: args{
				pushIntent: eventstore.NewPushIntent(
					"instance",
					eventstore.AppendAggregate("owner", "testType", "testID"),
				),
				expectations: []mock.Expectation{
					mock.ExpectQuery(
						`WITH existing AS ((SELECT instance_id, aggregate_type, aggregate_id, "sequence" FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 AND aggregate_id = $3 AND owner = $4 ORDER BY "sequence" DESC LIMIT 1)) SELECT e.instance_id, e.owner, e.aggregate_type, e.aggregate_id, e.sequence FROM eventstore.events2 e JOIN existing ON e.instance_id = existing.instance_id AND e.aggregate_type = existing.aggregate_type AND e.aggregate_id = existing.aggregate_id AND e.sequence = existing.sequence FOR UPDATE`,
						mock.WithQueryArgs("instance", "testType", "testID", "owner"),
						mock.WithQueryResult(
							[]string{"instance_id", "owner", "aggregate_type", "aggregate_id", "sequence"},
							[][]driver.Value{
								{
									"instance",
									"owner",
									"testType",
									"testID",
									42,
								},
							},
						),
					),
				},
			},
			want: want{
				intents: []*intent{
					{
						PushAggregate: eventstore.NewPushIntent(
							"instance",
							eventstore.AppendAggregate("owner", "testType", "testID"),
						).Aggregates()[0],
						sequence: 42,
					},
				},
				assertErr: expectNoErr,
			},
		},
		{
			name: "two intents",
			args: args{
				pushIntent: eventstore.NewPushIntent(
					"instance",
					eventstore.AppendAggregate("owner", "testType", "testID"),
					eventstore.AppendAggregate("owner", "myType", "id"),
				),
				expectations: []mock.Expectation{
					mock.ExpectQuery(
						`WITH existing AS ((SELECT instance_id, aggregate_type, aggregate_id, "sequence" FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 AND aggregate_id = $3 AND owner = $4 ORDER BY "sequence" DESC LIMIT 1) UNION ALL (SELECT instance_id, aggregate_type, aggregate_id, "sequence" FROM eventstore.events2 WHERE instance_id = $5 AND aggregate_type = $6 AND aggregate_id = $7 AND owner = $8 ORDER BY "sequence" DESC LIMIT 1)) SELECT e.instance_id, e.owner, e.aggregate_type, e.aggregate_id, e.sequence FROM eventstore.events2 e JOIN existing ON e.instance_id = existing.instance_id AND e.aggregate_type = existing.aggregate_type AND e.aggregate_id = existing.aggregate_id AND e.sequence = existing.sequence FOR UPDATE`,
						mock.WithQueryArgs(
							"instance", "testType", "testID", "owner",
							"instance", "myType", "id", "owner",
						),
						mock.WithQueryResult(
							[]string{"instance_id", "owner", "aggregate_type", "aggregate_id", "sequence"},
							[][]driver.Value{
								{
									"instance",
									"owner",
									"testType",
									"testID",
									42,
								},
								{
									"instance",
									"owner",
									"myType",
									"id",
									17,
								},
							},
						),
					),
				},
			},
			want: want{
				intents: []*intent{
					{
						PushAggregate: eventstore.NewPushIntent(
							"instance",
							eventstore.AppendAggregate("owner", "testType", "testID"),
						).Aggregates()[0],
						sequence: 42,
					},
					{
						PushAggregate: eventstore.NewPushIntent(
							"instance",
							eventstore.AppendAggregate("owner", "myType", "id"),
						).Aggregates()[0],
						sequence: 17,
					},
				},
				assertErr: expectNoErr,
			},
		},
		{
			name: "1 intent aggregate not found",
			args: args{
				pushIntent: eventstore.NewPushIntent(
					"instance",
					eventstore.AppendAggregate("owner", "testType", "testID"),
				),
				expectations: []mock.Expectation{
					mock.ExpectQuery(
						`WITH existing AS ((SELECT instance_id, aggregate_type, aggregate_id, "sequence" FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 AND aggregate_id = $3 AND owner = $4 ORDER BY "sequence" DESC LIMIT 1)) SELECT e.instance_id, e.owner, e.aggregate_type, e.aggregate_id, e.sequence FROM eventstore.events2 e JOIN existing ON e.instance_id = existing.instance_id AND e.aggregate_type = existing.aggregate_type AND e.aggregate_id = existing.aggregate_id AND e.sequence = existing.sequence FOR UPDATE`,
						mock.WithQueryArgs("instance", "testType", "testID", "owner"),
						mock.WithQueryResult(
							[]string{"instance_id", "owner", "aggregate_type", "aggregate_id", "sequence"},
							[][]driver.Value{},
						),
					),
				},
			},
			want: want{
				intents: []*intent{
					{
						PushAggregate: eventstore.NewPushIntent(
							"instance",
							eventstore.AppendAggregate("owner", "testType", "testID"),
						).Aggregates()[0],
						sequence: 0,
					},
				},
				assertErr: expectNoErr,
			},
		},
		{
			name: "two intents none found",
			args: args{
				pushIntent: eventstore.NewPushIntent(
					"instance",
					eventstore.AppendAggregate("owner", "testType", "testID"),
					eventstore.AppendAggregate("owner", "myType", "id"),
				),
				expectations: []mock.Expectation{
					mock.ExpectQuery(
						`WITH existing AS ((SELECT instance_id, aggregate_type, aggregate_id, "sequence" FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 AND aggregate_id = $3 AND owner = $4 ORDER BY "sequence" DESC LIMIT 1) UNION ALL (SELECT instance_id, aggregate_type, aggregate_id, "sequence" FROM eventstore.events2 WHERE instance_id = $5 AND aggregate_type = $6 AND aggregate_id = $7 AND owner = $8 ORDER BY "sequence" DESC LIMIT 1)) SELECT e.instance_id, e.owner, e.aggregate_type, e.aggregate_id, e.sequence FROM eventstore.events2 e JOIN existing ON e.instance_id = existing.instance_id AND e.aggregate_type = existing.aggregate_type AND e.aggregate_id = existing.aggregate_id AND e.sequence = existing.sequence FOR UPDATE`,
						mock.WithQueryArgs(
							"instance", "testType", "testID", "owner",
							"instance", "myType", "id", "owner",
						),
						mock.WithQueryResult(
							[]string{"instance_id", "owner", "aggregate_type", "aggregate_id", "sequence"},
							[][]driver.Value{},
						),
					),
				},
			},
			want: want{
				intents: []*intent{
					{
						PushAggregate: eventstore.NewPushIntent(
							"instance",
							eventstore.AppendAggregate("owner", "testType", "testID"),
						).Aggregates()[0],
						sequence: 0,
					},
					{
						PushAggregate: eventstore.NewPushIntent(
							"instance",
							eventstore.AppendAggregate("owner", "myType", "id"),
						).Aggregates()[0],
						sequence: 0,
					},
				},
				assertErr: expectNoErr,
			},
		},
		{
			name: "two intents 1 found",
			args: args{
				pushIntent: eventstore.NewPushIntent(
					"instance",
					eventstore.AppendAggregate("owner", "testType", "testID"),
					eventstore.AppendAggregate("owner", "myType", "id"),
				),
				expectations: []mock.Expectation{
					mock.ExpectQuery(
						`WITH existing AS ((SELECT instance_id, aggregate_type, aggregate_id, "sequence" FROM eventstore.events2 WHERE instance_id = $1 AND aggregate_type = $2 AND aggregate_id = $3 AND owner = $4 ORDER BY "sequence" DESC LIMIT 1) UNION ALL (SELECT instance_id, aggregate_type, aggregate_id, "sequence" FROM eventstore.events2 WHERE instance_id = $5 AND aggregate_type = $6 AND aggregate_id = $7 AND owner = $8 ORDER BY "sequence" DESC LIMIT 1)) SELECT e.instance_id, e.owner, e.aggregate_type, e.aggregate_id, e.sequence FROM eventstore.events2 e JOIN existing ON e.instance_id = existing.instance_id AND e.aggregate_type = existing.aggregate_type AND e.aggregate_id = existing.aggregate_id AND e.sequence = existing.sequence FOR UPDATE`,
						mock.WithQueryArgs(
							"instance", "testType", "testID", "owner",
							"instance", "myType", "id", "owner",
						),
						mock.WithQueryResult(
							[]string{"instance_id", "owner", "aggregate_type", "aggregate_id", "sequence"},
							[][]driver.Value{
								{
									"instance",
									"owner",
									"myType",
									"id",
									17,
								},
							},
						),
					),
				},
			},
			want: want{
				intents: []*intent{
					{
						PushAggregate: eventstore.NewPushIntent(
							"instance",
							eventstore.AppendAggregate("owner", "testType", "testID"),
						).Aggregates()[0],
						sequence: 0,
					},
					{
						PushAggregate: eventstore.NewPushIntent(
							"instance",
							eventstore.AppendAggregate("owner", "myType", "id"),
						).Aggregates()[0],
						sequence: 17,
					},
				},
				assertErr: expectNoErr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMock := mock.NewSQLMock(t, append([]mock.Expectation{mock.ExpectBegin(nil)}, tt.args.expectations...)...)
			tx, err := dbMock.DB.Begin()
			if err != nil {
				t.Errorf("unexpected error in begin: %v", err)
				t.FailNow()
			}
			got, err := lockAggregates(context.Background(), tx, tt.args.pushIntent)
			tt.want.assertErr(t, err)
			dbMock.Assert(t)
			if len(got) != len(tt.want.intents) {
				t.Errorf("unexpected length of intents %d, want: %d", len(got), len(tt.want.intents))
				return
			}
			for i, gotten := range got {
				assertIntent(t, gotten, tt.want.intents[i])
			}
		})
	}
}

func assertIntent(t *testing.T, got, want *intent) {
	if got.sequence != want.sequence {
		t.Errorf("unexpected sequence %d want %d", got.sequence, want.sequence)
	}
	assertPushAggregate(t, got.PushAggregate, want.PushAggregate)
}

func assertPushAggregate(t *testing.T, got, want *eventstore.PushAggregate) {
	if !reflect.DeepEqual(got.Type(), want.Type()) {
		t.Errorf("unexpected Type %v, want: %v", got.Type(), want.Type())
	}
	if !reflect.DeepEqual(got.ID(), want.ID()) {
		t.Errorf("unexpected ID %v, want: %v", got.ID(), want.ID())
	}
	if !reflect.DeepEqual(got.Owner(), want.Owner()) {
		t.Errorf("unexpected Owner %v, want: %v", got.Owner(), want.Owner())
	}
	if !reflect.DeepEqual(got.Commands(), want.Commands()) {
		t.Errorf("unexpected Commands %v, want: %v", got.Commands(), want.Commands())
	}
	if !reflect.DeepEqual(got.Aggregate(), want.Aggregate()) {
		t.Errorf("unexpected Aggregate %v, want: %v", got.Aggregate(), want.Aggregate())
	}
	if !reflect.DeepEqual(got.CurrentSequence(), want.CurrentSequence()) {
		t.Errorf("unexpected CurrentSequence %v, want: %v", got.CurrentSequence(), want.CurrentSequence())
	}
}

func Test_push(t *testing.T) {
	type args struct {
		commands     []*command
		expectations []mock.Expectation
		reducer      *testReducer
	}
	type want struct {
		assertErr func(t *testing.T, err error) bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "1 aggregate 1 command",
			args: args{
				reducer: &testReducer{
					expectedReduces: 1,
				},
				commands: []*command{
					{
						intent: &intent{
							PushAggregate: eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0],
						},
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: *eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0].Aggregate(),
							Creator:  "gigi",
							Revision: 1,
							Type:     "test.type",
							Sequence: 1,
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectQuery(
						`INSERT INTO eventstore.events2 (instance_id, "owner", aggregate_type, aggregate_id, revision, creator, event_type, payload, "sequence", in_tx_order, created_at, "position") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp())) RETURNING created_at, "position"`,
						mock.WithQueryArgs(
							"instance",
							"owner",
							"testType",
							"testID",
							uint16(1),
							"gigi",
							"test.type",
							nil,
							uint32(1),
							0,
						),
						mock.WithQueryResult(
							[]string{"created_at", "position"},
							[][]driver.Value{
								{
									time.Now(),
									float64(123),
								},
							},
						),
					),
				},
			},
			want: want{
				assertErr: expectNoErr,
			},
		},
		{
			name: "1 aggregate 2 commands",
			args: args{
				reducer: &testReducer{
					expectedReduces: 2,
				},
				commands: []*command{
					{
						intent: &intent{
							PushAggregate: eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0],
						},
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: *eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0].Aggregate(),
							Creator:  "gigi",
							Revision: 1,
							Type:     "test.type",
							Sequence: 1,
						},
					},
					{
						intent: &intent{
							PushAggregate: eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0],
						},
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: *eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0].Aggregate(),
							Creator:  "gigi",
							Revision: 1,
							Type:     "test.type2",
							Sequence: 2,
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectQuery(
						`INSERT INTO eventstore.events2 (instance_id, "owner", aggregate_type, aggregate_id, revision, creator, event_type, payload, "sequence", in_tx_order, created_at, "position") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp())), ($11, $12, $13, $14, $15, $16, $17, $18, $19, $20, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp())) RETURNING created_at, "position"`,
						mock.WithQueryArgs(
							"instance",
							"owner",
							"testType",
							"testID",
							uint16(1),
							"gigi",
							"test.type",
							nil,
							uint32(1),
							0,
							"instance",
							"owner",
							"testType",
							"testID",
							uint16(1),
							"gigi",
							"test.type2",
							nil,
							uint32(2),
							1,
						),
						mock.WithQueryResult(
							[]string{"created_at", "position"},
							[][]driver.Value{
								{
									time.Now(),
									float64(123),
								},
								{
									time.Now(),
									float64(123.1),
								},
							},
						),
					),
				},
			},
			want: want{
				assertErr: expectNoErr,
			},
		},
		{
			name: "1 command per aggregate 2 aggregates",
			args: args{
				reducer: &testReducer{
					expectedReduces: 2,
				},
				commands: []*command{
					{
						intent: &intent{
							PushAggregate: eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0],
						},
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: *eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0].Aggregate(),
							Creator:  "gigi",
							Revision: 1,
							Type:     "test.type",
							Sequence: 1,
						},
					},
					{
						intent: &intent{
							PushAggregate: eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0],
						},
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: eventstore.Aggregate{
								ID:       "id2",
								Type:     "type2",
								Instance: "instance",
								Owner:    "owner",
							},
							Creator:  "gigi",
							Revision: 1,
							Type:     "test.type2",
							Sequence: 10,
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectQuery(
						`INSERT INTO eventstore.events2 (instance_id, "owner", aggregate_type, aggregate_id, revision, creator, event_type, payload, "sequence", in_tx_order, created_at, "position") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp())), ($11, $12, $13, $14, $15, $16, $17, $18, $19, $20, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp())) RETURNING created_at, "position"`,
						mock.WithQueryArgs(
							"instance",
							"owner",
							"testType",
							"testID",
							uint16(1),
							"gigi",
							"test.type",
							nil,
							uint32(1),
							0,
							"instance",
							"owner",
							"type2",
							"id2",
							uint16(1),
							"gigi",
							"test.type2",
							nil,
							uint32(10),
							1,
						),
						mock.WithQueryResult(
							[]string{"created_at", "position"},
							[][]driver.Value{
								{
									time.Now(),
									float64(123),
								},
								{
									time.Now(),
									float64(123.1),
								},
							},
						),
					),
				},
			},
			want: want{
				assertErr: expectNoErr,
			},
		},
		{
			name: "1 aggregate 1 command with payload",
			args: args{
				reducer: &testReducer{
					expectedReduces: 1,
				},
				commands: []*command{
					{
						intent: &intent{
							PushAggregate: eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0],
						},
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: *eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0].Aggregate(),
							Creator:  "gigi",
							Revision: 1,
							Type:     "test.type",
							Sequence: 1,
							Payload:  unmarshalPayload(`{"name": "gigi"}`),
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectQuery(
						`INSERT INTO eventstore.events2 (instance_id, "owner", aggregate_type, aggregate_id, revision, creator, event_type, payload, "sequence", in_tx_order, created_at, "position") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp())) RETURNING created_at, "position"`,
						mock.WithQueryArgs(
							"instance",
							"owner",
							"testType",
							"testID",
							uint16(1),
							"gigi",
							"test.type",
							unmarshalPayload(`{"name": "gigi"}`),
							uint32(1),
							0,
						),
						mock.WithQueryResult(
							[]string{"created_at", "position"},
							[][]driver.Value{
								{
									time.Now(),
									float64(123),
								},
							},
						),
					),
				},
			},
			want: want{
				assertErr: expectNoErr,
			},
		},
		{
			name: "command reducer",
			args: args{
				reducer: &testReducer{
					expectedReduces: 1,
				},
				commands: []*command{
					{
						intent: &intent{
							PushAggregate: eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0],
						},
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: *eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0].Aggregate(),
							Creator:  "gigi",
							Revision: 1,
							Type:     "test.type",
							Sequence: 1,
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectQuery(
						`INSERT INTO eventstore.events2 (instance_id, "owner", aggregate_type, aggregate_id, revision, creator, event_type, payload, "sequence", in_tx_order, created_at, "position") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp())) RETURNING created_at, "position"`,
						mock.WithQueryArgs(
							"instance",
							"owner",
							"testType",
							"testID",
							uint16(1),
							"gigi",
							"test.type",
							nil,
							uint32(1),
							0,
						),
						mock.WithQueryResult(
							[]string{"created_at", "position"},
							[][]driver.Value{
								{
									time.Now(),
									float64(123),
								},
							},
						),
					),
				},
			},
			want: want{
				assertErr: expectNoErr,
			},
		},
		{
			name: "command reducer err",
			args: args{
				reducer: &testReducer{
					expectedReduces: 1,
					shouldErr:       true,
				},
				commands: []*command{
					{
						intent: &intent{
							PushAggregate: eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0],
						},
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: *eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0].Aggregate(),
							Creator:  "gigi",
							Revision: 1,
							Type:     "test.type",
							Sequence: 1,
						},
					},
					{
						intent: &intent{
							PushAggregate: eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0],
						},
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: *eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0].Aggregate(),
							Creator:  "gigi",
							Revision: 1,
							Type:     "test.type2",
							Sequence: 2,
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectQuery(
						`INSERT INTO eventstore.events2 (instance_id, "owner", aggregate_type, aggregate_id, revision, creator, event_type, payload, "sequence", in_tx_order, created_at, "position") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp())), ($11, $12, $13, $14, $15, $16, $17, $18, $19, $20, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp())) RETURNING created_at, "position"`,
						mock.WithQueryArgs(
							"instance",
							"owner",
							"testType",
							"testID",
							uint16(1),
							"gigi",
							"test.type",
							nil,
							uint32(1),
							0,
							"instance",
							"owner",
							"testType",
							"testID",
							uint16(1),
							"gigi",
							"test.type2",
							nil,
							uint32(2),
							1,
						),
						mock.WithQueryResult(
							[]string{"created_at", "position"},
							[][]driver.Value{
								{
									time.Now(),
									float64(123),
								},
								{
									time.Now(),
									float64(123.1),
								},
							},
						),
					),
				},
			},
			want: want{
				assertErr: func(t *testing.T, err error) bool {
					is := errors.Is(err, errReduce)
					if !is {
						t.Errorf("no error expected got: %v", err)
					}
					return is
				},
			},
		},
		{
			name: "1 aggregate 2 commands",
			args: args{
				reducer: &testReducer{
					expectedReduces: 2,
				},
				commands: []*command{
					{
						intent: &intent{
							PushAggregate: eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0],
						},
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: *eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0].Aggregate(),
							Creator:  "gigi",
							Revision: 1,
							Type:     "test.type",
							Sequence: 1,
						},
					},
					{
						intent: &intent{
							PushAggregate: eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0],
						},
						Event: &eventstore.Event[eventstore.StoragePayload]{
							Aggregate: *eventstore.NewPushIntent(
								"instance",
								eventstore.AppendAggregate("owner", "testType", "testID"),
							).Aggregates()[0].Aggregate(),
							Creator:  "gigi",
							Revision: 1,
							Type:     "test.type2",
							Sequence: 2,
						},
					},
				},
				expectations: []mock.Expectation{
					mock.ExpectQuery(
						`INSERT INTO eventstore.events2 (instance_id, "owner", aggregate_type, aggregate_id, revision, creator, event_type, payload, "sequence", in_tx_order, created_at, "position") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp())), ($11, $12, $13, $14, $15, $16, $17, $18, $19, $20, statement_timestamp(), EXTRACT(EPOCH FROM clock_timestamp())) RETURNING created_at, "position"`,
						mock.WithQueryArgs(
							"instance",
							"owner",
							"testType",
							"testID",
							uint16(1),
							"gigi",
							"test.type",
							nil,
							uint32(1),
							0,
							"instance",
							"owner",
							"testType",
							"testID",
							uint16(1),
							"gigi",
							"test.type2",
							nil,
							uint32(2),
							1,
						),
						mock.WithQueryResult(
							[]string{"created_at", "position"},
							[][]driver.Value{
								{
									time.Now(),
									float64(123),
								},
								{
									time.Now(),
									float64(123.1),
								},
							},
						),
					),
				},
			},
			want: want{
				assertErr: expectNoErr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMock := mock.NewSQLMock(t, append([]mock.Expectation{mock.ExpectBegin(nil)}, tt.args.expectations...)...)
			tx, err := dbMock.DB.Begin()
			if err != nil {
				t.Errorf("unexpected error in begin: %v", err)
				t.FailNow()
			}
			err = push(context.Background(), tx, tt.args.reducer, tt.args.commands)
			tt.want.assertErr(t, err)
			dbMock.Assert(t)
			if tt.args.reducer != nil {
				tt.args.reducer.assert(t)
			}
		})
	}
}

func expectNoErr(t *testing.T, err error) bool {
	is := err == nil
	if !is {
		t.Errorf("no error expected got: %v", err)
	}
	return is
}
