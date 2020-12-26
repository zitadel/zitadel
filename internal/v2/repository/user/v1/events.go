package v1

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/user/human"
	"github.com/caos/zitadel/internal/v2/repository/user/human/address"
	"github.com/caos/zitadel/internal/v2/repository/user/human/email"
	"github.com/caos/zitadel/internal/v2/repository/user/human/mfa"
	"github.com/caos/zitadel/internal/v2/repository/user/human/mfa/otp"
	"github.com/caos/zitadel/internal/v2/repository/user/human/password"
	"github.com/caos/zitadel/internal/v2/repository/user/human/phone"
	"github.com/caos/zitadel/internal/v2/repository/user/human/profile"
	"golang.org/x/text/language"
)

const (
	userEventTypePrefix                 = eventstore.EventType("user.")
	UserV1AddedType                     = userEventTypePrefix + "added"
	UserV1RegisteredType                = userEventTypePrefix + "selfregistered"
	UserV1InitialCodeAddedType          = userEventTypePrefix + "initialization.code.added"
	UserV1InitialCodeSentType           = userEventTypePrefix + "initialization.code.sent"
	UserV1InitializedCheckSucceededType = userEventTypePrefix + "initialization.check.succeeded"
	UserV1InitializedCheckFailedType    = userEventTypePrefix + "initialization.check.failed"
	UserV1SignedOutType                 = userEventTypePrefix + "signed.out"

	userV1PasswordEventTypePrefix    = userEventTypePrefix + "password."
	UserV1PasswordChangedType        = userV1PasswordEventTypePrefix + "changed"
	UserV1PasswordCodeAddedType      = userV1PasswordEventTypePrefix + "code.added"
	UserV1PasswordCodeSentType       = userV1PasswordEventTypePrefix + "code.sent"
	UserV1PasswordCheckSucceededType = userV1PasswordEventTypePrefix + "check.succeeded"
	UserV1PasswordCheckFailedType    = userV1PasswordEventTypePrefix + "check.failed"

	userV1EmailEventTypePrefix        = userEventTypePrefix + "email."
	UserV1EmailChangedType            = userV1EmailEventTypePrefix + "changed"
	UserV1EmailVerifiedType           = userV1EmailEventTypePrefix + "verified"
	UserV1EmailVerificationFailedType = userV1EmailEventTypePrefix + "verification.failed"
	UserV1EmailCodeAddedType          = userV1EmailEventTypePrefix + "code.added"
	UserV1EmailCodeSentType           = userV1EmailEventTypePrefix + "code.sent"

	userV1PhoneEventTypePrefix        = userEventTypePrefix + "phone."
	UserV1PhoneChangedType            = userV1PhoneEventTypePrefix + "changed"
	UserV1PhoneRemovedType            = userV1PhoneEventTypePrefix + "removed"
	UserV1PhoneVerifiedType           = userV1PhoneEventTypePrefix + "verified"
	UserV1PhoneVerificationFailedType = userV1PhoneEventTypePrefix + "verification.failed"
	UserV1PhoneCodeAddedType          = userV1PhoneEventTypePrefix + "code.added"
	UserV1PhoneCodeSentType           = userV1PhoneEventTypePrefix + "code.sent"

	userV1ProfileEventTypePrefix = userEventTypePrefix + "profile."
	UserV1ProfileChangedType     = userV1ProfileEventTypePrefix + "changed"

	userV1AddressEventTypePrefix = userEventTypePrefix + "address."
	UserV1AddressChangedType     = userV1AddressEventTypePrefix + "changed"

	userV1MFAEventTypePrefix = userEventTypePrefix + "mfa."
	UserV1MFAInitSkippedType = userV1MFAOTPEventTypePrefix + "init.skipped"

	userV1MFAOTPEventTypePrefix    = userV1MFAEventTypePrefix + "otp."
	UserV1MFAOTPAddedType          = userV1MFAOTPEventTypePrefix + "added"
	UserV1MFAOTPRemovedType        = userV1MFAOTPEventTypePrefix + "removed"
	UserV1MFAOTPVerifiedType       = userV1MFAOTPEventTypePrefix + "verified"
	UserV1MFAOTPCheckSucceededType = userV1MFAOTPEventTypePrefix + "check.succeeded"
	UserV1MFAOTPCheckFailedType    = userV1MFAOTPEventTypePrefix + "check.failed"
)

func NewUserV1AddedEvent(
	ctx context.Context,
	userName,
	firstName,
	lastName,
	nickName,
	displayName string,
	preferredLanguage language.Tag,
	gender human.Gender,
	emailAddress,
	phoneNumber,
	country,
	locality,
	postalCode,
	region,
	streetAddress string,
) *human.AddedEvent {
	return &human.AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1AddedType,
		),
		UserName:          userName,
		FirstName:         firstName,
		LastName:          lastName,
		NickName:          nickName,
		DisplayName:       displayName,
		PreferredLanguage: preferredLanguage,
		Gender:            gender,
		EmailAddress:      emailAddress,
		PhoneNumber:       phoneNumber,
		Country:           country,
		Locality:          locality,
		PostalCode:        postalCode,
		Region:            region,
		StreetAddress:     streetAddress,
	}
}

func NewUserV1RegisteredEvent(
	ctx context.Context,
	userName,
	firstName,
	lastName,
	nickName,
	displayName string,
	preferredLanguage language.Tag,
	gender int32,
	emailAddress,
	phoneNumber,
	country,
	locality,
	postalCode,
	region,
	streetAddress string,
) *human.RegisteredEvent {
	return &human.RegisteredEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1RegisteredType,
		),
		UserName:          userName,
		FirstName:         firstName,
		LastName:          lastName,
		NickName:          nickName,
		DisplayName:       displayName,
		PreferredLanguage: preferredLanguage,
		Gender:            gender,
		EmailAddress:      emailAddress,
		PhoneNumber:       phoneNumber,
		Country:           country,
		Locality:          locality,
		PostalCode:        postalCode,
		Region:            region,
		StreetAddress:     streetAddress,
	}
}

func NewUserV1InitialCodeAddedEvent(
	ctx context.Context,
	code *crypto.CryptoValue,
	expiry time.Duration,
) *human.InitialCodeAddedEvent {
	return &human.InitialCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1InitialCodeAddedType,
		),
		Code:   code,
		Expiry: expiry,
	}
}

func NewUserV1InitialCodeSentEvent(ctx context.Context) *human.InitialCodeSentEvent {
	return &human.InitialCodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1InitialCodeSentType,
		),
	}
}

func NewUserV1InitializedCheckSucceededEvent(ctx context.Context) *human.InitializedCheckSucceededEvent {
	return &human.InitializedCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1InitializedCheckSucceededType,
		),
	}
}

func NewUserV1InitializedCheckFailedEvent(ctx context.Context) *human.InitializedCheckFailedEvent {
	return &human.InitializedCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1InitializedCheckFailedType,
		),
	}
}

func NewUserV1SignedOutEvent(ctx context.Context) *human.SignedOutEvent {
	return &human.SignedOutEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1SignedOutType,
		),
	}
}

func NewUserV1PasswordChangedEvent(
	ctx context.Context,
	secret *crypto.CryptoValue,
	changeRequired bool,
) *password.ChangedEvent {
	return &password.ChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PasswordChangedType,
		),
		Secret:         secret,
		ChangeRequired: changeRequired,
	}
}

func NewUserV1PasswordCodeAddedEvent(
	ctx context.Context,
	code *crypto.CryptoValue,
	expiry time.Duration,
	notificationType human.NotificationType,
) *password.CodeAddedEvent {
	return &password.CodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PasswordCodeAddedType,
		),
		Code:             code,
		Expiry:           expiry,
		NotificationType: notificationType,
	}
}

func NewUserV1PasswordCodeSentEvent(ctx context.Context) *password.CodeSentEvent {
	return &password.CodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PasswordCodeSentType,
		),
	}
}

func NewUserV1PasswordCheckSucceededEvent(ctx context.Context) *password.CheckSucceededEvent {
	return &password.CheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PasswordCheckSucceededType,
		),
	}
}

func NewUserV1PasswordCheckFailedEvent(ctx context.Context) *password.CheckFailedEvent {
	return &password.CheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PasswordCheckFailedType,
		),
	}
}

func NewUserV1EmailChangedEvent(ctx context.Context, emailAddress string) *email.ChangedEvent {
	return &email.ChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1EmailChangedType,
		),
		EmailAddress: emailAddress,
	}
}

func NewUserV1EmailVerifiedEvent(ctx context.Context) *email.VerifiedEvent {
	return &email.VerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1EmailVerifiedType,
		),
	}
}

func NewUserV1EmailVerificationFailedEvent(ctx context.Context) *email.VerificationFailedEvent {
	return &email.VerificationFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1EmailVerificationFailedType,
		),
	}
}

func NewUserV1EmailCodeAddedEvent(
	ctx context.Context,
	code *crypto.CryptoValue,
	expiry time.Duration,
) *email.CodeAddedEvent {
	return &email.CodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1EmailCodeAddedType,
		),
		Code:   code,
		Expiry: expiry,
	}
}

func NewUserV1EmailCodeSentEvent(ctx context.Context) *email.CodeSentEvent {
	return &email.CodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1EmailCodeSentType,
		),
	}
}

func NewUserV1PhoneChangedEvent(ctx context.Context, phoneNbr string) *phone.ChangedEvent {
	return &phone.ChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PhoneChangedType,
		),
		PhoneNumber: phoneNbr,
	}
}

func NewUserV1PhoneRemovedEvent(ctx context.Context) *phone.RemovedEvent {
	return &phone.RemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PhoneRemovedType,
		),
	}
}

func NewUserV1PhoneVerifiedEvent(ctx context.Context) *phone.VerifiedEvent {
	return &phone.VerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PhoneVerifiedType,
		),
	}
}

func NewUserV1PhoneVerificationFailedEvent(ctx context.Context) *phone.VerificationFailedEvent {
	return &phone.VerificationFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PhoneVerificationFailedType,
		),
	}
}

func NewUserV1PhoneCodeAddedEvent(
	ctx context.Context,
	code *crypto.CryptoValue,
	expiry time.Duration,
) *phone.CodeAddedEvent {
	return &phone.CodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PhoneCodeAddedType,
		),
		Code:   code,
		Expiry: expiry,
	}
}

func NewUserV1PhoneCodeSentEvent(ctx context.Context) *phone.CodeSentEvent {
	return &phone.CodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PhoneCodeSentType,
		),
	}
}

func NewUserV1ProfileChangedEvent(
	ctx context.Context,
	firstName,
	lastName,
	nickName,
	displayName string,
	preferredLanguage language.Tag,
	gender human.Gender,
) *profile.ChangedEvent {
	return &profile.ChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1ProfileChangedType,
		),
		FirstName:         firstName,
		LastName:          lastName,
		NickName:          nickName,
		DisplayName:       displayName,
		PreferredLanguage: preferredLanguage,
		Gender:            gender,
	}
}

func NewUserV1AddressChangedEvent(
	ctx context.Context,
	country,
	locality,
	postalCode,
	region,
	streetAddress string,
) *address.ChangedEvent {
	return &address.ChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1AddressChangedType,
		),
		Country:       country,
		Locality:      locality,
		PostalCode:    postalCode,
		Region:        region,
		StreetAddress: streetAddress,
	}
}

func NewUserV1MFAInitSkippedEvent(ctx context.Context) *mfa.InitSkippedEvent {
	return &mfa.InitSkippedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1MFAInitSkippedType,
		),
	}
}

func NewUserV1MFAOTPAddedEvent(
	ctx context.Context,
	secret *crypto.CryptoValue,
) *otp.AddedEvent {
	return &otp.AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1MFAOTPAddedType,
		),
		Secret: secret,
	}
}

func NewUserV1MFAOTPVerifiedEvent(ctx context.Context) *otp.VerifiedEvent {
	return &otp.VerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1MFAOTPVerifiedType,
		),
	}
}

func NewUserV1MFAOTPRemovedEvent(ctx context.Context) *otp.RemovedEvent {
	return &otp.RemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1MFAOTPRemovedType,
		),
	}
}

func NewUserV1MFAOTPCheckSucceededEvent(ctx context.Context) *otp.CheckSucceededEvent {
	return &otp.CheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1MFAOTPCheckSucceededType,
		),
	}
}

func NewUserV1MFAOTPCheckFailedEvent(ctx context.Context) *otp.CheckFailedEvent {
	return &otp.CheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1MFAOTPCheckFailedType,
		),
	}
}
