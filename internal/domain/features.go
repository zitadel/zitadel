package domain

import (
	"time"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

const (
	FeatureLoginPolicy              = "login_policy"
	FeatureLoginPolicyFactors       = FeatureLoginPolicy + ".factors"
	FeatureLoginPolicyIDP           = FeatureLoginPolicy + ".idp"
	FeatureLoginPolicyPasswordless  = FeatureLoginPolicy + ".passwordless"
	FeatureLoginPolicyRegistration  = FeatureLoginPolicy + ".registration"
	FeatureLoginPolicyUsernameLogin = FeatureLoginPolicy + ".username_login"
	FeaturePasswordComplexityPolicy = "password_complexity_policy"
)

type Features struct {
	es_models.ObjectRoot

	TierName             string
	TierDescription      string
	TierState            FeaturesState
	TierStateDescription string
	IsDefault            bool

	AuditLogRetention        time.Duration
	LoginPolicyFactors       bool
	LoginPolicyIDP           bool
	LoginPolicyPasswordless  bool
	LoginPolicyRegistration  bool
	LoginPolicyUsernameLogin bool
	PasswordComplexityPolicy bool
}

type FeaturesState int32

const (
	FeaturesStateUnspecified FeaturesState = iota
	FeaturesStateActive
	FeaturesStateActionRequired
	FeaturesStateCanceled
	FeaturesStateRemoved

	featuresStateCount
)

func (f FeaturesState) Valid() bool {
	return f >= 0 && f < featuresStateCount
}
