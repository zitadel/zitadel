//go:build integration

package oidc_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	oidc_api "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/user"
)

func TestServer_ClientCredentialsExchange(t *testing.T) {
	machine, name, clientID, clientSecret, err := Tester.CreateOIDCCredentialsClient(CTX)
	require.NoError(t, err)

	type claims struct {
		name                       string
		username                   string
		updated                    time.Time
		resourceOwnerID            any
		resourceOwnerName          any
		resourceOwnerPrimaryDomain any
		orgDomain                  any
	}
	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		scope        []string
		wantClaims   claims
		wantErr      bool
	}{
		{
			name:         "missing client ID error",
			clientID:     "",
			clientSecret: clientSecret,
			scope:        []string{oidc.ScopeOpenID},
			wantErr:      true,
		},
		{
			name:         "client not found error",
			clientID:     "foo",
			clientSecret: clientSecret,
			scope:        []string{oidc.ScopeOpenID},
			wantErr:      true,
		},
		{
			name: "machine user without secret error",
			clientID: func() string {
				name := gofakeit.Username()
				_, err := Tester.Client.Mgmt.AddMachineUser(CTX, &management.AddMachineUserRequest{
					Name:            name,
					UserName:        name,
					AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT,
				})
				require.NoError(t, err)
				return name
			}(),
			clientSecret: clientSecret,
			scope:        []string{oidc.ScopeOpenID},
			wantErr:      true,
		},
		{
			name:         "wrong secret error",
			clientID:     clientID,
			clientSecret: "bar",
			scope:        []string{oidc.ScopeOpenID},
			wantErr:      true,
		},
		{
			name:         "success",
			clientID:     clientID,
			clientSecret: clientSecret,
			scope:        []string{oidc.ScopeOpenID},
		},
		{
			name:         "openid, profile, email",
			clientID:     clientID,
			clientSecret: clientSecret,
			scope:        []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail},
			wantClaims: claims{
				name:     name,
				username: name,
				updated:  machine.GetDetails().GetChangeDate().AsTime(),
			},
		},
		{
			name:         "org id and domain scope",
			clientID:     clientID,
			clientSecret: clientSecret,
			scope: []string{
				oidc.ScopeOpenID,
				domain.OrgIDScope + Tester.Organisation.ID,
				domain.OrgDomainPrimaryScope + Tester.Organisation.Domain,
			},
			wantClaims: claims{
				resourceOwnerID:            Tester.Organisation.ID,
				resourceOwnerName:          Tester.Organisation.Name,
				resourceOwnerPrimaryDomain: Tester.Organisation.Domain,
				orgDomain:                  Tester.Organisation.Domain,
			},
		},
		{
			name:         "invalid org domain filtered",
			clientID:     clientID,
			clientSecret: clientSecret,
			scope: []string{
				oidc.ScopeOpenID,
				domain.OrgDomainPrimaryScope + Tester.Organisation.Domain,
				domain.OrgDomainPrimaryScope + "foo"},
			wantClaims: claims{
				orgDomain: Tester.Organisation.Domain,
			},
		},
		{
			name:         "invalid org id filtered",
			clientID:     clientID,
			clientSecret: clientSecret,
			scope: []string{oidc.ScopeOpenID,
				domain.OrgIDScope + Tester.Organisation.ID,
				domain.OrgIDScope + "foo",
			},
			wantClaims: claims{
				resourceOwnerID:            Tester.Organisation.ID,
				resourceOwnerName:          Tester.Organisation.Name,
				resourceOwnerPrimaryDomain: Tester.Organisation.Domain,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := rp.NewRelyingPartyOIDC(CTX, Tester.OIDCIssuer(), tt.clientID, tt.clientSecret, redirectURI, tt.scope)
			require.NoError(t, err)
			tokens, err := rp.ClientCredentials(CTX, provider, nil)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, tokens)
			userinfo, err := rp.Userinfo[*oidc.UserInfo](CTX, tokens.AccessToken, oidc.BearerToken, machine.GetUserId(), provider)
			require.NoError(t, err)
			assert.Equal(t, tt.wantClaims.resourceOwnerID, userinfo.Claims[oidc_api.ClaimResourceOwnerID])
			assert.Equal(t, tt.wantClaims.resourceOwnerName, userinfo.Claims[oidc_api.ClaimResourceOwnerName])
			assert.Equal(t, tt.wantClaims.resourceOwnerPrimaryDomain, userinfo.Claims[oidc_api.ClaimResourceOwnerPrimaryDomain])
			assert.Equal(t, tt.wantClaims.orgDomain, userinfo.Claims[domain.OrgDomainPrimaryClaim])
			assert.Equal(t, tt.wantClaims.name, userinfo.Name)
			assert.Equal(t, tt.wantClaims.username, userinfo.PreferredUsername)
			assertOIDCTime(t, userinfo.UpdatedAt, tt.wantClaims.updated)
			assert.Empty(t, userinfo.UserInfoProfile.FamilyName)
			assert.Empty(t, userinfo.UserInfoProfile.GivenName)
			assert.Empty(t, userinfo.UserInfoEmail)
			assert.Empty(t, userinfo.UserInfoPhone)
			assert.Empty(t, userinfo.Address)
		})
	}
}
