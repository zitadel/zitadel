package eventsourcing

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	"golang.org/x/text/language"
)

const (
	userVersion = "v1"
)

type User struct {
	es_models.ObjectRoot
	State int32 `json:"-"`
	*Profile
	*Email
	*Phone
	*Address
}

type Profile struct {
	es_models.ObjectRoot

	UserName          string       `json:"userName,omitempty"`
	FirstName         string       `json:"firstName,omitempty"`
	LastName          string       `json:"lastName,omitempty"`
	NickName          string       `json:"nickName,omitempty"`
	DisplayName       string       `json:"displayName,omitempty"`
	PreferredLanguage language.Tag `json:"preferredLanguage,omitempty"`
	Gender            int32        `json:"gender,omitempty"`
}

type Email struct {
	es_models.ObjectRoot

	Email           string `json:"email,omitempty"`
	IsEmailVerified bool   `json:"-"`
}

type Phone struct {
	es_models.ObjectRoot

	Phone           string `json:"phone,omitempty"`
	IsPhoneVerified bool   `json:"-"`
}

type Address struct {
	es_models.ObjectRoot

	Country       string `json:"country,omitempty"`
	Locality      string `json:"locality,omitempty"`
	PostalCode    string `json:"postalCode,omitempty"`
	Region        string `json:"region,omitempty"`
	StreetAddress string `json:"streetAddress,omitempty"`
}

func (p *Profile) Changes(changed *Profile) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.FirstName != "" && p.FirstName != changed.FirstName {
		changes["firstName"] = changed.FirstName
	}
	if changed.LastName != "" && p.LastName != changed.LastName {
		changes["lastName"] = changed.LastName
	}
	if changed.NickName != p.NickName {
		changes["nickName"] = changed.NickName
	}
	if changed.DisplayName != p.DisplayName {
		changes["displayName"] = changed.DisplayName
	}
	if p.PreferredLanguage != language.Und && changed.PreferredLanguage != p.PreferredLanguage {
		changes["preferredLanguage"] = changed.PreferredLanguage
	}
	if p.Gender > 0 && changed.Gender != p.Gender {
		changes["gender"] = changed.Gender
	}
	return changes
}

func (e *Email) Changes(changed *Email) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.Email != "" && e.Email != changed.Email {
		changes["email"] = changed.Email
	}
	return changes
}

func (p *Phone) Changes(changed *Phone) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.Phone != "" && p.Phone != changed.Phone {
		changes["phone"] = changed.Phone
	}
	return changes
}

func (a *Address) Changes(changed *Address) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if a.Country != changed.Country {
		changes["country"] = changed.Country
	}
	if a.Locality != changed.Locality {
		changes["locality"] = changed.Locality
	}
	if a.PostalCode != changed.PostalCode {
		changes["postalCode"] = changed.PostalCode
	}
	if a.Region != changed.Region {
		changes["region"] = changed.Region
	}
	if a.StreetAddress != changed.StreetAddress {
		changes["streetAddress"] = changed.StreetAddress
	}
	return changes
}

func UserFromModel(user *model.User) *User {
	return &User{
		ObjectRoot: es_models.ObjectRoot{
			ID:           user.ObjectRoot.ID,
			Sequence:     user.Sequence,
			ChangeDate:   user.ChangeDate,
			CreationDate: user.CreationDate,
		},
		Profile: ProfileFromModel(user.Profile),
	}
}

func UserToModel(user *User) *model.User {
	return &model.User{
		ObjectRoot: es_models.ObjectRoot{
			ID:           user.ObjectRoot.ID,
			Sequence:     user.Sequence,
			ChangeDate:   user.ChangeDate,
			CreationDate: user.CreationDate,
		},
		Profile: ProfileToModel(user.Profile),
	}
}

func ProfileFromModel(project *model.Profile) *Profile {
	return &Profile{
		ObjectRoot: es_models.ObjectRoot{
			ID:           project.ObjectRoot.ID,
			Sequence:     project.Sequence,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
		},
		UserName:          project.UserName,
		FirstName:         project.FirstName,
		LastName:          project.LastName,
		NickName:          project.NickName,
		DisplayName:       project.DisplayName,
		PreferredLanguage: project.PreferredLanguage,
		Gender:            int32(project.Gender),
	}
}

func ProfileToModel(project *Profile) *model.Profile {
	return &model.Profile{
		ObjectRoot: es_models.ObjectRoot{
			ID:           project.ObjectRoot.ID,
			Sequence:     project.Sequence,
			ChangeDate:   project.ChangeDate,
			CreationDate: project.CreationDate,
		},
		UserName:          project.UserName,
		FirstName:         project.FirstName,
		LastName:          project.LastName,
		NickName:          project.NickName,
		DisplayName:       project.DisplayName,
		PreferredLanguage: project.PreferredLanguage,
		Gender:            model.Gender(project.Gender),
	}
}

func EmailFromModel(email *model.Email) *Email {
	return &Email{
		ObjectRoot: es_models.ObjectRoot{
			ID:           email.ObjectRoot.ID,
			Sequence:     email.Sequence,
			ChangeDate:   email.ChangeDate,
			CreationDate: email.CreationDate,
		},
		Email:           email.Email,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func EmailToModel(email *Email) *model.Email {
	return &model.Email{
		ObjectRoot: es_models.ObjectRoot{
			ID:           email.ObjectRoot.ID,
			Sequence:     email.Sequence,
			ChangeDate:   email.ChangeDate,
			CreationDate: email.CreationDate,
		},
		Email:           email.Email,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func PhoneFromModel(phone *model.Phone) *Phone {
	return &Phone{
		ObjectRoot: es_models.ObjectRoot{
			ID:           phone.ObjectRoot.ID,
			Sequence:     phone.Sequence,
			ChangeDate:   phone.ChangeDate,
			CreationDate: phone.CreationDate,
		},
		Phone:           phone.Phone,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}

func PhoneToModel(phone *Phone) *model.Phone {
	return &model.Phone{
		ObjectRoot: es_models.ObjectRoot{
			ID:           phone.ObjectRoot.ID,
			Sequence:     phone.Sequence,
			ChangeDate:   phone.ChangeDate,
			CreationDate: phone.CreationDate,
		},
		Phone:           phone.Phone,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}

func AddressFromModel(address *model.Address) *Address {
	return &Address{
		ObjectRoot: es_models.ObjectRoot{
			ID:           address.ObjectRoot.ID,
			Sequence:     address.Sequence,
			ChangeDate:   address.ChangeDate,
			CreationDate: address.CreationDate,
		},
		Country:       address.Country,
		Locality:      address.Locality,
		PostalCode:    address.PostalCode,
		Region:        address.Region,
		StreetAddress: address.StreetAddress,
	}
}

func AddressToModel(address *Address) *model.Address {
	return &model.Address{
		ObjectRoot: es_models.ObjectRoot{
			ID:           address.ObjectRoot.ID,
			Sequence:     address.Sequence,
			ChangeDate:   address.ChangeDate,
			CreationDate: address.CreationDate,
		},
		Country:       address.Country,
		Locality:      address.Locality,
		PostalCode:    address.PostalCode,
		Region:        address.Region,
		StreetAddress: address.StreetAddress,
	}
}
