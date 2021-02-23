package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"testing"
)

func TestAppendAddIdpConfigEvent(t *testing.T) {
	type args struct {
		org   *Org
		idp   *iam_es_model.IDPConfig
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
				idp:   &iam_es_model.IDPConfig{Name: "IDPConfig"},
				event: &es_models.Event{},
			},
			result: &Org{IDPs: []*iam_es_model.IDPConfig{&iam_es_model.IDPConfig{Name: "IDPConfig"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.idp != nil {
				data, _ := json.Marshal(tt.args.idp)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddIDPConfigEvent(tt.args.event)
			if len(tt.args.org.IDPs) != 1 {
				t.Errorf("got wrong result should have one idpConfig actual: %v ", len(tt.args.org.IDPs))
			}
			if tt.args.org.IDPs[0] == tt.result.IDPs[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.IDPs[0], tt.args.org.IDPs[0])
			}
		})
	}
}

func TestAppendChangeIdpConfigEvent(t *testing.T) {
	type args struct {
		org       *Org
		idpConfig *iam_es_model.IDPConfig
		event     *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append change idp config event",
			args: args{
				org:       &Org{IDPs: []*iam_es_model.IDPConfig{&iam_es_model.IDPConfig{Name: "IDPConfig"}}},
				idpConfig: &iam_es_model.IDPConfig{Name: "IDPConfig Change"},
				event:     &es_models.Event{},
			},
			result: &Org{IDPs: []*iam_es_model.IDPConfig{&iam_es_model.IDPConfig{Name: "IDPConfig Change"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.idpConfig != nil {
				data, _ := json.Marshal(tt.args.idpConfig)
				tt.args.event.Data = data
			}
			tt.args.org.appendChangeIDPConfigEvent(tt.args.event)
			if len(tt.args.org.IDPs) != 1 {
				t.Errorf("got wrong result should have one idpConfig actual: %v ", len(tt.args.org.IDPs))
			}
			if tt.args.org.IDPs[0] == tt.result.IDPs[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.IDPs[0], tt.args.org.IDPs[0])
			}
		})
	}
}

func TestAppendRemoveIDPEvent(t *testing.T) {
	type args struct {
		org   *Org
		idp   *iam_es_model.IDPConfig
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
				org:   &Org{IDPs: []*iam_es_model.IDPConfig{&iam_es_model.IDPConfig{IDPConfigID: "IDPConfigID", Name: "IDPConfig"}}},
				idp:   &iam_es_model.IDPConfig{IDPConfigID: "IDPConfigID", Name: "IDPConfig"},
				event: &es_models.Event{},
			},
			result: &Org{IDPs: []*iam_es_model.IDPConfig{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.idp != nil {
				data, _ := json.Marshal(tt.args.idp)
				tt.args.event.Data = data
			}
			tt.args.org.appendRemoveIDPConfigEvent(tt.args.event)
			if len(tt.args.org.IDPs) != 0 {
				t.Errorf("got wrong result should have no apps actual: %v ", len(tt.args.org.IDPs))
			}
		})
	}
}

func TestAppendAppStateEvent(t *testing.T) {
	type args struct {
		org   *Org
		idp   *iam_es_model.IDPConfig
		event *es_models.Event
		state model.IDPConfigState
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append deactivate application event",
			args: args{
				org:   &Org{IDPs: []*iam_es_model.IDPConfig{&iam_es_model.IDPConfig{IDPConfigID: "IDPConfigID", Name: "IDPConfig", State: int32(model.IDPConfigStateActive)}}},
				idp:   &iam_es_model.IDPConfig{IDPConfigID: "IDPConfigID"},
				event: &es_models.Event{},
				state: model.IDPConfigStateInactive,
			},
			result: &Org{IDPs: []*iam_es_model.IDPConfig{&iam_es_model.IDPConfig{IDPConfigID: "IDPConfigID", Name: "IDPConfig", State: int32(model.IDPConfigStateInactive)}}},
		},
		{
			name: "append reactivate application event",
			args: args{
				org:   &Org{IDPs: []*iam_es_model.IDPConfig{&iam_es_model.IDPConfig{IDPConfigID: "IDPConfigID", Name: "IDPConfig", State: int32(model.IDPConfigStateInactive)}}},
				idp:   &iam_es_model.IDPConfig{IDPConfigID: "IDPConfigID"},
				event: &es_models.Event{},
				state: model.IDPConfigStateActive,
			},
			result: &Org{IDPs: []*iam_es_model.IDPConfig{&iam_es_model.IDPConfig{IDPConfigID: "IDPConfigID", Name: "IDPConfig", State: int32(model.IDPConfigStateActive)}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.idp != nil {
				data, _ := json.Marshal(tt.args.idp)
				tt.args.event.Data = data
			}
			tt.args.org.appendIDPConfigStateEvent(tt.args.event, tt.args.state)
			if len(tt.args.org.IDPs) != 1 {
				t.Errorf("got wrong result should have one idpConfig actual: %v ", len(tt.args.org.IDPs))
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
		config *iam_es_model.OIDCIDPConfig
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
				org:    &Org{IDPs: []*iam_es_model.IDPConfig{&iam_es_model.IDPConfig{IDPConfigID: "IDPConfigID"}}},
				config: &iam_es_model.OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID"},
				event:  &es_models.Event{},
			},
			result: &Org{IDPs: []*iam_es_model.IDPConfig{&iam_es_model.IDPConfig{IDPConfigID: "IDPConfigID", OIDCIDPConfig: &iam_es_model.OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.config != nil {
				data, _ := json.Marshal(tt.args.config)
				tt.args.event.Data = data
			}
			tt.args.org.appendAddOIDCIDPConfigEvent(tt.args.event)
			if len(tt.args.org.IDPs) != 1 {
				t.Errorf("got wrong result should have one idpConfig actual: %v ", len(tt.args.org.IDPs))
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
		config *iam_es_model.OIDCIDPConfig
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
				org:    &Org{IDPs: []*iam_es_model.IDPConfig{&iam_es_model.IDPConfig{IDPConfigID: "IDPConfigID", OIDCIDPConfig: &iam_es_model.OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID"}}}},
				config: &iam_es_model.OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID Changed"},
				event:  &es_models.Event{},
			},
			result: &Org{IDPs: []*iam_es_model.IDPConfig{&iam_es_model.IDPConfig{IDPConfigID: "IDPConfigID", OIDCIDPConfig: &iam_es_model.OIDCIDPConfig{IDPConfigID: "IDPConfigID", ClientID: "ClientID Changed"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.config != nil {
				data, _ := json.Marshal(tt.args.config)
				tt.args.event.Data = data
			}
			tt.args.org.appendChangeOIDCIDPConfigEvent(tt.args.event)
			if len(tt.args.org.IDPs) != 1 {
				t.Errorf("got wrong result should have one idpConfig actual: %v ", len(tt.args.org.IDPs))
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
