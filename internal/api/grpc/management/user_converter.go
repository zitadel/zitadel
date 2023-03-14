package management

import (
	"context"
	"time"

	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/pkg/grpc/user"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/authn"
	"github.com/zitadel/zitadel/internal/api/grpc/metadata"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	user_grpc "github.com/zitadel/zitadel/internal/api/grpc/user"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	user_model "github.com/zitadel/zitadel/internal/user/model"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func ListUsersRequestToModel(req *mgmt_pb.ListUsersRequest) (*query.UserSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := user_grpc.UserQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.UserSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: UserFieldNameToSortingColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func UserFieldNameToSortingColumn(field user.UserFieldName) query.Column {
	switch field {
	case user.UserFieldName_USER_FIELD_NAME_EMAIL:
		return query.HumanEmailCol
	case user.UserFieldName_USER_FIELD_NAME_FIRST_NAME:
		return query.HumanFirstNameCol
	case user.UserFieldName_USER_FIELD_NAME_LAST_NAME:
		return query.HumanLastNameCol
	case user.UserFieldName_USER_FIELD_NAME_DISPLAY_NAME:
		return query.HumanDisplayNameCol
	case user.UserFieldName_USER_FIELD_NAME_USER_NAME:
		return query.UserUsernameCol
	case user.UserFieldName_USER_FIELD_NAME_STATE:
		return query.UserStateCol
	case user.UserFieldName_USER_FIELD_NAME_TYPE:
		return query.UserTypeCol
	case user.UserFieldName_USER_FIELD_NAME_NICK_NAME:
		return query.HumanNickNameCol
	case user.UserFieldName_USER_FIELD_NAME_CREATION_DATE:
		return query.UserCreationDateCol
	default:
		return query.UserIDCol
	}
}

func BulkSetUserMetadataToDomain(req *mgmt_pb.BulkSetUserMetadataRequest) []*domain.Metadata {
	metadata := make([]*domain.Metadata, len(req.Metadata))
	for i, data := range req.Metadata {
		metadata[i] = &domain.Metadata{
			Key:   data.Key,
			Value: data.Value,
		}
	}
	return metadata
}

func ListUserMetadataToDomain(req *mgmt_pb.ListUserMetadataRequest) (*query.UserMetadataSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := metadata.MetadataQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.UserMetadataSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}

func ImportHumanUserRequestToDomain(req *mgmt_pb.ImportHumanUserRequest) (human *domain.Human, passwordless bool, links []*domain.UserIDPLink) {
	human = &domain.Human{
		Username: req.UserName,
	}
	preferredLanguage, err := language.Parse(req.Profile.PreferredLanguage)
	logging.Log("MANAG-3GUFJ").OnError(err).Debug("language malformed")
	human.Profile = &domain.Profile{
		FirstName:         req.Profile.FirstName,
		LastName:          req.Profile.LastName,
		NickName:          req.Profile.NickName,
		DisplayName:       req.Profile.DisplayName,
		PreferredLanguage: preferredLanguage,
		Gender:            user_grpc.GenderToDomain(req.Profile.Gender),
	}
	human.Email = &domain.Email{
		EmailAddress:    domain.EmailAddress(req.Email.Email),
		IsEmailVerified: req.Email.IsEmailVerified,
	}
	if req.Phone != nil {
		human.Phone = &domain.Phone{
			PhoneNumber:     domain.PhoneNumber(req.Phone.Phone),
			IsPhoneVerified: req.Phone.IsPhoneVerified,
		}
	}

	if req.Password != "" {
		human.Password = domain.NewPassword(req.Password)
		human.Password.ChangeRequired = req.PasswordChangeRequired
	}

	if req.HashedPassword != nil && req.HashedPassword.Value != "" && req.HashedPassword.Algorithm != "" {
		human.HashedPassword = domain.NewHashedPassword(req.HashedPassword.Value, req.HashedPassword.Algorithm)
	}
	links = make([]*domain.UserIDPLink, len(req.Idps))
	for i, idp := range req.Idps {
		links[i] = &domain.UserIDPLink{
			IDPConfigID:    idp.ConfigId,
			ExternalUserID: idp.ExternalUserId,
			DisplayName:    idp.DisplayName,
		}
	}

	return human, req.RequestPasswordlessRegistration, links
}

func AddMachineUserRequestToCommand(req *mgmt_pb.AddMachineUserRequest, resourceowner string) *command.Machine {
	return &command.Machine{
		ObjectRoot: models.ObjectRoot{
			ResourceOwner: resourceowner,
		},
		Username:        req.UserName,
		Name:            req.Name,
		Description:     req.Description,
		AccessTokenType: user_grpc.AccessTokenTypeToDomain(req.AccessTokenType),
	}
}

func UpdateHumanProfileRequestToDomain(req *mgmt_pb.UpdateHumanProfileRequest) *domain.Profile {
	preferredLanguage, err := language.Parse(req.PreferredLanguage)
	logging.Log("MANAG-GPcYv").OnError(err).Debug("language malformed")
	return &domain.Profile{
		ObjectRoot:        models.ObjectRoot{AggregateID: req.UserId},
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		NickName:          req.NickName,
		DisplayName:       req.DisplayName,
		PreferredLanguage: preferredLanguage,
		Gender:            user_grpc.GenderToDomain(req.Gender),
	}
}

func UpdateHumanEmailRequestToDomain(ctx context.Context, req *mgmt_pb.UpdateHumanEmailRequest) *domain.Email {
	return &domain.Email{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   req.UserId,
			ResourceOwner: authz.GetCtxData(ctx).OrgID,
		},
		EmailAddress:    domain.EmailAddress(req.Email),
		IsEmailVerified: req.IsEmailVerified,
	}
}

func UpdateHumanPhoneRequestToDomain(req *mgmt_pb.UpdateHumanPhoneRequest) *domain.Phone {
	return &domain.Phone{
		ObjectRoot:      models.ObjectRoot{AggregateID: req.UserId},
		PhoneNumber:     domain.PhoneNumber(req.Phone),
		IsPhoneVerified: req.IsPhoneVerified,
	}
}

func notifyTypeToDomain(state mgmt_pb.SendHumanResetPasswordNotificationRequest_Type) domain.NotificationType {
	switch state {
	case mgmt_pb.SendHumanResetPasswordNotificationRequest_TYPE_EMAIL:
		return domain.NotificationTypeEmail
	case mgmt_pb.SendHumanResetPasswordNotificationRequest_TYPE_SMS:
		return domain.NotificationTypeSms
	default:
		return domain.NotificationTypeEmail
	}
}

func UpdateMachineRequestToCommand(req *mgmt_pb.UpdateMachineRequest, orgID string) *command.Machine {
	return &command.Machine{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   req.UserId,
			ResourceOwner: orgID,
		},
		Name:            req.Name,
		Description:     req.Description,
		AccessTokenType: user_grpc.AccessTokenTypeToDomain(req.AccessTokenType),
	}
}

func ListMachineKeysRequestToQuery(ctx context.Context, req *mgmt_pb.ListMachineKeysRequest) (*query.AuthNKeySearchQueries, error) {
	resourcOwner, err := query.NewAuthNKeyResourceOwnerQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	userID, err := query.NewAuthNKeyAggregateIDQuery(req.UserId)
	if err != nil {
		return nil, err
	}
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &query.AuthNKeySearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: []query.SearchQuery{
			resourcOwner,
			userID,
		},
	}, nil

}

func AddMachineKeyRequestToCommand(req *mgmt_pb.AddMachineKeyRequest, resourceOwner string) *command.MachineKey {
	expDate := time.Time{}
	if req.ExpirationDate != nil {
		expDate = req.ExpirationDate.AsTime()
	}

	return &command.MachineKey{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   req.UserId,
			ResourceOwner: resourceOwner,
		},
		ExpirationDate: expDate,
		Type:           authn.KeyTypeToDomain(req.Type),
	}
}

func RemoveMachineKeyRequestToCommand(req *mgmt_pb.RemoveMachineKeyRequest, resourceOwner string) *command.MachineKey {
	return &command.MachineKey{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   req.UserId,
			ResourceOwner: resourceOwner,
		},
		KeyID: req.KeyId,
	}
}

func AddPersonalAccessTokenRequestToCommand(req *mgmt_pb.AddPersonalAccessTokenRequest, resourceOwner string, scopes []string, allowedUserType domain.UserType) *command.PersonalAccessToken {
	expDate := time.Time{}
	if req.ExpirationDate != nil {
		expDate = req.ExpirationDate.AsTime()
	}

	return &command.PersonalAccessToken{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   req.UserId,
			ResourceOwner: resourceOwner,
		},
		ExpirationDate:  expDate,
		Scopes:          scopes,
		AllowedUserType: allowedUserType,
	}
}

func RemovePersonalAccessTokenRequestToCommand(req *mgmt_pb.RemovePersonalAccessTokenRequest, resourceOwner string) *command.PersonalAccessToken {
	return &command.PersonalAccessToken{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   req.UserId,
			ResourceOwner: resourceOwner,
		},
		TokenID: req.TokenId,
	}
}

func ListPersonalAccessTokensRequestToQuery(ctx context.Context, req *mgmt_pb.ListPersonalAccessTokensRequest) (*query.PersonalAccessTokenSearchQueries, error) {
	resourceOwner, err := query.NewPersonalAccessTokenResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	userID, err := query.NewPersonalAccessTokenUserIDSearchQuery(req.UserId)
	if err != nil {
		return nil, err
	}
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &query.PersonalAccessTokenSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: []query.SearchQuery{
			resourceOwner,
			userID,
		},
	}, nil

}

func RemoveHumanLinkedIDPRequestToDomain(ctx context.Context, req *mgmt_pb.RemoveHumanLinkedIDPRequest) *domain.UserIDPLink {
	return &domain.UserIDPLink{
		ObjectRoot: models.ObjectRoot{
			AggregateID:   req.UserId,
			ResourceOwner: authz.GetCtxData(ctx).OrgID,
		},
		IDPConfigID:    req.IdpId,
		ExternalUserID: req.LinkedUserId,
	}
}

func ListHumanLinkedIDPsRequestToQuery(ctx context.Context, req *mgmt_pb.ListHumanLinkedIDPsRequest) (*query.IDPUserLinksSearchQuery, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	userQuery, err := query.NewIDPUserLinksUserIDSearchQuery(req.UserId)
	if err != nil {
		return nil, err
	}
	resourceOwnerQuery, err := query.NewIDPUserLinksResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &query.IDPUserLinksSearchQuery{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: []query.SearchQuery{userQuery, resourceOwnerQuery},
	}, nil
}

func ListUserMembershipsRequestToModel(ctx context.Context, req *mgmt_pb.ListUserMembershipsRequest) (*query.MembershipSearchQuery, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := user_grpc.MembershipQueriesToQuery(req.Queries)
	if err != nil {
		return nil, err
	}
	userQuery, err := query.NewMembershipUserIDQuery(req.UserId)
	if err != nil {
		return nil, err
	}
	ownerQuery, err := query.NewMembershipResourceOwnersSearchQuery(authz.GetInstance(ctx).InstanceID(), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	queries = append(queries, userQuery, ownerQuery)
	return &query.MembershipSearchQuery{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		//SortingColumn: //TODO: sorting
		Queries: queries,
	}, nil
}

func UserMembershipViewsToDomain(memberships []*user_model.UserMembershipView) []*domain.UserMembership {
	result := make([]*domain.UserMembership, len(memberships))
	for i, membership := range memberships {
		result[i] = &domain.UserMembership{
			UserID:            membership.UserID,
			MemberType:        MemberTypeToDomain(membership.MemberType),
			AggregateID:       membership.AggregateID,
			ObjectID:          membership.ObjectID,
			Roles:             membership.Roles,
			DisplayName:       membership.DisplayName,
			CreationDate:      membership.CreationDate,
			ChangeDate:        membership.ChangeDate,
			ResourceOwner:     membership.ResourceOwner,
			ResourceOwnerName: membership.ResourceOwnerName,
			Sequence:          membership.Sequence,
		}
	}
	return result
}

func MemberTypeToDomain(mType user_model.MemberType) domain.MemberType {
	switch mType {
	case user_model.MemberTypeIam:
		return domain.MemberTypeIam
	case user_model.MemberTypeOrganisation:
		return domain.MemberTypeOrganisation
	case user_model.MemberTypeProject:
		return domain.MemberTypeProject
	case user_model.MemberTypeProjectGrant:
		return domain.MemberTypeProjectGrant
	default:
		return domain.MemberTypeUnspecified
	}
}
