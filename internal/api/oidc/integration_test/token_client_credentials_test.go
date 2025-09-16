//go:build integration

package oidc_test

import (
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	oidc_api "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/user"
)

func TestServer_ClientCredentialsExchange(t *testing.T) {
	machine, name, clientID, clientSecret, err := Instance.CreateOIDCCredentialsClient(CTX)
	require.NoError(t, err)

	_, _, clientIDInactive, clientSecretInactive, err := Instance.CreateOIDCCredentialsClientInactive(CTX)
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
				name := integration.Username()
				_, err := Instance.Client.Mgmt.AddMachineUser(CTX, &management.AddMachineUserRequest{
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
			name:         "inactive machine user error",
			clientID:     clientIDInactive,
			clientSecret: clientSecretInactive,
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
			name:         "openid, profile, email, zitadel",
			clientID:     clientID,
			clientSecret: clientSecret,
			scope:        []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, domain.ProjectScopeZITADEL},
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
				domain.OrgIDScope + Instance.DefaultOrg.Id,
				domain.OrgDomainPrimaryScope + Instance.DefaultOrg.PrimaryDomain,
			},
			wantClaims: claims{
				resourceOwnerID:            Instance.DefaultOrg.Id,
				resourceOwnerName:          Instance.DefaultOrg.Name,
				resourceOwnerPrimaryDomain: Instance.DefaultOrg.PrimaryDomain,
				orgDomain:                  Instance.DefaultOrg.PrimaryDomain,
			},
		},
		{
			name:         "invalid org domain filtered",
			clientID:     clientID,
			clientSecret: clientSecret,
			scope: []string{
				oidc.ScopeOpenID,
				domain.OrgDomainPrimaryScope + Instance.DefaultOrg.PrimaryDomain,
				domain.OrgDomainPrimaryScope + "foo"},
			wantClaims: claims{
				orgDomain: Instance.DefaultOrg.PrimaryDomain,
			},
		},
		{
			name:         "invalid org id filtered",
			clientID:     clientID,
			clientSecret: clientSecret,
			scope: []string{oidc.ScopeOpenID,
				domain.OrgIDScope + Instance.DefaultOrg.Id,
				domain.OrgIDScope + "foo",
			},
			wantClaims: claims{
				resourceOwnerID:            Instance.DefaultOrg.Id,
				resourceOwnerName:          Instance.DefaultOrg.Name,
				resourceOwnerPrimaryDomain: Instance.DefaultOrg.PrimaryDomain,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := rp.NewRelyingPartyOIDC(CTX, Instance.OIDCIssuer(), tt.clientID, tt.clientSecret, redirectURI, tt.scope)
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

			_, err = Instance.Client.Auth.GetMyUser(integration.WithAuthorizationToken(CTX, tokens.AccessToken), &auth.GetMyUserRequest{})
			if slices.Contains(tt.scope, domain.ProjectScopeZITADEL) {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
