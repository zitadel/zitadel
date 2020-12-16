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
	logging.Log("SETUP-hwG32").Info("starting setup")

	iam, err := s.IamEvents.IAMByID(ctx, s.iamID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return err
	}
	if iam != nil && (iam.SetUpDone == iam_model.StepCount-1 || iam.SetUpStarted != iam.SetUpDone) {
		logging.Log("SETUP-cWEsn").Info("all steps done")
		return nil
	}

	if iam == nil {
		iam = &iam_model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: s.iamID}}
	}

	steps, err := setUpConfig.steps(iam.SetUpDone)
	if err != nil || len(steps) == 0 {
		return err
	}

	ctx = setSetUpContextData(ctx, s.iamID)

	for _, step := range steps {
		step.init(s)
		if step.step() != iam.SetUpDone+1 {
			logging.LogWithFields("SETUP-rxRM1", "step", step.step(), "previous", iam.SetUpDone).Warn("wrong step order")
			return errors.ThrowPreconditionFailed(nil, "SETUP-wwAqO", "too few steps for this zitadel verison")
		}
		iam, err = s.Commands.StartSetup(ctx, s.iamID, step.step())
		if err != nil {
			return err
		}

		iam, err = step.execute(ctx)
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
