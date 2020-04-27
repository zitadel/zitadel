package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetUserByID(ctx context.Context, id *UserID) (*User, error) {
	user, err := s.user.UserByID(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return userFromModel(user), nil
}

func (s *Server) GetUserByEmailGlobal(ctx context.Context, email *UserEmailID) (*User, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-9djSw", "Not implemented")
}

func (s *Server) SearchUsers(ctx context.Context, userSearch *UserSearchRequest) (*UserSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-as2Dc", "Not implemented")
}

func (s *Server) UserChanges(ctx context.Context, changesRequest *ChangeRequest) (*Changes, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-pl6Zu", "Not implemented")
}

func (s *Server) IsUserUnique(ctx context.Context, request *UniqueUserRequest) (*UniqueUserResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-olF56", "Not implemented")
}

func (s *Server) CreateUser(ctx context.Context, in *CreateUserRequest) (*User, error) {
	user, err := s.user.CreateUser(ctx, userCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return userFromModel(user), nil
}

func (s *Server) DeactivateUser(ctx context.Context, in *UserID) (*User, error) {
	user, err := s.user.DeactivateUser(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return userFromModel(user), nil
}

func (s *Server) ReactivateUser(ctx context.Context, in *UserID) (*User, error) {
	user, err := s.user.DeactivateUser(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return userFromModel(user), nil
}

func (s *Server) LockUser(ctx context.Context, in *UserID) (*User, error) {
	user, err := s.user.LockUser(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return userFromModel(user), nil
}

func (s *Server) UnlockUser(ctx context.Context, in *UserID) (*User, error) {
	user, err := s.user.UnlockUser(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return userFromModel(user), nil
}

func (s *Server) DeleteUser(ctx context.Context, ID *UserID) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-as4fg", "Not implemented")
}

func (s *Server) GetUserProfile(ctx context.Context, ID *UserID) (*UserProfile, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-mT67d", "Not implemented")
}

func (s *Server) UpdateUserProfile(ctx context.Context, request *UpdateUserProfileRequest) (*UserProfile, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-asje3", "Not implemented")
}

func (s *Server) GetUserEmail(ctx context.Context, ID *UserID) (*UserEmail, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-peo9d", "Not implemented")
}

func (s *Server) ChangeUserEmail(ctx context.Context, request *UpdateUserEmailRequest) (*UserEmail, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-cloeS", "Not implemented")
}

func (s *Server) ResendEmailVerificationMail(ctx context.Context, ID *UserID) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-dwsP9", "Not implemented")
}

func (s *Server) GetUserPhone(ctx context.Context, ID *UserID) (*UserPhone, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-wlf7f", "Not implemented")
}

func (s *Server) ChangeUserPhone(ctx context.Context, request *UpdateUserPhoneRequest) (*UserPhone, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-pld5g", "Not implemented")
}

func (s *Server) ResendPhoneVerificationCode(ctx context.Context, ID *UserID) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-98hdE", "Not implemented")
}

func (s *Server) GetUserAddress(ctx context.Context, ID *UserID) (*UserAddress, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-plt67", "Not implemented")
}

func (s *Server) UpdateUserAddress(ctx context.Context, request *UpdateUserAddressRequest) (*UserAddress, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-dleo3", "Not implemented")
}

func (s *Server) SendSetPasswordNotification(ctx context.Context, request *SetPasswordNotificationRequest) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-LSe7s", "Not implemented")
}

func (s *Server) SetInitialPassword(ctx context.Context, request *PasswordRequest) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-ldo3s", "Not implemented")
}

func (s *Server) GetUserMfas(ctx context.Context, userID *UserID) (*MultiFactors, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-ldmw3", "Not implemented")
}
