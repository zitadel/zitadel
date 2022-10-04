package model

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
)

type UserView struct {
	ID                 string
	UserName           string
	CreationDate       time.Time
	ChangeDate         time.Time
	State              UserState
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
	Phone                    string
	IsPhoneVerified          bool
	Country                  string
	Locality                 string
	PostalCode               string
	Region                   string
	StreetAddress            string
	OTPState                 MFAState
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
)

type UserSearchQuery struct {
	Key    UserSearchKey
	Method domain.SearchMethod
	Value  interface{}
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

func (u *UserView) MFATypesSetupPossible(level domain.MFALevel, policy *domain.LoginPolicy) []domain.MFAType {
	types := make([]domain.MFAType, 0)
	switch level {
	default:
		fallthrough
	case domain.MFALevelSecondFactor:
		if policy.HasSecondFactors() {
			for _, mfaType := range policy.SecondFactors {
				switch mfaType {
				case domain.SecondFactorTypeOTP:
					if u.OTPState != MFAStateReady {
						types = append(types, domain.MFATypeOTP)
					}
				case domain.SecondFactorTypeU2F:
					types = append(types, domain.MFATypeU2F)
				}
			}
		}
		//PLANNED: add sms
	}
	return types
}

func (u *UserView) MFATypesAllowed(level domain.MFALevel, policy *domain.LoginPolicy) ([]domain.MFAType, bool) {
	types := make([]domain.MFAType, 0)
	required := true
	switch level {
	default:
		required = policy.ForceMFA
		fallthrough
	case domain.MFALevelSecondFactor:
		if policy.HasSecondFactors() {
			for _, mfaType := range policy.SecondFactors {
				switch mfaType {
				case domain.SecondFactorTypeOTP:
					if u.OTPState == MFAStateReady {
						types = append(types, domain.MFATypeOTP)
					}
				case domain.SecondFactorTypeU2F:
					if u.IsU2FReady() {
						types = append(types, domain.MFATypeU2F)
					}
				}
			}
		}
		//PLANNED: add sms
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
