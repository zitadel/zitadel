package auth

import (
	"context"
	"testing"

	caos_errs "github.com/caos/zitadel/internal/errors"
)

func getTestCtx(userID, orgID string) context.Context {
	return context.WithValue(context.Background(), dataKey, CtxData{UserID: userID, OrgID: orgID})
}

type testVerifier struct {
	grants []*Grant
}

func (v *testVerifier) VerifyAccessToken(ctx context.Context, token string) (string, string, string, error) {
	return "userID", "clientID", "agentID", nil
}

func (v *testVerifier) ResolveGrants(ctx context.Context, sub, orgID string) ([]*Grant, error) {
	return v.grants, nil
}

func (v *testVerifier) GetProjectIDByClientID(ctx context.Context, clientID string) (string, error) {
	return "", nil
}

func equalStringArray(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func Test_GetUserMethodPermissions(t *testing.T) {
	type args struct {
		ctx          context.Context
		verifier     TokenVerifier
		requiredPerm string
		authConfig   *Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errFunc func(err error) bool
		result  []string
	}{
		{
			name: "Empty Context",
			args: args{
				ctx: getTestCtx("", ""),
				verifier: &testVerifier{grants: []*Grant{&Grant{
					Roles: []string{"ORG_OWNER"}}}},
				requiredPerm: "project.read",
				authConfig: &Config{
					RolePermissionMappings: []RoleMapping{
						RoleMapping{
							Role:        "IAM_OWNER",
							Permissions: []string{"project.read"},
						},
						RoleMapping{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
			},
			wantErr: true,
			errFunc: caos_errs.IsUnauthenticated,
			result:  []string{"project.read"},
		},
		{
			name: "No Grants",
			args: args{
				ctx:          getTestCtx("", ""),
				verifier:     &testVerifier{grants: []*Grant{}},
				requiredPerm: "project.read",
				authConfig: &Config{
					RolePermissionMappings: []RoleMapping{
						RoleMapping{
							Role:        "IAM_OWNER",
							Permissions: []string{"project.read"},
						},
						RoleMapping{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
			},
			result: make([]string, 0),
		},
		{
			name: "Get Permissions",
			args: args{
				ctx: getTestCtx("userID", "orgID"),
				verifier: &testVerifier{grants: []*Grant{&Grant{
					Roles: []string{"ORG_OWNER"}}}},
				requiredPerm: "project.read",
				authConfig: &Config{
					RolePermissionMappings: []RoleMapping{
						RoleMapping{
							Role:        "IAM_OWNER",
							Permissions: []string{"project.read"},
						},
						RoleMapping{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
			},
			result: []string{"project.read"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, perms, err := getUserMethodPermissions(tt.args.ctx, tt.args.verifier, tt.args.requiredPerm, tt.args.authConfig)

			if tt.wantErr && err == nil {
				t.Errorf("got wrong result, should get err: actual: %v ", err)
			}

			if tt.wantErr && !tt.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}

			if !tt.wantErr && !equalStringArray(perms, tt.result) {
				t.Errorf("got wrong result, expecting: %v, actual: %v ", tt.result, perms)
			}
		})
	}
}

func Test_MapGrantsToPermissions(t *testing.T) {
	type args struct {
		requiredPerm string
		grants       []*Grant
		authConfig   *Config
	}
	tests := []struct {
		name   string
		args   args
		result []string
	}{
		{
			name: "One Role existing perm",
			args: args{
				requiredPerm: "project.read",
				grants: []*Grant{&Grant{
					Roles: []string{"ORG_OWNER"}}},
				authConfig: &Config{
					RolePermissionMappings: []RoleMapping{
						RoleMapping{
							Role:        "IAM_OWNER",
							Permissions: []string{"project.read"},
						},
						RoleMapping{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
			},
			result: []string{"project.read"},
		},
		{
			name: "One Role not existing perm",
			args: args{
				requiredPerm: "project.write",
				grants: []*Grant{&Grant{
					Roles: []string{"ORG_OWNER"}}},
				authConfig: &Config{
					RolePermissionMappings: []RoleMapping{
						RoleMapping{
							Role:        "IAM_OWNER",
							Permissions: []string{"project.read"},
						},
						RoleMapping{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
			},
			result: []string{},
		},
		{
			name: "Multiple Roles one existing",
			args: args{
				requiredPerm: "project.read",
				grants: []*Grant{&Grant{
					Roles: []string{"ORG_OWNER", "IAM_OWNER"}}},
				authConfig: &Config{
					RolePermissionMappings: []RoleMapping{
						RoleMapping{
							Role:        "IAM_OWNER",
							Permissions: []string{"project.read"},
						},
						RoleMapping{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
			},
			result: []string{"project.read"},
		},
		{
			name: "Multiple Roles, global and specific",
			args: args{
				requiredPerm: "project.read",
				grants: []*Grant{&Grant{
					Roles: []string{"ORG_OWNER", "PROJECT_OWNER:1"}}},
				authConfig: &Config{
					RolePermissionMappings: []RoleMapping{
						RoleMapping{
							Role:        "PROJECT_OWNER",
							Permissions: []string{"project.read"},
						},
						RoleMapping{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
			},
			result: []string{"project.read", "project.read:1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapGrantsToPermissions(tt.args.requiredPerm, tt.args.grants, tt.args.authConfig)
			if !equalStringArray(result, tt.result) {
				t.Errorf("got wrong result, expecting: %v, actual: %v ", tt.result, result)
			}
		})
	}
}

func Test_MapRoleToPerm(t *testing.T) {
	type args struct {
		requiredPerm        string
		actualRole          string
		authConfig          *Config
		resolvedPermissions []string
	}
	tests := []struct {
		name   string
		args   args
		result []string
	}{
		{
			name: "first perm without context id",
			args: args{
				requiredPerm: "project.read",
				actualRole:   "ORG_OWNER",
				authConfig: &Config{
					RolePermissionMappings: []RoleMapping{
						RoleMapping{
							Role:        "IAM_OWNER",
							Permissions: []string{"project.read"},
						},
						RoleMapping{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
				resolvedPermissions: []string{},
			},
			result: []string{"project.read"},
		},
		{
			name: "existing perm without context id",
			args: args{
				requiredPerm: "project.read",
				actualRole:   "ORG_OWNER",
				authConfig: &Config{
					RolePermissionMappings: []RoleMapping{
						RoleMapping{
							Role:        "IAM_OWNER",
							Permissions: []string{"project.read"},
						},
						RoleMapping{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
				resolvedPermissions: []string{"project.read"},
			},
			result: []string{"project.read"},
		},
		{
			name: "first perm with context id",
			args: args{
				requiredPerm: "project.read",
				actualRole:   "PROJECT_OWNER:1",
				authConfig: &Config{
					RolePermissionMappings: []RoleMapping{
						RoleMapping{
							Role:        "PROJECT_OWNER",
							Permissions: []string{"project.read"},
						},
						RoleMapping{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
				resolvedPermissions: []string{},
			},
			result: []string{"project.read:1"},
		},
		{
			name: "perm with context id, existing global",
			args: args{
				requiredPerm: "project.read",
				actualRole:   "PROJECT_OWNER:1",
				authConfig: &Config{
					RolePermissionMappings: []RoleMapping{
						RoleMapping{
							Role:        "PROJECT_OWNER",
							Permissions: []string{"project.read"},
						},
						RoleMapping{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
				resolvedPermissions: []string{"project.read"},
			},
			result: []string{"project.read", "project.read:1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapRoleToPerm(tt.args.requiredPerm, tt.args.actualRole, tt.args.authConfig, tt.args.resolvedPermissions)
			if !equalStringArray(result, tt.result) {
				t.Errorf("got wrong result, expecting: %v, actual: %v ", tt.result, result)
			}
		})
	}
}

func Test_AddRoleContextIDToPerm(t *testing.T) {
	type args struct {
		perm  string
		ctxID string
	}
	tests := []struct {
		name   string
		args   args
		result string
	}{
		{
			name: "with ctx id",
			args: args{
				perm:  "perm1",
				ctxID: "2",
			},
			result: "perm1:2",
		},
		{
			name: "with ctx id",
			args: args{
				perm:  "perm1",
				ctxID: "",
			},
			result: "perm1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addRoleContextIDToPerm(tt.args.perm, tt.args.ctxID)
			if result != tt.result {
				t.Errorf("got wrong result, expecting: %v, actual: %v ", tt.result, result)
			}
		})
	}
}

func Test_ExistisPerm(t *testing.T) {
	type args struct {
		existing []string
		perm     string
	}
	tests := []struct {
		name   string
		args   args
		result bool
	}{
		{
			name: "not existing perm",
			args: args{
				existing: []string{"perm1", "perm2", "perm3"},
				perm:     "perm4",
			},
			result: false,
		},
		{
			name: "existing perm",
			args: args{
				existing: []string{"perm1", "perm2", "perm3"},
				perm:     "perm2",
			},
			result: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := existsPerm(tt.args.existing, tt.args.perm)
			if result != tt.result {
				t.Errorf("got wrong result, expecting: %v, actual: %v ", tt.result, result)
			}
		})
	}
}
