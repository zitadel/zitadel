package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type OrgDomainPolicyWriteModel struct {
	PolicyDomainWriteModel
}

func NewOrgDomainPolicyWriteModel(orgID string) *OrgDomainPolicyWriteModel {
	return &OrgDomainPolicyWriteModel{
		PolicyDomainWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
		},
	}
}

func (wm *OrgDomainPolicyWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.DomainPolicyAddedEvent:
			wm.PolicyDomainWriteModel.AppendEvents(&e.DomainPolicyAddedEvent)
		case *org.DomainPolicyChangedEvent:
			wm.PolicyDomainWriteModel.AppendEvents(&e.DomainPolicyChangedEvent)
		case *org.DomainPolicyRemovedEvent:
			wm.PolicyDomainWriteModel.AppendEvents(&e.DomainPolicyRemovedEvent)
		}
	}
}

func (wm *OrgDomainPolicyWriteModel) Reduce() error {
	return wm.PolicyDomainWriteModel.Reduce()
}

func (wm *OrgDomainPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.PolicyDomainWriteModel.AggregateID).
		EventTypes(org.DomainPolicyAddedEventType,
			org.DomainPolicyChangedEventType,
			org.DomainPolicyRemovedEventType).
		Builder()
}

func (wm *OrgDomainPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userLoginMustBeDomain,
	validateOrgDomains,
	smtpSenderAddressMatchesInstanceDomain bool) (changedEvent *org.DomainPolicyChangedEvent, usernameChange bool, err error) {
	changes := make([]policy.DomainPolicyChanges, 0)
	if wm.UserLoginMustBeDomain != userLoginMustBeDomain {
		usernameChange = true
		changes = append(changes, policy.ChangeUserLoginMustBeDomain(userLoginMustBeDomain))
	}
	if wm.ValidateOrgDomains != validateOrgDomains {
		changes = append(changes, policy.ChangeValidateOrgDomains(validateOrgDomains))
	}
	if wm.SMTPSenderAddressMatchesInstanceDomain != smtpSenderAddressMatchesInstanceDomain {
		changes = append(changes, policy.ChangeSMTPSenderAddressMatchesInstanceDomain(smtpSenderAddressMatchesInstanceDomain))
	}
	if len(changes) == 0 {
		return nil, false, zerrors.ThrowPreconditionFailed(nil, "ORG-3M9ds", "Errors.Org.LabelPolicy.NotChanged")
	}
	changedEvent, err = org.NewDomainPolicyChangedEvent(ctx, aggregate, changes)
	return changedEvent, usernameChange, err
}
