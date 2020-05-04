package eventsourcing

import (
	"context"

	req_model "github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/auth_request/repository/cache"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/user/model"
	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type AuthRequestRepo struct {
	UserEvents   *user_event.UserEventstore
	AuthRequests *cache.AuthRequestCache
	//view      *view.View
}

func (repo *AuthRequestRepo) CreateAuthRequest(ctx context.Context, request *req_model.AuthRequest) (*req_model.AuthRequest, error) {
	err := repo.AuthRequests.SaveAuthRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return nextSteps(request, nil)
}

func (repo *AuthRequestRepo) AuthRequestByID(ctx context.Context, id string) (*req_model.AuthRequest, error) {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return nextSteps(request, nil)
}

func (repo *AuthRequestRepo) CheckUsername(ctx context.Context, id, username string) (*req_model.AuthRequest, error) {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if request.PasswordChecked() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-52NGs", "user already chosen")
	}
	return nil, errors.ThrowUnimplemented(nil, "EVENT-asjod", "user by username not yet implemented")
}

func (repo *AuthRequestRepo) VerifyPassword(ctx context.Context, id, userID, password string, info *req_model.BrowserInfo) (*req_model.AuthRequest, error) {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if request.PasswordChecked() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-s6Gn3", "password already checked")
	}
	if err = repo.UserEvents.CheckPassword(ctx, userID, password); err != nil {

	}
	return nextSteps(request, nil)
}

func (repo *AuthRequestRepo) RequestPasswordReset(ctx context.Context, id, userID string, info *req_model.BrowserInfo) (*req_model.AuthRequest, error) { //?

}

func (repo *AuthRequestRepo) SkipMfaInit(ctx context.Context, id, userID string) (*req_model.AuthRequest, error) {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err = repo.UserEvents.SkipMfaInit(ctx, userID); err != nil {
		return nil, err
	}

}
func (repo *AuthRequestRepo) AddMfa(ctx context.Context, agentID, authRequestID string, mfa interface{}, info *req_model.BrowserInfo) (*req_model.AuthRequest, error)
func (repo *AuthRequestRepo) VerifyMfa(ctx context.Context, agentID, authRequestID string, mfa interface{}, info *req_model.BrowserInfo) (*req_model.AuthRequest, error)

func nextSteps(request *req_model.AuthRequest, user *model.User) (*req_model.AuthRequest, error) {
	if user == nil {
		if request.Prompt != req_model.PromptNone {
			request.AddPossibleStep(req_model.NewLoginStep(err))
		}
		//TODO: select account
		return request, nil
	}
	if user.Password == nil {
		request.AddPossibleStep(&req_model.InitPasswordStep{})
		return request, nil
	}
	if ok, count := user.PasswordVerified(request.AggregateID); !ok { //TODO: ???
		request.AddPossibleStep(req_model.NewPasswordStep(count))
		return request, nil
	}
	if len(user.MfaTypesReady() > 0) {
		if ok, count := user.MfaVerified(request.AggregateID); !ok { //TODO: ???
			request.AddPossibleStep(req_model.NewMfaVerificationStep(count, user.MfaTypesReady()))
			return request, nil
		}
	}
	if user.Password.ChangeRequired {
		request.AddPossibleStep(&req_model.ChangePasswordStep{})
	}
	if user.Email == nil || user.Email != nil && !user.Email.IsEmailVerified {
		request.AddPossibleStep(&req_model.VerifyEMailStep{})
		return request, nil
	}

	return request, nil
}
