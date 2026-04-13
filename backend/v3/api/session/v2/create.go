package v2

import (
	"context"
	"net"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func CreateSession(ctx context.Context, request *connect.Request[session.CreateSessionRequest]) (*connect.Response[session.CreateSessionResponse], error) {
	return defaultServer.CreateSession(ctx, request)
}

// CreateSession implements [sessionconnect.SessionServiceHandler].
func (s *server) CreateSession(ctx context.Context, request *connect.Request[session.CreateSessionRequest]) (*connect.Response[session.CreateSessionResponse], error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	creatorID := authz.GetCtxData(ctx).UserID

	create := createSessionRequestToCommand(instanceID, creatorID, request.Msg)
	checks := checksToCommands(create, request.Msg.GetChecks())

	batch := domain.BatchExecutors(create)
	for _, check := range checks {
		batch.Append(check)
	}

	if err := domain.Invoke(ctx, batch); err != nil {
		return nil, err
	}

	return connect.NewResponse(&session.CreateSessionResponse{
		SessionId:    create.Result().ID,
		SessionToken: create.Result().TokenID, // TODO(adlerhurst): return the correct value
		Details: &object.Details{
			ResourceOwner: create.Result().InstanceID,
			CreationDate:  timestamppb.New(create.Result().CreatedAt),
		},
		Challenges: &session.Challenges{}, // TODO(adlerhurst): return the correct values
	}), nil
}

func createSessionRequestToCommand(instanceID, userID string, request *session.CreateSessionRequest) *domain.CreateSessionCommand {
	opts := make([]domain.CreateSessionOption, 0, 2)
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
	return domain.NewCreateSessionCommand(
		instanceID,
		userID,
		userAgentToDomain(request.GetUserAgent()),
		opts...,
	)
}

func userAgentToDomain(userAgent *session.UserAgent) *domain.SessionUserAgent {
	header := make(http.Header, len(userAgent.GetHeader()))
	for k, v := range userAgent.GetHeader() {
		header[k] = v.GetValues()
	}
	return &domain.SessionUserAgent{
		FingerprintID: userAgent.FingerprintId,
		Description:   userAgent.Description,
		IP:            net.ParseIP(userAgent.GetIp()),
		Header:        header,
	}
}
