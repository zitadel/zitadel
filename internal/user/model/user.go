package model

import (
	policy_model "github.com/caos/zitadel/internal/policy/model"
	"strings"
	"time"

	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type User struct {
	es_models.ObjectRoot

	State UserState
	*Password
	*Profile
	*Email
	*Phone
	*Address
	InitCode     *InitUserCode
	EmailCode    *EmailCode
	PhoneCode    *PhoneCode
	PasswordCode *PasswordCode
	OTP          *OTP
}

type InitUserCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue
	Expiry time.Duration
}

type UserState int32

const (
	USERSTATE_UNSPECIFIED UserState = iota
	USERSTATE_ACTIVE
	USERSTATE_INACTIVE
	USERSTATE_DELETED
	USERSTATE_LOCKED
	USERSTATE_SUSPEND
	USERSTATE_INITIAL
)

type Gender int32

const (
	GENDER_UNDEFINED Gender = iota
	GENDER_FEMALE
	GENDER_MALE
	GENDER_DIVERSE
)

func (u *User) SetEmailAsUsername(orgDomain string) {
	if u.Profile != nil && u.UserName == "" && u.Email != nil {
		u.UserName = strings.Split(u.EmailAddress, "@")[0]
	}

	u.UserName += "@" + orgDomain
}

func (u *User) IsValid() bool {
	return u.Profile != nil && u.FirstName != "" && u.LastName != "" && u.UserName != "" && u.Email != nil && u.Email.IsValid() && u.Phone == nil || (u.Phone != nil && u.Phone.IsValid())
}

func (u *User) IsInitialState() bool {
	return u.Email == nil || !u.IsEmailVerified || u.Password == nil || u.SecretString == ""
}

func (u *User) IsActive() bool {
	return u.State == USERSTATE_ACTIVE
}

func (u *User) IsInitial() bool {
	return u.State == USERSTATE_INITIAL
}

func (u *User) IsInactive() bool {
	return u.State == USERSTATE_INACTIVE
}

func (u *User) IsLocked() bool {
	return u.State == USERSTATE_LOCKED
}

func (u *User) IsOTPReady() bool {
	return u.OTP != nil && u.OTP.State == MFASTATE_READY
}

func (u *User) HashPasswordIfExisting(policy *policy_model.PasswordComplexityPolicy, passwordAlg crypto.HashAlgorithm, onetime bool) error {
	if u.Password != nil {
		return u.Password.HashPasswordIfExisting(policy, passwordAlg, onetime)
	}
	return nil
}

func (u *User) GenerateInitCodeIfNeeded(initGenerator crypto.Generator) error {
	if !u.IsInitialState() {
		return nil
	}
	u.InitCode = new(InitUserCode)
	return u.InitCode.GenerateInitUserCode(initGenerator)
}

func (u *User) GeneratePhoneCodeIfNeeded(phoneGenerator crypto.Generator) error {
	if u.Phone == nil || u.IsPhoneVerified {
		return nil
	}
	u.PhoneCode = new(PhoneCode)
	return u.PhoneCode.GeneratePhoneCode(phoneGenerator)
}

func (u *User) GenerateEmailCodeIfNeeded(emailGenerator crypto.Generator) error {
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
