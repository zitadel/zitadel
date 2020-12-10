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
	es.RegisterFilterEventMapper(v1.UserV1AddedType, human.AddedEventMapper).
		RegisterFilterEventMapper(v1.UserV1RegisteredType, human.RegisteredEventMapper).
		RegisterFilterEventMapper(v1.UserV1InitialCodeAddedType, human.InitialCodeAddedEventMapper).
		RegisterFilterEventMapper(v1.UserV1InitialCodeSentType, human.InitialCodeSentEventMapper).
		RegisterFilterEventMapper(v1.UserV1InitializedCheckSucceededType, human.InitializedCheckSucceededEventMapper).
		RegisterFilterEventMapper(v1.UserV1InitializedCheckFailedType, human.InitializedCheckFailedEventMapper).
		RegisterFilterEventMapper(v1.UserV1SignedOutType, human.SignedOutEventMapper).
		RegisterFilterEventMapper(v1.UserV1PasswordChangedType, password.ChangedEventMapper).
		RegisterFilterEventMapper(v1.UserV1PasswordCodeAddedType, password.CodeAddedEventMapper).
		RegisterFilterEventMapper(v1.UserV1PasswordCodeSentType, password.CodeSentEventMapper).
		RegisterFilterEventMapper(v1.UserV1PasswordCheckSucceededType, password.CheckSucceededEventMapper).
		RegisterFilterEventMapper(v1.UserV1PasswordCheckFailedType, password.CheckFailedEventMapper).
		RegisterFilterEventMapper(v1.UserV1EmailChangedType, email.ChangedEventMapper).
		RegisterFilterEventMapper(v1.UserV1EmailVerifiedType, email.VerifiedEventMapper).
		RegisterFilterEventMapper(v1.UserV1EmailVerificationFailedType, email.VerificationFailedEventMapper).
		RegisterFilterEventMapper(v1.UserV1EmailCodeAddedType, email.CodeAddedEventMapper).
		RegisterFilterEventMapper(v1.UserV1EmailCodeSentType, email.CodeSentEventMapper).
		RegisterFilterEventMapper(v1.UserV1PhoneChangedType, phone.ChangedEventMapper).
		RegisterFilterEventMapper(v1.UserV1PhoneRemovedType, phone.RemovedEventMapper).
		RegisterFilterEventMapper(v1.UserV1PhoneVerifiedType, phone.VerifiedEventMapper).
		RegisterFilterEventMapper(v1.UserV1PhoneVerificationFailedType, phone.VerificationFailedEventMapper).
		RegisterFilterEventMapper(v1.UserV1PhoneCodeAddedType, phone.CodeAddedEventMapper).
		RegisterFilterEventMapper(v1.UserV1PhoneCodeSentType, phone.CodeSentEventMapper).
		RegisterFilterEventMapper(v1.UserV1ProfileChangedType, profile.ChangedEventMapper).
		RegisterFilterEventMapper(v1.UserV1AddressChangedType, address.ChangedEventMapper).
		RegisterFilterEventMapper(v1.UserV1MFAInitSkippedType, mfa.InitSkippedEventMapper).
		RegisterFilterEventMapper(v1.UserV1MFAOTPAddedType, otp.AddedEventMapper).
		RegisterFilterEventMapper(v1.UserV1MFAOTPVerifiedType, otp.VerifiedEventMapper).
		RegisterFilterEventMapper(v1.UserV1MFAOTPRemovedType, otp.RemovedEventMapper).
		RegisterFilterEventMapper(v1.UserV1MFAOTPCheckSucceededType, otp.CheckSucceededEventMapper).
		RegisterFilterEventMapper(v1.UserV1MFAOTPCheckFailedType, otp.CheckFailedEventMapper).
		RegisterFilterEventMapper(UserLockedType, LockedEventMapper).
		RegisterFilterEventMapper(UserUnlockedType, LockedEventMapper).
		RegisterFilterEventMapper(UserDeactivatedType, DeactivatedEventMapper).
		RegisterFilterEventMapper(UserReactivatedType, ReactivatedEventMapper).
		RegisterFilterEventMapper(UserRemovedType, RemovedEventMapper).
		RegisterFilterEventMapper(UserTokenAddedType, TokenAddedEventMapper).
		RegisterFilterEventMapper(UserDomainClaimedType, DomainClaimedEventMapper).
		RegisterFilterEventMapper(UserDomainClaimedSentType, DomainClaimedEventMapper).
		RegisterFilterEventMapper(UserUserNameChangedType, UsernameChangedEventMapper).
		RegisterFilterEventMapper(human.HumanAddedType, human.AddedEventMapper).
		RegisterFilterEventMapper(human.HumanRegisteredType, human.RegisteredEventMapper).
		RegisterFilterEventMapper(human.HumanInitialCodeAddedType, human.InitialCodeAddedEventMapper).
		RegisterFilterEventMapper(human.HumanInitialCodeSentType, human.InitialCodeSentEventMapper).
		RegisterFilterEventMapper(human.HumanInitializedCheckSucceededType, human.InitializedCheckSucceededEventMapper).
		RegisterFilterEventMapper(human.HumanInitializedCheckFailedType, human.InitializedCheckFailedEventMapper).
		RegisterFilterEventMapper(human.HumanSignedOutType, human.SignedOutEventMapper).
		RegisterFilterEventMapper(password.HumanPasswordChangedType, password.ChangedEventMapper).
		RegisterFilterEventMapper(password.HumanPasswordCodeAddedType, password.CodeAddedEventMapper).
		RegisterFilterEventMapper(password.HumanPasswordCodeSentType, password.CodeSentEventMapper).
		RegisterFilterEventMapper(password.HumanPasswordCheckSucceededType, password.CheckSucceededEventMapper).
		RegisterFilterEventMapper(password.HumanPasswordCheckFailedType, password.CheckFailedEventMapper).
		RegisterFilterEventMapper(external_idp.HumanExternalIDPAddedType, external_idp.AddedEventMapper).
		RegisterFilterEventMapper(external_idp.HumanExternalIDPRemovedType, external_idp.RemovedEventMapper).
		RegisterFilterEventMapper(external_idp.HumanExternalIDPCascadeRemovedType, external_idp.CascadeRemovedEventMapper).
		RegisterFilterEventMapper(external_idp.HumanExternalLoginCheckSucceededType, external_idp.CheckSucceededEventMapper).
		RegisterFilterEventMapper(email.HumanEmailChangedType, email.ChangedEventMapper).
		RegisterFilterEventMapper(email.HumanEmailVerifiedType, email.VerifiedEventMapper).
		RegisterFilterEventMapper(email.HumanEmailVerificationFailedType, email.VerificationFailedEventMapper).
		RegisterFilterEventMapper(email.HumanEmailCodeAddedType, email.CodeAddedEventMapper).
		RegisterFilterEventMapper(email.HumanEmailCodeSentType, email.CodeSentEventMapper).
		RegisterFilterEventMapper(phone.HumanPhoneChangedType, phone.ChangedEventMapper).
		RegisterFilterEventMapper(phone.HumanPhoneRemovedType, phone.RemovedEventMapper).
		RegisterFilterEventMapper(phone.HumanPhoneVerifiedType, phone.VerifiedEventMapper).
		RegisterFilterEventMapper(phone.HumanPhoneVerificationFailedType, phone.VerificationFailedEventMapper).
		RegisterFilterEventMapper(phone.HumanPhoneCodeAddedType, phone.CodeAddedEventMapper).
		RegisterFilterEventMapper(phone.HumanPhoneCodeSentType, phone.CodeSentEventMapper).
		RegisterFilterEventMapper(profile.HumanProfileChangedType, profile.ChangedEventMapper).
		RegisterFilterEventMapper(address.HumanAddressChangedType, address.ChangedEventMapper).
		RegisterFilterEventMapper(mfa.HumanMFAInitSkippedType, mfa.InitSkippedEventMapper).
		RegisterFilterEventMapper(otp.HumanMFAOTPAddedType, otp.AddedEventMapper).
		RegisterFilterEventMapper(otp.HumanMFAOTPVerifiedType, otp.VerifiedEventMapper).
		RegisterFilterEventMapper(otp.HumanMFAOTPRemovedType, otp.RemovedEventMapper).
		RegisterFilterEventMapper(otp.HumanMFAOTPCheckSucceededType, otp.CheckSucceededEventMapper).
		RegisterFilterEventMapper(otp.HumanMFAOTPCheckFailedType, otp.CheckFailedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanU2FTokenAddedType, web_auth_n.AddedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanU2FTokenVerifiedType, web_auth_n.VerifiedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanU2FTokenSignCountChangedType, web_auth_n.SignCountChangedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanU2FTokenRemovedType, web_auth_n.RemovedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanU2FTokenBeginLoginType, web_auth_n.BeginLoginEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanU2FTokenCheckSucceededType, web_auth_n.CheckSucceededEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanU2FTokenCheckFailedType, web_auth_n.CheckFailedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanPasswordlessTokenAddedType, web_auth_n.AddedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanPasswordlessTokenVerifiedType, web_auth_n.VerifiedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanPasswordlessTokenSignCountChangedType, web_auth_n.SignCountChangedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanPasswordlessTokenRemovedType, web_auth_n.RemovedEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanPasswordlessTokenBeginLoginType, web_auth_n.BeginLoginEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanPasswordlessTokenCheckSucceededType, web_auth_n.CheckSucceededEventMapper).
		RegisterFilterEventMapper(web_auth_n.HumanPasswordlessTokenCheckFailedType, web_auth_n.CheckFailedEventMapper).
		RegisterFilterEventMapper(machine.MachineAddedEventType, machine.AddedEventMapper).
		RegisterFilterEventMapper(machine.MachineChangedEventType, machine.ChangedEventMapper).
		RegisterFilterEventMapper(keys.MachineKeyAddedEventType, keys.AddedEventMapper).
		RegisterFilterEventMapper(keys.MachineKeyRemovedEventType, keys.RemovedEventMapper)
}
