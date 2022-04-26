package model

import (
	"encoding/json"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/iam/model"
)

const (
	IAMVersion = "v1"
)

type Step int

const (
	Step1     = Step(model.Step1)
	Step2     = Step(model.Step2)
	StepCount = Step(model.StepCount)
)

type IAM struct {
	es_models.ObjectRoot
	SetUpStarted Step   `json:"-"`
	SetUpDone    Step   `json:"-"`
	GlobalOrgID  string `json:"globalOrgId,omitempty"`
	IAMProjectID string `json:"iamProjectId,omitempty"`
}

func IAMToModel(iam *IAM) *model.IAM {
	converted := &model.IAM{
		ObjectRoot:   iam.ObjectRoot,
		SetUpStarted: domain.Step(iam.SetUpStarted),
		SetUpDone:    domain.Step(iam.SetUpDone),
		GlobalOrgID:  iam.GlobalOrgID,
		IAMProjectID: iam.IAMProjectID,
	}
	return converted
}

func (i *IAM) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := i.AppendEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (i *IAM) AppendEvent(event *es_models.Event) (err error) {
	i.ObjectRoot.AppendEvent(event)
	switch event.Type {
	case IAMSetupStarted:
		if len(event.Data) == 0 {
			i.SetUpStarted = Step(model.Step1)
			return
		}
		step := new(struct{ Step Step })
		err = json.Unmarshal(event.Data, step)
		if err != nil {
			return err
		}
		i.SetUpStarted = step.Step
	case IAMSetupDone:
		if len(event.Data) == 0 {
			i.SetUpDone = Step(model.Step1)
			return
		}
		step := new(struct{ Step Step })
		err = json.Unmarshal(event.Data, step)
		if err != nil {
			return err
		}
		i.SetUpDone = step.Step
	case IAMProjectSet,
		GlobalOrgSet:
		err = i.SetData(event)
	}

	return err
}

func (i *IAM) SetData(event *es_models.Event) error {
	i.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, i); err != nil {
		logging.Log("EVEN-9sie4").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-slwi3", "could not unmarshal event")
	}
	return nil
}
