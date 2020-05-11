package model

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	"testing"
)

func TestAppendMfaOTPAddedEvent(t *testing.T) {
	type args struct {
		user  *User
		otp   *OTP
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *User
	}{
		{
			name: "append user otp event",
			args: args{
				user:  &User{},
				otp:   &OTP{Secret: &crypto.CryptoValue{KeyID: "KeyID"}},
				event: &es_models.Event{},
			},
			result: &User{OTP: &OTP{Secret: &crypto.CryptoValue{KeyID: "KeyID"}, State: int32(model.MFASTATE_NOTREADY)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.otp != nil {
				data, _ := json.Marshal(tt.args.otp)
				tt.args.event.Data = data
			}
			tt.args.user.appendOtpAddedEvent(tt.args.event)
			if tt.args.user.OTP.State != tt.result.OTP.State {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.OTP.State, tt.args.user.OTP.State)
			}
		})
	}
}

func TestAppendMfaOTPVerifyEvent(t *testing.T) {
	type args struct {
		user  *User
		otp   *OTP
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *User
	}{
		{
			name: "append otp verify event",
			args: args{
				user:  &User{OTP: &OTP{Secret: &crypto.CryptoValue{KeyID: "KeyID"}}},
				otp:   &OTP{Secret: &crypto.CryptoValue{KeyID: "KeyID"}},
				event: &es_models.Event{},
			},
			result: &User{OTP: &OTP{Secret: &crypto.CryptoValue{KeyID: "KeyID"}, State: int32(model.MFASTATE_READY)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.otp != nil {
				data, _ := json.Marshal(tt.args.otp)
				tt.args.event.Data = data
			}
			tt.args.user.appendOtpVerifiedEvent()
			if tt.args.user.OTP.State != tt.result.OTP.State {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.OTP.State, tt.args.user.OTP.State)
			}
		})
	}
}

func TestAppendMfaOTPRemoveEvent(t *testing.T) {
	type args struct {
		user  *User
		otp   *OTP
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *User
	}{
		{
			name: "append otp verify event",
			args: args{
				user:  &User{OTP: &OTP{Secret: &crypto.CryptoValue{KeyID: "KeyID"}}},
				event: &es_models.Event{},
			},
			result: &User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.user.appendOtpRemovedEvent()
			if tt.args.user.OTP != nil {
				t.Errorf("got wrong result: actual: %v ", tt.result.OTP)
			}
		})
	}
}
