package auth

import (
	"github.com/caos/logging"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/grpc/auth"
	"github.com/golang/protobuf/ptypes"
)

func humanViewFromModel(user *usr_model.HumanView) *auth.HumanView {
	passwordChanged, err := ptypes.TimestampProto(user.PasswordChanged)
	logging.Log("MANAG-h4ByY").OnError(err).Debug("unable to parse date")

	return &auth.HumanView{
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		DisplayName:       user.DisplayName,
		NickName:          user.NickName,
		PreferredLanguage: user.PreferredLanguage,
		//TODO: add converter
		Gender:          auth.Gender(user.Gender),
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
		Phone:           user.Phone,
		IsPhoneVerified: user.IsPhoneVerified,
		Country:         user.Country,
		Locality:        user.Locality,
		PostalCode:      user.PostalCode,
		Region:          user.Region,
		StreetAddress:   user.StreetAddress,
		PasswordChanged: passwordChanged,
	}
}
