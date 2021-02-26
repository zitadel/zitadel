package model

import (
	"encoding/json"
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

func TestAppendAddMailTemplateEvent(t *testing.T) {
	type args struct {
		org    *Org
		policy *iam_es_model.MailTemplate
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append add label policy event",
			args: args{
				org:    &Org{},
				policy: &iam_es_model.MailTemplate{Template: []byte("<!doctype html>")},
				event:  &es_models.Event{},
			},
			result: &Org{MailTemplate: &iam_es_model.MailTemplate{Template: []byte("<!doctype html>")}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddMailTemplateEvent(tt.args.event)
			if string(tt.result.MailTemplate.Template) != string(tt.args.org.MailTemplate.Template) {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.MailTemplate.Template, tt.args.org.MailTemplate.Template)
			}
		})
	}
}

func TestAppendChangeMailTemplateEvent(t *testing.T) {
	type args struct {
		org    *Org
		policy *iam_es_model.MailTemplate
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append change label policy event",
			args: args{
				org: &Org{MailTemplate: &iam_es_model.MailTemplate{
					Template: []byte("<x!doctype html>"),
				}},
				policy: &iam_es_model.MailTemplate{Template: []byte("<!doctype html>")},
				event:  &es_models.Event{},
			},
			result: &Org{MailTemplate: &iam_es_model.MailTemplate{
				Template: []byte("<!doctype html>"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.org.appendChangeMailTemplateEvent(tt.args.event)
			if string(tt.result.MailTemplate.Template) != string(tt.args.org.MailTemplate.Template) {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.MailTemplate.Template, tt.args.org.MailTemplate.Template)
			}
		})
	}
}
