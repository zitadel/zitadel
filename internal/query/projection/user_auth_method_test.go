package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestUserAuthMethodProjection_reduces(t *testing.T) {
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
			name: "reduceAddedPasswordless",
			args: args{
				event: getEvent(
					testEvent(
						user.HumanPasswordlessTokenAddedType,
						user.AggregateType,
						[]byte(`{
						"webAuthNTokenId": "token-id"
					}`),
					), user.HumanPasswordlessAddedEventMapper),
			},
			reduce: (&userAuthMethodProjection{}).reduceInitAuthMethod,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.user_auth_methods4 (token_id, creation_date, change_date, resource_owner, instance_id, user_id, sequence, state, method_type, name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (instance_id, user_id, method_type, token_id) DO UPDATE SET (creation_date, change_date, resource_owner, sequence, state, name) = (EXCLUDED.creation_date, EXCLUDED.change_date, EXCLUDED.resource_owner, EXCLUDED.sequence, EXCLUDED.state, EXCLUDED.name)",
							expectedArgs: []interface{}{
								"token-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								"agg-id",
								uint64(15),
								domain.MFAStateNotReady,
								domain.UserAuthMethodTypePasswordless,
								"",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAddedU2F",
			args: args{
				event: getEvent(
					testEvent(
						user.HumanU2FTokenAddedType,
						user.AggregateType,
						[]byte(`{
						"webAuthNTokenId": "token-id"
					}`),
					), user.HumanU2FAddedEventMapper),
			},
			reduce: (&userAuthMethodProjection{}).reduceInitAuthMethod,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.user_auth_methods4 (token_id, creation_date, change_date, resource_owner, instance_id, user_id, sequence, state, method_type, name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (instance_id, user_id, method_type, token_id) DO UPDATE SET (creation_date, change_date, resource_owner, sequence, state, name) = (EXCLUDED.creation_date, EXCLUDED.change_date, EXCLUDED.resource_owner, EXCLUDED.sequence, EXCLUDED.state, EXCLUDED.name)",
							expectedArgs: []interface{}{
								"token-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								"agg-id",
								uint64(15),
								domain.MFAStateNotReady,
								domain.UserAuthMethodTypeU2F,
								"",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAddedTOTP",
			args: args{
				event: getEvent(
					testEvent(
						user.HumanMFAOTPAddedType,
						user.AggregateType,
						[]byte(`{
					}`),
					), user.HumanOTPAddedEventMapper),
			},
			reduce: (&userAuthMethodProjection{}).reduceInitAuthMethod,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.user_auth_methods4 (token_id, creation_date, change_date, resource_owner, instance_id, user_id, sequence, state, method_type, name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (instance_id, user_id, method_type, token_id) DO UPDATE SET (creation_date, change_date, resource_owner, sequence, state, name) = (EXCLUDED.creation_date, EXCLUDED.change_date, EXCLUDED.resource_owner, EXCLUDED.sequence, EXCLUDED.state, EXCLUDED.name)",
							expectedArgs: []interface{}{
								"",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								"agg-id",
								uint64(15),
								domain.MFAStateNotReady,
								domain.UserAuthMethodTypeTOTP,
								"",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceVerifiedPasswordless",
			args: args{
				event: getEvent(
					testEvent(
						user.HumanPasswordlessTokenVerifiedType,
						user.AggregateType,
						[]byte(`{
						"webAuthNTokenId": "token-id",
						"webAuthNTokenName": "name"
					}`),
					), user.HumanPasswordlessVerifiedEventMapper),
			},
			reduce: (&userAuthMethodProjection{}).reduceActivateEvent,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.user_auth_methods4 SET (change_date, sequence, name, state) = ($1, $2, $3, $4) WHERE (user_id = $5) AND (method_type = $6) AND (resource_owner = $7) AND (token_id = $8) AND (instance_id = $9)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"name",
								domain.MFAStateReady,
								"agg-id",
								domain.UserAuthMethodTypePasswordless,
								"ro-id",
								"token-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceVerifiedU2F",
			args: args{
				event: getEvent(
					testEvent(
						user.HumanU2FTokenVerifiedType,
						user.AggregateType,
						[]byte(`{
						"webAuthNTokenId": "token-id",
						"webAuthNTokenName": "name"
					}`),
					), user.HumanU2FVerifiedEventMapper),
			},
			reduce: (&userAuthMethodProjection{}).reduceActivateEvent,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.user_auth_methods4 SET (change_date, sequence, name, state) = ($1, $2, $3, $4) WHERE (user_id = $5) AND (method_type = $6) AND (resource_owner = $7) AND (token_id = $8) AND (instance_id = $9)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"name",
								domain.MFAStateReady,
								"agg-id",
								domain.UserAuthMethodTypeU2F,
								"ro-id",
								"token-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceVerifiedTOTP",
			args: args{
				event: getEvent(
					testEvent(
						user.HumanMFAOTPVerifiedType,
						user.AggregateType,
						[]byte(`{
					}`),
					), user.HumanOTPVerifiedEventMapper),
			},
			reduce: (&userAuthMethodProjection{}).reduceActivateEvent,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.user_auth_methods4 SET (change_date, sequence, name, state) = ($1, $2, $3, $4) WHERE (user_id = $5) AND (method_type = $6) AND (resource_owner = $7) AND (token_id = $8) AND (instance_id = $9)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"",
								domain.MFAStateReady,
								"agg-id",
								domain.UserAuthMethodTypeTOTP,
								"ro-id",
								"",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAddedOTPSMS",
			args: args{
				event: getEvent(testEvent(
					user.HumanOTPSMSAddedType,
					user.AggregateType,
					nil,
				), eventstore.GenericEventMapper[user.HumanOTPSMSAddedEvent]),
			},
			reduce: (&userAuthMethodProjection{}).reduceAddAuthMethod,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.user_auth_methods4 (token_id, creation_date, change_date, resource_owner, instance_id, user_id, sequence, state, method_type, name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								"agg-id",
								uint64(15),
								domain.MFAStateReady,
								domain.UserAuthMethodTypeOTPSMS,
								"",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAddedOTPEmail",
			args: args{
				event: getEvent(testEvent(
					user.HumanOTPEmailAddedType,
					user.AggregateType,
					nil,
				), eventstore.GenericEventMapper[user.HumanOTPEmailAddedEvent]),
			},
			reduce: (&userAuthMethodProjection{}).reduceAddAuthMethod,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.user_auth_methods4 (token_id, creation_date, change_date, resource_owner, instance_id, user_id, sequence, state, method_type, name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								"agg-id",
								uint64(15),
								domain.MFAStateReady,
								domain.UserAuthMethodTypeOTPEmail,
								"",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceRemoveOTPPasswordless",
			args: args{
				event: getEvent(testEvent(
					user.HumanPasswordlessTokenRemovedType,
					user.AggregateType,
					[]byte(`{
						"webAuthNTokenId": "token-id"
					}`),
				), user.HumanPasswordlessRemovedEventMapper),
			},
			reduce: (&userAuthMethodProjection{}).reduceRemoveAuthMethod,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_auth_methods4 WHERE (user_id = $1) AND (method_type = $2) AND (resource_owner = $3) AND (instance_id = $4) AND (token_id = $5)",
							expectedArgs: []interface{}{
								"agg-id",
								domain.UserAuthMethodTypePasswordless,
								"ro-id",
								"instance-id",
								"token-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceRemoveOTPU2F",
			args: args{
				event: getEvent(testEvent(
					user.HumanU2FTokenRemovedType,
					user.AggregateType,
					[]byte(`{
						"webAuthNTokenId": "token-id"
					}`),
				), user.HumanU2FRemovedEventMapper),
			},
			reduce: (&userAuthMethodProjection{}).reduceRemoveAuthMethod,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_auth_methods4 WHERE (user_id = $1) AND (method_type = $2) AND (resource_owner = $3) AND (instance_id = $4) AND (token_id = $5)",
							expectedArgs: []interface{}{
								"agg-id",
								domain.UserAuthMethodTypeU2F,
								"ro-id",
								"instance-id",
								"token-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceRemoveTOTP",
			args: args{
				event: getEvent(testEvent(
					user.HumanMFAOTPRemovedType,
					user.AggregateType,
					nil,
				), user.HumanOTPRemovedEventMapper),
			},
			reduce: (&userAuthMethodProjection{}).reduceRemoveAuthMethod,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_auth_methods4 WHERE (user_id = $1) AND (method_type = $2) AND (resource_owner = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								"agg-id",
								domain.UserAuthMethodTypeTOTP,
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceRemoveOTPSMS",
			args: args{
				event: getEvent(testEvent(
					user.HumanOTPSMSRemovedType,
					user.AggregateType,
					nil,
				), eventstore.GenericEventMapper[user.HumanOTPSMSRemovedEvent]),
			},
			reduce: (&userAuthMethodProjection{}).reduceRemoveAuthMethod,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_auth_methods4 WHERE (user_id = $1) AND (method_type = $2) AND (resource_owner = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								"agg-id",
								domain.UserAuthMethodTypeOTPSMS,
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceRemovePhone",
			args: args{
				event: getEvent(testEvent(
					user.HumanPhoneRemovedType,
					user.AggregateType,
					nil,
				), user.HumanPhoneRemovedEventMapper),
			},
			reduce: (&userAuthMethodProjection{}).reduceRemoveAuthMethod,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_auth_methods4 WHERE (user_id = $1) AND (method_type = $2) AND (resource_owner = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								"agg-id",
								domain.UserAuthMethodTypeOTPSMS,
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceRemoveOTPEmail",
			args: args{
				event: getEvent(testEvent(
					user.HumanOTPEmailRemovedType,
					user.AggregateType,
					nil,
				), eventstore.GenericEventMapper[user.HumanOTPEmailRemovedEvent]),
			},
			reduce: (&userAuthMethodProjection{}).reduceRemoveAuthMethod,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_auth_methods4 WHERE (user_id = $1) AND (method_type = $2) AND (resource_owner = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								"agg-id",
								domain.UserAuthMethodTypeOTPEmail,
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceOwnerRemoved",
			reduce: (&userAuthMethodProjection{}).reduceOwnerRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						nil,
					), org.OrgRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_auth_methods4 WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
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
			reduce: reduceInstanceRemovedHelper(UserAuthMethodInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_auth_methods4 WHERE (instance_id = $1)",
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
			assertReduce(t, got, err, UserAuthMethodTable, tt.want)
		})
	}
}
