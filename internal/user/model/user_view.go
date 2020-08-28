package model

import (
	"time"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/eventstore/models"

	req_model "github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/model"
)

type UserView struct {
	ID                     string
	CreationDate           time.Time
	ChangeDate             time.Time
	State                  UserState
	ResourceOwner          string
	PasswordSet            bool
	PasswordChangeRequired bool
	UsernameChangeRequired bool
	PasswordChanged        time.Time
	LastLogin              time.Time
	UserName               string
	PreferredLoginName     string
	LoginNames             []string
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
	Sequence               uint64
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

func (u *UserView) GetProfile() *Profile {
	return &Profile{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   u.ID,
			Sequence:      u.Sequence,
			ResourceOwner: u.ResourceOwner,
			CreationDate:  u.CreationDate,
			ChangeDate:    u.ChangeDate,
		},
		UserName:           u.UserName,
		FirstName:          u.FirstName,
		LastName:           u.LastName,
		NickName:           u.NickName,
		DisplayName:        u.DisplayName,
		PreferredLanguage:  language.Make(u.PreferredLanguage),
		Gender:             u.Gender,
		PreferredLoginName: u.PreferredLoginName,
		LoginNames:         u.LoginNames,
	}
}

func (u *UserView) GetPhone() *Phone {
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
	}
}

func (u *UserView) GetEmail() *Email {
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
	}
}

func (u *UserView) GetAddress() *Address {
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
	}
}
