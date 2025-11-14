package connect_middleware

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const anAPIRole = "AN_API_ROLE"

type authzRepoMock struct{}

func (v *authzRepoMock) VerifyAccessToken(ctx context.Context, token, clientID, projectID string) (string, string, string, string, string, error) {
	return "", "", "", "", "", nil
}

func (v *authzRepoMock) SearchMyMemberships(ctx context.Context, orgID string, _ bool) ([]*authz.Membership, error) {
	return authz.Memberships{{
		MemberType:  authz.MemberTypeOrganization,
		AggregateID: orgID,
		Roles:       []string{anAPIRole},
	}}, nil
}

func (v *authzRepoMock) ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (string, []string, error) {
	return "", nil, nil
}

func (v *authzRepoMock) ExistsOrg(ctx context.Context, orgID, domain string) (string, error) {
	return orgID, nil
}

func (v *authzRepoMock) VerifierClientID(ctx context.Context, appName string) (string, string, error) {
	return "", "", nil
}

var (
	accessTokenOK = authz.AccessTokenVerifierFunc(func(ctx context.Context, token string) (userID string, clientID string, agentID string, prefLan string, resourceOwner string, err error) {
		return "user1", "", "", "", "org1", nil
	})
	accessTokenNOK = authz.AccessTokenVerifierFunc(func(ctx context.Context, token string) (userID string, clientID string, agentID string, prefLan string, resourceOwner string, err error) {
		return "", "", "", "", "", zerrors.ThrowUnauthenticated(nil, "TEST-fQHDI", "unauthenticaded")
	})
	systemTokenNOK = authz.SystemTokenVerifierFunc(func(ctx context.Context, token string, orgID string) (memberships authz.Memberships, userID string, err error) {
		return nil, "", errors.New("system token error")
	})
)

type mockOrgFromRequest struct {
	id string
}

func (m *mockOrgFromRequest) OrganizationFromRequestConnect() *Organization {
	return &Organization{
		ID:     m.id,
		Domain: "",
	}
}

func Test_authorize(t *testing.T) {
	type args struct {
		ctx        context.Context
		req        connect.AnyRequest
		handler    func(t *testing.T) connect.UnaryFunc
		verifier   func() authz.APITokenVerifier
		authConfig authz.Config
	}
	type res struct {
		want    interface{}
		wantErr bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"no token needed ok",
			args{
				ctx:     context.Background(),
				req:     &mockReq[struct{}]{procedure: "/no/token/needed"},
				handler: emptyMockHandler(&connect.Response[struct{}]{}, authz.CtxData{}),
				verifier: func() authz.APITokenVerifier {
					verifier := authz.StartAPITokenVerifier(&authzRepoMock{}, accessTokenOK, systemTokenNOK)
					verifier.RegisterServer("need", "need", authz.MethodMapping{})
					return verifier
				},
			},
			res{
				&connect.Response[struct{}]{},
				false,
			},
		},
		{
			"auth header missing error",
			args{
				ctx:     context.Background(),
				req:     &mockReq[struct{}]{procedure: "/need/authentication"},
				handler: emptyMockHandler(&connect.Response[struct{}]{}, authz.CtxData{}),
				verifier: func() authz.APITokenVerifier {
					verifier := authz.StartAPITokenVerifier(&authzRepoMock{}, accessTokenOK, systemTokenNOK)
					verifier.RegisterServer("need", "need", authz.MethodMapping{"/need/authentication": authz.Option{Permission: "authenticated"}})
					return verifier
				},
				authConfig: authz.Config{},
			},
			res{
				nil,
				true,
			},
		},
		{
			"unauthorized error",
			args{
				ctx:     context.Background(),
				req:     &mockReq[struct{}]{procedure: "/need/authentication", header: http.Header{"Authorization": []string{"wrong"}}},
				handler: emptyMockHandler(&connect.Response[struct{}]{}, authz.CtxData{}),
				verifier: func() authz.APITokenVerifier {
					verifier := authz.StartAPITokenVerifier(&authzRepoMock{}, accessTokenOK, systemTokenNOK)
					verifier.RegisterServer("need", "need", authz.MethodMapping{"/need/authentication": authz.Option{Permission: "authenticated"}})
					return verifier
				},
				authConfig: authz.Config{},
			},
			res{
				nil,
				true,
			},
		},
		{
			"authorized ok",
			args{
				ctx: context.Background(),
				req: &mockReq[struct{}]{procedure: "/need/authentication", header: http.Header{"Authorization": []string{"Bearer token"}}},
				handler: emptyMockHandler(&connect.Response[struct{}]{}, authz.CtxData{
					UserID:        "user1",
					OrgID:         "org1",
					ResourceOwner: "org1",
				}),
				verifier: func() authz.APITokenVerifier {
					verifier := authz.StartAPITokenVerifier(&authzRepoMock{}, accessTokenOK, systemTokenNOK)
					verifier.RegisterServer("need", "need", authz.MethodMapping{"/need/authentication": authz.Option{Permission: "authenticated"}})
					return verifier
				},
				authConfig: authz.Config{},
			},
			res{
				&connect.Response[struct{}]{},
				false,
			},
		},
		{
			"authorized ok, org by request",
			args{
				ctx: context.Background(),
				req: &mockReq[mockOrgFromRequest]{
					Request:   connect.Request[mockOrgFromRequest]{Msg: &mockOrgFromRequest{"id"}},
					procedure: "/need/authentication",
					header:    http.Header{"Authorization": []string{"Bearer token"}},
				},
				handler: emptyMockHandler(&connect.Response[mockOrgFromRequest]{Msg: &mockOrgFromRequest{"id"}}, authz.CtxData{
					UserID:        "user1",
					OrgID:         "id",
					ResourceOwner: "org1",
				}),
				verifier: func() authz.APITokenVerifier {
					verifier := authz.StartAPITokenVerifier(&authzRepoMock{}, accessTokenOK, systemTokenNOK)
					verifier.RegisterServer("need", "need", authz.MethodMapping{"/need/authentication": authz.Option{Permission: "authenticated"}})
					return verifier
				},
				authConfig: authz.Config{},
			},
			res{
				&connect.Response[mockOrgFromRequest]{Msg: &mockOrgFromRequest{"id"}},
				false,
			},
		},
		{
			"permission denied error",
			args{
				ctx: context.Background(),
				req: &mockReq[struct{}]{procedure: "/need/authentication", header: http.Header{"Authorization": []string{"Bearer token"}}},
				handler: emptyMockHandler(&connect.Response[struct{}]{}, authz.CtxData{
					UserID:        "user1",
					OrgID:         "org1",
					ResourceOwner: "org1",
				}),
				verifier: func() authz.APITokenVerifier {
					verifier := authz.StartAPITokenVerifier(&authzRepoMock{}, accessTokenOK, systemTokenNOK)
					verifier.RegisterServer("need", "need", authz.MethodMapping{"/need/authentication": authz.Option{Permission: "to.do.something"}})
					return verifier
				},
				authConfig: authz.Config{
					RolePermissionMappings: []authz.RoleMapping{{
						Role:        anAPIRole,
						Permissions: []string{"to.do.something.else"},
					}},
				},
			},
			res{
				nil,
				true,
			},
		},
		{
			"permission ok",
			args{
				ctx: context.Background(),
				req: &mockReq[struct{}]{procedure: "/need/authentication", header: http.Header{"Authorization": []string{"Bearer token"}}},
				handler: emptyMockHandler(&connect.Response[struct{}]{}, authz.CtxData{
					UserID:        "user1",
					OrgID:         "org1",
					ResourceOwner: "org1",
				}),
				verifier: func() authz.APITokenVerifier {
					verifier := authz.StartAPITokenVerifier(&authzRepoMock{}, accessTokenOK, systemTokenNOK)
					verifier.RegisterServer("need", "need", authz.MethodMapping{"/need/authentication": authz.Option{Permission: "to.do.something"}})
					return verifier
				},
				authConfig: authz.Config{
					RolePermissionMappings: []authz.RoleMapping{{
						Role:        anAPIRole,
						Permissions: []string{"to.do.something"},
					}},
				},
			},
			res{
				&connect.Response[struct{}]{},
				false,
			},
		},
		{
			"system token permission denied error",
			args{
				ctx:     context.Background(),
				req:     &mockReq[struct{}]{procedure: "/need/authentication", header: http.Header{"Authorization": []string{"Bearer token"}}},
				handler: emptyMockHandler(&connect.Response[struct{}]{}, authz.CtxData{}),
				verifier: func() authz.APITokenVerifier {
					verifier := authz.StartAPITokenVerifier(&authzRepoMock{}, accessTokenNOK, authz.SystemTokenVerifierFunc(func(ctx context.Context, token string, orgID string) (memberships authz.Memberships, userID string, err error) {
						return authz.Memberships{{
							MemberType: authz.MemberTypeSystem,
							Roles:      []string{"A_SYSTEM_ROLE"},
						}}, "systemuser", nil
					}))
					verifier.RegisterServer("need", "need", authz.MethodMapping{"/need/authentication": authz.Option{Permission: "to.do.something"}})
					return verifier
				},
				authConfig: authz.Config{
					RolePermissionMappings: []authz.RoleMapping{{
						Role:        "A_SYSTEM_ROLE",
						Permissions: []string{"to.do.something.else"},
					}},
				},
			},
			res{
				nil,
				true,
			},
		},
		{
			"system token permission denied error",
			args{
				ctx: context.Background(),
				req: &mockReq[struct{}]{procedure: "/need/authentication", header: http.Header{"Authorization": []string{"Bearer token"}}},
				handler: emptyMockHandler(&connect.Response[struct{}]{}, authz.CtxData{
					UserID: "systemuser",
					SystemMemberships: authz.Memberships{{
						MemberType: authz.MemberTypeSystem,
						Roles:      []string{"A_SYSTEM_ROLE"},
					}},
					SystemUserPermissions: []authz.SystemUserPermissions{{
						MemberType:  authz.MemberTypeSystem,
						Permissions: []string{"to.do.something"},
					}},
				}),
				verifier: func() authz.APITokenVerifier {
					verifier := authz.StartAPITokenVerifier(&authzRepoMock{}, accessTokenNOK, authz.SystemTokenVerifierFunc(func(ctx context.Context, token string, orgID string) (memberships authz.Memberships, userID string, err error) {
						return authz.Memberships{{
							MemberType: authz.MemberTypeSystem,
							Roles:      []string{"A_SYSTEM_ROLE"},
						}}, "systemuser", nil
					}))
					verifier.RegisterServer("need", "need", authz.MethodMapping{"/need/authentication": authz.Option{Permission: "to.do.something"}})
					return verifier
				},
				authConfig: authz.Config{
					RolePermissionMappings: []authz.RoleMapping{{
						Role:        "A_SYSTEM_ROLE",
						Permissions: []string{"to.do.something"},
					}},
				},
			},
			res{
				&connect.Response[struct{}]{},
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := authorize(tt.args.ctx, tt.args.req, tt.args.handler(t), tt.args.verifier(), tt.args.authConfig, tt.args.authConfig)
			if (err != nil) != tt.res.wantErr {
				t.Errorf("authorize() error = %v, wantErr %v", err, tt.res.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.res.want) {
				t.Errorf("authorize() got = %v, want %v", got, tt.res.want)
			}
		})
	}
}
