package admin

import (
	"testing"

	"github.com/zitadel/zitadel/internal/test"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/idp"
)

func Test_addOIDCIDPRequestToDomain(t *testing.T) {
	type args struct {
		req *admin_pb.AddOIDCIDPRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "all fields filled",
			args: args{
				req: &admin_pb.AddOIDCIDPRequest{
					Name:               "ZITADEL",
					StylingType:        idp.IDPStylingType_STYLING_TYPE_GOOGLE,
					ClientId:           "test1234",
					ClientSecret:       "test4321",
					Issuer:             "zitadel.ch",
					Scopes:             []string{"email", "profile"},
					DisplayNameMapping: idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
					UsernameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME,
					AutoRegister:       true,
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
				"OIDCConfig.AuthorizationEndpoint",
				"OIDCConfig.TokenEndpoint",
				"Type",
				"JWTConfig",
			)
		})
	}
}

func Test_addOIDCIDPRequestToDomainOIDCIDPConfig(t *testing.T) {
	type args struct {
		req *admin_pb.AddOIDCIDPRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "all fields filled",
			args: args{
				req: &admin_pb.AddOIDCIDPRequest{
					ClientId:           "test1234",
					ClientSecret:       "test4321",
					Issuer:             "zitadel.ch",
					Scopes:             []string{"email", "profile"},
					DisplayNameMapping: idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
					UsernameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addOIDCIDPRequestToDomainOIDCIDPConfig(tt.args.req)
			test.AssertFieldsMapped(t, got,
				"ObjectRoot",
				"ClientSecret",
				"IDPConfigID",
				"AuthorizationEndpoint",
				"TokenEndpoint",
			)
		})
	}
}

func Test_updateIDPToDomain(t *testing.T) {
	type args struct {
		req *admin_pb.UpdateIDPRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "all fields filled",
			args: args{
				req: &admin_pb.UpdateIDPRequest{
					IdpId:        "13523",
					Name:         "new name",
					StylingType:  idp.IDPStylingType_STYLING_TYPE_GOOGLE,
					AutoRegister: true,
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
				"JWTConfig",
				"State",
				"Type",
			)
		})
	}
}

func Test_updateOIDCConfigToDomain(t *testing.T) {
	type args struct {
		req *admin_pb.UpdateIDPOIDCConfigRequest
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "all fields filled",
			args: args{
				req: &admin_pb.UpdateIDPOIDCConfigRequest{
					IdpId:              "4208",
					Issuer:             "zitadel.ch",
					ClientId:           "ZITEADEL",
					ClientSecret:       "i'm so secret",
					Scopes:             []string{"profile"},
					DisplayNameMapping: idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
					UsernameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME,
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
				"AuthorizationEndpoint",
				"TokenEndpoint",
			)
		})
	}
}
