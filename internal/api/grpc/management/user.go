package management

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetUserByID(ctx context.Context, id *management.UserID) (*management.UserView, error) {
	user, err := s.user.UserByID(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return userViewFromModel(user), nil
}

func (s *Server) GetUserByLoginNameGlobal(ctx context.Context, loginName *management.LoginName) (*management.UserView, error) {
	user, err := s.user.GetUserByLoginNameGlobal(ctx, loginName.LoginName)
	if err != nil {
		return nil, err
	}
	return userViewFromModel(user), nil
}

func (s *Server) SearchUsers(ctx context.Context, in *management.UserSearchRequest) (*management.UserSearchResponse, error) {
	request := userSearchRequestsToModel(in)
	request.AppendMyOrgQuery(authz.GetCtxData(ctx).OrgID)
	response, err := s.user.SearchUsers(ctx, request)
	if err != nil {
		return nil, err
	}
	return userSearchResponseFromModel(response), nil
}

func (s *Server) UserChanges(ctx context.Context, changesRequest *management.ChangeRequest) (*management.Changes, error) {
	response, err := s.user.UserChanges(ctx, changesRequest.Id, changesRequest.SequenceOffset, changesRequest.Limit, changesRequest.Asc)
	if err != nil {
		return nil, err
	}
	return userChangesToResponse(response, changesRequest.GetSequenceOffset(), changesRequest.GetLimit()), nil
}

func (s *Server) IsUserUnique(ctx context.Context, request *management.UniqueUserRequest) (*management.UniqueUserResponse, error) {
	unique, err := s.user.IsUserUnique(ctx, request.UserName, request.Email)
	if err != nil {
		return nil, err
	}
	return &management.UniqueUserResponse{IsUnique: unique}, nil
}

func (s *Server) CreateUser(ctx context.Context, in *management.CreateUserRequest) (*management.UserResponse, error) {
	human, machine := userCreateToDomain(in)
	if human != nil {
		h, err := s.command.AddHuman(ctx, authz.GetCtxData(ctx).OrgID, human)
		if err != nil {
			return nil, err
		}
		return userHumanFromDomain(h), nil
	}
	m, err := s.command.AddMachine(ctx, authz.GetCtxData(ctx).OrgID, machine)
	if err != nil {
		return nil, err
	}
	return userMachineFromDomain(m), nil
}

func (s *Server) DeactivateUser(ctx context.Context, in *management.UserID) (*empty.Empty, error) {
	err := s.command.DeactivateUser(ctx, in.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) ReactivateUser(ctx context.Context, in *management.UserID) (*empty.Empty, error) {
	err := s.command.ReactivateUser(ctx, in.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) LockUser(ctx context.Context, in *management.UserID) (*empty.Empty, error) {
	err := s.command.LockUser(ctx, in.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) UnlockUser(ctx context.Context, in *management.UserID) (*empty.Empty, error) {
	err := s.command.UnlockUser(ctx, in.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) DeleteUser(ctx context.Context, in *management.UserID) (*empty.Empty, error) {
	err := s.command.RemoveUser(ctx, in.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) UpdateUserMachine(ctx context.Context, in *management.UpdateMachineRequest) (*management.MachineResponse, error) {
	machine, err := s.command.ChangeMachine(ctx, updateMachineToDomain(authz.GetCtxData(ctx), in))
	if err != nil {
		return nil, err
	}
	return machineFromDomain(machine), nil
}

func (s *Server) GetUserProfile(ctx context.Context, in *management.UserID) (*management.UserProfileView, error) {
	profile, err := s.user.ProfileByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return profileViewFromModel(profile), nil
}

func (s *Server) ChangeUserUserName(ctx context.Context, request *management.UpdateUserUserNameRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.command.ChangeUsername(ctx, authz.GetCtxData(ctx).OrgID, request.Id, request.UserName)
}

func (s *Server) UpdateUserProfile(ctx context.Context, request *management.UpdateUserProfileRequest) (*management.UserProfile, error) {
	profile, err := s.command.ChangeHumanProfile(ctx, updateProfileToDomain(request))
	if err != nil {
		return nil, err
	}
	return profileFromDomain(profile), nil
}

func (s *Server) GetUserEmail(ctx context.Context, in *management.UserID) (*management.UserEmailView, error) {
	email, err := s.user.EmailByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return emailViewFromModel(email), nil
}

func (s *Server) ChangeUserEmail(ctx context.Context, request *management.UpdateUserEmailRequest) (*management.UserEmail, error) {
	email, err := s.command.ChangeHumanEmail(ctx, updateEmailToDomain(request))
	if err != nil {
		return nil, err
	}
	return emailFromDomain(email), nil
}

func (s *Server) ResendEmailVerificationMail(ctx context.Context, in *management.UserID) (*empty.Empty, error) {
	err := s.command.CreateHumanEmailVerificationCode(ctx, in.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) GetUserPhone(ctx context.Context, in *management.UserID) (*management.UserPhoneView, error) {
	phone, err := s.user.PhoneByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return phoneViewFromModel(phone), nil
}

func (s *Server) ChangeUserPhone(ctx context.Context, request *management.UpdateUserPhoneRequest) (*management.UserPhone, error) {
	phone, err := s.command.ChangeHumanPhone(ctx, updatePhoneToDomain(request))
	if err != nil {
		return nil, err
	}
	return phoneFromDomain(phone), nil
}

func (s *Server) RemoveUserPhone(ctx context.Context, userID *management.UserID) (*empty.Empty, error) {
	err := s.command.RemoveHumanPhone(ctx, userID.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) ResendPhoneVerificationCode(ctx context.Context, in *management.UserID) (*empty.Empty, error) {
	err := s.command.CreateHumanPhoneVerificationCode(ctx, in.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) GetUserAddress(ctx context.Context, in *management.UserID) (*management.UserAddressView, error) {
	address, err := s.user.AddressByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return addressViewFromModel(address), nil
}

func (s *Server) UpdateUserAddress(ctx context.Context, request *management.UpdateUserAddressRequest) (*management.UserAddress, error) {
	address, err := s.command.ChangeHumanAddress(ctx, updateAddressToDomain(authz.GetCtxData(ctx), request))
	if err != nil {
		return nil, err
	}
	return addressFromDomain(address), nil
}

func (s *Server) SendSetPasswordNotification(ctx context.Context, request *management.SetPasswordNotificationRequest) (*empty.Empty, error) {
	err := s.command.RequestSetPassword(ctx, request.Id, authz.GetCtxData(ctx).OrgID, notifyTypeToDomain(request.Type))
	return &empty.Empty{}, err
}

func (s *Server) SetInitialPassword(ctx context.Context, request *management.PasswordRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.command.SetOneTimePassword(ctx, authz.GetCtxData(ctx).OrgID, request.Id, request.Password)
}

func (s *Server) ResendInitialMail(ctx context.Context, request *management.InitialMailRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.command.ResendInitialMail(ctx, request.Id, request.Email, authz.GetCtxData(ctx).OrgID)
}

func (s *Server) SearchUserExternalIDPs(ctx context.Context, request *management.ExternalIDPSearchRequest) (*management.ExternalIDPSearchResponse, error) {
	externalIDP, err := s.user.SearchExternalIDPs(ctx, externalIDPSearchRequestToModel(request))
	if err != nil {
		return nil, err
	}
	return externalIDPSearchResponseFromModel(externalIDP), nil
}

func (s *Server) RemoveExternalIDP(ctx context.Context, request *management.ExternalIDPRemoveRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.command.RemoveHumanExternalIDP(ctx, externalIDPRemoveToDomain(authz.GetCtxData(ctx), request))
}

func (s *Server) GetUserMfas(ctx context.Context, userID *management.UserID) (*management.UserMultiFactors, error) {
	mfas, err := s.user.UserMFAs(ctx, userID.Id)
	if err != nil {
		return nil, err
	}
	return &management.UserMultiFactors{Mfas: mfasFromModel(mfas)}, nil
}

func (s *Server) RemoveMfaOTP(ctx context.Context, userID *management.UserID) (*empty.Empty, error) {
	return &empty.Empty{}, s.command.HumanRemoveOTP(ctx, userID.Id, authz.GetCtxData(ctx).OrgID)
}

func (s *Server) RemoveMfaU2F(ctx context.Context, webAuthNTokenID *management.WebAuthNTokenID) (*empty.Empty, error) {
	return &empty.Empty{}, s.command.HumanRemoveU2F(ctx, webAuthNTokenID.UserId, webAuthNTokenID.Id, authz.GetCtxData(ctx).OrgID)
}

func (s *Server) GetPasswordless(ctx context.Context, userID *management.UserID) (_ *management.WebAuthNTokens, err error) {
	tokens, err := s.user.GetPasswordless(ctx, userID.Id)
	if err != nil {
		return nil, err
	}
	return webAuthNTokensFromModel(tokens), err
}

func (s *Server) RemovePasswordless(ctx context.Context, id *management.WebAuthNTokenID) (*empty.Empty, error) {
	return &empty.Empty{}, s.command.HumanRemovePasswordless(ctx, id.UserId, id.Id, authz.GetCtxData(ctx).OrgID)
}

func (s *Server) SearchUserMemberships(ctx context.Context, in *management.UserMembershipSearchRequest) (*management.UserMembershipSearchResponse, error) {
	request := userMembershipSearchRequestsToModel(in)
	request.AppendUserIDQuery(in.UserId)
	response, err := s.user.SearchUserMemberships(ctx, request)
	if err != nil {
		return nil, err
	}
	return userMembershipSearchResponseFromModel(response), nil
}
