package policy

import (
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	loginPolicyPrefix           = "policy.login."
	LoginPolicyAddedEventType   = loginPolicyPrefix + "added"
	LoginPolicyChangedEventType = loginPolicyPrefix + "changed"
	LoginPolicyRemovedEventType = loginPolicyPrefix + "removed"
)

type LoginPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AllowUserNamePassword      bool                    `json:"allowUsernamePassword,omitempty"`
	AllowRegister              bool                    `json:"allowRegister,omitempty"`
	AllowExternalIDP           bool                    `json:"allowExternalIdp,omitempty"`
	ForceMFA                   bool                    `json:"forceMFA,omitempty"`
	ForceMFALocalOnly          bool                    `json:"forceMFALocalOnly,omitempty"`
	HidePasswordReset          bool                    `json:"hidePasswordReset,omitempty"`
	IgnoreUnknownUsernames     bool                    `json:"ignoreUnknownUsernames,omitempty"`
	AllowDomainDiscovery       bool                    `json:"allowDomainDiscovery,omitempty"`
	DisableLoginWithEmail      bool                    `json:"disableLoginWithEmail,omitempty"`
	DisableLoginWithPhone      bool                    `json:"disableLoginWithPhone,omitempty"`
	PasswordlessType           domain.PasswordlessType `json:"passwordlessType,omitempty"`
	DefaultRedirectURI         string                  `json:"defaultRedirectURI,omitempty"`
	PasswordCheckLifetime      time.Duration           `json:"passwordCheckLifetime,omitempty"`
	ExternalLoginCheckLifetime time.Duration           `json:"externalLoginCheckLifetime,omitempty"`
	MFAInitSkipLifetime        time.Duration           `json:"mfaInitSkipLifetime,omitempty"`
	SecondFactorCheckLifetime  time.Duration           `json:"secondFactorCheckLifetime,omitempty"`
	MultiFactorCheckLifetime   time.Duration           `json:"multiFactorCheckLifetime,omitempty"`
}

func (e *LoginPolicyAddedEvent) Data() interface{} {
	return e
}

func (e *LoginPolicyAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLoginPolicyAddedEvent(
	base *eventstore.BaseEvent,
	allowUserNamePassword,
	allowRegister,
	allowExternalIDP,
	forceMFA,
	forceMFALocalOnly,
	hidePasswordReset,
	ignoreUnknownUsernames,
	allowDomainDiscovery,
	disableLoginWithEmail,
	disableLoginWithPhone bool,
	passwordlessType domain.PasswordlessType,
	defaultRedirectURI string,
	passwordCheckLifetime,
	externalLoginCheckLifetime,
	mfaInitSkipLifetime,
	secondFactorCheckLifetime,
	multiFactorCheckLifetime time.Duration,
) *LoginPolicyAddedEvent {
	return &LoginPolicyAddedEvent{
		BaseEvent:                  *base,
		AllowExternalIDP:           allowExternalIDP,
		AllowRegister:              allowRegister,
		AllowUserNamePassword:      allowUserNamePassword,
		ForceMFA:                   forceMFA,
		ForceMFALocalOnly:          forceMFALocalOnly,
		PasswordlessType:           passwordlessType,
		HidePasswordReset:          hidePasswordReset,
		IgnoreUnknownUsernames:     ignoreUnknownUsernames,
		AllowDomainDiscovery:       allowDomainDiscovery,
		DefaultRedirectURI:         defaultRedirectURI,
		PasswordCheckLifetime:      passwordCheckLifetime,
		ExternalLoginCheckLifetime: externalLoginCheckLifetime,
		MFAInitSkipLifetime:        mfaInitSkipLifetime,
		SecondFactorCheckLifetime:  secondFactorCheckLifetime,
		MultiFactorCheckLifetime:   multiFactorCheckLifetime,
		DisableLoginWithEmail:      disableLoginWithEmail,
		DisableLoginWithPhone:      disableLoginWithPhone,
	}
}

func LoginPolicyAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &LoginPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-nWndT", "unable to unmarshal policy")
	}

	return e, nil
}

type LoginPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AllowUserNamePassword      *bool                    `json:"allowUsernamePassword,omitempty"`
	AllowRegister              *bool                    `json:"allowRegister,omitempty"`
	AllowExternalIDP           *bool                    `json:"allowExternalIdp,omitempty"`
	ForceMFA                   *bool                    `json:"forceMFA,omitempty"`
	ForceMFALocalOnly          *bool                    `json:"forceMFALocalOnly,omitempty"`
	HidePasswordReset          *bool                    `json:"hidePasswordReset,omitempty"`
	IgnoreUnknownUsernames     *bool                    `json:"ignoreUnknownUsernames,omitempty"`
	AllowDomainDiscovery       *bool                    `json:"allowDomainDiscovery,omitempty"`
	DisableLoginWithEmail      *bool                    `json:"disableLoginWithEmail,omitempty"`
	DisableLoginWithPhone      *bool                    `json:"disableLoginWithPhone,omitempty"`
	PasswordlessType           *domain.PasswordlessType `json:"passwordlessType,omitempty"`
	DefaultRedirectURI         *string                  `json:"defaultRedirectURI,omitempty"`
	PasswordCheckLifetime      *time.Duration           `json:"passwordCheckLifetime,omitempty"`
	ExternalLoginCheckLifetime *time.Duration           `json:"externalLoginCheckLifetime,omitempty"`
	MFAInitSkipLifetime        *time.Duration           `json:"mfaInitSkipLifetime,omitempty"`
	SecondFactorCheckLifetime  *time.Duration           `json:"secondFactorCheckLifetime,omitempty"`
	MultiFactorCheckLifetime   *time.Duration           `json:"multiFactorCheckLifetime,omitempty"`
}

func (e *LoginPolicyChangedEvent) Data() interface{} {
	return e
}

func (e *LoginPolicyChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLoginPolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []LoginPolicyChanges,
) (*LoginPolicyChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "POLICY-ADg34", "Errors.NoChangesFound")
	}
	changeEvent := &LoginPolicyChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type LoginPolicyChanges func(*LoginPolicyChangedEvent)

func ChangeAllowUserNamePassword(allowUserNamePassword bool) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.AllowUserNamePassword = &allowUserNamePassword
	}
}

func ChangeAllowRegister(allowRegister bool) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.AllowRegister = &allowRegister
	}
}

func ChangeAllowExternalIDP(allowExternalIDP bool) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.AllowExternalIDP = &allowExternalIDP
	}
}

func ChangeForceMFA(forceMFA bool) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.ForceMFA = &forceMFA
	}
}

func ChangeForceMFALocalOnly(forceMFALocalOnly bool) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.ForceMFALocalOnly = &forceMFALocalOnly
	}
}

func ChangePasswordlessType(passwordlessType domain.PasswordlessType) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.PasswordlessType = &passwordlessType
	}
}

func ChangeHidePasswordReset(hidePasswordReset bool) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.HidePasswordReset = &hidePasswordReset
	}
}

func ChangePasswordCheckLifetime(passwordCheckLifetime time.Duration) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.PasswordCheckLifetime = &passwordCheckLifetime
	}
}

func ChangeExternalLoginCheckLifetime(externalLoginCheckLifetime time.Duration) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.ExternalLoginCheckLifetime = &externalLoginCheckLifetime
	}
}

func ChangeMFAInitSkipLifetime(mfaInitSkipLifetime time.Duration) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.MFAInitSkipLifetime = &mfaInitSkipLifetime
	}
}

func ChangeSecondFactorCheckLifetime(secondFactorCheckLifetime time.Duration) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.SecondFactorCheckLifetime = &secondFactorCheckLifetime
	}
}

func ChangeMultiFactorCheckLifetime(multiFactorCheckLifetime time.Duration) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.MultiFactorCheckLifetime = &multiFactorCheckLifetime
	}
}

func ChangeIgnoreUnknownUsernames(ignoreUnknownUsernames bool) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.IgnoreUnknownUsernames = &ignoreUnknownUsernames
	}
}

func ChangeAllowDomainDiscovery(allowDomainDiscovery bool) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.AllowDomainDiscovery = &allowDomainDiscovery
	}
}

func ChangeDefaultRedirectURI(defaultRedirectURI string) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.DefaultRedirectURI = &defaultRedirectURI
	}
}

func ChangeDisableLoginWithEmail(disableLoginWithEmail bool) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.DisableLoginWithEmail = &disableLoginWithEmail
	}
}

func ChangeDisableLoginWithPhone(DisableLoginWithPhone bool) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.DisableLoginWithPhone = &DisableLoginWithPhone
	}
}

func LoginPolicyChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &LoginPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-ehssl", "unable to unmarshal policy")
	}

	return e, nil
}

type LoginPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *LoginPolicyRemovedEvent) Data() interface{} {
	return nil
}

func (e *LoginPolicyRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewLoginPolicyRemovedEvent(base *eventstore.BaseEvent) *LoginPolicyRemovedEvent {
	return &LoginPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func LoginPolicyRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	return &LoginPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
