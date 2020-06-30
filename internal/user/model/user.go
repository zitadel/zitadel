package model

import (
	caos_errors "github.com/caos/zitadel/internal/errors"
	org_model "github.com/caos/zitadel/internal/org/model"
	policy_model "github.com/caos/zitadel/internal/policy/model"
	"github.com/golang/protobuf/ptypes/timestamp"
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
type UserChanges struct {
	Changes      []*UserChange
	LastSequence uint64
}

type UserChange struct {
	ChangeDate   *timestamp.Timestamp `json:"changeDate,omitempty"`
	EventType    string               `json:"eventType,omitempty"`
	Sequence     uint64               `json:"sequence,omitempty"`
	ModifierId   string               `json:"modifierUser,omitempty"`
	ModifierName string               `json:"-"`
	Data         interface{}          `json:"data,omitempty"`
}

type InitUserCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue
	Expiry time.Duration
}

type UserState int32

const (
	UserStateUnspecified UserState = iota
	UserStateActive
	UserStateInactive
	UserStateDeleted
	UserStateLocked
	UserStateSuspend
	UserStateInitial
)

type Gender int32

const (
	GenderUnspecified Gender = iota
	GenderFemale
	GenderMale
	GenderDiverse
)

func (u *User) CheckOrgIamPolicy(policy *org_model.OrgIamPolicy) error {
	if policy == nil {
		return caos_errors.ThrowPreconditionFailed(nil, "MODEL-zSH7j", "Errors.Users.OrgIamPolicyNil")
	}
	if policy.UserLoginMustBeDomain && strings.Contains(u.UserName, "@") {
		return caos_errors.ThrowPreconditionFailed(nil, "MODEL-se4sJ", "Errors.User.EmailAsUsernameNotAllowed")
	}
	if !policy.UserLoginMustBeDomain && u.Profile != nil && u.UserName == "" && u.Email != nil {
		u.UserName = u.EmailAddress
	}
	return nil
}

func (u *User) SetNamesAsDisplayname() {
	if u.Profile != nil && u.DisplayName == "" && u.FirstName != "" && u.LastName != "" {
		u.DisplayName = u.FirstName + " " + u.LastName
	}
}

func (u *User) IsValid() bool {
	return u.Profile != nil && u.FirstName != "" && u.LastName != "" && u.UserName != "" && u.Email != nil && u.Email.IsValid() && u.Phone == nil || (u.Phone != nil && u.Phone.IsValid())
}

func (u *User) IsInitialState() bool {
	return u.Email == nil || !u.IsEmailVerified || u.Password == nil || u.SecretString == ""
}

func (u *User) IsActive() bool {
	return u.State == UserStateActive
}

func (u *User) IsInitial() bool {
	return u.State == UserStateInitial
}

func (u *User) IsInactive() bool {
	return u.State == UserStateInactive
}

func (u *User) IsLocked() bool {
	return u.State == UserStateLocked
}

func (u *User) IsOTPReady() bool {
	return u.OTP != nil && u.OTP.State == MfaStateReady
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
