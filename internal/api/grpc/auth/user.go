package auth

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) GetMyUser(ctx context.Context, _ *empty.Empty) (*auth.UserView, error) {
	user, err := s.repo.MyUser(ctx)
	if err != nil {
		return nil, err
	}
	return userViewFromModel(user), nil
}

func (s *Server) GetMyUserProfile(ctx context.Context, _ *empty.Empty) (*auth.UserProfileView, error) {
	profile, err := s.repo.MyProfile(ctx)
	if err != nil {
		return nil, err
	}
	return profileViewFromModel(profile), nil
}

func (s *Server) GetMyUserEmail(ctx context.Context, _ *empty.Empty) (*auth.UserEmailView, error) {
	email, err := s.repo.MyEmail(ctx)
	if err != nil {
		return nil, err
	}
	return emailViewFromModel(email), nil
}

func (s *Server) GetMyUserPhone(ctx context.Context, _ *empty.Empty) (*auth.UserPhoneView, error) {
	phone, err := s.repo.MyPhone(ctx)
	if err != nil {
		return nil, err
	}
	return phoneViewFromModel(phone), nil
}

func (s *Server) RemoveMyUserPhone(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.repo.RemoveMyPhone(ctx)
	return &empty.Empty{}, err
}

func (s *Server) GetMyUserAddress(ctx context.Context, _ *empty.Empty) (*auth.UserAddressView, error) {
	address, err := s.repo.MyAddress(ctx)
	if err != nil {
		return nil, err
	}
	return addressViewFromModel(address), nil
}

func (s *Server) GetMyMfas(ctx context.Context, _ *empty.Empty) (*auth.MultiFactors, error) {
	mfas, err := s.repo.MyUserMFAs(ctx)
	if err != nil {
		return nil, err
	}
	return &auth.MultiFactors{Mfas: mfasFromModel(mfas)}, nil
}

func (s *Server) UpdateMyUserProfile(ctx context.Context, request *auth.UpdateUserProfileRequest) (*auth.UserProfile, error) {
	profile, err := s.command.ChangeHumanProfile(ctx, updateProfileToDomain(ctx, request))
	if err != nil {
		return nil, err
	}
	return profileFromDomain(profile), nil
}

func (s *Server) ChangeMyUserName(ctx context.Context, request *auth.ChangeUserNameRequest) (*empty.Empty, error) {
	ctxData := authz.GetCtxData(ctx)
	return &empty.Empty{}, s.command.ChangeUsername(ctx, ctxData.ResourceOwner, ctxData.UserID, request.UserName)
}

func (s *Server) ChangeMyUserEmail(ctx context.Context, request *auth.UpdateUserEmailRequest) (*auth.UserEmail, error) {
	email, err := s.command.ChangeHumanEmail(ctx, updateEmailToDomain(ctx, request))
	if err != nil {
		return nil, err
	}
	return emailFromDomain(email), nil
}

func (s *Server) VerifyMyUserEmail(ctx context.Context, request *auth.VerifyMyUserEmailRequest) (*empty.Empty, error) {
	err := s.repo.VerifyMyEmail(ctx, request.Code)
	return &empty.Empty{}, err
}

func (s *Server) ResendMyEmailVerificationMail(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.repo.ResendMyEmailVerificationMail(ctx)
	return &empty.Empty{}, err
}

func (s *Server) ChangeMyUserPhone(ctx context.Context, request *auth.UpdateUserPhoneRequest) (*auth.UserPhone, error) {
	phone, err := s.command.ChangeHumanPhone(ctx, updatePhoneToDomain(ctx, request))
	if err != nil {
		return nil, err
	}
	return phoneFromDomain(phone), nil
}

func (s *Server) VerifyMyUserPhone(ctx context.Context, request *auth.VerifyUserPhoneRequest) (*empty.Empty, error) {
	err := s.repo.VerifyMyPhone(ctx, request.Code)
	return &empty.Empty{}, err
}

func (s *Server) ResendMyPhoneVerificationCode(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.repo.ResendMyPhoneVerificationCode(ctx)
	return &empty.Empty{}, err
}

func (s *Server) UpdateMyUserAddress(ctx context.Context, request *auth.UpdateUserAddressRequest) (*auth.UserAddress, error) {
	address, err := s.repo.ChangeMyAddress(ctx, updateAddressToModel(ctx, request))
	if err != nil {
		return nil, err
	}
	return addressFromModel(address), nil
}

func (s *Server) ChangeMyPassword(ctx context.Context, request *auth.PasswordChange) (*empty.Empty, error) {
	err := s.repo.ChangeMyPassword(ctx, request.OldPassword, request.NewPassword)
	return &empty.Empty{}, err
}

func (s *Server) SearchMyExternalIDPs(ctx context.Context, request *auth.ExternalIDPSearchRequest) (*auth.ExternalIDPSearchResponse, error) {
	externalIDP, err := s.repo.SearchMyExternalIDPs(ctx, externalIDPSearchRequestToModel(request))
	if err != nil {
		return nil, err
	}
	return externalIDPSearchResponseFromModel(externalIDP), nil
}

func (s *Server) RemoveMyExternalIDP(ctx context.Context, request *auth.ExternalIDPRemoveRequest) (*empty.Empty, error) {
	err := s.repo.RemoveMyExternalIDP(ctx, externalIDPRemoveToModel(ctx, request))
	return &empty.Empty{}, err
}

func (s *Server) GetMyPasswordComplexityPolicy(ctx context.Context, _ *empty.Empty) (*auth.PasswordComplexityPolicy, error) {
	policy, err := s.repo.GetMyPasswordComplexityPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return passwordComplexityPolicyFromModel(policy), nil
}

func (s *Server) AddMfaOTP(ctx context.Context, _ *empty.Empty) (_ *auth.MfaOtpResponse, err error) {
	otp, err := s.repo.AddMyMFAOTP(ctx)
	if err != nil {
		return nil, err
	}
	return otpFromModel(otp), nil
}

func (s *Server) VerifyMfaOTP(ctx context.Context, request *auth.VerifyMfaOtp) (*empty.Empty, error) {
	err := s.repo.VerifyMyMFAOTPSetup(ctx, request.Code)
	return &empty.Empty{}, err
}

func (s *Server) RemoveMfaOTP(ctx context.Context, _ *empty.Empty) (_ *empty.Empty, err error) {
	err = s.repo.RemoveMyMFAOTP(ctx)
	return &empty.Empty{}, err
}

func (s *Server) AddMyMfaU2F(ctx context.Context, _ *empty.Empty) (_ *auth.WebAuthNResponse, err error) {
	u2f, err := s.repo.AddMyMFAU2F(ctx)
	if err != nil {
		return nil, err
	}
	return verifyWebAuthNFromModel(u2f), err
}

func (s *Server) VerifyMyMfaU2F(ctx context.Context, request *auth.VerifyWebAuthN) (*empty.Empty, error) {
	err := s.repo.VerifyMyMFAU2FSetup(ctx, request.TokenName, request.PublicKeyCredential)
	return &empty.Empty{}, err
}

func (s *Server) RemoveMyMfaU2F(ctx context.Context, id *auth.WebAuthNTokenID) (*empty.Empty, error) {
	err := s.repo.RemoveMyMFAU2F(ctx, id.Id)
	return &empty.Empty{}, err
}

func (s *Server) GetMyPasswordless(ctx context.Context, _ *empty.Empty) (_ *auth.WebAuthNTokens, err error) {
	tokens, err := s.repo.GetMyPasswordless(ctx)
	if err != nil {
		return nil, err
	}
	return webAuthNTokensFromModel(tokens), err
}

func (s *Server) AddMyPasswordless(ctx context.Context, _ *empty.Empty) (_ *auth.WebAuthNResponse, err error) {
	u2f, err := s.repo.AddMyPasswordless(ctx)
	if err != nil {
		return nil, err
	}
	return verifyWebAuthNFromModel(u2f), err
}

func (s *Server) VerifyMyPasswordless(ctx context.Context, request *auth.VerifyWebAuthN) (*empty.Empty, error) {
	err := s.repo.VerifyMyPasswordlessSetup(ctx, request.TokenName, request.PublicKeyCredential)
	return &empty.Empty{}, err
}

func (s *Server) RemoveMyPasswordless(ctx context.Context, id *auth.WebAuthNTokenID) (*empty.Empty, error) {
	err := s.repo.RemoveMyPasswordless(ctx, id.Id)
	return &empty.Empty{}, err
}

func (s *Server) GetMyUserChanges(ctx context.Context, request *auth.ChangesRequest) (*auth.Changes, error) {
	changes, err := s.repo.MyUserChanges(ctx, request.SequenceOffset, request.Limit, request.Asc)
	if err != nil {
		return nil, err
	}
	return userChangesToResponse(changes, request.GetSequenceOffset(), request.GetLimit()), nil
}
