package model

import (
	"encoding/json"
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func TestAppendAddMailTextEvent(t *testing.T) {
	type args struct {
		org      *Org
		mailText *iam_es_model.MailText
		event    *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append add mail text event",
			args: args{
				org:      &Org{},
				mailText: &iam_es_model.MailText{MailTextType: "Type", Language: "DE"},
				event:    &es_models.Event{},
			},
			result: &Org{MailTexts: []*iam_es_model.MailText{&iam_es_model.MailText{MailTextType: "Type", Language: "DE"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mailText != nil {
				data, _ := json.Marshal(tt.args.mailText)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddMailTextEvent(tt.args.event)
			if len(tt.args.org.MailTexts) != 1 {
				t.Errorf("got wrong result should have one mailtext actual: %v ", len(tt.args.org.MailTexts))
			}
			if tt.result.MailTexts[0].Language != tt.args.org.MailTexts[0].Language {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.MailTexts[0].Language, tt.args.org.MailTexts[0].Language)
			}
		})
	}
}

func TestAppendChangeMailTextEvent(t *testing.T) {
	type args struct {
		org      *Org
		mailText *iam_es_model.MailText
		event    *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append change mail text event",
			args: args{
				org: &Org{MailTexts: []*iam_es_model.MailText{&iam_es_model.MailText{
					Language:     "DE",
					MailTextType: "TypeX",
				}}},
				mailText: &iam_es_model.MailText{MailTextType: "Type", Language: "DE"},
				event:    &es_models.Event{},
			},
			result: &Org{MailTexts: []*iam_es_model.MailText{&iam_es_model.MailText{
				Language:     "DE",
				MailTextType: "Type",
			}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mailText != nil {
				data, _ := json.Marshal(tt.args.mailText)
				tt.args.event.Data = data
			}
			tt.args.org.appendChangeMailTextEvent(tt.args.event)
			if len(tt.args.org.MailTexts) != 1 {
				t.Errorf("got wrong result should have one mailtext actual: %v ", len(tt.args.org.MailTexts))
			}
			if tt.result.MailTexts[0].Language != tt.args.org.MailTexts[0].Language {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.MailTexts[0].Language, tt.args.org.MailTexts[0].Language)
			}
		})
	}
}
