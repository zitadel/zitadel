package model

// ToDo Michi
// import (
// 	"encoding/json"
// 	"testing"

// 	es_models "github.com/caos/zitadel/internal/eventstore/models"
// 	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
// )

// func TestAppendAddMailTextEvent(t *testing.T) {
// 	type args struct {
// 		org    *Org
// 		policy *iam_es_model.MailText
// 		event  *es_models.Event
// 	}
// 	tests := []struct {
// 		name   string
// 		args   args
// 		result *Org
// 	}{
// 		{
// 			name: "append add label policy event",
// 			args: args{
// 				org:    &Org{},
// 				policy: &iam_es_model.MailText{MailTextType: "Type", Language: "DE"},
// 				event:  &es_models.Event{},
// 			},
// 			result: &Org{MailText: &iam_es_model.MailText{MailTextType: "Type", Language: "DE"}},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.args.policy != nil {
// 				data, _ := json.Marshal(tt.args.policy)
// 				tt.args.event.Data = data
// 			}
// 			tt.args.org.appendAddMailTextEvent(tt.args.event)
// 			if tt.result.MailText.MailTextType != tt.args.org.MailText.MailTextType {
// 				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.MailText.MailTextType, tt.args.org.MailText.MailTextType)
// 			}
// 			if tt.result.MailText.Language != tt.args.org.MailText.Language {
// 				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.MailText.Language, tt.args.org.MailText.Language)
// 			}
// 		})
// 	}
// }

// func TestAppendChangeMailTextEvent(t *testing.T) {
// 	type args struct {
// 		org    *Org
// 		policy *iam_es_model.MailText
// 		event  *es_models.Event
// 	}
// 	tests := []struct {
// 		name   string
// 		args   args
// 		result *Org
// 	}{
// 		{
// 			name: "append change label policy event",
// 			args: args{
// 				org: &Org{MailText: &iam_es_model.MailText{
// 					Language:     "EN",
// 					MailTextType: "TypeX",
// 				}},
// 				policy: &iam_es_model.MailText{MailTextType: "Type", Language: "DE"},
// 				event:  &es_models.Event{},
// 			},
// 			result: &Org{MailText: &iam_es_model.MailText{
// 				Language:     "DE",
// 				MailTextType: "Type",
// 			}},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if tt.args.policy != nil {
// 				data, _ := json.Marshal(tt.args.policy)
// 				tt.args.event.Data = data
// 			}
// 			tt.args.org.appendChangeMailTextEvent(tt.args.event)
// 			if tt.result.MailText.MailTextType != tt.args.org.MailText.MailTextType {
// 				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.MailText.MailTextType, tt.args.org.MailText.MailTextType)
// 			}
// 			if tt.result.MailText.Language != tt.args.org.MailText.Language {
// 				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.MailText.Language, tt.args.org.MailText.Language)
// 			}
// 		})
// 	}
// }
