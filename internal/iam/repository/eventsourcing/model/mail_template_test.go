package model

import (
	"encoding/json"
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

func TestMailTemplateChanges(t *testing.T) {
	type args struct {
		existing *MailTemplate
		new      *MailTemplate
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
			name: "mailtemplate all attributes change",
			args: args{
				existing: &MailTemplate{Template: []byte("<doctype html>")},
				new:      &MailTemplate{Template: []byte("<!doctype html>")},
			},
			res: res{
				changesLen: 1,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &MailTemplate{Template: []byte("<!doctype html>")},
				new:      &MailTemplate{Template: []byte("<!doctype html>")},
			},
			res: res{
				changesLen: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := tt.args.existing.Changes(tt.args.new)
			if len(changes) != tt.res.changesLen {
				t.Errorf("got wrong changes len: expected: %v, actual: %v ", tt.res.changesLen, len(changes))
			}
		})
	}
}

func TestAppendAddMailTemplateEvent(t *testing.T) {
	type args struct {
		iam    *IAM
		policy *MailTemplate
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append add label policy event",
			args: args{
				iam:    new(IAM),
				policy: &MailTemplate{Template: []byte("<!doctype html>")},
				event:  new(es_models.Event),
			},
			result: &IAM{DefaultMailTemplate: &MailTemplate{Template: []byte("<!doctype html>")}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.policy != nil {
				data, _ := json.Marshal(tt.args.policy)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddMailTemplateEvent(tt.args.event)
			if string(tt.result.DefaultMailTemplate.Template) != string(tt.args.iam.DefaultMailTemplate.Template) {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultMailTemplate.Template, tt.args.iam.DefaultMailTemplate.Template)
			}
		})
	}
}

func TestAppendChangeMailTemplateEvent(t *testing.T) {
	type args struct {
		iam    *IAM
		policy *MailTemplate
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append change label policy event",
			args: args{
				iam: &IAM{DefaultMailTemplate: &MailTemplate{
					Template: []byte("<doctype html>"),
				}},
				policy: &MailTemplate{Template: []byte("<!doctype html>")},
				event:  &es_models.Event{},
			},
			result: &IAM{DefaultMailTemplate: &MailTemplate{
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
			tt.args.iam.appendChangeMailTemplateEvent(tt.args.event)
			if string(tt.result.DefaultMailTemplate.Template) != string(tt.args.iam.DefaultMailTemplate.Template) {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultMailTemplate.Template, tt.args.iam.DefaultMailTemplate.Template)
			}
		})
	}
}
