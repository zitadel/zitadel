package user

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func (s *Server) AddUser(ctx context.Context, req *user.AddUserRequest) (_ *user.AddUserResponse, err error) {
	human, err := addUserRequestToAddHuman(req)
	if err != nil {
		return nil, err
	}
	err = s.command.AddHuman(ctx, req.GetOrganisation().GetOrgId(), human, false)
	if err != nil {
		return nil, err
	}
	var emailCode *string
	if req.GetEmail().GetReturnCode() != nil {
		emailCode = &human.Email.VerificationCode
	}
	return &user.AddUserResponse{
		UserId:    human.ID,
		Details:   object.DomainToDetailsPb(human.Details),
		EmailCode: emailCode,
	}, nil
}

func addUserRequestToAddHuman(req *user.AddUserRequest) (*command.AddHuman, error) {
	username := req.GetUsername()
	if username == "" {
		username = req.GetEmail().GetEmail()
	}
	metadata := make([]*command.AddMetadataEntry, len(req.Metadata))
	for i, metadataEntry := range req.Metadata {
		metadata[i] = &command.AddMetadataEntry{
			Key:   metadataEntry.GetKey(),
			Value: metadataEntry.GetValue(),
		}
	}
	bcryptedPassword, err := setPasswordHashedPasswordToCommand(req.GetHashedPassword())
	if err != nil {
		return nil, err
	}
	passwordChangeRequired := req.GetPassword().GetChangeRequired() || req.GetHashedPassword().GetChangeRequired()
	return &command.AddHuman{
		ID:          req.GetUserId(),
		Username:    username,
		FirstName:   req.GetProfile().GetFirstName(),
		LastName:    req.GetProfile().GetLastName(),
		NickName:    req.GetProfile().GetNickName(),
		DisplayName: req.GetProfile().GetDisplayName(),
		Email: command.Email{
			Address:    domain.EmailAddress(req.GetEmail().GetEmail()),
			Verified:   req.GetEmail().GetIsVerified(), //TODO: oneof
			ReturnCode: req.GetEmail().GetReturnCode() != nil,
		},
		PreferredLanguage:      language.Make(req.GetProfile().GetPreferredLanguage()),
		Gender:                 genderToDomain(req.GetProfile().GetGender()),
		Phone:                  command.Phone{}, // TODO: add as soon as possible
		Password:               req.GetPassword().GetPassword(),
		BcryptedPassword:       bcryptedPassword,
		PasswordChangeRequired: passwordChangeRequired,
		Passwordless:           false,
		ExternalIDP:            false,
		Register:               false,
		Metadata:               metadata,
	}, nil
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

func setPasswordHashedPasswordToCommand(hashed *user.HashedPassword) (string, error) {
	if hashed == nil {
		return "", nil
	}
	// we currently only handle bcrypt
	if hashed.GetAlgorithm() != "bcrypt" {
		return "", errors.ThrowInvalidArgument(nil, "USER-JDk4t", "Errors.InvalidArgument")
	}
	return hashed.GetHash(), nil
}
