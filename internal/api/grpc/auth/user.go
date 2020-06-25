package auth

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/auth/grpc"
)

func (s *Server) GetMyUser(ctx context.Context, _ *empty.Empty) (*grpc.UserView, error) {
	user, err := s.repo.MyUser(ctx)
	if err != nil {
		return nil, err
	}
	return userViewFromModel(user), nil
}

func (s *Server) GetMyUserProfile(ctx context.Context, _ *empty.Empty) (*grpc.UserProfileView, error) {
	profile, err := s.repo.MyProfile(ctx)
	if err != nil {
		return nil, err
	}
	return profileViewFromModel(profile), nil
}

func (s *Server) GetMyUserEmail(ctx context.Context, _ *empty.Empty) (*grpc.UserEmailView, error) {
	email, err := s.repo.MyEmail(ctx)
	if err != nil {
		return nil, err
	}
	return emailViewFromModel(email), nil
}

func (s *Server) GetMyUserPhone(ctx context.Context, _ *empty.Empty) (*grpc.UserPhoneView, error) {
	phone, err := s.repo.MyPhone(ctx)
	if err != nil {
		return nil, err
	}
	return phoneViewFromModel(phone), nil
}

func (s *Server) GetMyUserAddress(ctx context.Context, _ *empty.Empty) (*grpc.UserAddressView, error) {
	address, err := s.repo.MyAddress(ctx)
	if err != nil {
		return nil, err
	}
	return addressViewFromModel(address), nil
}

func (s *Server) GetMyMfas(ctx context.Context, _ *empty.Empty) (*grpc.MultiFactors, error) {
	mfas, err := s.repo.MyUserMfas(ctx)
	if err != nil {
		return nil, err
	}
	return &grpc.MultiFactors{Mfas: mfasFromModel(mfas)}, nil
}

func (s *Server) UpdateMyUserProfile(ctx context.Context, request *grpc.UpdateUserProfileRequest) (*grpc.UserProfile, error) {
	profile, err := s.repo.ChangeMyProfile(ctx, updateProfileToModel(ctx, request))
	if err != nil {
		return nil, err
	}
	return profileFromModel(profile), nil
}

func (s *Server) ChangeMyUserEmail(ctx context.Context, request *grpc.UpdateUserEmailRequest) (*grpc.UserEmail, error) {
	email, err := s.repo.ChangeMyEmail(ctx, updateEmailToModel(ctx, request))
	if err != nil {
		return nil, err
	}
	return emailFromModel(email), nil
}

func (s *Server) VerifyMyUserEmail(ctx context.Context, request *grpc.VerifyMyUserEmailRequest) (*empty.Empty, error) {
	err := s.repo.VerifyMyEmail(ctx, request.Code)
	return &empty.Empty{}, err
}

func (s *Server) ResendMyEmailVerificationMail(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.repo.ResendMyEmailVerificationMail(ctx)
	return &empty.Empty{}, err
}

func (s *Server) ChangeMyUserPhone(ctx context.Context, request *grpc.UpdateUserPhoneRequest) (*grpc.UserPhone, error) {
	phone, err := s.repo.ChangeMyPhone(ctx, updatePhoneToModel(ctx, request))
	if err != nil {
		return nil, err
	}
	return phoneFromModel(phone), nil
}

func (s *Server) VerifyMyUserPhone(ctx context.Context, request *grpc.VerifyUserPhoneRequest) (*empty.Empty, error) {
	err := s.repo.VerifyMyPhone(ctx, request.Code)
	return &empty.Empty{}, err
}

func (s *Server) ResendMyPhoneVerificationCode(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.repo.ResendMyPhoneVerificationCode(ctx)
	return &empty.Empty{}, err
}

func (s *Server) UpdateMyUserAddress(ctx context.Context, request *grpc.UpdateUserAddressRequest) (*grpc.UserAddress, error) {
	address, err := s.repo.ChangeMyAddress(ctx, updateAddressToModel(ctx, request))
	if err != nil {
		return nil, err
	}
	return addressFromModel(address), nil
}

func (s *Server) ChangeMyPassword(ctx context.Context, request *grpc.PasswordChange) (*empty.Empty, error) {
	err := s.repo.ChangeMyPassword(ctx, request.OldPassword, request.NewPassword)
	return &empty.Empty{}, err
}

func (s *Server) AddMfaOTP(ctx context.Context, _ *empty.Empty) (_ *grpc.MfaOtpResponse, err error) {
	otp, err := s.repo.AddMyMfaOTP(ctx)
	if err != nil {
		return nil, err
	}
	return otpFromModel(otp), nil
}

func (s *Server) VerifyMfaOTP(ctx context.Context, request *grpc.VerifyMfaOtp) (*empty.Empty, error) {
	err := s.repo.VerifyMyMfaOTPSetup(ctx, request.Code)
	return &empty.Empty{}, err
}

func (s *Server) RemoveMfaOTP(ctx context.Context, _ *empty.Empty) (_ *empty.Empty, err error) {
	s.repo.RemoveMyMfaOTP(ctx)
	return &empty.Empty{}, err
}
