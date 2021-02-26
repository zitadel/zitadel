package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/project/model"
	es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	"testing"
)

func mockAppData(app *es_model.Application) []byte {
	data, _ := json.Marshal(app)
	return data
}

func mockOIDCConfigData(config *es_model.OIDCConfig) []byte {
	data, _ := json.Marshal(config)
	return data
}

func TestApplicationAppendEvent(t *testing.T) {
	type args struct {
		event *es_models.Event
		app   *ApplicationView
	}
	tests := []struct {
		name   string
		args   args
		result *ApplicationView
	}{
		{
			name: "append added app event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.ApplicationAdded, Data: mockAppData(&es_model.Application{Name: "AppName"})},
				app:   &ApplicationView{},
			},
			result: &ApplicationView{ProjectID: "AggregateID", Name: "AppName", State: int32(model.AppStateActive)},
		},
		{
			name: "append changed app event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.ApplicationChanged, Data: mockAppData(&es_model.Application{Name: "AppNameChanged"})},
				app:   &ApplicationView{ProjectID: "AggregateID", Name: "AppName", State: int32(model.AppStateActive)},
			},
			result: &ApplicationView{ProjectID: "AggregateID", Name: "AppNameChanged", State: int32(model.AppStateActive)},
		},
		{
			name: "append deactivate app event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.ApplicationDeactivated},
				app:   &ApplicationView{ProjectID: "AggregateID", Name: "AppName", State: int32(model.AppStateActive)},
			},
			result: &ApplicationView{ProjectID: "AggregateID", Name: "AppName", State: int32(model.AppStateInactive)},
		},
		{
			name: "append reactivate app event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.ApplicationReactivated},
				app:   &ApplicationView{ProjectID: "AggregateID", Name: "AppName", State: int32(model.AppStateInactive)},
			},
			result: &ApplicationView{ProjectID: "AggregateID", Name: "AppName", State: int32(model.AppStateActive)},
		},
		{
			name: "append added oidc config event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.OIDCConfigAdded, Data: mockOIDCConfigData(&es_model.OIDCConfig{ClientID: "clientID"})},
				app:   &ApplicationView{ProjectID: "AggregateID", Name: "AppName", State: int32(model.AppStateActive)},
			},
			result: &ApplicationView{ProjectID: "AggregateID", Name: "AppName", IsOIDC: true, OIDCClientID: "clientID", State: int32(model.AppStateActive)},
		},
		{
			name: "append changed oidc config event",
			args: args{
				event: &es_models.Event{AggregateID: "AggregateID", Sequence: 1, Type: es_model.OIDCConfigAdded, Data: mockOIDCConfigData(&es_model.OIDCConfig{ClientID: "clientIDChanged"})},
				app:   &ApplicationView{ProjectID: "AggregateID", Name: "AppName", OIDCClientID: "clientID", State: int32(model.AppStateActive)},
			},
			result: &ApplicationView{ProjectID: "AggregateID", Name: "AppName", IsOIDC: true, OIDCClientID: "clientIDChanged", State: int32(model.AppStateActive)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.app.AppendEvent(tt.args.event)
			if tt.args.app.ProjectID != tt.result.ProjectID {
				t.Errorf("got wrong result projectID: expected: %v, actual: %v ", tt.result.ProjectID, tt.args.app.ProjectID)
			}
			if tt.args.app.Name != tt.result.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.result.Name, tt.args.app.Name)
			}
			if tt.args.app.State != tt.result.State {
				t.Errorf("got wrong result state: expected: %v, actual: %v ", tt.result.State, tt.args.app.State)
			}
			if tt.args.app.IsOIDC != tt.result.IsOIDC {
				t.Errorf("got wrong result IsOIDC: expected: %v, actual: %v ", tt.result.IsOIDC, tt.args.app.IsOIDC)
			}
			if tt.args.app.OIDCClientID != tt.result.OIDCClientID {
				t.Errorf("got wrong result OIDCClientID: expected: %v, actual: %v ", tt.result.OIDCClientID, tt.args.app.OIDCClientID)
			}
		})
	}
}
