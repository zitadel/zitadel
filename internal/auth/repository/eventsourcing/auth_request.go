package eventsourcing

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/auth_request_cache"
	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type AuthRequestRepo struct {
	UserEvents   *user_event.UserEventstore
	AuthRequests *cache.AuthRequestCache
	//view      *view.View
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
	return nextSteps(request, nil)
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
	return nil, errors.ThrowUnimplemented(nil, "EVENT-asjod", "user by username not yet implemented")
	//if err != nil {
	//	return nextStepsNoUserSelected(request, true)
	//}
	//return nextSteps(request, user)
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
	authReq := &model.AuthRequest{
		BrowserInfo: info,
	}
	return repo.UserEvents.CheckPassword(ctx, userID, password, authReq)
}

func (repo *AuthRequestRepo) SkipMfaInit(ctx context.Context, authRequestID, userID string) error {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, authRequestID)
	if err != nil {
		return err
	}
	return repo.UserEvents.SkipMfaInit(ctx, userID)
}

func (repo *AuthRequestRepo) VerifyMfaOTP(ctx context.Context, authRequestID, userID string, code string, info *model.BrowserInfo) error {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, authRequestID)
	if err != nil {
		return err
	}
	if request.UserID != userID {
		return errors.ThrowPreconditionFailed(nil, "EVENT-ADJ26", "user id does not match request id")
	}
	return repo.UserEvents.CheckMfaOTP(ctx, userID, code)
}

type UserSession struct {
	PasswordVerification    time.Time
	MfaVerification         time.Time
	MfaSoftwareVerification time.Time
	MfaHardwareVerification time.Time
}

func nextStepsNoUserSelected(request *model.AuthRequest, notFound bool) (*model.AuthRequest, error) {
	if request.Prompt != model.PromptNone {
		request.AddPossibleStep(&model.LoginStep{NotFound: notFound})
	}
	//TODO: select account
	return request, nil
}

func (repo *AuthRequestRepo) nextSteps(request *model.AuthRequest) (*model.AuthRequest, error) {
	if request == nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-ds27a", "request must not be nil")
	}
	if request.UserID == "" {
		return nextStepsNoUserSelected(request, false)
	}
	//userSession, err := repo.view.GetUserSessionByIDs(request.UserAgentID, request.UserID)
	var userSession *UserSession
	var user *User

	if user.Password == nil {
		request.AddPossibleStep(&model.InitPasswordStep{})
		return request, nil
	}

	PasswordCheckLifeTime := 30 * 24 * time.Hour
	if !checkVerificationTime(userSession.PasswordVerification, PasswordCheckLifeTime) {
		request.AddPossibleStep(&model.PasswordStep{})
		return request, nil
	}

	if !mfaChecked(userSession, request, user) {
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

func checkVerificationTime(verificationTime time.Time, lifetime time.Duration) bool {
	return verificationTime.Add(lifetime).After(time.Now().UTC())
}

type User struct {
	MfaInitSkipped  time.Time
	MfaMaxSetup     model.MfaLevel
	Password        *Password
	OTP             interface{}
	IsEmailVerified bool
}

type Password struct {
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
	case model.MfaLevelSoftware:
		if u.OTP != nil {
			types = append(types, model.MfaTypeOTP)
		}
		fallthrough
	case model.MfaLevelHardware:
	}
	return types
}

const MfaInitSkippedLifeTime = 30 * 24 * time.Hour
const MfaSoftwareCheckLifeTime = 18 * time.Hour
const MfaHardwareCheckLifeTime = 12 * time.Hour

func mfaChecked(userSession *UserSession, request *model.AuthRequest, user *User) bool {
	mfaLevel := request.MfaLevel()
	required := user.MfaMaxSetup < mfaLevel
	if required || MfaNotSkipped(user) {
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
		if checkVerificationTime(userSession.MfaSoftwareVerification, MfaSoftwareCheckLifeTime) {
			return true
		}
		//fallthrough?
	case model.MfaLevelHardware:
		if checkVerificationTime(userSession.MfaHardwareVerification, MfaHardwareCheckLifeTime) {
			return true
		}
	}
	request.AddPossibleStep(&model.MfaVerificationStep{
		MfaProviders: user.MfaTypesAllowed(mfaLevel),
	})
	return false
}

func MfaNotSkipped(user *User) bool {
	if user.MfaMaxSetup >= 0 {
		return false
	}
	return !checkVerificationTime(user.MfaInitSkipped, MfaInitSkippedLifeTime)
}
