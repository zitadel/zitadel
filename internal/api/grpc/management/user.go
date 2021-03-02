package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/authn"
	change_grpc "github.com/caos/zitadel/internal/api/grpc/change"
	idp_grpc "github.com/caos/zitadel/internal/api/grpc/idp"
	"github.com/caos/zitadel/internal/api/grpc/object"
	obj_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/api/grpc/user"
	user_grpc "github.com/caos/zitadel/internal/api/grpc/user"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetUserByID(ctx context.Context, req *mgmt_pb.GetUserByIDRequest) (*mgmt_pb.GetUserByIDResponse, error) {
	user, err := s.user.UserByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetUserByIDResponse{
		User: user_grpc.UserToPb(user),
	}, nil
}

func (s *Server) GetUserByLoginNameGlobal(ctx context.Context, req *mgmt_pb.GetUserByLoginNameGlobalRequest) (*mgmt_pb.GetUserByLoginNameGlobalResponse, error) {
	user, err := s.user.GetUserByLoginNameGlobal(ctx, req.LoginName)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetUserByLoginNameGlobalResponse{
		User: user_grpc.UserToPb(user),
	}, nil
}

func (s *Server) ListUsers(ctx context.Context, req *mgmt_pb.ListUsersRequest) (*mgmt_pb.ListUsersResponse, error) {
	r := ListUsersRequestToModel(ctx, req)
	res, err := s.user.SearchUsers(ctx, r)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListUsersResponse{
		Result: user_grpc.UsersToPb(res.Result),
		MetaData: obj_grpc.ToListDetails(
			res.TotalResult,
			res.Sequence,
			res.Timestamp,
		),
	}, nil
}

func (s *Server) ListUserChanges(ctx context.Context, req *mgmt_pb.ListUserChangesRequest) (*mgmt_pb.ListUserChangesResponse, error) {
	res, err := s.user.UserChanges(ctx, req.UserId, req.Query.Offset, uint64(req.Query.Limit), req.Query.Asc)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListUserChangesResponse{
		Result: change_grpc.UserChangesToPb(res.Changes),
	}, nil
}

func (s *Server) IsUserUnique(ctx context.Context, req *mgmt_pb.IsUserUniqueRequest) (*mgmt_pb.IsUserUniqueResponse, error) {
	unique, err := s.user.IsUserUnique(ctx, req.UserName, req.Email)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.IsUserUniqueResponse{
		IsUnique: unique,
	}, nil
}

func (s *Server) AddHumanUser(ctx context.Context, req *mgmt_pb.AddHumanUserRequest) (*mgmt_pb.AddHumanUserResponse, error) {
	human, err := s.command.AddHuman(ctx, authz.GetCtxData(ctx).OrgID, AddHumanUserRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddHumanUserResponse{
		UserId: human.AggregateID,
		Details: obj_grpc.ToDetailsPb(
			human.Sequence,
			human.ChangeDate,
			human.ResourceOwner,
		),
	}, nil
}

func (s *Server) AddMachineUser(ctx context.Context, req *mgmt_pb.AddMachineUserRequest) (*mgmt_pb.AddMachineUserResponse, error) {
	machine, err := s.command.AddMachine(ctx, authz.GetCtxData(ctx).OrgID, AddMachineUserRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddMachineUserResponse{
		UserId: machine.AggregateID,
		Details: obj_grpc.ToDetailsPb(
			machine.Sequence,
			machine.ChangeDate,
			machine.ResourceOwner,
		),
	}, nil
}

func (s *Server) DeactivateUser(ctx context.Context, req *mgmt_pb.DeactivateUserRequest) (*mgmt_pb.DeactivateUserResponse, error) {
	err := s.command.DeactivateUser(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeactivateUserResponse{
		//TODO: details
	}, nil
}

func (s *Server) ReactivateUser(ctx context.Context, req *mgmt_pb.ReactivateUserRequest) (*mgmt_pb.ReactivateUserResponse, error) {
	err := s.command.ReactivateUser(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateUserResponse{
		//TODO: details
	}, nil
}

func (s *Server) LockUser(ctx context.Context, req *mgmt_pb.LockUserRequest) (*mgmt_pb.LockUserResponse, error) {
	err := s.command.LockUser(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.LockUserResponse{
		//TODO: details
	}, nil
}

func (s *Server) UnlockUser(ctx context.Context, req *mgmt_pb.UnlockUserRequest) (*mgmt_pb.UnlockUserResponse, error) {
	err := s.command.UnlockUser(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UnlockUserResponse{
		//TODO: details
	}, nil
}

func (s *Server) RemoveUser(ctx context.Context, req *mgmt_pb.RemoveUserRequest) (*mgmt_pb.RemoveUserResponse, error) {
	grants, err := s.usergrant.UserGrantsByUserID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	err = s.command.RemoveUser(ctx, req.Id, authz.GetCtxData(ctx).OrgID, userGrantsToIDs(grants)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveUserResponse{
		//TODO: details
	}, nil
}

func userGrantsToIDs(userGrants []*grant_model.UserGrantView) []string {
	converted := make([]string, len(userGrants))
	for i, grant := range userGrants {
		converted[i] = grant.ID
	}
	return converted
}

func (s *Server) UpdateUserName(ctx context.Context, req *mgmt_pb.UpdateUserNameRequest) (*mgmt_pb.UpdateUserNameResponse, error) {
	err := s.command.ChangeUsername(ctx, authz.GetCtxData(ctx).OrgID, req.UserId, req.UserName)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateUserNameResponse{
		//TODO: details
	}, nil
}

func (s *Server) GetHumanProfile(ctx context.Context, req *mgmt_pb.GetHumanProfileRequest) (*mgmt_pb.GetHumanProfileResponse, error) {
	profile, err := s.user.ProfileByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetHumanProfileResponse{
		Profile: user_grpc.ProfileToPb(profile),
		Details: obj_grpc.ToDetailsPb(
			profile.Sequence,
			profile.ChangeDate,
			profile.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateHumanProfile(ctx context.Context, req *mgmt_pb.UpdateHumanProfileRequest) (*mgmt_pb.UpdateHumanProfileResponse, error) {
	profile, err := s.command.ChangeHumanProfile(ctx, UpdateHumanProfileRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateHumanProfileResponse{
		Details: obj_grpc.ToDetailsPb(
			profile.Sequence,
			profile.ChangeDate,
			profile.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetHumanEmail(ctx context.Context, req *mgmt_pb.GetHumanEmailRequest) (*mgmt_pb.GetHumanEmailResponse, error) {
	email, err := s.user.EmailByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetHumanEmailResponse{
		Email: user_grpc.EmailToPb(email),
		Details: obj_grpc.ToDetailsPb(
			email.Sequence,
			email.ChangeDate,
			email.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateHumanEmail(ctx context.Context, req *mgmt_pb.UpdateHumanEmailRequest) (*mgmt_pb.UpdateHumanEmailResponse, error) {
	email, err := s.command.ChangeHumanEmail(ctx, UpdateHumanEmailRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateHumanEmailResponse{
		Details: obj_grpc.ToDetailsPb(
			email.Sequence,
			email.ChangeDate,
			email.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResendHumanInitialization(ctx context.Context, req *mgmt_pb.ResendHumanInitializationRequest) (*mgmt_pb.ResendHumanInitializationResponse, error) {
	//TODO: why do we need the email again?
	err := s.command.ResendInitialMail(ctx, req.UserId, "email", authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResendHumanInitializationResponse{
		//TODO: details
	}, nil
}

func (s *Server) ResendHumanEmailVerification(ctx context.Context, req *mgmt_pb.ResendHumanEmailVerificationRequest) (*mgmt_pb.ResendHumanEmailVerificationResponse, error) {
	details, err := s.command.CreateHumanEmailVerificationCode(ctx, req.UserId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResendHumanEmailVerificationResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) GetHumanPhone(ctx context.Context, req *mgmt_pb.GetHumanPhoneRequest) (*mgmt_pb.GetHumanPhoneResponse, error) {
	phone, err := s.user.PhoneByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetHumanPhoneResponse{
		Phone: user_grpc.PhoneToPb(phone),
		Details: obj_grpc.ToDetailsPb(
			phone.Sequence,
			phone.ChangeDate,
			phone.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateHumanPhone(ctx context.Context, req *mgmt_pb.UpdateHumanPhoneRequest) (*mgmt_pb.UpdateHumanPhoneResponse, error) {
	phone, err := s.command.ChangeHumanPhone(ctx, UpdateHumanPhoneRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateHumanPhoneResponse{
		Details: obj_grpc.ToDetailsPb(
			phone.Sequence,
			phone.ChangeDate,
			phone.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveHumanPhone(ctx context.Context, req *mgmt_pb.RemoveHumanPhoneRequest) (*mgmt_pb.RemoveHumanPhoneResponse, error) {
	err := s.command.RemoveHumanPhone(ctx, req.UserId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveHumanPhoneResponse{
		//TODO: details
	}, nil
}

func (s *Server) ResendHumanPhoneVerification(ctx context.Context, req *mgmt_pb.ResendHumanPhoneVerificationRequest) (*mgmt_pb.ResendHumanPhoneVerificationResponse, error) {
	err := s.command.CreateHumanPhoneVerificationCode(ctx, req.UserId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResendHumanPhoneVerificationResponse{
		//TODO: details
	}, nil
}

func (s *Server) SetHumanInitialPassword(ctx context.Context, req *mgmt_pb.SetHumanInitialPasswordRequest) (*mgmt_pb.SetHumanInitialPasswordResponse, error) {
	err := s.command.SetOneTimePassword(ctx, authz.GetCtxData(ctx).OrgID, req.UserId, req.Password)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetHumanInitialPasswordResponse{
		//TODO: details
	}, nil
}

func (s *Server) SendHumanResetPasswordNotification(ctx context.Context, req *mgmt_pb.SendHumanResetPasswordNotificationRequest) (*mgmt_pb.SendHumanResetPasswordNotificationResponse, error) {
	err := s.command.RequestSetPassword(ctx, req.UserId, authz.GetCtxData(ctx).OrgID, notifyTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SendHumanResetPasswordNotificationResponse{
		// TODO: details
	}, nil
}

func (s *Server) ListHumanMultiFactors(ctx context.Context, req *mgmt_pb.ListHumanMultiFactorsRequest) (*mgmt_pb.ListHumanMultiFactorsResponse, error) {
	mfas, err := s.user.UserMFAs(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListHumanMultiFactorsResponse{
		Result: user_grpc.MultiFactorsToPb(mfas),
	}, nil
}

func (s *Server) RemoveHumanMultiFactorOTP(ctx context.Context, req *mgmt_pb.RemoveHumanMultiFactorOTPRequest) (*mgmt_pb.RemoveHumanMultiFactorOTPResponse, error) {
	details, err := s.command.HumanRemoveOTP(ctx, req.UserId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveHumanMultiFactorOTPResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) RemoveHumanMultiFactorU2F(ctx context.Context, req *mgmt_pb.RemoveHumanMultiFactorU2FRequest) (*mgmt_pb.RemoveHumanMultiFactorU2FResponse, error) {
	details, err := s.command.HumanRemoveU2F(ctx, req.UserId, req.TokenId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveHumanMultiFactorU2FResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) ListHumanPasswordless(ctx context.Context, req *mgmt_pb.ListHumanPasswordlessRequest) (*mgmt_pb.ListHumanPasswordlessResponse, error) {
	tokens, err := s.user.GetPasswordless(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListHumanPasswordlessResponse{
		Result: user.WebAuthNTokensViewToPb(tokens),
	}, nil
}

func (s *Server) RemoveHumanPasswordless(ctx context.Context, req *mgmt_pb.RemoveHumanPasswordlessRequest) (*mgmt_pb.RemoveHumanPasswordlessResponse, error) {
	details, err := s.command.HumanRemovePasswordless(ctx, req.UserId, req.TokenId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveHumanPasswordlessResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) UpdateMachine(ctx context.Context, req *mgmt_pb.UpdateMachineRequest) (*mgmt_pb.UpdateMachineResponse, error) {
	machine, err := s.command.ChangeMachine(ctx, UpdateMachineRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateMachineResponse{
		Details: obj_grpc.ToDetailsPb(
			machine.Sequence,
			machine.ChangeDate,
			machine.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetMachineKeyByIDs(ctx context.Context, req *mgmt_pb.GetMachineKeyByIDsRequest) (*mgmt_pb.GetMachineKeyByIDsResponse, error) {
	key, err := s.user.GetMachineKey(ctx, req.UserId, req.KeyId)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetMachineKeyByIDsResponse{
		Key: authn.KeyToPb(key),
	}, nil
}

func (s *Server) ListMachineKeys(ctx context.Context, req *mgmt_pb.ListMachineKeysRequest) (*mgmt_pb.ListMachineKeysResponse, error) {
	result, err := s.user.SearchMachineKeys(ctx, ListMachineKeysRequestToModel(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListMachineKeysResponse{
		Result: authn.KeyViewsToPb(result.Result),
		MetaData: obj_grpc.ToListDetails(
			result.TotalResult,
			result.Sequence,
			result.Timestamp,
		),
	}, nil
}

func (s *Server) AddMachineKey(ctx context.Context, req *mgmt_pb.AddMachineKeyRequest) (*mgmt_pb.AddMachineKeyResponse, error) {
	key, err := s.command.AddUserMachineKey(ctx, AddMachineKeyRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddMachineKeyResponse{
		KeyId:      key.KeyID,
		KeyDetails: authn.KeyDetailsToPb(key),
		Details: object.ToDetailsPb(
			key.Sequence,
			key.ChangeDate,
			key.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveMachineKey(ctx context.Context, req *mgmt_pb.RemoveMachineKeyRequest) (*mgmt_pb.RemoveMachineKeyResponse, error) {
	err := s.command.RemoveUserMachineKey(ctx, req.UserId, req.KeyId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveMachineKeyResponse{
		//TODO: details
	}, nil
}

func (s *Server) ListUserIDPs(ctx context.Context, req *mgmt_pb.ListUserIDPsRequest) (*mgmt_pb.ListUserIDPsResponse, error) {
	res, err := s.user.SearchExternalIDPs(ctx, ListUserIDPsRequestToModel(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListUserIDPsResponse{
		Result: idp_grpc.IDPsToUserLinkPb(res.Result),
		MetaData: obj_grpc.ToListDetails(
			res.TotalResult,
			res.Sequence,
			res.Timestamp,
		),
	}, nil
}
func (s *Server) RemoveUserIDP(ctx context.Context, req *mgmt_pb.RemoveUserIDPRequest) (*mgmt_pb.RemoveUserIDPResponse, error) {
	details, err := s.command.RemoveHumanExternalIDP(ctx, RemoveUserIDPRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveUserIDPResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func RemoveUserIDPRequestToDomain(ctx context.Context, req *mgmt_pb.RemoveUserIDPRequest) *domain.ExternalIDP {
	return &domain.ExternalIDP{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   req.UserId,
			ResourceOwner: authz.GetCtxData(ctx).OrgID,
		},
		IDPConfigID:    req.IdpId,
		ExternalUserID: req.LinkedUserId,
	}
}

func (s *Server) ListUserMemberships(ctx context.Context, req *mgmt_pb.ListUserMembershipsRequest) (*mgmt_pb.ListUserMembershipsResponse, error) {
	request, err := ListUserMembershipsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	response, err := s.user.SearchUserMemberships(ctx, request)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListUserMembershipsResponse{
		Result: user_grpc.MembershipsToMembershipsPb(response.Result),
		MetaData: obj_grpc.ToListDetails(
			response.TotalResult,
			response.Sequence,
			response.Timestamp,
		),
	}, nil
}
