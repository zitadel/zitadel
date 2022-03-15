package projection

import (
	"testing"
	"time"

	"github.com/lib/pq"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestPersonalAccessTokenProjection_reduces(t *testing.T) {
	type args struct {
		event func(t *testing.T) eventstore.Event
	}
	tests := []struct {
		name   string
		args   args
		reduce func(event eventstore.Event) (*handler.Statement, error)
		want   wantReduce
	}{
		{
			name: "reducePersonalAccessTokenAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.PersonalAccessTokenAddedType),
					user.AggregateType,
					[]byte(`{"tokenId": "tokenID", "expiration": "9999-12-31T23:59:59Z", "scopes": ["openid"]}`),
				), user.PersonalAccessTokenAddedEventMapper),
			},
			reduce: (&PersonalAccessTokenProjection{}).reducePersonalAccessTokenAdded,
			want: wantReduce{
				projection:       PersonalAccessTokenProjectionTable,
				aggregateType:    eventstore.AggregateType("user"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.personal_access_tokens (id, creation_date, change_date, resource_owner, instance_id, sequence, user_id, expiration, scopes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"tokenID",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								uint64(15),
								"agg-id",
								time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
								pq.StringArray{"openid"},
							},
						},
					},
				},
			},
		},
		{
			name: "reducePersonalAccessTokenRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.PersonalAccessTokenRemovedType),
					user.AggregateType,
					[]byte(`{"tokenId": "tokenID"}`),
				), user.PersonalAccessTokenRemovedEventMapper),
			},
			reduce: (&PersonalAccessTokenProjection{}).reducePersonalAccessTokenRemoved,
			want: wantReduce{
				projection:       PersonalAccessTokenProjectionTable,
				aggregateType:    eventstore.AggregateType("user"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.personal_access_tokens WHERE (id = $1)",
							expectedArgs: []interface{}{
								"tokenID",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.PersonalAccessTokenRemovedType),
					user.AggregateType,
					nil,
				), user.UserRemovedEventMapper),
			},
			reduce: (&PersonalAccessTokenProjection{}).reduceUserRemoved,
			want: wantReduce{
				projection:       PersonalAccessTokenProjectionTable,
				aggregateType:    eventstore.AggregateType("user"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.personal_access_tokens WHERE (user_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := baseEvent(t)
			got, err := tt.reduce(event)
			if _, ok := err.(errors.InvalidArgument); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, tt.want)
		})
	}
}
