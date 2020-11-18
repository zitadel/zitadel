package model

// import (
// 	"encoding/json"
// 	"testing"

// 	es_models "github.com/caos/zitadel/internal/eventstore/models"
// )

// func TestMailTextChanges(t *testing.T) {
// 	type args struct {
// 		existing *MailText
// 		new      *MailText
// 	}
// 	type res struct {
// 		changesLen int
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		res  res
// 	}{
// 		{
// 			name: "mailtext all attributes change",
// 			args: args{
// 				existing: &MailText{
// 					MailTextType: "PasswordReset",
// 					Language:     "DE",
// 					Title:        "Zitadel - User initialisieren",
// 					PreHeader:    "User initialisieren",
// 					Subject:      "User initialisieren",
// 					Greeting:     "Hallo {{.FirstName}} {{.LastName}}",
// 					Text:         "Dieser Benutzer wurde soeben im Zitadel erstellt.",
// 					ButtonText:   "Initialisierung abschliessen"},
// 				new: &MailText{
// 					MailTextType: "InitCode",
// 					Language:     "DE",
// 					Title:        "Zitadel - User initialisieren",
// 					PreHeader:    "User initialisieren",
// 					Subject:      "User initialisieren",
// 					Greeting:     "Hallo {{.FirstName}} {{.LastName}}",
// 					Text:         "Dieser Benutzer wurde soeben im Zitadel erstellt.",
// 					ButtonText:   "Initialisierung abschliessen"},
// 			},
// 			res: res{
// 				changesLen: 2,
// 			},
// 		},
// 		{
// 			name: "no changes",
// 			args: args{
// 				existing: &MailText{
// 					MailTextType: "InitCode",
// 					Language:     "DE",
// 					Title:        "Zitadel - User initialisieren",
// 					PreHeader:    "User initialisieren",
// 					Subject:      "User initialisieren",
// 					Greeting:     "Hallo {{.FirstName}} {{.LastName}}",
// 					Text:         "Dieser Benutzer wurde soeben im Zitadel erstellt.",
// 					ButtonText:   "Initialisierung abschliessen"},
// 				new: &MailText{
// 					MailTextType: "InitCode",
// 					Language:     "DE",
// 					Title:        "Zitadel - User initialisieren",
// 					PreHeader:    "User initialisieren",
// 					Subject:      "User initialisieren",
// 					Greeting:     "Hallo {{.FirstName}} {{.LastName}}",
// 					Text:         "Dieser Benutzer wurde soeben im Zitadel erstellt.",
// 					ButtonText:   "Initialisierung abschliessen"},
// 			},
// 			res: res{
// 				changesLen: 0,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			changes := tt.args.existing.Changes(tt.args.new)
// 			if len(changes) != tt.res.changesLen {
// 				t.Errorf("got wrong changes len: expected: %v, actual: %v ", tt.res.changesLen, len(changes))
// 			}
// 		})
// 	}
// }

// func TestAppendAddMailTextEvent(t *testing.T) {
// 	type args struct {
// 		iam    *IAM
// 		policy *MailText
// 		event  *es_models.Event
// 	}
// 	tests := []struct {
// 		name   string
// 		args   args
// 		result *IAM
// 	}{
// 		{
// 			name: "append add label policy event",
// 			args: args{
// 				iam: new(IAM),
// 				policy: &MailText{
// 					MailTextType: "InitCode",
// 					Language:     "DE",
// 					Title:        "Zitadel - User initialisieren",
// 					PreHeader:    "User initialisieren",
// 					Subject:      "User initialisieren",
// 					Greeting:     "Hallo {{.FirstName}} {{.LastName}}",
// 					Text:         "Dieser Benutzer wurde soeben im Zitadel erstellt.",
// 					ButtonText:   "Initialisierung abschliessen"},
// 				event: new(es_models.Event),
// 			},
// 			result: &IAM{DefaultMailTexts: &MailText{
// 				MailTextType: "InitCode",
// 				Language:     "DE",
// 				Title:        "Zitadel - User initialisieren",
// 				PreHeader:    "User initialisieren",
// 				Subject:      "User initialisieren",
// 				Greeting:     "Hallo {{.FirstName}} {{.LastName}}",
// 				Text:         "Dieser Benutzer wurde soeben im Zitadel erstellt.",
// 				ButtonText:   "Initialisierung abschliessen"}},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.args.policy != nil {
// 				data, _ := json.Marshal(tt.args.policy)
// 				tt.args.event.Data = data
// 			}
// 			tt.args.iam.appendAddMailTextEvent(tt.args.event)
// 			if tt.result.DefaultMailTexts.MailTextType != tt.args.iam.DefaultMailTexts.MailTextType {
// 				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultMailTexts.MailTextType, tt.args.iam.DefaultMailText.MailTextType)
// 			}
// 		})
// 	}
// }

// func TestAppendChangeMailTextEvent(t *testing.T) {
// 	type args struct {
// 		iam    *IAM
// 		policy *MailText
// 		event  *es_models.Event
// 	}
// 	tests := []struct {
// 		name   string
// 		args   args
// 		result *IAM
// 	}{
// 		{
// 			name: "append change label policy event",
// 			args: args{
// 				iam: &IAM{DefaultMailTexts: &MailText{
// 					MailTextType: "PasswordReset",
// 					Language:     "DE",
// 					Title:        "Zitadel - User initialisieren",
// 					PreHeader:    "User initialisieren",
// 					Subject:      "User initialisieren",
// 					Greeting:     "Hallo {{.FirstName}} {{.LastName}}",
// 					Text:         "Dieser Benutzer wurde soeben im Zitadel erstellt.",
// 					ButtonText:   "Initialisierung abschliessen",
// 				}},
// 				policy: &MailText{
// 					MailTextType: "InitCode",
// 					Language:     "DE",
// 					Title:        "Zitadel - User initialisieren",
// 					PreHeader:    "User initialisieren",
// 					Subject:      "User initialisieren",
// 					Greeting:     "Hallo {{.FirstName}} {{.LastName}}",
// 					Text:         "Dieser Benutzer wurde soeben im Zitadel erstellt.",
// 					ButtonText:   "Initialisierung abschliessen",
// 				},
// 				event: &es_models.Event{},
// 			},
// 			result: &IAM{DefaultMailTexts: &MailText{
// 				MailTextType: "InitCode",
// 				Language:     "DE",
// 				Title:        "Zitadel - User initialisieren",
// 				PreHeader:    "User initialisieren",
// 				Subject:      "User initialisieren",
// 				Greeting:     "Hallo {{.FirstName}} {{.LastName}}",
// 				Text:         "Dieser Benutzer wurde soeben im Zitadel erstellt.",
// 				ButtonText:   "Initialisierung abschliessen",
// 			}},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.args.policy != nil {
// 				data, _ := json.Marshal(tt.args.policy)
// 				tt.args.event.Data = data
// 			}
// 			tt.args.iam.appendChangeMailTextEvent(tt.args.event)
// 			if tt.result.DefaultMailTexts.MailTextType != tt.args.iam.DefaultMailTexts.MailTextType {
// 				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.DefaultMailTexts.Mai, tt.args.iam.DefaultMailTexts.MailTextType)
// 			}
// 		})
// 	}
// }
