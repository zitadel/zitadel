package model

import (
	"bytes"
	caos_errors "github.com/caos/zitadel/internal/errors"
	"strings"
	"time"

	iam_model "github.com/caos/zitadel/internal/iam/model"

	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type Human struct {
	es_models.ObjectRoot

	*Password
	*Profile
	*Email
	*Phone
	*Address
	ExternalIDPs       []*ExternalIDP
	InitCode           *InitUserCode
	EmailCode          *EmailCode
	PhoneCode          *PhoneCode
	PasswordCode       *PasswordCode
	OTP                *OTP
	U2FTokens          []*WebAuthNToken
	PasswordlessTokens []*WebAuthNToken
	U2FLogins          []*WebAuthNLogin
	PasswordlessLogins []*WebAuthNLogin
}

type InitUserCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue
	Expiry time.Duration
}

type Gender int32

const (
	GenderUnspecified Gender = iota
	GenderFemale
	GenderMale
	GenderDiverse
)

func (u *Human) CheckOrgIAMPolicy(userName string, policy *iam_model.OrgIAMPolicy) error {
	if policy == nil {
		return caos_errors.ThrowPreconditionFailed(nil, "MODEL-zSH7j", "Errors.Users.OrgIamPolicyNil")
	}
	if policy.UserLoginMustBeDomain && strings.Contains(userName, "@") {
		return caos_errors.ThrowPreconditionFailed(nil, "MODEL-se4sJ", "Errors.User.EmailAsUsernameNotAllowed")
	}
	if !policy.UserLoginMustBeDomain && u.Profile != nil && userName == "" && u.Email != nil {
		userName = u.EmailAddress
	}
	return nil
}

func (u *Human) SetNamesAsDisplayname() {
	if u.Profile != nil && u.DisplayName == "" && u.FirstName != "" && u.LastName != "" {
		u.DisplayName = u.FirstName + " " + u.LastName
	}
}

func (u *Human) IsValid() bool {
	return u.Profile != nil && u.FirstName != "" && u.LastName != "" && u.Email != nil && u.Email.IsValid() && u.Phone == nil || (u.Phone != nil && u.Phone.PhoneNumber != "" && u.Phone.IsValid())
}

func (u *Human) IsInitialState() bool {
	return u.Email == nil || !u.IsEmailVerified || (u.ExternalIDPs == nil || len(u.ExternalIDPs) == 0) && (u.Password == nil || u.SecretString == "")
}

func (u *Human) IsOTPReady() bool {
	return u.OTP != nil && u.OTP.State == MFAStateReady
}

func (u *Human) HashPasswordIfExisting(policy *iam_model.PasswordComplexityPolicyView, passwordAlg crypto.HashAlgorithm, onetime bool) error {
	if u.Password != nil {
		return u.Password.HashPasswordIfExisting(policy, passwordAlg, onetime)
	}
	return nil
}

func (u *Human) GenerateInitCodeIfNeeded(initGenerator crypto.Generator) error {
	if !u.IsInitialState() {
		return nil
	}
	u.InitCode = new(InitUserCode)
	return u.InitCode.GenerateInitUserCode(initGenerator)
}

func (u *Human) GeneratePhoneCodeIfNeeded(phoneGenerator crypto.Generator) error {
	if u.Phone == nil || u.IsPhoneVerified {
		return nil
	}
	u.PhoneCode = new(PhoneCode)
	return u.PhoneCode.GeneratePhoneCode(phoneGenerator)
}

func (u *Human) GenerateEmailCodeIfNeeded(emailGenerator crypto.Generator) error {
	if u.Email == nil || u.IsEmailVerified {
		return nil
	}
	u.EmailCode = new(EmailCode)
	return u.EmailCode.GenerateEmailCode(emailGenerator)
}

func (init *InitUserCode) GenerateInitUserCode(generator crypto.Generator) error {
	initCodeCrypto, _, err := crypto.NewCode(generator)
	if err != nil {
		return err
	}
	init.Code = initCodeCrypto
	init.Expiry = generator.Expiry()
	return nil
}

func (u *Human) GetExternalIDP(externalIDP *ExternalIDP) (int, *ExternalIDP) {
	for i, idp := range u.ExternalIDPs {
		if idp.UserID == externalIDP.UserID {
			return i, idp
		}
	}
	return -1, nil
}

func (u *Human) GetU2F(webAuthNTokenID string) (int, *WebAuthNToken) {
	for i, u2f := range u.U2FTokens {
		if u2f.WebAuthNTokenID == webAuthNTokenID {
			return i, u2f
		}
	}
	return -1, nil
}

func (u *Human) GetU2FByKeyID(keyID []byte) (int, *WebAuthNToken) {
	for i, u2f := range u.U2FTokens {
		if bytes.Compare(u2f.KeyID, keyID) == 0 {
			return i, u2f
		}
	}
	return -1, nil
}

func (u *Human) GetU2FToVerify() (int, *WebAuthNToken) {
	for i, u2f := range u.U2FTokens {
		if u2f.State == MFAStateNotReady {
			return i, u2f
		}
	}
	return -1, nil
}

func (u *Human) GetPasswordless(webAuthNTokenID string) (int, *WebAuthNToken) {
	for i, u2f := range u.PasswordlessTokens {
		if u2f.WebAuthNTokenID == webAuthNTokenID {
			return i, u2f
		}
	}
	return -1, nil
}

func (u *Human) GetPasswordlessByKeyID(keyID []byte) (int, *WebAuthNToken) {
	for i, pwl := range u.PasswordlessTokens {
		if bytes.Compare(pwl.KeyID, keyID) == 0 {
			return i, pwl
		}
	}
	return -1, nil
}

func (u *Human) GetPasswordlessToVerify() (int, *WebAuthNToken) {
	for i, u2f := range u.PasswordlessTokens {
		if u2f.State == MFAStateNotReady {
			return i, u2f
		}
	}
	return -1, nil
}

func (u *Human) GetU2FLogin(authReqID string) (int, *WebAuthNLogin) {
	for i, u2f := range u.U2FLogins {
		if u2f.AuthRequest.ID == authReqID {
			return i, u2f
		}
	}
	return -1, nil
}

func (u *Human) GetPasswordlessLogin(authReqID string) (int, *WebAuthNLogin) {
	for i, pw := range u.PasswordlessLogins {
		if pw.AuthRequest.ID == authReqID {
			return i, pw
		}
	}
	return -1, nil
}
