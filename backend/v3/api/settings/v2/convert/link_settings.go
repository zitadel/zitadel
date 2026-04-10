package convert

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

func DomainLinksModelToGRPCResponse(models []domain.Link) ([]*settings.Link, error) {
	ls := make([]*settings.Link, len(models), len(models))
	for i, l := range models {
		tp, err := domainLinkTypeToGrpc(l.Type)
		if err != nil {
			return nil, err
		}

		target, err := domainLinkTargetToGrpc(l.Target)
		if err != nil {
			return nil, err
		}

		ls[i] = &settings.Link{
			Type:           tp,
			Url:            l.Url,
			TranslationKey: l.TranslationKey,
			Target:         target,
		}
	}
	return ls, nil
}

func GrpcLinksToDomain(ls []*settings.Link) ([]domain.Link, error) {
	models := make([]domain.Link, len(ls), len(ls))
	for i, l := range ls {
		tp, err := grpcLinkTypeToDomain(l.Type)
		if err != nil {
			return nil, err
		}

		target, err := grpcLinkTargetToDomain(l.Target)
		if err != nil {
			return nil, err
		}

		models[i] = domain.Link{
			Type:           tp,
			Url:            l.GetUrl(),
			TranslationKey: l.GetTranslationKey(),
			Target:         target,
		}
	}
	return models, nil
}

func domainLinkTypeToGrpc(tp domain.LinkType) (settings.LinkType, error) {
	switch tp {
	case domain.LinkTypeUnspecified:
		return settings.LinkType_LINK_TYPE_UNSPECIFIED, nil
	case domain.LinkTypeTermsOfService:
		return settings.LinkType_LINK_TYPE_TERMS_OF_SERVICE, nil
	case domain.LinkTypePrivacyPolicy:
		return settings.LinkType_LINK_TYPE_PRIVACY_POLICY, nil
	case domain.LinkTypeHelp:
		return settings.LinkType_LINK_TYPE_HELP, nil
	case domain.LinkTypeSupport:
		return settings.LinkType_LINK_TYPE_SUPPORT, nil
	case domain.LinkTypeDocs:
		return settings.LinkType_LINK_TYPE_DOCS, nil
	case domain.LinkTypeCustom:
		return settings.LinkType_LINK_TYPE_CUSTOM, nil
	default:
		return settings.LinkType_LINK_TYPE_UNSPECIFIED, zerrors.ThrowInvalidArgumentf(nil, "BCgEF3", "unknown link type %v", tp)
	}
}

func grpcLinkTypeToDomain(tp settings.LinkType) (domain.LinkType, error) {
	switch tp {
	case settings.LinkType_LINK_TYPE_UNSPECIFIED:
		return domain.LinkTypeUnspecified, nil
	case settings.LinkType_LINK_TYPE_TERMS_OF_SERVICE:
		return domain.LinkTypeTermsOfService, nil
	case settings.LinkType_LINK_TYPE_PRIVACY_POLICY:
		return domain.LinkTypePrivacyPolicy, nil
	case settings.LinkType_LINK_TYPE_HELP:
		return domain.LinkTypeHelp, nil
	case settings.LinkType_LINK_TYPE_SUPPORT:
		return domain.LinkTypeSupport, nil
	case settings.LinkType_LINK_TYPE_DOCS:
		return domain.LinkTypeDocs, nil
	case settings.LinkType_LINK_TYPE_CUSTOM:
		return domain.LinkTypeCustom, nil
	default:
		return domain.LinkTypeUnspecified, zerrors.ThrowInvalidArgumentf(nil, "BCgEF3", "unknown link type %v", tp)
	}
}

func domainLinkTargetToGrpc(target domain.LinkTarget) (settings.LinkTarget, error) {
	switch target {
	case domain.LinkTargetUnspecified:
		return settings.LinkTarget_LINK_TARGET_UNSPECIFIED, nil
	case domain.LinkTargetSelf:
		return settings.LinkTarget_LINK_TARGET_SELF, nil
	case domain.LinkTypeBlank:
		return settings.LinkTarget_LINK_TARGET_BLANK, nil
	default:
		return settings.LinkTarget_LINK_TARGET_UNSPECIFIED, zerrors.ThrowInvalidArgumentf(nil, "W35qZF", "unknown target type %v", target)
	}
}

func grpcLinkTargetToDomain(target settings.LinkTarget) (domain.LinkTarget, error) {
	switch target {
	case settings.LinkTarget_LINK_TARGET_UNSPECIFIED:
		return domain.LinkTargetUnspecified, nil
	case settings.LinkTarget_LINK_TARGET_SELF:
		return domain.LinkTargetSelf, nil
	case settings.LinkTarget_LINK_TARGET_BLANK:
		return domain.LinkTypeBlank, nil
	default:
		return domain.LinkTargetUnspecified, zerrors.ThrowInvalidArgumentf(nil, "yJJTA3", "unknown target type %v", target)
	}
}

func DomainSettingsSourceToSourceToGrpc(s domain.SettingsSource) (settings.Source, error) {
	switch s {
	// TODO(wim): add cases
	default:
		return settings.Source_SOURCE_UNSPECIFIED, zerrors.ThrowInvalidArgumentf(nil, "fylcBu", "unknown source %v", s)
	}
}
