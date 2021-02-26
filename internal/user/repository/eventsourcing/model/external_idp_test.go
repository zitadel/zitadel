package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"testing"
)

func TestAppendExternalIDPAddedEvent(t *testing.T) {
	type args struct {
		user        *Human
		externalIDP *ExternalIDP
		event       *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Human
	}{
		{
			name: "append external idp added event",
			args: args{
				user:        &Human{},
				externalIDP: &ExternalIDP{IDPConfigID: "IDPConfigID", UserID: "UserID", DisplayName: "DisplayName"},
				event:       &es_models.Event{},
			},
			result: &Human{ExternalIDPs: []*ExternalIDP{{IDPConfigID: "IDPConfigID", UserID: "UserID", DisplayName: "DisplayName"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.externalIDP != nil {
				data, _ := json.Marshal(tt.args.externalIDP)
				tt.args.event.Data = data
			}
			tt.args.user.appendExternalIDPAddedEvent(tt.args.event)
			if len(tt.args.user.ExternalIDPs) == 0 {
				t.Error("got wrong result expected external idps on user ")
			}
			if tt.args.user.ExternalIDPs[0].UserID != tt.result.ExternalIDPs[0].UserID {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.ExternalIDPs[0].UserID, tt.args.user.ExternalIDPs[0].UserID)
			}
			if tt.args.user.ExternalIDPs[0].IDPConfigID != tt.result.ExternalIDPs[0].IDPConfigID {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.ExternalIDPs[0].IDPConfigID, tt.args.user.ExternalIDPs[0].IDPConfigID)
			}
			if tt.args.user.ExternalIDPs[0].DisplayName != tt.result.ExternalIDPs[0].DisplayName {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.ExternalIDPs[0].DisplayName, tt.args.user.ExternalIDPs[0].IDPConfigID)
			}
		})
	}
}

func TestAppendExternalIDPRemovedEvent(t *testing.T) {
	type args struct {
		user        *Human
		externalIDP *ExternalIDP
		event       *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Human
	}{
		{
			name: "append external idp removed event",
			args: args{
				user: &Human{
					ExternalIDPs: []*ExternalIDP{
						{IDPConfigID: "IDPConfigID", UserID: "UserID", DisplayName: "DisplayName"},
					}},
				externalIDP: &ExternalIDP{UserID: "UserID"},
				event:       &es_models.Event{},
			},
			result: &Human{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.externalIDP != nil {
				data, _ := json.Marshal(tt.args.externalIDP)
				tt.args.event.Data = data
			}
			tt.args.user.appendExternalIDPRemovedEvent(tt.args.event)
			if len(tt.args.user.ExternalIDPs) != 0 {
				t.Error("got wrong result expected 0 external idps on user ")
			}
		})
	}
}
