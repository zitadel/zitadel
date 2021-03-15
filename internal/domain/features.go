package domain

import es_models "github.com/caos/zitadel/internal/eventstore/v1/models"

type Features struct {
	es_models.ObjectRoot

	TierName             string
	TierDescription      string
	TierState            TierState
	TierStateDescription string
	IsDefault            bool

	LoginPolicyFactors       bool
	LoginPolicyIDP           bool
	LoginPolicyPasswordless  bool
	LoginPolicyRegistration  bool
	LoginPolicyUsernameLogin bool
}

type TierState int32

const (
	TierStateActive TierState = iota
	TierStateActionRequired
	TierStateCanceled
)

type FeaturesState int32

const (
	FeaturesStateUnspecified FeaturesState = iota
	FeaturesStateActive
	FeaturesStateRemoved

	featuresStateCount
)

func (f FeaturesState) Valid() bool {
	return f >= 0 && f < featuresStateCount
}
