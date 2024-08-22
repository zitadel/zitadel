package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/user/v3alpha"
)

func (s *Server) CreateUser(ctx context.Context, req *user.CreateUserRequest) (_ *user.CreateUserResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	schemauser, err := createUserRequestToCreateSchemaUser(req)
	if err != nil {
		return nil, err
	}
	orgID := authz.GetCtxData(ctx).OrgID
	id, details, err := s.command.CreateSchemaUser(ctx, orgID, schemauser)
	if err != nil {
		return nil, err
	}
	return &user.CreateUserResponse{
		UserId:    id,
		Details:   object.DomainToDetailsPb(details),
		EmailCode: human.EmailCode,
		PhoneCode: human.PhoneCode,
	}, nil
}

func createUserRequestToCreateSchemaUser(req *user.CreateUserRequest) (*command.CreateSchemaUser, error) {
	data, err := req.GetData().MarshalJSON()
	if err != nil {
		return nil, err
	}
	return &command.CreateSchemaUser{
		SchemaID: req.SchemaId,
		ID:       req.GetUserId(),
		Data:     data,
	}, nil
}

func checkUserSchemaEnabled(ctx context.Context) error {
	if authz.GetInstance(ctx).Features().UserSchema {
		return nil
	}
	return zerrors.ThrowPreconditionFailed(nil, "TODO", "Errors.UserSchema.NotEnabled")
}
