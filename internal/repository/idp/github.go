package idp

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

type GitHubIDPAddedEvent struct {
	OAuthIDPAddedEvent
}

func NewGitHubIDPAddedEvent(
	base *eventstore.BaseEvent,
	id,
	clientID string,
	clientSecret *crypto.CryptoValue,
	options Options,
) *GitHubIDPAddedEvent {
	return &GitHubIDPAddedEvent{
		OAuthIDPAddedEvent: OAuthIDPAddedEvent{
			BaseEvent:    *base,
			ID:           id,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Options:      options,
		},
	}
}

func (e *GitHubIDPAddedEvent) Data() interface{} {
	return e
}

func (e *GitHubIDPAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
	//return []*eventstore.EventUniqueConstraint{idpconfig.NewAddIDPConfigNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner)}
}

func GitHubIDPAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := OAuthIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubIDPAddedEvent{OAuthIDPAddedEvent: *e.(*OAuthIDPAddedEvent)}, nil
}

type GitHubIDPChangedEvent struct {
	OAuthIDPChangedEvent
}

func NewGitHubIDPChangedEvent(
	base *eventstore.BaseEvent,
	id string,
	changes []OAuthIDPChanges,
) (*GitHubIDPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-BH3dl", "Errors.NoChangesFound")
	}
	changedEvent := &GitHubIDPChangedEvent{
		OAuthIDPChangedEvent: OAuthIDPChangedEvent{
			BaseEvent: *base,
			ID:        id,
		},
	}
	for _, change := range changes {
		change(&changedEvent.OAuthIDPChangedEvent)
	}
	return changedEvent, nil
}

//
//type OAuthIDPChanges func(*OAuthIDPChangedEvent)
//
//func ChangeOAuthName(name string) func(*OAuthIDPChangedEvent) {
//	return func(e *OAuthIDPChangedEvent) {
//		e.Name = &name
//	}
//}
//func ChangeOAuthClientID(clientID string) func(*OAuthIDPChangedEvent) {
//	return func(e *OAuthIDPChangedEvent) {
//		e.ClientID = &clientID
//	}
//}
//
//func ChangeOAuthClientSecret(clientSecret *crypto.CryptoValue) func(*OAuthIDPChangedEvent) {
//	return func(e *OAuthIDPChangedEvent) {
//		e.ClientSecret = clientSecret
//	}
//}
//
//func ChangeOAuthOptions(options OptionChanges) func(*OAuthIDPChangedEvent) {
//	return func(e *OAuthIDPChangedEvent) {
//		e.OptionChanges = options
//	}
//}
//
//func ChangeOAuthAuthorizationEndpoint(authorizationEndpoint string) func(*OAuthIDPChangedEvent) {
//	return func(e *OAuthIDPChangedEvent) {
//		e.AuthorizationEndpoint = &authorizationEndpoint
//	}
//}
//
//func ChangeOAuthTokenEndpoint(tokenEndpoint string) func(*OAuthIDPChangedEvent) {
//	return func(e *OAuthIDPChangedEvent) {
//		e.TokenEndpoint = &tokenEndpoint
//	}
//}
//
//func ChangeOAuthUserEndpoint(userEndpoint string) func(*OAuthIDPChangedEvent) {
//	return func(e *OAuthIDPChangedEvent) {
//		e.UserEndpoint = &userEndpoint
//	}
//}
//
//func ChangeOAuthScopes(scopes []string) func(*OAuthIDPChangedEvent) {
//	return func(e *OAuthIDPChangedEvent) {
//		e.Scopes = scopes
//	}
//}

func (e *GitHubIDPChangedEvent) Data() interface{} {
	return e
}

func (e *GitHubIDPChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
	//if e.Name == nil || e.oldName == *e.Name { // TODO: nil check should be enough
	//	return nil
	//}
	//return []*eventstore.EventUniqueConstraint{
	//	idpconfig.NewRemoveIDPConfigNameUniqueConstraint(e.oldName, e.Aggregate().ResourceOwner),
	//	idpconfig.NewAddIDPConfigNameUniqueConstraint(*e.Name, e.Aggregate().ResourceOwner),
	//}
}

func GitHubIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := OAuthIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubIDPChangedEvent{OAuthIDPChangedEvent: *e.(*OAuthIDPChangedEvent)}, nil
}

//
//func OAuthIDPChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
//	e := &OAuthIDPChangedEvent{
//		BaseEvent: *eventstore.BaseEventFromRepo(event),
//	}
//
//	err := json.Unmarshal(event.Data, e)
//	if err != nil {
//		return nil, errors.ThrowInternal(err, "IDP-D3gjzh", "unable to unmarshal event")
//	}
//
//	return e, nil
//}
