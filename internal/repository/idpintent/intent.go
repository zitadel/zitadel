package idpintent

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	StartedEventType        = instanceEventTypePrefix + "started"
	OAuthSucceededEventType = instanceEventTypePrefix + "oauth.succeeded"
	LDAPSucceededEventType  = instanceEventTypePrefix + "ldap.succeeded"
	FailedEventType         = instanceEventTypePrefix + "failed"
)

type StartedEvent struct {
	eventstore.BaseEvent `json:"-"`

	SuccessURL *url.URL `json:"successURL"`
	FailureURL *url.URL `json:"failureURL"`
	IDPID      string   `json:"idpId"`
}

func NewStartedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	successURL,
	failureURL *url.URL,
	idpID string,
) *StartedEvent {
	return &StartedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			StartedEventType,
		),
		SuccessURL: successURL,
		FailureURL: failureURL,
		IDPID:      idpID,
	}
}

func (e *StartedEvent) Data() interface{} {
	return e
}

func (e *StartedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func StartedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &StartedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-Sf3f1", "unable to unmarshal event")
	}

	return e, nil
}

type OAuthSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPUser     []byte `json:"idpUser"`
	IDPUserID   string `json:"idpUserId,omitempty"`
	IDPUserName string `json:"idpUserName,omitempty"`
	UserID      string `json:"userId,omitempty"`

	IDPAccessToken *crypto.CryptoValue `json:"idpAccessToken,omitempty"`
	IDPIDToken     string              `json:"idpIdToken,omitempty"`
}

func NewOAuthSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpUser []byte,
	idpUserID,
	idpUserName,
	userID string,
	idpAccessToken *crypto.CryptoValue,
	idpIDToken string,
) *OAuthSucceededEvent {
	return &OAuthSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OAuthSucceededEventType,
		),
		IDPUser:        idpUser,
		IDPUserID:      idpUserID,
		IDPUserName:    idpUserName,
		UserID:         userID,
		IDPAccessToken: idpAccessToken,
		IDPIDToken:     idpIDToken,
	}
}

func (e *OAuthSucceededEvent) Data() interface{} {
	return e
}

func (e *OAuthSucceededEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func OAuthSucceededEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &OAuthSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-HBreq", "unable to unmarshal event")
	}

	return e, nil
}

type LDAPSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPUser     []byte `json:"idpUser"`
	IDPUserID   string `json:"idpUserId,omitempty"`
	IDPUserName string `json:"idpUserName,omitempty"`
	UserID      string `json:"userId,omitempty"`

	EntryAttributes map[string][]string `json:"user,omitempty"`
}

func NewLDAPSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpUser []byte,
	idpUserID,
	idpUserName,
	userID string,
	attributes map[string][]string,
) *LDAPSucceededEvent {
	return &LDAPSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OAuthSucceededEventType,
		),
		IDPUser:         idpUser,
		IDPUserID:       idpUserID,
		IDPUserName:     idpUserName,
		UserID:          userID,
		EntryAttributes: attributes,
	}
}

func (e *LDAPSucceededEvent) Data() interface{} {
	return e
}

func (e *LDAPSucceededEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func LDAPSucceededEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &LDAPSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-HBreq", "unable to unmarshal event")
	}

	return e, nil
}

type FailedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Reason string `json:"reason,omitempty"`
}

func NewFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	reason string,
) *FailedEvent {
	return &FailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			FailedEventType,
		),
		Reason: reason,
	}
}

func (e *FailedEvent) Data() interface{} {
	return e
}

func (e *FailedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func FailedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &FailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-Sfer3", "unable to unmarshal event")
	}

	return e, nil
}
