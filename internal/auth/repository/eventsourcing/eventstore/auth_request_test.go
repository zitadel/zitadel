package eventstore

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/auth_request/repository/cache"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	user_model "github.com/caos/zitadel/internal/user/model"
	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	user_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type mockViewNoUserSession struct{}

func (m *mockViewNoUserSession) UserSessionByIDs(string, string) (*view_model.UserSessionView, error) {
	return nil, errors.ThrowNotFound(nil, "id", "user session not found")
}

func (m *mockViewNoUserSession) UserSessionsByAgentID(string) ([]*view_model.UserSessionView, error) {
	return nil, errors.ThrowInternal(nil, "id", "internal error")
}

type mockViewErrUserSession struct{}

func (m *mockViewErrUserSession) UserSessionByIDs(string, string) (*view_model.UserSessionView, error) {
	return nil, errors.ThrowInternal(nil, "id", "internal error")
}

func (m *mockViewErrUserSession) UserSessionsByAgentID(string) ([]*view_model.UserSessionView, error) {
	return nil, errors.ThrowInternal(nil, "id", "internal error")
}

type mockViewUserSession struct {
	PasswordVerification    time.Time
	MfaSoftwareVerification time.Time
	Users                   []mockUser
}

type mockUser struct {
	UserID   string
	UserName string
}

func (m *mockViewUserSession) UserSessionByIDs(string, string) (*view_model.UserSessionView, error) {
	return &view_model.UserSessionView{
		PasswordVerification:    m.PasswordVerification,
		MfaSoftwareVerification: m.MfaSoftwareVerification,
	}, nil
}

func (m *mockViewUserSession) UserSessionsByAgentID(string) ([]*view_model.UserSessionView, error) {
	sessions := make([]*view_model.UserSessionView, len(m.Users))
	for i, user := range m.Users {
		sessions[i] = &view_model.UserSessionView{
			UserID:   user.UserID,
			UserName: user.UserName,
		}
	}
	return sessions, nil
}

type mockViewNoUser struct{}

func (m *mockViewNoUser) UserByID(string) (*view_model.UserView, error) {
	return nil, errors.ThrowNotFound(nil, "id", "user not found")
}

type mockEventUser struct {
	Event *es_models.Event
}

func (m *mockEventUser) UserEventsByID(ctx context.Context, id string, sequence uint64) ([]*es_models.Event, error) {
	events := make([]*es_models.Event, 0)
	if m.Event != nil {
		events = append(events, m.Event)
	}
	return events, nil
}

type mockEventErrUser struct{}

func (m *mockEventErrUser) UserEventsByID(ctx context.Context, id string, sequence uint64) ([]*es_models.Event, error) {
	return nil, errors.ThrowInternal(nil, "id", "internal error")
}

type mockViewUser struct {
	PasswordSet            bool
	PasswordChangeRequired bool
	IsEmailVerified        bool
	OTPState               int32
	MfaMaxSetUp            int32
	MfaInitSkipped         time.Time
}

func (m *mockViewUser) UserByID(string) (*view_model.UserView, error) {
	return &view_model.UserView{
		PasswordSet:            m.PasswordSet,
		PasswordChangeRequired: m.PasswordChangeRequired,
		IsEmailVerified:        m.IsEmailVerified,
		OTPState:               m.OTPState,
		MfaMaxSetUp:            m.MfaMaxSetUp,
		MfaInitSkipped:         m.MfaInitSkipped,
	}, nil
}

func TestAuthRequestRepo_nextSteps(t *testing.T) {
	type fields struct {
		UserEvents               *user_event.UserEventstore
		AuthRequests             *cache.AuthRequestCache
		View                     *view.View
		userSessionViewProvider  userSessionViewProvider
		userViewProvider         userViewProvider
		userEventProvider        userEventProvider
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
		{
			"user not set and prompt none, no step",
			fields{},
			args{&model.AuthRequest{Prompt: model.PromptNone}},
			[]model.NextStep{},
			nil,
		},
		{
			"user not set, prompt select account and internal error, internal error",
			fields{
				userSessionViewProvider: &mockViewNoUserSession{},
			},
			args{&model.AuthRequest{Prompt: model.PromptSelectAccount}},
			nil,
			errors.IsInternal,
		},
		{
			"user not set, prompt select account, login and select account steps",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					Users: []mockUser{
						{
							"id1",
							"username1",
						},
						{
							"id2",
							"username2",
						},
					},
				},
				userEventProvider: &mockEventUser{},
			},
			args{&model.AuthRequest{Prompt: model.PromptSelectAccount}},
			[]model.NextStep{
				&model.LoginStep{},
				&model.SelectUserStep{
					Users: []model.UserSelection{
						{
							UserID:   "id1",
							UserName: "username1",
						},
						{
							UserID:   "id2",
							UserName: "username2",
						},
					},
				}},
			nil,
		},
		{
			"usersession not found, new user session, password step",
			fields{
				userSessionViewProvider: &mockViewNoUserSession{},
				userViewProvider: &mockViewUser{
					PasswordSet: true,
				},
				userEventProvider: &mockEventUser{},
			},
			args{&model.AuthRequest{UserID: "UserID"}},
			[]model.NextStep{&model.PasswordStep{}},
			nil,
		},
		{
			"usersession error, internal error",
			fields{
				userSessionViewProvider: &mockViewErrUserSession{},
				userEventProvider:       &mockEventUser{},
			},
			args{&model.AuthRequest{UserID: "UserID"}},
			nil,
			errors.IsInternal,
		},
		{
			"user not not found, not found error",
			fields{
				userSessionViewProvider: &mockViewUserSession{},
				userViewProvider:        &mockViewNoUser{},
				userEventProvider:       &mockEventUser{},
			},
			args{&model.AuthRequest{UserID: "UserID"}},
			nil,
			errors.IsNotFound,
		},
		{
			"password not set, init password step",
			fields{
				userSessionViewProvider: &mockViewUserSession{},
				userViewProvider:        &mockViewUser{},
				userEventProvider:       &mockEventUser{},
			},
			args{&model.AuthRequest{UserID: "UserID"}},
			[]model.NextStep{&model.InitPasswordStep{}},
			nil,
		},
		{
			"password not verified, password check step",
			fields{
				userSessionViewProvider: &mockViewUserSession{},
				userViewProvider: &mockViewUser{
					PasswordSet: true,
				},
				userEventProvider:     &mockEventUser{},
				PasswordCheckLifeTime: 10 * 24 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID"}},
			[]model.NextStep{&model.PasswordStep{}},
			nil,
		},
		{
			"mfa not verified, mfa check step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet: true,
					OTPState:    int32(user_model.MFASTATE_READY),
					MfaMaxSetUp: int32(model.MfaLevelSoftware),
				},
				userEventProvider:        &mockEventUser{},
				PasswordCheckLifeTime:    10 * 24 * time.Hour,
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID"}},
			[]model.NextStep{&model.MfaVerificationStep{
				MfaProviders: []model.MfaType{model.MfaTypeOTP},
			}},
			nil,
		},
		{
			"password change required and email verified, password change step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:    time.Now().UTC().Add(-5 * time.Minute),
					MfaSoftwareVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:            true,
					PasswordChangeRequired: true,
					IsEmailVerified:        true,
					MfaMaxSetUp:            int32(model.MfaLevelSoftware),
				},
				userEventProvider:        &mockEventUser{},
				PasswordCheckLifeTime:    10 * 24 * time.Hour,
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID"}},
			[]model.NextStep{&model.ChangePasswordStep{}},
			nil,
		},
		{
			"email not verified and no password change required, mail verification step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:    time.Now().UTC().Add(-5 * time.Minute),
					MfaSoftwareVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet: true,
					MfaMaxSetUp: int32(model.MfaLevelSoftware),
				},
				userEventProvider:        &mockEventUser{},
				PasswordCheckLifeTime:    10 * 24 * time.Hour,
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID"}},
			[]model.NextStep{&model.VerifyEMailStep{}},
			nil,
		},
		{
			"email not verified and password change required, mail verification step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:    time.Now().UTC().Add(-5 * time.Minute),
					MfaSoftwareVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:            true,
					PasswordChangeRequired: true,
					MfaMaxSetUp:            int32(model.MfaLevelSoftware),
				},
				userEventProvider:        &mockEventUser{},
				PasswordCheckLifeTime:    10 * 24 * time.Hour,
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID"}},
			[]model.NextStep{&model.ChangePasswordStep{}, &model.VerifyEMailStep{}},
			nil,
		},
		{
			"email verified and no password change required, redirect to callback step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:    time.Now().UTC().Add(-5 * time.Minute),
					MfaSoftwareVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MfaMaxSetUp:     int32(model.MfaLevelSoftware),
				},
				userEventProvider:        &mockEventUser{},
				PasswordCheckLifeTime:    10 * 24 * time.Hour,
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID"}},
			[]model.NextStep{&model.RedirectToCallbackStep{}},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AuthRequestRepo{
				UserEvents:               tt.fields.UserEvents,
				AuthRequests:             tt.fields.AuthRequests,
				View:                     tt.fields.View,
				UserSessionViewProvider:  tt.fields.userSessionViewProvider,
				UserViewProvider:         tt.fields.userViewProvider,
				UserEventProvider:        tt.fields.userEventProvider,
				PasswordCheckLifeTime:    tt.fields.PasswordCheckLifeTime,
				MfaInitSkippedLifeTime:   tt.fields.MfaInitSkippedLifeTime,
				MfaSoftwareCheckLifeTime: tt.fields.MfaSoftwareCheckLifeTime,
				MfaHardwareCheckLifeTime: tt.fields.MfaHardwareCheckLifeTime,
			}
			got, err := repo.nextSteps(context.Background(), tt.args.request)
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
		userSession *user_model.UserSessionView
		request     *model.AuthRequest
		user        *user_model.UserView
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
		//		user: &user_model.UserView{
		//			OTPState: user_model.MFASTATE_READY,
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
				user: &user_model.UserView{
					MfaMaxSetUp: model.MfaLevelNotSetUp,
				},
			},
			&model.MfaPromptStep{
				MfaProviders: []model.MfaType{
					model.MfaTypeOTP,
				},
			},
			false,
		},
		{
			"not set up and skipped, true",
			fields{
				MfaInitSkippedLifeTime: 30 * 24 * time.Hour,
			},
			args{
				request: &model.AuthRequest{},
				user: &user_model.UserView{
					MfaMaxSetUp:    model.MfaLevelNotSetUp,
					MfaInitSkipped: time.Now().UTC(),
				},
			},
			nil,
			true,
		},
		{
			"checked mfa software, true",
			fields{
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
			},
			args{
				request: &model.AuthRequest{},
				user: &user_model.UserView{
					MfaMaxSetUp: model.MfaLevelSoftware,
					OTPState:    user_model.MFASTATE_READY,
				},
				userSession: &user_model.UserSessionView{MfaSoftwareVerification: time.Now().UTC().Add(-5 * time.Hour)},
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
				user: &user_model.UserView{
					MfaMaxSetUp: model.MfaLevelSoftware,
					OTPState:    user_model.MFASTATE_READY,
				},
				userSession: &user_model.UserSessionView{},
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
		user *user_model.UserView
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
			args{&user_model.UserView{
				MfaMaxSetUp: model.MfaLevelSoftware,
			}},
			true,
		},
		{
			"mfa skipped active, true",
			fields{
				MfaInitSkippedLifeTime: 30 * 24 * time.Hour,
			},
			args{&user_model.UserView{
				MfaMaxSetUp:    -1,
				MfaInitSkipped: time.Now().UTC().Add(-10 * time.Hour),
			}},
			true,
		},
		{
			"mfa skipped inactive, false",
			fields{
				MfaInitSkippedLifeTime: 30 * 24 * time.Hour,
			},
			args{&user_model.UserView{
				MfaMaxSetUp:    -1,
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

func Test_userSessionByIDs(t *testing.T) {
	type args struct {
		userProvider  userSessionViewProvider
		eventProvider userEventProvider
		agentID       string
		userID        string
	}
	tests := []struct {
		name    string
		args    args
		want    *user_model.UserSessionView
		wantErr func(error) bool
	}{
		{
			"not found, new session",
			args{
				userProvider:  &mockViewNoUserSession{},
				eventProvider: &mockEventErrUser{},
			},
			&user_model.UserSessionView{},
			nil,
		},
		{
			"internal error, internal error",
			args{
				userProvider: &mockViewErrUserSession{},
			},
			nil,
			errors.IsInternal,
		},
		{
			"error user events, old view model state",
			args{
				userProvider: &mockViewUserSession{
					PasswordVerification: time.Now().UTC().Round(1 * time.Second),
				},
				eventProvider: &mockEventErrUser{},
			},
			&user_model.UserSessionView{
				PasswordVerification:    time.Now().UTC().Round(1 * time.Second),
				MfaSoftwareVerification: time.Time{},
				MfaHardwareVerification: time.Time{},
			},
			nil,
		},
		{
			"new user events, new view model state",
			args{
				userProvider: &mockViewUserSession{
					PasswordVerification: time.Now().UTC().Round(1 * time.Second),
				},
				eventProvider: &mockEventUser{
					&es_models.Event{
						AggregateType: user_es_model.UserAggregate,
						Type:          user_es_model.MfaOtpCheckSucceeded,
						CreationDate:  time.Now().UTC().Round(1 * time.Second),
					},
				},
			},
			&user_model.UserSessionView{
				PasswordVerification:    time.Now().UTC().Round(1 * time.Second),
				MfaSoftwareVerification: time.Now().UTC().Round(1 * time.Second),
				ChangeDate:              time.Now().UTC().Round(1 * time.Second),
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userSessionByIDs(context.Background(), tt.args.userProvider, tt.args.eventProvider, tt.args.agentID, tt.args.userID)
			if (err != nil && tt.wantErr == nil) || (tt.wantErr != nil && !tt.wantErr(err)) {
				t.Errorf("nextSteps() wrong error = %v", err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_userByID(t *testing.T) {
	type args struct {
		ctx           context.Context
		viewProvider  userViewProvider
		eventProvider userEventProvider
		userID        string
	}
	tests := []struct {
		name    string
		args    args
		want    *user_model.UserView
		wantErr func(error) bool
	}{

		{
			"not found, not found error",
			args{
				viewProvider: &mockViewNoUser{},
			},
			nil,
			errors.IsNotFound,
		},
		{
			"error user events, old view model state",
			args{
				viewProvider: &mockViewUser{
					PasswordChangeRequired: true,
				},
				eventProvider: &mockEventErrUser{},
			},
			&user_model.UserView{
				PasswordChangeRequired: true,
			},
			nil,
		},
		{
			"new user events but error, old view model state",
			args{
				viewProvider: &mockViewUser{
					PasswordChangeRequired: true,
				},
				eventProvider: &mockEventUser{
					&es_models.Event{
						AggregateType: user_es_model.UserAggregate,
						Type:          user_es_model.UserPasswordChanged,
						CreationDate:  time.Now().UTC().Round(1 * time.Second),
						Data:          nil,
					},
				},
			},
			&user_model.UserView{
				PasswordChangeRequired: true,
			},
			nil,
		},
		{
			"new user events, new view model state",
			args{
				viewProvider: &mockViewUser{
					PasswordChangeRequired: true,
				},
				eventProvider: &mockEventUser{
					&es_models.Event{
						AggregateType: user_es_model.UserAggregate,
						Type:          user_es_model.UserPasswordChanged,
						CreationDate:  time.Now().UTC().Round(1 * time.Second),
						Data: func() []byte {
							data, _ := json.Marshal(user_es_model.Password{ChangeRequired: false})
							return data
						}(),
					},
				},
			},
			&user_model.UserView{
				PasswordChangeRequired: false,
				ChangeDate:             time.Now().UTC().Round(1 * time.Second),
				State:                  user_model.USERSTATE_INITIAL,
				PasswordChanged:        time.Now().UTC().Round(1 * time.Second),
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userByID(tt.args.ctx, tt.args.viewProvider, tt.args.eventProvider, tt.args.userID)
			if (err != nil && tt.wantErr == nil) || (tt.wantErr != nil && !tt.wantErr(err)) {
				t.Errorf("nextSteps() wrong error = %v", err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
