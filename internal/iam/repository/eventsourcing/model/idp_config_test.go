package model

import (
	"encoding/json"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/iam/model"
	"testing"
)

func TestIdpConfigChanges(t *testing.T) {
	type args struct {
		existing *IdpConfig
		new      *IdpConfig
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
				existing: &IdpConfig{IDPConfigID: "IdpConfigID", Name: "Name"},
				new:      &IdpConfig{IDPConfigID: "IdpConfigID", Name: "NameChanged"},
			},
			res: res{
				changesLen: 2,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &IdpConfig{IDPConfigID: "IdpConfigID", Name: "Name"},
				new:      &IdpConfig{IDPConfigID: "IdpConfigID", Name: "Name"},
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
		iam   *Iam
		idp   *IdpConfig
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Iam
	}{
		{
			name: "append add idp config event",
			args: args{
				iam:   &Iam{},
				idp:   &IdpConfig{Name: "IdpConfig"},
				event: &es_models.Event{},
			},
			result: &Iam{IDPs: []*IdpConfig{&IdpConfig{Name: "IdpConfig"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.idp != nil {
				data, _ := json.Marshal(tt.args.idp)
				tt.args.event.Data = data
			}
			tt.args.iam.appendAddIdpConfigEvent(tt.args.event)
			if len(tt.args.iam.IDPs) != 1 {
				t.Errorf("got wrong result should have one app actual: %v ", len(tt.args.iam.IDPs))
			}
			if tt.args.iam.IDPs[0] == tt.result.IDPs[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.IDPs[0], tt.args.iam.IDPs[0])
			}
		})
	}
}

func TestAppendChangeIdpConfigEvent(t *testing.T) {
	type args struct {
		project *Iam
		app     *IdpConfig
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Iam
	}{
		{
			name: "append change idp config event",
			args: args{
				project: &Iam{IDPs: []*IdpConfig{&IdpConfig{Name: "IdpConfig"}}},
				app:     &IdpConfig{Name: "IdpConfig Change"},
				event:   &es_models.Event{},
			},
			result: &Iam{IDPs: []*IdpConfig{&IdpConfig{Name: "IdpConfig Change"}}},
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
		iam   *Iam
		idp   *IdpConfig
		event *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Iam
	}{
		{
			name: "append remove idp config event",
			args: args{
				iam:   &Iam{IDPs: []*IdpConfig{&IdpConfig{IDPConfigID: "IdpConfigID", Name: "IdpConfig"}}},
				idp:   &IdpConfig{IDPConfigID: "IdpConfigID", Name: "IdpConfig"},
				event: &es_models.Event{},
			},
			result: &Iam{IDPs: []*IdpConfig{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.idp != nil {
				data, _ := json.Marshal(tt.args.idp)
				tt.args.event.Data = data
			}
			tt.args.iam.appendRemoveIdpConfigEvent(tt.args.event)
			if len(tt.args.iam.IDPs) != 0 {
				t.Errorf("got wrong result should have no apps actual: %v ", len(tt.args.iam.IDPs))
			}
		})
	}
}

func TestAppendAppStateEvent(t *testing.T) {
	type args struct {
		iam   *Iam
		idp   *IdpConfig
		event *es_models.Event
		state model.IdpConfigState
	}
	tests := []struct {
		name   string
		args   args
		result *Iam
	}{
		{
			name: "append deactivate application event",
			args: args{
				iam:   &Iam{IDPs: []*IdpConfig{&IdpConfig{IDPConfigID: "IdpConfigID", Name: "IdpConfig", State: int32(model.IdpConfigStateActive)}}},
				idp:   &IdpConfig{IDPConfigID: "IdpConfigID"},
				event: &es_models.Event{},
				state: model.IdpConfigStateInactive,
			},
			result: &Iam{IDPs: []*IdpConfig{&IdpConfig{IDPConfigID: "IdpConfigID", Name: "IdpConfig", State: int32(model.IdpConfigStateInactive)}}},
		},
		{
			name: "append reactivate application event",
			args: args{
				iam:   &Iam{IDPs: []*IdpConfig{&IdpConfig{IDPConfigID: "IdpConfigID", Name: "IdpConfig", State: int32(model.IdpConfigStateInactive)}}},
				idp:   &IdpConfig{IDPConfigID: "IdpConfigID"},
				event: &es_models.Event{},
				state: model.IdpConfigStateActive,
			},
			result: &Iam{IDPs: []*IdpConfig{&IdpConfig{IDPConfigID: "IdpConfigID", Name: "IdpConfig", State: int32(model.IdpConfigStateActive)}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.idp != nil {
				data, _ := json.Marshal(tt.args.idp)
				tt.args.event.Data = data
			}
			tt.args.iam.appendIdpConfigStateEvent(tt.args.event, tt.args.state)
			if len(tt.args.iam.IDPs) != 1 {
				t.Errorf("got wrong result should have one app actual: %v ", len(tt.args.iam.IDPs))
			}
			if tt.args.iam.IDPs[0] == tt.result.IDPs[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.IDPs[0], tt.args.iam.IDPs[0])
			}
		})
	}
}
