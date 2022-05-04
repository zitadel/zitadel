package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

type LoginPolicyWriteModel struct {
	eventstore.WriteModel

	AllowUserNamePassword  bool
	AllowRegister          bool
	AllowExternalIDP       bool
	ForceMFA               bool
	HidePasswordReset      bool
	IgnoreUnknownUsernames bool
	DefaultRedirectURI     string
	PasswordlessType       domain.PasswordlessType
	State                  domain.PolicyState
}

func (wm *LoginPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.LoginPolicyAddedEvent:
			wm.AllowRegister = e.AllowRegister
			wm.AllowUserNamePassword = e.AllowUserNamePassword
			wm.AllowExternalIDP = e.AllowExternalIDP
			wm.ForceMFA = e.ForceMFA
			wm.PasswordlessType = e.PasswordlessType
			wm.HidePasswordReset = e.HidePasswordReset
			wm.IgnoreUnknownUsernames = e.IgnoreUnknownUsernames
			wm.DefaultRedirectURI = e.DefaultRedirectURI
			wm.State = domain.PolicyStateActive
		case *policy.LoginPolicyChangedEvent:
			if e.AllowRegister != nil {
				wm.AllowRegister = *e.AllowRegister
			}
			if e.AllowUserNamePassword != nil {
				wm.AllowUserNamePassword = *e.AllowUserNamePassword
			}
			if e.AllowExternalIDP != nil {
				wm.AllowExternalIDP = *e.AllowExternalIDP
			}
			if e.ForceMFA != nil {
				wm.ForceMFA = *e.ForceMFA
			}
			if e.HidePasswordReset != nil {
				wm.HidePasswordReset = *e.HidePasswordReset
			}
			if e.IgnoreUnknownUsernames != nil {
				wm.IgnoreUnknownUsernames = *e.IgnoreUnknownUsernames
			}
			if e.PasswordlessType != nil {
				wm.PasswordlessType = *e.PasswordlessType
			}
			if e.DefaultRedirectURI != nil {
				wm.DefaultRedirectURI = *e.DefaultRedirectURI
			}
		case *policy.LoginPolicyRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
