package model

import (
	"encoding/json"
	"testing"

	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/user/model"
)

func TestAppendMFAOTPAddedEvent(t *testing.T) {
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
			result: &Human{OTP: &OTP{Secret: &crypto.CryptoValue{KeyID: "KeyID"}, State: int32(model.MFAStateNotReady)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.otp != nil {
				data, _ := json.Marshal(tt.args.otp)
				tt.args.event.Data = data
			}
			tt.args.user.appendOTPAddedEvent(tt.args.event)
			if tt.args.user.OTP.State != tt.result.OTP.State {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.OTP.State, tt.args.user.OTP.State)
			}
		})
	}
}

func TestAppendMFAOTPVerifyEvent(t *testing.T) {
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
			result: &Human{OTP: &OTP{Secret: &crypto.CryptoValue{KeyID: "KeyID"}, State: int32(model.MFAStateReady)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.otp != nil {
				data, _ := json.Marshal(tt.args.otp)
				tt.args.event.Data = data
			}
			tt.args.user.appendOTPVerifiedEvent()
			if tt.args.user.OTP.State != tt.result.OTP.State {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.OTP.State, tt.args.user.OTP.State)
			}
		})
	}
}

func TestAppendMFAOTPRemoveEvent(t *testing.T) {
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
			tt.args.user.appendOTPRemovedEvent()
			if tt.args.user.OTP != nil {
				t.Errorf("got wrong result: actual: %v ", tt.result.OTP)
			}
		})
	}
}
