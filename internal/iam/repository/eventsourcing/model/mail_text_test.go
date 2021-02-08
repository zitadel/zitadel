package model

import (
	"encoding/json"
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

func TestAppendAddMailTextEvent(t *testing.T) {
	type args struct {
		iam      *IAM
		mailText *MailText
		event    *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append add mailText event",
			args: args{
				iam: &IAM{},
				mailText: &MailText{
					MailTextType: "PasswordReset",
					Language:     "DE"},
				event: &es_models.Event{},
			},
			result: &IAM{DefaultMailTexts: []*MailText{&MailText{
				MailTextType: "PasswordReset",
				Language:     "DE"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mailText != nil {
				data, _ := json.Marshal(tt.args.mailText)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddMailTextEvent(tt.args.event)
			if len(tt.args.iam.DefaultMailTexts) != 1 {
				t.Errorf("got wrong result should have one mailText actual: %v ", len(tt.args.iam.DefaultMailTexts))
			}
			if tt.args.iam.DefaultMailTexts[0] == tt.result.DefaultMailTexts[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultMailTexts[0], tt.args.iam.DefaultMailTexts[0])
			}
		})
	}
}

func TestAppendChangeMailTextEvent(t *testing.T) {
	type args struct {
		iam      *IAM
		mailText *MailText
		event    *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append change mailText event",
			args: args{
				iam: &IAM{DefaultMailTexts: []*MailText{&MailText{
					MailTextType: "PasswordReset",
					Language:     "DE"}}},
				mailText: &MailText{
					MailTextType: "ChangedPasswordReset",
					Language:     "DE"},
				event: &es_models.Event{},
			},
			result: &IAM{DefaultMailTexts: []*MailText{&MailText{
				MailTextType: "PasswordReset",
				Language:     "ChangedDE"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mailText != nil {
				data, _ := json.Marshal(tt.args.mailText)
				tt.args.event.Data = data
			}
			tt.args.iam.appendChangeMailTextEvent(tt.args.event)
			if len(tt.args.iam.DefaultMailTexts) != 1 {
				t.Errorf("got wrong result should have one mailText actual: %v ", len(tt.args.iam.DefaultMailTexts))
			}
			if tt.args.iam.DefaultMailTexts[0] == tt.result.DefaultMailTexts[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultMailTexts[0], tt.args.iam.DefaultMailTexts[0])
			}
		})
	}
}

func TestAppendRemoveMailTextEvent(t *testing.T) {
	type args struct {
		iam      *IAM
		mailText *MailText
		event    *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append remove mailText event",
			args: args{
				iam: &IAM{DefaultMailTexts: []*MailText{&MailText{
					MailTextType: "PasswordReset",
					Language:     "DE",
					Subject:      "Subject"}}},
				mailText: &MailText{
					MailTextType: "PasswordReset",
					Language:     "DE"},
				event: &es_models.Event{},
			},
			result: &IAM{DefaultMailTexts: []*MailText{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.mailText != nil {
				data, _ := json.Marshal(tt.args.mailText)
				tt.args.event.Data = data
			}
			tt.args.iam.appendRemoveMailTextEvent(tt.args.event)
			if len(tt.args.iam.DefaultMailTexts) != 0 {
				t.Errorf("got wrong result should have no mailText actual: %v ", len(tt.args.iam.DefaultMailTexts))
			}
		})
	}
}
