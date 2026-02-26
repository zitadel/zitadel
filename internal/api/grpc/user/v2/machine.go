package user

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (s *Server) createUserTypeMachine(ctx context.Context, machinePb *user.CreateUserRequest_Machine, orgId, userName, userId string) (*connect.Response[user.CreateUserResponse], error) {
	cmd := &command.Machine{
		Username:        userName,
		Name:            machinePb.Name,
		Description:     machinePb.GetDescription(),
		AccessTokenType: accessTokenTypeToDomain(machinePb.GetAccessTokenType()),
		ObjectRoot: models.ObjectRoot{
			ResourceOwner: orgId,
			AggregateID:   userId,
		},
	}
	details, err := s.command.AddMachine(
		ctx,
		cmd,
		nil,
		s.command.NewPermissionCheckUserWrite(ctx, true),
		command.AddMachineWithUsernameToIDFallback(),
	)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.CreateUserResponse{
		Id:           cmd.AggregateID,
		CreationDate: timestamppb.New(details.EventDate),
	}), nil
}

func (s *Server) updateUserTypeMachine(ctx context.Context, machinePb *user.UpdateUserRequest_Machine, userId string, userName *string) (*connect.Response[user.UpdateUserResponse], error) {
	cmd := updateMachineUserToCommand(userId, userName, machinePb)
	err := s.command.ChangeUserMachine(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&user.UpdateUserResponse{
		ChangeDate: timestamppb.New(cmd.Details.EventDate),
	}), nil
}

func updateMachineUserToCommand(userId string, userName *string, machine *user.UpdateUserRequest_Machine) *command.ChangeMachine {
	var accessTokenType *domain.OIDCTokenType
	if machine.AccessTokenType != nil {
		tokenType := accessTokenTypeToDomain(*machine.AccessTokenType)
		accessTokenType = &tokenType
	}
	return &command.ChangeMachine{
		ID:              userId,
		Username:        userName,
		Name:            machine.Name,
		Description:     machine.Description,
		AccessTokenType: accessTokenType,
	}
}

func accessTokenTypeToDomain(accessTokenType user.AccessTokenType) domain.OIDCTokenType {
	switch accessTokenType {
	case user.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER:
		return domain.OIDCTokenTypeBearer
	case user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT:
		return domain.OIDCTokenTypeJWT
	default:
		return domain.OIDCTokenTypeBearer
	}
}
