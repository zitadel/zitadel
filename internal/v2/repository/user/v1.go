package user

import (
	"context"
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
	"time"
)

const (
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
) *human.HumanAddedEvent {
	return &human.HumanAddedEvent{
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
) *human.HumanRegisteredEvent {
	return &human.HumanRegisteredEvent{
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
) *human.HumanInitialCodeAddedEvent {
	return &human.HumanInitialCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1InitialCodeAddedType,
		),
		Code:   code,
		Expiry: expiry,
	}
}

func NewUserV1InitialCodeSentEvent(ctx context.Context) *human.HumanInitialCodeSentEvent {
	return &human.HumanInitialCodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1InitialCodeSentType,
		),
	}
}

func NewUserV1InitializedCheckSucceededEvent(ctx context.Context) *human.HumanInitializedCheckSucceededEvent {
	return &human.HumanInitializedCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1InitializedCheckSucceededType,
		),
	}
}

func NewUserV1InitializedCheckFailedEvent(ctx context.Context) *human.HumanInitializedCheckFailedEvent {
	return &human.HumanInitializedCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1InitializedCheckFailedType,
		),
	}
}

func NewUserV1SignedOutEvent(ctx context.Context) *human.HumanSignedOutEvent {
	return &human.HumanSignedOutEvent{
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
) *password.HumanPasswordChangedEvent {
	return &password.HumanPasswordChangedEvent{
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
	notificationType NotificationType,
) *password.HumanPasswordCodeAddedEvent {
	return &password.HumanPasswordCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PasswordCodeAddedType,
		),
		Code:             code,
		Expiry:           expiry,
		NotificationType: notificationType,
	}
}

func NewUserV1PasswordCodeSentEvent(ctx context.Context) *password.HumanPasswordCodeSentEvent {
	return &password.HumanPasswordCodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PasswordCodeSentType,
		),
	}
}

func NewUserV1PasswordCheckSucceededEvent(ctx context.Context) *password.HumanPasswordCheckSucceededEvent {
	return &password.HumanPasswordCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PasswordCheckSucceededType,
		),
	}
}

func NewUserV1PasswordCheckFailedEvent(ctx context.Context) *password.HumanPasswordCheckFailedEvent {
	return &password.HumanPasswordCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PasswordCheckFailedType,
		),
	}
}

func NewUserV1EmailChangedEvent(ctx context.Context, emailAddress string) *email.HumanEmailChangedEvent {
	return &email.HumanEmailChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1EmailChangedType,
		),
		EmailAddress: emailAddress,
	}
}

func NewUserV1EmailVerifiedEvent(ctx context.Context) *email.HumanEmailVerifiedEvent {
	return &email.HumanEmailVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1EmailVerifiedType,
		),
	}
}

func NewUserV1EmailVerificationFailedEvent(ctx context.Context) *email.HumanEmailVerificationFailedEvent {
	return &email.HumanEmailVerificationFailedEvent{
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
) *email.HumanEmailCodeAddedEvent {
	return &email.HumanEmailCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1EmailCodeAddedType,
		),
		Code:   code,
		Expiry: expiry,
	}
}

func NewUserV1EmailCodeSentEvent(ctx context.Context) *email.HumanEmailCodeSentEvent {
	return &email.HumanEmailCodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1EmailCodeSentType,
		),
	}
}

func NewUserV1PhoneChangedEvent(ctx context.Context, phone string) *phone.HumanPhoneChangedEvent {
	return phone.HumanPhoneChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PhoneChangedType,
		),
		PhoneNumber: phone,
	}
}

func NewUserV1PhoneRemovedEvent(ctx context.Context) *phone.HumanPhoneRemovedEvent {
	return &phone.HumanPhoneRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PhoneRemovedType,
		),
	}
}

func NewUserV1PhoneVerifiedEvent(ctx context.Context) *phone.HumanPhoneVerifiedEvent {
	return &phone.HumanPhoneVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PhoneVerifiedType,
		),
	}
}

func NewUserV1PhoneVerificationFailedEvent(ctx context.Context) *phone.HumanPhoneVerificationFailedEvent {
	return &phone.HumanPhoneVerificationFailedEvent{
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
) *phone.HumanPhoneCodeAddedEvent {
	return &phone.HumanPhoneCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1PhoneCodeAddedType,
		),
		Code:   code,
		Expiry: expiry,
	}
}

func NewUserV1PhoneCodeSentEvent(ctx context.Context) *phone.HumanPhoneCodeSentEvent {
	return &phone.HumanPhoneCodeSentEvent{
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
) *profile.HumanProfileChangedEvent {
	return &profile.HumanProfileChangedEvent{
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
) *address.HumanAddressChangedEvent {
	return &address.HumanAddressChangedEvent{
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

func NewUserV1MFAInitSkippedEvent(ctx context.Context) *mfa.HumanMFAInitSkippedEvent {
	return &mfa.HumanMFAInitSkippedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1MFAInitSkippedType,
		),
	}
}

func NewUserV1MFAOTPAddedEvent(
	ctx context.Context,
	secret *crypto.CryptoValue,
) *otp.HumanMFAOTPAddedEvent {
	return &otp.HumanMFAOTPAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1MFAOTPAddedType,
		),
		Secret: secret,
	}
}

func NewUserV1MFAOTPVerifiedEvent(ctx context.Context) *otp.HumanMFAOTPVerifiedEvent {
	return &otp.HumanMFAOTPVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1MFAOTPVerifiedType,
		),
	}
}

func NewUserV1MFAOTPRemovedEvent(ctx context.Context) *otp.HumanMFAOTPRemovedEvent {
	return &otp.HumanMFAOTPRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1MFAOTPRemovedType,
		),
	}
}

func NewUserV1MFAOTPCheckSucceededEvent(ctx context.Context) *otp.HumanMFAOTPCheckSucceededEvent {
	return &otp.HumanMFAOTPCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1MFAOTPCheckSucceededType,
		),
	}
}

func NewUserV1MFAOTPCheckFailedEvent(ctx context.Context) *otp.HumanMFAOTPCheckFailedEvent {
	return &otp.HumanMFAOTPCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserV1MFAOTPCheckFailedType,
		),
	}
}
