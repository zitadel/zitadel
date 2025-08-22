package management

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/test"
	"github.com/zitadel/zitadel/pkg/grpc/idp"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
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
			got := AddOIDCIDPRequestToDomain(tt.args.req)
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

func Test_signatureAlgorithmToCommand(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                   string
		signatureAlgorithm     idp.SAMLSignatureAlgorithm
		wantSignatureAlgorithm string
	}{
		{
			name:                   "signature algorithm default value",
			signatureAlgorithm:     11,
			wantSignatureAlgorithm: "",
		},
		{
			name:                   "RSA_SHA1",
			signatureAlgorithm:     idp.SAMLSignatureAlgorithm_SAML_SIGNATURE_RSA_SHA1,
			wantSignatureAlgorithm: "http://www.w3.org/2000/09/xmldsig#rsa-sha1",
		},
		{
			name:                   "RSA_SHA256",
			signatureAlgorithm:     idp.SAMLSignatureAlgorithm_SAML_SIGNATURE_RSA_SHA256,
			wantSignatureAlgorithm: "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
		},
		{
			name:                   "RSA_SHA512",
			signatureAlgorithm:     idp.SAMLSignatureAlgorithm_SAML_SIGNATURE_RSA_SHA512,
			wantSignatureAlgorithm: "http://www.w3.org/2001/04/xmldsig-more#rsa-sha512",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := signatureAlgorithmToCommand(tt.signatureAlgorithm)
			require.Equal(t, tt.wantSignatureAlgorithm, got)
		})
	}
}
