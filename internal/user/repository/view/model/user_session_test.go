package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

func now() time.Time {
	return time.Now().UTC().Round(1 * time.Second)
}

func TestAppendEvent(t *testing.T) {
	type args struct {
		event    *es_models.Event
		userView *UserSessionView
	}
	tests := []struct {
		name   string
		args   args
		result *UserSessionView
	}{
		{
			name: "append password check succeeded event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.UserPasswordCheckSucceeded},
				userView: &UserSessionView{},
			},
			result: &UserSessionView{ChangeDate: now(), PasswordVerification: now()},
		},
		{
			name: "append password check failed event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.UserPasswordCheckFailed},
				userView: &UserSessionView{PasswordVerification: now()},
			},
			result: &UserSessionView{ChangeDate: now(), PasswordVerification: time.Time{}},
		},
		{
			name: "append password changed event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.UserPasswordChanged},
				userView: &UserSessionView{PasswordVerification: now()},
			},
			result: &UserSessionView{ChangeDate: now(), PasswordVerification: time.Time{}},
		},
		{
			name: "append otp check succeeded event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.MfaOtpCheckSucceeded},
				userView: &UserSessionView{},
			},
			result: &UserSessionView{ChangeDate: now(), MfaSoftwareVerification: now()},
		},
		{
			name: "append otp check failed event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.MfaOtpCheckFailed},
				userView: &UserSessionView{MfaSoftwareVerification: now()},
			},
			result: &UserSessionView{ChangeDate: now(), MfaSoftwareVerification: time.Time{}},
		},
		{
			name: "append otp removed event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.MfaOtpCheckFailed},
				userView: &UserSessionView{MfaSoftwareVerification: now()},
			},
			result: &UserSessionView{ChangeDate: now(), MfaSoftwareVerification: time.Time{}},
		},
		{
			name: "append otp removed event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.SignedOut},
				userView: &UserSessionView{PasswordVerification: now(), MfaSoftwareVerification: now()},
			},
			result: &UserSessionView{ChangeDate: now(), PasswordVerification: time.Time{}, MfaSoftwareVerification: time.Time{}, State: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.userView.AppendEvent(tt.args.event)
			assert.Equal(t, tt.result, tt.args.userView)
		})
	}
}
