package setup

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/business/command"
	"github.com/caos/zitadel/internal/v2/business/domain"
)

func StartSetupV2(esConfig es_int.Config, sd systemdefaults.SystemDefaults) (*Setup, error) {
	setup := &Setup{
		iamID: sd.IamID,
	}
	es, err := es_int.Start(esConfig)
	if err != nil {
		return nil, err
	}

	setup.Commands, err = command.StartCommandSide(&command.Config{
		Eventstore:     es.V2(),
		SystemDefaults: sd,
	})
	if err != nil {
		return nil, err
	}

	return setup, nil
}

func (s *Setup) ExecuteV2(ctx context.Context, setUpConfig IAMSetUp) error {
	logging.Log("SETUP-JAK2q").Info("starting setup")

	iam, err := s.IamEvents.IAMByID(ctx, s.iamID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return err
	}
	if iam != nil && (iam.SetUpDone == domain.StepCount-1 || iam.SetUpStarted != iam.SetUpDone) {
		logging.Log("SETUP-VA2k1").Info("all steps done")
		return nil
	}

	if iam == nil {
		iam = &iam_model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: s.iamID}}
	}

	steps, err := setUpConfig.steps(iam_model.Step(iam.SetUpDone))
	if err != nil || len(steps) == 0 {
		return err
	}

	ctx = setSetUpContextData(ctx, s.iamID)

	for _, step := range steps {
		//step.init(s)
		if step.step() != iam_model.Step(iam.SetUpDone+1) {
			logging.LogWithFields("SETUP-rxRM1", "step", step.step(), "previous", iam.SetUpDone).Warn("wrong step order")
			return caos_errs.ThrowPreconditionFailed(nil, "SETUP-wwAqO", "too few steps for this zitadel verison")
		}
		iam, err = s.Commands.StartSetup(ctx, s.iamID, domain.Step(step.step()))
		if err != nil {
			return err
		}

		err = step.execute(ctx, *s.Commands)
		if err != nil {
			return err
		}

		err = s.validateExecutedStep(ctx)
		if err != nil {
			return err
		}
	}

	logging.Log("SETUP-ds31h").Info("setup done")
	return nil
}

func (s *Setup) validateExecutedStep(ctx context.Context) error {
	iam, err := s.IamEvents.IAMByID(ctx, s.iamID)
	if err != nil {
		return err
	}
	if iam.SetUpStarted != iam.SetUpDone {
		return caos_errs.ThrowInternal(nil, "SETUP-QeukK", "started step is not equal to done")
	}
	return nil
}
