package user

import (
	"context"

	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	z_ctx "github.com/zitadel/zitadel/pkg/grpc/context/v2alpha"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func (s *Server) AddUser(ctx context.Context, req *user.AddUserRequest) (_ *user.AddUserResponse, err error) {
	human := addUserRequestToAddHuman(req)
	var details *domain.HumanDetails
	if req.UserId != nil {
		details, err = s.command.AddHumanWithID(ctx, authz.GetCtxData(ctx).OrgID, req.GetUserId(), human)
	} else {
		details, err = s.command.AddHuman(ctx, authz.GetCtxData(ctx).OrgID, human)
	}
	if err != nil {
		return nil, err
	}
	return &user.AddUserResponse{
		UserId: details.ID,
		Details: &z_ctx.ObjectDetails{
			Sequence:      details.Sequence,
			CreationDate:  timestamppb.New(details.EventDate),
			ChangeDate:    timestamppb.New(details.EventDate),
			ResourceOwner: details.ResourceOwner,
		},
	}, nil
}

func addUserRequestToAddHuman(req *user.AddUserRequest) *command.AddHuman {
	username := req.GetUsername()
	if username == "" {
		username = req.GetEmail().GetEmail()
	}
	return &command.AddHuman{
		Username:    username,
		FirstName:   req.GetProfile().GetFirstName(),
		LastName:    req.GetProfile().GetLastName(),
		NickName:    req.GetProfile().GetNickName(),
		DisplayName: req.GetProfile().GetDisplayName(),
		Email: command.Email{
			Address:  domain.EmailAddress(req.GetEmail().GetEmail()),
			Verified: req.GetEmail().GetIsVerified(), //TODO: oneof
		},
		PreferredLanguage:      language.Make(req.GetProfile().GetPreferredLanguage()),
		Gender:                 genderToDomain(req.GetProfile().GetGender()),
		Phone:                  command.Phone{}, // TODO: add
		Password:               req.GetPassword().GetPassword().GetPassword(),
		BcryptedPassword:       req.GetPassword().GetHashedPassword().GetHash(),
		PasswordChangeRequired: false,
		Passwordless:           false,
		ExternalIDP:            false,
		Register:               false,
		Metadata:
	}
}

func genderToDomain(gender user.Gender) domain.Gender {
	switch gender {
	case user.Gender_GENDER_UNSPECIFIED:
		return domain.GenderUnspecified
	case user.Gender_GENDER_FEMALE:
		return domain.GenderFemale
	case user.Gender_GENDER_MALE:
		return domain.GenderMale
	case user.Gender_GENDER_DIVERSE:
		return domain.GenderDiverse
	default:
		return domain.GenderUnspecified
	}
}
