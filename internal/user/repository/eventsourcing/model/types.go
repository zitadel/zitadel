package model

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	UserAggregate         models.AggregateType = "user"
	UserUserNameAggregate models.AggregateType = "user.username"
	UserEmailAggregate    models.AggregateType = "user.email"

	UserAdded                models.EventType = "user.added"
	UserRegistered           models.EventType = "user.selfregistered"
	InitializedUserCodeAdded models.EventType = "user.initialization.code.added"
	InitializedUserCodeSent  models.EventType = "user.initialization.code.sent"

	UserUserNameReserved models.EventType = "user.username.reserved"
	UserUserNameReleased models.EventType = "user.username.released"
	UserEmailReserved    models.EventType = "user.email.reserved"
	UserEmailReleased    models.EventType = "user.email.released"

	UserLocked      models.EventType = "user.locked"
	UserUnlocked    models.EventType = "user.unlocked"
	UserDeactivated models.EventType = "user.deactivated"
	UserReactivated models.EventType = "user.reactivated"
	UserDeleted     models.EventType = "user.deleted"

	UserPasswordChanged        models.EventType = "user.password.changed"
	UserPasswordCodeAdded      models.EventType = "user.password.code.added"
	UserPasswordCodeSent       models.EventType = "user.password.code.sent"
	UserPasswordCheckSucceeded models.EventType = "user.password.check.succeeded"
	UserPasswordCheckFailed    models.EventType = "user.password.check.failed"

	UserEmailChanged   models.EventType = "user.email.changed"
	UserEmailVerified  models.EventType = "user.email.verified"
	UserEmailCodeAdded models.EventType = "user.email.code.added"
	UserEmailCodeSent  models.EventType = "user.email.code.sent"

	UserPhoneChanged   models.EventType = "user.phone.changed"
	UserPhoneVerified  models.EventType = "user.phone.verified"
	UserPhoneCodeAdded models.EventType = "user.phone.code.added"
	UserPhoneCodeSent  models.EventType = "user.phone.code.sent"

	UserProfileChanged models.EventType = "user.profile.changed"
	UserAddressChanged models.EventType = "user.address.changed"

	MfaOtpAdded          models.EventType = "user.mfa.otp.added"
	MfaOtpVerified       models.EventType = "user.mfa.otp.verified"
	MfaOtpRemoved        models.EventType = "user.mfa.otp.removed"
	MfaOtpCheckSucceeded models.EventType = "user.mfa.otp.check.succeeded"
	MfaOtpCheckFailed    models.EventType = "user.mfa.otp.check.failed"
	MfaInitSkipped       models.EventType = "user.mfa.init.skipped"

	SignedOut models.EventType = "user.signed.out"
)
