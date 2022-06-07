package projection

import (
	"database/sql"
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
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
						"email": "email@zitadel.ch",
						"phone": "+41 00 000 00 00"
					}`),
				), user.HumanAddedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanAdded,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.users (id, creation_date, change_date, resource_owner, state, sequence, username, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								domain.UserStateActive,
								uint64(15),
								"user-name",
								domain.UserTypeHuman,
							},
						},
						{
							expectedStmt: "INSERT INTO zitadel.projections.users_humans (user_id, first_name, last_name, nick_name, display_name, preferred_language, gender, email, phone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								"first-name",
								"last-name",
								&sql.NullString{String: "nick-name", Valid: true},
								&sql.NullString{String: "display-name", Valid: true},
								&sql.NullString{String: "ch-DE", Valid: true},
								&sql.NullInt16{Int16: int16(domain.GenderFemale), Valid: true},
								"email@zitadel.ch",
								&sql.NullString{String: "+41 00 000 00 00", Valid: true},
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
						"email": "email@zitadel.ch",
						"phone": "+41 00 000 00 00"
					}`),
				), user.HumanAddedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanAdded,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.users (id, creation_date, change_date, resource_owner, state, sequence, username, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								domain.UserStateActive,
								uint64(15),
								"user-name",
								domain.UserTypeHuman,
							},
						},
						{
							expectedStmt: "INSERT INTO zitadel.projections.users_humans (user_id, first_name, last_name, nick_name, display_name, preferred_language, gender, email, phone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								"first-name",
								"last-name",
								&sql.NullString{String: "nick-name", Valid: true},
								&sql.NullString{String: "display-name", Valid: true},
								&sql.NullString{String: "ch-DE", Valid: true},
								&sql.NullInt16{Int16: int16(domain.GenderFemale), Valid: true},
								"email@zitadel.ch",
								&sql.NullString{String: "+41 00 000 00 00", Valid: true},
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
						"email": "email@zitadel.ch"
					}`),
				), user.HumanAddedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanAdded,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.users (id, creation_date, change_date, resource_owner, state, sequence, username, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								domain.UserStateActive,
								uint64(15),
								"user-name",
								domain.UserTypeHuman,
							},
						},
						{
							expectedStmt: "INSERT INTO zitadel.projections.users_humans (user_id, first_name, last_name, nick_name, display_name, preferred_language, gender, email, phone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								"first-name",
								"last-name",
								&sql.NullString{},
								&sql.NullString{},
								&sql.NullString{String: "und", Valid: false},
								&sql.NullInt16{},
								"email@zitadel.ch",
								&sql.NullString{},
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
						"email": "email@zitadel.ch",
						"phone": "+41 00 000 00 00"
					}`),
				), user.HumanRegisteredEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanRegistered,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.users (id, creation_date, change_date, resource_owner, state, sequence, username, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								domain.UserStateActive,
								uint64(15),
								"user-name",
								domain.UserTypeHuman,
							},
						},
						{
							expectedStmt: "INSERT INTO zitadel.projections.users_humans (user_id, first_name, last_name, nick_name, display_name, preferred_language, gender, email, phone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								"first-name",
								"last-name",
								&sql.NullString{String: "nick-name", Valid: true},
								&sql.NullString{String: "display-name", Valid: true},
								&sql.NullString{String: "ch-DE", Valid: true},
								&sql.NullInt16{Int16: int16(domain.GenderFemale), Valid: true},
								"email@zitadel.ch",
								&sql.NullString{String: "+41 00 000 00 00", Valid: true},
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
						"email": "email@zitadel.ch",
						"phone": "+41 00 000 00 00"
					}`),
				), user.HumanRegisteredEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanRegistered,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.users (id, creation_date, change_date, resource_owner, state, sequence, username, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								domain.UserStateActive,
								uint64(15),
								"user-name",
								domain.UserTypeHuman,
							},
						},
						{
							expectedStmt: "INSERT INTO zitadel.projections.users_humans (user_id, first_name, last_name, nick_name, display_name, preferred_language, gender, email, phone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								"first-name",
								"last-name",
								&sql.NullString{String: "nick-name", Valid: true},
								&sql.NullString{String: "display-name", Valid: true},
								&sql.NullString{String: "ch-DE", Valid: true},
								&sql.NullInt16{Int16: int16(domain.GenderFemale), Valid: true},
								"email@zitadel.ch",
								&sql.NullString{String: "+41 00 000 00 00", Valid: true},
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
						"email": "email@zitadel.ch"
					}`),
				), user.HumanRegisteredEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanRegistered,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.users (id, creation_date, change_date, resource_owner, state, sequence, username, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								domain.UserStateActive,
								uint64(15),
								"user-name",
								domain.UserTypeHuman,
							},
						},
						{
							expectedStmt: "INSERT INTO zitadel.projections.users_humans (user_id, first_name, last_name, nick_name, display_name, preferred_language, gender, email, phone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								"first-name",
								"last-name",
								&sql.NullString{},
								&sql.NullString{},
								&sql.NullString{String: "und", Valid: false},
								&sql.NullInt16{},
								"email@zitadel.ch",
								&sql.NullString{},
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (state) = ($1) WHERE (id = $2)",
							expectedArgs: []interface{}{
								domain.UserStateInitial,
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (state) = ($1) WHERE (id = $2)",
							expectedArgs: []interface{}{
								domain.UserStateInitial,
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (state) = ($1) WHERE (id = $2)",
							expectedArgs: []interface{}{
								domain.UserStateActive,
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (state) = ($1) WHERE (id = $2)",
							expectedArgs: []interface{}{
								domain.UserStateActive,
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, state, sequence) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								domain.UserStateLocked,
								uint64(15),
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, state, sequence) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								domain.UserStateActive,
								uint64(15),
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, state, sequence) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								domain.UserStateInactive,
								uint64(15),
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, state, sequence) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								domain.UserStateActive,
								uint64(15),
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.users WHERE (id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, username, sequence) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								"username",
								uint64(15),
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_humans SET (first_name, last_name, nick_name, display_name, preferred_language, gender) = ($1, $2, $3, $4, $5, $6) WHERE (user_id = $7)",
							expectedArgs: []interface{}{
								"first-name",
								"last-name",
								"nick-name",
								"display-name",
								"ch-DE",
								domain.GenderDiverse,
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_humans SET (first_name, last_name, nick_name, display_name, preferred_language, gender) = ($1, $2, $3, $4, $5, $6) WHERE (user_id = $7)",
							expectedArgs: []interface{}{
								"first-name",
								"last-name",
								"nick-name",
								"display-name",
								"ch-DE",
								domain.GenderDiverse,
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_humans SET (phone, is_phone_verified) = ($1, $2) WHERE (user_id = $3)",
							expectedArgs: []interface{}{
								"+41 00 000 00 00",
								false,
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_humans SET (phone, is_phone_verified) = ($1, $2) WHERE (user_id = $3)",
							expectedArgs: []interface{}{
								"+41 00 000 00 00",
								false,
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_humans SET (phone, is_phone_verified) = ($1, $2) WHERE (user_id = $3)",
							expectedArgs: []interface{}{
								nil,
								nil,
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_humans SET (phone, is_phone_verified) = ($1, $2) WHERE (user_id = $3)",
							expectedArgs: []interface{}{
								nil,
								nil,
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_humans SET (is_phone_verified) = ($1) WHERE (user_id = $2)",
							expectedArgs: []interface{}{
								true,
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_humans SET (is_phone_verified) = ($1) WHERE (user_id = $2)",
							expectedArgs: []interface{}{
								true,
								"agg-id",
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
						"email": "email@zitadel.ch"
					}`),
				), user.HumanEmailChangedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanEmailChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_humans SET (email, is_email_verified) = ($1, $2) WHERE (user_id = $3)",
							expectedArgs: []interface{}{
								"email@zitadel.ch",
								false,
								"agg-id",
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
						"email": "email@zitadel.ch"
					}`),
				), user.HumanEmailChangedEventMapper),
			},
			reduce: (&userProjection{}).reduceHumanEmailChanged,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_humans SET (email, is_email_verified) = ($1, $2) WHERE (user_id = $3)",
							expectedArgs: []interface{}{
								"email@zitadel.ch",
								false,
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_humans SET (is_email_verified) = ($1) WHERE (user_id = $2)",
							expectedArgs: []interface{}{
								true,
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_humans SET (is_email_verified) = ($1) WHERE (user_id = $2)",
							expectedArgs: []interface{}{
								true,
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_humans SET (avatar_key) = ($1) WHERE (user_id = $2)",
							expectedArgs: []interface{}{
								"users/agg-id/avatar",
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_humans SET (avatar_key) = ($1) WHERE (user_id = $2)",
							expectedArgs: []interface{}{
								nil,
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.users (id, creation_date, change_date, resource_owner, state, sequence, username, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								domain.UserStateActive,
								uint64(15),
								"username",
								domain.UserTypeMachine,
							},
						},
						{
							expectedStmt: "INSERT INTO zitadel.projections.users_machines (user_id, name, description) VALUES ($1, $2, $3)",
							expectedArgs: []interface{}{
								"agg-id",
								"machine-name",
								&sql.NullString{},
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.users (id, creation_date, change_date, resource_owner, state, sequence, username, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								domain.UserStateActive,
								uint64(15),
								"username",
								domain.UserTypeMachine,
							},
						},
						{
							expectedStmt: "INSERT INTO zitadel.projections.users_machines (user_id, name, description) VALUES ($1, $2, $3)",
							expectedArgs: []interface{}{
								"agg-id",
								"machine-name",
								&sql.NullString{String: "description", Valid: true},
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_machines SET (name, description) = ($1, $2) WHERE (user_id = $3)",
							expectedArgs: []interface{}{
								"machine-name",
								"description",
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_machines SET (name) = ($1) WHERE (user_id = $2)",
							expectedArgs: []interface{}{
								"machine-name",
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.users SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"agg-id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.users_machines SET (description) = ($1) WHERE (user_id = $2)",
							expectedArgs: []interface{}{
								"description",
								"agg-id",
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
				projection:       UserTable,
				executer: &testExecuter{
					executions: []execution{},
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
