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
	org_model "github.com/caos/zitadel/internal/org/model"
	org_view_model "github.com/caos/zitadel/internal/org/repository/view/model"
	user_model "github.com/caos/zitadel/internal/user/model"
	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	user_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	user_view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type mockViewNoUserSession struct{}

func (m *mockViewNoUserSession) UserSessionByIDs(string, string) (*user_view_model.UserSessionView, error) {
	return nil, errors.ThrowNotFound(nil, "id", "user session not found")
}

func (m *mockViewNoUserSession) UserSessionsByAgentID(string) ([]*user_view_model.UserSessionView, error) {
	return nil, nil
}

type mockViewErrUserSession struct{}

func (m *mockViewErrUserSession) UserSessionByIDs(string, string) (*user_view_model.UserSessionView, error) {
	return nil, errors.ThrowInternal(nil, "id", "internal error")
}

func (m *mockViewErrUserSession) UserSessionsByAgentID(string) ([]*user_view_model.UserSessionView, error) {
	return nil, errors.ThrowInternal(nil, "id", "internal error")
}

type mockViewUserSession struct {
	ExternalLoginVerification time.Time
	PasswordVerification      time.Time
	MfaSoftwareVerification   time.Time
	Users                     []mockUser
}

type mockUser struct {
	UserID    string
	LoginName string
}

func (m *mockViewUserSession) UserSessionByIDs(string, string) (*user_view_model.UserSessionView, error) {
	return &user_view_model.UserSessionView{
		ExternalLoginVerification: m.ExternalLoginVerification,
		PasswordVerification:      m.PasswordVerification,
		MfaSoftwareVerification:   m.MfaSoftwareVerification,
	}, nil
}

func (m *mockViewUserSession) UserSessionsByAgentID(string) ([]*user_view_model.UserSessionView, error) {
	sessions := make([]*user_view_model.UserSessionView, len(m.Users))
	for i, user := range m.Users {
		sessions[i] = &user_view_model.UserSessionView{
			UserID:    user.UserID,
			LoginName: user.LoginName,
		}
	}
	return sessions, nil
}

type mockViewNoUser struct{}

func (m *mockViewNoUser) UserByID(string) (*user_view_model.UserView, error) {
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

func (m *mockEventUser) BulkAddExternalIDPs(ctx context.Context, userID string, externalIDPs []*user_model.ExternalIDP) error {
	return nil
}

type mockEventErrUser struct{}

func (m *mockEventErrUser) UserEventsByID(ctx context.Context, id string, sequence uint64) ([]*es_models.Event, error) {
	return nil, errors.ThrowInternal(nil, "id", "internal error")
}

func (m *mockEventErrUser) BulkAddExternalIDPs(ctx context.Context, userID string, externalIDPs []*user_model.ExternalIDP) error {
	return errors.ThrowInternal(nil, "id", "internal error")
}

type mockViewUser struct {
	InitRequired           bool
	PasswordSet            bool
	PasswordChangeRequired bool
	IsEmailVerified        bool
	OTPState               int32
	MfaMaxSetUp            int32
	MfaInitSkipped         time.Time
}

func (m *mockViewUser) UserByID(string) (*user_view_model.UserView, error) {
	return &user_view_model.UserView{
		State:    int32(user_model.UserStateActive),
		UserName: "UserName",
		HumanView: &user_view_model.HumanView{
			FirstName:              "FirstName",
			InitRequired:           m.InitRequired,
			PasswordSet:            m.PasswordSet,
			PasswordChangeRequired: m.PasswordChangeRequired,
			IsEmailVerified:        m.IsEmailVerified,
			OTPState:               m.OTPState,
			MfaMaxSetUp:            m.MfaMaxSetUp,
			MfaInitSkipped:         m.MfaInitSkipped,
		},
	}, nil
}

type mockViewOrg struct {
	State org_model.OrgState
}

func (m *mockViewOrg) OrgByID(string) (*org_view_model.OrgView, error) {
	return &org_view_model.OrgView{
		State: int32(m.State),
	}, nil
}

func (m *mockViewOrg) OrgByPrimaryDomain(string) (*org_view_model.OrgView, error) {
	return &org_view_model.OrgView{
		State: int32(m.State),
	}, nil
}

type mockViewErrOrg struct{}

func (m *mockViewErrOrg) OrgByID(string) (*org_view_model.OrgView, error) {
	return nil, errors.ThrowInternal(nil, "id", "internal error")
}

func (m *mockViewErrOrg) OrgByPrimaryDomain(string) (*org_view_model.OrgView, error) {
	return nil, errors.ThrowInternal(nil, "id", "internal error")
}

func TestAuthRequestRepo_nextSteps(t *testing.T) {
	type fields struct {
		UserEvents                 *user_event.UserEventstore
		AuthRequests               *cache.AuthRequestCache
		View                       *view.View
		userSessionViewProvider    userSessionViewProvider
		userViewProvider           userViewProvider
		userEventProvider          userEventProvider
		orgViewProvider            orgViewProvider
		PasswordCheckLifeTime      time.Duration
		ExternalLoginCheckLifeTime time.Duration
		MfaInitSkippedLifeTime     time.Duration
		MfaSoftwareCheckLifeTime   time.Duration
		MfaHardwareCheckLifeTime   time.Duration
	}
	type args struct {
		request       *model.AuthRequest
		checkLoggedIn bool
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
			args{nil, false},
			nil,
			errors.IsErrorInvalidArgument,
		},
		{
			"prompt none and checkLoggedIn false, callback step",
			fields{},
			args{&model.AuthRequest{Prompt: model.PromptNone}, false},
			[]model.NextStep{&model.RedirectToCallbackStep{}},
			nil,
		},
		{
			"user not set no active session, login step",
			fields{
				userSessionViewProvider: &mockViewNoUserSession{},
			},
			args{&model.AuthRequest{}, false},
			[]model.NextStep{&model.LoginStep{}},
			nil,
		},
		{
			"user not set no active session, linking users, external user not found option",
			fields{
				userSessionViewProvider: &mockViewNoUserSession{},
			},
			args{&model.AuthRequest{LinkingUsers: []*model.ExternalUser{{IDPConfigID: "IDPConfigID", ExternalUserID: "ExternalUserID"}}}, false},
			[]model.NextStep{&model.ExternalNotFoundOptionStep{}},
			nil,
		},
		{
			"user not set, prompt select account and internal error, internal error",
			fields{
				userSessionViewProvider: &mockViewErrUserSession{},
			},
			args{&model.AuthRequest{Prompt: model.PromptSelectAccount}, false},
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
							"loginname1",
						},
						{
							"id2",
							"loginname2",
						},
					},
				},
				userEventProvider: &mockEventUser{},
			},
			args{&model.AuthRequest{Prompt: model.PromptSelectAccount}, false},
			[]model.NextStep{
				&model.LoginStep{},
				&model.SelectUserStep{
					Users: []model.UserSelection{
						{
							UserID:    "id1",
							LoginName: "loginname1",
						},
						{
							UserID:    "id2",
							LoginName: "loginname2",
						},
					},
				}},
			nil,
		},
		{
			"user not set, prompt select account, no active session, login and select account steps",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					Users: nil,
				},
				userEventProvider: &mockEventUser{},
			},
			args{&model.AuthRequest{Prompt: model.PromptSelectAccount}, false},
			[]model.NextStep{
				&model.LoginStep{},
				&model.SelectUserStep{
					Users: []model.UserSelection{},
				}},
			nil,
		},
		{
			"user not found, not found error",
			fields{
				userViewProvider:  &mockViewNoUser{},
				userEventProvider: &mockEventUser{},
			},
			args{&model.AuthRequest{UserID: "UserID"}, false},
			nil,
			errors.IsNotFound,
		},
		{
			"user not active, precondition failed error",
			fields{
				userViewProvider: &mockViewUser{},
				userEventProvider: &mockEventUser{
					&es_models.Event{
						AggregateType: user_es_model.UserAggregate,
						Type:          user_es_model.UserDeactivated,
					},
				},
				orgViewProvider: &mockViewOrg{State: org_model.OrgStateActive},
			},
			args{&model.AuthRequest{UserID: "UserID"}, false},
			nil,
			errors.IsPreconditionFailed,
		},
		{
			"user locked, precondition failed error",
			fields{
				userViewProvider: &mockViewUser{},
				userEventProvider: &mockEventUser{
					&es_models.Event{
						AggregateType: user_es_model.UserAggregate,
						Type:          user_es_model.UserLocked,
					},
				},
				orgViewProvider: &mockViewOrg{State: org_model.OrgStateActive},
			},
			args{&model.AuthRequest{UserID: "UserID"}, false},
			nil,
			errors.IsPreconditionFailed,
		},
		{
			"org error, internal error",
			fields{
				userViewProvider:  &mockViewUser{},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewErrOrg{},
			},
			args{&model.AuthRequest{UserID: "UserID"}, false},
			nil,
			errors.IsInternal,
		},
		{
			"org not active, precondition failed error",
			fields{
				userViewProvider:  &mockViewUser{},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: org_model.OrgStateInactive},
			},
			args{&model.AuthRequest{UserID: "UserID"}, false},
			nil,
			errors.IsPreconditionFailed,
		},
		{
			"usersession not found, new user session, password step",
			fields{
				userSessionViewProvider: &mockViewNoUserSession{},
				userViewProvider: &mockViewUser{
					PasswordSet: true,
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: org_model.OrgStateActive},
			},
			args{&model.AuthRequest{UserID: "UserID"}, false},
			[]model.NextStep{&model.PasswordStep{}},
			nil,
		},
		{
			"usersession error, internal error",
			fields{
				userSessionViewProvider: &mockViewErrUserSession{},
				userViewProvider:        &mockViewUser{},
				userEventProvider:       &mockEventUser{},
				orgViewProvider:         &mockViewOrg{State: org_model.OrgStateActive},
			},
			args{&model.AuthRequest{UserID: "UserID"}, false},
			nil,
			errors.IsInternal,
		},
		{
			"user not initialized, init user step",
			fields{
				userSessionViewProvider: &mockViewUserSession{},
				userViewProvider: &mockViewUser{
					InitRequired: true,
					PasswordSet:  true,
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: org_model.OrgStateActive},
			},
			args{&model.AuthRequest{UserID: "UserID"}, false},
			[]model.NextStep{&model.InitUserStep{
				PasswordSet: true,
			}},
			nil,
		},
		{
			"password not set, init password step",
			fields{
				userSessionViewProvider: &mockViewUserSession{},
				userViewProvider:        &mockViewUser{},
				userEventProvider:       &mockEventUser{},
				orgViewProvider:         &mockViewOrg{State: org_model.OrgStateActive},
			},
			args{&model.AuthRequest{UserID: "UserID"}, false},
			[]model.NextStep{&model.InitPasswordStep{}},
			nil,
		},
		{
			"external user (no external verification), external login step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					MfaSoftwareVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					IsEmailVerified: true,
					MfaMaxSetUp:     int32(model.MfaLevelSoftware),
				},
				userEventProvider:        &mockEventUser{},
				orgViewProvider:          &mockViewOrg{State: org_model.OrgStateActive},
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID", SelectedIDPConfigID: "IDPConfigID"}, false},
			[]model.NextStep{&model.ExternalLoginStep{}},
			nil,
		},
		{
			"external user (external verification set), callback",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					ExternalLoginVerification: time.Now().UTC().Add(-5 * time.Minute),
					MfaSoftwareVerification:   time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					IsEmailVerified: true,
					MfaMaxSetUp:     int32(model.MfaLevelSoftware),
				},
				userEventProvider:          &mockEventUser{},
				orgViewProvider:            &mockViewOrg{State: org_model.OrgStateActive},
				ExternalLoginCheckLifeTime: 10 * 24 * time.Hour,
				MfaSoftwareCheckLifeTime:   18 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID", SelectedIDPConfigID: "IDPConfigID"}, false},
			[]model.NextStep{&model.RedirectToCallbackStep{}},
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
				orgViewProvider:       &mockViewOrg{State: org_model.OrgStateActive},
				PasswordCheckLifeTime: 10 * 24 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID"}, false},
			[]model.NextStep{&model.PasswordStep{}},
			nil,
		},
		{
			"external user (no password check needed), callback",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					MfaSoftwareVerification:   time.Now().UTC().Add(-5 * time.Minute),
					ExternalLoginVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MfaMaxSetUp:     int32(model.MfaLevelSoftware),
				},
				userEventProvider:          &mockEventUser{},
				orgViewProvider:            &mockViewOrg{State: org_model.OrgStateActive},
				MfaSoftwareCheckLifeTime:   18 * time.Hour,
				ExternalLoginCheckLifeTime: 10 * 24 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID", SelectedIDPConfigID: "IDPConfigID"}, false},
			[]model.NextStep{&model.RedirectToCallbackStep{}},
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
					OTPState:    int32(user_model.MfaStateReady),
					MfaMaxSetUp: int32(model.MfaLevelSoftware),
				},
				userEventProvider:        &mockEventUser{},
				orgViewProvider:          &mockViewOrg{State: org_model.OrgStateActive},
				PasswordCheckLifeTime:    10 * 24 * time.Hour,
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID"}, false},
			[]model.NextStep{&model.MfaVerificationStep{
				MfaProviders: []model.MfaType{model.MfaTypeOTP},
			}},
			nil,
		},
		{
			"external user, mfa not verified, mfa check step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:      time.Now().UTC().Add(-5 * time.Minute),
					ExternalLoginVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet: true,
					OTPState:    int32(user_model.MfaStateReady),
					MfaMaxSetUp: int32(model.MfaLevelSoftware),
				},
				userEventProvider:          &mockEventUser{},
				orgViewProvider:            &mockViewOrg{State: org_model.OrgStateActive},
				PasswordCheckLifeTime:      10 * 24 * time.Hour,
				ExternalLoginCheckLifeTime: 10 * 24 * time.Hour,
				MfaSoftwareCheckLifeTime:   18 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID", SelectedIDPConfigID: "IDPConfigID"}, false},
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
				orgViewProvider:          &mockViewOrg{State: org_model.OrgStateActive},
				PasswordCheckLifeTime:    10 * 24 * time.Hour,
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID"}, false},
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
				orgViewProvider:          &mockViewOrg{State: org_model.OrgStateActive},
				PasswordCheckLifeTime:    10 * 24 * time.Hour,
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID"}, false},
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
				orgViewProvider:          &mockViewOrg{State: org_model.OrgStateActive},
				PasswordCheckLifeTime:    10 * 24 * time.Hour,
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID"}, false},
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
				orgViewProvider:          &mockViewOrg{State: org_model.OrgStateActive},
				PasswordCheckLifeTime:    10 * 24 * time.Hour,
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID"}, false},
			[]model.NextStep{&model.RedirectToCallbackStep{}},
			nil,
		},
		{
			"prompt none, checkLoggedIn true and authenticated, redirect to callback step",
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
				orgViewProvider:          &mockViewOrg{State: org_model.OrgStateActive},
				PasswordCheckLifeTime:    10 * 24 * time.Hour,
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
			},
			args{&model.AuthRequest{UserID: "UserID", Prompt: model.PromptNone}, true},
			[]model.NextStep{&model.RedirectToCallbackStep{}},
			nil,
		},
		{
			"linking users, password step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					MfaSoftwareVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MfaMaxSetUp:     int32(model.MfaLevelSoftware),
				},
				userEventProvider:        &mockEventUser{},
				orgViewProvider:          &mockViewOrg{State: org_model.OrgStateActive},
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
			},
			args{
				&model.AuthRequest{
					UserID:              "UserID",
					SelectedIDPConfigID: "IDPConfigID",
					LinkingUsers:        []*model.ExternalUser{{IDPConfigID: "IDPConfigID", ExternalUserID: "UserID", DisplayName: "DisplayName"}},
				}, false},
			[]model.NextStep{&model.PasswordStep{}},
			nil,
		},
		{
			"linking users, linking step",
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
				orgViewProvider:          &mockViewOrg{State: org_model.OrgStateActive},
				MfaSoftwareCheckLifeTime: 18 * time.Hour,
				PasswordCheckLifeTime:    10 * 24 * time.Hour,
			},
			args{
				&model.AuthRequest{
					UserID:              "UserID",
					SelectedIDPConfigID: "IDPConfigID",
					LinkingUsers:        []*model.ExternalUser{{IDPConfigID: "IDPConfigID", ExternalUserID: "UserID", DisplayName: "DisplayName"}},
				}, false},
			[]model.NextStep{&model.LinkUsersStep{}},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AuthRequestRepo{
				UserEvents:                 tt.fields.UserEvents,
				AuthRequests:               tt.fields.AuthRequests,
				View:                       tt.fields.View,
				UserSessionViewProvider:    tt.fields.userSessionViewProvider,
				UserViewProvider:           tt.fields.userViewProvider,
				UserEventProvider:          tt.fields.userEventProvider,
				OrgViewProvider:            tt.fields.orgViewProvider,
				PasswordCheckLifeTime:      tt.fields.PasswordCheckLifeTime,
				ExternalLoginCheckLifeTime: tt.fields.ExternalLoginCheckLifeTime,
				MfaInitSkippedLifeTime:     tt.fields.MfaInitSkippedLifeTime,
				MfaSoftwareCheckLifeTime:   tt.fields.MfaSoftwareCheckLifeTime,
				MfaHardwareCheckLifeTime:   tt.fields.MfaHardwareCheckLifeTime,
			}
			got, err := repo.nextSteps(context.Background(), tt.args.request, tt.args.checkLoggedIn)
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
		//			OTPState: user_model.MfaStateReady,
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
					HumanView: &user_model.HumanView{
						MfaMaxSetUp: model.MfaLevelNotSetUp,
					},
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
					HumanView: &user_model.HumanView{
						MfaMaxSetUp:    model.MfaLevelNotSetUp,
						MfaInitSkipped: time.Now().UTC(),
					},
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
					HumanView: &user_model.HumanView{
						MfaMaxSetUp: model.MfaLevelSoftware,
						OTPState:    user_model.MfaStateReady,
					},
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
					HumanView: &user_model.HumanView{
						MfaMaxSetUp: model.MfaLevelSoftware,
						OTPState:    user_model.MfaStateReady,
					},
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
				HumanView: &user_model.HumanView{
					MfaMaxSetUp: model.MfaLevelSoftware,
				},
			}},
			true,
		},
		{
			"mfa skipped active, true",
			fields{
				MfaInitSkippedLifeTime: 30 * 24 * time.Hour,
			},
			args{&user_model.UserView{
				HumanView: &user_model.HumanView{
					MfaMaxSetUp:    -1,
					MfaInitSkipped: time.Now().UTC().Add(-10 * time.Hour),
				},
			}},
			true,
		},
		{
			"mfa skipped inactive, false",
			fields{
				MfaInitSkippedLifeTime: 30 * 24 * time.Hour,
			},
			args{&user_model.UserView{
				HumanView: &user_model.HumanView{
					MfaMaxSetUp:    -1,
					MfaInitSkipped: time.Now().UTC().Add(-40 * 24 * time.Hour),
				},
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
		user          *user_model.UserView
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
				user:          &user_model.UserView{ID: "id"},
			},
			&user_model.UserSessionView{},
			nil,
		},
		{
			"internal error, internal error",
			args{
				userProvider: &mockViewErrUserSession{},
				user:         &user_model.UserView{ID: "id"},
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
				user:          &user_model.UserView{ID: "id", HumanView: &user_model.HumanView{FirstName: "FirstName"}},
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
			"new user events but error, old view model state",
			args{
				userProvider: &mockViewUserSession{
					PasswordVerification: time.Now().UTC().Round(1 * time.Second),
				},
				agentID: "agentID",
				user:    &user_model.UserView{ID: "id", HumanView: &user_model.HumanView{FirstName: "FirstName"}},
				eventProvider: &mockEventUser{
					&es_models.Event{
						AggregateType: user_es_model.UserAggregate,
						Type:          user_es_model.MFAOTPCheckSucceeded,
						CreationDate:  time.Now().UTC().Round(1 * time.Second),
					},
				},
			},
			&user_model.UserSessionView{
				PasswordVerification:    time.Now().UTC().Round(1 * time.Second),
				MfaSoftwareVerification: time.Time{},
				MfaHardwareVerification: time.Time{},
			},
			nil,
		},
		{
			"new user events but other agentID, old view model state",
			args{
				userProvider: &mockViewUserSession{
					PasswordVerification: time.Now().UTC().Round(1 * time.Second),
				},
				agentID: "agentID",
				user:    &user_model.UserView{ID: "id"},
				eventProvider: &mockEventUser{
					&es_models.Event{
						AggregateType: user_es_model.UserAggregate,
						Type:          user_es_model.MFAOTPCheckSucceeded,
						CreationDate:  time.Now().UTC().Round(1 * time.Second),
						Data: func() []byte {
							data, _ := json.Marshal(&user_es_model.AuthRequest{UserAgentID: "otherID"})
							return data
						}(),
					},
				},
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
				agentID: "agentID",
				user:    &user_model.UserView{ID: "id", HumanView: &user_model.HumanView{FirstName: "FirstName"}},
				eventProvider: &mockEventUser{
					&es_models.Event{
						AggregateType: user_es_model.UserAggregate,
						Type:          user_es_model.MFAOTPCheckSucceeded,
						CreationDate:  time.Now().UTC().Round(1 * time.Second),
						Data: func() []byte {
							data, _ := json.Marshal(&user_es_model.AuthRequest{UserAgentID: "agentID"})
							return data
						}(),
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
		{
			"new user events (user deleted), precondition failed error",
			args{
				userProvider: &mockViewUserSession{
					PasswordVerification: time.Now().UTC().Round(1 * time.Second),
				},
				agentID: "agentID",
				user:    &user_model.UserView{ID: "id"},
				eventProvider: &mockEventUser{
					&es_models.Event{
						AggregateType: user_es_model.UserAggregate,
						Type:          user_es_model.UserRemoved,
					},
				},
			},
			nil,
			errors.IsPreconditionFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userSessionByIDs(context.Background(), tt.args.userProvider, tt.args.eventProvider, tt.args.agentID, tt.args.user)
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
				State:    user_model.UserStateActive,
				UserName: "UserName",
				HumanView: &user_model.HumanView{
					PasswordChangeRequired: true,
					FirstName:              "FirstName",
				},
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
				State:    user_model.UserStateActive,
				UserName: "UserName",
				HumanView: &user_model.HumanView{
					PasswordChangeRequired: true,
					FirstName:              "FirstName",
				},
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
				ChangeDate: time.Now().UTC().Round(1 * time.Second),
				State:      user_model.UserStateActive,
				UserName:   "UserName",
				HumanView: &user_model.HumanView{
					PasswordChangeRequired: false,
					PasswordChanged:        time.Now().UTC().Round(1 * time.Second),
					FirstName:              "FirstName",
				},
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
