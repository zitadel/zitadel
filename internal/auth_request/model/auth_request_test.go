package model

import (
	"net"
	"reflect"
	"testing"
)

func TestAuthRequest_IsValid(t *testing.T) {
	type fields struct {
		ID            string
		AgentID       string
		BrowserInfo   *BrowserInfo
		ApplicationID string
		CallbackURI   string
		Request       Request
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"missing id, false",
			fields{},
			false,
		},
		{
			"missing agent id, false",
			fields{
				ID: "id",
			},
			false,
		},
		{
			"missing browser info, false",
			fields{
				ID:      "id",
				AgentID: "agentID",
			},
			false,
		},
		{
			"browser info invalid, false",
			fields{
				ID:          "id",
				AgentID:     "agentID",
				BrowserInfo: &BrowserInfo{},
			},
			false,
		},
		{
			"missing application id, false",
			fields{
				ID:      "id",
				AgentID: "agentID",
				BrowserInfo: &BrowserInfo{
					UserAgent:      "user agent",
					AcceptLanguage: "accept language",
					RemoteIP:       net.IPv4(29, 4, 20, 19),
				},
			},
			false,
		},
		{
			"missing callback uri, false",
			fields{
				ID:      "id",
				AgentID: "agentID",
				BrowserInfo: &BrowserInfo{
					UserAgent:      "user agent",
					AcceptLanguage: "accept language",
					RemoteIP:       net.IPv4(29, 4, 20, 19),
				},
				ApplicationID: "appID",
			},
			false,
		},
		{
			"missing request, false",
			fields{
				ID:      "id",
				AgentID: "agentID",
				BrowserInfo: &BrowserInfo{
					UserAgent:      "user agent",
					AcceptLanguage: "accept language",
					RemoteIP:       net.IPv4(29, 4, 20, 19),
				},
				ApplicationID: "appID",
				CallbackURI:   "schema://callback",
			},
			false,
		},
		{
			"request invalid, false",
			fields{
				ID:      "id",
				AgentID: "agentID",
				BrowserInfo: &BrowserInfo{
					UserAgent:      "user agent",
					AcceptLanguage: "accept language",
					RemoteIP:       net.IPv4(29, 4, 20, 19),
				},
				ApplicationID: "appID",
				CallbackURI:   "schema://callback",
				Request:       &AuthRequestOIDC{},
			},
			false,
		},
		{
			"valid auth request, true",
			fields{
				ID:      "id",
				AgentID: "agentID",
				BrowserInfo: &BrowserInfo{
					UserAgent:      "user agent",
					AcceptLanguage: "accept language",
					RemoteIP:       net.IPv4(29, 4, 20, 19),
				},
				ApplicationID: "appID",
				CallbackURI:   "schema://callback",
				Request: &AuthRequestOIDC{
					Scopes: []string{"openid"},
					CodeChallenge: &OIDCCodeChallenge{
						Challenge: "challenge",
						Method:    CodeChallengeMethodS256,
					},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthRequest{
				ID:            tt.fields.ID,
				AgentID:       tt.fields.AgentID,
				BrowserInfo:   tt.fields.BrowserInfo,
				ApplicationID: tt.fields.ApplicationID,
				CallbackURI:   tt.fields.CallbackURI,
				Request:       tt.fields.Request,
			}
			if got := a.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthRequest_MFALevel(t *testing.T) {
	type fields struct {
		Prompt       Prompt
		PossibleLOAs []LevelOfAssurance
	}
	tests := []struct {
		name   string
		fields fields
		want   MFALevel
	}{
		//PLANNED: Add / replace test cases when LOA is set
		{"-1",
			fields{},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthRequest{
				Prompt:       tt.fields.Prompt,
				PossibleLOAs: tt.fields.PossibleLOAs,
			}
			if got := a.MFALevel(); got != tt.want {
				t.Errorf("MFALevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthRequest_WithCurrentInfo(t *testing.T) {
	type fields struct {
		ID          string
		AgentID     string
		BrowserInfo *BrowserInfo
	}
	type args struct {
		info *BrowserInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *AuthRequest
	}{
		{
			"unchanged",
			fields{
				ID:      "id",
				AgentID: "agentID",
				BrowserInfo: &BrowserInfo{
					UserAgent:      "ua",
					AcceptLanguage: "de",
					RemoteIP:       net.IPv4(29, 4, 20, 19),
				},
			},
			args{
				&BrowserInfo{
					UserAgent:      "ua",
					AcceptLanguage: "de",
					RemoteIP:       net.IPv4(29, 4, 20, 19),
				},
			},
			&AuthRequest{
				ID:      "id",
				AgentID: "agentID",
				BrowserInfo: &BrowserInfo{
					UserAgent:      "ua",
					AcceptLanguage: "de",
					RemoteIP:       net.IPv4(29, 4, 20, 19),
				},
			},
		},
		{
			"changed",
			fields{
				ID:      "id",
				AgentID: "agentID",
				BrowserInfo: &BrowserInfo{
					UserAgent:      "ua",
					AcceptLanguage: "de",
					RemoteIP:       net.IPv4(29, 4, 20, 19),
				},
			},
			args{
				&BrowserInfo{
					UserAgent:      "ua",
					AcceptLanguage: "de",
					RemoteIP:       net.IPv4(16, 12, 20, 19),
				},
			},
			&AuthRequest{
				ID:      "id",
				AgentID: "agentID",
				BrowserInfo: &BrowserInfo{
					UserAgent:      "ua",
					AcceptLanguage: "de",
					RemoteIP:       net.IPv4(16, 12, 20, 19),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthRequest{
				ID:          tt.fields.ID,
				AgentID:     tt.fields.AgentID,
				BrowserInfo: tt.fields.BrowserInfo,
			}
			if got := a.WithCurrentInfo(tt.args.info); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCurrentInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
