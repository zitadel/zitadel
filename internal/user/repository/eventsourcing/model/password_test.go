package model

import (
	"encoding/json"
	"testing"
	"time"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

func TestAppendUserPasswordChangedEvent(t *testing.T) {
	type args struct {
		user  *Human
		pw    *Password
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Human
	}{
		{
			name: "append init user code event",
			args: args{
				user:  &Human{},
				pw:    &Password{ChangeRequired: true},
				event: &es_models.Event{},
			},
			result: &Human{Password: &Password{ChangeRequired: true}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.pw != nil {
				data, _ := json.Marshal(tt.args.pw)
				tt.args.event.Data = data
			}
			tt.args.user.appendUserPasswordChangedEvent(tt.args.event)
			if tt.args.user.Password.ChangeRequired != tt.result.Password.ChangeRequired {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.user)
			}
		})
	}
}

func TestAppendPasswordSetRequestedEvent(t *testing.T) {
	type args struct {
		user  *Human
		code  *PasswordCode
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Human
	}{
		{
			name: "append user phone code added event",
			args: args{
				user:  &Human{Phone: &Phone{PhoneNumber: "PhoneNumber"}},
				code:  &PasswordCode{Expiry: time.Hour * 1},
				event: &es_models.Event{},
			},
			result: &Human{PasswordCode: &PasswordCode{Expiry: time.Hour * 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.code != nil {
				data, _ := json.Marshal(tt.args.code)
				tt.args.event.Data = data
			}
			tt.args.user.appendPasswordSetRequestedEvent(tt.args.event)
			if tt.args.user.PasswordCode.Expiry != tt.result.PasswordCode.Expiry {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.user)
			}
		})
	}
}
