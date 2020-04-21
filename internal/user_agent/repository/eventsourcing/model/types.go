package model

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	UserAgentAggregate models.AggregateType = "useragent"
	Key                models.AggregateType = "key"

	UserAgentAdded   models.EventType = "useragent.added"
	UserAgentChanged models.EventType = "useragent.changed"
	UserAgentRevoked models.EventType = "useragent.revoked"

	UserSessionAdded       models.EventType = "useragent.usersession.added"
	UserSessionTerminated  models.EventType = "useragent.usersession.terminated"
	UserNameCheckSucceeded models.EventType = "useragent.usersession.username.check.succeeded"
	UserNameCheckFailed    models.EventType = "useragent.usersession.username.check.failed"
	PasswordCheckSucceeded models.EventType = "useragent.usersession.password.check.succeeded"
	PasswordCheckFailed    models.EventType = "useragent.usersession.password.check.failed"
	MfaCheckSucceeded      models.EventType = "useragent.usersession.mfa.check.succeeded"
	MfaCheckFailed         models.EventType = "useragent.usersession.mfa.ckeck.failed"
	ReAuthRequested        models.EventType = "useragent.usersession.reauth.requested"

	AuthSessionAdded models.EventType = "useragent.usersession.authsession.added"
	AuthSessionSet   models.EventType = "useragent.usersession.authsession.set"

	TokenAdded models.EventType = "useragent.usersession.authsession.token.added"

	KeyPairAdded      models.EventType = "key.pair.added"
	KeyPairRevoked    models.EventType = "key.pair.revoked"
	PrivateKeyRevoked models.EventType = "key.private.revoked"
)
