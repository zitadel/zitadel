package model

import (
	"encoding/json"
	"testing"

	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
)

func TestAppendMfaOTPAddedEvent(t *testing.T) {
	type args struct {
		user  *Human
		otp   *OTP
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Human
	}{
		{
			name: "append user otp event",
			args: args{
				user:  &Human{},
				otp:   &OTP{Secret: &crypto.CryptoValue{KeyID: "KeyID"}},
				event: &es_models.Event{},
			},
			result: &Human{OTP: &OTP{Secret: &crypto.CryptoValue{KeyID: "KeyID"}, State: int32(model.MfaStateNotReady)}},
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
		user  *Human
		otp   *OTP
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Human
	}{
		{
			name: "append otp verify event",
			args: args{
				user:  &Human{OTP: &OTP{Secret: &crypto.CryptoValue{KeyID: "KeyID"}}},
				otp:   &OTP{Secret: &crypto.CryptoValue{KeyID: "KeyID"}},
				event: &es_models.Event{},
			},
			result: &Human{OTP: &OTP{Secret: &crypto.CryptoValue{KeyID: "KeyID"}, State: int32(model.MfaStateReady)}},
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
		user  *Human
		otp   *OTP
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Human
	}{
		{
			name: "append otp verify event",
			args: args{
				user:  &Human{OTP: &OTP{Secret: &crypto.CryptoValue{KeyID: "KeyID"}}},
				event: &es_models.Event{},
			},
			result: &Human{},
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
