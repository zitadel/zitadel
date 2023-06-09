package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
)

func TestCommands_AllIDPWriteModel(t *testing.T) {
	type args struct {
		resourceOwner string
		instanceBool  bool
		id            string
		idpType       domain.IDPType
	}
	type res struct {
		writeModelType interface{}
		err            error
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "writemodel instance oidc",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeOIDC,
			},
			res: res{
				writeModelType: &InstanceOIDCIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance jwt",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeJWT,
			},
			res: res{
				writeModelType: &InstanceJWTIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance oauth",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeOAuth,
			},
			res: res{
				writeModelType: &InstanceOAuthIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance ldap",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeLDAP,
			},
			res: res{
				writeModelType: &InstanceLDAPIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance azureAD",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeAzureAD,
			},
			res: res{
				writeModelType: &InstanceAzureADIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance github",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeGitHub,
			},
			res: res{
				writeModelType: &InstanceGitHubIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance github enterprise",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeGitHubEnterprise,
			},
			res: res{
				writeModelType: &InstanceGitHubEnterpriseIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance gitlab",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeGitLab,
			},
			res: res{
				writeModelType: &InstanceGitLabIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance gitlab self hosted",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeGitLabSelfHosted,
			},
			res: res{
				writeModelType: &InstanceGitLabSelfHostedIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance google",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeGoogle,
			},
			res: res{
				writeModelType: &InstanceGoogleIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance unspecified",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeUnspecified,
			},
			res: res{
				err: errors.ThrowInternal(nil, "COMMAND-xw921211", "Errors.IDPConfig.NotExisting"),
			},
		},
		{
			name: "writemodel org oidc",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeOIDC,
			},
			res: res{
				writeModelType: &OrgOIDCIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org jwt",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeJWT,
			},
			res: res{
				writeModelType: &OrgJWTIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org oauth",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeOAuth,
			},
			res: res{
				writeModelType: &OrgOAuthIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org ldap",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeLDAP,
			},
			res: res{
				writeModelType: &OrgLDAPIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org azureAD",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeAzureAD,
			},
			res: res{
				writeModelType: &OrgAzureADIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org github",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeGitHub,
			},
			res: res{
				writeModelType: &OrgGitHubIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org github enterprise",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeGitHubEnterprise,
			},
			res: res{
				writeModelType: &OrgGitHubEnterpriseIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org gitlab",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeGitLab,
			},
			res: res{
				writeModelType: &OrgGitLabIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org gitlab self hosted",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeGitLabSelfHosted,
			},
			res: res{
				writeModelType: &OrgGitLabSelfHostedIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org google",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeGoogle,
			},
			res: res{
				writeModelType: &OrgGoogleIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org unspecified",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeUnspecified,
			},
			res: res{
				err: errors.ThrowInternal(nil, "COMMAND-xw921111", "Errors.IDPConfig.NotExisting"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm, err := NewAllIDPWriteModel(tt.args.resourceOwner, tt.args.instanceBool, tt.args.id, tt.args.idpType)
			require.ErrorIs(t, err, tt.res.err)
			if wm != nil {
				assert.IsType(t, tt.res.writeModelType, wm.model)
			}
		})
	}
}
