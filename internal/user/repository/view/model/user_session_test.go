package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/crypto"
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
			name: "append user password check succeeded event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.UserPasswordCheckSucceeded},
				userView: &UserSessionView{},
			},
			result: &UserSessionView{ChangeDate: now(), PasswordVerification: now()},
		},
		{
			name: "append human password check succeeded event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.HumanPasswordCheckSucceeded},
				userView: &UserSessionView{},
			},
			result: &UserSessionView{ChangeDate: now(), PasswordVerification: now()},
		},
		{
			name: "append user password check failed event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.UserPasswordCheckFailed},
				userView: &UserSessionView{PasswordVerification: now()},
			},
			result: &UserSessionView{ChangeDate: now(), PasswordVerification: time.Time{}},
		},
		{
			name: "append human password check failed event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.HumanPasswordCheckFailed},
				userView: &UserSessionView{PasswordVerification: now()},
			},
			result: &UserSessionView{ChangeDate: now(), PasswordVerification: time.Time{}},
		},
		{
			name: "append user password changed event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.UserPasswordChanged},
				userView: &UserSessionView{PasswordVerification: now()},
			},
			result: &UserSessionView{ChangeDate: now(), PasswordVerification: time.Time{}},
		},
		{
			name: "append human password changed event",
			args: args{
				event: &es_models.Event{
					CreationDate: now(),
					Type:         es_model.HumanPasswordChanged,
					Data: func() []byte {
						d, _ := json.Marshal(&es_model.PasswordChange{
							Password: &es_model.Password{
								Secret: &crypto.CryptoValue{Crypted: []byte("test")},
							},
						})
						return d
					}(),
				},
				userView: &UserSessionView{PasswordVerification: now(), UserAgentID: "id"},
			},
			result: &UserSessionView{ChangeDate: now(), UserAgentID: "id", PasswordVerification: time.Time{}},
		},
		{
			name: "append human password changed event same user agent",
			args: args{
				event: &es_models.Event{
					AggregateID:  "id",
					CreationDate: now(),
					Type:         es_model.HumanPasswordChanged,
					Data: func() []byte {
						d, _ := json.Marshal(&es_model.PasswordChange{
							Password: &es_model.Password{
								Secret: &crypto.CryptoValue{Crypted: []byte("test")},
							},
							UserAgentID: "id",
						})
						return d
					}(),
				},
				userView: &UserSessionView{PasswordVerification: now(), UserAgentID: "id"},
			},
			result: &UserSessionView{ChangeDate: now(), UserAgentID: "id", PasswordVerification: now()},
		},
		{
			name: "append user otp check succeeded event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.MFAOTPCheckSucceeded},
				userView: &UserSessionView{},
			},
			result: &UserSessionView{ChangeDate: now(), SecondFactorVerification: now()},
		},
		{
			name: "append human otp check succeeded event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.HumanMFAOTPCheckSucceeded},
				userView: &UserSessionView{},
			},
			result: &UserSessionView{ChangeDate: now(), SecondFactorVerification: now()},
		},
		{
			name: "append user otp check failed event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.MFAOTPCheckFailed},
				userView: &UserSessionView{SecondFactorVerification: now()},
			},
			result: &UserSessionView{ChangeDate: now(), SecondFactorVerification: time.Time{}},
		},
		{
			name: "append human otp check failed event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.HumanMFAOTPCheckFailed},
				userView: &UserSessionView{SecondFactorVerification: now()},
			},
			result: &UserSessionView{ChangeDate: now(), SecondFactorVerification: time.Time{}},
		},
		{
			name: "append user otp removed event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.MFAOTPRemoved},
				userView: &UserSessionView{SecondFactorVerification: now()},
			},
			result: &UserSessionView{ChangeDate: now(), SecondFactorVerification: time.Time{}},
		},
		{
			name: "append human otp removed event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.HumanMFAOTPRemoved},
				userView: &UserSessionView{SecondFactorVerification: now()},
			},
			result: &UserSessionView{ChangeDate: now(), SecondFactorVerification: time.Time{}},
		},
		{
			name: "append user signed out event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.SignedOut},
				userView: &UserSessionView{PasswordVerification: now(), SecondFactorVerification: now()},
			},
			result: &UserSessionView{ChangeDate: now(), PasswordVerification: time.Time{}, SecondFactorVerification: time.Time{}, State: 1},
		},
		{
			name: "append human signed out event",
			args: args{
				event:    &es_models.Event{CreationDate: now(), Type: es_model.HumanSignedOut},
				userView: &UserSessionView{PasswordVerification: now(), SecondFactorVerification: now()},
			},
			result: &UserSessionView{ChangeDate: now(), PasswordVerification: time.Time{}, SecondFactorVerification: time.Time{}, State: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.userView.AppendEvent(tt.args.event)
			assert.Equal(t, tt.result, tt.args.userView)
		})
	}
}
