package grpc

import (
	"context"

	"github.com/caos/zitadel/internal/api"
	grpc_util "github.com/caos/zitadel/internal/api/grpc"
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

func (s *Server) GetUserByEmailGlobal(ctx context.Context, email *UserEmailID) (*UserView, error) {
	user, err := s.user.GetGlobalUserByEmail(ctx, email.Email)
	if err != nil {
		return nil, err
	}
	return userViewFromModel(user), nil
}

func (s *Server) SearchUsers(ctx context.Context, in *UserSearchRequest) (*UserSearchResponse, error) {
	request := userSearchRequestsToModel(in)
	orgID := grpc_util.GetHeader(ctx, api.ZitadelOrgID)
	request.AppendMyOrgQuery(orgID)
	response, err := s.user.SearchUsers(ctx, request)
	if err != nil {
		return nil, err
	}
	return userSearchResponseFromModel(response), nil
}

func (s *Server) UserChanges(ctx context.Context, changesRequest *ChangeRequest) (*Changes, error) {
	response, err := s.user.UserChanges(ctx, changesRequest.Id, 0, 0)
	if err != nil {
		return nil, err
	}
	return changesToResponse(response, changesRequest.GetSequenceOffset(), changesRequest.GetLimit()), nil
}

func (s *Server) IsUserUnique(ctx context.Context, request *UniqueUserRequest) (*UniqueUserResponse, error) {
	unique, err := s.user.IsUserUnique(ctx, request.UserName, request.Email)
	if err != nil {
		return nil, err
	}
	return &UniqueUserResponse{IsUnique: unique}, nil
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
	user, err := s.user.ReactivateUser(ctx, in.Id)
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

func (s *Server) DeleteUser(ctx context.Context, in *UserID) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-as4fg", "Not implemented")
}

func (s *Server) GetUserProfile(ctx context.Context, in *UserID) (*UserProfile, error) {
	profile, err := s.user.ProfileByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return profileFromModel(profile), nil
}

func (s *Server) UpdateUserProfile(ctx context.Context, request *UpdateUserProfileRequest) (*UserProfile, error) {
	profile, err := s.user.ChangeProfile(ctx, updateProfileToModel(request))
	if err != nil {
		return nil, err
	}
	return profileFromModel(profile), nil
}

func (s *Server) GetUserEmail(ctx context.Context, in *UserID) (*UserEmail, error) {
	email, err := s.user.EmailByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return emailFromModel(email), nil
}

func (s *Server) ChangeUserEmail(ctx context.Context, request *UpdateUserEmailRequest) (*UserEmail, error) {
	email, err := s.user.ChangeEmail(ctx, updateEmailToModel(request))
	if err != nil {
		return nil, err
	}
	return emailFromModel(email), nil
}

func (s *Server) ResendEmailVerificationMail(ctx context.Context, in *UserID) (*empty.Empty, error) {
	err := s.user.CreateEmailVerificationCode(ctx, in.Id)
	return &empty.Empty{}, err
}

func (s *Server) GetUserPhone(ctx context.Context, in *UserID) (*UserPhone, error) {
	phone, err := s.user.PhoneByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return phoneFromModel(phone), nil
}

func (s *Server) ChangeUserPhone(ctx context.Context, request *UpdateUserPhoneRequest) (*UserPhone, error) {
	phone, err := s.user.ChangePhone(ctx, updatePhoneToModel(request))
	if err != nil {
		return nil, err
	}
	return phoneFromModel(phone), nil
}

func (s *Server) ResendPhoneVerificationCode(ctx context.Context, in *UserID) (*empty.Empty, error) {
	err := s.user.CreatePhoneVerificationCode(ctx, in.Id)
	return &empty.Empty{}, err
}

func (s *Server) GetUserAddress(ctx context.Context, in *UserID) (*UserAddress, error) {
	address, err := s.user.AddressByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return addressFromModel(address), nil
}

func (s *Server) UpdateUserAddress(ctx context.Context, request *UpdateUserAddressRequest) (*UserAddress, error) {
	address, err := s.user.ChangeAddress(ctx, updateAddressToModel(request))
	if err != nil {
		return nil, err
	}
	return addressFromModel(address), nil
}

func (s *Server) SendSetPasswordNotification(ctx context.Context, request *SetPasswordNotificationRequest) (*empty.Empty, error) {
	err := s.user.RequestSetPassword(ctx, request.Id, notifyTypeToModel(request.Type))
	return &empty.Empty{}, err
}

func (s *Server) SetInitialPassword(ctx context.Context, request *PasswordRequest) (*empty.Empty, error) {
	_, err := s.user.SetOneTimePassword(ctx, passwordRequestToModel(request))
	return &empty.Empty{}, err
}

func (s *Server) GetUserMfas(ctx context.Context, userID *UserID) (*MultiFactors, error) {
	mfas, err := s.user.UserMfas(ctx, userID.Id)
	if err != nil {
		return nil, err
	}
	return &MultiFactors{Mfas: mfasFromModel(mfas)}, nil
}
