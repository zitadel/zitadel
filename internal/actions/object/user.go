package object

import (
	"time"

	"github.com/dop251/goja"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func UserFromExternalUser(c *actions.FieldConfig, user *domain.ExternalUser) goja.Value {
	return c.Runtime.ToValue(externalUserFromDomain(user))
}

func externalUsersFromDomain(users []*domain.ExternalUser) []*externalUser {
	externalUsers := make([]*externalUser, len(users))

	for i, user := range users {
		externalUsers[i] = externalUserFromDomain(user)
	}

	return externalUsers
}

func externalUserFromDomain(user *domain.ExternalUser) *externalUser {
	return &externalUser{
		ExternalId:    user.ExternalUserID,
		ExternalIdpId: user.IDPConfigID,
		Human: human{
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			NickName:          user.NickName,
			DisplayName:       user.DisplayName,
			PreferredLanguage: user.PreferredLanguage.String(),
			Email:             user.Email,
			IsEmailVerified:   user.IsEmailVerified,
			Phone:             user.Phone,
			IsPhoneVerified:   user.IsPhoneVerified,
		},
	}
}

func UserFromHuman(c *actions.FieldConfig, user *domain.Human) goja.Value {
	u := &humanUser{
		Id:                 user.AggregateID,
		CreationDate:       user.CreationDate,
		ChangeDate:         user.ChangeDate,
		ResourceOwner:      user.ResourceOwner,
		Sequence:           user.Sequence,
		State:              user.State,
		Username:           user.Username,
		LoginNames:         user.LoginNames,
		PreferredLoginName: user.PreferredLoginName,
	}

	if user.Profile != nil {
		u.Human.FirstName = user.Profile.FirstName
		u.Human.LastName = user.Profile.LastName
		u.Human.NickName = user.Profile.NickName
		u.Human.DisplayName = user.Profile.DisplayName
		u.Human.PreferredLanguage = user.Profile.PreferredLanguage.String()
	}

	if user.Email != nil {
		u.Human.Email = user.Email.EmailAddress
		u.Human.IsEmailVerified = user.Email.IsEmailVerified
	}

	if user.Phone != nil {
		u.Human.Phone = user.Phone.PhoneNumber
		u.Human.IsPhoneVerified = user.Phone.IsPhoneVerified
	}

	return c.Runtime.ToValue(u)
}

func UserFromQuery(c *actions.FieldConfig, user *query.User) goja.Value {
	if user.Human != nil {
		return humanFromQuery(c, user)
	}
	return machineFromQuery(c, user)
}

func humanFromQuery(c *actions.FieldConfig, user *query.User) goja.Value {
	return c.Runtime.ToValue(&humanUser{
		Id:                 user.ID,
		CreationDate:       user.CreationDate,
		ChangeDate:         user.ChangeDate,
		ResourceOwner:      user.ResourceOwner,
		Sequence:           user.Sequence,
		State:              user.State,
		Username:           user.Username,
		LoginNames:         user.LoginNames,
		PreferredLoginName: user.PreferredLoginName,
		Human: human{
			FirstName:         user.Human.FirstName,
			LastName:          user.Human.LastName,
			NickName:          user.Human.NickName,
			DisplayName:       user.Human.DisplayName,
			AvatarKey:         user.Human.AvatarKey,
			PreferredLanguage: user.Human.PreferredLanguage.String(),
			Gender:            user.Human.Gender,
			Email:             user.Human.Email,
			IsEmailVerified:   user.Human.IsEmailVerified,
			Phone:             user.Human.Phone,
			IsPhoneVerified:   user.Human.IsPhoneVerified,
		},
	})
}

func machineFromQuery(c *actions.FieldConfig, user *query.User) goja.Value {
	return c.Runtime.ToValue(&machineUser{
		Id:                 user.ID,
		CreationDate:       user.CreationDate,
		ChangeDate:         user.ChangeDate,
		ResourceOwner:      user.ResourceOwner,
		Sequence:           user.Sequence,
		State:              user.State,
		Username:           user.Username,
		LoginNames:         user.LoginNames,
		PreferredLoginName: user.PreferredLoginName,
		Machine: machine{
			Name:        user.Machine.Name,
			Description: user.Machine.Description,
		},
	})
}

type externalUser struct {
	ExternalId    string
	ExternalIdpId string
	Human         human
}

type humanUser struct {
	Id                 string
	CreationDate       time.Time
	ChangeDate         time.Time
	ResourceOwner      string
	Sequence           uint64
	State              domain.UserState
	Username           string
	LoginNames         database.StringArray
	PreferredLoginName string
	Human              human
}

type human struct {
	FirstName         string
	LastName          string
	NickName          string
	DisplayName       string
	AvatarKey         string
	PreferredLanguage string
	Gender            domain.Gender
	Email             domain.EmailAddress
	IsEmailVerified   bool
	Phone             domain.PhoneNumber
	IsPhoneVerified   bool
}

type machineUser struct {
	Id                 string
	CreationDate       time.Time
	ChangeDate         time.Time
	ResourceOwner      string
	Sequence           uint64
	State              domain.UserState
	Username           string
	LoginNames         database.StringArray
	PreferredLoginName string
	Machine            machine
}

type machine struct {
	Name        string
	Description string
}
