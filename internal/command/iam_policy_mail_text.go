package command

//
//func (c *Commands) AddDefaultMailText(ctx context.Context, policy *domain.MailText) (*domain.MailText, error) {
//	addedPolicy := NewIAMMailTextWriteModel(policy.MessageTextType, policy.Language)
//	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.MailTextWriteModel.WriteModel)
//	event, err := c.addDefaultMailText(ctx, iamAgg, addedPolicy, policy)
//	if err != nil {
//		return nil, err
//	}
//
//	pushedEvents, err := c.eventstore.PushEvents(ctx, event)
//	if err != nil {
//		return nil, err
//	}
//	err = AppendAndReduce(addedPolicy, pushedEvents...)
//	if err != nil {
//		return nil, err
//	}
//	return writeModelToMailTextPolicy(&addedPolicy.MailTextWriteModel), nil
//}
//
//func (c *Commands) addDefaultMailText(ctx context.Context, iamAgg *eventstore.Aggregate, addedPolicy *IAMMailTextWriteModel, mailText *domain.MailText) (eventstore.EventPusher, error) {
//	if !mailText.IsValid() {
//		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-3n8fs", "Errors.IAM.CustomMailText.Invalid")
//	}
//	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
//	if err != nil {
//		return nil, err
//	}
//	if addedPolicy.State == domain.PolicyStateActive {
//		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-9o0pM", "Errors.IAM.CustomMailText.AlreadyExists")
//	}
//
//	return iam_repo.NewMailTextAddedEvent(
//		ctx,
//		iamAgg,
//		mailText.MessageTextType,
//		mailText.Language,
//		mailText.Title,
//		mailText.PreHeader,
//		mailText.Subject,
//		mailText.Greeting,
//		mailText.Text,
//		mailText.ButtonText), nil
//}
//
//func (c *Commands) ChangeDefaultMailText(ctx context.Context, mailText *domain.MailText) (*domain.MailText, error) {
//	if !mailText.IsValid() {
//		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-kd9fs", "Errors.IAM.CustomMailText.Invalid")
//	}
//	existingPolicy, err := c.defaultMailTextWriteModelByID(ctx, mailText.MessageTextType, mailText.Language)
//	if err != nil {
//		return nil, err
//	}
//
//	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
//		return nil, caos_errs.ThrowNotFound(nil, "IAM-0k2Rf", "Errors.IAM.CustomMailText.NotFound")
//	}
//
//	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.MailTextWriteModel.WriteModel)
//	changedEvent, hasChanged := existingPolicy.NewChangedEvent(
//		ctx,
//		iamAgg,
//		mailText.MessageTextType,
//		mailText.Language,
//		mailText.Title,
//		mailText.PreHeader,
//		mailText.Subject,
//		mailText.Greeting,
//		mailText.Text,
//		mailText.ButtonText)
//	if !hasChanged {
//		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-m9L0s", "Errors.IAM.CustomMailText.NotChanged")
//	}
//
//	pushedEvents, err := c.eventstore.PushEvents(ctx, changedEvent)
//	if err != nil {
//		return nil, err
//	}
//	err = AppendAndReduce(existingPolicy, pushedEvents...)
//	if err != nil {
//		return nil, err
//	}
//	return writeModelToMailTextPolicy(&existingPolicy.MailTextWriteModel), nil
//}
//
//func (c *Commands) defaultMailTextWriteModelByID(ctx context.Context, mailTextType, language string) (policy *IAMMailTextWriteModel, err error) {
//	ctx, span := tracing.NewSpan(ctx)
//	defer func() { span.EndWithError(err) }()
//
//	writeModel := NewIAMMailTextWriteModel(mailTextType, language)
//	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
//	if err != nil {
//		return nil, err
//	}
//	return writeModel, nil
//}
