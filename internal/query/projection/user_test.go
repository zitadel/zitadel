package projection

import (
	"database/sql"
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestUserProjection_reduces(t *testing.T) {
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
			name: "reduceHumanAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanAddedType),
					user.AggregateType,
					[]byte(`{
						"username": "user-name",
						"firstName": "first-name",
						"lastName": "last-name",
						"nickName": "nick-name",
						"displayName": "display-name",
						"preferredLanguage": "ch-DE",
						"gender": 1,
						"email": "email@zitadel.com",
						"phone": "+41 00 000 00 00"
					}`),
				), user.HumanAddedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanAdded,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.users8 (id, creation_date, change_date, resource_owner, instance_id, state, sequence, username, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								domain.UserStateActive,
								uint64(15),
								"user-name",
								domain.UserTypeHuman,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.users8_humans (user_id, instance_id, first_name, last_name, nick_name, display_name, preferred_language, gender, email, phone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								"first-name",
								"last-name",
								&sql.NullString{String: "nick-name", Valid: true},
								&sql.NullString{String: "display-name", Valid: true},
								&sql.NullString{String: "ch-DE", Valid: true},
								&sql.NullInt16{Int16: int16(domain.GenderFemale), Valid: true},
								"email@zitadel.com",
								&sql.NullString{String: "+41 00 000 00 00", Valid: true},
							},
						},
						{
							expectedStmt: "INSERT INTO projections.users8_notifications (user_id, instance_id, last_email, last_phone, password_set) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								"email@zitadel.com",
								&sql.NullString{String: "+41 00 000 00 00", Valid: true},
								false,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserV1Added",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserV1AddedType),
					user.AggregateType,
					[]byte(`{
						"username": "user-name",
						"firstName": "first-name",
						"lastName": "last-name",
						"nickName": "nick-name",
						"displayName": "display-name",
						"preferredLanguage": "ch-DE",
						"gender": 1,
						"email": "email@zitadel.com",
						"phone": "+41 00 000 00 00"
					}`),
				), user.HumanAddedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanAdded,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.users8 (id, creation_date, change_date, resource_owner, instance_id, state, sequence, username, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								domain.UserStateActive,
								uint64(15),
								"user-name",
								domain.UserTypeHuman,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.users8_humans (user_id, instance_id, first_name, last_name, nick_name, display_name, preferred_language, gender, email, phone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								"first-name",
								"last-name",
								&sql.NullString{String: "nick-name", Valid: true},
								&sql.NullString{String: "display-name", Valid: true},
								&sql.NullString{String: "ch-DE", Valid: true},
								&sql.NullInt16{Int16: int16(domain.GenderFemale), Valid: true},
								"email@zitadel.com",
								&sql.NullString{String: "+41 00 000 00 00", Valid: true},
							},
						},
						{
							expectedStmt: "INSERT INTO projections.users8_notifications (user_id, instance_id, last_email, last_phone, password_set) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								"email@zitadel.com",
								&sql.NullString{String: "+41 00 000 00 00", Valid: true},
								false,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceHumanAdded NULLs",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanAddedType),
					user.AggregateType,
					[]byte(`{
						"username": "user-name",
						"firstName": "first-name",
						"lastName": "last-name",
						"email": "email@zitadel.com"
					}`),
				), user.HumanAddedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanAdded,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.users8 (id, creation_date, change_date, resource_owner, instance_id, state, sequence, username, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								domain.UserStateActive,
								uint64(15),
								"user-name",
								domain.UserTypeHuman,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.users8_humans (user_id, instance_id, first_name, last_name, nick_name, display_name, preferred_language, gender, email, phone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								"first-name",
								"last-name",
								&sql.NullString{},
								&sql.NullString{},
								&sql.NullString{String: "und", Valid: false},
								&sql.NullInt16{},
								"email@zitadel.com",
								&sql.NullString{},
							},
						},
						{
							expectedStmt: "INSERT INTO projections.users8_notifications (user_id, instance_id, last_email, last_phone, password_set) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								"email@zitadel.com",
								&sql.NullString{String: "", Valid: false},
								false,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceHumanRegistered",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanRegisteredType),
					user.AggregateType,
					[]byte(`{
						"username": "user-name",
						"firstName": "first-name",
						"lastName": "last-name",
						"nickName": "nick-name",
						"displayName": "display-name",
						"preferredLanguage": "ch-DE",
						"gender": 1,
						"email": "email@zitadel.com",
						"phone": "+41 00 000 00 00"
					}`),
				), user.HumanRegisteredEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanRegistered,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.users8 (id, creation_date, change_date, resource_owner, instance_id, state, sequence, username, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								domain.UserStateActive,
								uint64(15),
								"user-name",
								domain.UserTypeHuman,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.users8_humans (user_id, instance_id, first_name, last_name, nick_name, display_name, preferred_language, gender, email, phone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								"first-name",
								"last-name",
								&sql.NullString{String: "nick-name", Valid: true},
								&sql.NullString{String: "display-name", Valid: true},
								&sql.NullString{String: "ch-DE", Valid: true},
								&sql.NullInt16{Int16: int16(domain.GenderFemale), Valid: true},
								"email@zitadel.com",
								&sql.NullString{String: "+41 00 000 00 00", Valid: true},
							},
						},
						{
							expectedStmt: "INSERT INTO projections.users8_notifications (user_id, instance_id, last_email, last_phone, password_set) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								"email@zitadel.com",
								&sql.NullString{String: "+41 00 000 00 00", Valid: true},
								false,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserV1Registered",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserV1RegisteredType),
					user.AggregateType,
					[]byte(`{
						"username": "user-name",
						"firstName": "first-name",
						"lastName": "last-name",
						"nickName": "nick-name",
						"displayName": "display-name",
						"preferredLanguage": "ch-DE",
						"gender": 1,
						"email": "email@zitadel.com",
						"phone": "+41 00 000 00 00"
					}`),
				), user.HumanRegisteredEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanRegistered,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.users8 (id, creation_date, change_date, resource_owner, instance_id, state, sequence, username, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								domain.UserStateActive,
								uint64(15),
								"user-name",
								domain.UserTypeHuman,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.users8_humans (user_id, instance_id, first_name, last_name, nick_name, display_name, preferred_language, gender, email, phone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								"first-name",
								"last-name",
								&sql.NullString{String: "nick-name", Valid: true},
								&sql.NullString{String: "display-name", Valid: true},
								&sql.NullString{String: "ch-DE", Valid: true},
								&sql.NullInt16{Int16: int16(domain.GenderFemale), Valid: true},
								"email@zitadel.com",
								&sql.NullString{String: "+41 00 000 00 00", Valid: true},
							},
						},
						{
							expectedStmt: "INSERT INTO projections.users8_notifications (user_id, instance_id, last_email, last_phone, password_set) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								"email@zitadel.com",
								&sql.NullString{String: "+41 00 000 00 00", Valid: true},
								false,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceHumanRegistered NULLs",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanRegisteredType),
					user.AggregateType,
					[]byte(`{
						"username": "user-name",
						"firstName": "first-name",
						"lastName": "last-name",
						"email": "email@zitadel.com"
					}`),
				), user.HumanRegisteredEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanRegistered,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.users8 (id, creation_date, change_date, resource_owner, instance_id, state, sequence, username, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								domain.UserStateActive,
								uint64(15),
								"user-name",
								domain.UserTypeHuman,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.users8_humans (user_id, instance_id, first_name, last_name, nick_name, display_name, preferred_language, gender, email, phone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								"first-name",
								"last-name",
								&sql.NullString{},
								&sql.NullString{},
								&sql.NullString{String: "und", Valid: false},
								&sql.NullInt16{},
								"email@zitadel.com",
								&sql.NullString{},
							},
						},
						{
							expectedStmt: "INSERT INTO projections.users8_notifications (user_id, instance_id, last_email, last_phone, password_set) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								"email@zitadel.com",
								&sql.NullString{String: "", Valid: false},
								false,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceHumanInitCodeAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanInitialCodeAddedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.HumanInitialCodeAddedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanInitCodeAdded,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET state = $1 WHERE (id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								domain.UserStateInitial,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserV1InitCodeAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserV1InitialCodeAddedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.HumanInitialCodeAddedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanInitCodeAdded,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET state = $1 WHERE (id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								domain.UserStateInitial,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceHumanInitCodeSucceeded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanInitializedCheckSucceededType),
					user.AggregateType,
					[]byte(`{}`),
				), user.HumanInitializedCheckSucceededEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanInitCodeSucceeded,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET state = $1 WHERE (id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								domain.UserStateActive,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserV1InitCodeAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserV1InitializedCheckSucceededType),
					user.AggregateType,
					[]byte(`{}`),
				), user.HumanInitializedCheckSucceededEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanInitCodeSucceeded,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET state = $1 WHERE (id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								domain.UserStateActive,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserLocked",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserLockedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.UserLockedEventMapper),
			},
			reduce: (&userProjection{}).reduceUserLocked,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, state, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								domain.UserStateLocked,
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserUnlocked",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserUnlockedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.UserUnlockedEventMapper),
			},
			reduce: (&userProjection{}).reduceUserUnlocked,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, state, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								domain.UserStateActive,
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserDeactivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserDeactivatedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.UserDeactivatedEventMapper),
			},
			reduce: (&userProjection{}).reduceUserDeactivated,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, state, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								domain.UserStateInactive,
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserReactivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserReactivatedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.UserReactivatedEventMapper),
			},
			reduce: (&userProjection{}).reduceUserReactivated,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, state, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								domain.UserStateActive,
								uint64(15),
								"agg-id",
								"instance-id",
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
					repository.EventType(user.UserRemovedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.UserRemovedEventMapper),
			},
			reduce: (&userProjection{}).reduceUserRemoved,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.users8 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserUserNameChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserUserNameChangedType),
					user.AggregateType,
					[]byte(`{
						"username": "username"
					}`),
				), user.UsernameChangedEventMapper),
			},
			reduce: (&userProjection{}).reduceUserNameChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, username, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								"username",
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceDomainClaimed",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserDomainClaimedType),
					user.AggregateType,
					[]byte(`{
						"username": "id@temporary.domain"
					}`),
				), user.DomainClaimedEventMapper),
			},
			reduce: (&userProjection{}).reduceDomainClaimed,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, username, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								"id@temporary.domain",
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceHumanProfileChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanProfileChangedType),
					user.AggregateType,
					[]byte(`{
						"firstName": "first-name",
						"lastName": "last-name",
						"nickName": "nick-name",
						"displayName": "display-name",
						"preferredLanguage": "ch-DE",
						"gender": 3
					}`),
				), user.HumanProfileChangedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanProfileChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_humans SET (first_name, last_name, nick_name, display_name, preferred_language, gender) = ($1, $2, $3, $4, $5, $6) WHERE (user_id = $7) AND (instance_id = $8)",
							expectedArgs: []interface{}{
								"first-name",
								"last-name",
								"nick-name",
								"display-name",
								"ch-DE",
								domain.GenderDiverse,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserV1ProfileChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserV1ProfileChangedType),
					user.AggregateType,
					[]byte(`{
						"firstName": "first-name",
						"lastName": "last-name",
						"nickName": "nick-name",
						"displayName": "display-name",
						"preferredLanguage": "ch-DE",
						"gender": 3
					}`),
				), user.HumanProfileChangedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanProfileChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_humans SET (first_name, last_name, nick_name, display_name, preferred_language, gender) = ($1, $2, $3, $4, $5, $6) WHERE (user_id = $7) AND (instance_id = $8)",
							expectedArgs: []interface{}{
								"first-name",
								"last-name",
								"nick-name",
								"display-name",
								"ch-DE",
								domain.GenderDiverse,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceHumanPhoneChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanPhoneChangedType),
					user.AggregateType,
					[]byte(`{
						"phone": "+41 00 000 00 00"
						}`),
				), user.HumanPhoneChangedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanPhoneChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_humans SET (phone, is_phone_verified) = ($1, $2) WHERE (user_id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								"+41 00 000 00 00",
								false,
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_notifications SET last_phone = $1 WHERE (user_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								&sql.NullString{String: "+41 00 000 00 00", Valid: true},
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserV1PhoneChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserV1PhoneChangedType),
					user.AggregateType,
					[]byte(`{
						"phone": "+41 00 000 00 00"
						}`),
				), user.HumanPhoneChangedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanPhoneChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_humans SET (phone, is_phone_verified) = ($1, $2) WHERE (user_id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								"+41 00 000 00 00",
								false,
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_notifications SET last_phone = $1 WHERE (user_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								&sql.NullString{String: "+41 00 000 00 00", Valid: true},
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceHumanPhoneRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanPhoneRemovedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.HumanPhoneRemovedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanPhoneRemoved,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_humans SET (phone, is_phone_verified) = ($1, $2) WHERE (user_id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								nil,
								nil,
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_notifications SET (last_phone, verified_phone) = ($1, $2) WHERE (user_id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								nil,
								nil,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserV1PhoneRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserV1PhoneRemovedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.HumanPhoneRemovedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanPhoneRemoved,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_humans SET (phone, is_phone_verified) = ($1, $2) WHERE (user_id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								nil,
								nil,
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_notifications SET (last_phone, verified_phone) = ($1, $2) WHERE (user_id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								nil,
								nil,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceHumanPhoneVerified",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanPhoneVerifiedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.HumanPhoneVerifiedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanPhoneVerified,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_humans SET is_phone_verified = $1 WHERE (user_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								true,
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_notifications SET verified_phone = last_phone WHERE (user_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserV1PhoneVerified",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserV1PhoneVerifiedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.HumanPhoneVerifiedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanPhoneVerified,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_humans SET is_phone_verified = $1 WHERE (user_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								true,
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_notifications SET verified_phone = last_phone WHERE (user_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceHumanEmailChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanEmailChangedType),
					user.AggregateType,
					[]byte(`{
						"email": "email@zitadel.com"
					}`),
				), user.HumanEmailChangedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanEmailChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_humans SET (email, is_email_verified) = ($1, $2) WHERE (user_id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								"email@zitadel.com",
								false,
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_notifications SET last_email = $1 WHERE (user_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								&sql.NullString{String: "email@zitadel.com", Valid: true},
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserV1EmailChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserV1EmailChangedType),
					user.AggregateType,
					[]byte(`{
						"email": "email@zitadel.com"
					}`),
				), user.HumanEmailChangedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanEmailChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_humans SET (email, is_email_verified) = ($1, $2) WHERE (user_id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								"email@zitadel.com",
								false,
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_notifications SET last_email = $1 WHERE (user_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								&sql.NullString{String: "email@zitadel.com", Valid: true},
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceHumanEmailVerified",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanEmailVerifiedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.HumanEmailVerifiedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanEmailVerified,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_humans SET is_email_verified = $1 WHERE (user_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								true,
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_notifications SET verified_email = last_email WHERE (user_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserV1EmailVerified",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserV1EmailVerifiedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.HumanEmailVerifiedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanEmailVerified,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_humans SET is_email_verified = $1 WHERE (user_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								true,
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_notifications SET verified_email = last_email WHERE (user_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceHumanAvatarAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanAvatarAddedType),
					user.AggregateType,
					[]byte(`{
						"storeKey": "users/agg-id/avatar"
					}`),
				), user.HumanAvatarAddedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanAvatarAdded,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_humans SET avatar_key = $1 WHERE (user_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"users/agg-id/avatar",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceHumanAvatarRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.HumanAvatarRemovedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.HumanAvatarRemovedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanAvatarRemoved,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_humans SET avatar_key = $1 WHERE (user_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								nil,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceMachineAddedEvent no description",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.MachineAddedEventType),
					user.AggregateType,
					[]byte(`{
						"username": "username",
						"name": "machine-name"
					}`),
				), user.MachineAddedEventMapper),
			},
			reduce: (&userProjection{}).reduceMachineAdded,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.users8 (id, creation_date, change_date, resource_owner, instance_id, state, sequence, username, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								domain.UserStateActive,
								uint64(15),
								"username",
								domain.UserTypeMachine,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.users8_machines (user_id, instance_id, name, description, access_token_type) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								"machine-name",
								&sql.NullString{},
								domain.OIDCTokenTypeBearer,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceMachineAddedEvent",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.MachineAddedEventType),
					user.AggregateType,
					[]byte(`{
						"username": "username",
						"name": "machine-name",
						"description": "description"
					}`),
				), user.MachineAddedEventMapper),
			},
			reduce: (&userProjection{}).reduceMachineAdded,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.users8 (id, creation_date, change_date, resource_owner, instance_id, state, sequence, username, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								domain.UserStateActive,
								uint64(15),
								"username",
								domain.UserTypeMachine,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.users8_machines (user_id, instance_id, name, description, access_token_type) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								"machine-name",
								&sql.NullString{String: "description", Valid: true},
								domain.OIDCTokenTypeBearer,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceMachineChangedEvent",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.MachineChangedEventType),
					user.AggregateType,
					[]byte(`{
						"name": "machine-name",
						"description": "description"
					}`),
				), user.MachineChangedEventMapper),
			},
			reduce: (&userProjection{}).reduceMachineChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_machines SET (name, description) = ($1, $2) WHERE (user_id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								"machine-name",
								"description",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceMachineChangedEvent name",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.MachineChangedEventType),
					user.AggregateType,
					[]byte(`{
						"name": "machine-name"
					}`),
				), user.MachineChangedEventMapper),
			},
			reduce: (&userProjection{}).reduceMachineChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_machines SET name = $1 WHERE (user_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"machine-name",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceMachineChangedEvent description",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.MachineChangedEventType),
					user.AggregateType,
					[]byte(`{
						"description": "description"
					}`),
				), user.MachineChangedEventMapper),
			},
			reduce: (&userProjection{}).reduceMachineChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_machines SET description = $1 WHERE (user_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"description",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceMachineChangedEvent no values",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.MachineChangedEventType),
					user.AggregateType,
					[]byte(`{}`),
				), user.MachineChangedEventMapper),
			},
			reduce: (&userProjection{}).reduceMachineChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{},
				},
			},
		},
		{
			name: "reduceMachineSecretSet",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.MachineSecretSetType),
					user.AggregateType,
					[]byte(`{
						"client_secret": {}
					}`),
				), user.MachineSecretSetEventMapper),
			},
			reduce: (&userProjection{}).reduceMachineSecretSet,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_machines SET has_secret = $1 WHERE (user_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								true,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceMachineSecretSet",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.MachineSecretRemovedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.MachineSecretRemovedEventMapper),
			},
			reduce: (&userProjection{}).reduceMachineSecretRemoved,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.users8_machines SET has_secret = $1 WHERE (user_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								false,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceOwnerRemoved",
			reduce: (&userProjection{}).reduceOwnerRemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgRemovedEventType),
					org.AggregateType,
					nil,
				), org.OrgRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.users8 SET (change_date, sequence, owner_removed) = ($1, $2, $3) WHERE (instance_id = $4) AND (resource_owner = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
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
				event: getEvent(testEvent(
					repository.EventType(instance.InstanceRemovedEventType),
					instance.AggregateType,
					nil,
				), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(UserInstanceIDCol),
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.users8 WHERE (instance_id = $1)",
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
			assertReduce(t, got, err, UserTable, tt.want)
		})
	}
}
