package v2

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func SetSession(ctx context.Context, request *connect.Request[session.SetSessionRequest]) (*connect.Response[session.SetSessionResponse], error) {
	return defaultServer.SetSession(ctx, request)
}

// SetSession implements [sessionconnect.SessionServiceHandler].
func (s *server) SetSession(ctx context.Context, request *connect.Request[session.SetSessionRequest]) (*connect.Response[session.SetSessionResponse], error) {
	instanceID := authz.GetInstance(ctx).InstanceID()

	set := setSessionRequestToCommand(instanceID, request.Msg)
	checks := checksToCommands(set, request.Msg.GetChecks())

	batch := domain.BatchExecutors(set)
	for _, check := range checks {
		batch.Append(check)
	}

	if err := domain.Invoke(ctx, batch); err != nil {
		return nil, err
	}

	return connect.NewResponse(&session.SetSessionResponse{
		SessionToken: set.Result().TokenID,
		Details: &object.Details{
			ResourceOwner: set.Result().InstanceID,
			CreationDate:  timestamppb.New(set.Result().CreatedAt),
			ChangeDate:    timestamppb.New(set.Result().UpdatedAt),
		},
		Challenges: &session.Challenges{}, // TODO(adlerhurst): return the correct values
	}), nil
}

func setSessionRequestToCommand(instanceID string, request *session.SetSessionRequest) *domain.SetSessionCommand {
	opts := make([]domain.SetSessionOption, 0, 2)
	if len(request.GetMetadata()) > 0 {
		metadata := make([]*domain.SessionMetadata, 0, len(request.GetMetadata()))
		for key, value := range request.GetMetadata() {
			metadata = append(metadata, &domain.SessionMetadata{
				Metadata: domain.Metadata{
					InstanceID: instanceID,
					Key:        key,
					Value:      value,
				},
			})
		}
		opts = append(opts, domain.WithSessionMetadata(metadata...))
	}
	if request.GetLifetime().IsValid() {
		opts = append(opts, domain.WithSessionLifetime(request.GetLifetime().AsDuration()))
	}
	return domain.NewSetSessionCommand(instanceID, request.GetSessionId(), opts...)
}
