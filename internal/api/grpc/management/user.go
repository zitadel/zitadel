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
		Details: obj_grpc.ToListDetails(
			res.TotalResult,
			res.Sequence,
			res.Timestamp,
		),
	}, nil
}

func (s *Server) ListUserChanges(ctx context.Context, req *mgmt_pb.ListUserChangesRequest) (*mgmt_pb.ListUserChangesResponse, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	res, err := s.user.UserChanges(ctx, req.UserId, offset, limit, asc)
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
	objectDetails, err := s.command.DeactivateUser(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeactivateUserResponse{
		Details: obj_grpc.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ReactivateUser(ctx context.Context, req *mgmt_pb.ReactivateUserRequest) (*mgmt_pb.ReactivateUserResponse, error) {
	objectDetails, err := s.command.ReactivateUser(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateUserResponse{
		Details: obj_grpc.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) LockUser(ctx context.Context, req *mgmt_pb.LockUserRequest) (*mgmt_pb.LockUserResponse, error) {
	objectDetails, err := s.command.LockUser(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.LockUserResponse{
		Details: obj_grpc.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) UnlockUser(ctx context.Context, req *mgmt_pb.UnlockUserRequest) (*mgmt_pb.UnlockUserResponse, error) {
	objectDetails, err := s.command.UnlockUser(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UnlockUserResponse{
		Details: obj_grpc.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveUser(ctx context.Context, req *mgmt_pb.RemoveUserRequest) (*mgmt_pb.RemoveUserResponse, error) {
	grants, err := s.usergrant.UserGrantsByUserID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	objectDetails, err := s.command.RemoveUser(ctx, req.Id, authz.GetCtxData(ctx).OrgID, userGrantsToIDs(grants)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveUserResponse{
		Details: obj_grpc.DomainToDetailsPb(objectDetails),
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
	objectDetails, err := s.command.ChangeUsername(ctx, authz.GetCtxData(ctx).OrgID, req.UserId, req.UserName)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateUserNameResponse{
		Details: obj_grpc.DomainToDetailsPb(objectDetails),
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
	details, err := s.command.ResendInitialMail(ctx, req.UserId, req.Email, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResendHumanInitializationResponse{
		Details: obj_grpc.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) ResendHumanEmailVerification(ctx context.Context, req *mgmt_pb.ResendHumanEmailVerificationRequest) (*mgmt_pb.ResendHumanEmailVerificationResponse, error) {
	objectDetails, err := s.command.CreateHumanEmailVerificationCode(ctx, req.UserId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResendHumanEmailVerificationResponse{
		Details: obj_grpc.DomainToDetailsPb(objectDetails),
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
	objectDetails, err := s.command.RemoveHumanPhone(ctx, req.UserId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveHumanPhoneResponse{
		Details: obj_grpc.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ResendHumanPhoneVerification(ctx context.Context, req *mgmt_pb.ResendHumanPhoneVerificationRequest) (*mgmt_pb.ResendHumanPhoneVerificationResponse, error) {
	objectDetails, err := s.command.CreateHumanPhoneVerificationCode(ctx, req.UserId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResendHumanPhoneVerificationResponse{
		Details: obj_grpc.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) SetHumanInitialPassword(ctx context.Context, req *mgmt_pb.SetHumanInitialPasswordRequest) (*mgmt_pb.SetHumanInitialPasswordResponse, error) {
	objectDetails, err := s.command.SetOneTimePassword(ctx, authz.GetCtxData(ctx).OrgID, req.UserId, req.Password)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetHumanInitialPasswordResponse{
		Details: obj_grpc.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) SendHumanResetPasswordNotification(ctx context.Context, req *mgmt_pb.SendHumanResetPasswordNotificationRequest) (*mgmt_pb.SendHumanResetPasswordNotificationResponse, error) {
	objectDetails, err := s.command.RequestSetPassword(ctx, req.UserId, authz.GetCtxData(ctx).OrgID, notifyTypeToDomain(req.Type))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SendHumanResetPasswordNotificationResponse{
		Details: obj_grpc.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ListHumanAuthFactors(ctx context.Context, req *mgmt_pb.ListHumanAuthFactorsRequest) (*mgmt_pb.ListHumanAuthFactorsResponse, error) {
	mfas, err := s.user.UserMFAs(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListHumanAuthFactorsResponse{
		Result: user_grpc.AuthFactorsToPb(mfas),
	}, nil
}

func (s *Server) RemoveHumanAuthFactorOTP(ctx context.Context, req *mgmt_pb.RemoveHumanAuthFactorOTPRequest) (*mgmt_pb.RemoveHumanAuthFactorOTPResponse, error) {
	objectDetails, err := s.command.HumanRemoveOTP(ctx, req.UserId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveHumanAuthFactorOTPResponse{
		Details: obj_grpc.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveHumanAuthFactorU2F(ctx context.Context, req *mgmt_pb.RemoveHumanAuthFactorU2FRequest) (*mgmt_pb.RemoveHumanAuthFactorU2FResponse, error) {
	objectDetails, err := s.command.HumanRemoveU2F(ctx, req.UserId, req.TokenId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveHumanAuthFactorU2FResponse{
		Details: obj_grpc.DomainToDetailsPb(objectDetails),
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
	objectDetails, err := s.command.HumanRemovePasswordless(ctx, req.UserId, req.TokenId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveHumanPasswordlessResponse{
		Details: obj_grpc.DomainToDetailsPb(objectDetails),
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
		Details: obj_grpc.ToListDetails(
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
	keyDetails, err := key.Detail()
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddMachineKeyResponse{
		KeyId:      key.KeyID,
		KeyDetails: keyDetails,
		Details: object.ToDetailsPb(
			key.Sequence,
			key.ChangeDate,
			key.ResourceOwner,
		),
	}, nil
}

func (s *Server) RemoveMachineKey(ctx context.Context, req *mgmt_pb.RemoveMachineKeyRequest) (*mgmt_pb.RemoveMachineKeyResponse, error) {
	objectDetails, err := s.command.RemoveUserMachineKey(ctx, req.UserId, req.KeyId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveMachineKeyResponse{
		Details: obj_grpc.DomainToDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ListHumanLinkedIDPs(ctx context.Context, req *mgmt_pb.ListHumanLinkedIDPsRequest) (*mgmt_pb.ListHumanLinkedIDPsResponse, error) {
	res, err := s.user.SearchExternalIDPs(ctx, ListHumanLinkedIDPsRequestToModel(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListHumanLinkedIDPsResponse{
		Result: idp_grpc.IDPsToUserLinkPb(res.Result),
		Details: obj_grpc.ToListDetails(
			res.TotalResult,
			res.Sequence,
			res.Timestamp,
		),
	}, nil
}
func (s *Server) RemoveHumanLinkedIDP(ctx context.Context, req *mgmt_pb.RemoveHumanLinkedIDPRequest) (*mgmt_pb.RemoveHumanLinkedIDPResponse, error) {
	objectDetails, err := s.command.RemoveHumanExternalIDP(ctx, RemoveHumanLinkedIDPRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveHumanLinkedIDPResponse{
		Details: obj_grpc.DomainToDetailsPb(objectDetails),
	}, nil
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
		Details: obj_grpc.ToListDetails(
			response.TotalResult,
			response.Sequence,
			response.Timestamp,
		),
	}, nil
}
