package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/features"
)

type FeaturesWriteModel struct {
	eventstore.WriteModel

	State domain.FeaturesState

	TierName                 string
	TierDescription          string
	TierStatus               domain.TierStatus
	TierStatusDescription    string
	LoginPolicyFactors       bool
	LoginPolicyIDP           bool
	LoginPolicyPasswordless  bool
	LoginPolicyRegistration  bool
	LoginPolicyUsernameLogin bool
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
			if e.TierStatus != nil {
				wm.TierStatus = *e.TierStatus
			}
			if e.TierStatusDescription != nil {
				wm.TierStatusDescription = *e.TierStatusDescription
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
		case *features.FeaturesRemovedEvent:
			wm.State = domain.FeaturesStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
