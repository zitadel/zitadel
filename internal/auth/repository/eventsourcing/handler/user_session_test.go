package handler

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/user"
	view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

// recordingExecuter records every statement a reduced event executes.
type recordingExecuter struct {
	stmts []string
	args  [][]any
}

func (e *recordingExecuter) Exec(stmt string, args ...any) (sql.Result, error) {
	if stmt == "SAVEPOINT exec_stmt" {
		return nil, nil
	}
	e.stmts = append(e.stmts, stmt)
	e.args = append(e.args, args)
	return nil, nil
}

type staticIDGenerator string

func (g staticIDGenerator) Next() (string, error) { return string(g), nil }

func userSessionTestEvent(t *testing.T, typ eventstore.EventType, createdAt time.Time, mapper func(eventstore.Event) (eventstore.Event, error)) eventstore.Event {
	t.Helper()
	event, err := mapper(&repository.Event{
		Seq:           15,
		CreationDate:  createdAt,
		Typ:           typ,
		AggregateType: user.AggregateType,
		Data:          []byte(`{"id":"auth-request-id","userAgentID":"agent-id"}`),
		Version:       user.AggregateVersion,
		AggregateID:   "user-id",
		ResourceOwner: sql.NullString{String: "ro-id", Valid: true},
		InstanceID:    "instance-id",
		ID:            "event-id",
		EditorUser:    "editor-user",
	})
	require.NoError(t, err)
	return event
}

// The login v1 user session only skips the second factor while
// second_factor_verification is within the second factor check lifetime, so
// OTP SMS and OTP Email checks must be written like TOTP checks.
func TestUserSession_Reduce_OTPSMSAndEmailChecks(t *testing.T) {
	createdAt := time.Date(2026, 11, 24, 7, 33, 27, 0, time.UTC)
	tests := []struct {
		name             string
		eventType        eventstore.EventType
		mapper           func(eventstore.Event) (eventstore.Event, error)
		wantVerification time.Time
		wantType         *domain.MFAType
	}{
		{
			name:             "otp sms check succeeded",
			eventType:        user.HumanOTPSMSCheckSucceededType,
			mapper:           eventstore.GenericEventMapper[user.HumanOTPSMSCheckSucceededEvent],
			wantVerification: createdAt,
			wantType:         mfaTypePtr(domain.MFATypeOTPSMS),
		},
		{
			name:             "otp email check succeeded",
			eventType:        user.HumanOTPEmailCheckSucceededType,
			mapper:           eventstore.GenericEventMapper[user.HumanOTPEmailCheckSucceededEvent],
			wantVerification: createdAt,
			wantType:         mfaTypePtr(domain.MFATypeOTPEmail),
		},
		{
			name:             "otp sms check failed",
			eventType:        user.HumanOTPSMSCheckFailedType,
			mapper:           eventstore.GenericEventMapper[user.HumanOTPSMSCheckFailedEvent],
			wantVerification: time.Time{},
		},
		{
			name:             "otp email check failed",
			eventType:        user.HumanOTPEmailCheckFailedType,
			mapper:           eventstore.GenericEventMapper[user.HumanOTPEmailCheckFailedEvent],
			wantVerification: time.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserSession{idGenerator: staticIDGenerator("session-id")}

			stmt, err := u.Reduce(userSessionTestEvent(t, tt.eventType, createdAt, tt.mapper))
			require.NoError(t, err)
			require.NotNil(t, stmt.Execute, "event must not be reduced to a no-op")

			ex := new(recordingExecuter)
			require.NoError(t, stmt.Execute(t.Context(), ex, userSessionTable))
			require.Len(t, ex.stmts, 1)
			assert.Contains(t, ex.stmts[0], view_model.UserSessionKeySecondFactorVerification)
			assert.Contains(t, ex.args[0], "agent-id", "the row of the user agent that made the check")
			assert.Contains(t, ex.args[0], "user-id")
			assert.Contains(t, ex.args[0], tt.wantVerification)
			if tt.wantType != nil {
				assert.Contains(t, ex.stmts[0], view_model.UserSessionKeySecondFactorVerificationType)
				assert.Contains(t, ex.args[0], *tt.wantType)
			}
		})
	}
}

func TestUserSession_Reducers_OTPSMSAndEmailChecks(t *testing.T) {
	var subscribed []eventstore.EventType
	for _, aggregate := range (&UserSession{}).Reducers() {
		if aggregate.Aggregate != user.AggregateType {
			continue
		}
		for _, reducer := range aggregate.EventReducers {
			subscribed = append(subscribed, reducer.Event)
		}
	}
	for _, eventType := range []eventstore.EventType{
		user.HumanOTPSMSCheckSucceededType,
		user.HumanOTPSMSCheckFailedType,
		user.HumanOTPEmailCheckSucceededType,
		user.HumanOTPEmailCheckFailedType,
	} {
		assert.Contains(t, subscribed, eventType)
	}
}

func mfaTypePtr(t domain.MFAType) *domain.MFAType { return &t }
