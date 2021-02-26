package model

import (
	"encoding/json"
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

func TestOIDCConfigChanges(t *testing.T) {
	type args struct {
		existingConfig *OIDCConfig
		newConfig      *OIDCConfig
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
			name: "all possible values change",
			args: args{
				existingConfig: &OIDCConfig{
					AppID:                  "AppID",
					RedirectUris:           []string{"RedirectUris"},
					ResponseTypes:          []int32{1},
					GrantTypes:             []int32{1},
					ApplicationType:        1,
					AuthMethodType:         1,
					PostLogoutRedirectUris: []string{"PostLogoutRedirectUris"},
				},
				newConfig: &OIDCConfig{
					AppID:                  "AppID",
					RedirectUris:           []string{"RedirectUrisChanged"},
					ResponseTypes:          []int32{2},
					GrantTypes:             []int32{2},
					ApplicationType:        2,
					AuthMethodType:         2,
					PostLogoutRedirectUris: []string{"PostLogoutRedirectUrisChanged"},
				},
			},
			res: res{
				changesLen: 7,
			},
		},
		{
			name: "no changes",
			args: args{
				existingConfig: &OIDCConfig{
					AppID:                  "AppID",
					RedirectUris:           []string{"RedirectUris"},
					ResponseTypes:          []int32{1},
					GrantTypes:             []int32{1},
					ApplicationType:        1,
					AuthMethodType:         1,
					PostLogoutRedirectUris: []string{"PostLogoutRedirectUris"},
				},
				newConfig: &OIDCConfig{
					AppID:                  "AppID",
					RedirectUris:           []string{"RedirectUris"},
					ResponseTypes:          []int32{1},
					GrantTypes:             []int32{1},
					ApplicationType:        1,
					AuthMethodType:         1,
					PostLogoutRedirectUris: []string{"PostLogoutRedirectUris"},
				},
			},
			res: res{
				changesLen: 1,
			},
		},
		{
			name: "change not changeable attributes",
			args: args{
				existingConfig: &OIDCConfig{
					AppID:    "AppID",
					ClientID: "ClientID",
				},
				newConfig: &OIDCConfig{
					AppID:    "AppIDChange",
					ClientID: "ClientIDChange",
				},
			},
			res: res{
				changesLen: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := tt.args.existingConfig.Changes(tt.args.newConfig)
			if len(changes) != tt.res.changesLen {
				t.Errorf("got wrong changes len: expected: %v, actual: %v ", tt.res.changesLen, len(changes))
			}
		})
	}
}

func TestAppendAddOIDCConfigEvent(t *testing.T) {
	type args struct {
		project *Project
		config  *OIDCConfig
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append add application event",
			args: args{
				project: &Project{
					Applications: []*Application{
						{AppID: "AppID"},
					},
				},
				config: &OIDCConfig{AppID: "AppID", ClientID: "ClientID"},
				event:  &es_models.Event{},
			},
			result: &Project{
				Applications: []*Application{
					{AppID: "AppID", OIDCConfig: &OIDCConfig{AppID: "AppID", ClientID: "ClientID"}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.config != nil {
				data, _ := json.Marshal(tt.args.config)
				tt.args.event.Data = data
			}
			tt.args.project.appendAddOIDCConfigEvent(tt.args.event)
			if len(tt.args.project.Applications) != 1 {
				t.Errorf("got wrong result should have one app actual: %v ", len(tt.args.project.Applications))
			}
			if tt.args.project.Applications[0].OIDCConfig == nil {
				t.Errorf("got wrong result should have oidc config actual: %v ", tt.args.project.Applications[0].OIDCConfig)
			}
			if tt.args.project.Applications[0] == tt.result.Applications[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Applications[0], tt.args.project.Applications[0])
			}
		})
	}
}

func TestAppendChangeOIDCConfigEvent(t *testing.T) {
	type args struct {
		project *Project
		config  *OIDCConfig
		event   *es_models.Event
	}
	tests := []struct {
		name   string
		args   args
		result *Project
	}{
		{
			name: "append change application event",
			args: args{
				project: &Project{
					Applications: []*Application{
						{AppID: "AppID", OIDCConfig: &OIDCConfig{AppID: "AppID", ClientID: "ClientID"}},
					},
				},
				config: &OIDCConfig{AppID: "AppID", ClientID: "ClientID Changed"},
				event:  &es_models.Event{},
			},
			result: &Project{
				Applications: []*Application{
					{AppID: "AppID", OIDCConfig: &OIDCConfig{AppID: "AppID", ClientID: "ClientID Changed"}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.config != nil {
				data, _ := json.Marshal(tt.args.config)
				tt.args.event.Data = data
			}
			tt.args.project.appendChangeOIDCConfigEvent(tt.args.event)
			if len(tt.args.project.Applications) != 1 {
				t.Errorf("got wrong result should have one app actual: %v ", len(tt.args.project.Applications))
			}
			if tt.args.project.Applications[0].OIDCConfig == nil {
				t.Errorf("got wrong result should have oidc config actual: %v ", tt.args.project.Applications[0].OIDCConfig)
			}
			if tt.args.project.Applications[0] == tt.result.Applications[0] {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.Applications[0], tt.args.project.Applications[0])
			}
		})
	}
}
