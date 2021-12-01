package idp

import (
	obj_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/query"
	idp_pb "github.com/caos/zitadel/pkg/grpc/idp"
)

func IDPViewsToPb(idps []*query.IDP) []*idp_pb.IDP {
	resp := make([]*idp_pb.IDP, len(idps))
	for i, idp := range idps {
		resp[i] = ModelIDPViewToPb(idp)
	}
	return resp
}

func ModelIDPViewToPb(idp *query.IDP) *idp_pb.IDP {
	return &idp_pb.IDP{
		Id:           idp.ID,
		State:        ModelIDPStateToPb(idp.State),
		Name:         idp.Name,
		StylingType:  ModelIDPStylingTypeToPb(idp.StylingType),
		AutoRegister: idp.AutoRegister,
		Owner:        ModelIDPProviderTypeToPb(idp.OwnerType),
		Config:       ModelIDPViewToConfigPb(idp),
		Details: obj_grpc.ToViewDetailsPb(
			idp.Sequence,
			idp.CreationDate,
			idp.ChangeDate,
			idp.ID,
		),
	}
}

func IDPViewToPb(idp *query.IDP) *idp_pb.IDP {
	mapped := &idp_pb.IDP{
		Owner:        ownerTypeToPB(idp.OwnerType),
		Id:           idp.ID,
		State:        IDPStateToPb(idp.State),
		Name:         idp.Name,
		StylingType:  IDPStylingTypeToPb(idp.StylingType),
		AutoRegister: idp.AutoRegister,
		Config:       IDPViewToConfigPb(idp),
		Details:      obj_grpc.ToViewDetailsPb(idp.Sequence, idp.CreationDate, idp.ChangeDate, idp.ID),
	}
	return mapped
}

func IDPLoginPolicyLinksToPb(links []*query.IDPLoginPolicyLink) []*idp_pb.IDPLoginPolicyLink {
	l := make([]*idp_pb.IDPLoginPolicyLink, len(links))
	for i, link := range links {
		l[i] = IDPLoginPolicyLinkToPb(link)
	}
	return l
}

func IDPLoginPolicyLinkToPb(link *query.IDPLoginPolicyLink) *idp_pb.IDPLoginPolicyLink {
	return &idp_pb.IDPLoginPolicyLink{
		IdpId:   link.IDPID,
		IdpName: link.IDPName,
		IdpType: IDPTypeToPb(link.IDPType),
	}
}

func IDPUserLinksToPb(res []*query.UserIDPLink) []*idp_pb.IDPUserLink {
	links := make([]*idp_pb.IDPUserLink, len(res))
	for i, link := range res {
		links[i] = IDPUserLinkToPb(link)
	}
	return links
}

func IDPUserLinkToPb(link *query.UserIDPLink) *idp_pb.IDPUserLink {
	return &idp_pb.IDPUserLink{
		UserId:           link.UserID,
		IdpId:            link.IDPID,
		IdpName:          link.IDPName,
		ProvidedUserId:   link.ProvidedUserID,
		ProvidedUserName: link.ProvidedUsername,
		//TODO: as soon as saml is implemented we need to switch here
		//IdpType: IDPTypeToPb(link.Type),
	}
}

func IDPTypeToPb(idpType domain.IDPConfigType) idp_pb.IDPType {
	switch idpType {
	case domain.IDPConfigTypeOIDC:
		return idp_pb.IDPType_IDP_TYPE_OIDC
	case domain.IDPConfigTypeSAML:
		return idp_pb.IDPType_IDP_TYPE_UNSPECIFIED
	case domain.IDPConfigTypeJWT:
		return idp_pb.IDPType_IDP_TYPE_JWT
	default:
		return idp_pb.IDPType_IDP_TYPE_UNSPECIFIED
	}
}

func IDPStateToPb(state domain.IDPConfigState) idp_pb.IDPState {
	switch state {
	case domain.IDPConfigStateActive:
		return idp_pb.IDPState_IDP_STATE_ACTIVE
	case domain.IDPConfigStateInactive:
		return idp_pb.IDPState_IDP_STATE_INACTIVE
	default:
		return idp_pb.IDPState_IDP_STATE_UNSPECIFIED
	}
}

func ModelIDPStateToPb(state domain.IDPConfigState) idp_pb.IDPState {
	switch state {
	case domain.IDPConfigStateActive:
		return idp_pb.IDPState_IDP_STATE_ACTIVE
	case domain.IDPConfigStateInactive:
		return idp_pb.IDPState_IDP_STATE_INACTIVE
	default:
		return idp_pb.IDPState_IDP_STATE_UNSPECIFIED
	}
}

func IDPStylingTypeToDomain(stylingType idp_pb.IDPStylingType) domain.IDPConfigStylingType {
	switch stylingType {
	case idp_pb.IDPStylingType_STYLING_TYPE_GOOGLE:
		return domain.IDPConfigStylingTypeGoogle
	default:
		return domain.IDPConfigStylingTypeUnspecified
	}
}

func ModelIDPStylingTypeToPb(stylingType domain.IDPConfigStylingType) idp_pb.IDPStylingType {
	switch stylingType {
	case domain.IDPConfigStylingTypeGoogle:
		return idp_pb.IDPStylingType_STYLING_TYPE_GOOGLE
	default:
		return idp_pb.IDPStylingType_STYLING_TYPE_UNSPECIFIED
	}
}

func IDPStylingTypeToPb(stylingType domain.IDPConfigStylingType) idp_pb.IDPStylingType {
	switch stylingType {
	case domain.IDPConfigStylingTypeGoogle:
		return idp_pb.IDPStylingType_STYLING_TYPE_GOOGLE
	default:
		return idp_pb.IDPStylingType_STYLING_TYPE_UNSPECIFIED
	}
}

func ModelIDPViewToConfigPb(config *query.IDP) idp_pb.IDPConfig {
	if config.OIDCIDP != nil {
		return &idp_pb.IDP_OidcConfig{
			OidcConfig: &idp_pb.OIDCConfig{
				ClientId:           config.ClientID,
				Issuer:             config.OIDCIDP.Issuer,
				Scopes:             config.Scopes,
				DisplayNameMapping: ModelMappingFieldToPb(config.DisplayNameMapping),
				UsernameMapping:    ModelMappingFieldToPb(config.UsernameMapping),
			},
		}
	}
	return &idp_pb.IDP_JwtConfig{
		JwtConfig: &idp_pb.JWTConfig{
			JwtEndpoint:  config.Endpoint,
			Issuer:       config.JWTIDP.Issuer,
			KeysEndpoint: config.KeysEndpoint,
			HeaderName:   config.HeaderName,
		},
	}
}

func IDPViewToConfigPb(config *query.IDP) idp_pb.IDPConfig {
	if config.OIDCIDP != nil {
		return &idp_pb.IDP_OidcConfig{
			OidcConfig: &idp_pb.OIDCConfig{
				ClientId:           config.ClientID,
				Issuer:             config.OIDCIDP.Issuer,
				Scopes:             config.Scopes,
				DisplayNameMapping: MappingFieldToPb(config.DisplayNameMapping),
				UsernameMapping:    MappingFieldToPb(config.UsernameMapping),
			},
		}
	}
	return &idp_pb.IDP_JwtConfig{
		JwtConfig: &idp_pb.JWTConfig{
			JwtEndpoint:  config.JWTIDP.Endpoint,
			Issuer:       config.JWTIDP.Issuer,
			KeysEndpoint: config.JWTIDP.KeysEndpoint,
		},
	}
}

func FieldNameToModel(fieldName idp_pb.IDPFieldName) query.Column {
	switch fieldName {
	// case admin.IdpSearchKey_IDPSEARCHKEY_IDP_CONFIG_ID: //TODO: not implemented in proto
	// 	return iam_model.IDPConfigSearchKeyIdpConfigID
	case idp_pb.IDPFieldName_IDP_FIELD_NAME_NAME:
		return query.IDPNameCol
	default:
		return query.Column{}
	}
}

func ModelMappingFieldToPb(mappingField domain.OIDCMappingField) idp_pb.OIDCMappingField {
	switch mappingField {
	case domain.OIDCMappingFieldEmail:
		return idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL
	case domain.OIDCMappingFieldPreferredLoginName:
		return idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME
	default:
		return idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_UNSPECIFIED
	}
}

func MappingFieldToPb(mappingField domain.OIDCMappingField) idp_pb.OIDCMappingField {
	switch mappingField {
	case domain.OIDCMappingFieldEmail:
		return idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL
	case domain.OIDCMappingFieldPreferredLoginName:
		return idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME
	default:
		return idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_UNSPECIFIED
	}
}

func MappingFieldToDomain(mappingField idp_pb.OIDCMappingField) domain.OIDCMappingField {
	switch mappingField {
	case idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL:
		return domain.OIDCMappingFieldEmail
	case idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME:
		return domain.OIDCMappingFieldPreferredLoginName
	default:
		return domain.OIDCMappingFieldUnspecified
	}
}

func ModelIDPProviderTypeToPb(typ domain.IdentityProviderType) idp_pb.IDPOwnerType {
	switch typ {
	case domain.IdentityProviderTypeOrg:
		return idp_pb.IDPOwnerType_IDP_OWNER_TYPE_ORG
	case domain.IdentityProviderTypeSystem:
		return idp_pb.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM
	default:
		return idp_pb.IDPOwnerType_IDP_OWNER_TYPE_UNSPECIFIED
	}
}

func IDPProviderTypeFromPb(typ idp_pb.IDPOwnerType) domain.IdentityProviderType {
	switch typ {
	case idp_pb.IDPOwnerType_IDP_OWNER_TYPE_ORG:
		return domain.IdentityProviderTypeOrg
	case idp_pb.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM:
		return domain.IdentityProviderTypeSystem
	default:
		return domain.IdentityProviderTypeOrg
	}
}

func IDPProviderTypeModelFromPb(typ idp_pb.IDPOwnerType) iam_model.IDPProviderType {
	switch typ {
	case idp_pb.IDPOwnerType_IDP_OWNER_TYPE_ORG:
		return iam_model.IDPProviderTypeOrg
	case idp_pb.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM:
		return iam_model.IDPProviderTypeSystem
	default:
		return iam_model.IDPProviderTypeOrg
	}
}

func IDPIDQueryToModel(query *idp_pb.IDPIDQuery) *iam_model.IDPConfigSearchQuery {
	return &iam_model.IDPConfigSearchQuery{
		Key:    iam_model.IDPConfigSearchKeyIdpConfigID,
		Method: domain.SearchMethodEquals,
		Value:  query.Id,
	}
}

func IDPNameQueryToModel(query *idp_pb.IDPNameQuery) *iam_model.IDPConfigSearchQuery {
	return &iam_model.IDPConfigSearchQuery{
		Key:    iam_model.IDPConfigSearchKeyName,
		Method: obj_grpc.TextMethodToModel(query.Method),
		Value:  query.Name,
	}
}

func IDPOwnerTypeQueryToModel(query *idp_pb.IDPOwnerTypeQuery) *iam_model.IDPConfigSearchQuery {
	return &iam_model.IDPConfigSearchQuery{
		Key:    iam_model.IDPConfigSearchKeyIdpProviderType,
		Method: domain.SearchMethodEquals,
		Value:  IDPProviderTypeModelFromPb(query.OwnerType),
	}
}
func ownerTypeToPB(typ domain.IdentityProviderType) idp_pb.IDPOwnerType {
	switch typ {
	case domain.IdentityProviderTypeOrg:
		return idp_pb.IDPOwnerType_IDP_OWNER_TYPE_ORG
	case domain.IdentityProviderTypeSystem:
		return idp_pb.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM
	default:
		return idp_pb.IDPOwnerType_IDP_OWNER_TYPE_UNSPECIFIED
	}
}
