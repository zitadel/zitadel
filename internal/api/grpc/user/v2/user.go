package user

import (
	"context"
	"io"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/idp"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func (s *Server) AddHumanUser(ctx context.Context, req *user.AddHumanUserRequest) (_ *user.AddHumanUserResponse, err error) {
	human, err := addUserRequestToAddHuman(req)
	if err != nil {
		return nil, err
	}
	orgID := authz.GetCtxData(ctx).OrgID
	err = s.command.AddHuman(ctx, orgID, human, false)
	if err != nil {
		return nil, err
	}
	return &user.AddHumanUserResponse{
		UserId:    human.ID,
		Details:   object.DomainToDetailsPb(human.Details),
		EmailCode: human.EmailCode,
	}, nil
}

func addUserRequestToAddHuman(req *user.AddHumanUserRequest) (*command.AddHuman, error) {
	username := req.GetUsername()
	if username == "" {
		username = req.GetEmail().GetEmail()
	}
	var urlTemplate string
	if req.GetEmail().GetSendCode() != nil {
		urlTemplate = req.GetEmail().GetSendCode().GetUrlTemplate()
		// test the template execution so the async notification will not fail because of it and the user won't realize
		if err := domain.RenderConfirmURLTemplate(io.Discard, urlTemplate, req.GetUserId(), "code", "orgID"); err != nil {
			return nil, err
		}
	}
	bcryptedPassword, err := hashedPasswordToCommand(req.GetHashedPassword())
	if err != nil {
		return nil, err
	}
	passwordChangeRequired := req.GetPassword().GetChangeRequired() || req.GetHashedPassword().GetChangeRequired()
	metadata := make([]*command.AddMetadataEntry, len(req.Metadata))
	for i, metadataEntry := range req.Metadata {
		metadata[i] = &command.AddMetadataEntry{
			Key:   metadataEntry.GetKey(),
			Value: metadataEntry.GetValue(),
		}
	}
	links := make([]*command.AddLink, len(req.Links))
	for i, link := range req.Links {
		links[i] = &command.AddLink{
			IDPID:         link.IdpId,
			IDPExternalID: link.IdpExternalId,
			DisplayName:   link.DisplayName,
		}
	}
	return &command.AddHuman{
		ID:          req.GetUserId(),
		Username:    username,
		FirstName:   req.GetProfile().GetFirstName(),
		LastName:    req.GetProfile().GetLastName(),
		NickName:    req.GetProfile().GetNickName(),
		DisplayName: req.GetProfile().GetDisplayName(),
		Email: command.Email{
			Address:     domain.EmailAddress(req.GetEmail().GetEmail()),
			Verified:    req.GetEmail().GetIsVerified(),
			ReturnCode:  req.GetEmail().GetReturnCode() != nil,
			URLTemplate: urlTemplate,
		},
		PreferredLanguage:      language.Make(req.GetProfile().GetPreferredLanguage()),
		Gender:                 genderToDomain(req.GetProfile().GetGender()),
		Phone:                  command.Phone{}, // TODO: add as soon as possible
		Password:               req.GetPassword().GetPassword(),
		BcryptedPassword:       bcryptedPassword,
		PasswordChangeRequired: passwordChangeRequired,
		Passwordless:           false,
		Register:               false,
		Metadata:               metadata,
		Links:                  links,
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

func hashedPasswordToCommand(hashed *user.HashedPassword) (string, error) {
	if hashed == nil {
		return "", nil
	}
	// we currently only handle bcrypt
	if hashed.GetAlgorithm() != "bcrypt" {
		return "", errors.ThrowInvalidArgument(nil, "USER-JDk4t", "Errors.InvalidArgument")
	}
	return hashed.GetHash(), nil
}

func (s *Server) AddLink(ctx context.Context, req *user.AddLinkRequest) (_ *user.AddLinkResponse, err error) {
	orgID := authz.GetCtxData(ctx).OrgID
	details, err := s.command.AddUserIDPLink(ctx, req.UserId, orgID, &domain.UserIDPLink{
		IDPConfigID:    req.Link.IdpId,
		ExternalUserID: req.Link.IdpExternalId,
		DisplayName:    "",
	})
	if err != nil {
		return nil, err
	}
	return &user.AddLinkResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) StartIdentityProviderFlow(ctx context.Context, req *user.StartIdentityProviderFlowRequest) (_ *user.StartIdentityProviderFlowResponse, err error) {


	createEventIntentPers(
		intentID,
		req.IdpId,
		intentToken = ""
		req.SuccessUrl,
		req.FailureUrl,
		userID = "",
		userInfo = {}
		)

	return &user.StartIdentityProviderFlowResponse{
		Details: object.DomainToDetailsPb(&domain.ObjectDetails{
			Sequence:      identityProvider.Sequence,
			EventDate:     identityProvider.ChangeDate,
			ResourceOwner: identityProvider.ResourceOwner,
		}),
		NextStep: &user.StartIdentityProviderFlowResponse_AuthUrl{AuthUrl: session.GetAuthURL()},
	}, nil
}
