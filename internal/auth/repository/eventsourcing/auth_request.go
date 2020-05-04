package eventsourcing

import (
	"context"
	"time"

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
	return nextStepsNoUserSelected(request, false)
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
	//if request.PasswordChecked() {
	//	return nil, errors.ThrowPreconditionFailed(nil, "EVENT-52NGs", "user already chosen")
	//}
	return nil, errors.ThrowUnimplemented(nil, "EVENT-asjod", "user by username not yet implemented")
	if err != nil {
		return nextStepsNoUserSelected(request, true)
	}
	return nextSteps(request, user)
}

func (repo *AuthRequestRepo) VerifyPassword(ctx context.Context, id, userID, password string, info *req_model.BrowserInfo) (*req_model.AuthRequest, error) {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return nil, err
	}
	//if request.UserID == 0 {
	//
	//}
	//if request.PasswordChecked() {
	//	return nil, errors.ThrowPreconditionFailed(nil, "EVENT-s6Gn3", "password already checked")
	//}
	user, err := repo.UserEvents.CheckPassword(ctx, userID, password, request.AggregateID)
	return nextSteps(request, user)
}

func (repo *AuthRequestRepo) SkipMfaInit(ctx context.Context, authRequestID, userID string) (*req_model.AuthRequest, error) {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, authRequestID)
	if err != nil {
		return nil, err
	}
	if err = repo.UserEvents.SkipMfaInit(ctx, userID); err != nil {
		return nil, err
	}
	user, err := repo.UserEvents.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return nextSteps(request, user)
}

func (repo *AuthRequestRepo) VerifyMfaOTP(ctx context.Context, authRequestID, userID string, code string, info *req_model.BrowserInfo) (*req_model.AuthRequest, error) {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, authRequestID)
	if err != nil {
		return nil, err
	}
	if err = repo.UserEvents.CheckMfaOTP(ctx, userID, code); err != nil {
		return nil, err
	}
	user, err := repo.UserEvents.UserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return nextSteps(request, user)
}

func nextStepsNoUserSelected(request *req_model.AuthRequest, notFound bool) (*req_model.AuthRequest, error) {
	if request.Prompt != req_model.PromptNone {
		request.AddPossibleStep(&req_model.LoginStep{NotFound: notFound})
	}
	//TODO: select account
	return request, nil
}

func nextSteps(request *req_model.AuthRequest, user *model.User) (*req_model.AuthRequest, error) {
	if user == nil {
		return nextStepsNoUserSelected(request, true)
	}
	if user.Password == nil {
		request.AddPossibleStep(&req_model.InitPasswordStep{})
		return request, nil
	}
	if ok, count := user.PasswordVerified(request.AggregateID); !ok {
		request.AddPossibleStep(&req_model.PasswordStep{FailureCount: count})
		return request, nil
	}
	minimalLevel := MfaLevel(request)
	if len(user.MfaTypesReady(minimalLevel)) > 0 {
		if ok, count := user.MfaVerified(request.AggregateID); !ok {
			request.AddPossibleStep(&req_model.MfaVerificationStep{
				FailureCount: count,
				MfaProviders: user.MfaTypesReady(minimalLevel),
			})
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
	mfaRequired := MfaRequired(request)
	if MfaNotSkippedAndNotReady(user) {
		request.AddPossibleStep(&req_model.MfaPromptStep{
			Required:     mfaRequired,
			MfaProviders: user.MfaTypesPossible(),
		})
		if mfaRequired {
			return request, nil
		}
	}
	//TODO: consent step
	if authenticated() {
		request.AddPossibleStep(&req_model.RedirectToCallbackStep{})
		return request, nil
	}
	return request, nil
}

func MfaLevel(request *req_model.AuthRequest) model.MfaLevel {
	return model.MfaLevelSoftware //TODO: map acr_values
}

func MfaRequired(request *req_model.AuthRequest) bool {
	return request.IsMfaRequired() //TODO: add policies (org requires mfa, org whitelist?)
}

func MfaNotSkippedAndNotReady(user *model.User) bool {
	skipDuration := 30 * 24 * time.Hour
	return user.SkippedMfaInit.Add(skipDuration).Before(time.Now().UTC()) &&
		len(user.MfaTypesReady()) == 0
}

func authenticated() bool {

}
