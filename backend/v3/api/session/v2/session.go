package sessionv2

import (
	"context"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func CreateSession(ctx context.Context, request *connect.Request[session.CreateSessionRequest]) (*connect.Response[session.CreateSessionResponse], error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	executors := make([]domain.Executor, 12)
	creator := domain.NewCreateSessionCommand(
		instanceID,
		request.Msg.GetUserAgent(),
		request.Msg.GetMetadata(),
		request.Msg.GetLifetime(),
		nil,
	)
	sessionID := *creator.SessionID
	executors[0] = creator

	executors[1] = domain.NewUserCheckCommand(sessionID, instanceID)
	executors[2] = domain.NewPasswordCheckCommand(sessionID, instanceID, nil, nil, request.Msg.GetChecks().GetPassword())
	executors[3] = domain.NewIDPIntentCheckCommand(sessionID, instanceID, request.Msg.GetChecks().GetIdpIntent(), nil)
	executors[4] = domain.NewPasskeyCheckCommand(sessionID, instanceID, request.Msg.GetChecks().GetWebAuthN(), nil)
	executors[5] = domain.NewTOTPCheckCommand(sessionID, instanceID, nil, nil, nil, request.Msg.GetChecks().GetTotp())
	executors[6] = domain.NewOTPCheckCommand(sessionID, instanceID, nil, nil, nil, nil, request.Msg.GetChecks().GetOtpSms(), domain.OTPSMSRequestType)
	executors[7] = domain.NewOTPCheckCommand(sessionID, instanceID, nil, nil, nil, nil, request.Msg.GetChecks().GetOtpEmail(), domain.OTPEmailRequestType)
	executors[8] = domain.NewRecoveryCodeCheckCommand(sessionID, instanceID, request.Msg.GetChecks().GetRecoveryCode())

	// TODO(IAM-Marco): Add challenges
	// executors[9] = domain.NewPasskeyChallengeCommand()
	// executors[10] = domain.NewOTPChallengeCommand(domain.OTPSMSRequestType)
	// executors[11] = domain.NewOTPChallengeCommand(domain.OTPEmailRequestType)

	batcher := domain.BatchExecutors(executors...)

	err := domain.Invoke(ctx, batcher,
		// TODO(IAM-Marco): Uncomment when session repository is available
		// domain.WithSessionRepo(repository.SessionRepository()),
		domain.WithLockoutSettingsRepo(repository.LockoutSettingsRepository()),
		// TODO(IAM-Marco): Uncomment when user repository is available
		// domain.WithUserRepo(repository.UserRepository()),
	)
	if err != nil {
		return nil, err
	}

	return &connect.Response[session.CreateSessionResponse]{
		Msg: &session.CreateSessionResponse{
			Details:      &object.Details{},
			SessionId:    sessionID,
			SessionToken: "", // TODO(IAM-Marco): Where do I take this from?
			Challenges:   &session.Challenges{},
		},
	}, nil
}

// GetSession returns the session for the given session ID and token (optional)
func GetSession(ctx context.Context, req *connect.Request[session.GetSessionRequest]) (*connect.Response[session.GetSessionResponse], error) {
	sessionGetQuery := domain.NewGetSessionQuery(
		req.Msg.GetSessionId(),
		authz.GetInstance(ctx).InstanceID(),
		req.Msg.GetSessionToken(),
		nil, // todo: check passing the verifier fn here?
	)
	err := domain.Invoke(ctx,
		sessionGetQuery,
		// domain.WithSessionRepo(repository.SessionRepository()), // todo: uncomment when these repos are available
		// domain.WithUserRepo(repository.UserRepository()),
	)
	if err != nil {
		return nil, err
	}

	return &connect.Response[session.GetSessionResponse]{
		Msg: &session.GetSessionResponse{
			Session: sessionToPb(sessionGetQuery.Result()),
		},
	}, nil
}
