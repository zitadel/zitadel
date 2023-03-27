package management

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/authn"
	change_grpc "github.com/zitadel/zitadel/internal/api/grpc/change"
	idp_grpc "github.com/zitadel/zitadel/internal/api/grpc/idp"
	"github.com/zitadel/zitadel/internal/api/grpc/metadata"
	obj_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	user_grpc "github.com/zitadel/zitadel/internal/api/grpc/user"
	"github.com/zitadel/zitadel/internal/api/http"
	z_oidc "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/user"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Server) GetUserByID(ctx context.Context, req *mgmt_pb.GetUserByIDRequest) (*mgmt_pb.GetUserByIDResponse, error) {
	owner, err := query.NewUserResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID, query.TextEquals)
	if err != nil {
		return nil, err
	}
	user, err := s.query.GetUserByID(ctx, true, req.Id, false, owner)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetUserByIDResponse{
		User: user_grpc.UserToPb(user, s.assetAPIPrefix(ctx)),
	}, nil
}

func (s *Server) GetUserByLoginNameGlobal(ctx context.Context, req *mgmt_pb.GetUserByLoginNameGlobalRequest) (*mgmt_pb.GetUserByLoginNameGlobalResponse, error) {
	loginName, err := query.NewUserPreferredLoginNameSearchQuery(req.LoginName, query.TextEquals)
	if err != nil {
		return nil, err
	}
	user, err := s.query.GetUser(ctx, true, false, loginName)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetUserByLoginNameGlobalResponse{
		User: user_grpc.UserToPb(user, s.assetAPIPrefix(ctx)),
	}, nil
}

func (s *Server) ListUsers(ctx context.Context, req *mgmt_pb.ListUsersRequest) (*mgmt_pb.ListUsersResponse, error) {
	queries, err := ListUsersRequestToModel(req)
	if err != nil {
		return nil, err
	}

	err = queries.AppendMyResourceOwnerQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	res, err := s.query.SearchUsers(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListUsersResponse{
		Result:  user_grpc.UsersToPb(res.Users, s.assetAPIPrefix(ctx)),
		Details: obj_grpc.ToListDetails(res.Count, res.Sequence, res.Timestamp),
	}, nil
}

func (s *Server) ListUserChanges(ctx context.Context, req *mgmt_pb.ListUserChangesRequest) (*mgmt_pb.ListUserChangesResponse, error) {
	var (
		limit    uint64
		sequence uint64
		asc      bool
	)
	if req.Query != nil {
		limit = uint64(req.Query.Limit)
		sequence = req.Query.Sequence
		asc = req.Query.Asc
	}

	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AllowTimeTravel().
		Limit(limit).
		OrderDesc().
		ResourceOwner(authz.GetCtxData(ctx).OrgID).
		AddQuery().
		SequenceGreater(sequence).
		AggregateTypes(user.AggregateType).
		AggregateIDs(req.UserId).
		Builder()
	if asc {
		query.OrderAsc()
	}

	changes, err := s.query.SearchEvents(ctx, query, s.auditLogRetention)
	if err != nil {
		return nil, err
	}

	return &mgmt_pb.ListUserChangesResponse{
		Result: change_grpc.EventsToChangesPb(changes, s.assetAPIPrefix(ctx)),
	}, nil
}

func (s *Server) IsUserUnique(ctx context.Context, req *mgmt_pb.IsUserUniqueRequest) (*mgmt_pb.IsUserUniqueResponse, error) {
	orgID := authz.GetCtxData(ctx).OrgID
	policy, err := s.query.DomainPolicyByOrg(ctx, true, orgID, false)
	if err != nil {
		return nil, err
	}
	if !policy.UserLoginMustBeDomain {
		orgID = ""
	}
	unique, err := s.query.IsUserUnique(ctx, req.UserName, req.Email, orgID, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.IsUserUniqueResponse{
		IsUnique: unique,
	}, nil
}

func (s *Server) ListUserMetadata(ctx context.Context, req *mgmt_pb.ListUserMetadataRequest) (*mgmt_pb.ListUserMetadataResponse, error) {
	metadataQueries, err := ListUserMetadataToDomain(req)
	if err != nil {
		return nil, err
	}
	err = metadataQueries.AppendMyResourceOwnerQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	res, err := s.query.SearchUserMetadata(ctx, true, req.Id, metadataQueries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListUserMetadataResponse{
		Result:  metadata.UserMetadataListToPb(res.Metadata),
		Details: obj_grpc.ToListDetails(res.Count, res.Sequence, res.Timestamp),
	}, nil
}

func (s *Server) GetUserMetadata(ctx context.Context, req *mgmt_pb.GetUserMetadataRequest) (*mgmt_pb.GetUserMetadataResponse, error) {
	owner, err := query.NewUserMetadataResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	data, err := s.query.GetUserMetadataByKey(ctx, true, req.Id, req.Key, false, owner)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetUserMetadataResponse{
		Metadata: metadata.UserMetadataToPb(data),
	}, nil
}

func (s *Server) SetUserMetadata(ctx context.Context, req *mgmt_pb.SetUserMetadataRequest) (*mgmt_pb.SetUserMetadataResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	result, err := s.command.SetUserMetadata(ctx, &domain.Metadata{Key: req.Key, Value: req.Value}, req.Id, ctxData.OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetUserMetadataResponse{
		Details: obj_grpc.AddToDetailsPb(
			result.Sequence,
			result.ChangeDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) BulkSetUserMetadata(ctx context.Context, req *mgmt_pb.BulkSetUserMetadataRequest) (*mgmt_pb.BulkSetUserMetadataResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	result, err := s.command.BulkSetUserMetadata(ctx, req.Id, ctxData.OrgID, BulkSetUserMetadataToDomain(req)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.BulkSetUserMetadataResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(result),
	}, nil
}

func (s *Server) RemoveUserMetadata(ctx context.Context, req *mgmt_pb.RemoveUserMetadataRequest) (*mgmt_pb.RemoveUserMetadataResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	result, err := s.command.RemoveUserMetadata(ctx, req.Key, req.Id, ctxData.OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveUserMetadataResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(result),
	}, nil
}

func (s *Server) BulkRemoveUserMetadata(ctx context.Context, req *mgmt_pb.BulkRemoveUserMetadataRequest) (*mgmt_pb.BulkRemoveUserMetadataResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	result, err := s.command.BulkRemoveUserMetadata(ctx, req.Id, ctxData.OrgID, req.Keys...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.BulkRemoveUserMetadataResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(result),
	}, nil
}

func (s *Server) AddHumanUser(ctx context.Context, req *mgmt_pb.AddHumanUserRequest) (*mgmt_pb.AddHumanUserResponse, error) {
	details, err := s.command.AddHuman(ctx, authz.GetCtxData(ctx).OrgID, AddHumanUserRequestToAddHuman(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddHumanUserResponse{
		UserId: details.ID,
		Details: obj_grpc.AddToDetailsPb(
			details.Sequence,
			details.EventDate,
			details.ResourceOwner,
		),
	}, nil
}

func AddHumanUserRequestToAddHuman(req *mgmt_pb.AddHumanUserRequest) *command.AddHuman {
	lang, err := language.Parse(req.Profile.PreferredLanguage)
	logging.OnError(err).Debug("unable to parse language")

	human := &command.AddHuman{
		Username:    req.UserName,
		FirstName:   req.Profile.FirstName,
		LastName:    req.Profile.LastName,
		NickName:    req.Profile.NickName,
		DisplayName: req.Profile.DisplayName,
		Email: command.Email{
			Address:  domain.EmailAddress(req.Email.Email),
			Verified: req.Email.IsEmailVerified,
		},
		PreferredLanguage:      lang,
		Gender:                 user_grpc.GenderToDomain(req.Profile.Gender),
		Password:               req.InitialPassword,
		PasswordChangeRequired: true,
		Passwordless:           false,
		Register:               false,
		ExternalIDP:            false,
	}
	if req.Phone != nil {
		human.Phone = command.Phone{
			Number:   domain.PhoneNumber(req.Phone.Phone),
			Verified: req.Phone.IsPhoneVerified,
		}
	}
	return human
}

func (s *Server) ImportHumanUser(ctx context.Context, req *mgmt_pb.ImportHumanUserRequest) (*mgmt_pb.ImportHumanUserResponse, error) {
	human, passwordless, links := ImportHumanUserRequestToDomain(req)
	initCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeInitCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	emailCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyEmailCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	phoneCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyPhoneCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	passwordlessInitCode, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypePasswordlessInitCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	addedHuman, code, err := s.command.ImportHuman(ctx, authz.GetCtxData(ctx).OrgID, human, passwordless, links, initCodeGenerator, phoneCodeGenerator, emailCodeGenerator, passwordlessInitCode)
	if err != nil {
		return nil, err
	}
	resp := &mgmt_pb.ImportHumanUserResponse{
		UserId: addedHuman.AggregateID,
		Details: obj_grpc.AddToDetailsPb(
			addedHuman.Sequence,
			addedHuman.ChangeDate,
			addedHuman.ResourceOwner,
		),
	}
	if code != nil {
		origin := http.BuildOrigin(authz.GetInstance(ctx).RequestedHost(), s.externalSecure)
		resp.PasswordlessRegistration = &mgmt_pb.ImportHumanUserResponse_PasswordlessRegistration{
			Link:       code.Link(origin + login.HandlerPrefix + login.EndpointPasswordlessRegistration),
			Lifetime:   durationpb.New(code.Expiration),
			Expiration: durationpb.New(code.Expiration),
		}
	}
	return resp, nil
}

func (s *Server) AddMachineUser(ctx context.Context, req *mgmt_pb.AddMachineUserRequest) (*mgmt_pb.AddMachineUserResponse, error) {
	machine := AddMachineUserRequestToCommand(req, authz.GetCtxData(ctx).OrgID)
	objectDetails, err := s.command.AddMachine(ctx, machine)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddMachineUserResponse{
		UserId:  machine.AggregateID,
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) DeactivateUser(ctx context.Context, req *mgmt_pb.DeactivateUserRequest) (*mgmt_pb.DeactivateUserResponse, error) {
	objectDetails, err := s.command.DeactivateUser(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeactivateUserResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ReactivateUser(ctx context.Context, req *mgmt_pb.ReactivateUserRequest) (*mgmt_pb.ReactivateUserResponse, error) {
	objectDetails, err := s.command.ReactivateUser(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateUserResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) LockUser(ctx context.Context, req *mgmt_pb.LockUserRequest) (*mgmt_pb.LockUserResponse, error) {
	objectDetails, err := s.command.LockUser(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.LockUserResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) UnlockUser(ctx context.Context, req *mgmt_pb.UnlockUserRequest) (*mgmt_pb.UnlockUserResponse, error) {
	objectDetails, err := s.command.UnlockUser(ctx, req.Id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UnlockUserResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveUser(ctx context.Context, req *mgmt_pb.RemoveUserRequest) (*mgmt_pb.RemoveUserResponse, error) {
	memberships, grants, err := s.removeUserDependencies(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	objectDetails, err := s.command.RemoveUser(ctx, req.Id, authz.GetCtxData(ctx).OrgID, memberships, grants...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveUserResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) removeUserDependencies(ctx context.Context, userID string) ([]*command.CascadingMembership, []string, error) {
	userGrantUserQuery, err := query.NewUserGrantUserIDSearchQuery(userID)
	if err != nil {
		return nil, nil, err
	}
	grants, err := s.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: []query.SearchQuery{userGrantUserQuery},
	}, true, true)
	if err != nil {
		return nil, nil, err
	}
	membershipsUserQuery, err := query.NewMembershipUserIDQuery(userID)
	if err != nil {
		return nil, nil, err
	}
	memberships, err := s.query.Memberships(ctx, &query.MembershipSearchQuery{
		Queries: []query.SearchQuery{membershipsUserQuery},
	}, true)
	if err != nil {
		return nil, nil, err
	}
	return cascadingMemberships(memberships.Memberships), userGrantsToIDs(grants.UserGrants), nil
}

func (s *Server) UpdateUserName(ctx context.Context, req *mgmt_pb.UpdateUserNameRequest) (*mgmt_pb.UpdateUserNameResponse, error) {
	objectDetails, err := s.command.ChangeUsername(ctx, authz.GetCtxData(ctx).OrgID, req.UserId, req.UserName)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateUserNameResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) GetHumanProfile(ctx context.Context, req *mgmt_pb.GetHumanProfileRequest) (*mgmt_pb.GetHumanProfileResponse, error) {
	owner, err := query.NewUserResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID, query.TextEquals)
	if err != nil {
		return nil, err
	}
	profile, err := s.query.GetHumanProfile(ctx, req.UserId, false, owner)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetHumanProfileResponse{
		Profile: user_grpc.ProfileToPb(profile, s.assetAPIPrefix(ctx)),
		Details: obj_grpc.ToViewDetailsPb(
			profile.Sequence,
			profile.CreationDate,
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
		Details: obj_grpc.ChangeToDetailsPb(
			profile.Sequence,
			profile.ChangeDate,
			profile.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetHumanEmail(ctx context.Context, req *mgmt_pb.GetHumanEmailRequest) (*mgmt_pb.GetHumanEmailResponse, error) {
	owner, err := query.NewUserResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID, query.TextEquals)
	if err != nil {
		return nil, err
	}
	email, err := s.query.GetHumanEmail(ctx, req.UserId, false, owner)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetHumanEmailResponse{
		Email: user_grpc.EmailToPb(email),
		Details: obj_grpc.ToViewDetailsPb(
			email.Sequence,
			email.CreationDate,
			email.ChangeDate,
			email.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateHumanEmail(ctx context.Context, req *mgmt_pb.UpdateHumanEmailRequest) (*mgmt_pb.UpdateHumanEmailResponse, error) {
	emailCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyEmailCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	email, err := s.command.ChangeHumanEmail(ctx, UpdateHumanEmailRequestToDomain(ctx, req), emailCodeGenerator)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateHumanEmailResponse{
		Details: obj_grpc.ChangeToDetailsPb(
			email.Sequence,
			email.ChangeDate,
			email.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResendHumanInitialization(ctx context.Context, req *mgmt_pb.ResendHumanInitializationRequest) (*mgmt_pb.ResendHumanInitializationResponse, error) {
	initCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeInitCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	details, err := s.command.ResendInitialMail(ctx, req.UserId, domain.EmailAddress(req.Email), authz.GetCtxData(ctx).OrgID, initCodeGenerator)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResendHumanInitializationResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ResendHumanEmailVerification(ctx context.Context, req *mgmt_pb.ResendHumanEmailVerificationRequest) (*mgmt_pb.ResendHumanEmailVerificationResponse, error) {
	emailCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyEmailCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	objectDetails, err := s.command.CreateHumanEmailVerificationCode(ctx, req.UserId, authz.GetCtxData(ctx).OrgID, emailCodeGenerator)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResendHumanEmailVerificationResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) GetHumanPhone(ctx context.Context, req *mgmt_pb.GetHumanPhoneRequest) (*mgmt_pb.GetHumanPhoneResponse, error) {
	owner, err := query.NewUserResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID, query.TextEquals)
	if err != nil {
		return nil, err
	}
	phone, err := s.query.GetHumanPhone(ctx, req.UserId, false, owner)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetHumanPhoneResponse{
		Phone: user_grpc.PhoneToPb(phone),
		Details: obj_grpc.ToViewDetailsPb(
			phone.Sequence,
			phone.CreationDate,
			phone.ChangeDate,
			phone.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateHumanPhone(ctx context.Context, req *mgmt_pb.UpdateHumanPhoneRequest) (*mgmt_pb.UpdateHumanPhoneResponse, error) {
	phoneCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyPhoneCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	phone, err := s.command.ChangeHumanPhone(ctx, UpdateHumanPhoneRequestToDomain(req), authz.GetCtxData(ctx).OrgID, phoneCodeGenerator)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateHumanPhoneResponse{
		Details: obj_grpc.ChangeToDetailsPb(
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
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ResendHumanPhoneVerification(ctx context.Context, req *mgmt_pb.ResendHumanPhoneVerificationRequest) (*mgmt_pb.ResendHumanPhoneVerificationResponse, error) {
	phoneCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypeVerifyPhoneCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	objectDetails, err := s.command.CreateHumanPhoneVerificationCode(ctx, req.UserId, authz.GetCtxData(ctx).OrgID, phoneCodeGenerator)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResendHumanPhoneVerificationResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveHumanAvatar(ctx context.Context, req *mgmt_pb.RemoveHumanAvatarRequest) (*mgmt_pb.RemoveHumanAvatarResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	objectDetails, err := s.command.RemoveHumanAvatar(ctx, ctxData.OrgID, req.UserId)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveHumanAvatarResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) SetHumanInitialPassword(ctx context.Context, req *mgmt_pb.SetHumanInitialPasswordRequest) (*mgmt_pb.SetHumanInitialPasswordResponse, error) {
	objectDetails, err := s.command.SetPassword(ctx, authz.GetCtxData(ctx).OrgID, req.UserId, req.Password, true)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetHumanInitialPasswordResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) SetHumanPassword(ctx context.Context, req *mgmt_pb.SetHumanPasswordRequest) (*mgmt_pb.SetHumanPasswordResponse, error) {
	objectDetails, err := s.command.SetPassword(ctx, authz.GetCtxData(ctx).OrgID, req.UserId, req.Password, !req.NoChangeRequired)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SetHumanPasswordResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) SendHumanResetPasswordNotification(ctx context.Context, req *mgmt_pb.SendHumanResetPasswordNotificationRequest) (*mgmt_pb.SendHumanResetPasswordNotificationResponse, error) {
	passwordCodeGenerator, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypePasswordResetCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	objectDetails, err := s.command.RequestSetPassword(ctx, req.UserId, authz.GetCtxData(ctx).OrgID, notifyTypeToDomain(req.Type), passwordCodeGenerator)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SendHumanResetPasswordNotificationResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ListHumanAuthFactors(ctx context.Context, req *mgmt_pb.ListHumanAuthFactorsRequest) (*mgmt_pb.ListHumanAuthFactorsResponse, error) {
	query := new(query.UserAuthMethodSearchQueries)
	err := query.AppendUserIDQuery(req.UserId)
	if err != nil {
		return nil, err
	}
	err = query.AppendAuthMethodsQuery(domain.UserAuthMethodTypeU2F, domain.UserAuthMethodTypeOTP)
	if err != nil {
		return nil, err
	}
	err = query.AppendStateQuery(domain.MFAStateReady)
	if err != nil {
		return nil, err
	}
	authMethods, err := s.query.SearchUserAuthMethods(ctx, query, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListHumanAuthFactorsResponse{
		Result: user_grpc.AuthMethodsToPb(authMethods),
	}, nil
}

func (s *Server) RemoveHumanAuthFactorOTP(ctx context.Context, req *mgmt_pb.RemoveHumanAuthFactorOTPRequest) (*mgmt_pb.RemoveHumanAuthFactorOTPResponse, error) {
	objectDetails, err := s.command.HumanRemoveOTP(ctx, req.UserId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveHumanAuthFactorOTPResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) RemoveHumanAuthFactorU2F(ctx context.Context, req *mgmt_pb.RemoveHumanAuthFactorU2FRequest) (*mgmt_pb.RemoveHumanAuthFactorU2FResponse, error) {
	objectDetails, err := s.command.HumanRemoveU2F(ctx, req.UserId, req.TokenId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveHumanAuthFactorU2FResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ListHumanPasswordless(ctx context.Context, req *mgmt_pb.ListHumanPasswordlessRequest) (*mgmt_pb.ListHumanPasswordlessResponse, error) {
	query := new(query.UserAuthMethodSearchQueries)
	err := query.AppendUserIDQuery(req.UserId)
	if err != nil {
		return nil, err
	}
	err = query.AppendAuthMethodQuery(domain.UserAuthMethodTypePasswordless)
	if err != nil {
		return nil, err
	}
	err = query.AppendStateQuery(domain.MFAStateReady)
	if err != nil {
		return nil, err
	}
	authMethods, err := s.query.SearchUserAuthMethods(ctx, query, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListHumanPasswordlessResponse{
		Result: user_grpc.UserAuthMethodsToWebAuthNTokenPb(authMethods),
	}, nil
}

func (s *Server) AddPasswordlessRegistration(ctx context.Context, req *mgmt_pb.AddPasswordlessRegistrationRequest) (*mgmt_pb.AddPasswordlessRegistrationResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	passwordlessInitCode, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypePasswordlessInitCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	initCode, err := s.command.HumanAddPasswordlessInitCode(ctx, req.UserId, ctxData.OrgID, passwordlessInitCode)
	if err != nil {
		return nil, err
	}
	origin := http.BuildOrigin(authz.GetInstance(ctx).RequestedHost(), s.externalSecure)
	return &mgmt_pb.AddPasswordlessRegistrationResponse{
		Details:    obj_grpc.AddToDetailsPb(initCode.Sequence, initCode.ChangeDate, initCode.ResourceOwner),
		Link:       initCode.Link(origin + login.HandlerPrefix + login.EndpointPasswordlessRegistration),
		Expiration: durationpb.New(initCode.Expiration),
	}, nil
}

func (s *Server) SendPasswordlessRegistration(ctx context.Context, req *mgmt_pb.SendPasswordlessRegistrationRequest) (*mgmt_pb.SendPasswordlessRegistrationResponse, error) {
	ctxData := authz.GetCtxData(ctx)
	passwordlessInitCode, err := s.query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypePasswordlessInitCode, s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	initCode, err := s.command.HumanSendPasswordlessInitCode(ctx, req.UserId, ctxData.OrgID, passwordlessInitCode)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.SendPasswordlessRegistrationResponse{
		Details: obj_grpc.AddToDetailsPb(initCode.Sequence, initCode.ChangeDate, initCode.ResourceOwner),
	}, nil
}

func (s *Server) RemoveHumanPasswordless(ctx context.Context, req *mgmt_pb.RemoveHumanPasswordlessRequest) (*mgmt_pb.RemoveHumanPasswordlessResponse, error) {
	objectDetails, err := s.command.HumanRemovePasswordless(ctx, req.UserId, req.TokenId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveHumanPasswordlessResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) UpdateMachine(ctx context.Context, req *mgmt_pb.UpdateMachineRequest) (*mgmt_pb.UpdateMachineResponse, error) {
	machine := UpdateMachineRequestToCommand(req, authz.GetCtxData(ctx).OrgID)
	objectDetails, err := s.command.ChangeMachine(ctx, machine)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateMachineResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) GetMachineKeyByIDs(ctx context.Context, req *mgmt_pb.GetMachineKeyByIDsRequest) (*mgmt_pb.GetMachineKeyByIDsResponse, error) {
	resourceOwner, err := query.NewAuthNKeyResourceOwnerQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	aggregateID, err := query.NewAuthNKeyAggregateIDQuery(req.UserId)
	if err != nil {
		return nil, err
	}
	key, err := s.query.GetAuthNKeyByID(ctx, true, req.KeyId, false, resourceOwner, aggregateID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetMachineKeyByIDsResponse{
		Key: authn.KeyToPb(key),
	}, nil
}

func (s *Server) ListMachineKeys(ctx context.Context, req *mgmt_pb.ListMachineKeysRequest) (*mgmt_pb.ListMachineKeysResponse, error) {
	query, err := ListMachineKeysRequestToQuery(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := s.query.SearchAuthNKeys(ctx, query, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListMachineKeysResponse{
		Result:  authn.KeysToPb(result.AuthNKeys),
		Details: obj_grpc.ToListDetails(result.Count, result.Sequence, result.Timestamp),
	}, nil
}

func (s *Server) AddMachineKey(ctx context.Context, req *mgmt_pb.AddMachineKeyRequest) (*mgmt_pb.AddMachineKeyResponse, error) {
	machineKey := AddMachineKeyRequestToCommand(req, authz.GetCtxData(ctx).OrgID)
	details, err := s.command.AddUserMachineKey(ctx, machineKey)
	if err != nil {
		return nil, err
	}
	keyDetails, err := machineKey.Detail()
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddMachineKeyResponse{
		KeyId:      machineKey.KeyID,
		KeyDetails: keyDetails,
		Details:    obj_grpc.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) RemoveMachineKey(ctx context.Context, req *mgmt_pb.RemoveMachineKeyRequest) (*mgmt_pb.RemoveMachineKeyResponse, error) {
	objectDetails, err := s.command.RemoveUserMachineKey(ctx, RemoveMachineKeyRequestToCommand(req, authz.GetCtxData(ctx).OrgID))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveMachineKeyResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) GenerateMachineSecret(ctx context.Context, req *mgmt_pb.GenerateMachineSecretRequest) (*mgmt_pb.GenerateMachineSecretResponse, error) {
	// use SecretGeneratorTypeAppSecret as the secrets will be used in the client_credentials grant like a client secret
	secretGenerator, err := s.query.InitHashGenerator(ctx, domain.SecretGeneratorTypeAppSecret, s.passwordHashAlg)
	if err != nil {
		return nil, err
	}
	set := new(command.GenerateMachineSecret)
	details, err := s.command.GenerateMachineSecret(ctx, req.UserId, authz.GetCtxData(ctx).OrgID, secretGenerator, set)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GenerateMachineSecretResponse{
		ClientId:     set.ClientID,
		ClientSecret: set.ClientSecret,
		Details:      obj_grpc.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) RemoveMachineSecret(ctx context.Context, req *mgmt_pb.RemoveMachineSecretRequest) (*mgmt_pb.RemoveMachineSecretResponse, error) {
	objectDetails, err := s.command.RemoveMachineSecret(ctx, req.UserId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveMachineSecretResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) GetPersonalAccessTokenByIDs(ctx context.Context, req *mgmt_pb.GetPersonalAccessTokenByIDsRequest) (*mgmt_pb.GetPersonalAccessTokenByIDsResponse, error) {
	resourceOwner, err := query.NewPersonalAccessTokenResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	aggregateID, err := query.NewPersonalAccessTokenUserIDSearchQuery(req.UserId)
	if err != nil {
		return nil, err
	}
	token, err := s.query.PersonalAccessTokenByID(ctx, true, req.TokenId, false, resourceOwner, aggregateID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetPersonalAccessTokenByIDsResponse{
		Token: user_grpc.PersonalAccessTokenToPb(token),
	}, nil
}

func (s *Server) ListPersonalAccessTokens(ctx context.Context, req *mgmt_pb.ListPersonalAccessTokensRequest) (*mgmt_pb.ListPersonalAccessTokensResponse, error) {
	queries, err := ListPersonalAccessTokensRequestToQuery(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := s.query.SearchPersonalAccessTokens(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListPersonalAccessTokensResponse{
		Result:  user_grpc.PersonalAccessTokensToPb(result.PersonalAccessTokens),
		Details: obj_grpc.ToListDetails(result.Count, result.Sequence, result.Timestamp),
	}, nil
}

func (s *Server) AddPersonalAccessToken(ctx context.Context, req *mgmt_pb.AddPersonalAccessTokenRequest) (*mgmt_pb.AddPersonalAccessTokenResponse, error) {
	scopes := []string{oidc.ScopeOpenID, z_oidc.ScopeUserMetaData, z_oidc.ScopeResourceOwner}
	pat := AddPersonalAccessTokenRequestToCommand(req, authz.GetCtxData(ctx).OrgID, scopes, domain.UserTypeMachine)
	details, err := s.command.AddPersonalAccessToken(ctx, pat)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddPersonalAccessTokenResponse{
		TokenId: pat.TokenID,
		Token:   pat.Token,
		Details: obj_grpc.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) RemovePersonalAccessToken(ctx context.Context, req *mgmt_pb.RemovePersonalAccessTokenRequest) (*mgmt_pb.RemovePersonalAccessTokenResponse, error) {
	objectDetails, err := s.command.RemovePersonalAccessToken(ctx, RemovePersonalAccessTokenRequestToCommand(req, authz.GetCtxData(ctx).OrgID))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemovePersonalAccessTokenResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ListHumanLinkedIDPs(ctx context.Context, req *mgmt_pb.ListHumanLinkedIDPsRequest) (*mgmt_pb.ListHumanLinkedIDPsResponse, error) {
	queries, err := ListHumanLinkedIDPsRequestToQuery(ctx, req)
	if err != nil {
		return nil, err
	}
	res, err := s.query.IDPUserLinks(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListHumanLinkedIDPsResponse{
		Result:  idp_grpc.IDPUserLinksToPb(res.Links),
		Details: obj_grpc.ToListDetails(res.Count, res.Sequence, res.Timestamp),
	}, nil
}
func (s *Server) RemoveHumanLinkedIDP(ctx context.Context, req *mgmt_pb.RemoveHumanLinkedIDPRequest) (*mgmt_pb.RemoveHumanLinkedIDPResponse, error) {
	objectDetails, err := s.command.RemoveUserIDPLink(ctx, RemoveHumanLinkedIDPRequestToDomain(ctx, req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveHumanLinkedIDPResponse{
		Details: obj_grpc.DomainToChangeDetailsPb(objectDetails),
	}, nil
}

func (s *Server) ListUserMemberships(ctx context.Context, req *mgmt_pb.ListUserMembershipsRequest) (*mgmt_pb.ListUserMembershipsResponse, error) {
	request, err := ListUserMembershipsRequestToModel(ctx, req)
	if err != nil {
		return nil, err
	}
	response, err := s.query.Memberships(ctx, request, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListUserMembershipsResponse{
		Result:  user_grpc.MembershipsToMembershipsPb(response.Memberships),
		Details: obj_grpc.ToListDetails(response.Count, response.Sequence, response.Timestamp),
	}, nil
}

func cascadingMemberships(memberships []*query.Membership) []*command.CascadingMembership {
	cascades := make([]*command.CascadingMembership, len(memberships))
	for i, membership := range memberships {
		cascades[i] = &command.CascadingMembership{
			UserID:        membership.UserID,
			ResourceOwner: membership.ResourceOwner,
			IAM:           cascadingIAMMembership(membership.IAM),
			Org:           cascadingOrgMembership(membership.Org),
			Project:       cascadingProjectMembership(membership.Project),
			ProjectGrant:  cascadingProjectGrantMembership(membership.ProjectGrant),
		}
	}
	return cascades
}

func cascadingIAMMembership(membership *query.IAMMembership) *command.CascadingIAMMembership {
	if membership == nil {
		return nil
	}
	return &command.CascadingIAMMembership{IAMID: membership.IAMID}
}
func cascadingOrgMembership(membership *query.OrgMembership) *command.CascadingOrgMembership {
	if membership == nil {
		return nil
	}
	return &command.CascadingOrgMembership{OrgID: membership.OrgID}
}
func cascadingProjectMembership(membership *query.ProjectMembership) *command.CascadingProjectMembership {
	if membership == nil {
		return nil
	}
	return &command.CascadingProjectMembership{ProjectID: membership.ProjectID}
}
func cascadingProjectGrantMembership(membership *query.ProjectGrantMembership) *command.CascadingProjectGrantMembership {
	if membership == nil {
		return nil
	}
	return &command.CascadingProjectGrantMembership{ProjectID: membership.ProjectID, GrantID: membership.GrantID}
}

func userGrantsToIDs(userGrants []*query.UserGrant) []string {
	converted := make([]string, len(userGrants))
	for i, grant := range userGrants {
		converted[i] = grant.ID
	}
	return converted
}
