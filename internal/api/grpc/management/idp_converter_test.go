package management

import (
	"testing"

	"github.com/caos/zitadel/internal/test"
	"github.com/caos/zitadel/pkg/grpc/idp"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func Test_addOIDCIDPRequestToDomain(t *testing.T) {
	type args struct {
		req *mgmt_pb.AddOrgOIDCIDPRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "all fields filled",
			args: args{
				req: &mgmt_pb.AddOrgOIDCIDPRequest{
					Name:                  "ZITADEL",
					StylingType:           idp.IDPStylingType_STYLING_TYPE_GOOGLE,
					ClientId:              "test1234",
					ClientSecret:          "test4321",
					Issuer:                "zitadel.ch",
					AuthorizationEndpoint: "https://accounts.zitadel.ch/oauth/v2/authorize",
					TokenEndpoint:         "https://api.zitadel.ch/oauth/v2/token",
					Scopes:                []string{"email", "profile"},
					DisplayNameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
					UsernameMapping:       idp.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addOIDCIDPRequestToDomain(tt.args.req)
			test.AssertFieldsMapped(t, got,
				"ObjectRoot",
				"OIDCConfig.ClientSecret",
				"OIDCConfig.ObjectRoot",
				"OIDCConfig.IDPConfigID",
				"IDPConfigID",
				"State",
				"Type", //TODO: default (0) is oidc
			)
		})
	}
}

func Test_addOIDCIDPRequestToDomainOIDCIDPConfig(t *testing.T) {
	type args struct {
		req *mgmt_pb.AddOrgOIDCIDPRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "all fields filled",
			args: args{
				req: &mgmt_pb.AddOrgOIDCIDPRequest{
					ClientId:              "test1234",
					ClientSecret:          "test4321",
					Issuer:                "zitadel.ch",
					AuthorizationEndpoint: "https://accounts.zitadel.ch/oauth/v2/authorize",
					TokenEndpoint:         "https://api.zitadel.ch/oauth/v2/token",
					Scopes:                []string{"email", "profile"},
					DisplayNameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
					UsernameMapping:       idp.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addOIDCIDPRequestToDomainOIDCIDPConfig(tt.args.req)
			test.AssertFieldsMapped(t, got,
				"ObjectRoot",
				"ClientSecret", //TODO: is client secret string enough for backend?
				"IDPConfigID",
			)
		})
	}
}

func Test_updateIDPToDomain(t *testing.T) {
	type args struct {
		req *mgmt_pb.UpdateOrgIDPRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "all fields filled",
			args: args{
				req: &mgmt_pb.UpdateOrgIDPRequest{
					IdpId:       "13523",
					Name:        "new name",
					StylingType: idp.IDPStylingType_STYLING_TYPE_GOOGLE,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateIDPToDomain(tt.args.req)
			test.AssertFieldsMapped(t, got,
				"ObjectRoot",
				"OIDCConfig",
				"State",
				"Type", //TODO: type should not be changeable
			)
		})
	}
}

func Test_updateOIDCConfigToDomain(t *testing.T) {
	type args struct {
		req *mgmt_pb.UpdateOrgIDPOIDCConfigRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "all fields filled",
			args: args{
				req: &mgmt_pb.UpdateOrgIDPOIDCConfigRequest{
					IdpId:                 "4208",
					Issuer:                "zitadel.ch",
					AuthorizationEndpoint: "https://accounts.zitadel.ch/oauth/v2/authorize",
					TokenEndpoint:         "https://api.zitadel.ch/oauth/v2/token",
					ClientId:              "ZITEADEL",
					ClientSecret:          "i'm so secret",
					Scopes:                []string{"profile"},
					DisplayNameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
					UsernameMapping:       idp.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateOIDCConfigToDomain(tt.args.req)
			test.AssertFieldsMapped(t, got,
				"ObjectRoot",
				"ClientSecret",
			)
		})
	}
}
