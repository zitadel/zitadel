package authz

import (
	"context"
	"testing"

	"github.com/zitadel/zitadel/internal/zerrors"
)

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

type membershipsResolverFunc func(ctx context.Context, orgID string, shouldTriggerBulk bool) ([]*Membership, error)

func (m membershipsResolverFunc) SearchMyMemberships(ctx context.Context, orgID string, shouldTriggerBulk bool) ([]*Membership, error) {
	return m(ctx, orgID, shouldTriggerBulk)
}

func Test_GetUserPermissions(t *testing.T) {
	type args struct {
		ctxData             CtxData
		membershipsResolver MembershipsResolver
		requiredPerm        string
		authConfig          Config
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
				membershipsResolver: membershipsResolverFunc(func(ctx context.Context, orgID string, shouldTriggerBulk bool) ([]*Membership, error) {
					return []*Membership{{Roles: []string{"ORG_OWNER"}}}, nil
				}),
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
			errFunc: zerrors.IsUnauthenticated,
			result:  []string{"project.read"},
		},
		{
			name: "No Grants",
			args: args{
				ctxData: CtxData{},
				membershipsResolver: membershipsResolverFunc(func(ctx context.Context, orgID string, shouldTriggerBulk bool) ([]*Membership, error) {
					return []*Membership{}, nil
				}),
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
				membershipsResolver: membershipsResolverFunc(func(ctx context.Context, orgID string, shouldTriggerBulk bool) ([]*Membership, error) {
					return []*Membership{
						{
							AggregateID: "IAM",
							ObjectID:    "IAM",
							MemberType:  MemberTypeIAM,
							Roles:       []string{"IAM_OWNER"},
						},
					}, nil
				}),
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
			_, perms, err := getUserPermissions(context.Background(), tt.args.membershipsResolver, tt.args.requiredPerm, tt.args.authConfig.RolePermissionMappings, tt.args.ctxData, tt.args.ctxData.OrgID)

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
						MemberType:  MemberTypeOrganization,
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
						MemberType:  MemberTypeOrganization,
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
						MemberType:  MemberTypeOrganization,
						Roles:       []string{"ORG_OWNER"},
					},
					{
						AggregateID: "IAM",
						ObjectID:    "IAM",
						MemberType:  MemberTypeIAM,
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
						MemberType:  MemberTypeOrganization,
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
			requestPerms, allPerms := mapMembershipsToPermissions(tt.args.requiredPerm, tt.args.membership, tt.args.authConfig.RolePermissionMappings)
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
					MemberType:  MemberTypeOrganization,
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
					MemberType:  MemberTypeOrganization,
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
			requestPerms, allPerms := mapMembershipToPerm(tt.args.requiredPerm, tt.args.membership, tt.args.authConfig.RolePermissionMappings, tt.args.requestPerms, tt.args.allPerms)
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

func Test_CheckUserResourcePermissions(t *testing.T) {
	type args struct {
		perms      []string
		resourceID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no permissions",
			args: args{
				perms:      []string{},
				resourceID: "",
			},
			wantErr: true,
		},
		{
			name: "has permission and no context requested",
			args: args{
				perms:      []string{"project.read"},
				resourceID: "",
			},
			wantErr: false,
		},
		{
			name: "context requested and has global permission",
			args: args{
				perms:      []string{"project.read", "project.read:1"},
				resourceID: "Test",
			},
			wantErr: false,
		},
		{
			name: "context requested and has specific permission",
			args: args{
				perms:      []string{"project.read:Test"},
				resourceID: "Test",
			},
			wantErr: false,
		},
		{
			name: "context requested and has no permission",
			args: args{
				perms:      []string{"project.read:Test"},
				resourceID: "Hodor",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkUserResourcePermissions(tt.args.perms, tt.args.resourceID)
			if tt.wantErr && err == nil {
				t.Errorf("got wrong result, should get err: actual: %v ", err)
			}

			if !tt.wantErr && err != nil {
				t.Errorf("shouldn't get err: %v ", err)
			}

			if tt.wantErr && !zerrors.IsPermissionDenied(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func Test_HasContextResourcePermission(t *testing.T) {
	type args struct {
		perms      []string
		resourceID string
	}
	tests := []struct {
		name   string
		args   args
		result bool
	}{
		{
			name: "existing context permission",
			args: args{
				perms:      []string{"test:wrong", "test:right"},
				resourceID: "right",
			},
			result: true,
		},
		{
			name: "not existing context permission",
			args: args{
				perms:      []string{"test:wrong", "test:wrong2"},
				resourceID: "test",
			},
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasContextResourcePermission(tt.args.perms, tt.args.resourceID)
			if result != tt.result {
				t.Errorf("got wrong result, expecting: %v, actual: %v ", tt.result, result)
			}
		})
	}
}
