package eventstore

import (
	"context"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/auth_request/repository/cache"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/id"
	user_model "github.com/caos/zitadel/internal/user/model"
	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type AuthRequestRepo struct {
	UserEvents   *user_event.UserEventstore
	AuthRequests *cache.AuthRequestCache
	View         *view.View

	UserSessionViewProvider userSessionViewProvider
	UserViewProvider        userViewProvider
	UserEventProvider       userEventProvider

	IdGenerator id.Generator

	PasswordCheckLifeTime    time.Duration
	MfaInitSkippedLifeTime   time.Duration
	MfaSoftwareCheckLifeTime time.Duration
	MfaHardwareCheckLifeTime time.Duration
}

type userSessionViewProvider interface {
	UserSessionByIDs(string, string) (*view_model.UserSessionView, error)
	UserSessionsByAgentID(string) ([]*view_model.UserSessionView, error)
}
type userViewProvider interface {
	UserByID(string) (*view_model.UserView, error)
}

type userEventProvider interface {
	UserEventsByID(ctx context.Context, id string, sequence uint64) ([]*es_models.Event, error)
}

func (repo *AuthRequestRepo) Health(ctx context.Context) error {
	if err := repo.UserEvents.Health(ctx); err != nil {
		return err
	}
	return repo.AuthRequests.Health(ctx)
}

func (repo *AuthRequestRepo) CreateAuthRequest(ctx context.Context, request *model.AuthRequest) (*model.AuthRequest, error) {
	reqID, err := repo.IdGenerator.Next()
	if err != nil {
		return nil, err
	}
	request.ID = reqID
	err = repo.AuthRequests.SaveAuthRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (repo *AuthRequestRepo) AuthRequestByID(ctx context.Context, id string) (*model.AuthRequest, error) {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return nil, err
	}
	steps, err := repo.nextSteps(ctx, request)
	if err != nil {
		return nil, err
	}
	request.PossibleSteps = steps
	return request, nil
}

func (repo *AuthRequestRepo) CheckUsername(ctx context.Context, id, username string) error {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return err
	}
	user, err := repo.View.UserByUsername(username)
	if err != nil {
		return err
	}
	request.SetUserInfo(user.ID, user.UserName, user.ResourceOwner)
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) VerifyPassword(ctx context.Context, id, userID, password string, info *model.BrowserInfo) error {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return err
	}
	if request.UserID != userID {
		return errors.ThrowPreconditionFailed(nil, "EVENT-ds35D", "user id does not match request id")
	}
	return repo.UserEvents.CheckPassword(ctx, userID, password, request.WithCurrentInfo(info))
}

func (repo *AuthRequestRepo) VerifyMfaOTP(ctx context.Context, authRequestID, userID string, code string, info *model.BrowserInfo) error {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, authRequestID)
	if err != nil {
		return err
	}
	if request.UserID != userID {
		return errors.ThrowPreconditionFailed(nil, "EVENT-ADJ26", "user id does not match request id")
	}
	return repo.UserEvents.CheckMfaOTP(ctx, userID, code, request.WithCurrentInfo(info))
}

func (repo *AuthRequestRepo) nextSteps(ctx context.Context, request *model.AuthRequest) ([]model.NextStep, error) {
	if request == nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-ds27a", "request must not be nil")
	}
	steps := make([]model.NextStep, 0)
	if request.UserID == "" {
		if request.Prompt != model.PromptNone {
			steps = append(steps, &model.LoginStep{})
		}
		if request.Prompt == model.PromptSelectAccount {
			users, err := repo.usersForUserSelection(request)
			if err != nil {
				return nil, err
			}
			steps = append(steps, &model.SelectUserStep{Users: users})
		}
		return steps, nil
	}
	userSession, err := userSessionByIDs(ctx, repo.UserSessionViewProvider, repo.UserEventProvider, request.AgentID, request.UserID)
	if err != nil {
		return nil, err
	}
	user, err := userByID(ctx, repo.UserViewProvider, repo.UserEventProvider, request.UserID)
	if err != nil {
		return nil, err
	}

	if !user.PasswordSet {
		return append(steps, &model.InitPasswordStep{}), nil
	}

	if !checkVerificationTime(userSession.PasswordVerification, repo.PasswordCheckLifeTime) {
		return append(steps, &model.PasswordStep{}), nil
	}

	if step, ok := repo.mfaChecked(userSession, request, user); !ok {
		return append(steps, step), nil
	}

	if user.PasswordChangeRequired {
		steps = append(steps, &model.ChangePasswordStep{})
	}
	if !user.IsEmailVerified {
		steps = append(steps, &model.VerifyEMailStep{})
	}

	if user.PasswordChangeRequired || !user.IsEmailVerified {
		return steps, nil
	}

	//PLANNED: consent step
	return append(steps, &model.RedirectToCallbackStep{}), nil
}

func (repo *AuthRequestRepo) usersForUserSelection(request *model.AuthRequest) ([]model.UserSelection, error) {
	userSessions, err := userSessionsByUserAgentID(repo.UserSessionViewProvider, request.AgentID)
	if err != nil {
		return nil, err
	}
	users := make([]model.UserSelection, len(userSessions))
	for i, session := range userSessions {
		users[i] = model.UserSelection{
			UserID:           session.UserID,
			UserName:         session.UserName,
			UserSessionState: session.State,
		}
	}
	return users, nil
}

func (repo *AuthRequestRepo) mfaChecked(userSession *user_model.UserSessionView, request *model.AuthRequest, user *user_model.UserView) (model.NextStep, bool) {
	mfaLevel := request.MfaLevel()
	promptRequired := user.MfaMaxSetUp < mfaLevel
	if promptRequired || !repo.mfaSkippedOrSetUp(user) {
		return &model.MfaPromptStep{
			Required:     promptRequired,
			MfaProviders: user.MfaTypesSetupPossible(mfaLevel),
		}, false
	}
	switch mfaLevel {
	default:
		fallthrough
	case model.MfaLevelNotSetUp:
		if user.MfaMaxSetUp == model.MfaLevelNotSetUp {
			return nil, true
		}
		fallthrough
	case model.MfaLevelSoftware:
		if checkVerificationTime(userSession.MfaSoftwareVerification, repo.MfaSoftwareCheckLifeTime) {
			return nil, true
		}
		fallthrough
	case model.MfaLevelHardware:
		if checkVerificationTime(userSession.MfaHardwareVerification, repo.MfaHardwareCheckLifeTime) {
			return nil, true
		}
	}
	return &model.MfaVerificationStep{
		MfaProviders: user.MfaTypesAllowed(mfaLevel),
	}, false
}

func (repo *AuthRequestRepo) mfaSkippedOrSetUp(user *user_model.UserView) bool {
	if user.MfaMaxSetUp > model.MfaLevelNotSetUp {
		return true
	}
	return checkVerificationTime(user.MfaInitSkipped, repo.MfaInitSkippedLifeTime)
}

func checkVerificationTime(verificationTime time.Time, lifetime time.Duration) bool {
	return verificationTime.Add(lifetime).After(time.Now().UTC())
}

func userSessionsByUserAgentID(provider userSessionViewProvider, agentID string) ([]*user_model.UserSessionView, error) {
	session, err := provider.UserSessionsByAgentID(agentID)
	if err != nil {
		return nil, err
	}
	return view_model.UserSessionsToModel(session), nil
}

func userSessionByIDs(ctx context.Context, provider userSessionViewProvider, eventProvider userEventProvider, agentID, userID string) (*user_model.UserSessionView, error) {
	session, err := provider.UserSessionByIDs(agentID, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			return &user_model.UserSessionView{}, nil
		}
		return nil, err
	}
	events, err := eventProvider.UserEventsByID(ctx, userID, session.Sequence)
	if err != nil {
		logging.Log("EVENT-Hse6s").WithError(err).Debug("error retrieving new events")
		return view_model.UserSessionToModel(session), nil
	}
	for _, event := range events {
		session.AppendEvent(event)
	}
	return view_model.UserSessionToModel(session), nil
}

func userByID(ctx context.Context, viewProvider userViewProvider, eventProvider userEventProvider, userID string) (*user_model.UserView, error) {
	user, err := viewProvider.UserByID(userID)
	if err != nil {
		return nil, err
	}
	events, err := eventProvider.UserEventsByID(ctx, userID, user.Sequence)
	if err != nil {
		logging.Log("EVENT-dfg42").WithError(err).Debug("error retrieving new events")
		return view_model.UserToModel(user), nil
	}
	userCopy := *user
	for _, event := range events {
		if err := userCopy.AppendEvent(event); err != nil {
			return view_model.UserToModel(user), nil
		}
	}
	return view_model.UserToModel(&userCopy), nil
}
