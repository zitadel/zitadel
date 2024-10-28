package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestNotificationOIDCSession_reduces(t *testing.T) {
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
			name: "reduceOIDCSessionAdded",
			args: args{
				event: getEvent(
					testEvent(
						oidcsession.AddedType,
						oidcsession.AggregateType,
						[]byte(`{
									"audience": ["client-id"],
									"authMethods": [4],
									"authTime": "2024-10-28T07:58:03.956834+01:00",
									"clientID": "client-id",
									"nonce": "nonce",
									"preferredLanguage": "en",
									"scope": ["openid", "email", "profile"],
									"sessionID": "session-id",
									"userID": "user-id",
									"userResourceOwner": "user-resourceOwner"
}`),
					), eventstore.GenericEventMapper[oidcsession.AddedEvent]),
			},
			reduce: (&notificationOIDCSessionProjection{}).reduceOIDCSessionAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("oidc_session"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.notification_oidc_sessions (id, creation_date, change_date, resource_owner, instance_id, sequence, session_id, client_id, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								uint64(15),
								"session-id",
								"client-id",
								"user-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(NotificationOIDCSessionColumnInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.notification_oidc_sessions WHERE (instance_id = $1)",
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

			if ok := zerrors.IsErrorInvalidArgument(err); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, NotificationOIDCSessionProjectionTable, tt.want)
		})
	}
}
