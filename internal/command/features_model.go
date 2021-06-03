package command

import (
	"time"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/features"
)

type FeaturesWriteModel struct {
	eventstore.WriteModel

	TierName                 string
	TierDescription          string
	State                    domain.FeaturesState
	StateDescription         string
	AuditLogRetention        time.Duration
	LoginPolicyFactors       bool
	LoginPolicyIDP           bool
	LoginPolicyPasswordless  bool
	LoginPolicyRegistration  bool
	LoginPolicyUsernameLogin bool
	LoginPolicyPasswordReset bool
	PasswordComplexityPolicy bool
	LabelPolicy              bool
	CustomDomain             bool
	CustomText               bool
}

func (wm *FeaturesWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *features.FeaturesSetEvent:
			if e.TierName != nil {
				wm.TierName = *e.TierName
			}
			if e.TierDescription != nil {
				wm.TierDescription = *e.TierDescription
			}
			wm.State = domain.FeaturesStateActive
			if e.State != nil {
				wm.State = *e.State
			}
			if e.StateDescription != nil {
				wm.StateDescription = *e.StateDescription
			}
			if e.AuditLogRetention != nil {
				wm.AuditLogRetention = *e.AuditLogRetention
			}
			if e.LoginPolicyFactors != nil {
				wm.LoginPolicyFactors = *e.LoginPolicyFactors
			}
			if e.LoginPolicyIDP != nil {
				wm.LoginPolicyIDP = *e.LoginPolicyIDP
			}
			if e.LoginPolicyPasswordless != nil {
				wm.LoginPolicyPasswordless = *e.LoginPolicyPasswordless
			}
			if e.LoginPolicyRegistration != nil {
				wm.LoginPolicyRegistration = *e.LoginPolicyRegistration
			}
			if e.LoginPolicyUsernameLogin != nil {
				wm.LoginPolicyUsernameLogin = *e.LoginPolicyUsernameLogin
			}
			if e.LoginPolicyPasswordReset != nil {
				wm.LoginPolicyPasswordReset = *e.LoginPolicyPasswordReset
			}
			if e.PasswordComplexityPolicy != nil {
				wm.PasswordComplexityPolicy = *e.PasswordComplexityPolicy
			}
			if e.LabelPolicy != nil {
				wm.LabelPolicy = *e.LabelPolicy
			}
			if e.CustomDomain != nil {
				wm.CustomDomain = *e.CustomDomain
			}
			if e.CustomText != nil {
				wm.CustomText = *e.CustomText
			}
		case *features.FeaturesRemovedEvent:
			wm.State = domain.FeaturesStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
