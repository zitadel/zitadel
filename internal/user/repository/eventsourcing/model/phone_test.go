package model

import (
	"encoding/json"
	"testing"
	"time"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

func TestPhoneChanges(t *testing.T) {
	type args struct {
		existingPhone *Phone
		newPhone      *Phone
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
				existingPhone: &Phone{PhoneNumber: "Phone", IsPhoneVerified: true},
				newPhone:      &Phone{PhoneNumber: "PhoneChanged", IsPhoneVerified: false},
			},
			res: res{
				changesLen: 1,
			},
		},
		{
			name: "no fields changed",
			args: args{
				existingPhone: &Phone{PhoneNumber: "Phone", IsPhoneVerified: true},
				newPhone:      &Phone{PhoneNumber: "Phone", IsPhoneVerified: false},
			},
			res: res{
				changesLen: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := tt.args.existingPhone.Changes(tt.args.newPhone)
			if len(changes) != tt.res.changesLen {
				t.Errorf("got wrong changes len: expected: %v, actual: %v ", tt.res.changesLen, len(changes))
			}
		})
	}
}

func TestAppendUserPhoneChangedEvent(t *testing.T) {
	type args struct {
		user  *Human
		phone *Phone
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Human
	}{
		{
			name: "append user phone event",
			args: args{
				user:  &Human{Phone: &Phone{PhoneNumber: "PhoneNumber"}},
				phone: &Phone{PhoneNumber: "PhoneNumberChanged"},
				event: &es_models.Event{},
			},
			result: &Human{Phone: &Phone{PhoneNumber: "PhoneNumberChanged"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.phone != nil {
				data, _ := json.Marshal(tt.args.phone)
				tt.args.event.Data = data
			}
			tt.args.user.appendUserPhoneChangedEvent(tt.args.event)
			if tt.args.user.Phone.PhoneNumber != tt.result.Phone.PhoneNumber {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.user)
			}
		})
	}
}

func TestAppendUserPhoneCodeAddedEvent(t *testing.T) {
	type args struct {
		user  *Human
		code  *PhoneCode
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
				code:  &PhoneCode{Expiry: time.Hour * 1},
				event: &es_models.Event{},
			},
			result: &Human{PhoneCode: &PhoneCode{Expiry: time.Hour * 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.code != nil {
				data, _ := json.Marshal(tt.args.code)
				tt.args.event.Data = data
			}
			tt.args.user.appendUserPhoneCodeAddedEvent(tt.args.event)
			if tt.args.user.PhoneCode.Expiry != tt.result.PhoneCode.Expiry {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.user)
			}
		})
	}
}

func TestAppendUserPhoneVerifiedEvent(t *testing.T) {
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
			name: "append user phone event",
			args: args{
				user:  &Human{Phone: &Phone{PhoneNumber: "PhoneNumber"}},
				event: &es_models.Event{},
			},
			result: &Human{Phone: &Phone{PhoneNumber: "PhoneNumber", IsPhoneVerified: true}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.args.user.appendUserPhoneVerifiedEvent()
			if tt.args.user.Phone.IsPhoneVerified != tt.result.Phone.IsPhoneVerified {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.user)
			}
		})
	}
}
