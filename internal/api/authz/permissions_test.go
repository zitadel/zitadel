package authz

import (
	"context"
	"testing"

	caos_errs "github.com/caos/zitadel/internal/errors"
)

func getTestCtx(userID, orgID string) context.Context {
	return context.WithValue(context.Background(), dataKey, CtxData{UserID: userID, OrgID: orgID})
}

type testVerifier struct {
	memberships []*Membership
}

func (v *testVerifier) VerifyAccessToken(ctx context.Context, token, clientID string) (string, string, string, string, error) {
	return "userID", "agentID", "de", "orgID", nil
}
func (v *testVerifier) SearchMyMemberships(ctx context.Context) ([]*Membership, error) {
	return v.memberships, nil
}

func (v *testVerifier) ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (string, []string, error) {
	return "", nil, nil
}

func (v *testVerifier) ExistsOrg(ctx context.Context, orgID string) error {
	return nil
}

func (v *testVerifier) VerifierClientID(ctx context.Context, appName string) (string, error) {
	return "clientID", nil
}

func (v *testVerifier) CheckOrgFeatures(context.Context, string, ...string) error {
	return nil
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
		ctxData      CtxData
		verifier     *TokenVerifier
		requiredPerm string
		authConfig   Config
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
				ctxData: CtxData{},
				verifier: Start(&testVerifier{memberships: []*Membership{
					{
						Roles: []string{"ORG_OWNER"},
					},
				}}),
				requiredPerm: "project.read",
				authConfig: Config{
					RolePermissionMappings: []RoleMapping{
						{
							Role:        "IAM_OWNER",
							Permissions: []string{"project.read"},
						},
						{
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
				ctxData:      CtxData{},
				verifier:     Start(&testVerifier{memberships: []*Membership{}}),
				requiredPerm: "project.read",
				authConfig: Config{
					RolePermissionMappings: []RoleMapping{
						{
							Role:        "IAM_OWNER",
							Permissions: []string{"project.read"},
						},
						{
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
				ctxData: CtxData{UserID: "userID", OrgID: "orgID"},
				verifier: Start(&testVerifier{memberships: []*Membership{
					{
						AggregateID: "IAM",
						ObjectID:    "IAM",
						MemberType:  MemberTypeIam,
						Roles:       []string{"IAM_OWNER"},
					},
				}}),
				requiredPerm: "project.read",
				authConfig: Config{
					RolePermissionMappings: []RoleMapping{
						{
							Role:        "IAM_OWNER",
							Permissions: []string{"project.read"},
						},
						{
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
			_, perms, err := getUserMethodPermissions(context.Background(), tt.args.verifier, tt.args.requiredPerm, tt.args.authConfig, tt.args.ctxData)

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

func Test_MapMembershipToPermissions(t *testing.T) {
	type args struct {
		requiredPerm string
		membership   []*Membership
		authConfig   Config
	}
	tests := []struct {
		name         string
		args         args
		requestPerms []string
		allPerms     []string
	}{
		{
			name: "One Role existing perm",
			args: args{
				requiredPerm: "project.read",
				membership: []*Membership{
					{
						AggregateID: "1",
						ObjectID:    "1",
						MemberType:  MemberTypeOrganisation,
						Roles:       []string{"ORG_OWNER"},
					},
				},
				authConfig: Config{
					RolePermissionMappings: []RoleMapping{
						{
							Role:        "IAM_OWNER",
							Permissions: []string{"project.read"},
						},
						{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
			},
			requestPerms: []string{"project.read"},
			allPerms:     []string{"org.read", "project.read"},
		},
		{
			name: "One Role not existing perm",
			args: args{
				requiredPerm: "project.write",
				membership: []*Membership{
					{
						AggregateID: "1",
						ObjectID:    "1",
						MemberType:  MemberTypeOrganisation,
						Roles:       []string{"ORG_OWNER"},
					},
				},
				authConfig: Config{
					RolePermissionMappings: []RoleMapping{
						{
							Role:        "IAM_OWNER",
							Permissions: []string{"project.read"},
						},
						{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
			},
			requestPerms: []string{},
			allPerms:     []string{"org.read", "project.read"},
		},
		{
			name: "Multiple Roles one existing",
			args: args{
				requiredPerm: "project.read",
				membership: []*Membership{
					{
						AggregateID: "1",
						ObjectID:    "1",
						MemberType:  MemberTypeOrganisation,
						Roles:       []string{"ORG_OWNER"},
					},
					{
						AggregateID: "IAM",
						ObjectID:    "IAM",
						MemberType:  MemberTypeIam,
						Roles:       []string{"IAM_OWNER"},
					},
				},
				authConfig: Config{
					RolePermissionMappings: []RoleMapping{
						{
							Role:        "IAM_OWNER",
							Permissions: []string{"project.read"},
						},
						{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
			},
			requestPerms: []string{"project.read"},
			allPerms:     []string{"org.read", "project.read"},
		},
		{
			name: "Multiple Roles, global and specific",
			args: args{
				requiredPerm: "project.read",
				membership: []*Membership{
					{
						AggregateID: "2",
						ObjectID:    "2",
						MemberType:  MemberTypeOrganisation,
						Roles:       []string{"ORG_OWNER"},
					},
					{
						AggregateID: "1",
						ObjectID:    "1",
						MemberType:  MemberTypeProject,
						Roles:       []string{"PROJECT_OWNER"},
					},
				},
				authConfig: Config{
					RolePermissionMappings: []RoleMapping{
						{
							Role:        "PROJECT_OWNER",
							Permissions: []string{"project.read"},
						},
						{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
			},
			requestPerms: []string{"project.read", "project.read:1"},
			allPerms:     []string{"org.read", "project.read", "project.read:1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestPerms, allPerms := mapMembershipsToPermissions(tt.args.requiredPerm, tt.args.membership, tt.args.authConfig)
			if !equalStringArray(requestPerms, tt.requestPerms) {
				t.Errorf("got wrong requestPerms, expecting: %v, actual: %v ", tt.requestPerms, requestPerms)
			}
			if !equalStringArray(allPerms, tt.allPerms) {
				t.Errorf("got wrong allPerms, expecting: %v, actual: %v ", tt.allPerms, allPerms)
			}
		})
	}
}

func Test_MapMembershipToPerm(t *testing.T) {
	type args struct {
		requiredPerm string
		membership   *Membership
		authConfig   Config
		requestPerms []string
		allPerms     []string
	}
	tests := []struct {
		name         string
		args         args
		requestPerms []string
		allPerms     []string
	}{
		{
			name: "first perm without context id",
			args: args{
				requiredPerm: "project.read",
				membership: &Membership{
					AggregateID: "Org",
					ObjectID:    "Org",
					MemberType:  MemberTypeOrganisation,
					Roles:       []string{"ORG_OWNER"},
				},
				authConfig: Config{
					RolePermissionMappings: []RoleMapping{
						{
							Role:        "IAM_OWNER",
							Permissions: []string{"project.read"},
						},
						{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
				requestPerms: []string{},
				allPerms:     []string{},
			},
			requestPerms: []string{"project.read"},
			allPerms:     []string{"org.read", "project.read"},
		},
		{
			name: "existing perm without context id",
			args: args{
				requiredPerm: "project.read",
				membership: &Membership{
					AggregateID: "Org",
					ObjectID:    "Org",
					MemberType:  MemberTypeOrganisation,
					Roles:       []string{"ORG_OWNER"},
				},
				authConfig: Config{
					RolePermissionMappings: []RoleMapping{
						{
							Role:        "IAM_OWNER",
							Permissions: []string{"project.read"},
						},
						{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
				requestPerms: []string{"project.read"},
				allPerms:     []string{"org.read", "project.read"},
			},
			requestPerms: []string{"project.read"},
			allPerms:     []string{"org.read", "project.read"},
		},
		{
			name: "first perm with context id",
			args: args{
				requiredPerm: "project.read",
				membership: &Membership{
					AggregateID: "1",
					ObjectID:    "1",
					MemberType:  MemberTypeProject,
					Roles:       []string{"PROJECT_OWNER"},
				},
				authConfig: Config{
					RolePermissionMappings: []RoleMapping{
						{
							Role:        "PROJECT_OWNER",
							Permissions: []string{"project.read"},
						},
						{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
				requestPerms: []string{},
				allPerms:     []string{},
			},
			requestPerms: []string{"project.read:1"},
			allPerms:     []string{"project.read:1"},
		},
		{
			name: "perm with context id, existing global",
			args: args{
				requiredPerm: "project.read",
				membership: &Membership{
					AggregateID: "1",
					ObjectID:    "1",
					MemberType:  MemberTypeProject,
					Roles:       []string{"PROJECT_OWNER"},
				},
				authConfig: Config{
					RolePermissionMappings: []RoleMapping{
						{
							Role:        "PROJECT_OWNER",
							Permissions: []string{"project.read"},
						},
						{
							Role:        "ORG_OWNER",
							Permissions: []string{"org.read", "project.read"},
						},
					},
				},
				requestPerms: []string{"project.read"},
				allPerms:     []string{"org.read", "project.read"},
			},
			requestPerms: []string{"project.read", "project.read:1"},
			allPerms:     []string{"org.read", "project.read", "project.read:1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestPerms, allPerms := mapMembershipToPerm(tt.args.requiredPerm, tt.args.membership, tt.args.authConfig, tt.args.requestPerms, tt.args.allPerms)
			if !equalStringArray(requestPerms, tt.requestPerms) {
				t.Errorf("got wrong requestPerms, expecting: %v, actual: %v ", tt.requestPerms, requestPerms)
			}
			if !equalStringArray(allPerms, tt.allPerms) {
				t.Errorf("got wrong allPerms, expecting: %v, actual: %v ", tt.allPerms, allPerms)
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
		existingPermissions []string
		perm                string
	}
	tests := []struct {
		name   string
		args   args
		result bool
	}{
		{
			name: "not existing perm",
			args: args{
				existingPermissions: []string{"perm1", "perm2", "perm3"},
				perm:                "perm4",
			},
			result: false,
		},
		{
			name: "existing perm",
			args: args{
				existingPermissions: []string{"perm1", "perm2", "perm3"},
				perm:                "perm2",
			},
			result: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExistsPerm(tt.args.existingPermissions, tt.args.perm)
			if result != tt.result {
				t.Errorf("got wrong result, expecting: %v, actual: %v ", tt.result, result)
			}
		})
	}
}
