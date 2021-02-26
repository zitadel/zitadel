package model

import (
	"encoding/json"
	"testing"
	"time"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

func TestEmailChanges(t *testing.T) {
	type args struct {
		existingEmail *Email
		new           *Email
	}
	type res struct {
		changesLen int
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "all fields changed",
			args: args{
				existingEmail: &Email{EmailAddress: "Email", IsEmailVerified: true},
				new:           &Email{EmailAddress: "EmailChanged", IsEmailVerified: false},
			},
			res: res{
				changesLen: 1,
			},
		},
		{
			name: "no fields changed",
			args: args{
				existingEmail: &Email{EmailAddress: "Email", IsEmailVerified: true},
				new:           &Email{EmailAddress: "Email", IsEmailVerified: false},
			},
			res: res{
				changesLen: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := tt.args.existingEmail.Changes(tt.args.new)
			if len(changes) != tt.res.changesLen {
				t.Errorf("got wrong changes len: expected: %v, actual: %v ", tt.res.changesLen, len(changes))
			}
		})
	}
}

func TestAppendUserEmailChangedEvent(t *testing.T) {
	type args struct {
		user  *Human
		email *Email
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Human
	}{
		{
			name: "append user email event",
			args: args{
				user:  &Human{Email: &Email{EmailAddress: "EmailAddress"}},
				email: &Email{EmailAddress: "EmailAddressChanged"},
				event: &es_models.Event{},
			},
			result: &Human{Email: &Email{EmailAddress: "EmailAddressChanged"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.email != nil {
				data, _ := json.Marshal(tt.args.email)
				tt.args.event.Data = data
			}
			tt.args.user.appendUserEmailChangedEvent(tt.args.event)
			if tt.args.user.Email.EmailAddress != tt.result.Email.EmailAddress {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.user)
			}
		})
	}
}

func TestAppendUserEmailCodeAddedEvent(t *testing.T) {
	type args struct {
		user  *Human
		code  *EmailCode
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Human
	}{
		{
			name: "append user email code added event",
			args: args{
				user:  &Human{Email: &Email{EmailAddress: "EmailAddress"}},
				code:  &EmailCode{Expiry: time.Hour * 1},
				event: &es_models.Event{},
			},
			result: &Human{EmailCode: &EmailCode{Expiry: time.Hour * 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.code != nil {
				data, _ := json.Marshal(tt.args.code)
				tt.args.event.Data = data
			}
			tt.args.user.appendUserEmailCodeAddedEvent(tt.args.event)
			if tt.args.user.EmailCode.Expiry != tt.result.EmailCode.Expiry {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.user)
			}
		})
	}
}

func TestAppendUserEmailVerifiedEvent(t *testing.T) {
	type args struct {
		user  *Human
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Human
	}{
		{
			name: "append user email event",
			args: args{
				user:  &Human{Email: &Email{EmailAddress: "EmailAddress"}},
				event: &es_models.Event{},
			},
			result: &Human{Email: &Email{EmailAddress: "EmailAddress", IsEmailVerified: true}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.args.user.appendUserEmailVerifiedEvent()
			if tt.args.user.Email.IsEmailVerified != tt.result.Email.IsEmailVerified {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.user)
			}
		})
	}
}
