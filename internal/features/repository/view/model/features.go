package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	features_model "github.com/caos/zitadel/internal/features/model"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	org_repo "github.com/caos/zitadel/internal/repository/org"
)

const (
	FeaturesKeyAggregateID = "aggregate_id"
	FeaturesKeyDefault     = "default_features"
)

type FeaturesView struct {
	AggregateID  string    `json:"-" gorm:"column:aggregate_id;primary_key"`
	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
	Sequence     uint64    `json:"-" gorm:"column:sequence"`
	Default      bool      `json:"-" gorm:"column:default_features"`

	TierName                 string `json:"tierName" gorm:"column:tier_name"`
	TierDescription          string `json:"tierDescription" gorm:"column:tier_description"`
	TierState                int32  `json:"tierState" gorm:"column:tier_state"`
	TierStateDescription     string `json:"tierStateDescription" gorm:"column:tier_state_description"`
	LoginPolicyFactors       bool   `json:"loginPolicyFactors" gorm:"column:login_policy_factors"`
	LoginPolicyIDP           bool   `json:"loginPolicyIDP" gorm:"column:login_policy_idp"`
	LoginPolicyPasswordless  bool   `json:"loginPolicyPasswordless" gorm:"column:login_policy_passwordless"`
	LoginPolicyRegistration  bool   `json:"loginPolicyRegistration" gorm:"column:login_policy_registration"`
	LoginPolicyUsernameLogin bool   `json:"loginPolicyUsernameLogin" gorm:"column:login_policy_username_login"`
}

func FeaturesToModel(features *FeaturesView) *features_model.FeaturesView {
	return &features_model.FeaturesView{
		AggregateID:              features.AggregateID,
		CreationDate:             features.CreationDate,
		ChangeDate:               features.ChangeDate,
		Sequence:                 features.Sequence,
		Default:                  features.Default,
		TierName:                 features.TierName,
		TierDescription:          features.TierDescription,
		TierState:                domain.TierState(features.TierState),
		TierStateDescription:     features.TierStateDescription,
		LoginPolicyFactors:       features.LoginPolicyFactors,
		LoginPolicyIDP:           features.LoginPolicyIDP,
		LoginPolicyPasswordless:  features.LoginPolicyPasswordless,
		LoginPolicyRegistration:  features.LoginPolicyRegistration,
		LoginPolicyUsernameLogin: features.LoginPolicyUsernameLogin,
	}
}

func (f *FeaturesView) AppendEvent(event *models.Event) (err error) {
	f.Sequence = event.Sequence
	f.ChangeDate = event.CreationDate
	switch string(event.Type) {
	case string(iam_repo.FeaturesSetEventType):
		f.SetRootData(event)
		f.CreationDate = event.CreationDate
		f.Default = true
		err = f.SetData(event)
	case string(org_repo.FeaturesSetEventType):
		f.SetRootData(event)
		f.CreationDate = event.CreationDate
		err = f.SetData(event)
		f.Default = false
	}
	return err
}

func (f *FeaturesView) SetRootData(event *models.Event) {
	if f.AggregateID == "" {
		f.AggregateID = event.AggregateID
	}
}

func (f *FeaturesView) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, f); err != nil {
		logging.Log("EVEN-DVsf2").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Bfg31", "Could not unmarshal data")
	}
	return nil
}
