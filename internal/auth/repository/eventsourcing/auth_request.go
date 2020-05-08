package eventsourcing

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/auth_request/repository/cache"
	"github.com/caos/zitadel/internal/errors"
	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type AuthRequestRepo struct {
	UserEvents   *user_event.UserEventstore
	AuthRequests *cache.AuthRequestCache
	//view      *view.View

	PasswordCheckLifeTime    time.Duration
	MfaInitSkippedLifeTime   time.Duration
	MfaSoftwareCheckLifeTime time.Duration
	MfaHardwareCheckLifeTime time.Duration
}

func (repo *AuthRequestRepo) Health(ctx context.Context) error {
	if err := repo.UserEvents.Health(ctx); err != nil {
		return err
	}
	return repo.AuthRequests.Health(ctx)
}

func (repo *AuthRequestRepo) CreateAuthRequest(ctx context.Context, request *model.AuthRequest) (*model.AuthRequest, error) {
	err := repo.AuthRequests.SaveAuthRequest(ctx, request)
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
	//query view
	return repo.nextSteps(request)
}

func (repo *AuthRequestRepo) CheckUsername(ctx context.Context, id, username string) error {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return err
	}
	_ = request
	//if request.PasswordChecked() {
	//	return nil, errors.ThrowPreconditionFailed(nil, "EVENT-52NGs", "user already chosen")
	//}
	return errors.ThrowUnimplemented(nil, "EVENT-asjod", "user by username not yet implemented")
}

func (repo *AuthRequestRepo) VerifyPassword(ctx context.Context, id, userID, password string, info *model.BrowserInfo) error {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return err
	}
	//if request.UserID == 0 {
	//
	//}
	//if request.PasswordChecked() {
	//	return nil, errors.ThrowPreconditionFailed(nil, "EVENT-s6Gn3", "password already checked")
	//}
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

func (repo *AuthRequestRepo) nextSteps(request *model.AuthRequest) (*model.AuthRequest, error) {
	if request == nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-ds27a", "request must not be nil")
	}
	if request.UserID == "" {
		return repo.nextStepsNoUserSelected(request)
	}
	//userSession, err := repo.view.GetUserSessionByIDs(request.UserAgentID, request.UserID)
	var userSession *UserSession
	var user *User

	if user.Password == nil {
		request.AddPossibleStep(&model.InitPasswordStep{})
		return request, nil
	}

	if !checkVerificationTime(userSession.PasswordVerification, repo.PasswordCheckLifeTime) {
		request.AddPossibleStep(&model.PasswordStep{})
		return request, nil
	}

	if !repo.mfaChecked(userSession, request, user) {
		return request, nil
	}

	if user.Password.ChangeRequired {
		request.AddPossibleStep(&model.ChangePasswordStep{})
		return request, nil
	}
	if !user.IsEmailVerified {
		request.AddPossibleStep(&model.VerifyEMailStep{})
		return request, nil
	}

	//TODO: consent step
	request.AddPossibleStep(&model.RedirectToCallbackStep{})
	return request, nil
}

func (repo *AuthRequestRepo) nextStepsNoUserSelected(request *model.AuthRequest) (*model.AuthRequest, error) {
	if request.Prompt != model.PromptNone {
		request.AddPossibleStep(&model.LoginStep{})
	}
	//TODO: select account
	return request, nil
}

func (repo *AuthRequestRepo) mfaChecked(userSession *UserSession, request *model.AuthRequest, user *User) bool {
	mfaLevel := request.MfaLevel()
	required := user.MfaMaxSetup < mfaLevel
	if required || repo.mfaNotSkipped(user) {
		request.AddPossibleStep(&model.MfaPromptStep{
			Required:     required,
			MfaProviders: user.MfaTypesSetupPossible(mfaLevel),
		})
		return false
	}
	switch mfaLevel {
	default:
		fallthrough
	case model.MfaLevelSoftware:
		if checkVerificationTime(userSession.MfaSoftwareVerification, repo.MfaSoftwareCheckLifeTime) {
			return true
		}
		//fallthrough?
	case model.MfaLevelHardware:
		if checkVerificationTime(userSession.MfaHardwareVerification, repo.MfaHardwareCheckLifeTime) {
			return true
		}
	}
	request.AddPossibleStep(&model.MfaVerificationStep{
		MfaProviders: user.MfaTypesAllowed(mfaLevel),
	})
	return false
}

func (repo *AuthRequestRepo) mfaNotSkipped(user *User) bool {
	if user.MfaMaxSetup >= 0 {
		return false
	}
	return !checkVerificationTime(user.MfaInitSkipped, repo.MfaInitSkippedLifeTime)
}

func checkVerificationTime(verificationTime time.Time, lifetime time.Duration) bool {
	return verificationTime.Add(lifetime).After(time.Now().UTC())
}

//TODO: into view

type UserSession struct {
	PasswordVerification    time.Time
	MfaVerification         time.Time
	MfaSoftwareVerification time.Time
	MfaHardwareVerification time.Time
}

type User struct {
	MfaInitSkipped  time.Time
	MfaMaxSetup     model.MfaLevel
	Password        *Password
	OTP             interface{}
	IsEmailVerified bool
}

type Password struct {
	ChangeRequired bool
}

func (u *User) MfaTypesSetupPossible(level model.MfaLevel) []model.MfaType {
	types := make([]model.MfaType, 0)
	switch level {
	case model.MfaLevelSoftware:
		if u.OTP == nil {
			types = append(types, model.MfaTypeOTP)
		}
		fallthrough
	case model.MfaLevelHardware:
	}
	return types
}

func (u *User) MfaTypesAllowed(level model.MfaLevel) []model.MfaType {
	types := make([]model.MfaType, 0)
	switch level {
	case model.MfaLevelSoftware:
		if u.OTP != nil {
			types = append(types, model.MfaTypeOTP)
		}
		fallthrough
	case model.MfaLevelHardware:
	}
	return types
}
