package command

import (
	"context"
	"slices"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type InstanceDomainPolicyWriteModel struct {
	PolicyDomainWriteModel
}

func NewInstanceDomainPolicyWriteModel(ctx context.Context) *InstanceDomainPolicyWriteModel {
	return &InstanceDomainPolicyWriteModel{
		PolicyDomainWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   authz.GetInstance(ctx).InstanceID(),
				ResourceOwner: authz.GetInstance(ctx).InstanceID(),
			},
		},
	}
}

func (wm *InstanceDomainPolicyWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.DomainPolicyAddedEvent:
			wm.PolicyDomainWriteModel.AppendEvents(&e.DomainPolicyAddedEvent)
		case *instance.DomainPolicyChangedEvent:
			wm.PolicyDomainWriteModel.AppendEvents(&e.DomainPolicyChangedEvent)
		}
	}
}

func (wm *InstanceDomainPolicyWriteModel) Reduce() error {
	return wm.PolicyDomainWriteModel.Reduce()
}

func (wm *InstanceDomainPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.PolicyDomainWriteModel.AggregateID).
		EventTypes(
			instance.DomainPolicyAddedEventType,
			instance.DomainPolicyChangedEventType).
		Builder()
}

func (wm *InstanceDomainPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userLoginMustBeDomain,
	validateOrgDomain,
	smtpSenderAddresssMatchesInstanceDomain bool) (changedEvent *instance.DomainPolicyChangedEvent, usernameChange bool, err error) {
	changes := make([]policy.DomainPolicyChanges, 0)
	if wm.UserLoginMustBeDomain != userLoginMustBeDomain {
		usernameChange = true
		changes = append(changes, policy.ChangeUserLoginMustBeDomain(userLoginMustBeDomain))
	}
	if wm.ValidateOrgDomains != validateOrgDomain {
		changes = append(changes, policy.ChangeValidateOrgDomains(validateOrgDomain))
	}
	if wm.SMTPSenderAddressMatchesInstanceDomain != smtpSenderAddresssMatchesInstanceDomain {
		changes = append(changes, policy.ChangeSMTPSenderAddressMatchesInstanceDomain(smtpSenderAddresssMatchesInstanceDomain))
	}
	if len(changes) == 0 {
		return nil, false, zerrors.ThrowPreconditionFailed(nil, "INSTANCE-pl9fN", "Errors.IAM.DomainPolicy.NotChanged")
	}
	changedEvent, err = instance.NewDomainPolicyChangedEvent(ctx, aggregate, changes)
	return changedEvent, usernameChange, err
}

type DomainPolicyOrgsWriteModel struct {
	eventstore.WriteModel

	OrgIDs []string
}

func NewDomainPolicyOrgsWriteModel() *DomainPolicyOrgsWriteModel {
	return &DomainPolicyOrgsWriteModel{
		WriteModel: eventstore.WriteModel{},
	}
}

func (wm *DomainPolicyOrgsWriteModel) AppendEvents(events ...eventstore.Event) {
	wm.WriteModel.AppendEvents(events...)
}

func (wm *DomainPolicyOrgsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		// organization irrelevant if removed or with custom policy
		case *org.OrgRemovedEvent, *org.DomainPolicyAddedEvent:
			wm.OrgIDs = slices.DeleteFunc(wm.OrgIDs, func(orgID string) bool { return orgID == e.Aggregate().ID })
		// organization relevant if added without custom policy or custom policy is removed
		case *org.OrgAddedEvent, *org.DomainPolicyRemovedEvent:
			wm.OrgIDs = append(wm.OrgIDs, e.Aggregate().ID)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *DomainPolicyOrgsWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(org.AggregateType).
		EventTypes(
			org.OrgAddedEventType,
			org.OrgRemovedEventType,
			org.DomainPolicyAddedEventType,
			org.DomainPolicyRemovedEventType).
		Builder()
}
