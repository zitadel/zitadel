package v2

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	v2_object "github.com/zitadel/zitadel/pkg/grpc/object/v2"
	v2_session "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func DeleteSession(ctx context.Context, request *connect.Request[v2_session.DeleteSessionRequest]) (*connect.Response[v2_session.DeleteSessionResponse], error) {
	// TODO stebenz cryptoalgorithm
	sessionDeleteCmd := domain.NewDeleteSessionCommand(request.Msg.GetSessionId(), authz.SessionTokenVerifier(nil))

	err := domain.Invoke(ctx, sessionDeleteCmd) //domain.WithSessionRepo(repository.SessionRepository()),

	if err != nil {
		var notFoundError *database.NoRowFoundError
		if errors.As(err, &notFoundError) {
			return connect.NewResponse(&v2_session.DeleteSessionResponse{}), nil
		}
		return nil, err
	}

	return &connect.Response[v2_session.DeleteSessionResponse]{
		Msg: &v2_session.DeleteSessionResponse{
			Details: &v2_object.Details{
				Sequence:      0,
				ChangeDate:    nil,
				ResourceOwner: "",
				CreationDate:  nil,
			},
		},
	}, nil
}
