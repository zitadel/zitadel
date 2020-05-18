package model

import (
	"time"

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
	PasswordChanged        time.Time
	LastLogin              time.Time
	UserName               string
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
	USERSEARCHKEY_UNSPECIFIED UserSearchKey = iota
	USERSEARCHKEY_USER_ID
	USERSEARCHKEY_USER_NAME
	USERSEARCHKEY_FIRST_NAME
	USERSEARCHKEY_LAST_NAME
	USERSEARCHKEY_NICK_NAME
	USERSEARCHKEY_DISPLAY_NAME
	USERSEARCHKEY_EMAIL
	USERSEARCHKEY_STATE
	USERSEARCHKEY_RESOURCEOWNER
)

type UserSearchQuery struct {
	Key    UserSearchKey
	Method model.SearchMethod
	Value  string
}

type UserSearchResponse struct {
	Offset      uint64
	Limit       uint64
	TotalResult uint64
	Result      []*UserView
}

func (r *UserSearchRequest) EnsureLimit(limit uint64) {
	if r.Limit == 0 || r.Limit > limit {
		r.Limit = limit
	}
}

func (r *UserSearchRequest) AppendMyOrgQuery(orgID string) {
	r.Queries = append(r.Queries, &UserSearchQuery{Key: USERSEARCHKEY_RESOURCEOWNER, Method: model.SEARCHMETHOD_EQUALS, Value: orgID})
}

func (u *UserView) MfaTypesSetupPossible(level req_model.MfaLevel) []req_model.MfaType {
	types := make([]req_model.MfaType, 0)
	switch level {
	case req_model.MfaLevelSoftware:
		if u.OTPState != MFASTATE_READY {
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
		if u.OTPState == MFASTATE_READY {
			types = append(types, req_model.MfaTypeOTP)
		}
		//PLANNED: add sms
		fallthrough
	case req_model.MfaLevelHardware:
		//PLANNED: add token
	}
	return types
}
