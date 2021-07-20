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
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_view_model "github.com/caos/zitadel/internal/org/repository/view/model"
	proj_view_model "github.com/caos/zitadel/internal/project/repository/view/model"
	user_model "github.com/caos/zitadel/internal/user/model"
	user_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	user_view_model "github.com/caos/zitadel/internal/user/repository/view/model"
	grant_view_model "github.com/caos/zitadel/internal/usergrant/repository/view/model"
)

type mockViewNoUserSession struct{}

func (m *mockViewNoUserSession) UserSessionByIDs(string, string) (*user_view_model.UserSessionView, error) {
	return nil, errors.ThrowNotFound(nil, "id", "user session not found")
}

func (m *mockViewNoUserSession) UserSessionsByAgentID(string) ([]*user_view_model.UserSessionView, error) {
	return nil, nil
}

func (m *mockViewNoUserSession) PrefixAvatarURL() string {
	return ""
}

type mockViewErrUserSession struct{}

func (m *mockViewErrUserSession) UserSessionByIDs(string, string) (*user_view_model.UserSessionView, error) {
	return nil, errors.ThrowInternal(nil, "id", "internal error")
}

func (m *mockViewErrUserSession) UserSessionsByAgentID(string) ([]*user_view_model.UserSessionView, error) {
	return nil, errors.ThrowInternal(nil, "id", "internal error")
}

func (m *mockViewErrUserSession) PrefixAvatarURL() string {
	return ""
}

type mockViewUserSession struct {
	ExternalLoginVerification time.Time
	PasswordlessVerification  time.Time
	PasswordVerification      time.Time
	SecondFactorVerification  time.Time
	MultiFactorVerification   time.Time
	Users                     []mockUser
}

type mockUser struct {
	UserID        string
	LoginName     string
	ResourceOwner string
}

func (m *mockViewUserSession) UserSessionByIDs(string, string) (*user_view_model.UserSessionView, error) {
	return &user_view_model.UserSessionView{
		ExternalLoginVerification: m.ExternalLoginVerification,
		PasswordlessVerification:  m.PasswordlessVerification,
		PasswordVerification:      m.PasswordVerification,
		SecondFactorVerification:  m.SecondFactorVerification,
		MultiFactorVerification:   m.MultiFactorVerification,
	}, nil
}

func (m *mockViewUserSession) UserSessionsByAgentID(string) ([]*user_view_model.UserSessionView, error) {
	sessions := make([]*user_view_model.UserSessionView, len(m.Users))
	for i, user := range m.Users {
		sessions[i] = &user_view_model.UserSessionView{
			UserID:        user.UserID,
			LoginName:     user.LoginName,
			ResourceOwner: user.ResourceOwner,
		}
	}
	return sessions, nil
}

func (m *mockViewUserSession) PrefixAvatarURL() string {
	return "prefix/"
}

type mockViewNoUser struct{}

func (m *mockViewNoUser) UserByID(string) (*user_view_model.UserView, error) {
	return nil, errors.ThrowNotFound(nil, "id", "user not found")
}

func (m *mockViewNoUser) PrefixAvatarURL() string {
	return ""
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
	MFAMaxSetUp            int32
	MFAInitSkipped         time.Time
	PasswordlessTokens     user_view_model.WebAuthNTokens
}

type mockLoginPolicy struct {
	policy *iam_view_model.LoginPolicyView
}

func (m *mockLoginPolicy) LoginPolicyByAggregateID(id string) (*iam_view_model.LoginPolicyView, error) {
	return m.policy, nil
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
			MFAMaxSetUp:            m.MFAMaxSetUp,
			MFAInitSkipped:         m.MFAInitSkipped,
			PasswordlessTokens:     m.PasswordlessTokens,
		},
	}, nil
}

func (m *mockViewUser) PrefixAvatarURL() string {
	return ""
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

type mockUserGrants struct {
	roleCheck  bool
	userGrants int
}

func (m *mockUserGrants) ApplicationByClientID(ctx context.Context, s string) (*proj_view_model.ApplicationView, error) {
	return &proj_view_model.ApplicationView{ProjectRoleCheck: m.roleCheck}, nil
}

func (m *mockUserGrants) UserGrantsByProjectAndUserID(s string, s2 string) ([]*grant_view_model.UserGrantView, error) {
	var grants []*grant_view_model.UserGrantView
	if m.userGrants > 0 {
		grants = make([]*grant_view_model.UserGrantView, m.userGrants)
	}
	return grants, nil
}

func TestAuthRequestRepo_nextSteps(t *testing.T) {
	type fields struct {
		AuthRequests               *cache.AuthRequestCache
		View                       *view.View
		userSessionViewProvider    userSessionViewProvider
		userViewProvider           userViewProvider
		userEventProvider          userEventProvider
		orgViewProvider            orgViewProvider
		userGrantProvider          userGrantProvider
		loginPolicyProvider        loginPolicyViewProvider
		PasswordCheckLifeTime      time.Duration
		ExternalLoginCheckLifeTime time.Duration
		MFAInitSkippedLifeTime     time.Duration
		SecondFactorCheckLifeTime  time.Duration
		MultiFactorCheckLifeTime   time.Duration
	}
	type args struct {
		request       *domain.AuthRequest
		checkLoggedIn bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.NextStep
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
			args{&domain.AuthRequest{Prompt: []domain.Prompt{domain.PromptNone}}, false},
			[]domain.NextStep{&domain.RedirectToCallbackStep{}},
			nil,
		},
		{
			"user not set no active session, login step",
			fields{
				userSessionViewProvider: &mockViewNoUserSession{},
			},
			args{&domain.AuthRequest{}, false},
			[]domain.NextStep{&domain.LoginStep{}},
			nil,
		},
		{
			"user not set no active session, linking users, external user not found option",
			fields{
				userSessionViewProvider: &mockViewNoUserSession{},
			},
			args{&domain.AuthRequest{LinkingUsers: []*domain.ExternalUser{{IDPConfigID: "IDPConfigID", ExternalUserID: "ExternalUserID"}}}, false},
			[]domain.NextStep{&domain.ExternalNotFoundOptionStep{}},
			nil,
		},
		{
			"user not set, prompt select account and internal error, internal error",
			fields{
				userSessionViewProvider: &mockViewErrUserSession{},
			},
			args{&domain.AuthRequest{Prompt: []domain.Prompt{domain.PromptSelectAccount}}, false},
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
							"orgID1",
						},
						{
							"id2",
							"loginname2",
							"orgID2",
						},
					},
				},
				userEventProvider: &mockEventUser{},
			},
			args{&domain.AuthRequest{Prompt: []domain.Prompt{domain.PromptSelectAccount}}, false},
			[]domain.NextStep{
				&domain.LoginStep{},
				&domain.SelectUserStep{
					Users: []domain.UserSelection{
						{
							UserID:            "id1",
							LoginName:         "loginname1",
							SelectionPossible: true,
							ResourceOwner:     "orgID1",
						},
						{
							UserID:            "id2",
							LoginName:         "loginname2",
							SelectionPossible: true,
							ResourceOwner:     "orgID2",
						},
					},
				}},
			nil,
		},
		{
			"user not set, primary domain set, prompt select account, login and select account steps",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					Users: []mockUser{
						{
							"id1",
							"loginname1",
							"orgID1",
						},
						{
							"id2",
							"loginname2",
							"orgID2",
						},
					},
				},
				userEventProvider: &mockEventUser{},
			},
			args{&domain.AuthRequest{Prompt: []domain.Prompt{domain.PromptSelectAccount}, RequestedOrgID: "orgID1"}, false},
			[]domain.NextStep{
				&domain.LoginStep{},
				&domain.SelectUserStep{
					Users: []domain.UserSelection{
						{
							UserID:            "id1",
							LoginName:         "loginname1",
							SelectionPossible: true,
							ResourceOwner:     "orgID1",
						},
						{
							UserID:            "id2",
							LoginName:         "loginname2",
							SelectionPossible: false,
							ResourceOwner:     "orgID2",
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
			args{&domain.AuthRequest{Prompt: []domain.Prompt{domain.PromptSelectAccount}}, false},
			[]domain.NextStep{
				&domain.LoginStep{},
				&domain.SelectUserStep{
					Users: []domain.UserSelection{},
				}},
			nil,
		},
		{
			"user not found, not found error",
			fields{
				userViewProvider:  &mockViewNoUser{},
				userEventProvider: &mockEventUser{},
			},
			args{&domain.AuthRequest{UserID: "UserID"}, false},
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
			args{&domain.AuthRequest{UserID: "UserID"}, false},
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
			args{&domain.AuthRequest{UserID: "UserID"}, false},
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
			args{&domain.AuthRequest{UserID: "UserID"}, false},
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
			args{&domain.AuthRequest{UserID: "UserID"}, false},
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
			args{&domain.AuthRequest{UserID: "UserID", LoginPolicy: &domain.LoginPolicy{}}, false},
			[]domain.NextStep{&domain.PasswordStep{}},
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
			args{&domain.AuthRequest{UserID: "UserID"}, false},
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
			args{&domain.AuthRequest{UserID: "UserID"}, false},
			[]domain.NextStep{&domain.InitUserStep{
				PasswordSet: true,
			}},
			nil,
		},
		{
			"passwordless not verified, passwordless check step",
			fields{
				userSessionViewProvider: &mockViewUserSession{},
				userViewProvider: &mockViewUser{
					PasswordSet:        true,
					PasswordlessTokens: user_view_model.WebAuthNTokens{&user_view_model.WebAuthNView{ID: "id", State: int32(user_model.MFAStateReady)}},
				},
				userEventProvider:        &mockEventUser{},
				orgViewProvider:          &mockViewOrg{State: org_model.OrgStateActive},
				MultiFactorCheckLifeTime: 10 * time.Hour,
			},
			args{&domain.AuthRequest{UserID: "UserID", LoginPolicy: &domain.LoginPolicy{PasswordlessType: domain.PasswordlessTypeAllowed}}, false},
			[]domain.NextStep{&domain.PasswordlessStep{}},
			nil,
		},
		{
			"passwordless verified, email not verified, email verification step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordlessVerification: time.Now().Add(-5 * time.Minute),
					MultiFactorVerification:  time.Now().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:            true,
					PasswordlessTokens:     user_view_model.WebAuthNTokens{&user_view_model.WebAuthNView{ID: "id", State: int32(user_model.MFAStateReady)}},
					PasswordChangeRequired: false,
					IsEmailVerified:        false,
					MFAMaxSetUp:            int32(model.MFALevelMultiFactor),
				},
				userEventProvider:        &mockEventUser{},
				orgViewProvider:          &mockViewOrg{State: org_model.OrgStateActive},
				MultiFactorCheckLifeTime: 10 * time.Hour,
			},
			args{&domain.AuthRequest{
				UserID: "UserID",
				LoginPolicy: &domain.LoginPolicy{
					PasswordlessType: domain.PasswordlessTypeAllowed,
					MultiFactors:     []domain.MultiFactorType{domain.MultiFactorTypeU2FWithPIN},
				},
			}, false},
			[]domain.NextStep{&domain.VerifyEMailStep{}},
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
			args{&domain.AuthRequest{UserID: "UserID", LoginPolicy: &domain.LoginPolicy{}}, false},
			[]domain.NextStep{&domain.InitPasswordStep{}},
			nil,
		},
		{
			"external user (no external verification), external login step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					SecondFactorVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(model.MFALevelSecondFactor),
				},
				userEventProvider:         &mockEventUser{},
				orgViewProvider:           &mockViewOrg{State: org_model.OrgStateActive},
				SecondFactorCheckLifeTime: 18 * time.Hour,
			},
			args{&domain.AuthRequest{UserID: "UserID", SelectedIDPConfigID: "IDPConfigID"}, false},
			[]domain.NextStep{&domain.ExternalLoginStep{SelectedIDPConfigID: "IDPConfigID"}},
			nil,
		},
		{
			"external user (external verification set), callback",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					ExternalLoginVerification: time.Now().UTC().Add(-5 * time.Minute),
					SecondFactorVerification:  time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(model.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: org_model.OrgStateActive},
				userGrantProvider: &mockUserGrants{},
				loginPolicyProvider: &mockLoginPolicy{
					policy: &iam_view_model.LoginPolicyView{},
				},
				ExternalLoginCheckLifeTime: 10 * 24 * time.Hour,
				SecondFactorCheckLifeTime:  18 * time.Hour,
			},
			args{
				&domain.AuthRequest{
					UserID:              "UserID",
					SelectedIDPConfigID: "IDPConfigID",
					Request:             &domain.AuthRequestOIDC{},
					LoginPolicy:         &domain.LoginPolicy{},
				},
				false},
			[]domain.NextStep{&domain.RedirectToCallbackStep{}},
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
			args{&domain.AuthRequest{UserID: "UserID", LoginPolicy: &domain.LoginPolicy{}}, false},
			[]domain.NextStep{&domain.PasswordStep{}},
			nil,
		},
		{
			"external user (no password check needed), callback",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					SecondFactorVerification:  time.Now().UTC().Add(-5 * time.Minute),
					ExternalLoginVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(model.MFALevelSecondFactor),
				},
				userEventProvider:          &mockEventUser{},
				orgViewProvider:            &mockViewOrg{State: org_model.OrgStateActive},
				userGrantProvider:          &mockUserGrants{},
				SecondFactorCheckLifeTime:  18 * time.Hour,
				ExternalLoginCheckLifeTime: 10 * 24 * time.Hour,
			},
			args{
				&domain.AuthRequest{
					UserID:              "UserID",
					SelectedIDPConfigID: "IDPConfigID",
					Request:             &domain.AuthRequestOIDC{},
					LoginPolicy:         &domain.LoginPolicy{},
				}, false},
			[]domain.NextStep{&domain.RedirectToCallbackStep{}},
			nil,
		},
		{
			"password verified, passwordless set up, mfa not verified, mfa check step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:        true,
					PasswordlessTokens: user_view_model.WebAuthNTokens{&user_view_model.WebAuthNView{ID: "id", State: int32(user_model.MFAStateReady)}},
					OTPState:           int32(user_model.MFAStateReady),
					MFAMaxSetUp:        int32(model.MFALevelMultiFactor),
				},
				userEventProvider:         &mockEventUser{},
				orgViewProvider:           &mockViewOrg{State: org_model.OrgStateActive},
				PasswordCheckLifeTime:     10 * 24 * time.Hour,
				SecondFactorCheckLifeTime: 18 * time.Hour,
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors: []domain.SecondFactorType{domain.SecondFactorTypeOTP},
					},
				}, false},
			[]domain.NextStep{&domain.MFAVerificationStep{
				MFAProviders: []domain.MFAType{domain.MFATypeOTP},
			}},
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
					OTPState:    int32(user_model.MFAStateReady),
					MFAMaxSetUp: int32(model.MFALevelSecondFactor),
				},
				userEventProvider:         &mockEventUser{},
				orgViewProvider:           &mockViewOrg{State: org_model.OrgStateActive},
				PasswordCheckLifeTime:     10 * 24 * time.Hour,
				SecondFactorCheckLifeTime: 18 * time.Hour,
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors: []domain.SecondFactorType{domain.SecondFactorTypeOTP},
					},
				}, false},
			[]domain.NextStep{&domain.MFAVerificationStep{
				MFAProviders: []domain.MFAType{domain.MFATypeOTP},
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
					OTPState:    int32(user_model.MFAStateReady),
					MFAMaxSetUp: int32(model.MFALevelSecondFactor),
				},
				userEventProvider:          &mockEventUser{},
				orgViewProvider:            &mockViewOrg{State: org_model.OrgStateActive},
				PasswordCheckLifeTime:      10 * 24 * time.Hour,
				ExternalLoginCheckLifeTime: 10 * 24 * time.Hour,
				SecondFactorCheckLifeTime:  18 * time.Hour,
			},
			args{
				&domain.AuthRequest{
					UserID:              "UserID",
					SelectedIDPConfigID: "IDPConfigID",
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors: []domain.SecondFactorType{domain.SecondFactorTypeOTP},
					},
				}, false},
			[]domain.NextStep{&domain.MFAVerificationStep{
				MFAProviders: []domain.MFAType{domain.MFATypeOTP},
			}},
			nil,
		},
		{
			"password change required and email verified, password change step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     time.Now().UTC().Add(-5 * time.Minute),
					SecondFactorVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:            true,
					PasswordChangeRequired: true,
					IsEmailVerified:        true,
					MFAMaxSetUp:            int32(model.MFALevelSecondFactor),
				},
				userEventProvider:         &mockEventUser{},
				orgViewProvider:           &mockViewOrg{State: org_model.OrgStateActive},
				PasswordCheckLifeTime:     10 * 24 * time.Hour,
				SecondFactorCheckLifeTime: 18 * time.Hour,
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors: []domain.SecondFactorType{domain.SecondFactorTypeOTP},
					},
				}, false},
			[]domain.NextStep{&domain.ChangePasswordStep{}},
			nil,
		},
		{
			"email not verified and no password change required, mail verification step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     time.Now().UTC().Add(-5 * time.Minute),
					SecondFactorVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet: true,
					MFAMaxSetUp: int32(model.MFALevelSecondFactor),
				},
				userEventProvider:         &mockEventUser{},
				orgViewProvider:           &mockViewOrg{State: org_model.OrgStateActive},
				PasswordCheckLifeTime:     10 * 24 * time.Hour,
				SecondFactorCheckLifeTime: 18 * time.Hour,
			},
			args{&domain.AuthRequest{
				UserID: "UserID",
				LoginPolicy: &domain.LoginPolicy{
					SecondFactors: []domain.SecondFactorType{domain.SecondFactorTypeOTP},
				},
			}, false},
			[]domain.NextStep{&domain.VerifyEMailStep{}},
			nil,
		},
		{
			"email not verified and password change required, mail verification step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     time.Now().UTC().Add(-5 * time.Minute),
					SecondFactorVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:            true,
					PasswordChangeRequired: true,
					MFAMaxSetUp:            int32(model.MFALevelSecondFactor),
				},
				userEventProvider:         &mockEventUser{},
				orgViewProvider:           &mockViewOrg{State: org_model.OrgStateActive},
				PasswordCheckLifeTime:     10 * 24 * time.Hour,
				SecondFactorCheckLifeTime: 18 * time.Hour,
			},
			args{&domain.AuthRequest{
				UserID: "UserID",
				LoginPolicy: &domain.LoginPolicy{
					SecondFactors: []domain.SecondFactorType{domain.SecondFactorTypeOTP},
				},
			}, false},
			[]domain.NextStep{&domain.ChangePasswordStep{}, &domain.VerifyEMailStep{}},
			nil,
		},
		{
			"email verified and no password change required, redirect to callback step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     time.Now().UTC().Add(-5 * time.Minute),
					SecondFactorVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(model.MFALevelSecondFactor),
				},
				userEventProvider:         &mockEventUser{},
				orgViewProvider:           &mockViewOrg{State: org_model.OrgStateActive},
				userGrantProvider:         &mockUserGrants{},
				PasswordCheckLifeTime:     10 * 24 * time.Hour,
				SecondFactorCheckLifeTime: 18 * time.Hour,
			},
			args{&domain.AuthRequest{
				UserID:  "UserID",
				Request: &domain.AuthRequestOIDC{},
				LoginPolicy: &domain.LoginPolicy{
					SecondFactors: []domain.SecondFactorType{domain.SecondFactorTypeOTP},
				},
			}, false},
			[]domain.NextStep{&domain.RedirectToCallbackStep{}},
			nil,
		},
		{
			"prompt none, checkLoggedIn true and authenticated, redirect to callback step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     time.Now().UTC().Add(-5 * time.Minute),
					SecondFactorVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(model.MFALevelSecondFactor),
				},
				userEventProvider:         &mockEventUser{},
				orgViewProvider:           &mockViewOrg{State: org_model.OrgStateActive},
				userGrantProvider:         &mockUserGrants{},
				PasswordCheckLifeTime:     10 * 24 * time.Hour,
				SecondFactorCheckLifeTime: 18 * time.Hour,
			},
			args{&domain.AuthRequest{
				UserID:  "UserID",
				Prompt:  []domain.Prompt{domain.PromptNone},
				Request: &domain.AuthRequestOIDC{},
				LoginPolicy: &domain.LoginPolicy{
					SecondFactors: []domain.SecondFactorType{domain.SecondFactorTypeOTP},
				},
			}, true},
			[]domain.NextStep{&domain.RedirectToCallbackStep{}},
			nil,
		},
		{
			"prompt none, checkLoggedIn true, authenticated and required user grants missing, grant required step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     time.Now().UTC().Add(-5 * time.Minute),
					SecondFactorVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(model.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: org_model.OrgStateActive},
				userGrantProvider: &mockUserGrants{
					roleCheck:  true,
					userGrants: 0,
				},
				PasswordCheckLifeTime:     10 * 24 * time.Hour,
				SecondFactorCheckLifeTime: 18 * time.Hour,
			},
			args{&domain.AuthRequest{
				UserID:  "UserID",
				Prompt:  []domain.Prompt{domain.PromptNone},
				Request: &domain.AuthRequestOIDC{},
				LoginPolicy: &domain.LoginPolicy{
					SecondFactors: []domain.SecondFactorType{domain.SecondFactorTypeOTP},
				},
			}, true},
			[]domain.NextStep{&domain.GrantRequiredStep{}},
			nil,
		},
		{
			"prompt none, checkLoggedIn true, authenticated and required user grants exist, redirect to callback step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     time.Now().UTC().Add(-5 * time.Minute),
					SecondFactorVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(model.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: org_model.OrgStateActive},
				userGrantProvider: &mockUserGrants{
					roleCheck:  true,
					userGrants: 2,
				},
				PasswordCheckLifeTime:     10 * 24 * time.Hour,
				SecondFactorCheckLifeTime: 18 * time.Hour,
			},
			args{&domain.AuthRequest{
				UserID:  "UserID",
				Prompt:  []domain.Prompt{domain.PromptNone},
				Request: &domain.AuthRequestOIDC{},
				LoginPolicy: &domain.LoginPolicy{
					SecondFactors: []domain.SecondFactorType{domain.SecondFactorTypeOTP},
				},
			}, true},
			[]domain.NextStep{&domain.RedirectToCallbackStep{}},
			nil,
		},
		{
			"linking users, password step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					SecondFactorVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(model.MFALevelSecondFactor),
				},
				userEventProvider:         &mockEventUser{},
				orgViewProvider:           &mockViewOrg{State: org_model.OrgStateActive},
				SecondFactorCheckLifeTime: 18 * time.Hour,
			},
			args{
				&domain.AuthRequest{
					UserID:              "UserID",
					LoginPolicy:         &domain.LoginPolicy{},
					SelectedIDPConfigID: "IDPConfigID",
					LinkingUsers:        []*domain.ExternalUser{{IDPConfigID: "IDPConfigID", ExternalUserID: "UserID", DisplayName: "DisplayName"}},
				}, false},
			[]domain.NextStep{&domain.PasswordStep{}},
			nil,
		},
		{
			"linking users, linking step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     time.Now().UTC().Add(-5 * time.Minute),
					SecondFactorVerification: time.Now().UTC().Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(model.MFALevelSecondFactor),
				},
				userEventProvider:         &mockEventUser{},
				orgViewProvider:           &mockViewOrg{State: org_model.OrgStateActive},
				SecondFactorCheckLifeTime: 18 * time.Hour,
				PasswordCheckLifeTime:     10 * 24 * time.Hour,
			},
			args{
				&domain.AuthRequest{
					UserID:              "UserID",
					SelectedIDPConfigID: "IDPConfigID",
					LinkingUsers:        []*domain.ExternalUser{{IDPConfigID: "IDPConfigID", ExternalUserID: "UserID", DisplayName: "DisplayName"}},
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors: []domain.SecondFactorType{domain.SecondFactorTypeOTP},
					},
				}, false},
			[]domain.NextStep{&domain.LinkUsersStep{}},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AuthRequestRepo{
				AuthRequests:               tt.fields.AuthRequests,
				View:                       tt.fields.View,
				UserSessionViewProvider:    tt.fields.userSessionViewProvider,
				UserViewProvider:           tt.fields.userViewProvider,
				UserEventProvider:          tt.fields.userEventProvider,
				OrgViewProvider:            tt.fields.orgViewProvider,
				UserGrantProvider:          tt.fields.userGrantProvider,
				LoginPolicyViewProvider:    tt.fields.loginPolicyProvider,
				PasswordCheckLifeTime:      tt.fields.PasswordCheckLifeTime,
				ExternalLoginCheckLifeTime: tt.fields.ExternalLoginCheckLifeTime,
				MFAInitSkippedLifeTime:     tt.fields.MFAInitSkippedLifeTime,
				SecondFactorCheckLifeTime:  tt.fields.SecondFactorCheckLifeTime,
				MultiFactorCheckLifeTime:   tt.fields.MultiFactorCheckLifeTime,
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
		MFAInitSkippedLifeTime    time.Duration
		SecondFactorCheckLifeTime time.Duration
		MultiFactorCheckLifeTime  time.Duration
	}
	type args struct {
		userSession *user_model.UserSessionView
		request     *domain.AuthRequest
		user        *user_model.UserView
		policy      *iam_model.LoginPolicyView
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        domain.NextStep
		wantChecked bool
		errFunc     func(err error) bool
	}{
		//{
		//	"required, prompt and false", //TODO: enable when LevelsOfAssurance is checked
		//	fields{},
		//	args{
		//		request: &domain.AuthRequest{PossibleLOAs: []model.LevelOfAssurance{}},
		//		user: &user_model.UserView{
		//			OTPState: user_model.MFAStateReady,
		//		},
		//	},
		//	false,
		//},
		{
			"not set up, forced by policy, no mfas configured, error",
			fields{
				MFAInitSkippedLifeTime: 30 * 24 * time.Hour,
			},
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						ForceMFA: true,
					},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: model.MFALevelNotSetUp,
					},
				},
			},
			nil,
			false,
			errors.IsPreconditionFailed,
		},
		{
			"not set up, no mfas configured, no prompt and true",
			fields{
				MFAInitSkippedLifeTime: 30 * 24 * time.Hour,
			},
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: model.MFALevelNotSetUp,
					},
				},
			},
			nil,
			true,
			nil,
		},
		{
			"not set up, prompt and false",
			fields{
				MFAInitSkippedLifeTime: 30 * 24 * time.Hour,
			},
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors: []domain.SecondFactorType{domain.SecondFactorTypeOTP},
					},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: model.MFALevelNotSetUp,
					},
				},
			},
			&domain.MFAPromptStep{
				MFAProviders: []domain.MFAType{
					domain.MFATypeOTP,
				},
			},
			false,
			nil,
		},
		{
			"not set up, forced by org, true",
			fields{
				MFAInitSkippedLifeTime: 30 * 24 * time.Hour,
			},
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						ForceMFA:      true,
						SecondFactors: []domain.SecondFactorType{domain.SecondFactorTypeOTP},
					},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: model.MFALevelNotSetUp,
					},
				},
			},
			&domain.MFAPromptStep{
				Required: true,
				MFAProviders: []domain.MFAType{
					domain.MFATypeOTP,
				},
			},
			false,
			nil,
		},
		{
			"not set up and skipped, true",
			fields{
				MFAInitSkippedLifeTime: 30 * 24 * time.Hour,
			},
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp:    model.MFALevelNotSetUp,
						MFAInitSkipped: time.Now().UTC(),
					},
				},
			},
			nil,
			true,
			nil,
		},
		{
			"checked second factor, true",
			fields{
				SecondFactorCheckLifeTime: 18 * time.Hour,
			},
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors: []domain.SecondFactorType{domain.SecondFactorTypeOTP},
					},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: model.MFALevelSecondFactor,
						OTPState:    user_model.MFAStateReady,
					},
				},
				userSession: &user_model.UserSessionView{SecondFactorVerification: time.Now().UTC().Add(-5 * time.Hour)},
			},
			nil,
			true,
			nil,
		},
		{
			"not checked, check and false",
			fields{
				SecondFactorCheckLifeTime: 18 * time.Hour,
			},
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors: []domain.SecondFactorType{domain.SecondFactorTypeOTP},
					},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: model.MFALevelSecondFactor,
						OTPState:    user_model.MFAStateReady,
					},
				},
				userSession: &user_model.UserSessionView{},
			},

			&domain.MFAVerificationStep{
				MFAProviders: []domain.MFAType{domain.MFATypeOTP},
			},
			false,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AuthRequestRepo{
				MFAInitSkippedLifeTime:    tt.fields.MFAInitSkippedLifeTime,
				SecondFactorCheckLifeTime: tt.fields.SecondFactorCheckLifeTime,
				MultiFactorCheckLifeTime:  tt.fields.MultiFactorCheckLifeTime,
			}
			got, ok, err := repo.mfaChecked(tt.args.userSession, tt.args.request, tt.args.user)
			if (tt.errFunc != nil && !tt.errFunc(err)) || (err != nil && tt.errFunc == nil) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if ok != tt.wantChecked {
				t.Errorf("mfaChecked() checked = %v, want %v", ok, tt.wantChecked)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAuthRequestRepo_mfaSkippedOrSetUp(t *testing.T) {
	type fields struct {
		MFAInitSkippedLifeTime time.Duration
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
			args{
				&user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: model.MFALevelSecondFactor,
					},
				},
			},
			true,
		},
		{
			"mfa skipped active, true",
			fields{
				MFAInitSkippedLifeTime: 30 * 24 * time.Hour,
			},
			args{
				&user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp:    -1,
						MFAInitSkipped: time.Now().UTC().Add(-10 * time.Hour),
					},
				},
			},
			true,
		},
		{
			"mfa skipped inactive, false",
			fields{
				MFAInitSkippedLifeTime: 30 * 24 * time.Hour,
			},
			args{
				&user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp:    -1,
						MFAInitSkipped: time.Now().UTC().Add(-40 * 24 * time.Hour),
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AuthRequestRepo{
				MFAInitSkippedLifeTime: tt.fields.MFAInitSkippedLifeTime,
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
			&user_model.UserSessionView{UserID: "id"},
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
				PasswordVerification:     time.Now().UTC().Round(1 * time.Second),
				SecondFactorVerification: time.Time{},
				MultiFactorVerification:  time.Time{},
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
				PasswordVerification:     time.Now().UTC().Round(1 * time.Second),
				SecondFactorVerification: time.Time{},
				MultiFactorVerification:  time.Time{},
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
				PasswordVerification:     time.Now().UTC().Round(1 * time.Second),
				SecondFactorVerification: time.Time{},
				MultiFactorVerification:  time.Time{},
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
				PasswordVerification:     time.Now().UTC().Round(1 * time.Second),
				SecondFactorVerification: time.Now().UTC().Round(1 * time.Second),
				ChangeDate:               time.Now().UTC().Round(1 * time.Second),
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
				userID:        "userID",
				viewProvider:  &mockViewNoUser{},
				eventProvider: &mockEventUser{},
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
			got, err := userByID(context.Background(), tt.args.viewProvider, tt.args.eventProvider, tt.args.userID)
			if (err != nil && tt.wantErr == nil) || (tt.wantErr != nil && !tt.wantErr(err)) {
				t.Errorf("nextSteps() wrong error = %v", err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
