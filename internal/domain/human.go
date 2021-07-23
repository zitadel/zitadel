package domain

import (
	"time"

	"github.com/caos/zitadel/internal/crypto"
	caos_errors "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

type Human struct {
	es_models.ObjectRoot

	Username string
	State    UserState
	*Password
	*Profile
	*Email
	*Phone
	*Address
	ExternalIDPs []*ExternalIDP
}

func (h Human) GetUsername() string {
	return h.Username
}

func (h Human) GetState() UserState {
	return h.State
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

	genderCount
)

func (f Gender) Valid() bool {
	return f >= 0 && f < genderCount
}

func (u *Human) IsValid() bool {
	return u.Username != "" && u.Profile != nil && u.Profile.IsValid() && u.Email != nil && u.Email.IsValid() && u.Phone == nil || (u.Phone != nil && u.Phone.PhoneNumber != "" && u.Phone.IsValid())
}

func (u *Human) CheckOrgIAMPolicy(policy *OrgIAMPolicy) error {
	if policy == nil {
		return caos_errors.ThrowPreconditionFailed(nil, "DOMAIN-zSH7j", "Errors.Users.OrgIamPolicyNil")
	}
	if !policy.UserLoginMustBeDomain && u.Profile != nil && u.Username == "" && u.Email != nil {
		u.Username = u.EmailAddress
	}
	return nil
}

func (u *Human) SetNamesAsDisplayname() {
	if u.Profile != nil && u.DisplayName == "" && u.FirstName != "" && u.LastName != "" {
		u.DisplayName = u.FirstName + " " + u.LastName
	}
}

func (u *Human) HashPasswordIfExisting(policy *PasswordComplexityPolicy, passwordAlg crypto.HashAlgorithm, onetime bool) error {
	if u.Password != nil {
		u.Password.ChangeRequired = onetime
		return u.Password.HashPasswordIfExisting(policy, passwordAlg)
	}
	return nil
}

func (u *Human) IsInitialState(passwordless bool) bool {
	return u.Email == nil || !u.IsEmailVerified || (u.ExternalIDPs == nil || len(u.ExternalIDPs) == 0) && !passwordless && (u.Password == nil || u.SecretString == "")
}

func NewInitUserCode(generator crypto.Generator) (*InitUserCode, error) {
	initCodeCrypto, _, err := crypto.NewCode(generator)
	if err != nil {
		return nil, err
	}
	return &InitUserCode{
		Code:   initCodeCrypto,
		Expiry: generator.Expiry(),
	}, nil
}

func GenerateLoginName(username, domain string, appendDomain bool) string {
	if !appendDomain {
		return username
	}
	return username + "@" + domain
}
