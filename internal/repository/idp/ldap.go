package idp

import (
	"encoding/json"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
)

type LDAPIDPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	Host                string              `json:"host"`
	Port                string              `json:"port,omitempty"`
	TLS                 bool                `json:"tls"`
	BaseDN              string              `json:"baseDN"`
	UserObjectClass     string              `json:"userObjectClass"`
	UserUniqueAttribute string              `json:"userUniqueAttribute"`
	Admin               string              `json:"admin"`
	Password            *crypto.CryptoValue `json:"password"`

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
	id,
	name,
	host,
	port string,
	tls bool,
	baseDN,
	userObjectClass,
	userUniqueAttribute,
	admin string,
	password *crypto.CryptoValue,
	attributes LDAPAttributes,
	options Options,
) *LDAPIDPAddedEvent {
	return &LDAPIDPAddedEvent{
		BaseEvent:           *base,
		ID:                  id,
		Name:                name,
		Host:                host,
		Port:                port,
		TLS:                 tls,
		BaseDN:              baseDN,
		UserObjectClass:     userObjectClass,
		UserUniqueAttribute: userUniqueAttribute,
		Admin:               admin,
		Password:            password,
		LDAPAttributes:      attributes,
		Options:             options,
	}
}

func (e *LDAPIDPAddedEvent) Data() interface{} {
	return e
}

func (e *LDAPIDPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{idpconfig.NewAddIDPConfigNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner)}
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

	oldName string

	ID                  string              `json:"id"`
	Name                *string             `json:"name,omitempty"`
	Host                *string             `json:"host,omitempty"`
	Port                *string             `json:"port,omitempty"`
	TLS                 *bool               `json:"tls,omitempty"`
	BaseDN              *string             `json:"baseDN,omitempty"`
	UserObjectClass     *string             `json:"userObjectClass,omitempty"`
	UserUniqueAttribute *string             `json:"userUniqueAttribute,omitempty"`
	Admin               *string             `json:"admin,omitempty"`
	Password            *crypto.CryptoValue `json:"password,omitempty"`

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
	oldName string,
	changes []LDAPIDPChanges,
) (*LDAPIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-SDf3f", "Errors.NoChangesFound")
	}
	changedEvent := &LDAPIDPChangedEvent{
		BaseEvent: *base,
		ID:        id,
		oldName:   oldName,
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

func ChangeLDAPHost(host string) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.Host = &host
	}
}

func ChangeLDAPPort(port string) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.Port = &port
	}
}

func ChangeLDAPTLS(tls bool) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.TLS = &tls
	}
}

func ChangeLDAPBaseDN(basDN string) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.BaseDN = &basDN
	}
}

func ChangeLDAPUserObjectClass(userObjectClass string) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.UserObjectClass = &userObjectClass
	}
}

func ChangeLDAPUserUniqueAttribute(userUniqueAttribute string) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.UserUniqueAttribute = &userUniqueAttribute
	}
}

func ChangeLDAPAdmin(admin string) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.Admin = &admin
	}
}

func ChangeLDAPPassword(password *crypto.CryptoValue) func(*LDAPIDPChangedEvent) {
	return func(e *LDAPIDPChangedEvent) {
		e.Password = password
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
	if e.Name == nil || e.oldName == *e.Name { // TODO: nil check should be enough?
		return nil
	}
	return []*eventstore.EventUniqueConstraint{
		idpconfig.NewRemoveIDPConfigNameUniqueConstraint(e.oldName, e.Aggregate().ResourceOwner),
		idpconfig.NewAddIDPConfigNameUniqueConstraint(*e.Name, e.Aggregate().ResourceOwner),
	}
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
