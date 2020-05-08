package eventsourcing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/auth_request/repository/cache"
	"github.com/caos/zitadel/internal/errors"
	user_model "github.com/caos/zitadel/internal/user/model"
	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

func TestAuthRequestRepo_nextSteps(t *testing.T) {
	type fields struct {
		UserEvents               *user_event.UserEventstore
		AuthRequests             *cache.AuthRequestCache
		PasswordCheckLifeTime    time.Duration
		MfaInitSkippedLifeTime   time.Duration
		MfaSoftwareCheckLifeTime time.Duration
		MfaHardwareCheckLifeTime time.Duration
	}
	type args struct {
		request *model.AuthRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.NextStep
		wantErr func(error) bool
	}{
		{
			"request nil, error",
			fields{},
			args{nil},
			nil,
			errors.IsErrorInvalidArgument,
		},
		{
			"user not set, login step",
			fields{},
			args{&model.AuthRequest{}},
			[]model.NextStep{&model.LoginStep{}},
			nil,
		},
		//{ //TODO: view
		//	"password not set, init password step",
		//	fields{},
		//	args{&model.AuthRequest{UserID: "UserID"}},
		//	[]model.NextStep{&model.InitPasswordStep{}},
		//	nil,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AuthRequestRepo{
				UserEvents:               tt.fields.UserEvents,
				AuthRequests:             tt.fields.AuthRequests,
				PasswordCheckLifeTime:    tt.fields.PasswordCheckLifeTime,
				MfaInitSkippedLifeTime:   tt.fields.MfaInitSkippedLifeTime,
				MfaSoftwareCheckLifeTime: tt.fields.MfaSoftwareCheckLifeTime,
				MfaHardwareCheckLifeTime: tt.fields.MfaHardwareCheckLifeTime,
			}
			got, err := repo.nextSteps(tt.args.request)
			if (err != nil && tt.wantErr == nil) || (tt.wantErr != nil && !tt.wantErr(err)) {
				t.Errorf("nextSteps() wrong error = %v", err)
				return
			}
			assert.ElementsMatch(t, got, tt.want)
		})
	}
}

func TestAuthRequestRepo_mfaChecked(t *testing.T) {
	type fields struct {
		MfaInitSkippedLifeTime   time.Duration
		MfaSoftwareCheckLifeTime time.Duration
		MfaHardwareCheckLifeTime time.Duration
	}
	type args struct {
		userSession *UserSession
		request     *model.AuthRequest
		user        *User
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        model.NextStep
		wantChecked bool
	}{
		//{
		//	"required, prompt and false", //TODO: enable when LevelsOfAssurance is checked
		//	fields{},
		//	args{
		//		request: &model.AuthRequest{PossibleLOAs: []model.LevelOfAssurance{}},
		//		user: &User{
		//			OTP: nil,
		//		},
		//	},
		//	false,
		//},
		{
			"not set up, prompt and false",
			fields{
				MfaInitSkippedLifeTime: 30 * 24 * time.Hour,
			},
			args{
				request: &model.AuthRequest{},
				user: &User{
					MfaMaxSetup: -1,
				},
			},
			&model.MfaPromptStep{
				MfaProviders: []model.MfaType{},
			},
			false,
		},
		{
			"checked mfa software, true",
			fields{
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
			},
			args{
				request: &model.AuthRequest{},
				user: &User{
					OTP: &user_model.OTP{State: user_model.MFASTATE_READY},
				},
				userSession: &UserSession{MfaSoftwareVerification: time.Now().UTC().Add(-5 * time.Hour)},
			},
			nil,
			true,
		},
		{
			"not checked, check and false",
			fields{
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
			},
			args{
				request: &model.AuthRequest{},
				user: &User{
					OTP: &user_model.OTP{State: user_model.MFASTATE_READY},
				},
				userSession: &UserSession{},
			},

			&model.MfaVerificationStep{
				MfaProviders: []model.MfaType{model.MfaTypeOTP},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AuthRequestRepo{
				MfaInitSkippedLifeTime:   tt.fields.MfaInitSkippedLifeTime,
				MfaSoftwareCheckLifeTime: tt.fields.MfaSoftwareCheckLifeTime,
				MfaHardwareCheckLifeTime: tt.fields.MfaHardwareCheckLifeTime,
			}
			got, ok := repo.mfaChecked(tt.args.userSession, tt.args.request, tt.args.user)
			if ok != tt.wantChecked {
				t.Errorf("mfaChecked() checked = %v, want %v", ok, tt.wantChecked)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAuthRequestRepo_mfaSkippedOrSetUp(t *testing.T) {
	type fields struct {
		MfaInitSkippedLifeTime time.Duration
	}
	type args struct {
		user *User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"mfa set up, true",
			fields{},
			args{&User{
				MfaMaxSetup: model.MfaLevelSoftware,
			}},
			true,
		},
		{
			"mfa skipped active, true",
			fields{
				MfaInitSkippedLifeTime: 30 * 24 * time.Hour,
			},
			args{&User{
				MfaMaxSetup:    -1,
				MfaInitSkipped: time.Now().UTC().Add(-10 * time.Hour),
			}},
			true,
		},
		{
			"mfa skipped inactive, false",
			fields{
				MfaInitSkippedLifeTime: 30 * 24 * time.Hour,
			},
			args{&User{
				MfaMaxSetup:    -1,
				MfaInitSkipped: time.Now().UTC().Add(-40 * 24 * time.Hour),
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AuthRequestRepo{
				MfaInitSkippedLifeTime: tt.fields.MfaInitSkippedLifeTime,
			}
			if got := repo.mfaSkippedOrSetUp(tt.args.user); got != tt.want {
				t.Errorf("mfaSkippedOrSetUp() = %v, want %v", got, tt.want)
			}
		})
	}
}
