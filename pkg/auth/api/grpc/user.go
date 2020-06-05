package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/internal/errors"
)

func (s *Server) GetMyUserProfile(ctx context.Context, _ *empty.Empty) (*UserProfile, error) {
	profile, err := s.repo.MyProfile(ctx)
	if err != nil {
		return nil, err
	}
	return profileFromModel(profile), nil
}

func (s *Server) GetMyUserEmail(ctx context.Context, _ *empty.Empty) (*UserEmail, error) {
	email, err := s.repo.MyEmail(ctx)
	if err != nil {
		return nil, err
	}
	return emailFromModel(email), nil
}

func (s *Server) GetMyUserPhone(ctx context.Context, _ *empty.Empty) (*UserPhone, error) {
	phone, err := s.repo.MyPhone(ctx)
	if err != nil {
		return nil, err
	}
	return phoneFromModel(phone), nil
}

func (s *Server) GetMyUserAddress(ctx context.Context, _ *empty.Empty) (*UserAddress, error) {
	address, err := s.repo.MyAddress(ctx)
	if err != nil {
		return nil, err
	}
	return addressFromModel(address), nil
}

func (s *Server) GetMyMfas(ctx context.Context, _ *empty.Empty) (*MultiFactors, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-vkl9i", "Not implemented")
}

func (s *Server) UpdateMyUserProfile(ctx context.Context, request *UpdateUserProfileRequest) (*UserProfile, error) {
	profile, err := s.repo.ChangeMyProfile(ctx, updateProfileToModel(ctx, request))
	if err != nil {
		return nil, err
	}
	return profileFromModel(profile), nil
}

func (s *Server) ChangeMyUserEmail(ctx context.Context, request *UpdateUserEmailRequest) (*UserEmail, error) {
	email, err := s.repo.ChangeMyEmail(ctx, updateEmailToModel(ctx, request))
	if err != nil {
		return nil, err
	}
	return emailFromModel(email), nil
}

func (s *Server) VerifyMyUserEmail(ctx context.Context, request *VerifyMyUserEmailRequest) (*empty.Empty, error) {
	err := s.repo.VerifyMyEmail(ctx, request.Code)
	return &empty.Empty{}, err
}

func (s *Server) ResendMyEmailVerificationMail(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.repo.ResendMyEmailVerificationMail(ctx)
	return &empty.Empty{}, err
}

func (s *Server) ChangeMyUserPhone(ctx context.Context, request *UpdateUserPhoneRequest) (*UserPhone, error) {
	phone, err := s.repo.ChangeMyPhone(ctx, updatePhoneToModel(ctx, request))
	if err != nil {
		return nil, err
	}
	return phoneFromModel(phone), nil
}

func (s *Server) VerifyMyUserPhone(ctx context.Context, request *VerifyUserPhoneRequest) (*empty.Empty, error) {
	err := s.repo.VerifyMyPhone(ctx, request.Code)
	return &empty.Empty{}, err
}

func (s *Server) ResendMyPhoneVerificationCode(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.repo.ResendMyPhoneVerificationCode(ctx)
	return &empty.Empty{}, err
}

func (s *Server) UpdateMyUserAddress(ctx context.Context, request *UpdateUserAddressRequest) (*UserAddress, error) {
	address, err := s.repo.ChangeMyAddress(ctx, updateAddressToModel(ctx, request))
	if err != nil {
		return nil, err
	}
	return addressFromModel(address), nil
}

func (s *Server) ChangeMyPassword(ctx context.Context, request *PasswordChange) (*empty.Empty, error) {
	err := s.repo.ChangeMyPassword(ctx, request.OldPassword, request.NewPassword)
	return &empty.Empty{}, err
}

func (s *Server) AddMfaOTP(ctx context.Context, _ *empty.Empty) (_ *MfaOtpResponse, err error) {
	otp, err := s.repo.AddMyMfaOTP(ctx)
	if err != nil {
		return nil, err
	}
	return otpFromModel(otp), nil
}

func (s *Server) VerifyMfaOTP(ctx context.Context, request *VerifyMfaOtp) (*empty.Empty, error) {
	err := s.repo.VerifyMyMfaOTPSetup(ctx, request.Code)
	return &empty.Empty{}, err
}

func (s *Server) RemoveMfaOTP(ctx context.Context, _ *empty.Empty) (_ *empty.Empty, err error) {
	s.repo.RemoveMyMfaOTP(ctx)
	return &empty.Empty{}, err
}
