package oidc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/domain"
)

func Test_cacheUserInfo(t *testing.T) {
	type call struct {
		roleAssertion bool
		triggerType   domain.TriggerType
	}
	errAction := errors.New("action failed")

	tests := []struct {
		name      string
		calls     []call
		failOn    *call
		wantCalls int
	}{
		{
			name: "same role assertion and trigger type calls once",
			calls: []call{
				{roleAssertion: true, triggerType: domain.TriggerTypePreAccessTokenCreation},
				{roleAssertion: true, triggerType: domain.TriggerTypePreAccessTokenCreation},
			},
			wantCalls: 1,
		},
		{
			name: "different trigger types call once each",
			calls: []call{
				{roleAssertion: true, triggerType: domain.TriggerTypePreAccessTokenCreation},
				{roleAssertion: true, triggerType: domain.TriggerTypePreUserinfoCreation},
				{roleAssertion: true, triggerType: domain.TriggerTypePreAccessTokenCreation},
				{roleAssertion: true, triggerType: domain.TriggerTypePreUserinfoCreation},
			},
			wantCalls: 2,
		},
		{
			name: "different role assertions call once each",
			calls: []call{
				{roleAssertion: true, triggerType: domain.TriggerTypePreUserinfoCreation},
				{roleAssertion: false, triggerType: domain.TriggerTypePreUserinfoCreation},
			},
			wantCalls: 2,
		},
		{
			name: "error is returned again without another call",
			calls: []call{
				{roleAssertion: false, triggerType: domain.TriggerTypePreAccessTokenCreation},
				{roleAssertion: false, triggerType: domain.TriggerTypePreAccessTokenCreation},
			},
			failOn:    &call{roleAssertion: false, triggerType: domain.TriggerTypePreAccessTokenCreation},
			wantCalls: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotCalls int
			getUserInfo := cacheUserInfo(func(_ context.Context, roleAssertion bool, triggerType domain.TriggerType, _ *domain.TokenActor) (*oidc.UserInfo, error) {
				gotCalls++
				if tt.failOn != nil && *tt.failOn == (call{roleAssertion: roleAssertion, triggerType: triggerType}) {
					return nil, errAction
				}
				return &oidc.UserInfo{Subject: "user1"}, nil
			})

			first := make(map[call]*oidc.UserInfo)
			for _, c := range tt.calls {
				userInfo, err := getUserInfo(context.Background(), c.roleAssertion, c.triggerType, nil)
				if tt.failOn != nil && *tt.failOn == c {
					require.ErrorIs(t, err, errAction)
					assert.Nil(t, userInfo)
					continue
				}
				require.NoError(t, err)
				if prev, ok := first[c]; ok {
					assert.Same(t, prev, userInfo, "a repeated call must return the stored user info")
				}
				first[c] = userInfo
			}
			assert.Equal(t, tt.wantCalls, gotCalls)
		})
	}
}
