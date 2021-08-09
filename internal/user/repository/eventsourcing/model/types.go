package model

import "github.com/caos/zitadel/internal/eventstore/v1/models"

//aggregates
const (
	UserAggregate            models.AggregateType = "user"
	UserUserNameAggregate    models.AggregateType = "user.username"
	UserExternalIDPAggregate models.AggregateType = "user.human.externalidp"
)

// the following consts are for user v1 events
const (
	UserAdded                     models.EventType = "user.added"
	UserRegistered                models.EventType = "user.selfregistered"
	InitializedUserCodeAdded      models.EventType = "user.initialization.code.added"
	InitializedUserCodeSent       models.EventType = "user.initialization.code.sent"
	InitializedUserCheckSucceeded models.EventType = "user.initialization.check.succeeded"
	InitializedUserCheckFailed    models.EventType = "user.initialization.check.failed"

	UserUserNameReserved models.EventType = "user.username.reserved"
	UserUserNameReleased models.EventType = "user.username.released"

	UserPasswordChanged        models.EventType = "user.password.changed"
	UserPasswordCodeAdded      models.EventType = "user.password.code.added"
	UserPasswordCodeSent       models.EventType = "user.password.code.sent"
	UserPasswordCheckSucceeded models.EventType = "user.password.check.succeeded"
	UserPasswordCheckFailed    models.EventType = "user.password.check.failed"

	UserEmailChanged            models.EventType = "user.email.changed"
	UserEmailVerified           models.EventType = "user.email.verified"
	UserEmailVerificationFailed models.EventType = "user.email.verification.failed"
	UserEmailCodeAdded          models.EventType = "user.email.code.added"
	UserEmailCodeSent           models.EventType = "user.email.code.sent"

	UserPhoneChanged            models.EventType = "user.phone.changed"
	UserPhoneRemoved            models.EventType = "user.phone.removed"
	UserPhoneVerified           models.EventType = "user.phone.verified"
	UserPhoneVerificationFailed models.EventType = "user.phone.verification.failed"
	UserPhoneCodeAdded          models.EventType = "user.phone.code.added"
	UserPhoneCodeSent           models.EventType = "user.phone.code.sent"

	UserProfileChanged  models.EventType = "user.profile.changed"
	UserAddressChanged  models.EventType = "user.address.changed"
	UserUserNameChanged models.EventType = "user.username.changed"

	MFAOTPAdded          models.EventType = "user.mfa.otp.added"
	MFAOTPVerified       models.EventType = "user.mfa.otp.verified"
	MFAOTPRemoved        models.EventType = "user.mfa.otp.removed"
	MFAOTPCheckSucceeded models.EventType = "user.mfa.otp.check.succeeded"
	MFAOTPCheckFailed    models.EventType = "user.mfa.otp.check.failed"
	MFAInitSkipped       models.EventType = "user.mfa.init.skipped"

	SignedOut models.EventType = "user.signed.out"
)

//the following consts are for user(v2)
const (
	UserNameReserved models.EventType = "user.username.reserved"
	UserNameReleased models.EventType = "user.username.released"

	UserLocked      models.EventType = "user.locked"
	UserUnlocked    models.EventType = "user.unlocked"
	UserDeactivated models.EventType = "user.deactivated"
	UserReactivated models.EventType = "user.reactivated"
	UserRemoved     models.EventType = "user.removed"

	UserTokenAdded models.EventType = "user.token.added"

	DomainClaimed     models.EventType = "user.domain.claimed"
	DomainClaimedSent models.EventType = "user.domain.claimed.sent"

	UserMetadataSet     models.EventType = "user.metadata.set"
	UserMetadataRemoved models.EventType = "user.metadata.removed"
)

// the following consts are for user(v2).human
const (
	HumanAdded                     models.EventType = "user.human.added"
	HumanRegistered                models.EventType = "user.human.selfregistered"
	InitializedHumanCodeAdded      models.EventType = "user.human.initialization.code.added"
	InitializedHumanCodeSent       models.EventType = "user.human.initialization.code.sent"
	InitializedHumanCheckSucceeded models.EventType = "user.human.initialization.check.succeeded"
	InitializedHumanCheckFailed    models.EventType = "user.human.initialization.check.failed"

	HumanPasswordChanged        models.EventType = "user.human.password.changed"
	HumanPasswordCodeAdded      models.EventType = "user.human.password.code.added"
	HumanPasswordCodeSent       models.EventType = "user.human.password.code.sent"
	HumanPasswordCheckSucceeded models.EventType = "user.human.password.check.succeeded"
	HumanPasswordCheckFailed    models.EventType = "user.human.password.check.failed"

	HumanExternalLoginCheckSucceeded models.EventType = "user.human.externallogin.check.succeeded"

	HumanExternalIDPReserved models.EventType = "user.human.externalidp.reserved"
	HumanExternalIDPReleased models.EventType = "user.human.externalidp.released"

	HumanExternalIDPAdded          models.EventType = "user.human.externalidp.added"
	HumanExternalIDPRemoved        models.EventType = "user.human.externalidp.removed"
	HumanExternalIDPCascadeRemoved models.EventType = "user.human.externalidp.cascade.removed"

	HumanAvatarAdded   models.EventType = "user.human.avatar.added"
	HumanAvatarRemoved models.EventType = "user.human.avatar.removed"

	HumanEmailChanged            models.EventType = "user.human.email.changed"
	HumanEmailVerified           models.EventType = "user.human.email.verified"
	HumanEmailVerificationFailed models.EventType = "user.human.email.verification.failed"
	HumanEmailCodeAdded          models.EventType = "user.human.email.code.added"
	HumanEmailCodeSent           models.EventType = "user.human.email.code.sent"

	HumanPhoneChanged            models.EventType = "user.human.phone.changed"
	HumanPhoneRemoved            models.EventType = "user.human.phone.removed"
	HumanPhoneVerified           models.EventType = "user.human.phone.verified"
	HumanPhoneVerificationFailed models.EventType = "user.human.phone.verification.failed"
	HumanPhoneCodeAdded          models.EventType = "user.human.phone.code.added"
	HumanPhoneCodeSent           models.EventType = "user.human.phone.code.sent"

	HumanProfileChanged models.EventType = "user.human.profile.changed"
	HumanAddressChanged models.EventType = "user.human.address.changed"

	HumanMFAOTPAdded          models.EventType = "user.human.mfa.otp.added"
	HumanMFAOTPVerified       models.EventType = "user.human.mfa.otp.verified"
	HumanMFAOTPRemoved        models.EventType = "user.human.mfa.otp.removed"
	HumanMFAOTPCheckSucceeded models.EventType = "user.human.mfa.otp.check.succeeded"
	HumanMFAOTPCheckFailed    models.EventType = "user.human.mfa.otp.check.failed"
	HumanMFAInitSkipped       models.EventType = "user.human.mfa.init.skipped"

	HumanMFAU2FTokenAdded            models.EventType = "user.human.mfa.u2f.token.added"
	HumanMFAU2FTokenVerified         models.EventType = "user.human.mfa.u2f.token.verified"
	HumanMFAU2FTokenSignCountChanged models.EventType = "user.human.mfa.u2f.token.signcount.changed"
	HumanMFAU2FTokenRemoved          models.EventType = "user.human.mfa.u2f.token.removed"
	HumanMFAU2FTokenBeginLogin       models.EventType = "user.human.mfa.u2f.token.begin.login"
	HumanMFAU2FTokenCheckSucceeded   models.EventType = "user.human.mfa.u2f.token.check.succeeded"
	HumanMFAU2FTokenCheckFailed      models.EventType = "user.human.mfa.u2f.token.check.failed"

	HumanPasswordlessTokenAdded           models.EventType = "user.human.passwordless.token.added"
	HumanPasswordlessTokenVerified        models.EventType = "user.human.passwordless.token.verified"
	HumanPasswordlessTokenChangeSignCount models.EventType = "user.human.passwordless.token.signcount.changed"
	HumanPasswordlessTokenRemoved         models.EventType = "user.human.passwordless.token.removed"
	HumanPasswordlessTokenBeginLogin      models.EventType = "user.human.passwordless.token.begin.login"
	HumanPasswordlessTokenCheckSucceeded  models.EventType = "user.human.passwordless.token.check.succeeded"
	HumanPasswordlessTokenCheckFailed     models.EventType = "user.human.passwordless.token.check.failed"

	HumanSignedOut models.EventType = "user.human.signed.out"
)

// the following consts are for user(v2).machines
const (
	MachineAdded   models.EventType = "user.machine.added"
	MachineChanged models.EventType = "user.machine.changed"

	MachineKeyAdded   models.EventType = "user.machine.key.added"
	MachineKeyRemoved models.EventType = "user.machine.key.removed"
)
