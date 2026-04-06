package v2

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/api/session/v2/convert"
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func ListSessions(ctx context.Context, request *connect.Request[session.ListSessionsRequest]) (*connect.Response[session.ListSessionsResponse], error) {
	req, err := convert.ListSessionsRequestGRPCToDomain(request.Msg)
	if err != nil {
		return nil, err
	}

	q := domain.NewListSessionsQuery(req)
	err = domain.Invoke(ctx, q, domain.WithSessionRepo(repository.SessionRepository()))
	if err != nil {
		return nil, err
	}

	sessions := convert.DomainSessionListToGRPCResponse(q.Result())
	return connect.NewResponse(&session.ListSessionsResponse{
		Sessions: sessions,
		Details: &object.ListDetails{
			TotalResult: uint64(len(sessions)),

			// TODO(IAM-Marco): Put something meaningful once we have a reliable timestamp from DB
			Timestamp: timestamppb.Now(),
		},
	}), nil
}

func DeleteSession(ctx context.Context, request *connect.Request[session.DeleteSessionRequest]) (*connect.Response[session.DeleteSessionResponse], error) {
	sessionDeleteCmd := domain.NewDeleteSessionCommand(request.Msg.GetSessionId(), request.Msg.GetSessionToken(), true)

	err := domain.Invoke(ctx, sessionDeleteCmd,
		domain.WithSessionRepo(repository.SessionRepository()),
	)

	if err != nil {
		return nil, err
	}

	details := &object.Details{}
	if !sessionDeleteCmd.DeletedAt.IsZero() {
		details.ChangeDate = timestamppb.New(sessionDeleteCmd.DeletedAt)
	}

	return connect.NewResponse(&session.DeleteSessionResponse{
		Details: details,
	}), nil
}
