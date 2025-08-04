package user

import (
	"context"
	"io"

	"connectrpc.com/connect"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	legacyobject "github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) createUserTypeHuman(ctx context.Context, humanPb *user.CreateUserRequest_Human, orgId string, userName, userId *string) (*connect.Response[user.CreateUserResponse], error) {
	metadataEntries := make([]*user.SetMetadataEntry, len(humanPb.Metadata))
	for i, metadataEntry := range humanPb.Metadata {
		metadataEntries[i] = &user.SetMetadataEntry{
			Key:   metadataEntry.GetKey(),
			Value: metadataEntry.GetValue(),
		}
	}
	addHumanPb := &user.AddHumanUserRequest{
		Username: userName,
		UserId:   userId,
		Organization: &legacyobject.Organization{
			Org: &legacyobject.Organization_OrgId{OrgId: orgId},
		},
		Profile:    humanPb.Profile,
		Email:      humanPb.Email,
		Phone:      humanPb.Phone,
		IdpLinks:   humanPb.IdpLinks,
		TotpSecret: humanPb.TotpSecret,
		Metadata:   metadataEntries,
	}
	switch pwType := humanPb.GetPasswordType().(type) {
	case *user.CreateUserRequest_Human_HashedPassword:
		addHumanPb.PasswordType = &user.AddHumanUserRequest_HashedPassword{
			HashedPassword: pwType.HashedPassword,
		}
	case *user.CreateUserRequest_Human_Password:
		addHumanPb.PasswordType = &user.AddHumanUserRequest_Password{
			Password: pwType.Password,
		}
	default:
		// optional password is not set
	}
	newHuman, err := AddUserRequestToAddHuman(addHumanPb)
	if err != nil {
		return nil, err
	}
	if err = s.command.AddUserHuman(
		ctx,
		orgId,
		newHuman,
		false,
		s.userCodeAlg,
	); err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.CreateUserResponse{
		Id:           newHuman.ID,
		CreationDate: timestamppb.New(newHuman.Details.EventDate),
		EmailCode:    newHuman.EmailCode,
		PhoneCode:    newHuman.PhoneCode,
	}), nil
}

func (s *Server) updateUserTypeHuman(ctx context.Context, humanPb *user.UpdateUserRequest_Human, userId string, userName *string) (*connect.Response[user.UpdateUserResponse], error) {
	cmd, err := updateHumanUserToCommand(userId, userName, humanPb)
	if err != nil {
		return nil, err
	}
	if err = s.command.ChangeUserHuman(ctx, cmd, s.userCodeAlg); err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.UpdateUserResponse{
		ChangeDate: timestamppb.New(cmd.Details.EventDate),
		EmailCode:  cmd.EmailCode,
		PhoneCode:  cmd.PhoneCode,
	}), nil
}

func updateHumanUserToCommand(userId string, userName *string, human *user.UpdateUserRequest_Human) (*command.ChangeHuman, error) {
	phone := human.GetPhone()
	if phone != nil && phone.Phone == "" && phone.GetVerification() != nil {
		return nil, zerrors.ThrowInvalidArgument(nil, "USERv2-4f3d6", "Errors.User.Phone.VerifyingRemovalIsNotSupported")
	}
	email, err := setHumanEmailToEmail(human.Email, userId)
	if err != nil {
		return nil, err
	}
	return &command.ChangeHuman{
		ID:       userId,
		Username: userName,
		Profile:  SetHumanProfileToProfile(human.Profile),
		Email:    email,
		Phone:    setHumanPhoneToPhone(human.Phone, true),
		Password: setHumanPasswordToPassword(human.Password),
	}, nil
}

func updateHumanUserRequestToChangeHuman(req *user.UpdateHumanUserRequest) (*command.ChangeHuman, error) {
	email, err := setHumanEmailToEmail(req.Email, req.GetUserId())
	if err != nil {
		return nil, err
	}
	changeHuman := &command.ChangeHuman{
		ID:       req.GetUserId(),
		Username: req.Username,
		Email:    email,
		Phone:    setHumanPhoneToPhone(req.Phone, false),
		Password: setHumanPasswordToPassword(req.Password),
	}
	if profile := req.GetProfile(); profile != nil {
		var firstName *string
		if profile.GivenName != "" {
			firstName = &profile.GivenName
		}
		var lastName *string
		if profile.FamilyName != "" {
			lastName = &profile.FamilyName
		}
		changeHuman.Profile = SetHumanProfileToProfile(&user.UpdateUserRequest_Human_Profile{
			GivenName:         firstName,
			FamilyName:        lastName,
			NickName:          profile.NickName,
			DisplayName:       profile.DisplayName,
			PreferredLanguage: profile.PreferredLanguage,
			Gender:            profile.Gender,
		})
	}
	return changeHuman, nil
}

func SetHumanProfileToProfile(profile *user.UpdateUserRequest_Human_Profile) *command.Profile {
	if profile == nil {
		return nil
	}
	return &command.Profile{
		FirstName:         profile.GivenName,
		LastName:          profile.FamilyName,
		NickName:          profile.NickName,
		DisplayName:       profile.DisplayName,
		PreferredLanguage: ifNotNilPtr(profile.PreferredLanguage, language.Make),
		Gender:            ifNotNilPtr(profile.Gender, genderToDomain),
	}
}

func setHumanEmailToEmail(email *user.SetHumanEmail, userID string) (*command.Email, error) {
	if email == nil {
		return nil, nil
	}
	var urlTemplate string
	if email.GetSendCode() != nil && email.GetSendCode().UrlTemplate != nil {
		urlTemplate = *email.GetSendCode().UrlTemplate
		if err := domain.RenderConfirmURLTemplate(io.Discard, urlTemplate, userID, "code", "orgID"); err != nil {
			return nil, err
		}
	}
	return &command.Email{
		Address:     domain.EmailAddress(email.Email),
		Verified:    email.GetIsVerified(),
		ReturnCode:  email.GetReturnCode() != nil,
		URLTemplate: urlTemplate,
	}, nil
}

func setHumanPhoneToPhone(phone *user.SetHumanPhone, withRemove bool) *command.Phone {
	if phone == nil {
		return nil
	}
	number := phone.GetPhone()
	return &command.Phone{
		Number:     domain.PhoneNumber(number),
		Verified:   phone.GetIsVerified(),
		ReturnCode: phone.GetReturnCode() != nil,
		Remove:     withRemove && number == "",
	}
}

func setHumanPasswordToPassword(password *user.SetPassword) *command.Password {
	if password == nil {
		return nil
	}
	return &command.Password{
		PasswordCode:        password.GetVerificationCode(),
		OldPassword:         password.GetCurrentPassword(),
		Password:            password.GetPassword().GetPassword(),
		EncodedPasswordHash: password.GetHashedPassword().GetHash(),
		ChangeRequired:      password.GetPassword().GetChangeRequired() || password.GetHashedPassword().GetChangeRequired(),
	}
}
