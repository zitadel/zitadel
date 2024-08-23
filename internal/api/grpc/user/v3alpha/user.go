package user

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/user/v3alpha"
)

func (s *Server) CreateUser(ctx context.Context, req *user.CreateUserRequest) (_ *user.CreateUserResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	schemauser, err := createUserRequestToCreateSchemaUser(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := s.command.CreateSchemaUser(ctx, schemauser, s.userCodeAlg, s.userCodeAlg); err != nil {
		return nil, err
	}
	return &user.CreateUserResponse{
		UserId:    schemauser.ID,
		Details:   object.DomainToDetailsPb(schemauser.Details),
		EmailCode: gu.Ptr(schemauser.ReturnCodeEmail),
		PhoneCode: gu.Ptr(schemauser.ReturnCodePhone),
	}, nil
}

func createUserRequestToCreateSchemaUser(ctx context.Context, req *user.CreateUserRequest) (*command.CreateSchemaUser, error) {
	data, err := req.GetData().MarshalJSON()
	if err != nil {
		return nil, err
	}
	return &command.CreateSchemaUser{
		ResourceOwner: authz.GetCtxData(ctx).OrgID,
		SchemaID:      req.SchemaId,
		ID:            req.GetUserId(),
		Data:          data,
	}, nil
}

func checkUserSchemaEnabled(ctx context.Context) error {
	if authz.GetInstance(ctx).Features().UserSchema {
		return nil
	}
	return zerrors.ThrowPreconditionFailed(nil, "TODO", "Errors.UserSchema.NotEnabled")
}
