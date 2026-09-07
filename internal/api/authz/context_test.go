package authz

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type tokenVerifierMock struct {
	AccessTokenVerifier
	SystemTokenVerifier
	existsOrg       func(ctx context.Context, id, domain string) (string, error)
	checkOrgActive  func(ctx context.Context, orgID string) error
	projectIDOrigin func(ctx context.Context, clientID string) (string, []string, error)
}

func (m *tokenVerifierMock) RegisterServer(string, string, MethodMapping) {}

func (m *tokenVerifierMock) CheckAuthMethod(string) (Option, bool) {
	return Option{}, false
}

func (m *tokenVerifierMock) SearchMyMemberships(context.Context, string, bool) ([]*Membership, error) {
	return nil, nil
}

func (m *tokenVerifierMock) ExistsOrg(ctx context.Context, id, domain string) (string, error) {
	if m.existsOrg != nil {
		return m.existsOrg(ctx, id, domain)
	}
	if id != "" {
		return id, nil
	}
	return domain, nil
}

func (m *tokenVerifierMock) CheckOrgActive(ctx context.Context, orgID string) error {
	if m.checkOrgActive != nil {
		return m.checkOrgActive(ctx, orgID)
	}
	return nil
}

func (m *tokenVerifierMock) ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (string, []string, error) {
	if m.projectIDOrigin != nil {
		return m.projectIDOrigin(ctx, clientID)
	}
	return "", nil, nil
}

func TestVerifyTokenAndCreateCtxData_orgActive(t *testing.T) {
	activeCaller := AccessTokenVerifierFunc(func(context.Context, string) (string, string, string, string, string, error) {
		return "user1", "", "", "", "caller-org", nil
	})
	systemTokenNOK := SystemTokenVerifierFunc(func(context.Context, string, string) (Memberships, string, error) {
		return nil, "", zerrors.ThrowUnauthenticated(nil, "TEST-sys", "unauthenticated")
	})

	tests := []struct {
		name           string
		orgID          string
		orgDomain      string
		accessToken    AccessTokenVerifier
		checkOrgActive func(context.Context, string) error
		wantErr        func(error) bool
		wantResource   string
		wantOrgID      string
	}{
		{
			name:        "active caller org",
			accessToken: activeCaller,
			wantResource: "caller-org",
			wantOrgID:    "caller-org",
		},
		{
			name:        "deactivated caller org",
			accessToken: activeCaller,
			checkOrgActive: func(_ context.Context, orgID string) error {
				assert.Equal(t, "caller-org", orgID)
				return zerrors.ThrowPreconditionFailed(nil, "QUERY-oR9nA", "Errors.Org.NotActive")
			},
			wantErr: zerrors.IsUnauthenticated,
		},
		{
			name:   "deactivated target org with active caller",
			orgID:  "target-org",
			accessToken: activeCaller,
			checkOrgActive: func(_ context.Context, orgID string) error {
				// Must check the caller's org, not the target.
				assert.Equal(t, "caller-org", orgID)
				return nil
			},
			wantResource: "caller-org",
			wantOrgID:    "target-org",
		},
		{
			name: "empty resource owner for system token",
			accessToken: AccessTokenVerifierFunc(func(context.Context, string) (string, string, string, string, string, error) {
				return "", "", "", "", "", zerrors.ThrowUnauthenticated(nil, "TEST-tok", "unauthenticated")
			}),
			orgID: "system-org",
			checkOrgActive: func(context.Context, string) error {
				t.Fatal("CheckOrgActive must not be called for system tokens without resource owner")
				return nil
			},
			wantOrgID: "system-org",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier := &tokenVerifierMock{
				AccessTokenVerifier: tt.accessToken,
				SystemTokenVerifier: SystemTokenVerifierFunc(func(ctx context.Context, token, orgID string) (Memberships, string, error) {
					if tt.name == "empty resource owner for system token" {
						return Memberships{{MemberType: MemberTypeSystem, Roles: []string{"SYSTEM_OWNER"}}}, "system-user", nil
					}
					return systemTokenNOK.VerifySystemToken(ctx, token, orgID)
				}),
				checkOrgActive: tt.checkOrgActive,
			}
			got, err := VerifyTokenAndCreateCtxData(context.Background(), BearerPrefix+"token", tt.orgID, tt.orgDomain, verifier, nil)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.True(t, tt.wantErr(err))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantResource, got.ResourceOwner)
			assert.Equal(t, tt.wantOrgID, got.OrgID)
		})
	}
}
