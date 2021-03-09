package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"

	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/org"
)

func (c *Commands) AddOrgDomain(ctx context.Context, orgDomain *domain.OrgDomain) (*domain.OrgDomain, error) {
	domainWriteModel := NewOrgDomainWriteModel(orgDomain.AggregateID, orgDomain.Domain)
	orgAgg := OrgAggregateFromWriteModel(&domainWriteModel.WriteModel)
	events, err := c.addOrgDomain(ctx, orgAgg, domainWriteModel, orgDomain)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(domainWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return orgDomainWriteModelToOrgDomain(domainWriteModel), nil
}

func (c *Commands) GenerateOrgDomainValidation(ctx context.Context, orgDomain *domain.OrgDomain) (token, url string, err error) {
	if orgDomain == nil || !orgDomain.IsValid() {
		return "", "", caos_errs.ThrowPreconditionFailed(nil, "ORG-R24hb", "Errors.Org.InvalidDomain")
	}
	checkType, ok := orgDomain.ValidationType.CheckType()
	if !ok {
		return "", "", caos_errs.ThrowPreconditionFailed(nil, "ORG-Gsw31", "Errors.Org.DomainVerificationTypeInvalid")
	}
	domainWriteModel, err := c.getOrgDomainWriteModel(ctx, orgDomain.AggregateID, orgDomain.Domain)
	if err != nil {
		return "", "", err
	}
	if domainWriteModel.State != domain.OrgDomainStateActive {
		return "", "", caos_errs.ThrowPreconditionFailed(nil, "ORG-AGD31", "Errors.Org.DomainNotOnOrg")
	}
	if domainWriteModel.Verified {
		return "", "", caos_errs.ThrowPreconditionFailed(nil, "ORG-HGw21", "Errors.Org.DomainAlreadyVerified")
	}
	token, err = orgDomain.GenerateVerificationCode(c.domainVerificationGenerator)
	if err != nil {
		return "", "", err
	}
	url, err = http_utils.TokenUrl(orgDomain.Domain, token, checkType)
	if err != nil {
		return "", "", caos_errs.ThrowPreconditionFailed(err, "ORG-Bae21", "Errors.Org.DomainVerificationTypeInvalid")
	}

	orgAgg := OrgAggregateFromWriteModel(&domainWriteModel.WriteModel)

	_, err = c.eventstore.PushEvents(
		ctx,
		org.NewDomainVerificationAddedEvent(ctx, orgAgg, orgDomain.Domain, orgDomain.ValidationType, orgDomain.ValidationCode))
	if err != nil {
		return "", "", err
	}
	return token, url, nil
}

func (c *Commands) ValidateOrgDomain(ctx context.Context, orgDomain *domain.OrgDomain, claimedUserIDs ...string) (*domain.ObjectDetails, error) {
	if orgDomain == nil || !orgDomain.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-R24hb", "Errors.Org.InvalidDomain")
	}
	domainWriteModel, err := c.getOrgDomainWriteModel(ctx, orgDomain.AggregateID, orgDomain.Domain)
	if err != nil {
		return nil, err
	}
	if domainWriteModel.State != domain.OrgDomainStateActive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-Sjdi3", "Errors.Org.DomainNotOnOrg")
	}
	if domainWriteModel.Verified {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-HGw21", "Errors.Org.DomainAlreadyVerified")
	}
	if domainWriteModel.ValidationCode == nil || domainWriteModel.ValidationType == domain.OrgDomainValidationTypeUnspecified {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-SFBB3", "Errors.Org.DomainVerificationMissing")
	}

	validationCode, err := crypto.DecryptString(domainWriteModel.ValidationCode, c.domainVerificationAlg)
	if err != nil {
		return nil, err
	}
	checkType, _ := domainWriteModel.ValidationType.CheckType()
	err = c.domainVerificationValidator(domainWriteModel.Domain, validationCode, validationCode, checkType)
	orgAgg := OrgAggregateFromWriteModel(&domainWriteModel.WriteModel)
	var events []eventstore.EventPusher
	if err == nil {
		events = append(events, org.NewDomainVerifiedEvent(ctx, orgAgg, orgDomain.Domain))

		for _, userID := range claimedUserIDs {
			userEvents, _, err := c.userDomainClaimed(ctx, userID)
			if err != nil {
				logging.LogWithFields("COMMAND-5m8fs", "userid", userID).WithError(err).Warn("could not claim user")
				continue
			}
			events = append(events, userEvents...)
		}
		pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
		if err != nil {
			return nil, err
		}
		err = AppendAndReduce(domainWriteModel, pushedEvents...)
		if err != nil {
			return nil, err
		}
		return writeModelToObjectDetails(&domainWriteModel.WriteModel), nil
	}
	events = append(events, org.NewDomainVerificationFailedEvent(ctx, orgAgg, orgDomain.Domain))
	_, err = c.eventstore.PushEvents(ctx, events...)
	logging.LogWithFields("ORG-dhTE", "orgID", orgAgg.ID, "domain", orgDomain.Domain).OnError(err).Error("NewDomainVerificationFailedEvent push failed")
	return nil, caos_errs.ThrowInvalidArgument(err, "ORG-GH3s", "Errors.Org.DomainVerificationFailed")
}

func (c *Commands) SetPrimaryOrgDomain(ctx context.Context, orgDomain *domain.OrgDomain) (*domain.ObjectDetails, error) {
	if orgDomain == nil || !orgDomain.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-SsDG2", "Errors.Org.InvalidDomain")
	}
	domainWriteModel, err := c.getOrgDomainWriteModel(ctx, orgDomain.AggregateID, orgDomain.Domain)
	if err != nil {
		return nil, err
	}
	if domainWriteModel.State != domain.OrgDomainStateActive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-GDfA3", "Errors.Org.DomainNotOnOrg")
	}
	if !domainWriteModel.Verified {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-Ggd32", "Errors.Org.DomainNotVerified")
	}
	orgAgg := OrgAggregateFromWriteModel(&domainWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewDomainPrimarySetEvent(ctx, orgAgg, orgDomain.Domain))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(domainWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&domainWriteModel.WriteModel), nil
}

func (c *Commands) RemoveOrgDomain(ctx context.Context, orgDomain *domain.OrgDomain) (*domain.ObjectDetails, error) {
	if orgDomain == nil || !orgDomain.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-SJsK3", "Errors.Org.InvalidDomain")
	}
	domainWriteModel, err := c.getOrgDomainWriteModel(ctx, orgDomain.AggregateID, orgDomain.Domain)
	if err != nil {
		return nil, err
	}
	if domainWriteModel.State != domain.OrgDomainStateActive {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-GDfA3", "Errors.Org.DomainNotOnOrg")
	}
	if domainWriteModel.Primary {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-Sjdi3", "Errors.Org.PrimaryDomainNotDeletable")
	}
	orgAgg := OrgAggregateFromWriteModel(&domainWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewDomainRemovedEvent(ctx, orgAgg, orgDomain.Domain))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(domainWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&domainWriteModel.WriteModel), nil
}

func (c *Commands) addOrgDomain(ctx context.Context, orgAgg *eventstore.Aggregate, addedDomain *OrgDomainWriteModel, orgDomain *domain.OrgDomain, claimedUserIDs ...string) ([]eventstore.EventPusher, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedDomain)
	if err != nil {
		return nil, err
	}
	if addedDomain.State == domain.OrgDomainStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMA-Bd2jj", "Errors.Org.Domain.AlreadyExists")
	}

	events := []eventstore.EventPusher{
		org.NewDomainAddedEvent(ctx, orgAgg, orgDomain.Domain),
	}

	if orgDomain.Verified {
		events = append(events, org.NewDomainVerifiedEvent(ctx, orgAgg, orgDomain.Domain))
		for _, userID := range claimedUserIDs {
			userEvents, _, err := c.userDomainClaimed(ctx, userID)
			if err != nil {
				logging.LogWithFields("COMMAND-nn8Jf", "userid", userID).WithError(err).Warn("could not claim user")
				continue
			}
			events = append(events, userEvents...)
		}
	}
	if orgDomain.Primary {
		events = append(events, org.NewDomainPrimarySetEvent(ctx, orgAgg, orgDomain.Domain))
	}
	return events, nil
}

func (c *Commands) getOrgDomainWriteModel(ctx context.Context, orgID, domain string) (*OrgDomainWriteModel, error) {
	domainWriteModel := NewOrgDomainWriteModel(orgID, domain)
	err := c.eventstore.FilterToQueryReducer(ctx, domainWriteModel)
	if err != nil {
		return nil, err
	}
	return domainWriteModel, nil
}
