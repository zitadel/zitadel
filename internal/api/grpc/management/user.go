package management

import (
	"context"

	grpc_util "github.com/caos/zitadel/internal/api/grpc"
	"github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/pkg/management/grpc"

	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetUserByID(ctx context.Context, id *grpc.UserID) (*grpc.UserView, error) {
	user, err := s.user.UserByID(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return userViewFromModel(user), nil
}

func (s *Server) GetUserByEmailGlobal(ctx context.Context, email *grpc.Email) (*grpc.UserView, error) {
	user, err := s.user.GetGlobalUserByEmail(ctx, email.Email)
	if err != nil {
		return nil, err
	}
	return userViewFromModel(user), nil
}

func (s *Server) SearchUsers(ctx context.Context, in *grpc.UserSearchRequest) (*grpc.UserSearchResponse, error) {
	request := userSearchRequestsToModel(in)
	orgID := grpc_util.GetHeader(ctx, http.ZitadelOrgID)
	request.AppendMyOrgQuery(orgID)
	response, err := s.user.SearchUsers(ctx, request)
	if err != nil {
		return nil, err
	}
	return userSearchResponseFromModel(response), nil
}

func (s *Server) UserChanges(ctx context.Context, changesRequest *grpc.ChangeRequest) (*grpc.Changes, error) {
	response, err := s.user.UserChanges(ctx, changesRequest.Id, 0, 0)
	if err != nil {
		return nil, err
	}
	return userChangesToResponse(response, changesRequest.GetSequenceOffset(), changesRequest.GetLimit()), nil
}

func (s *Server) IsUserUnique(ctx context.Context, request *grpc.UniqueUserRequest) (*grpc.UniqueUserResponse, error) {
	unique, err := s.user.IsUserUnique(ctx, request.UserName, request.Email)
	if err != nil {
		return nil, err
	}
	return &grpc.UniqueUserResponse{IsUnique: unique}, nil
}

func (s *Server) CreateUser(ctx context.Context, in *grpc.CreateUserRequest) (*grpc.User, error) {
	user, err := s.user.CreateUser(ctx, userCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return userFromModel(user), nil
}

func (s *Server) DeactivateUser(ctx context.Context, in *grpc.UserID) (*grpc.User, error) {
	user, err := s.user.DeactivateUser(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return userFromModel(user), nil
}

func (s *Server) ReactivateUser(ctx context.Context, in *grpc.UserID) (*grpc.User, error) {
	user, err := s.user.ReactivateUser(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return userFromModel(user), nil
}

func (s *Server) LockUser(ctx context.Context, in *grpc.UserID) (*grpc.User, error) {
	user, err := s.user.LockUser(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return userFromModel(user), nil
}

func (s *Server) UnlockUser(ctx context.Context, in *grpc.UserID) (*grpc.User, error) {
	user, err := s.user.UnlockUser(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return userFromModel(user), nil
}

func (s *Server) DeleteUser(ctx context.Context, in *grpc.UserID) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-as4fg", "Not implemented")
}

func (s *Server) GetUserProfile(ctx context.Context, in *grpc.UserID) (*grpc.UserProfileView, error) {
	profile, err := s.user.ProfileByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return profileViewFromModel(profile), nil
}

func (s *Server) UpdateUserProfile(ctx context.Context, request *grpc.UpdateUserProfileRequest) (*grpc.UserProfile, error) {
	profile, err := s.user.ChangeProfile(ctx, updateProfileToModel(request))
	if err != nil {
		return nil, err
	}
	return profileFromModel(profile), nil
}

func (s *Server) GetUserEmail(ctx context.Context, in *grpc.UserID) (*grpc.UserEmailView, error) {
	email, err := s.user.EmailByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return emailViewFromModel(email), nil
}

func (s *Server) ChangeUserEmail(ctx context.Context, request *grpc.UpdateUserEmailRequest) (*grpc.UserEmail, error) {
	email, err := s.user.ChangeEmail(ctx, updateEmailToModel(request))
	if err != nil {
		return nil, err
	}
	return emailFromModel(email), nil
}

func (s *Server) ResendEmailVerificationMail(ctx context.Context, in *grpc.UserID) (*empty.Empty, error) {
	err := s.user.CreateEmailVerificationCode(ctx, in.Id)
	return &empty.Empty{}, err
}

func (s *Server) GetUserPhone(ctx context.Context, in *grpc.UserID) (*grpc.UserPhoneView, error) {
	phone, err := s.user.PhoneByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return phoneViewFromModel(phone), nil
}

func (s *Server) ChangeUserPhone(ctx context.Context, request *grpc.UpdateUserPhoneRequest) (*grpc.UserPhone, error) {
	phone, err := s.user.ChangePhone(ctx, updatePhoneToModel(request))
	if err != nil {
		return nil, err
	}
	return phoneFromModel(phone), nil
}

func (s *Server) ResendPhoneVerificationCode(ctx context.Context, in *grpc.UserID) (*empty.Empty, error) {
	err := s.user.CreatePhoneVerificationCode(ctx, in.Id)
	return &empty.Empty{}, err
}

func (s *Server) GetUserAddress(ctx context.Context, in *grpc.UserID) (*grpc.UserAddressView, error) {
	address, err := s.user.AddressByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return addressViewFromModel(address), nil
}

func (s *Server) UpdateUserAddress(ctx context.Context, request *grpc.UpdateUserAddressRequest) (*grpc.UserAddress, error) {
	address, err := s.user.ChangeAddress(ctx, updateAddressToModel(request))
	if err != nil {
		return nil, err
	}
	return addressFromModel(address), nil
}

func (s *Server) SendSetPasswordNotification(ctx context.Context, request *grpc.SetPasswordNotificationRequest) (*empty.Empty, error) {
	err := s.user.RequestSetPassword(ctx, request.Id, notifyTypeToModel(request.Type))
	return &empty.Empty{}, err
}

func (s *Server) SetInitialPassword(ctx context.Context, request *grpc.PasswordRequest) (*empty.Empty, error) {
	_, err := s.user.SetOneTimePassword(ctx, passwordRequestToModel(request))
	return &empty.Empty{}, err
}

func (s *Server) GetUserMfas(ctx context.Context, userID *grpc.UserID) (*grpc.MultiFactors, error) {
	mfas, err := s.user.UserMfas(ctx, userID.Id)
	if err != nil {
		return nil, err
	}
	return &grpc.MultiFactors{Mfas: mfasFromModel(mfas)}, nil
}
