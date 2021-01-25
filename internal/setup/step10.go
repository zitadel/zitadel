package setup

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

type Step10 struct {
	DefaultMailTemplate iam_model.MailTemplate
	DefaultMailTexts    []iam_model.MailText

	setup *Setup
}

func (s *Step10) isNil() bool {
	return s == nil
}

func (step *Step10) step() iam_model.Step {
	return iam_model.Step10
}

func (step *Step10) init(setup *Setup) {
	step.setup = setup
}

func (step *Step10) execute(ctx context.Context) (*iam_model.IAM, error) {
	iam, agg, err := step.mailTemplate(ctx, &step.DefaultMailTemplate)
	if err != nil {
		logging.Log("SETUP-1UYCt").WithField("step", step.step()).WithError(err).Error("unable to finish setup (Mail template)")
		return nil, err
	}
	iam, agg, push, err := step.setup.IamEvents.PrepareSetupDone(ctx, iam, agg, step.step())
	if err != nil {
		logging.Log("SETUP-fMLsb").WithField("step", step.step()).WithError(err).Error("unable to finish setup (prepare setup done)")
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, push, iam.AppendEvents, agg)
	if err != nil {
		logging.Log("SETUP-GuS3f").WithField("step", step.step()).WithError(err).Error("unable to finish setup")
		return nil, err
	}

	iam, agg, err = step.defaultMailTexts(ctx, &step.DefaultMailTexts)
	if err != nil {
		logging.Log("SETUP-p4oWq").WithError(err).Error("unable to set up defaultMailTexts")
		return nil, err
	}
	iam, agg, push, err = step.setup.IamEvents.PrepareSetupDone(ctx, iam, agg, step.step())
	if err != nil {
		logging.Log("SETUP-fMLsb").WithField("step", step.step()).WithError(err).Error("unable to finish setup (prepare setup done)")
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, push, iam.AppendEvents, agg)
	if err != nil {
		logging.Log("SETUP-GuS3f").WithField("step", step.step()).WithError(err).Error("unable to finish setup")
		return nil, err
	}

	return iam_es_model.IAMToModel(iam), nil
}

func (step *Step10) mailTemplate(ctx context.Context, mailTemplate *iam_model.MailTemplate) (*iam_es_model.IAM, *models.Aggregate, error) {
	logging.Log("SETUP-cNrF3").Info("setting up mail template")
	mailTemplate.AggregateID = step.setup.iamID
	iam, aggregate, err := step.setup.IamEvents.PrepareAddMailTemplate(ctx, mailTemplate)
	if err != nil {
		return nil, nil, err
	}
	return iam, aggregate, nil
}

func (step *Step10) defaultMailTexts(ctx context.Context, defaultMailTexts *[]iam_model.MailText) (*iam_es_model.IAM, *models.Aggregate, error) {
	logging.Log("SETUP-dsTh3").Info("setting up defaultMailTexts")
	iam := &iam_es_model.IAM{}
	var aggregate *models.Aggregate
	for index, iamDefaultMailText := range *defaultMailTexts {
		iaml, aggregatel, err := step.defaultMailText(ctx, &iamDefaultMailText)
		if err != nil {
			logging.LogWithFields("SETUP-IlLif", "DefaultMailText", iamDefaultMailText.MailTextType).WithError(err).Error("unable to create defaultMailText")
			return nil, nil, err
		}
		if index == 0 {
			aggregate = aggregatel
		} else {
			aggregate.Events = append(aggregate.Events, aggregatel.Events...)
		}
		iam = iaml
	}
	logging.Log("SETUP-dgjT4").Info("defaultMailTexts set up")
	return iam, aggregate, nil
}

func (step *Step10) defaultMailText(ctx context.Context, mailText *iam_model.MailText) (*iam_es_model.IAM, *models.Aggregate, error) {
	logging.Log("SETUP-cNrF3").Info("setting up mail text")
	mailText.AggregateID = step.setup.iamID
	iam, aggregate, err := step.setup.IamEvents.PrepareAddMailText(ctx, mailText)
	if err != nil {
		return nil, nil, err
	}
	return iam, aggregate, nil
}
