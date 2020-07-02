package grpc

import (
	"context"

	"github.com/caos/zitadel/internal/user/model"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetMyUser(ctx context.Context, _ *empty.Empty) (*UserView, error) {
	user, err := s.repo.MyUser(ctx)
	if err != nil {
		return nil, err
	}
	return userViewFromModel(user), nil
}

func (s *Server) GetMyUserProfile(ctx context.Context, _ *empty.Empty) (*UserProfileView, error) {
	profile, err := s.repo.MyProfile(ctx)
	if err != nil {
		return nil, err
	}
	return profileViewFromModel(profile), nil
}

func (s *Server) GetMyUserEmail(ctx context.Context, _ *empty.Empty) (*UserEmailView, error) {
	email, err := s.repo.MyEmail(ctx)
	if err != nil {
		return nil, err
	}
	return emailViewFromModel(email), nil
}

func (s *Server) GetMyUserPhone(ctx context.Context, _ *empty.Empty) (*UserPhoneView, error) {
	phone, err := s.repo.MyPhone(ctx)
	if err != nil {
		return nil, err
	}
	return phoneViewFromModel(phone), nil
}

func (s *Server) GetMyUserAddress(ctx context.Context, _ *empty.Empty) (*UserAddressView, error) {
	address, err := s.repo.MyAddress(ctx)
	if err != nil {
		return nil, err
	}
	return addressViewFromModel(address), nil
}

func (s *Server) GetMyMfas(ctx context.Context, _ *empty.Empty) (*MultiFactors, error) {
	mfas, err := s.repo.MyUserMfas(ctx)
	if err != nil {
		return nil, err
	}
	return &MultiFactors{Mfas: mfasFromModel(mfas)}, nil
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

func (s *Server) GetMyUserChanges(ctx context.Context, request *ChangesRequest) (*Changes, error) {
	changes, err := s.repo.MyUserChanges(ctx, request.SequenceOffset, request.Limit, request.Asc)
	if err != nil {
		return nil, err
	}
	changes = &model.UserChanges{
		Changes: []*model.UserChange{
			{
				ChangeDate: ptypes.TimestampNow(),
				EventType:  "user.added",
			},
		},
	}
	return userChangesToResponse(changes, request.GetSequenceOffset(), request.GetLimit()), nil
}
