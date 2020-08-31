package model

import (
	"time"

	req_model "github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/model"
	"golang.org/x/text/language"
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
	PasswordSet            bool
	PasswordChangeRequired bool
	UsernameChangeRequired bool
	PasswordChanged        time.Time
	FirstName              string
	LastName               string
	NickName               string
	DisplayName            string
	PreferredLanguage      string
	Gender                 Gender
	Email                  string
	IsEmailVerified        bool
	Phone                  string
	IsPhoneVerified        bool
	Country                string
	Locality               string
	PostalCode             string
	Region                 string
	StreetAddress          string
	OTPState               MfaState
	MfaMaxSetUp            req_model.MfaLevel
	MfaInitSkipped         time.Time
	InitRequired           bool
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
)

type UserSearchQuery struct {
	Key    UserSearchKey
	Method model.SearchMethod
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

func (r *UserSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}

func (r *UserSearchRequest) AppendMyOrgQuery(orgID string) {
	r.Queries = append(r.Queries, &UserSearchQuery{Key: UserSearchKeyResourceOwner, Method: model.SearchMethodEquals, Value: orgID})
}

func (u *UserView) MfaTypesSetupPossible(level req_model.MfaLevel) []req_model.MfaType {
	types := make([]req_model.MfaType, 0)
	switch level {
	default:
		fallthrough
	case req_model.MfaLevelSoftware:
		if u.OTPState != MfaStateReady {
			types = append(types, req_model.MfaTypeOTP)
		}
		//PLANNED: add sms
		fallthrough
	case req_model.MfaLevelHardware:
		//PLANNED: add token
	}
	return types
}

func (u *UserView) MfaTypesAllowed(level req_model.MfaLevel) []req_model.MfaType {
	types := make([]req_model.MfaType, 0)
	switch level {
	default:
		fallthrough
	case req_model.MfaLevelSoftware:
		if u.OTPState == MfaStateReady {
			types = append(types, req_model.MfaTypeOTP)
		}
		//PLANNED: add sms
		fallthrough
	case req_model.MfaLevelHardware:
		//PLANNED: add token
	}
	return types
}

func (u *UserView) GetProfile() (*Profile, error) {
	if u.HumanView == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "MODEL-WLTce", "Errors.User.NotHuman")
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
	}, nil
}

func (u *UserView) GetPhone() (*Phone, error) {
	if u.HumanView == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "MODEL-him4a", "Errors.User.NotHuman")
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
		return nil, errors.ThrowPreconditionFailed(nil, "MODEL-PWd6K", "Errors.User.NotHuman")
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
		return nil, errors.ThrowPreconditionFailed(nil, "MODEL-DN61m", "Errors.User.NotHuman")
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
