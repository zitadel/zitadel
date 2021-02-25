package idp

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/v2/domain"
	idp_pb "github.com/caos/zitadel/pkg/grpc/idp"
)

func IDPViewToPb(idp *domain.IDPConfigView) *idp_pb.IDP {
	mapped := &idp_pb.IDP{
		Id:          idp.AggregateID,
		State:       IDPStateToPb(idp.State),
		Name:        idp.Name,
		StylingType: IDPStylingTypeToPb(idp.StylingType),
		Config:      IDPViewToConfigPb(idp),
		Details:     object.ToDetailsPb(idp.Sequence, idp.CreationDate, idp.ChangeDate, "idp.ResourceOwner"), //TODO: resource owner in view
	}
	return mapped
}

func IDPsToUserLinkPb(resp *model.ExternalIDPSearchResponse) []*idp_pb.IDPUserLink {
	links := make([]*idp_pb.IDPUserLink, len(resp.Result))
	for i, link := range resp.Result {
		links[i] = ExternalIDPViewToUserLinkPb(link)
	}
	return links
}

func ExternalIDPViewToUserLinkPb(link *model.ExternalIDPView) *idp_pb.IDPUserLink {
	return &idp_pb.IDPUserLink{
		UserId:           link.UserID,
		IdpId:            link.IDPConfigID,
		IdpName:          link.IDPName,
		ProvidedUserId:   link.ExternalUserID,
		ProvidedUserName: link.UserDisplayName,
		//TODO: as soon as saml is implemented we need to switch here
		IdpType: idp_pb.IDPType_IDP_TYPE_OIDC,
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

func IDPStylingTypeToDomain(stylingType idp_pb.IDPStylingType) domain.IDPConfigStylingType {
	switch stylingType {
	case idp_pb.IDPStylingType_STYLING_TYPE_GOOGLE:
		return domain.IDPConfigStylingTypeGoogle
	default:
		return domain.IDPConfigStylingTypeUnspecified
	}
}

func ModelIDPStylingTypeToPb(stylingType iam_model.IDPStylingType) idp_pb.IDPStylingType {
	switch stylingType {
	case iam_model.IDPStylingTypeGoogle:
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

func ModelIDPViewToConfigPb(config *iam_model.IDPConfigView) *idp_pb.IDP_OidcConfig {
	return &idp_pb.IDP_OidcConfig{
		OidcConfig: &idp_pb.OIDCConfig{
			ClientId:           config.OIDCClientID,
			Issuer:             config.OIDCIssuer,
			Scopes:             config.OIDCScopes,
			DisplayNameMapping: ModelMappingFieldToPb(config.OIDCIDPDisplayNameMapping),
			UsernameMapping:    ModelMappingFieldToPb(config.OIDCUsernameMapping),
		},
	}
}

func IDPViewToConfigPb(config *domain.IDPConfigView) *idp_pb.IDP_OidcConfig {
	return &idp_pb.IDP_OidcConfig{
		OidcConfig: &idp_pb.OIDCConfig{
			ClientId:           config.OIDCClientID,
			Issuer:             config.OIDCIssuer,
			Scopes:             config.OIDCScopes,
			DisplayNameMapping: MappingFieldToPb(config.OIDCIDPDisplayNameMapping),
			UsernameMapping:    MappingFieldToPb(config.OIDCUsernameMapping),
		},
	}
}

func OIDCConfigToPb(config *domain.OIDCIDPConfig) *idp_pb.IDP_OidcConfig {
	return &idp_pb.IDP_OidcConfig{
		OidcConfig: &idp_pb.OIDCConfig{
			ClientId: config.ClientID,
			// ClientSecret:       config.ClientSecretString,
			Issuer:             config.Issuer,
			Scopes:             config.Scopes,
			DisplayNameMapping: MappingFieldToPb(config.IDPDisplayNameMapping),
			UsernameMapping:    MappingFieldToPb(config.UsernameMapping),
		},
	}
}

func ModelMappingFieldToPb(mappingField iam_model.OIDCMappingField) idp_pb.OIDCMappingField {
	switch mappingField {
	case iam_model.OIDCMappingFieldEmail:
		return idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL
	case iam_model.OIDCMappingFieldPreferredLoginName:
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
