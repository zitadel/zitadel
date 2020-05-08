package model

import (
	"context"
	"github.com/caos/zitadel/internal/api"
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/model"
	"time"
)

type UserView struct {
	ID                string
	CreationDate      time.Time
	ChangeDate        time.Time
	State             UserState
	ResourceOwner     string
	PasswordChanged   time.Time
	LastLogin         time.Time
	UserName          string
	FirstName         string
	LastName          string
	NickName          string
	DisplayName       string
	PreferredLanguage string
	Gender            Gender
	Email             string
	IsEmailVerified   bool
	Phone             string
	IsPhoneVerified   bool
	Country           string
	Locality          string
	PostalCode        string
	Region            string
	StreetAddress     string
	OTPState          MfaState
	Sequence          uint64
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

func (r *UserSearchRequest) AppendMyOrgQuery(ctx context.Context) {
	orgID := grpc_util.GetHeader(ctx, api.ZitadelOrgID)
	r.Queries = append(r.Queries, &UserSearchQuery{Key: USERSEARCHKEY_RESOURCEOWNER, Method: model.SEARCHMETHOD_EQUALS, Value: orgID})
}
