package idp

import (
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

type LDAPIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID                string              `json:"id"`
	Name              string              `json:"name"`
	Servers           []string            `json:"servers"`
	StartTLS          bool                `json:"startTLS"`
	BaseDN            string              `json:"baseDN"`
	BindDN            string              `json:"bindDN"`
	BindPassword      *crypto.CryptoValue `json:"bindPassword"`
	UserBase          string              `json:"userBase"`
	UserObjectClasses []string            `json:"userObjectClasses"`
	UserFilters       []string            `json:"userFilters"`
	Timeout           time.Duration       `json:"timeout"`

	LDAPAttributes
	Options
}

type LDAPAttributes struct {
	IDAttribute                string `json:"idAttribute,omitempty"`
	FirstNameAttribute         string `json:"firstNameAttribute,omitempty"`
	LastNameAttribute          string `json:"lastNameAttribute,omitempty"`
	DisplayNameAttribute       string `json:"displayNameAttribute,omitempty"`
	NickNameAttribute          string `json:"nickNameAttribute,omitempty"`
	PreferredUsernameAttribute string `json:"preferredUsernameAttribute,omitempty"`
	EmailAttribute             string `json:"emailAttribute,omitempty"`
	EmailVerifiedAttribute     string `json:"emailVerifiedAttribute,omitempty"`
	PhoneAttribute             string `json:"phoneAttribute,omitempty"`
	PhoneVerifiedAttribute     string `json:"phoneVerifiedAttribute,omitempty"`
	PreferredLanguageAttribute string `json:"preferredLanguageAttribute,omitempty"`
	AvatarURLAttribute         string `json:"avatarURLAttribute,omitempty"`
	ProfileAttribute           string `json:"profileAttribute,omitempty"`
}

func (o *LDAPAttributes) Changes(attributes LDAPAttributes) LDAPAttributeChanges {
	attrs := LDAPAttributeChanges{}
	if o.IDAttribute != attributes.IDAttribute {
		attrs.IDAttribute = &attributes.IDAttribute
	}
	if o.FirstNameAttribute != attributes.FirstNameAttribute {
		attrs.FirstNameAttribute = &attributes.FirstNameAttribute
	}
	if o.LastNameAttribute != attributes.LastNameAttribute {
		attrs.LastNameAttribute = &attributes.LastNameAttribute
	}
	if o.DisplayNameAttribute != attributes.DisplayNameAttribute {
		attrs.DisplayNameAttribute = &attributes.DisplayNameAttribute
	}
	if o.NickNameAttribute != attributes.NickNameAttribute {
		attrs.NickNameAttribute = &attributes.NickNameAttribute
	}
	if o.PreferredUsernameAttribute != attributes.PreferredUsernameAttribute {
		attrs.PreferredUsernameAttribute = &attributes.PreferredUsernameAttribute
	}
	if o.EmailAttribute != attributes.EmailAttribute {
		attrs.EmailAttribute = &attributes.EmailAttribute
	}
	if o.EmailVerifiedAttribute != attributes.EmailVerifiedAttribute {
		attrs.EmailVerifiedAttribute = &attributes.EmailVerifiedAttribute
	}
	if o.PhoneAttribute != attributes.PhoneAttribute {
		attrs.PhoneAttribute = &attributes.PhoneAttribute
	}
	if o.PhoneVerifiedAttribute != attributes.PhoneVerifiedAttribute {
		attrs.PhoneVerifiedAttribute = &attributes.PhoneVerifiedAttribute
	}
	if o.PreferredLanguageAttribute != attributes.PreferredLanguageAttribute {
		attrs.PreferredLanguageAttribute = &attributes.PreferredLanguageAttribute
	}
	if o.AvatarURLAttribute != attributes.AvatarURLAttribute {
		attrs.AvatarURLAttribute = &attributes.AvatarURLAttribute
	}
	if o.ProfileAttribute != attributes.ProfileAttribute {
		attrs.ProfileAttribute = &attributes.ProfileAttribute
	}
	return attrs
}

func (o *LDAPAttributes) ReduceChanges(changes LDAPAttributeChanges) {
	if changes.IDAttribute != nil {
		o.IDAttribute = *changes.IDAttribute
	}
	if changes.FirstNameAttribute != nil {
		o.FirstNameAttribute = *changes.FirstNameAttribute
	}
	if changes.LastNameAttribute != nil {
		o.LastNameAttribute = *changes.LastNameAttribute
	}
	if changes.DisplayNameAttribute != nil {
		o.DisplayNameAttribute = *changes.DisplayNameAttribute
	}
	if changes.NickNameAttribute != nil {
		o.NickNameAttribute = *changes.NickNameAttribute
	}
	if changes.PreferredUsernameAttribute != nil {
		o.PreferredUsernameAttribute = *changes.PreferredUsernameAttribute
	}
	if changes.EmailAttribute != nil {
		o.EmailAttribute = *changes.EmailAttribute
	}
	if changes.EmailVerifiedAttribute != nil {
		o.EmailVerifiedAttribute = *changes.EmailVerifiedAttribute
	}
	if changes.PhoneAttribute != nil {
		o.PhoneAttribute = *changes.PhoneAttribute
	}
	if changes.PhoneVerifiedAttribute != nil {
		o.PhoneVerifiedAttribute = *changes.PhoneVerifiedAttribute
	}
	if changes.PreferredLanguageAttribute != nil {
		o.PreferredLanguageAttribute = *changes.PreferredLanguageAttribute
	}
	if changes.AvatarURLAttribute != nil {
		o.AvatarURLAttribute = *changes.AvatarURLAttribute
	}
	if changes.ProfileAttribute != nil {
		o.ProfileAttribute = *changes.ProfileAttribute
	}
}

func NewLDAPIDPAddedEvent(
	base *eventstore.BaseEvent,
	id string,
	name string,
	servers []string,
	startTLS bool,
	baseDN string,
	bindDN string,
	bindPassword *crypto.CryptoValue,
	userBase string,
	userObjectClasses []string,
	userFilters []string,
	timeout time.Duration,
	attributes LDAPAttributes,
	options Options,
) *LDAPIDPAddedEvent {
	return &LDAPIDPAddedEvent{
		BaseEvent:         *base,
		ID:                id,
		Name:              name,
		Servers:           servers,
		StartTLS:          startTLS,
		BaseDN:            baseDN,
		BindDN:            bindDN,
		BindPassword:      bindPassword,
		UserBase:          userBase,
		UserObjectClasses: userObjectClasses,
		UserFilters:       userFilters,
		Timeout:           timeout,
		LDAPAttributes:    attributes,
		Options:           options,
	}
}

func (e *LDAPIDPAddedEvent) Data() interface{} {
	return e
}

func (e *LDAPIDPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func LDAPIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &LDAPIDPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-Dgh42", "unable to unmarshal event")
	}

	return e, nil
}

type LDAPIDPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID                string              `json:"id"`
	Name              *string             `json:"name,omitempty"`
	Servers           []string            `json:"servers,omitempty"`
	StartTLS          *bool               `json:"startTLS,omitempty"`
	BaseDN            *string             `json:"baseDN,omitempty"`
	BindDN            *string             `json:"bindDN,omitempty"`
	BindPassword      *crypto.CryptoValue `json:"bindPassword,omitempty"`
	UserBase          *string             `json:"userBase,omitempty"`
	UserObjectClasses []string            `json:"userObjectClasses,omitempty"`
	UserFilters       []string            `json:"userFilters,omitempty"`
	Timeout           *time.Duration      `json:"timeout,omitempty"`

	LDAPAttributeChanges
	OptionChanges
}

type LDAPAttributeChanges struct {
	IDAttribute                *string `json:"idAttribute,omitempty"`
	FirstNameAttribute         *string `json:"firstNameAttribute,omitempty"`
	LastNameAttribute          *string `json:"lastNameAttribute,omitempty"`
	DisplayNameAttribute       *string `json:"displayNameAttribute,omitempty"`
	NickNameAttribute          *string `json:"nickNameAttribute,omitempty"`
	PreferredUsernameAttribute *string `json:"preferredUsernameAttribute,omitempty"`
	EmailAttribute             *string `json:"emailAttribute,omitempty"`
	EmailVerifiedAttribute     *string `json:"emailVerifiedAttribute,omitempty"`
	PhoneAttribute             *string `json:"phoneAttribute,omitempty"`
	PhoneVerifiedAttribute     *string `json:"phoneVerifiedAttribute,omitempty"`
	PreferredLanguageAttribute *string `json:"preferredLanguageAttribute,omitempty"`
	AvatarURLAttribute         *string `json:"avatarURLAttribute,omitempty"`
	ProfileAttribute           *string `json:"profileAttribute,omitempty"`
}

func (o LDAPAttributeChanges) IsZero() bool {
	return o.IDAttribute == nil &&
		o.FirstNameAttribute == nil &&
		o.LastNameAttribute == nil &&
		o.DisplayNameAttribute == nil &&
		o.NickNameAttribute == nil &&
		o.PreferredUsernameAttribute == nil &&
		o.EmailAttribute == nil &&
		o.EmailVerifiedAttribute == nil &&
		o.PhoneAttribute == nil &&
		o.PhoneVerifiedAttribute == nil &&
		o.PreferredLanguageAttribute == nil &&
		o.AvatarURLAttribute == nil &&
		o.ProfileAttribute == nil
}

func NewLDAPIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []LDAPIDPChanges,
) (*LDAPIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-SDf3f", "Errors.NoChangesFound")
	}
	changedEvent := &LDAPIDPChangedEvent{
		BaseEvent: *base,
		ID:        id,
	}
	for _, change := range changes {
		change(changedEvent)
	}
	return changedEvent, nil
}

type LDAPIDPChanges func(*LDAPIDPChangedEvent)

func ChangeLDAPName(name string) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.Name = &name
	}
}

func ChangeLDAPServers(servers []string) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.Servers = servers
	}
}

func ChangeLDAPStartTLS(startTls bool) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.StartTLS = &startTls
	}
}

func ChangeLDAPBaseDN(baseDN string) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.BaseDN = &baseDN
	}
}

func ChangeLDAPBindDN(bindDN string) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.BindDN = &bindDN
	}
}

func ChangeLDAPBindPassword(password *crypto.CryptoValue) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.BindPassword = password
	}
}

func ChangeLDAPUserBase(userBase string) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.UserBase = &userBase
	}
}

func ChangeLDAPUserObjectClasses(objectClasses []string) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.UserObjectClasses = objectClasses
	}
}

func ChangeLDAPUserFilters(userFilters []string) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.UserFilters = userFilters
	}
}

func ChangeLDAPTimeout(timeout time.Duration) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.Timeout = &timeout
	}
}

func ChangeLDAPAttributes(attributes LDAPAttributeChanges) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.LDAPAttributeChanges = attributes
	}
}

func ChangeLDAPOptions(options OptionChanges) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.OptionChanges = options
	}
}

func (e *LDAPIDPChangedEvent) Data() interface{} {
	return e
}

func (e *LDAPIDPChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func LDAPIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &LDAPIDPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-Sfth3", "unable to unmarshal event")
	}

	return e, nil
}
