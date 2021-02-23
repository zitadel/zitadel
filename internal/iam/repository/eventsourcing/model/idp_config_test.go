package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/model"
	"testing"
)

func TestIdpConfigChanges(t *testing.T) {
	type args struct {
		existing *IDPConfig
		new      *IDPConfig
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
			name: "idp config name changes",
			args: args{
				existing: &IDPConfig{IDPConfigID: "IDPConfigID", Name: "Name"},
				new:      &IDPConfig{IDPConfigID: "IDPConfigID", Name: "NameChanged"},
			},
			res: res{
				changesLen: 2,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &IDPConfig{IDPConfigID: "IDPConfigID", Name: "Name"},
				new:      &IDPConfig{IDPConfigID: "IDPConfigID", Name: "Name"},
			},
			res: res{
				changesLen: 1,
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

func TestAppendAddIdpConfigEvent(t *testing.T) {
	type args struct {
		iam   *IAM
		idp   *IDPConfig
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append add idp config event",
			args: args{
				iam:   &IAM{},
				idp:   &IDPConfig{Name: "IDPConfig"},
				event: &es_models.Event{},
			},
			result: &IAM{IDPs: []*IDPConfig{&IDPConfig{Name: "IDPConfig"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.idp != nil {
				data, _ := json.Marshal(tt.args.idp)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddIDPConfigEvent(tt.args.event)
			if len(tt.args.iam.IDPs) != 1 {
				t.Errorf("got wrong result should have one idpConfig actual: %v ", len(tt.args.iam.IDPs))
			}
			if tt.args.iam.IDPs[0] == tt.result.IDPs[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.IDPs[0], tt.args.iam.IDPs[0])
			}
		})
	}
}

func TestAppendChangeIdpConfigEvent(t *testing.T) {
	type args struct {
		iam       *IAM
		idpConfig *IDPConfig
		event     *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append change idp config event",
			args: args{
				iam:       &IAM{IDPs: []*IDPConfig{&IDPConfig{Name: "IDPConfig"}}},
				idpConfig: &IDPConfig{Name: "IDPConfig Change"},
				event:     &es_models.Event{},
			},
			result: &IAM{IDPs: []*IDPConfig{&IDPConfig{Name: "IDPConfig Change"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.idpConfig != nil {
				data, _ := json.Marshal(tt.args.idpConfig)
				tt.args.event.Data = data
			}
			tt.args.iam.appendChangeIDPConfigEvent(tt.args.event)
			if len(tt.args.iam.IDPs) != 1 {
				t.Errorf("got wrong result should have one idpConfig actual: %v ", len(tt.args.iam.IDPs))
			}
			if tt.args.iam.IDPs[0] == tt.result.IDPs[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.IDPs[0], tt.args.iam.IDPs[0])
			}
		})
	}
}

func TestAppendRemoveIDPEvent(t *testing.T) {
	type args struct {
		iam   *IAM
		idp   *IDPConfig
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append remove idp config event",
			args: args{
				iam:   &IAM{IDPs: []*IDPConfig{&IDPConfig{IDPConfigID: "IDPConfigID", Name: "IDPConfig"}}},
				idp:   &IDPConfig{IDPConfigID: "IDPConfigID", Name: "IDPConfig"},
				event: &es_models.Event{},
			},
			result: &IAM{IDPs: []*IDPConfig{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.idp != nil {
				data, _ := json.Marshal(tt.args.idp)
				tt.args.event.Data = data
			}
			tt.args.iam.appendRemoveIDPConfigEvent(tt.args.event)
			if len(tt.args.iam.IDPs) != 0 {
				t.Errorf("got wrong result should have no apps actual: %v ", len(tt.args.iam.IDPs))
			}
		})
	}
}

func TestAppendAppStateEvent(t *testing.T) {
	type args struct {
		iam   *IAM
		idp   *IDPConfig
		event *es_models.Event
		state model.IDPConfigState
	}
	tests := []struct {
		name   string
		args   args
		result *IAM
	}{
		{
			name: "append deactivate application event",
			args: args{
				iam:   &IAM{IDPs: []*IDPConfig{&IDPConfig{IDPConfigID: "IDPConfigID", Name: "IDPConfig", State: int32(model.IDPConfigStateActive)}}},
				idp:   &IDPConfig{IDPConfigID: "IDPConfigID"},
				event: &es_models.Event{},
				state: model.IDPConfigStateInactive,
			},
			result: &IAM{IDPs: []*IDPConfig{&IDPConfig{IDPConfigID: "IDPConfigID", Name: "IDPConfig", State: int32(model.IDPConfigStateInactive)}}},
		},
		{
			name: "append reactivate application event",
			args: args{
				iam:   &IAM{IDPs: []*IDPConfig{&IDPConfig{IDPConfigID: "IDPConfigID", Name: "IDPConfig", State: int32(model.IDPConfigStateInactive)}}},
				idp:   &IDPConfig{IDPConfigID: "IDPConfigID"},
				event: &es_models.Event{},
				state: model.IDPConfigStateActive,
			},
			result: &IAM{IDPs: []*IDPConfig{&IDPConfig{IDPConfigID: "IDPConfigID", Name: "IDPConfig", State: int32(model.IDPConfigStateActive)}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.idp != nil {
				data, _ := json.Marshal(tt.args.idp)
				tt.args.event.Data = data
			}
			tt.args.iam.appendIDPConfigStateEvent(tt.args.event, tt.args.state)
			if len(tt.args.iam.IDPs) != 1 {
				t.Errorf("got wrong result should have one idpConfig actual: %v ", len(tt.args.iam.IDPs))
			}
			if tt.args.iam.IDPs[0] == tt.result.IDPs[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.IDPs[0], tt.args.iam.IDPs[0])
			}
		})
	}
}
