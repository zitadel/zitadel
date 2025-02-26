package model

import (
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type UserView struct {
	ID                 string
	UserName           string
	CreationDate       time.Time
	ChangeDate         time.Time
	State              UserState
	Sequence           uint64
	ResourceOwner      string
	LastLogin          time.Time
	PreferredLoginName string
	LoginNames         []string
	*MachineView
	*HumanView
}

type HumanView struct {
	PasswordSet              bool
	PasswordInitRequired     bool
	PasswordChangeRequired   bool
	UsernameChangeRequired   bool
	PasswordChanged          time.Time
	FirstName                string
	LastName                 string
	NickName                 string
	DisplayName              string
	AvatarKey                string
	PreferredLanguage        string
	Gender                   Gender
	Email                    string
	IsEmailVerified          bool
	VerifiedEmail            string
	Phone                    string
	IsPhoneVerified          bool
	Country                  string
	Locality                 string
	PostalCode               string
	Region                   string
	StreetAddress            string
	OTPState                 MFAState
	OTPSMSAdded              bool
	OTPEmailAdded            bool
	U2FTokens                []*WebAuthNView
	PasswordlessTokens       []*WebAuthNView
	MFAMaxSetUp              domain.MFALevel
	MFAInitSkipped           time.Time
	InitRequired             bool
	PasswordlessInitRequired bool
}

type WebAuthNView struct {
	TokenID string
	Name    string
	State   MFAState
}

type MachineView struct {
	LastKeyAdded time.Time
	Name         string
	Description  string
}

type UserSearchRequest struct {
	Offset        uint64
	Limit         uint64
	SortingColumn UserSearchKey
	Asc           bool
	Queries       []*UserSearchQuery
}

type UserSearchKey int32

const (
	UserSearchKeyUnspecified UserSearchKey = iota
	UserSearchKeyUserID
	UserSearchKeyUserName
	UserSearchKeyFirstName
	UserSearchKeyLastName
	UserSearchKeyNickName
	UserSearchKeyDisplayName
	UserSearchKeyEmail
	UserSearchKeyState
	UserSearchKeyResourceOwner
	UserSearchKeyLoginNames
	UserSearchKeyType
	UserSearchKeyPreferredLoginName
	UserSearchKeyInstanceID
	UserSearchOwnerRemoved
)

type UserSearchQuery struct {
	Key    UserSearchKey
	Method domain.SearchMethod
	Value  interface{}
}

type UserSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*UserView
	Sequence    uint64
	Timestamp   time.Time
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

func (r *UserSearchRequest) EnsureLimit(limit uint64) error {
	if r.Limit > limit {
		return zerrors.ThrowInvalidArgument(nil, "SEARCH-zz62F", "Errors.Limit.ExceedsDefault")
	}
	if r.Limit == 0 {
		r.Limit = limit
	}
	return nil
}

func (r *UserSearchRequest) AppendMyOrgQuery(orgID string) {
	r.Queries = append(r.Queries, &UserSearchQuery{Key: UserSearchKeyResourceOwner, Method: domain.SearchMethodEquals, Value: orgID})
}

func (u *UserView) MFATypesSetupPossible(level domain.MFALevel, policy *domain.LoginPolicy) []domain.MFAType {
	types := make([]domain.MFAType, 0)
	switch level {
	default:
		fallthrough
	case domain.MFALevelSecondFactor:
		if policy.HasSecondFactors() {
			for _, mfaType := range policy.SecondFactors {
				switch mfaType {
				case domain.SecondFactorTypeTOTP:
					if u.OTPState != MFAStateReady {
						types = append(types, domain.MFATypeTOTP)
					}
				case domain.SecondFactorTypeU2F:
					types = append(types, domain.MFATypeU2F)
				case domain.SecondFactorTypeOTPSMS:
					if !u.OTPSMSAdded {
						types = append(types, domain.MFATypeOTPSMS)
					}
				case domain.SecondFactorTypeOTPEmail:
					if !u.OTPEmailAdded {
						types = append(types, domain.MFATypeOTPEmail)
					}
				}
			}
		}
	}
	return types
}

func (u *UserView) MFATypesAllowed(level domain.MFALevel, policy *domain.LoginPolicy, isInternalAuthentication bool) ([]domain.MFAType, bool) {
	types := make([]domain.MFAType, 0)
	required := true
	switch level {
	default:
		required = domain.RequiresMFA(policy.ForceMFA, policy.ForceMFALocalOnly, isInternalAuthentication)
		fallthrough
	case domain.MFALevelSecondFactor:
		if policy.HasSecondFactors() {
			for _, mfaType := range policy.SecondFactors {
				switch mfaType {
				case domain.SecondFactorTypeTOTP:
					if u.OTPState == MFAStateReady {
						types = append(types, domain.MFATypeTOTP)
					}
				case domain.SecondFactorTypeU2F:
					if u.IsU2FReady() {
						types = append(types, domain.MFATypeU2F)
					}
				case domain.SecondFactorTypeOTPSMS:
					if u.OTPSMSAdded {
						types = append(types, domain.MFATypeOTPSMS)
					}
				case domain.SecondFactorTypeOTPEmail:
					if u.OTPEmailAdded {
						types = append(types, domain.MFATypeOTPEmail)
					}
				}
			}
		}
	}
	return types, required
}

func (u *UserView) IsU2FReady() bool {
	for _, token := range u.U2FTokens {
		if token.State == MFAStateReady {
			return true
		}
	}
	return false
}

func (u *UserView) IsPasswordlessReady() bool {
	for _, token := range u.PasswordlessTokens {
		if token.State == MFAStateReady {
			return true
		}
	}
	return false
}

func (u *UserView) GetProfile() (*Profile, error) {
	if u.HumanView == nil {
		return nil, zerrors.ThrowPreconditionFailed(nil, "MODEL-WLTce", "Errors.User.NotHuman")
	}
	return &Profile{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   u.ID,
			Sequence:      u.Sequence,
			ResourceOwner: u.ResourceOwner,
			CreationDate:  u.CreationDate,
			ChangeDate:    u.ChangeDate,
		},
		FirstName:          u.FirstName,
		LastName:           u.LastName,
		NickName:           u.NickName,
		DisplayName:        u.DisplayName,
		PreferredLanguage:  language.Make(u.PreferredLanguage),
		Gender:             u.Gender,
		PreferredLoginName: u.PreferredLoginName,
		LoginNames:         u.LoginNames,
		AvatarKey:          u.AvatarKey,
	}, nil
}

func (u *UserView) GetPhone() (*Phone, error) {
	if u.HumanView == nil {
		return nil, zerrors.ThrowPreconditionFailed(nil, "MODEL-him4a", "Errors.User.NotHuman")
	}
	return &Phone{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   u.ID,
			Sequence:      u.Sequence,
			ResourceOwner: u.ResourceOwner,
			CreationDate:  u.CreationDate,
			ChangeDate:    u.ChangeDate,
		},
		PhoneNumber:     u.Phone,
		IsPhoneVerified: u.IsPhoneVerified,
	}, nil
}

func (u *UserView) GetEmail() (*Email, error) {
	if u.HumanView == nil {
		return nil, zerrors.ThrowPreconditionFailed(nil, "MODEL-PWd6K", "Errors.User.NotHuman")
	}
	return &Email{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   u.ID,
			Sequence:      u.Sequence,
			ResourceOwner: u.ResourceOwner,
			CreationDate:  u.CreationDate,
			ChangeDate:    u.ChangeDate,
		},
		EmailAddress:    u.Email,
		IsEmailVerified: u.IsEmailVerified,
	}, nil
}

func (u *UserView) GetAddress() (*Address, error) {
	if u.HumanView == nil {
		return nil, zerrors.ThrowPreconditionFailed(nil, "MODEL-DN61m", "Errors.User.NotHuman")
	}
	return &Address{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   u.ID,
			Sequence:      u.Sequence,
			ResourceOwner: u.ResourceOwner,
			CreationDate:  u.CreationDate,
			ChangeDate:    u.ChangeDate,
		},
		Country:       u.Country,
		Locality:      u.Locality,
		PostalCode:    u.PostalCode,
		Region:        u.Region,
		StreetAddress: u.StreetAddress,
	}, nil
}
