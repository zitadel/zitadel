package command

import (
	"context"

	"github.com/caos/logging"

	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

func (r *CommandSide) AddOrgDomain(ctx context.Context, orgDomain *domain.OrgDomain) (*domain.OrgDomain, error) {
	domainWriteModel := NewOrgDomainWriteModel(orgDomain.AggregateID, orgDomain.Domain)
	orgAgg := OrgAggregateFromWriteModel(&domainWriteModel.WriteModel)
	err := r.addOrgDomain(ctx, orgAgg, domainWriteModel, orgDomain)
	if err != nil {
		return nil, err
	}
	err = r.eventstore.PushAggregate(ctx, domainWriteModel, orgAgg)
	if err != nil {
		return nil, err
	}
	return orgDomainWriteModelToOrgDomain(domainWriteModel), nil
}

func (r *CommandSide) GenerateOrgDomainValidation(ctx context.Context, orgDomain *domain.OrgDomain) (token, url string, err error) {
	if orgDomain == nil || !orgDomain.IsValid() {
		return "", "", caos_errs.ThrowPreconditionFailed(nil, "ORG-R24hb", "Errors.Org.InvalidDomain")
	}
	checkType, ok := orgDomain.ValidationType.CheckType()
	if !ok {
		return "", "", caos_errs.ThrowPreconditionFailed(nil, "ORG-Gsw31", "Errors.Org.DomainVerificationTypeInvalid")
	}
	domainWriteModel, err := r.getOrgDomainWriteModel(ctx, orgDomain.AggregateID, orgDomain.Domain)
	if err != nil {
		return "", "", err
	}
	if domainWriteModel.State != domain.OrgDomainStateActive {
		return "", "", caos_errs.ThrowPreconditionFailed(nil, "ORG-AGD31", "Errors.Org.DomainNotOnOrg")
	}
	if domainWriteModel.Verified {
		return "", "", caos_errs.ThrowPreconditionFailed(nil, "ORG-HGw21", "Errors.Org.DomainAlreadyVerified")
	}
	token, err = orgDomain.GenerateVerificationCode(r.domainVerificationGenerator)
	if err != nil {
		return "", "", err
	}
	url, err = http_utils.TokenUrl(orgDomain.Domain, token, checkType)
	if err != nil {
		return "", "", caos_errs.ThrowPreconditionFailed(err, "ORG-Bae21", "Errors.Org.DomainVerificationTypeInvalid")
	}

	orgAgg := OrgAggregateFromWriteModel(&domainWriteModel.WriteModel)
	orgAgg.PushEvents(org.NewDomainVerificationAddedEvent(ctx, orgDomain.Domain, orgDomain.ValidationType, orgDomain.ValidationCode))

	err = r.eventstore.PushAggregate(ctx, domainWriteModel, orgAgg)
	if err != nil {
		return "", "", err
	}
	return token, url, nil
}

func (r *CommandSide) ValidateOrgDomain(ctx context.Context, orgDomain *domain.OrgDomain) error {
	if orgDomain == nil || !orgDomain.IsValid() {
		return caos_errs.ThrowPreconditionFailed(nil, "ORG-R24hb", "Errors.Org.InvalidDomain")
	}
	domainWriteModel, err := r.getOrgDomainWriteModel(ctx, orgDomain.AggregateID, orgDomain.Domain)
	if err != nil {
		return err
	}
	if domainWriteModel.State != domain.OrgDomainStateActive {
		return caos_errs.ThrowPreconditionFailed(nil, "ORG-Sjdi3", "Errors.Org.DomainNotOnOrg")
	}
	if domainWriteModel.Verified {
		return caos_errs.ThrowPreconditionFailed(nil, "ORG-HGw21", "Errors.Org.DomainAlreadyVerified")
	}
	if domainWriteModel.ValidationCode == nil || domainWriteModel.ValidationType == domain.OrgDomainValidationTypeUnspecified {
		return caos_errs.ThrowPreconditionFailed(nil, "ORG-SFBB3", "Errors.Org.DomainVerificationMissing")
	}

	validationCode, err := crypto.DecryptString(domainWriteModel.ValidationCode, r.domainVerificationAlg)
	if err != nil {
		return err
	}
	checkType, _ := domainWriteModel.ValidationType.CheckType()
	err = r.domainVerificationValidator(domainWriteModel.Domain, validationCode, validationCode, checkType)
	orgAgg := OrgAggregateFromWriteModel(&domainWriteModel.WriteModel)
	if err == nil {
		orgAgg.PushEvents(org.NewDomainVerifiedEvent(ctx, orgDomain.Domain))
		return r.eventstore.PushAggregate(ctx, domainWriteModel, orgAgg)
	}
	orgAgg.PushEvents(org.NewDomainVerificationFailedEvent(ctx, orgDomain.Domain))
	err = r.eventstore.PushAggregate(ctx, domainWriteModel, orgAgg)
	logging.LogWithFields("ORG-dhTE", "orgID", orgAgg.ID(), "domain", orgDomain.Domain).OnError(err).Error("NewDomainVerificationFailedEvent push failed")
	return caos_errs.ThrowInvalidArgument(err, "ORG-GH3s", "Errors.Org.DomainVerificationFailed")
}

func (r *CommandSide) SetPrimaryOrgDomain(ctx context.Context, orgDomain *domain.OrgDomain) error {
	if orgDomain == nil || !orgDomain.IsValid() {
		return caos_errs.ThrowPreconditionFailed(nil, "ORG-SsDG2", "Errors.Org.InvalidDomain")
	}
	domainWriteModel, err := r.getOrgDomainWriteModel(ctx, orgDomain.AggregateID, orgDomain.Domain)
	if err != nil {
		return err
	}
	if domainWriteModel.State != domain.OrgDomainStateActive {
		return caos_errs.ThrowPreconditionFailed(nil, "ORG-GDfA3", "Errors.Org.DomainNotOnOrg")
	}
	if !domainWriteModel.Verified {
		return caos_errs.ThrowPreconditionFailed(nil, "ORG-Ggd32", "Errors.Org.DomainNotVerified")
	}
	orgAgg := OrgAggregateFromWriteModel(&domainWriteModel.WriteModel)
	orgAgg.PushEvents(org.NewDomainPrimarySetEvent(ctx, orgDomain.Domain))
	return r.eventstore.PushAggregate(ctx, domainWriteModel, orgAgg)
}

func (r *CommandSide) RemoveOrgDomain(ctx context.Context, orgDomain *domain.OrgDomain) error {
	if orgDomain == nil || !orgDomain.IsValid() {
		return caos_errs.ThrowPreconditionFailed(nil, "ORG-SJsK3", "Errors.Org.InvalidDomain")
	}
	domainWriteModel, err := r.getOrgDomainWriteModel(ctx, orgDomain.AggregateID, orgDomain.Domain)
	if err != nil {
		return err
	}
	if domainWriteModel.State != domain.OrgDomainStateActive {
		return caos_errs.ThrowPreconditionFailed(nil, "ORG-GDfA3", "Errors.Org.DomainNotOnOrg")
	}
	if domainWriteModel.Primary {
		return caos_errs.ThrowPreconditionFailed(nil, "ORG-Sjdi3", "Errors.Org.PrimaryDomainNotDeletable")
	}
	orgAgg := OrgAggregateFromWriteModel(&domainWriteModel.WriteModel)
	orgAgg.PushEvents(org.NewDomainRemovedEvent(ctx, orgDomain.Domain))
	return r.eventstore.PushAggregate(ctx, domainWriteModel, orgAgg)
}

func (r *CommandSide) addOrgDomain(ctx context.Context, orgAgg *org.Aggregate, addedDomain *OrgDomainWriteModel, orgDomain *domain.OrgDomain) error {
	err := r.eventstore.FilterToQueryReducer(ctx, addedDomain)
	if err != nil {
		return err
	}
	if addedDomain.State == domain.OrgDomainStateActive {
		return caos_errs.ThrowAlreadyExists(nil, "COMMA-Bd2jj", "Errors.Org.Domain.AlreadyExists")
	}
	orgAgg.PushEvents(org.NewDomainAddedEvent(ctx, orgDomain.Domain))
	if orgDomain.Verified {
		//TODO: uniqueness verified domain
		//TODO: users with verified domain -> domain claimed
		orgAgg.PushEvents(org.NewDomainVerifiedEvent(ctx, orgDomain.Domain))
	}
	if orgDomain.Primary {
		orgAgg.PushEvents(org.NewDomainPrimarySetEvent(ctx, orgDomain.Domain))
	}
	return nil
}

func (r *CommandSide) getOrgDomainWriteModel(ctx context.Context, orgID, domain string) (*OrgDomainWriteModel, error) {
	domainWriteModel := NewOrgDomainWriteModel(orgID, domain)
	err := r.eventstore.FilterToQueryReducer(ctx, domainWriteModel)
	if err != nil {
		return nil, err
	}
	return domainWriteModel, nil
}
