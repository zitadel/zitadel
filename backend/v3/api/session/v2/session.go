package v2

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func DeleteSession(ctx context.Context, request *connect.Request[session.DeleteSessionRequest]) (*connect.Response[session.DeleteSessionResponse], error) {
	sessionDeleteCmd := domain.NewDeleteSessionCommand(request.Msg.GetSessionId(), request.Msg.GetSessionToken(), true)

	err := domain.Invoke(ctx, sessionDeleteCmd,
		domain.WithSessionRepo(repository.SessionRepository()),
	)

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&session.DeleteSessionResponse{
		Details: &object.Details{
			ChangeDate: timestamppb.New(sessionDeleteCmd.DeletedAt),
		},
	}), nil
}
