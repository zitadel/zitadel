package model

import (
	"encoding/json"
	"testing"
	"time"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

func TestAppendInitUserCodeEvent(t *testing.T) {
	type args struct {
		user  *Human
		code  *InitUserCode
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
				code:  &InitUserCode{Expiry: time.Hour * 30},
				event: &es_models.Event{},
			},
			result: &Human{InitCode: &InitUserCode{Expiry: time.Hour * 30}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.code != nil {
				data, _ := json.Marshal(tt.args.code)
				tt.args.event.Data = data
			}
			tt.args.user.appendInitUsercodeCreatedEvent(tt.args.event)
			if tt.args.user.InitCode.Expiry != tt.result.InitCode.Expiry {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.user)
			}
		})
	}
}
