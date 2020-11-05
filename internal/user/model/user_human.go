package model

import (
	"github.com/duo-labs/webauthn/webauthn"
	"time"

	iam_model "github.com/caos/zitadel/internal/iam/model"

	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type Human struct {
	*Password
	*Profile
	*Email
	*Phone
	*Address
	ExternalIDPs []*ExternalIDP
	InitCode     *InitUserCode
	EmailCode    *EmailCode
	PhoneCode    *PhoneCode
	PasswordCode *PasswordCode
	OTP          *OTP
	U2Fs         []*WebauthNToken
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
	return u.OTP != nil && u.OTP.State == MfaStateReady
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

func (u *Human) GetU2F(webAuthNTokenID string) (int, *WebauthNToken) {
	for i, u2f := range u.U2Fs {
		if u2f.SessionID == webAuthNTokenID {
			return i, u2f
		}
	}
	return -1, nil
}

func (u *Human) GetU2FToVerify() (int, *WebauthNToken) {
	for i, u2f := range u.U2Fs {
		if u2f.State == MfaStateNotReady {
			return i, u2f
		}
	}
	return -1, nil
}

func (u *Human) GetU2FCredentials() []webauthn.Credential {
	creds := make([]webauthn.Credential, len(u.U2Fs))
	for i, u2f := range u.U2Fs {
		creds[i] = webauthn.Credential{
			ID:              u2f.KeyID,
			PublicKey:       u2f.PublicKey,
			AttestationType: u2f.AttestationType,
			Authenticator: webauthn.Authenticator{
				AAGUID:    u2f.AAGUID,
				SignCount: u2f.SignCount,
			},
		}
	}
	return creds
}
