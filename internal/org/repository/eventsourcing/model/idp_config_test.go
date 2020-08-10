package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"testing"
)

func TestAppendAddIdpConfigEvent(t *testing.T) {
	type args struct {
		org   *Org
		idp   *iam_es_model.IdpConfig
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append add idp config event",
			args: args{
				org:   &Org{},
				idp:   &iam_es_model.IdpConfig{Name: "IdpConfig"},
				event: &es_models.Event{},
			},
			result: &Org{IDPs: []*iam_es_model.IdpConfig{&iam_es_model.IdpConfig{Name: "IdpConfig"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.idp != nil {
				data, _ := json.Marshal(tt.args.idp)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddIdpConfigEvent(tt.args.event)
			if len(tt.args.org.IDPs) != 1 {
				t.Errorf("got wrong result should have one app actual: %v ", len(tt.args.org.IDPs))
			}
			if tt.args.org.IDPs[0] == tt.result.IDPs[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.IDPs[0], tt.args.org.IDPs[0])
			}
		})
	}
}

func TestAppendChangeIdpConfigEvent(t *testing.T) {
	type args struct {
		project *Org
		app     *iam_es_model.IdpConfig
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append change idp config event",
			args: args{
				project: &Org{IDPs: []*iam_es_model.IdpConfig{&iam_es_model.IdpConfig{Name: "IdpConfig"}}},
				app:     &iam_es_model.IdpConfig{Name: "IdpConfig Change"},
				event:   &es_models.Event{},
			},
			result: &Org{IDPs: []*iam_es_model.IdpConfig{&iam_es_model.IdpConfig{Name: "IdpConfig Change"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.app != nil {
				data, _ := json.Marshal(tt.args.app)
				tt.args.event.Data = data
			}
			tt.args.project.appendChangeIdpConfigEvent(tt.args.event)
			if len(tt.args.project.IDPs) != 1 {
				t.Errorf("got wrong result should have one app actual: %v ", len(tt.args.project.IDPs))
			}
			if tt.args.project.IDPs[0] == tt.result.IDPs[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.IDPs[0], tt.args.project.IDPs[0])
			}
		})
	}
}

func TestAppendRemoveIDPEvent(t *testing.T) {
	type args struct {
		org   *Org
		idp   *iam_es_model.IdpConfig
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append remove idp config event",
			args: args{
				org:   &Org{IDPs: []*iam_es_model.IdpConfig{&iam_es_model.IdpConfig{IDPConfigID: "IdpConfigID", Name: "IdpConfig"}}},
				idp:   &iam_es_model.IdpConfig{IDPConfigID: "IdpConfigID", Name: "IdpConfig"},
				event: &es_models.Event{},
			},
			result: &Org{IDPs: []*iam_es_model.IdpConfig{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.idp != nil {
				data, _ := json.Marshal(tt.args.idp)
				tt.args.event.Data = data
			}
			tt.args.org.appendRemoveIdpConfigEvent(tt.args.event)
			if len(tt.args.org.IDPs) != 0 {
				t.Errorf("got wrong result should have no apps actual: %v ", len(tt.args.org.IDPs))
			}
		})
	}
}

func TestAppendAppStateEvent(t *testing.T) {
	type args struct {
		org   *Org
		idp   *iam_es_model.IdpConfig
		event *es_models.Event
		state model.IdpConfigState
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append deactivate application event",
			args: args{
				org:   &Org{IDPs: []*iam_es_model.IdpConfig{&iam_es_model.IdpConfig{IDPConfigID: "IdpConfigID", Name: "IdpConfig", State: int32(model.IdpConfigStateActive)}}},
				idp:   &iam_es_model.IdpConfig{IDPConfigID: "IdpConfigID"},
				event: &es_models.Event{},
				state: model.IdpConfigStateInactive,
			},
			result: &Org{IDPs: []*iam_es_model.IdpConfig{&iam_es_model.IdpConfig{IDPConfigID: "IdpConfigID", Name: "IdpConfig", State: int32(model.IdpConfigStateInactive)}}},
		},
		{
			name: "append reactivate application event",
			args: args{
				org:   &Org{IDPs: []*iam_es_model.IdpConfig{&iam_es_model.IdpConfig{IDPConfigID: "IdpConfigID", Name: "IdpConfig", State: int32(model.IdpConfigStateInactive)}}},
				idp:   &iam_es_model.IdpConfig{IDPConfigID: "IdpConfigID"},
				event: &es_models.Event{},
				state: model.IdpConfigStateActive,
			},
			result: &Org{IDPs: []*iam_es_model.IdpConfig{&iam_es_model.IdpConfig{IDPConfigID: "IdpConfigID", Name: "IdpConfig", State: int32(model.IdpConfigStateActive)}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.idp != nil {
				data, _ := json.Marshal(tt.args.idp)
				tt.args.event.Data = data
			}
			tt.args.org.appendIdpConfigStateEvent(tt.args.event, tt.args.state)
			if len(tt.args.org.IDPs) != 1 {
				t.Errorf("got wrong result should have one app actual: %v ", len(tt.args.org.IDPs))
			}
			if tt.args.org.IDPs[0] == tt.result.IDPs[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.IDPs[0], tt.args.org.IDPs[0])
			}
		})
	}
}

func TestAppendAddOIDCIdpConfigEvent(t *testing.T) {
	type args struct {
		org    *Org
		config *iam_es_model.OidcIdpConfig
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append add oidc idp config event",
			args: args{
				org:    &Org{IDPs: []*iam_es_model.IdpConfig{&iam_es_model.IdpConfig{IDPConfigID: "IdpConfigID"}}},
				config: &iam_es_model.OidcIdpConfig{IdpConfigID: "IdpConfigID", ClientID: "ClientID"},
				event:  &es_models.Event{},
			},
			result: &Org{IDPs: []*iam_es_model.IdpConfig{&iam_es_model.IdpConfig{IDPConfigID: "IdpConfigID", OIDCIDPConfig: &iam_es_model.OidcIdpConfig{IdpConfigID: "IdpConfigID", ClientID: "ClientID"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.config != nil {
				data, _ := json.Marshal(tt.args.config)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddOidcIdpConfigEvent(tt.args.event)
			if len(tt.args.org.IDPs) != 1 {
				t.Errorf("got wrong result should have one app actual: %v ", len(tt.args.org.IDPs))
			}
			if tt.args.org.IDPs[0].OIDCIDPConfig == nil {
				t.Errorf("got wrong result should have oidc config actual: %v ", tt.args.org.IDPs[0].OIDCIDPConfig)
			}
			if tt.args.org.IDPs[0] == tt.result.IDPs[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.IDPs[0], tt.args.org.IDPs[0])
			}
		})
	}
}

func TestAppendChangeOIDCIdpConfigEvent(t *testing.T) {
	type args struct {
		org    *Org
		config *iam_es_model.OidcIdpConfig
		event  *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append change oidc idp config event",
			args: args{
				org:    &Org{IDPs: []*iam_es_model.IdpConfig{&iam_es_model.IdpConfig{IDPConfigID: "IdpConfigID", OIDCIDPConfig: &iam_es_model.OidcIdpConfig{IdpConfigID: "IdpConfigID", ClientID: "ClientID"}}}},
				config: &iam_es_model.OidcIdpConfig{IdpConfigID: "IdpConfigID", ClientID: "ClientID Changed"},
				event:  &es_models.Event{},
			},
			result: &Org{IDPs: []*iam_es_model.IdpConfig{&iam_es_model.IdpConfig{IDPConfigID: "IdpConfigID", OIDCIDPConfig: &iam_es_model.OidcIdpConfig{IdpConfigID: "IdpConfigID", ClientID: "ClientID Changed"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.config != nil {
				data, _ := json.Marshal(tt.args.config)
				tt.args.event.Data = data
			}
			tt.args.org.appendChangeOidcIdpConfigEvent(tt.args.event)
			if len(tt.args.org.IDPs) != 1 {
				t.Errorf("got wrong result should have one app actual: %v ", len(tt.args.org.IDPs))
			}
			if tt.args.org.IDPs[0].OIDCIDPConfig == nil {
				t.Errorf("got wrong result should have oidc config actual: %v ", tt.args.org.IDPs[0].OIDCIDPConfig)
			}
			if tt.args.org.IDPs[0] == tt.result.IDPs[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.IDPs[0], tt.args.org.IDPs[0])
			}
		})
	}
}
