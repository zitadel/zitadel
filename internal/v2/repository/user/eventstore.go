package user

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/user/human"
	"github.com/caos/zitadel/internal/v2/repository/user/human/address"
	"github.com/caos/zitadel/internal/v2/repository/user/human/email"
	"github.com/caos/zitadel/internal/v2/repository/user/human/external_idp"
	"github.com/caos/zitadel/internal/v2/repository/user/human/mfa"
	"github.com/caos/zitadel/internal/v2/repository/user/human/mfa/otp"
	"github.com/caos/zitadel/internal/v2/repository/user/human/mfa/web_auth_n"
	"github.com/caos/zitadel/internal/v2/repository/user/human/password"
	"github.com/caos/zitadel/internal/v2/repository/user/human/phone"
	"github.com/caos/zitadel/internal/v2/repository/user/human/profile"
	"github.com/caos/zitadel/internal/v2/repository/user/machine"
	"github.com/caos/zitadel/internal/v2/repository/user/machine/keys"
	"github.com/caos/zitadel/internal/v2/repository/user/v1"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(v1.UserV1AddedType, human.HumanAddedEventMapper).
		RegisterFilterEventMapper(v1.UserV1RegisteredType, human.HumanRegisteredEventMapper).
		RegisterFilterEventMapper(v1.UserV1InitialCodeAddedType, human.HumanInitialCodeAddedEventMapper).
		RegisterFilterEventMapper(v1.UserV1InitialCodeSentType, human.HumanInitialCodeSentEventMapper).
		RegisterFilterEventMapper(v1.UserV1InitializedCheckSucceededType, human.HumanInitializedCheckSucceededEventMapper).
		RegisterFilterEventMapper(v1.UserV1InitializedCheckFailedType, human.HumanInitializedCheckFailedEventMapper).
		RegisterFilterEventMapper(v1.UserV1SignedOutType, human.HumanSignedOutEventMapper).
		RegisterFilterEventMapper(v1.UserV1PasswordChangedType, password.HumanPasswordChangedEventMapper).
		RegisterFilterEventMapper(v1.UserV1PasswordCodeAddedType, password.HumanPasswordCodeAddedEventMapper).
		RegisterFilterEventMapper(v1.UserV1PasswordCodeSentType, password.HumanPasswordCodeSentEventMapper).
		RegisterFilterEventMapper(v1.UserV1PasswordCheckSucceededType, password.HumanPasswordCheckSucceededEventMapper).
		RegisterFilterEventMapper(v1.UserV1PasswordCheckFailedType, password.HumanPasswordCheckFailedEventMapper).
		RegisterFilterEventMapper(v1.UserV1EmailChangedType, email.HumanEmailChangedEventMapper).
		RegisterFilterEventMapper(v1.UserV1EmailVerifiedType, email.HumanEmailVerifiedEventMapper).
		RegisterFilterEventMapper(v1.UserV1EmailVerificationFailedType, email.HumanEmailVerificationFailedEventMapper).
		RegisterFilterEventMapper(v1.UserV1EmailCodeAddedType, email.HumanEmailCodeAddedEventMapper).
		RegisterFilterEventMapper(v1.UserV1EmailCodeSentType, email.HumanEmailCodeSentEventMapper).
		RegisterFilterEventMapper(v1.UserV1PhoneChangedType, phone.HumanPhoneChangedEventMapper).
		RegisterFilterEventMapper(v1.UserV1PhoneRemovedType, phone.HumanPhoneRemovedEventMapper).
		RegisterFilterEventMapper(v1.UserV1PhoneVerifiedType, phone.HumanPhoneVerifiedEventMapper).
		RegisterFilterEventMapper(v1.UserV1PhoneVerificationFailedType, phone.HumanPhoneVerificationFailedEventMapper).
		RegisterFilterEventMapper(v1.UserV1PhoneCodeAddedType, phone.HumanPhoneCodeAddedEventMapper).
		RegisterFilterEventMapper(v1.UserV1PhoneCodeSentType, phone.HumanPhoneCodeSentEventMapper).
		RegisterFilterEventMapper(v1.UserV1ProfileChangedType, profile.HumanProfileChangedEventMapper).
		RegisterFilterEventMapper(v1.UserV1AddressChangedType, address.HumanAddressChangedEventMapper).
		RegisterFilterEventMapper(v1.UserV1MFAInitSkippedType, mfa.HumanMFAInitSkippedEventMapper).
		RegisterFilterEventMapper(v1.UserV1MFAOTPAddedType, otp.HumanMFAOTPAddedEventMapper).
		RegisterFilterEventMapper(v1.UserV1MFAOTPVerifiedType, otp.HumanMFAOTPVerifiedEventMapper).
		RegisterFilterEventMapper(v1.UserV1MFAOTPRemovedType, otp.HumanMFAOTPRemovedEventMapper).
		RegisterFilterEventMapper(v1.UserV1MFAOTPCheckSucceededType, otp.HumanMFAOTPCheckSucceededEventMapper).
		RegisterFilterEventMapper(v1.UserV1MFAOTPCheckFailedType, otp.HumanMFAOTPCheckFailedEventMapper).
		RegisterFilterEventMapper(UserLockedType, UserLockedEventMapper).
		RegisterFilterEventMapper(UserUnlockedType, UserLockedEventMapper).
		RegisterFilterEventMapper(UserDeactivatedType, UserDeactivatedEventMapper).
		RegisterFilterEventMapper(UserReactivatedType, UserReactivatedEventMapper).
		RegisterFilterEventMapper(UserRemovedType, UserRemovedEventMapper).
		RegisterFilterEventMapper(UserTokenAddedType, UserTokenAddedEventMapper).
		RegisterFilterEventMapper(UserDomainClaimedType, UserDomainClaimedEventMapper).
		RegisterFilterEventMapper(UserDomainClaimedSentType, UserDomainClaimedEventMapper).
		RegisterFilterEventMapper(UserUserNameChangedType, UserUsernameChangedEventMapper).
		RegisterFilterEventMapper(human.HumanAddedType, human.HumanAddedEventMapper).
		RegisterFilterEventMapper(human.HumanRegisteredType, human.HumanRegisteredEventMapper).
		RegisterFilterEventMapper(human.HumanInitialCodeAddedType, human.HumanInitialCodeAddedEventMapper).
		RegisterFilterEventMapper(human.HumanInitialCodeSentType, human.HumanInitialCodeSentEventMapper).
		RegisterFilterEventMapper(human.HumanInitializedCheckSucceededType, human.HumanInitializedCheckSucceededEventMapper).
		RegisterFilterEventMapper(human.HumanInitializedCheckFailedType, human.HumanInitializedCheckFailedEventMapper).
		RegisterFilterEventMapper(human.HumanSignedOutType, human.HumanSignedOutEventMapper).
		RegisterFilterEventMapper(password.HumanPasswordChangedType, password.HumanPasswordChangedEventMapper).
		RegisterFilterEventMapper(password.HumanPasswordCodeAddedType, password.HumanPasswordCodeAddedEventMapper).
		RegisterFilterEventMapper(password.HumanPasswordCodeSentType, password.HumanPasswordCodeSentEventMapper).
		RegisterFilterEventMapper(password.HumanPasswordCheckSucceededType, password.HumanPasswordCheckSucceededEventMapper).
		RegisterFilterEventMapper(password.HumanPasswordCheckFailedType, password.HumanPasswordCheckFailedEventMapper).
		RegisterFilterEventMapper(external_idp.HumanExternalIDPAddedType, external_idp.HumanExternalIDPAddedEventMapper).
		RegisterFilterEventMapper(external_idp.HumanExternalIDPRemovedType, external_idp.HumanExternalIDPRemovedEventMapper).
		RegisterFilterEventMapper(external_idp.HumanExternalIDPCascadeRemovedType, external_idp.HumanExternalIDPCascadeRemovedEventMapper).
		RegisterFilterEventMapper(external_idp.HumanExternalLoginCheckSucceededType, external_idp.HumanExternalLoginCheckSucceededEventMapper).
		RegisterFilterEventMapper(email.HumanEmailChangedType, email.HumanEmailChangedEventMapper).
		RegisterFilterEventMapper(email.HumanEmailVerifiedType, email.HumanEmailVerifiedEventMapper).
		RegisterFilterEventMapper(email.HumanEmailVerificationFailedType, email.HumanEmailVerificationFailedEventMapper).
		RegisterFilterEventMapper(email.HumanEmailCodeAddedType, email.HumanEmailCodeAddedEventMapper).
		RegisterFilterEventMapper(email.HumanEmailCodeSentType, email.HumanEmailCodeSentEventMapper).
		RegisterFilterEventMapper(phone.HumanPhoneChangedType, phone.HumanPhoneChangedEventMapper).
		RegisterFilterEventMapper(phone.HumanPhoneRemovedType, phone.HumanPhoneRemovedEventMapper).
		RegisterFilterEventMapper(phone.HumanPhoneVerifiedType, phone.HumanPhoneVerifiedEventMapper).
		RegisterFilterEventMapper(phone.HumanPhoneVerificationFailedType, phone.HumanPhoneVerificationFailedEventMapper).
		RegisterFilterEventMapper(phone.HumanPhoneCodeAddedType, phone.HumanPhoneCodeAddedEventMapper).
		RegisterFilterEventMapper(phone.HumanPhoneCodeSentType, phone.HumanPhoneCodeSentEventMapper).
		RegisterFilterEventMapper(profile.HumanProfileChangedType, profile.HumanProfileChangedEventMapper).
		RegisterFilterEventMapper(address.HumanAddressChangedType, address.HumanAddressChangedEventMapper).
		RegisterFilterEventMapper(mfa.HumanMFAInitSkippedType, mfa.HumanMFAInitSkippedEventMapper).
		RegisterFilterEventMapper(otp.HumanMFAOTPAddedType, otp.HumanMFAOTPAddedEventMapper).
		RegisterFilterEventMapper(otp.HumanMFAOTPVerifiedType, otp.HumanMFAOTPVerifiedEventMapper).
		RegisterFilterEventMapper(otp.HumanMFAOTPRemovedType, otp.HumanMFAOTPRemovedEventMapper).
		RegisterFilterEventMapper(otp.HumanMFAOTPCheckSucceededType, otp.HumanMFAOTPCheckSucceededEventMapper).
		RegisterFilterEventMapper(otp.HumanMFAOTPCheckFailedType, otp.HumanMFAOTPCheckFailedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanU2FTokenAddedType, web_auth_n.HumanWebAuthNAddedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanU2FTokenVerifiedType, web_auth_n.HumanWebAuthNVerifiedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanU2FTokenSignCountChangedType, web_auth_n.HumanWebAuthNSignCountChangedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanU2FTokenRemovedType, web_auth_n.HumanWebAuthNRemovedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanU2FTokenBeginLoginType, web_auth_n.HumanWebAuthNBeginLoginEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanU2FTokenCheckSucceededType, web_auth_n.HumanWebAuthNCheckSucceededEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanU2FTokenCheckFailedType, web_auth_n.HumanWebAuthNCheckFailedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanPasswordlessTokenAddedType, web_auth_n.HumanWebAuthNAddedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanPasswordlessTokenVerifiedType, web_auth_n.HumanWebAuthNVerifiedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanPasswordlessTokenSignCountChangedType, web_auth_n.HumanWebAuthNSignCountChangedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanPasswordlessTokenRemovedType, web_auth_n.HumanWebAuthNRemovedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanPasswordlessTokenBeginLoginType, web_auth_n.HumanWebAuthNBeginLoginEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanPasswordlessTokenCheckSucceededType, web_auth_n.HumanWebAuthNCheckSucceededEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanPasswordlessTokenCheckFailedType, web_auth_n.HumanWebAuthNCheckFailedEventMapper).
		RegisterFilterEventMapper(machine.MachineAddedEventType, machine.MachineAddedEventMapper).
		RegisterFilterEventMapper(machine.MachineChangedEventType, machine.MachineChangedEventMapper).
		RegisterFilterEventMapper(keys.MachineKeyAddedEventType, keys.MachineKeyAddedEventMapper).
		RegisterFilterEventMapper(keys.MachineKeyRemovedEventType, keys.MachineKeyRemovedEventMapper)
}
