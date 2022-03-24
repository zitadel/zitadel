package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"golang.org/x/text/language"
)

//TODO: private as soon as setup uses query
func (c *Commands) GetIAM(ctx context.Context) (*domain.IAM, error) {
	iamWriteModel := NewIAMWriteModel()
	err := c.eventstore.FilterToQueryReducer(ctx, iamWriteModel)
	if err != nil {
		return nil, err
	}
	return writeModelToIAM(iamWriteModel), nil
}

func (c *Commands) AddInstance(ctx context.Context, name string) (*domain.IAM, error) {
	_, addedInstance, events, err := c.addInstance(ctx, &domain.IAM{Name: name})
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedInstance, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToIAM(addedInstance), nil
}

func (c *Commands) ChangeInstance(ctx context.Context, instanceID, name string) (*domain.ObjectDetails, error) {
	if instanceID == "" || name == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "EVENT-a93lf", "Errors.Instance.Invalid")
	}
	orgWriteModel, err := c.getIAMWriteModel(ctx, instanceID)
	if err != nil {
		return nil, err
	}
	if orgWriteModel.State == domain.OrgStateUnspecified || orgWriteModel.State == domain.OrgStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-1MRds", "Errors.Org.NotFound")
	}
	if orgWriteModel.Name == name {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-4VSdf", "Errors.Org.NotChanged")
	}
	orgAgg := OrgAggregateFromWriteModel(&orgWriteModel.WriteModel)
	events := make([]eventstore.Command, 0)
	events = append(events, org.NewOrgChangedEvent(ctx, orgAgg, orgWriteModel.Name, name))
	changeDomainEvents, err := c.changeDefaultDomain(ctx, orgID, name)
	if err != nil {
		return nil, err
	}
	if len(changeDomainEvents) > 0 {
		events = append(events, changeDomainEvents...)
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(orgWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&orgWriteModel.WriteModel), nil
}

func (c *Commands) addInstance(ctx context.Context, instance *domain.IAM) (_ *eventstore.Aggregate, _ *IAMWriteModel, _ []eventstore.Command, err error) {
	if instance == nil || !instance.IsValid() {
		return nil, nil, nil, caos_errs.ThrowInvalidArgument(nil, "COMM-M9wjf", "Errors.Instance.Invalid")
	}

	instance.AggregateID, err = c.idGenerator.Next()
	if err != nil {
		return nil, nil, nil, caos_errs.ThrowInternal(err, "COMMA-f92lj", "Errors.Internal")
	}
	instance.AddGeneratedDomain(c.iamDomain)
	addedInstance := NewIAMWriteModel(instance.AggregateID)

	instanceAgg := IAMAggregateFromWriteModel(&addedInstance.WriteModel)
	events := []eventstore.Command{
		iam.NewInstanceAddedEvent(ctx, instanceAgg, instance.Name),
	}
	instanceDomainEvent, err := c.addInstanceDomain(ctx, instanceAgg, NewInstanceDomainWriteModel(instanceAgg.ID, instance.GeneratedDomain.Domain), instance.GeneratedDomain)
	if err != nil {
		return nil, nil, nil, err
	}
	events = append(events, instanceDomainEvent)
	return instanceAgg, addedInstance, events, nil
}

func (c *Commands) setGlobalOrg(ctx context.Context, iamAgg *eventstore.Aggregate, iamWriteModel *IAMWriteModel, orgID string) (eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, iamWriteModel)
	if err != nil {
		return nil, err
	}
	if iamWriteModel.GlobalOrgID != "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-HGG24", "Errors.IAM.GlobalOrgAlreadySet")
	}
	return iam.NewGlobalOrgSetEventEvent(ctx, iamAgg, orgID), nil
}

func (c *Commands) setIAMProject(ctx context.Context, iamAgg *eventstore.Aggregate, iamWriteModel *IAMWriteModel, projectID string) (eventstore.Command, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, iamWriteModel)
	if err != nil {
		return nil, err
	}
	if iamWriteModel.ProjectID != "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-EGbw2", "Errors.IAM.IAMProjectAlreadySet")
	}
	return iam.NewIAMProjectSetEvent(ctx, iamAgg, projectID), nil
}

func (c *Commands) SetDefaultLanguage(ctx context.Context, language language.Tag) (*domain.ObjectDetails, error) {
	iamWriteModel, err := c.getIAMWriteModel(ctx)
	if err != nil {
		return nil, err
	}
	iamAgg := IAMAggregateFromWriteModel(&iamWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, iam.NewDefaultLanguageSetEvent(ctx, iamAgg, language))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(iamWriteModel, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&iamWriteModel.WriteModel), nil
}

func (c *Commands) getIAMWriteModel(ctx context.Context, instanceID string) (_ *IAMWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMWriteModel(instanceID)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}
