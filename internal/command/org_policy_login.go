package command

import (
	"context"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type AddLoginPolicy struct {
	AllowUsernamePassword      bool
	AllowRegister              bool
	AllowExternalIDP           bool
	IDPProviders               []*AddLoginPolicyIDP
	ForceMFA                   bool
	SecondFactors              []domain.SecondFactorType
	MultiFactors               []domain.MultiFactorType
	PasswordlessType           domain.PasswordlessType
	HidePasswordReset          bool
	IgnoreUnknownUsernames     bool
	AllowDomainDiscovery       bool
	DefaultRedirectURI         string
	PasswordCheckLifetime      time.Duration
	ExternalLoginCheckLifetime time.Duration
	MFAInitSkipLifetime        time.Duration
	SecondFactorCheckLifetime  time.Duration
	MultiFactorCheckLifetime   time.Duration
	DisableLoginWithEmail      bool
	DisableLoginWithPhone      bool
}

type AddLoginPolicyIDP struct {
	ConfigID string
	Type     domain.IdentityProviderType
}

type ChangeLoginPolicy struct {
	AllowUsernamePassword      bool
	AllowRegister              bool
	AllowExternalIDP           bool
	ForceMFA                   bool
	PasswordlessType           domain.PasswordlessType
	HidePasswordReset          bool
	IgnoreUnknownUsernames     bool
	AllowDomainDiscovery       bool
	DefaultRedirectURI         string
	PasswordCheckLifetime      time.Duration
	ExternalLoginCheckLifetime time.Duration
	MFAInitSkipLifetime        time.Duration
	SecondFactorCheckLifetime  time.Duration
	MultiFactorCheckLifetime   time.Duration
	DisableLoginWithEmail      bool
	DisableLoginWithPhone      bool
}

func (c *Commands) AddLoginPolicy(ctx context.Context, resourceOwner string, policy *AddLoginPolicy) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareAddLoginPolicy(orgAgg, policy))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) orgLoginPolicyWriteModelByID(ctx context.Context, orgID string) (*OrgLoginPolicyWriteModel, error) {
	policyWriteModel := NewOrgLoginPolicyWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, policyWriteModel)
	if err != nil {
		return nil, err
	}
	return policyWriteModel, nil
}

func (c *Commands) getOrgLoginPolicy(ctx context.Context, orgID string) (*domain.LoginPolicy, error) {
	policy, err := c.orgLoginPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if policy.State == domain.PolicyStateActive {
		return writeModelToLoginPolicy(&policy.LoginPolicyWriteModel), nil
	}
	return c.getDefaultLoginPolicy(ctx)
}

func (c *Commands) ChangeLoginPolicy(ctx context.Context, resourceOwner string, policy *ChangeLoginPolicy) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareChangeLoginPolicy(orgAgg, policy))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) RemoveLoginPolicy(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-55Mg9", "Errors.ResourceOwnerMissing")
	}
	existingPolicy := NewOrgLoginPolicyWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-GHB37", "Errors.Org.LoginPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewLoginPolicyRemovedEvent(ctx, orgAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.WriteModel), nil
}

func (c *Commands) AddIDPToLoginPolicy(ctx context.Context, resourceOwner string, idpProvider *domain.IDPProvider) (*domain.IDPProvider, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-M0fs9", "Errors.ResourceOwnerMissing")
	}
	if !idpProvider.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-9nf88", "Errors.Org.LoginPolicy.IDP.")
	}
	existingPolicy, err := c.orgLoginPolicyWriteModelByID(ctx, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-Ffgw2", "Errors.Org.LoginPolicy.NotFound")
	}

	var exists bool
	if idpProvider.Type == domain.IdentityProviderTypeOrg {
		exists, err = ExistsOrgIDP(ctx, c.eventstore.Filter, idpProvider.IDPConfigID, resourceOwner)
	} else {
		exists, err = ExistsInstanceIDP(ctx, c.eventstore.Filter, idpProvider.IDPConfigID)
	}
	if !exists || err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "Org-3N9fs", "Errors.IDPConfig.NotExisting")
	}
	idpModel := NewOrgIdentityProviderWriteModel(resourceOwner, idpProvider.IDPConfigID)
	err = c.eventstore.FilterToQueryReducer(ctx, idpModel)
	if err != nil {
		return nil, err
	}
	if idpModel.State == domain.IdentityProviderStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "Org-2B0ps", "Errors.Org.LoginPolicy.IDP.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&idpModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewIdentityProviderAddedEvent(ctx, orgAgg, idpProvider.IDPConfigID, idpProvider.Type))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(idpModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPProvider(&idpModel.IdentityProviderWriteModel), nil
}

func (c *Commands) RemoveIDPFromLoginPolicy(ctx context.Context, resourceOwner string, idpProvider *domain.IDPProvider, cascadeExternalIDPs ...*domain.UserIDPLink) (*domain.ObjectDetails, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-M0fs9", "Errors.ResourceOwnerMissing")
	}
	if !idpProvider.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-66m9s", "Errors.Org.LoginPolicy.IDP.Invalid")
	}
	existingPolicy, err := c.orgLoginPolicyWriteModelByID(ctx, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-GVDfe", "Errors.Org.LoginPolicy.NotFound")
	}

	idpModel := NewOrgIdentityProviderWriteModel(resourceOwner, idpProvider.IDPConfigID)
	err = c.eventstore.FilterToQueryReducer(ctx, idpModel)
	if err != nil {
		return nil, err
	}
	if idpModel.State == domain.IdentityProviderStateUnspecified || idpModel.State == domain.IdentityProviderStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-39fjs", "Errors.Org.LoginPolicy.IDP.NotExisting")
	}

	orgAgg := OrgAggregateFromWriteModel(&idpModel.IdentityProviderWriteModel.WriteModel)
	events := c.removeIDPFromLoginPolicy(ctx, orgAgg, idpProvider.IDPConfigID, false, cascadeExternalIDPs...)

	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(idpModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&idpModel.WriteModel), nil
}

func (c *Commands) removeIDPFromLoginPolicy(ctx context.Context, orgAgg *eventstore.Aggregate, idpConfigID string, cascade bool, cascadeExternalIDPs ...*domain.UserIDPLink) []eventstore.Command {
	var events []eventstore.Command
	if cascade {
		events = append(events, org.NewIdentityProviderCascadeRemovedEvent(ctx, orgAgg, idpConfigID))
	} else {
		events = append(events, org.NewIdentityProviderRemovedEvent(ctx, orgAgg, idpConfigID))
	}

	for _, idp := range cascadeExternalIDPs {
		event, _, err := c.removeUserIDPLink(ctx, idp, true)
		if err != nil {
			logging.LogWithFields("COMMAND-n8RRf", "userid", idp.AggregateID, "idpconfigid", idp.IDPConfigID).WithError(err).Warn("could not cascade remove external idp")
			continue
		}
		events = append(events, event)
	}
	return events
}

func (c *Commands) AddSecondFactorToLoginPolicy(ctx context.Context, secondFactor domain.SecondFactorType, orgID string) (domain.SecondFactorType, *domain.ObjectDetails, error) {
	if orgID == "" {
		return domain.SecondFactorTypeUnspecified, nil, caos_errs.ThrowInvalidArgument(nil, "Org-M0fs9", "Errors.ResourceOwnerMissing")
	}
	if !secondFactor.Valid() {
		return domain.SecondFactorTypeUnspecified, nil, caos_errs.ThrowInvalidArgument(nil, "Org-5m9fs", "Errors.Org.LoginPolicy.MFA.Unspecified")
	}
	secondFactorModel := NewOrgSecondFactorWriteModel(orgID, secondFactor)
	addedEvent, err := c.addSecondFactorToLoginPolicy(ctx, secondFactorModel, secondFactor)
	if err != nil {
		return domain.SecondFactorTypeUnspecified, nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, addedEvent)
	if err != nil {
		return domain.SecondFactorTypeUnspecified, nil, err
	}

	err = AppendAndReduce(secondFactorModel, pushedEvents...)
	if err != nil {
		return domain.SecondFactorTypeUnspecified, nil, err
	}
	return secondFactorModel.MFAType, writeModelToObjectDetails(&secondFactorModel.WriteModel), nil
}

func (c *Commands) addSecondFactorToLoginPolicy(ctx context.Context, secondFactorModel *OrgSecondFactorWriteModel, secondFactor domain.SecondFactorType) (*org.LoginPolicySecondFactorAddedEvent, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, secondFactorModel)
	if err != nil {
		return nil, err
	}

	if secondFactorModel.State == domain.FactorStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "Org-2B0ps", "Errors.Org.LoginPolicy.MFA.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel)
	return org.NewLoginPolicySecondFactorAddedEvent(ctx, orgAgg, secondFactor), nil
}

func (c *Commands) RemoveSecondFactorFromLoginPolicy(ctx context.Context, secondFactor domain.SecondFactorType, orgID string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-fM0gs", "Errors.ResourceOwnerMissing")
	}
	if !secondFactor.Valid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-55n8s", "Errors.Org.LoginPolicy.MFA.Unspecified")
	}
	secondFactorModel := NewOrgSecondFactorWriteModel(orgID, secondFactor)
	removedEvent, err := c.removeSecondFactorFromLoginPolicy(ctx, secondFactorModel, secondFactor)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, removedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(secondFactorModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&secondFactorModel.WriteModel), nil
}

func (c *Commands) removeSecondFactorFromLoginPolicy(ctx context.Context, secondFactorModel *OrgSecondFactorWriteModel, secondFactor domain.SecondFactorType) (*org.LoginPolicySecondFactorRemovedEvent, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, secondFactorModel)
	if err != nil {
		return nil, err
	}
	if secondFactorModel.State == domain.FactorStateUnspecified || secondFactorModel.State == domain.FactorStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-3M9od", "Errors.Org.LoginPolicy.MFA.NotExisting")
	}
	orgAgg := OrgAggregateFromWriteModel(&secondFactorModel.SecondFactorWriteModel.WriteModel)
	return org.NewLoginPolicySecondFactorRemovedEvent(ctx, orgAgg, secondFactor), nil
}

func (c *Commands) AddMultiFactorToLoginPolicy(ctx context.Context, multiFactor domain.MultiFactorType, orgID string) (domain.MultiFactorType, *domain.ObjectDetails, error) {
	if orgID == "" {
		return domain.MultiFactorTypeUnspecified, nil, caos_errs.ThrowInvalidArgument(nil, "Org-M0fsf", "Errors.ResourceOwnerMissing")
	}
	if !multiFactor.Valid() {
		return domain.MultiFactorTypeUnspecified, nil, caos_errs.ThrowInvalidArgument(nil, "Org-5m9fs", "Errors.Org.LoginPolicy.MFA.Unspecified")
	}
	multiFactorModel := NewOrgMultiFactorWriteModel(orgID, multiFactor)
	addedEvent, err := c.addMultiFactorToLoginPolicy(ctx, multiFactorModel, multiFactor)
	if err != nil {
		return domain.MultiFactorTypeUnspecified, nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, addedEvent)
	if err != nil {
		return domain.MultiFactorTypeUnspecified, nil, err
	}
	err = AppendAndReduce(multiFactorModel, pushedEvents...)
	if err != nil {
		return domain.MultiFactorTypeUnspecified, nil, err
	}
	return multiFactorModel.MFAType, writeModelToObjectDetails(&multiFactorModel.WriteModel), nil
}

func (c *Commands) addMultiFactorToLoginPolicy(ctx context.Context, multiFactorModel *OrgMultiFactorWriteModel, multiFactor domain.MultiFactorType) (*org.LoginPolicyMultiFactorAddedEvent, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, multiFactorModel)
	if err != nil {
		return nil, err
	}
	if multiFactorModel.State == domain.FactorStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "Org-3M9od", "Errors.Org.LoginPolicy.MFA.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&multiFactorModel.WriteModel)
	return org.NewLoginPolicyMultiFactorAddedEvent(ctx, orgAgg, multiFactor), nil
}

func (c *Commands) RemoveMultiFactorFromLoginPolicy(ctx context.Context, multiFactor domain.MultiFactorType, orgID string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-M0fsf", "Errors.ResourceOwnerMissing")
	}
	if !multiFactor.Valid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-5m9fs", "Errors.Org.LoginPolicy.MFA.Unspecified")
	}
	multiFactorModel := NewOrgMultiFactorWriteModel(orgID, multiFactor)
	removedEvent, err := c.removeMultiFactorFromLoginPolicy(ctx, multiFactorModel, multiFactor)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, removedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(multiFactorModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&multiFactorModel.WriteModel), nil
}

func (c *Commands) removeMultiFactorFromLoginPolicy(ctx context.Context, multiFactorModel *OrgMultiFactorWriteModel, multiFactor domain.MultiFactorType) (*org.LoginPolicyMultiFactorRemovedEvent, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, multiFactorModel)
	if err != nil {
		return nil, err
	}
	if multiFactorModel.State == domain.FactorStateUnspecified || multiFactorModel.State == domain.FactorStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-3M9df", "Errors.Org.LoginPolicy.MFA.NotExisting")
	}
	orgAgg := OrgAggregateFromWriteModel(&multiFactorModel.MultiFactorWriteModel.WriteModel)

	return org.NewLoginPolicyMultiFactorRemovedEvent(ctx, orgAgg, multiFactor), nil
}

func (c *Commands) orgLoginPolicyAuthFactorsWriteModel(ctx context.Context, orgID string) (_ *OrgAuthFactorsAllowedWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewOrgAuthFactorsAllowedWriteModel(ctx, orgID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func prepareAddLoginPolicy(a *org.Aggregate, policy *AddLoginPolicy) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if ok := domain.ValidateDefaultRedirectURI(policy.DefaultRedirectURI); !ok {
			return nil, caos_errs.ThrowInvalidArgument(nil, "Org-WSfdq", "Errors.Org.LoginPolicy.RedirectURIInvalid")
		}
		for _, factor := range policy.SecondFactors {
			if !factor.Valid() {
				return nil, caos_errs.ThrowInvalidArgument(nil, "Org-SFeea", "Errors.Org.LoginPolicy.MFA.Unspecified")
			}
		}
		for _, factor := range policy.MultiFactors {
			if !factor.Valid() {
				return nil, caos_errs.ThrowInvalidArgument(nil, "Org-WSfrg", "Errors.Org.LoginPolicy.MFA.Unspecified")
			}
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			if exists, err := exists(ctx, filter, NewOrgLoginPolicyWriteModel(a.ID)); exists || err != nil {
				return nil, caos_errs.ThrowAlreadyExists(nil, "Org-Dgfb2", "Errors.Org.LoginPolicy.AlreadyExists")
			}
			for _, idp := range policy.IDPProviders {
				exists, err := idpExists(ctx, filter, idp)
				if !exists || err != nil {
					return nil, caos_errs.ThrowPreconditionFailed(err, "Org-FEd32", "Errors.IDPConfig.NotExisting")
				}
			}
			cmds := make([]eventstore.Command, 0, len(policy.SecondFactors)+len(policy.MultiFactors)+len(policy.IDPProviders)+1)
			cmds = append(cmds, org.NewLoginPolicyAddedEvent(ctx, &a.Aggregate,
				policy.AllowUsernamePassword,
				policy.AllowRegister,
				policy.AllowExternalIDP,
				policy.ForceMFA,
				policy.HidePasswordReset,
				policy.IgnoreUnknownUsernames,
				policy.AllowDomainDiscovery,
				policy.DisableLoginWithEmail,
				policy.DisableLoginWithPhone,
				policy.PasswordlessType,
				policy.DefaultRedirectURI,
				policy.PasswordCheckLifetime,
				policy.ExternalLoginCheckLifetime,
				policy.MFAInitSkipLifetime,
				policy.SecondFactorCheckLifetime,
				policy.MultiFactorCheckLifetime,
			))
			for _, factor := range policy.SecondFactors {
				cmds = append(cmds, org.NewLoginPolicySecondFactorAddedEvent(ctx, &a.Aggregate, factor))
			}
			for _, factor := range policy.MultiFactors {
				cmds = append(cmds, org.NewLoginPolicyMultiFactorAddedEvent(ctx, &a.Aggregate, factor))
			}
			for _, idp := range policy.IDPProviders {
				cmds = append(cmds, org.NewIdentityProviderAddedEvent(ctx, &a.Aggregate, idp.ConfigID, idp.Type))
			}
			return cmds, nil
		}, nil
	}
}

func prepareChangeLoginPolicy(a *org.Aggregate, policy *ChangeLoginPolicy) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if ok := domain.ValidateDefaultRedirectURI(policy.DefaultRedirectURI); !ok {
			return nil, caos_errs.ThrowInvalidArgument(nil, "Org-Sfd21", "Errors.Org.LoginPolicy.RedirectURIInvalid")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			wm := NewOrgLoginPolicyWriteModel(a.ID)
			if err := queryAndReduce(ctx, filter, wm); err != nil {
				return nil, err
			}
			if !wm.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "Org-M0sif", "Errors.Org.LoginPolicy.NotFound")
			}
			changedEvent, hasChanged := wm.NewChangedEvent(ctx, &a.Aggregate,
				policy.AllowUsernamePassword,
				policy.AllowRegister,
				policy.AllowExternalIDP,
				policy.ForceMFA,
				policy.HidePasswordReset,
				policy.IgnoreUnknownUsernames,
				policy.AllowDomainDiscovery,
				policy.DisableLoginWithEmail,
				policy.DisableLoginWithPhone,
				policy.PasswordlessType,
				policy.DefaultRedirectURI,
				policy.PasswordCheckLifetime,
				policy.ExternalLoginCheckLifetime,
				policy.MFAInitSkipLifetime,
				policy.SecondFactorCheckLifetime,
				policy.MultiFactorCheckLifetime)
			if !hasChanged {
				return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-5M9vdd", "Errors.Org.LoginPolicy.NotChanged")
			}
			return []eventstore.Command{changedEvent}, nil
		}, nil
	}
}

func idpExists(ctx context.Context, filter preparation.FilterToQueryReducer, idp *AddLoginPolicyIDP) (bool, error) {
	if idp.Type == domain.IdentityProviderTypeSystem {
		return exists(ctx, filter, NewInstanceIDPConfigWriteModel(ctx, idp.ConfigID))
	}
	return exists(ctx, filter, NewOrgIDPConfigWriteModel(idp.ConfigID, authz.GetCtxData(ctx).ResourceOwner))
}
