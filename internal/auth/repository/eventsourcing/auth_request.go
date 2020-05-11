package eventsourcing

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/auth_request/repository/cache"
	"github.com/caos/zitadel/internal/errors"
	user_model "github.com/caos/zitadel/internal/user/model"
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
	steps, err := repo.nextSteps(request)
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
	return errors.ThrowUnimplemented(nil, "EVENT-asjod", "user by username not yet implemented")
	//check username
	var userID string
	request.UserID = userID
	return repo.AuthRequests.SaveAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) VerifyPassword(ctx context.Context, id, userID, password string, info *model.BrowserInfo) error {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return err
	}
	if request.UserID == userID {
		return errors.ThrowPreconditionFailed(nil, "EVENT-ds35D", "user id does not match request id ")
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

func (repo *AuthRequestRepo) nextSteps(request *model.AuthRequest) ([]model.NextStep, error) {
	if request == nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-ds27a", "request must not be nil")
	}
	steps := make([]model.NextStep, 0)
	if request.UserID == "" {
		if request.Prompt != model.PromptNone {
			steps = append(steps, &model.LoginStep{})
		}
		//TODO: select account
		return steps, nil
	}
	//userSession, err := repo.view.GetUserSessionByIDs(request.UserAgentID, request.UserID)
	var userSession *UserSession
	var user *User

	if user.Password == nil {
		steps = append(steps, &model.InitPasswordStep{})
		return steps, nil
	}

	if !checkVerificationTime(userSession.PasswordVerification, repo.PasswordCheckLifeTime) {
		steps = append(steps, &model.PasswordStep{})
		return steps, nil
	}

	if step, ok := repo.mfaChecked(userSession, request, user); !ok {
		steps = append(steps, step)
		return steps, nil
	}

	if user.Password.ChangeRequired {
		steps = append(steps, &model.ChangePasswordStep{})
	}
	if !user.IsEmailVerified {
		steps = append(steps, &model.VerifyEMailStep{})
	}

	if user.Password.ChangeRequired || !user.IsEmailVerified {
		return steps, nil
	}

	//TODO: consent step
	steps = append(steps, &model.RedirectToCallbackStep{})
	return steps, nil
}

func (repo *AuthRequestRepo) mfaChecked(userSession *UserSession, request *model.AuthRequest, user *User) (model.NextStep, bool) {
	mfaLevel := request.MfaLevel()
	required := user.MfaMaxSetup < mfaLevel
	if required || !repo.mfaSkippedOrSetUp(user) {
		return &model.MfaPromptStep{
			Required:     required,
			MfaProviders: user.MfaTypesSetupPossible(mfaLevel),
		}, false
	}
	switch mfaLevel {
	default:
		fallthrough
	case model.MfaLevelSoftware:
		if checkVerificationTime(userSession.MfaSoftwareVerification, repo.MfaSoftwareCheckLifeTime) {
			return nil, true
		}
		//fallthrough?
	case model.MfaLevelHardware:
		if checkVerificationTime(userSession.MfaHardwareVerification, repo.MfaHardwareCheckLifeTime) {
			return nil, true
		}
	}
	return &model.MfaVerificationStep{
		MfaProviders: user.MfaTypesAllowed(mfaLevel),
	}, false
}

func (repo *AuthRequestRepo) mfaSkippedOrSetUp(user *User) bool {
	if user.MfaMaxSetup >= 0 {
		return true
	}
	return checkVerificationTime(user.MfaInitSkipped, repo.MfaInitSkippedLifeTime)
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
	Password        *user_model.Password
	OTP             *user_model.OTP
	IsEmailVerified bool
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
	default:
		fallthrough
	case model.MfaLevelSoftware:
		if u.OTP != nil {
			types = append(types, model.MfaTypeOTP)
		}
		fallthrough
	case model.MfaLevelHardware:
	}
	return types
}
