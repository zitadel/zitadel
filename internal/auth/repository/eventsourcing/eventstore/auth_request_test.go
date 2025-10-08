package eventstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	cache "github.com/zitadel/zitadel/internal/auth_request/repository"
	"github.com/zitadel/zitadel/internal/auth_request/repository/mock"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	user_repo "github.com/zitadel/zitadel/internal/repository/user"
	user_model "github.com/zitadel/zitadel/internal/user/model"
	user_es_model "github.com/zitadel/zitadel/internal/user/repository/eventsourcing/model"
	user_view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	testNow = time.Now()
)

type mockViewNoUserSession struct{}

func (m *mockViewNoUserSession) UserSessionByIDs(context.Context, string, string, string) (*user_view_model.UserSessionView, error) {
	return nil, zerrors.ThrowNotFound(nil, "id", "user session not found")
}

func (m *mockViewNoUserSession) UserSessionsByAgentID(context.Context, string, string) ([]*user_view_model.UserSessionView, error) {
	return nil, nil
}

func (m *mockViewNoUserSession) GetLatestUserSessionSequence(ctx context.Context, instanceID string) (*query.CurrentState, error) {
	return &query.CurrentState{State: query.State{Sequence: 0}}, nil
}

type mockViewErrUserSession struct{}

func (m *mockViewErrUserSession) UserSessionByIDs(context.Context, string, string, string) (*user_view_model.UserSessionView, error) {
	return nil, zerrors.ThrowInternal(nil, "id", "internal error")
}

func (m *mockViewErrUserSession) UserSessionsByAgentID(context.Context, string, string) ([]*user_view_model.UserSessionView, error) {
	return nil, zerrors.ThrowInternal(nil, "id", "internal error")
}

func (m *mockViewErrUserSession) GetLatestUserSessionSequence(ctx context.Context, instanceID string) (*query.CurrentState, error) {
	return &query.CurrentState{State: query.State{Sequence: 0}}, nil
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
	SessionState  domain.UserSessionState
}

func (m *mockViewUserSession) UserSessionByIDs(context.Context, string, string, string) (*user_view_model.UserSessionView, error) {
	return &user_view_model.UserSessionView{
		ExternalLoginVerification: sql.NullTime{Time: m.ExternalLoginVerification},
		PasswordlessVerification:  sql.NullTime{Time: m.PasswordlessVerification},
		PasswordVerification:      sql.NullTime{Time: m.PasswordVerification},
		SecondFactorVerification:  sql.NullTime{Time: m.SecondFactorVerification},
		MultiFactorVerification:   sql.NullTime{Time: m.MultiFactorVerification},
	}, nil
}

func (m *mockViewUserSession) UserSessionsByAgentID(context.Context, string, string) ([]*user_view_model.UserSessionView, error) {
	sessions := make([]*user_view_model.UserSessionView, len(m.Users))
	for i, user := range m.Users {
		sessions[i] = &user_view_model.UserSessionView{
			ResourceOwner: user.ResourceOwner,
			State:         sql.Null[domain.UserSessionState]{V: user.SessionState},
			UserID:        user.UserID,
			LoginName:     sql.NullString{String: user.LoginName},
		}
	}
	return sessions, nil
}

func (m *mockViewUserSession) GetLatestUserSessionSequence(ctx context.Context, instanceID string) (*query.CurrentState, error) {
	return &query.CurrentState{State: query.State{Sequence: 0}}, nil
}

type mockViewNoUser struct{}

func (m *mockViewNoUser) UserByID(context.Context, string, string) (*user_view_model.UserView, error) {
	return nil, zerrors.ThrowNotFound(nil, "id", "user not found")
}

type mockEventUser struct {
	Events               []eventstore.Event
	PwCodeExists         bool
	InvitationCodeExists bool
}

func (m *mockEventUser) UserEventsByID(ctx context.Context, id string, changeDate time.Time, types []eventstore.EventType) ([]eventstore.Event, error) {
	return m.Events, nil
}

func (m *mockEventUser) PasswordCodeExists(ctx context.Context, userID string) (bool, error) {
	return m.PwCodeExists, nil
}

func (m *mockEventUser) InviteCodeExists(ctx context.Context, userID string) (bool, error) {
	return m.InvitationCodeExists, nil
}

func (m *mockEventUser) GetLatestUserSessionSequence(ctx context.Context, instanceID string) (*query.CurrentState, error) {
	return &query.CurrentState{State: query.State{Sequence: 0}}, nil
}

func (m *mockEventUser) BulkAddExternalIDPs(ctx context.Context, userID string, externalIDPs []*user_model.ExternalIDP) error {
	return nil
}

type mockEventErrUser struct{}

func (m *mockEventErrUser) UserEventsByID(ctx context.Context, id string, changeDate time.Time, types []eventstore.EventType) ([]eventstore.Event, error) {
	return nil, zerrors.ThrowInternal(nil, "id", "internal error")
}

func (m *mockEventErrUser) PasswordCodeExists(ctx context.Context, userID string) (bool, error) {
	return false, zerrors.ThrowInternal(nil, "id", "internal error")
}

func (m *mockEventErrUser) InviteCodeExists(ctx context.Context, userID string) (bool, error) {
	return false, zerrors.ThrowInternal(nil, "id", "internal error")
}

type mockViewUser struct {
	InitRequired             bool
	PasswordInitRequired     bool
	PasswordSet              bool
	PasswordChanged          time.Time
	PasswordChangeRequired   bool
	IsEmailVerified          bool
	VerifiedEmail            string
	OTPState                 int32
	MFAMaxSetUp              int32
	MFAInitSkipped           time.Time
	PasswordlessInitRequired bool
	PasswordlessTokens       user_view_model.WebAuthNTokens
}

type mockLoginPolicy struct {
	policy *query.LoginPolicy
}

func (m *mockLoginPolicy) LoginPolicyByID(ctx context.Context, _ bool, id string, _ bool) (*query.LoginPolicy, error) {
	return m.policy, nil
}

type mockPrivacyPolicy struct {
	policy *query.PrivacyPolicy
}

func (m mockPrivacyPolicy) PrivacyPolicyByOrg(ctx context.Context, b bool, s string, b2 bool) (*query.PrivacyPolicy, error) {
	return m.policy, nil
}

type mockLabelPolicy struct {
	policy *query.LabelPolicy
}

func (m mockLabelPolicy) ActiveLabelPolicyByOrg(ctx context.Context, s string, b bool) (*query.LabelPolicy, error) {
	return m.policy, nil
}

type mockCustomText struct {
	texts *query.CustomTexts
}

func (m *mockCustomText) CustomTextListByTemplate(ctx context.Context, aggregateID string, text string, withOwnerRemoved bool) (texts *query.CustomTexts, err error) {
	return m.texts, nil
}

type mockLockoutPolicy struct {
	policy *query.LockoutPolicy
}

func (m *mockLockoutPolicy) LockoutPolicyByOrg(context.Context, bool, string) (*query.LockoutPolicy, error) {
	return m.policy, nil
}

type mockPasswordAgePolicy struct {
	policy *query.PasswordAgePolicy
}

func (m *mockPasswordAgePolicy) PasswordAgePolicyByOrg(context.Context, bool, string, bool) (*query.PasswordAgePolicy, error) {
	return m.policy, nil
}

func (m *mockViewUser) UserByID(context.Context, string, string) (*user_view_model.UserView, error) {
	return &user_view_model.UserView{
		State:    int32(user_model.UserStateActive),
		UserName: "UserName",
		HumanView: &user_view_model.HumanView{
			FirstName:                "FirstName",
			InitRequired:             m.InitRequired,
			PasswordInitRequired:     m.PasswordInitRequired,
			PasswordSet:              m.PasswordSet,
			PasswordChangeRequired:   m.PasswordChangeRequired,
			IsEmailVerified:          m.IsEmailVerified,
			VerifiedEmail:            m.VerifiedEmail,
			OTPState:                 m.OTPState,
			MFAMaxSetUp:              m.MFAMaxSetUp,
			MFAInitSkipped:           m.MFAInitSkipped,
			PasswordlessInitRequired: m.PasswordlessInitRequired,
			PasswordlessTokens:       m.PasswordlessTokens,
		},
	}, nil
}

type mockViewOrg struct {
	State domain.OrgState
}

func (m *mockViewOrg) OrgByID(context.Context, string) (*query.Org, error) {
	return &query.Org{
		State: m.State,
	}, nil
}

func (m *mockViewOrg) OrgByPrimaryDomain(context.Context, string) (*query.Org, error) {
	return &query.Org{
		State: m.State,
	}, nil
}

type mockViewErrOrg struct{}

func (m *mockViewErrOrg) OrgByID(context.Context, string) (*query.Org, error) {
	return nil, zerrors.ThrowInternal(nil, "id", "internal error")
}

func (m *mockViewErrOrg) OrgByPrimaryDomain(context.Context, string) (*query.Org, error) {
	return nil, zerrors.ThrowInternal(nil, "id", "internal error")
}

type mockUserGrants struct {
	roleCheck  bool
	userGrants int
}

func (m *mockUserGrants) ProjectByClientID(ctx context.Context, s string) (*query.Project, error) {
	return &query.Project{ProjectRoleCheck: m.roleCheck}, nil
}

func (m *mockUserGrants) UserGrantsByProjectAndUserID(ctx context.Context, s string, s2 string) ([]*query.UserGrant, error) {
	var grants []*query.UserGrant
	if m.userGrants > 0 {
		grants = make([]*query.UserGrant, m.userGrants)
	}
	return grants, nil
}

type mockProject struct {
	hasProject    bool
	projectCheck  bool
	resourceOwner string
}

func (m *mockProject) ProjectByClientID(ctx context.Context, s string) (*query.Project, error) {
	return &query.Project{ResourceOwner: m.resourceOwner, HasProjectCheck: m.projectCheck}, nil
}

func (m *mockProject) SearchProjectGrants(ctx context.Context, queries *query.ProjectGrantSearchQueries, permissionCheck domain.PermissionCheck) (*query.ProjectGrants, error) {
	if m.hasProject {
		mockProjectGrant := new(query.ProjectGrant)
		return &query.ProjectGrants{ProjectGrants: []*query.ProjectGrant{mockProjectGrant}}, nil
	}
	return &query.ProjectGrants{}, nil
}

type mockApp struct {
	app *query.App
}

func (m *mockApp) AppByOIDCClientID(ctx context.Context, id string) (*query.App, error) {
	if m.app != nil {
		return m.app, nil
	}
	return nil, zerrors.ThrowNotFound(nil, "ERROR", "error")
}

type mockIDPUserLinks struct {
	idps []*query.IDPUserLink
}

func (m *mockIDPUserLinks) IDPUserLinks(ctx context.Context, queries *query.IDPUserLinksSearchQuery, permissionCheck domain.PermissionCheck) (*query.IDPUserLinks, error) {
	return &query.IDPUserLinks{Links: m.idps}, nil
}

type mockPasswordReset struct {
	t          *testing.T
	expectCall bool
}

func newMockPasswordReset(expectCall bool) func(*testing.T) passwordReset {
	return func(t *testing.T) passwordReset {
		return &mockPasswordReset{
			t:          t,
			expectCall: expectCall,
		}
	}
}

func (m *mockPasswordReset) RequestSetPassword(ctx context.Context, userID, resourceOwner string, notifyType domain.NotificationType, authRequestID string) (objectDetails *domain.ObjectDetails, err error) {
	if !m.expectCall {
		m.t.Error("unexpected call to RequestSetPassword")
		return nil, nil
	}
	return nil, err
}

type mockPasswordChecker struct {
	err error
}

func (m *mockPasswordChecker) HumanCheckPassword(ctx context.Context, resourceOwner, userID, password string, authReq *domain.AuthRequest) error {
	return m.err
}

func TestAuthRequestRepo_nextSteps(t *testing.T) {
	type fields struct {
		AuthRequests              cache.AuthRequestCache
		View                      *view.View
		userSessionViewProvider   userSessionViewProvider
		userViewProvider          userViewProvider
		userEventProvider         userEventProvider
		orgViewProvider           orgViewProvider
		userGrantProvider         userGrantProvider
		projectProvider           projectProvider
		applicationProvider       applicationProvider
		loginPolicyProvider       loginPolicyViewProvider
		lockoutPolicyProvider     lockoutPolicyViewProvider
		idpUserLinksProvider      idpUserLinksProvider
		privacyPolicyProvider     privacyPolicyProvider
		labelPolicyProvider       labelPolicyProvider
		passwordAgePolicyProvider passwordAgePolicyProvider
		customTextProvider        customTextProvider
		passwordReset             func(t *testing.T) passwordReset
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
			zerrors.IsErrorInvalidArgument,
		},
		{
			"prompt none and checkLoggedIn false, callback step",
			fields{},
			args{&domain.AuthRequest{Prompt: []domain.Prompt{domain.PromptNone}}, false},
			[]domain.NextStep{&domain.RedirectToCallbackStep{}},
			nil,
		},
		{
			"user not set prompt create, registration step",
			fields{
				userSessionViewProvider: &mockViewNoUserSession{},
			},
			args{&domain.AuthRequest{
				Prompt: []domain.Prompt{domain.PromptCreate},
			}, false},
			[]domain.NextStep{&domain.RegistrationStep{}},
			nil,
		},
		{
			"user not set, prompts other than select account, create step",
			fields{
				userSessionViewProvider: &mockViewNoUserSession{},
			},
			args{&domain.AuthRequest{
				Prompt: []domain.Prompt{domain.PromptLogin, domain.PromptConsent},
			}, false},
			[]domain.NextStep{&domain.LoginStep{}},
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
			"user not set no active session selected idp, redirect to external idp step",
			fields{
				userSessionViewProvider: &mockViewNoUserSession{},
			},
			args{&domain.AuthRequest{SelectedIDPConfigID: "id"}, false},
			[]domain.NextStep{&domain.LoginStep{}, &domain.RedirectToExternalIDPStep{}},
			nil,
		},
		{
			"user not set, prompt select account and internal error, internal error",
			fields{
				userSessionViewProvider: &mockViewErrUserSession{},
			},
			args{&domain.AuthRequest{Prompt: []domain.Prompt{domain.PromptSelectAccount}}, false},
			nil,
			zerrors.IsInternal,
		},
		{
			"user not set, prompt select account, select account step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					Users: []mockUser{
						{
							"id1",
							"loginname1",
							"orgID1",
							domain.UserSessionStateActive,
						},
						{
							"id2",
							"loginname2",
							"orgID2",
							domain.UserSessionStateActive,
						},
					},
				},
				userEventProvider: &mockEventUser{},
			},
			args{&domain.AuthRequest{Prompt: []domain.Prompt{domain.PromptSelectAccount}}, false},
			[]domain.NextStep{
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
			"user not set, primary domain set, prompt select account, select account step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					Users: []mockUser{
						{
							"id1",
							"loginname1",
							"orgID1",
							domain.UserSessionStateActive,
						},
						{
							"id2",
							"loginname2",
							"orgID2",
							domain.UserSessionStateActive,
						},
					},
				},
				userEventProvider: &mockEventUser{},
			},
			args{&domain.AuthRequest{Prompt: []domain.Prompt{domain.PromptSelectAccount}, RequestedOrgID: "orgID1"}, false},
			[]domain.NextStep{
				&domain.SelectUserStep{
					Users: []domain.UserSelection{
						{
							UserID:            "id1",
							LoginName:         "loginname1",
							SelectionPossible: true,
							ResourceOwner:     "orgID1",
						},
					},
				}},
			nil,
		},
		{
			"user not set, prompt select account, no active session, login step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					Users: nil,
				},
				userEventProvider: &mockEventUser{},
			},
			args{&domain.AuthRequest{Prompt: []domain.Prompt{domain.PromptSelectAccount}}, false},
			[]domain.NextStep{
				&domain.LoginStep{}},
			nil,
		},
		{
			"user not set single active session, callback step",
			fields{
				AuthRequests: func() cache.AuthRequestCache {
					m := mock.NewMockAuthRequestCache(gomock.NewController(t))
					m.EXPECT().UpdateAuthRequest(gomock.Any(), gomock.Any())
					m.EXPECT().CacheAuthRequest(gomock.Any(), gomock.Any())
					return m
				}(),
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     time.Now().Add(-5 * time.Minute),
					SecondFactorVerification: time.Now().Add(-5 * time.Minute),
					Users: []mockUser{
						{
							"id1",
							"loginname1",
							"orgID1",
							domain.UserSessionStateActive,
						},
					},
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider:   &mockEventUser{},
				orgViewProvider:     &mockViewOrg{State: domain.OrgStateActive},
				userGrantProvider:   &mockUserGrants{},
				projectProvider:     &mockProject{},
				applicationProvider: &mockApp{app: &query.App{OIDCConfig: &query.OIDCApp{AppType: domain.OIDCApplicationTypeWeb}}},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
				loginPolicyProvider: &mockLoginPolicy{
					policy: &query.LoginPolicy{
						AllowUsernamePassword:     true,
						SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
						PasswordCheckLifetime:     database.Duration(10 * 24 * time.Hour),
						SecondFactorCheckLifetime: database.Duration(18 * time.Hour),
					},
				},
				privacyPolicyProvider: &mockPrivacyPolicy{
					policy: &query.PrivacyPolicy{},
				},
				labelPolicyProvider: &mockLabelPolicy{
					policy: &query.LabelPolicy{},
				},
				passwordAgePolicyProvider: &mockPasswordAgePolicy{
					policy: &query.PasswordAgePolicy{},
				},
				customTextProvider: &mockCustomText{
					texts: &query.CustomTexts{},
				},
			},
			args{&domain.AuthRequest{
				Request: &domain.AuthRequestOIDC{},
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     true,
					SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
					PasswordCheckLifetime:     10 * 24 * time.Hour,
					SecondFactorCheckLifetime: 18 * time.Hour,
				},
			}, false},
			[]domain.NextStep{&domain.RedirectToCallbackStep{}},
			nil,
		},
		{
			"user not set single inactive session, callback step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     time.Now().Add(-5 * time.Minute),
					SecondFactorVerification: time.Now().Add(-5 * time.Minute),
					Users: []mockUser{
						{
							"id1",
							"loginname1",
							"orgID1",
							domain.UserSessionStateTerminated,
						},
					},
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider:   &mockEventUser{},
				orgViewProvider:     &mockViewOrg{State: domain.OrgStateActive},
				userGrantProvider:   &mockUserGrants{},
				projectProvider:     &mockProject{},
				applicationProvider: &mockApp{app: &query.App{OIDCConfig: &query.OIDCApp{AppType: domain.OIDCApplicationTypeWeb}}},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{&domain.AuthRequest{
				Request: &domain.AuthRequestOIDC{},
				LoginPolicy: &domain.LoginPolicy{
					SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
					PasswordCheckLifetime:     10 * 24 * time.Hour,
					SecondFactorCheckLifetime: 18 * time.Hour,
				},
			}, false},
			[]domain.NextStep{&domain.SelectUserStep{Users: []domain.UserSelection{
				{
					UserID:            "id1",
					LoginName:         "loginname1",
					ResourceOwner:     "orgID1",
					UserSessionState:  domain.UserSessionStateTerminated,
					SelectionPossible: true,
				},
			}}},
			nil,
		},
		{
			"user not set multiple sessions, select account step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					Users: []mockUser{
						{
							"id1",
							"loginname1",
							"orgID1",
							domain.UserSessionStateActive,
						},
						{
							"id2",
							"loginname2",
							"orgID2",
							domain.UserSessionStateTerminated,
						},
					},
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider:   &mockEventUser{},
				orgViewProvider:     &mockViewOrg{State: domain.OrgStateActive},
				userGrantProvider:   &mockUserGrants{},
				projectProvider:     &mockProject{},
				applicationProvider: &mockApp{app: &query.App{OIDCConfig: &query.OIDCApp{AppType: domain.OIDCApplicationTypeWeb}}},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{&domain.AuthRequest{
				Request: &domain.AuthRequestOIDC{},
				LoginPolicy: &domain.LoginPolicy{
					SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
					PasswordCheckLifetime:     10 * 24 * time.Hour,
					SecondFactorCheckLifetime: 18 * time.Hour,
				},
			}, false},
			[]domain.NextStep{&domain.SelectUserStep{
				Users: []domain.UserSelection{
					{
						UserID:            "id1",
						LoginName:         "loginname1",
						SelectionPossible: true,
						ResourceOwner:     "orgID1",
						UserSessionState:  domain.UserSessionStateActive,
					},
					{
						UserID:            "id2",
						LoginName:         "loginname2",
						SelectionPossible: true,
						ResourceOwner:     "orgID2",
						UserSessionState:  domain.UserSessionStateTerminated,
					},
				},
			}},
			nil,
		},
		{
			"user not found, not found error",
			fields{
				userViewProvider:  &mockViewNoUser{},
				userEventProvider: &mockEventUser{},
			},
			args{&domain.AuthRequest{UserID: "UserID", LoginPolicy: &domain.LoginPolicy{}}, false},
			nil,
			zerrors.IsNotFound,
		},
		{
			"user not active, precondition failed error",
			fields{
				userViewProvider: &mockViewUser{},
				userEventProvider: &mockEventUser{
					Events: []eventstore.Event{
						&es_models.Event{
							AggregateType: user_repo.AggregateType,
							Typ:           user_repo.UserDeactivatedType,
						},
					},
				},
				orgViewProvider: &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{

						ShowFailures: true,
					},
				},
			},
			args{&domain.AuthRequest{UserID: "UserID", LoginPolicy: &domain.LoginPolicy{}}, false},
			nil,
			zerrors.IsPreconditionFailed,
		},
		{
			"user locked, precondition failed error",
			fields{
				userViewProvider: &mockViewUser{},
				userEventProvider: &mockEventUser{
					Events: []eventstore.Event{
						&es_models.Event{
							AggregateType: user_repo.AggregateType,
							Typ:           user_repo.UserLockedType,
						},
					},
				},
				orgViewProvider: &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
			},
			args{&domain.AuthRequest{UserID: "UserID", LoginPolicy: &domain.LoginPolicy{}}, false},
			nil,
			zerrors.IsPreconditionFailed,
		},
		{
			"org error, internal error",
			fields{
				userViewProvider:  &mockViewUser{},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewErrOrg{},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
			},
			args{&domain.AuthRequest{UserID: "UserID", LoginPolicy: &domain.LoginPolicy{}}, false},
			nil,
			zerrors.IsInternal,
		},
		{
			"org not active, precondition failed error",
			fields{
				userViewProvider:  &mockViewUser{},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateInactive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
			},
			args{&domain.AuthRequest{UserID: "UserID", LoginPolicy: &domain.LoginPolicy{}}, false},
			nil,
			zerrors.IsPreconditionFailed,
		},
		{
			"usersession not found, new user session, password step",
			fields{
				userSessionViewProvider: &mockViewNoUserSession{},
				userViewProvider: &mockViewUser{
					PasswordSet: true,
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword: true,
					},
				},
				false,
			},
			[]domain.NextStep{&domain.PasswordStep{}},
			nil,
		},
		{
			"usersession error, internal error",
			fields{
				userSessionViewProvider: &mockViewErrUserSession{},
				userViewProvider:        &mockViewUser{},
				userEventProvider:       &mockEventUser{},
				orgViewProvider:         &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{&domain.AuthRequest{UserID: "UserID", LoginPolicy: &domain.LoginPolicy{}}, false},
			nil,
			zerrors.IsInternal,
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
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				loginPolicyProvider: &mockLoginPolicy{
					policy: &query.LoginPolicy{
						AllowUsernamePassword: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword: true,
					},
				},
				false,
			},
			[]domain.NextStep{&domain.InitUserStep{
				PasswordSet: true,
			}},
			nil,
		},
		{
			"passwordless not initialised, passwordless prompt step",
			fields{
				userSessionViewProvider: &mockViewUserSession{},
				userViewProvider: &mockViewUser{
					PasswordlessInitRequired: true,
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				loginPolicyProvider: &mockLoginPolicy{
					policy: &query.LoginPolicy{
						MultiFactorCheckLifetime: database.Duration(10 * time.Hour),
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword: true,
						PasswordlessType:      domain.PasswordlessTypeAllowed,
					},
				},
				false,
			},
			[]domain.NextStep{&domain.PasswordlessRegistrationPromptStep{}},
			nil,
		},
		{
			"passwordless not verified, no password set, passwordless check step",
			fields{
				userSessionViewProvider: &mockViewUserSession{},
				userViewProvider: &mockViewUser{
					PasswordlessTokens: user_view_model.WebAuthNTokens{&user_view_model.WebAuthNView{ID: "id", State: int32(user_model.MFAStateReady)}},
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				loginPolicyProvider: &mockLoginPolicy{
					policy: &query.LoginPolicy{
						MultiFactorCheckLifetime: database.Duration(10 * time.Hour),
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword: true,
						PasswordlessType:      domain.PasswordlessTypeAllowed,
					},
				}, false,
			},
			[]domain.NextStep{&domain.PasswordlessStep{}},
			nil,
		},
		{
			"passwordless not verified, passwordless check step, downgrade possible",
			fields{
				userSessionViewProvider: &mockViewUserSession{},
				userViewProvider: &mockViewUser{
					PasswordSet:        true,
					PasswordlessTokens: user_view_model.WebAuthNTokens{&user_view_model.WebAuthNView{ID: "id", State: int32(user_model.MFAStateReady)}},
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				loginPolicyProvider: &mockLoginPolicy{
					policy: &query.LoginPolicy{
						MultiFactorCheckLifetime: database.Duration(10 * time.Hour),
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword: true,
						PasswordlessType:      domain.PasswordlessTypeAllowed,
					},
				}, false,
			},
			[]domain.NextStep{&domain.PasswordlessStep{PasswordSet: true}},
			nil,
		},
		{
			"passwordless verified, email not verified, email verification step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordlessVerification: testNow.Add(-5 * time.Minute),
					MultiFactorVerification:  testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:            true,
					PasswordlessTokens:     user_view_model.WebAuthNTokens{&user_view_model.WebAuthNView{ID: "id", State: int32(user_model.MFAStateReady)}},
					PasswordChangeRequired: false,
					IsEmailVerified:        false,
					MFAMaxSetUp:            int32(domain.MFALevelMultiFactor),
				},
				userEventProvider: &mockEventUser{},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				orgViewProvider:      &mockViewOrg{State: domain.OrgStateActive},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword:    true,
						PasswordlessType:         domain.PasswordlessTypeAllowed,
						MultiFactors:             []domain.MultiFactorType{domain.MultiFactorTypeU2FWithPIN},
						MultiFactorCheckLifetime: 10 * time.Hour,
					},
				},
				false,
			},
			[]domain.NextStep{&domain.VerifyEMailStep{}},
			nil,
		},
		{
			"password not set (email not verified), invite code exists, invite step",
			fields{
				userSessionViewProvider: &mockViewUserSession{},
				userViewProvider: &mockViewUser{
					PasswordInitRequired: true,
				},
				userEventProvider: &mockEventUser{
					InvitationCodeExists: true,
				},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				orgViewProvider:      &mockViewOrg{State: domain.OrgStateActive},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword: true,
					},
				},
				false,
			},
			[]domain.NextStep{&domain.VerifyInviteStep{}},
			nil,
		},
		{
			"password not set (email not verified), verify email with password step",
			fields{
				userSessionViewProvider: &mockViewUserSession{},
				userViewProvider: &mockViewUser{
					PasswordInitRequired: true,
				},
				userEventProvider: &mockEventUser{},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				orgViewProvider:      &mockViewOrg{State: domain.OrgStateActive},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword: true,
					},
				},
				false,
			},
			[]domain.NextStep{&domain.VerifyEMailStep{InitPassword: true}},
			nil,
		},
		{
			"password not set, but idp, email not verified, verify email step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					ExternalLoginVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider:  &mockViewUser{},
				userEventProvider: &mockEventUser{},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				orgViewProvider: &mockViewOrg{State: domain.OrgStateActive},
				idpUserLinksProvider: &mockIDPUserLinks{
					[]*query.IDPUserLink{
						{
							IDPID:            "idpID",
							UserID:           "userID",
							IDPName:          "idpName",
							ProvidedUserID:   "providedUserID",
							ProvidedUsername: "providedUsername",
						},
					},
				},
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword:      true,
						ExternalLoginCheckLifetime: 10 * 24 * time.Hour,
					},
					SelectedIDPConfigID: "idpID",
				},
				false,
			},
			[]domain.NextStep{&domain.VerifyEMailStep{}},
			nil,
		},
		{
			"password not set (email not verified), init password step",
			fields{
				userSessionViewProvider: &mockViewUserSession{},
				userViewProvider: &mockViewUser{
					PasswordInitRequired: true,
				},
				userEventProvider: &mockEventUser{},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				orgViewProvider:      &mockViewOrg{State: domain.OrgStateActive},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword: true,
					},
				},
				false,
			},
			[]domain.NextStep{&domain.VerifyEMailStep{InitPassword: true}},
			nil,
		},
		{
			"password not set (email verified), init password step",
			fields{
				userSessionViewProvider: &mockViewUserSession{},
				userViewProvider: &mockViewUser{
					PasswordInitRequired: true,
					IsEmailVerified:      true,
				},
				userEventProvider: &mockEventUser{
					PwCodeExists: true,
				},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				orgViewProvider:      &mockViewOrg{State: domain.OrgStateActive},
				idpUserLinksProvider: &mockIDPUserLinks{},
				passwordReset:        newMockPasswordReset(false),
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword: true,
					},
				},
				false,
			},
			[]domain.NextStep{&domain.InitPasswordStep{}},
			nil,
		},
		{
			"password not set (email verified, password code not exists), create code, init password step",
			fields{
				userSessionViewProvider: &mockViewUserSession{},
				userViewProvider: &mockViewUser{
					PasswordInitRequired: true,
					IsEmailVerified:      true,
				},
				userEventProvider: &mockEventUser{
					PwCodeExists: false,
				},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				orgViewProvider:      &mockViewOrg{State: domain.OrgStateActive},
				idpUserLinksProvider: &mockIDPUserLinks{},
				passwordReset:        newMockPasswordReset(true),
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword: true,
					},
				},
				false,
			},
			[]domain.NextStep{&domain.InitPasswordStep{}},
			nil,
		},
		{
			"external user (idp selected, no external verification), external login step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				orgViewProvider: &mockViewOrg{State: domain.OrgStateActive},
				loginPolicyProvider: &mockLoginPolicy{
					policy: &query.LoginPolicy{
						SecondFactorCheckLifetime: database.Duration(18 * time.Hour),
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{&domain.AuthRequest{
				UserID:              "UserID",
				SelectedIDPConfigID: "IDPConfigID",
				AllowedExternalIDPs: []*domain.IDPProvider{{IDPConfigID: "IDPConfigID"}},
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     false,
					SecondFactorCheckLifetime: 18 * time.Hour,
				}}, false},
			[]domain.NextStep{&domain.ExternalLoginStep{SelectedIDPConfigID: "IDPConfigID"}},
			nil,
		},
		{
			"external user (idp selected, not allowed, no external verification), error",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				orgViewProvider: &mockViewOrg{State: domain.OrgStateActive},
				loginPolicyProvider: &mockLoginPolicy{
					policy: &query.LoginPolicy{
						SecondFactorCheckLifetime: database.Duration(18 * time.Hour),
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{&domain.AuthRequest{
				UserID:              "UserID",
				SelectedIDPConfigID: "IDPConfigID",
				AllowedExternalIDPs: []*domain.IDPProvider{},
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     false,
					SecondFactorCheckLifetime: 18 * time.Hour,
				}}, false},
			nil,
			zerrors.IsPreconditionFailed,
		},
		{
			"external user (idp link, no external verification), external login step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				orgViewProvider: &mockViewOrg{State: domain.OrgStateActive},
				loginPolicyProvider: &mockLoginPolicy{
					policy: &query.LoginPolicy{
						SecondFactorCheckLifetime: database.Duration(18 * time.Hour),
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{
					[]*query.IDPUserLink{
						{IDPID: "IDPConfigID"},
					},
				},
			},
			args{&domain.AuthRequest{
				UserID:              "UserID",
				SelectedIDPConfigID: "",
				AllowedExternalIDPs: []*domain.IDPProvider{{IDPConfigID: "IDPConfigID"}},
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     false,
					SecondFactorCheckLifetime: 18 * time.Hour,
				}}, false},
			[]domain.NextStep{&domain.ExternalLoginStep{SelectedIDPConfigID: "IDPConfigID"}},
			nil,
		},
		{
			"external user (idp link not allowed, no external verification), external login step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				orgViewProvider: &mockViewOrg{State: domain.OrgStateActive},
				loginPolicyProvider: &mockLoginPolicy{
					policy: &query.LoginPolicy{
						SecondFactorCheckLifetime: database.Duration(18 * time.Hour),
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{
					[]*query.IDPUserLink{
						{IDPID: "IDPConfigID1"},
					},
				},
			},
			args{&domain.AuthRequest{
				UserID:              "UserID",
				SelectedIDPConfigID: "",
				AllowedExternalIDPs: []*domain.IDPProvider{{IDPConfigID: "IDPConfigID2"}},
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     false,
					SecondFactorCheckLifetime: 18 * time.Hour,
				}}, false},
			[]domain.NextStep{&domain.ExternalLoginStep{SelectedIDPConfigID: "IDPConfigID2"}},
			nil,
		},
		{
			"external user (idp link not allowed, none allowed, no external verification), external login step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				orgViewProvider: &mockViewOrg{State: domain.OrgStateActive},
				loginPolicyProvider: &mockLoginPolicy{
					policy: &query.LoginPolicy{
						SecondFactorCheckLifetime: database.Duration(18 * time.Hour),
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{
					[]*query.IDPUserLink{
						{IDPID: "IDPConfigID1"},
					},
				},
			},
			args{&domain.AuthRequest{
				UserID:              "UserID",
				SelectedIDPConfigID: "",
				AllowedExternalIDPs: []*domain.IDPProvider{},
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     false,
					SecondFactorCheckLifetime: 18 * time.Hour,
				}}, false},
			nil,
			zerrors.IsPreconditionFailed,
		},
		{
			"external user (no idp allowed, no external verification), error",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				orgViewProvider: &mockViewOrg{State: domain.OrgStateActive},
				loginPolicyProvider: &mockLoginPolicy{
					policy: &query.LoginPolicy{
						SecondFactorCheckLifetime: database.Duration(18 * time.Hour),
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{&domain.AuthRequest{
				UserID:              "UserID",
				SelectedIDPConfigID: "",
				AllowedExternalIDPs: []*domain.IDPProvider{},
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     false,
					SecondFactorCheckLifetime: 18 * time.Hour,
				}}, false},
			nil,
			zerrors.IsPreconditionFailed,
		},
		{
			"external user (only idp available, no external verification), external login step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				orgViewProvider: &mockViewOrg{State: domain.OrgStateActive},
				loginPolicyProvider: &mockLoginPolicy{
					policy: &query.LoginPolicy{
						SecondFactorCheckLifetime: database.Duration(18 * time.Hour),
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{
					idps: []*query.IDPUserLink{{IDPID: "IDPConfigID"}},
				},
			},
			args{&domain.AuthRequest{
				UserID:              "UserID",
				AllowedExternalIDPs: []*domain.IDPProvider{{IDPConfigID: "IDPConfigID"}},
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     false,
					SecondFactorCheckLifetime: 18 * time.Hour,
				}}, false},
			[]domain.NextStep{&domain.ExternalLoginStep{SelectedIDPConfigID: "IDPConfigID"}},
			nil,
		}, {
			"external user (only idp available, no allowed, no external verification), external login step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				orgViewProvider: &mockViewOrg{State: domain.OrgStateActive},
				loginPolicyProvider: &mockLoginPolicy{
					policy: &query.LoginPolicy{
						SecondFactorCheckLifetime: database.Duration(18 * time.Hour),
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{
					idps: []*query.IDPUserLink{{IDPID: "IDPConfigID"}},
				},
			},
			args{&domain.AuthRequest{
				UserID:              "UserID",
				AllowedExternalIDPs: []*domain.IDPProvider{},
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     false,
					SecondFactorCheckLifetime: 18 * time.Hour,
				}}, false},
			nil,
			zerrors.IsPreconditionFailed,
		},
		{
			"external user (external verification set), callback",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					ExternalLoginVerification: testNow.Add(-5 * time.Minute),
					SecondFactorVerification:  testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider:   &mockEventUser{},
				orgViewProvider:     &mockViewOrg{State: domain.OrgStateActive},
				userGrantProvider:   &mockUserGrants{},
				projectProvider:     &mockProject{},
				applicationProvider: &mockApp{app: &query.App{OIDCConfig: &query.OIDCApp{AppType: domain.OIDCApplicationTypeWeb}}},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID:              "UserID",
					SelectedIDPConfigID: "IDPConfigID",
					Request:             &domain.AuthRequestOIDC{},
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword:      false,
						ExternalLoginCheckLifetime: 10 * 24 * time.Hour,
						SecondFactorCheckLifetime:  18 * time.Hour,
					},
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
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				loginPolicyProvider: &mockLoginPolicy{
					policy: &query.LoginPolicy{
						PasswordCheckLifetime: database.Duration(10 * 24 * time.Hour),
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword: true,
					},
				},
				false,
			},
			[]domain.NextStep{&domain.PasswordStep{}},
			nil,
		},
		{
			"external user (no password check needed), callback",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					SecondFactorVerification:  testNow.Add(-5 * time.Minute),
					ExternalLoginVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider:   &mockEventUser{},
				orgViewProvider:     &mockViewOrg{State: domain.OrgStateActive},
				userGrantProvider:   &mockUserGrants{},
				projectProvider:     &mockProject{},
				applicationProvider: &mockApp{app: &query.App{OIDCConfig: &query.OIDCApp{AppType: domain.OIDCApplicationTypeWeb}}},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID:              "UserID",
					SelectedIDPConfigID: "IDPConfigID",
					Request:             &domain.AuthRequestOIDC{},
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword:      true,
						SecondFactorCheckLifetime:  18 * time.Hour,
						ExternalLoginCheckLifetime: 10 * 24 * time.Hour,
					},
				}, false},
			[]domain.NextStep{&domain.RedirectToCallbackStep{}},
			nil,
		},
		{
			"password verified, passwordless set up, mfa not verified, mfa check step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					VerifiedEmail:      "verified",
					PasswordSet:        true,
					PasswordlessTokens: user_view_model.WebAuthNTokens{&user_view_model.WebAuthNView{ID: "id", State: int32(user_model.MFAStateReady)}},
					OTPState:           int32(user_model.MFAStateReady),
					MFAMaxSetUp:        int32(domain.MFALevelMultiFactor),
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword:     true,
						SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
						PasswordCheckLifetime:     10 * 24 * time.Hour,
						SecondFactorCheckLifetime: 18 * time.Hour,
					},
				}, false},
			[]domain.NextStep{&domain.MFAVerificationStep{
				MFAProviders: []domain.MFAType{domain.MFATypeTOTP},
			}},
			nil,
		},
		{
			"mfa not verified, mfa check step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					VerifiedEmail: "verified",
					PasswordSet:   true,
					OTPState:      int32(user_model.MFAStateReady),
					MFAMaxSetUp:   int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword:     true,
						SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
						PasswordCheckLifetime:     10 * 24 * time.Hour,
						SecondFactorCheckLifetime: 18 * time.Hour,
					},
				}, false},
			[]domain.NextStep{&domain.MFAVerificationStep{
				MFAProviders: []domain.MFAType{domain.MFATypeTOTP},
			}},
			nil,
		},
		{
			"external user, mfa not verified, mfa check step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:      testNow.Add(-5 * time.Minute),
					ExternalLoginVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					VerifiedEmail: "verified",
					PasswordSet:   true,
					OTPState:      int32(user_model.MFAStateReady),
					MFAMaxSetUp:   int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID:              "UserID",
					SelectedIDPConfigID: "IDPConfigID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword:      true,
						SecondFactors:              []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
						PasswordCheckLifetime:      10 * 24 * time.Hour,
						ExternalLoginCheckLifetime: 10 * 24 * time.Hour,
						SecondFactorCheckLifetime:  18 * time.Hour,
					},
				}, false},
			[]domain.NextStep{&domain.MFAVerificationStep{
				MFAProviders: []domain.MFAType{domain.MFATypeTOTP},
			}},
			nil,
		},
		{
			"external user, mfa not verified, email never verified, email verification step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:      testNow.Add(-5 * time.Minute),
					ExternalLoginVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet: true,
					OTPState:    int32(user_model.MFAStateReady),
					MFAMaxSetUp: int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID:              "UserID",
					SelectedIDPConfigID: "IDPConfigID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword:      true,
						SecondFactors:              []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
						PasswordCheckLifetime:      10 * 24 * time.Hour,
						ExternalLoginCheckLifetime: 10 * 24 * time.Hour,
						SecondFactorCheckLifetime:  18 * time.Hour,
					},
				}, false},
			[]domain.NextStep{&domain.VerifyEMailStep{
				InitPassword: false,
			}},
			nil,
		},
		{
			"password change required and email verified, password change step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     testNow.Add(-5 * time.Minute),
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:            true,
					PasswordChangeRequired: true,
					IsEmailVerified:        true,
					MFAMaxSetUp:            int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID: "UserID",
					LoginPolicy: &domain.LoginPolicy{
						AllowUsernamePassword:     true,
						SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
						PasswordCheckLifetime:     10 * 24 * time.Hour,
						SecondFactorCheckLifetime: 18 * time.Hour,
					},
				}, false},
			[]domain.NextStep{&domain.ChangePasswordStep{}},
			nil,
		},
		{
			"email not verified and no password change required, mail verification step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     testNow.Add(-5 * time.Minute),
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet: true,
					MFAMaxSetUp: int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{&domain.AuthRequest{
				UserID: "UserID",
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     true,
					SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
					PasswordCheckLifetime:     10 * 24 * time.Hour,
					SecondFactorCheckLifetime: 18 * time.Hour,
				},
			}, false},
			[]domain.NextStep{&domain.VerifyEMailStep{}},
			nil,
		},
		{
			"email not verified (but had before) and password change required, mail verification step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     testNow.Add(-5 * time.Minute),
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					VerifiedEmail:          "verified",
					PasswordSet:            true,
					PasswordChangeRequired: true,
					MFAMaxSetUp:            int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{&domain.AuthRequest{
				UserID: "UserID",
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     true,
					SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
					PasswordCheckLifetime:     10 * 24 * time.Hour,
					SecondFactorCheckLifetime: 18 * time.Hour,
				},
			}, false},
			[]domain.NextStep{&domain.ChangePasswordStep{}, &domain.VerifyEMailStep{}},
			nil,
		},
		{
			"password change expired, password change step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     testNow.Add(-5 * time.Minute),
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:            true,
					PasswordChangeRequired: false,
					PasswordChanged:        testNow.Add(-50 * 24 * time.Hour),
					IsEmailVerified:        true,
					MFAMaxSetUp:            int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				passwordAgePolicyProvider: &mockPasswordAgePolicy{
					policy: &query.PasswordAgePolicy{
						MaxAgeDays: 30,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{&domain.AuthRequest{
				UserID: "UserID",
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     true,
					SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
					PasswordCheckLifetime:     10 * 24 * time.Hour,
					SecondFactorCheckLifetime: 18 * time.Hour,
				},
				PasswordAgePolicy: &domain.PasswordAgePolicy{
					MaxAgeDays: 30,
				},
			}, false},
			[]domain.NextStep{&domain.ChangePasswordStep{Expired: true}},
			nil,
		},
		{
			"email verified and no password change required, redirect to callback step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     testNow.Add(-5 * time.Minute),
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider:   &mockEventUser{},
				orgViewProvider:     &mockViewOrg{State: domain.OrgStateActive},
				userGrantProvider:   &mockUserGrants{},
				projectProvider:     &mockProject{},
				applicationProvider: &mockApp{app: &query.App{OIDCConfig: &query.OIDCApp{AppType: domain.OIDCApplicationTypeWeb}}},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{&domain.AuthRequest{
				UserID:  "UserID",
				Request: &domain.AuthRequestOIDC{},
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     true,
					SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
					PasswordCheckLifetime:     10 * 24 * time.Hour,
					SecondFactorCheckLifetime: 18 * time.Hour,
				},
			}, false},
			[]domain.NextStep{&domain.RedirectToCallbackStep{}},
			nil,
		},
		{
			"prompt none, checkLoggedIn true and authenticated, redirect to callback step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     testNow.Add(-5 * time.Minute),
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider:   &mockEventUser{},
				orgViewProvider:     &mockViewOrg{State: domain.OrgStateActive},
				userGrantProvider:   &mockUserGrants{},
				projectProvider:     &mockProject{},
				applicationProvider: &mockApp{app: &query.App{OIDCConfig: &query.OIDCApp{AppType: domain.OIDCApplicationTypeWeb}}},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{&domain.AuthRequest{
				UserID:  "UserID",
				Prompt:  []domain.Prompt{domain.PromptNone},
				Request: &domain.AuthRequestOIDC{},
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     true,
					SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
					PasswordCheckLifetime:     10 * 24 * time.Hour,
					SecondFactorCheckLifetime: 18 * time.Hour,
				},
			}, true},
			[]domain.NextStep{&domain.RedirectToCallbackStep{}},
			nil,
		},
		{
			"prompt none, checkLoggedIn true, authenticated and native, login succeeded step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     testNow.Add(-5 * time.Minute),
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider:   &mockEventUser{},
				orgViewProvider:     &mockViewOrg{State: domain.OrgStateActive},
				userGrantProvider:   &mockUserGrants{},
				projectProvider:     &mockProject{},
				applicationProvider: &mockApp{app: &query.App{OIDCConfig: &query.OIDCApp{AppType: domain.OIDCApplicationTypeNative}}},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{&domain.AuthRequest{
				UserID:  "UserID",
				Prompt:  []domain.Prompt{domain.PromptNone},
				Request: &domain.AuthRequestOIDC{},
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     true,
					SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
					PasswordCheckLifetime:     10 * 24 * time.Hour,
					SecondFactorCheckLifetime: 18 * time.Hour,
				},
			}, true},
			[]domain.NextStep{&domain.LoginSucceededStep{}, &domain.RedirectToCallbackStep{}},
			nil,
		},
		{
			"prompt none, checkLoggedIn true, authenticated and required user grants missing, grant required step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     testNow.Add(-5 * time.Minute),
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				userGrantProvider: &mockUserGrants{
					roleCheck:  true,
					userGrants: 0,
				},
				projectProvider: &mockProject{},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{&domain.AuthRequest{
				UserID:  "UserID",
				Prompt:  []domain.Prompt{domain.PromptNone},
				Request: &domain.AuthRequestOIDC{},
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     true,
					SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
					PasswordCheckLifetime:     10 * 24 * time.Hour,
					SecondFactorCheckLifetime: 18 * time.Hour,
				},
			}, true},
			[]domain.NextStep{&domain.GrantRequiredStep{}},
			nil,
		},
		{
			"prompt none, checkLoggedIn true, authenticated and required user grants exist, redirect to callback step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     testNow.Add(-5 * time.Minute),
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				userGrantProvider: &mockUserGrants{
					roleCheck:  true,
					userGrants: 2,
				},
				projectProvider:     &mockProject{},
				applicationProvider: &mockApp{app: &query.App{OIDCConfig: &query.OIDCApp{AppType: domain.OIDCApplicationTypeWeb}}},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{&domain.AuthRequest{
				UserID:  "UserID",
				Prompt:  []domain.Prompt{domain.PromptNone},
				Request: &domain.AuthRequestOIDC{},
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     true,
					SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
					PasswordCheckLifetime:     10 * 24 * time.Hour,
					SecondFactorCheckLifetime: 18 * time.Hour,
				},
			}, true},
			[]domain.NextStep{&domain.RedirectToCallbackStep{}},
			nil,
		},
		{
			"prompt none, checkLoggedIn true, authenticated and required project missing, project required step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     testNow.Add(-5 * time.Minute),
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				userGrantProvider: &mockUserGrants{},
				projectProvider: &mockProject{
					projectCheck:  true,
					hasProject:    false,
					resourceOwner: "other-org",
				},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{&domain.AuthRequest{
				UserID:  "UserID",
				Prompt:  []domain.Prompt{domain.PromptNone},
				Request: &domain.AuthRequestOIDC{},
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     true,
					SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
					PasswordCheckLifetime:     10 * 24 * time.Hour,
					SecondFactorCheckLifetime: 18 * time.Hour,
				},
			}, true},
			[]domain.NextStep{&domain.ProjectRequiredStep{}},
			nil,
		},
		{
			"prompt none, checkLoggedIn true, authenticated and required project exist, redirect to callback step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					PasswordVerification:     testNow.Add(-5 * time.Minute),
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				userGrantProvider: &mockUserGrants{},
				projectProvider: &mockProject{
					projectCheck:  true,
					hasProject:    true,
					resourceOwner: "other-org",
				},
				applicationProvider: &mockApp{app: &query.App{OIDCConfig: &query.OIDCApp{AppType: domain.OIDCApplicationTypeWeb}}},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{&domain.AuthRequest{
				UserID:  "UserID",
				Prompt:  []domain.Prompt{domain.PromptNone},
				Request: &domain.AuthRequestOIDC{},
				LoginPolicy: &domain.LoginPolicy{
					AllowUsernamePassword:     true,
					SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
					PasswordCheckLifetime:     10 * 24 * time.Hour,
					SecondFactorCheckLifetime: 18 * time.Hour,
				},
			}, true},
			[]domain.NextStep{&domain.RedirectToCallbackStep{}},
			nil,
		},
		{
			"linking users, password step",
			fields{
				userSessionViewProvider: &mockViewUserSession{
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				loginPolicyProvider: &mockLoginPolicy{
					policy: &query.LoginPolicy{
						SecondFactorCheckLifetime: database.Duration(18 * time.Hour),
					},
				},
				userEventProvider:    &mockEventUser{},
				orgViewProvider:      &mockViewOrg{State: domain.OrgStateActive},
				idpUserLinksProvider: &mockIDPUserLinks{},
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
					PasswordVerification:     testNow.Add(-5 * time.Minute),
					SecondFactorVerification: testNow.Add(-5 * time.Minute),
				},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
					MFAMaxSetUp:     int32(domain.MFALevelSecondFactor),
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{},
			},
			args{
				&domain.AuthRequest{
					UserID:              "UserID",
					SelectedIDPConfigID: "IDPConfigID",
					LinkingUsers:        []*domain.ExternalUser{{IDPConfigID: "IDPConfigID", ExternalUserID: "UserID", DisplayName: "DisplayName"}},
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
						SecondFactorCheckLifetime: 18 * time.Hour,
						PasswordCheckLifetime:     10 * 24 * time.Hour,
					},
				}, false},
			[]domain.NextStep{&domain.LinkUsersStep{}},
			nil,
		},
		{
			"local auth requested (passwordless and password set up), passwordless step",
			fields{
				userSessionViewProvider: &mockViewUserSession{},
				userViewProvider: &mockViewUser{
					PasswordSet:        true,
					IsEmailVerified:    true,
					MFAMaxSetUp:        int32(domain.MFALevelMultiFactor),
					PasswordlessTokens: user_view_model.WebAuthNTokens{&user_view_model.WebAuthNView{ID: "id", State: int32(user_model.MFAStateReady)}},
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{
					idps: []*query.IDPUserLink{{IDPID: "IDPConfigID"}},
				},
			},
			args{
				&domain.AuthRequest{
					UserID:              "UserID",
					SelectedIDPConfigID: "IDPConfigID",
					LoginPolicy: &domain.LoginPolicy{
						PasswordlessType: domain.PasswordlessTypeAllowed,
					},
					AllowedExternalIDPs: []*domain.IDPProvider{
						{
							IDPConfigID: "IDPConfigID",
						},
					},
					RequestLocalAuth: true,
				}, false},
			[]domain.NextStep{
				&domain.PasswordlessStep{
					PasswordSet: true,
				},
			},
			nil,
		},
		{
			"local auth requested (password set up), password step",
			fields{
				userSessionViewProvider: &mockViewUserSession{},
				userViewProvider: &mockViewUser{
					PasswordSet:     true,
					IsEmailVerified: true,
				},
				userEventProvider: &mockEventUser{},
				orgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				lockoutPolicyProvider: &mockLockoutPolicy{
					policy: &query.LockoutPolicy{
						ShowFailures: true,
					},
				},
				idpUserLinksProvider: &mockIDPUserLinks{
					idps: []*query.IDPUserLink{{IDPID: "IDPConfigID"}},
				},
			},
			args{
				&domain.AuthRequest{
					UserID:              "UserID",
					SelectedIDPConfigID: "IDPConfigID",
					LoginPolicy: &domain.LoginPolicy{
						PasswordlessType: domain.PasswordlessTypeAllowed,
					},
					AllowedExternalIDPs: []*domain.IDPProvider{
						{
							IDPConfigID: "IDPConfigID",
						},
					},
					RequestLocalAuth: true,
				}, false},
			[]domain.NextStep{
				&domain.PasswordStep{},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AuthRequestRepo{
				AuthRequests:              tt.fields.AuthRequests,
				View:                      tt.fields.View,
				UserSessionViewProvider:   tt.fields.userSessionViewProvider,
				UserViewProvider:          tt.fields.userViewProvider,
				UserEventProvider:         tt.fields.userEventProvider,
				OrgViewProvider:           tt.fields.orgViewProvider,
				UserGrantProvider:         tt.fields.userGrantProvider,
				ProjectProvider:           tt.fields.projectProvider,
				ApplicationProvider:       tt.fields.applicationProvider,
				LoginPolicyViewProvider:   tt.fields.loginPolicyProvider,
				LockoutPolicyViewProvider: tt.fields.lockoutPolicyProvider,
				IDPUserLinksProvider:      tt.fields.idpUserLinksProvider,
				PrivacyPolicyProvider:     tt.fields.privacyPolicyProvider,
				LabelPolicyProvider:       tt.fields.labelPolicyProvider,
				PasswordAgePolicyProvider: tt.fields.passwordAgePolicyProvider,
				CustomTextProvider:        tt.fields.customTextProvider,
			}
			if tt.fields.passwordReset != nil {
				repo.PasswordReset = tt.fields.passwordReset(t)
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
	type args struct {
		userSession *user_model.UserSessionView
		request     *domain.AuthRequest
		user        *user_model.UserView
		isInternal  bool
	}
	tests := []struct {
		name            string
		args            args
		want            domain.NextStep
		wantChecked     bool
		errFunc         func(err error) bool
		wantMFAVerified []domain.MFAType
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
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						ForceMFA:            true,
						MFAInitSkipLifetime: 30 * 24 * time.Hour,
					},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: domain.MFALevelNotSetUp,
					},
				},
				isInternal: true,
			},
			nil,
			false,
			zerrors.IsPreconditionFailed,
			nil,
		},
		{
			"not set up, no mfas configured, no prompt and true",
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						MFAInitSkipLifetime: 30 * 24 * time.Hour,
					},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: domain.MFALevelNotSetUp,
					},
				},
				isInternal: true,
			},
			nil,
			true,
			nil,
			nil,
		},
		{
			"not set up, prompt and false",
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors:       []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
						MFAInitSkipLifetime: 30 * 24 * time.Hour,
					},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: domain.MFALevelNotSetUp,
					},
				},
				isInternal: true,
			},
			&domain.MFAPromptStep{
				MFAProviders: []domain.MFAType{
					domain.MFATypeTOTP,
				},
			},
			false,
			nil,
			nil,
		},
		{
			"not set up, forced by org, true",
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						ForceMFA:            true,
						SecondFactors:       []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
						MFAInitSkipLifetime: 30 * 24 * time.Hour,
					},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: domain.MFALevelNotSetUp,
					},
				},
				isInternal: true,
			},
			&domain.MFAPromptStep{
				Required: true,
				MFAProviders: []domain.MFAType{
					domain.MFATypeTOTP,
				},
			},
			false,
			nil,
			nil,
		},
		{
			"not set up and skipped, true",
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						MFAInitSkipLifetime: 30 * 24 * time.Hour,
					},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp:    domain.MFALevelNotSetUp,
						MFAInitSkipped: testNow,
					},
				},
				isInternal: true,
			},
			nil,
			true,
			nil,
			nil,
		},
		{
			"checked second factor, true",
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
						SecondFactorCheckLifetime: 18 * time.Hour,
					},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: domain.MFALevelSecondFactor,
						OTPState:    user_model.MFAStateReady,
					},
				},
				userSession: &user_model.UserSessionView{SecondFactorVerification: testNow.Add(-5 * time.Hour)},
				isInternal:  true,
			},
			nil,
			true,
			nil,
			[]domain.MFAType{domain.MFATypeTOTP},
		},
		{
			"checked passwordless, true",
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
						SecondFactorCheckLifetime: 18 * time.Hour,
						MultiFactors:              []domain.MultiFactorType{domain.MultiFactorTypeU2FWithPIN},
						MultiFactorCheckLifetime:  18 * time.Hour,
					},
					MFAsVerified: []domain.MFAType{domain.MFATypeU2FUserVerification},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: domain.MFALevelMultiFactor,
						PasswordlessTokens: []*user_model.WebAuthNView{
							{
								TokenID: "tokenID",
								State:   user_model.MFAStateReady,
							},
						},
					},
				},
				userSession: &user_model.UserSessionView{PasswordlessVerification: testNow.Add(-5 * time.Hour)},
				isInternal:  true,
			},
			nil,
			true,
			nil,
			[]domain.MFAType{domain.MFATypeU2FUserVerification},
		},
		{
			"not checked, check and false",
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
						SecondFactorCheckLifetime: 18 * time.Hour,
					},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: domain.MFALevelSecondFactor,
						OTPState:    user_model.MFAStateReady,
					},
				},
				userSession: &user_model.UserSessionView{},
				isInternal:  true,
			},

			&domain.MFAVerificationStep{
				MFAProviders: []domain.MFAType{domain.MFATypeTOTP},
			},
			false,
			nil,
			nil,
		},
		{
			"external not checked or forced but set up, want step",
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
						SecondFactorCheckLifetime: 18 * time.Hour,
					},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: domain.MFALevelSecondFactor,
						OTPState:    user_model.MFAStateReady,
					},
				},
				userSession: &user_model.UserSessionView{},
				isInternal:  false,
			},
			&domain.MFAVerificationStep{
				MFAProviders: []domain.MFAType{domain.MFATypeTOTP},
			},
			false,
			nil,
			nil,
		},
		{
			"external not forced but checked",
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
						SecondFactorCheckLifetime: 18 * time.Hour,
					},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: domain.MFALevelSecondFactor,
						OTPState:    user_model.MFAStateReady,
					},
				},
				userSession: &user_model.UserSessionView{SecondFactorVerification: testNow.Add(-5 * time.Hour)},
				isInternal:  false,
			},
			nil,
			true,
			nil,
			[]domain.MFAType{domain.MFATypeTOTP},
		},
		{
			"external not checked but required, want step",
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
						SecondFactorCheckLifetime: 18 * time.Hour,
						ForceMFA:                  true,
					},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: domain.MFALevelNotSetUp,
					},
				},
				userSession: &user_model.UserSessionView{},
				isInternal:  false,
			},
			&domain.MFAPromptStep{
				Required: true,
				MFAProviders: []domain.MFAType{
					domain.MFATypeTOTP,
				},
			},
			false,
			nil,
			nil,
		},
		{
			"external not checked but local required",
			args{
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						SecondFactors:             []domain.SecondFactorType{domain.SecondFactorTypeTOTP},
						SecondFactorCheckLifetime: 18 * time.Hour,
						ForceMFA:                  true,
						ForceMFALocalOnly:         true,
					},
				},
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: domain.MFALevelNotSetUp,
					},
				},
				userSession: &user_model.UserSessionView{},
				isInternal:  false,
			},
			nil,
			true,
			nil,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AuthRequestRepo{}
			got, ok, err := repo.mfaChecked(tt.args.userSession, tt.args.request, tt.args.user, tt.args.isInternal)
			if (tt.errFunc != nil && !tt.errFunc(err)) || (err != nil && tt.errFunc == nil) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if ok != tt.wantChecked {
				t.Errorf("mfaChecked() checked = %v, want %v", ok, tt.wantChecked)
			}
			assert.Equal(t, tt.want, got)
			assert.ElementsMatch(t, tt.args.request.MFAsVerified, tt.wantMFAVerified)
		})
	}
}

func TestAuthRequestRepo_mfaSkippedOrSetUp(t *testing.T) {
	type fields struct {
		MFAInitSkippedLifeTime time.Duration
	}
	type args struct {
		user    *user_model.UserView
		request *domain.AuthRequest
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
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp: domain.MFALevelSecondFactor,
					},
				},
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{},
				},
			},
			true,
		},
		{
			"mfa skipped active, true",
			fields{},
			args{
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp:    -1,
						MFAInitSkipped: testNow.Add(-10 * time.Hour),
					},
				},
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						MFAInitSkipLifetime: 30 * 24 * time.Hour,
					},
				},
			},
			true,
		},
		{
			"mfa skipped inactive, false",
			fields{},
			args{
				user: &user_model.UserView{
					HumanView: &user_model.HumanView{
						MFAMaxSetUp:    -1,
						MFAInitSkipped: testNow.Add(-40 * 24 * time.Hour),
					},
				},
				request: &domain.AuthRequest{
					LoginPolicy: &domain.LoginPolicy{
						MFAInitSkipLifetime: 30 * 24 * time.Hour,
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AuthRequestRepo{}
			if got := repo.mfaSkippedOrSetUp(tt.args.user, tt.args.request); got != tt.want {
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
				userID:        "id",
			},
			&user_model.UserSessionView{UserID: "id"},
			nil,
		},
		{
			"internal error, internal error",
			args{
				userProvider: &mockViewErrUserSession{},
				userID:       "id",
			},
			nil,
			zerrors.IsInternal,
		},
		{
			"error user events, old view model state",
			args{
				userProvider: &mockViewUserSession{
					PasswordVerification: testNow,
				},
				userID:        "id",
				eventProvider: &mockEventErrUser{},
			},
			&user_model.UserSessionView{
				PasswordVerification:     testNow,
				SecondFactorVerification: time.Time{},
				MultiFactorVerification:  time.Time{},
			},
			nil,
		},
		{
			"new user events but error, old view model state",
			args{
				userProvider: &mockViewUserSession{
					PasswordVerification: testNow,
				},
				agentID: "agentID",
				userID:  "id",
				eventProvider: &mockEventUser{
					Events: []eventstore.Event{
						&es_models.Event{
							AggregateType: user_repo.AggregateType,
							Typ:           user_repo.UserV1MFAOTPCheckSucceededType,
							CreationDate:  testNow,
						},
					},
				},
			},
			&user_model.UserSessionView{
				PasswordVerification:     testNow,
				SecondFactorVerification: time.Time{},
				MultiFactorVerification:  time.Time{},
			},
			nil,
		},
		{
			"new user events but other agentID, old view model state",
			args{
				userProvider: &mockViewUserSession{
					PasswordVerification: testNow,
				},
				agentID: "agentID",
				userID:  "id",
				eventProvider: &mockEventUser{
					Events: []eventstore.Event{
						&es_models.Event{
							AggregateType: user_repo.AggregateType,
							Typ:           user_repo.UserV1MFAOTPCheckSucceededType,
							CreationDate:  testNow,
							Data: func() []byte {
								data, _ := json.Marshal(&user_es_model.AuthRequest{UserAgentID: "otherID"})
								return data
							}(),
						},
					},
				},
			},
			&user_model.UserSessionView{
				PasswordVerification:     testNow,
				SecondFactorVerification: time.Time{},
				MultiFactorVerification:  time.Time{},
			},
			nil,
		},
		{
			"new user events, new view model state",
			args{
				userProvider: &mockViewUserSession{
					PasswordVerification: testNow,
				},
				agentID: "agentID",
				userID:  "id",
				eventProvider: &mockEventUser{
					Events: []eventstore.Event{
						&es_models.Event{
							AggregateType: user_repo.AggregateType,
							Typ:           user_repo.UserV1MFAOTPCheckSucceededType,
							CreationDate:  testNow,
							Data: func() []byte {
								data, _ := json.Marshal(&user_es_model.AuthRequest{UserAgentID: "agentID"})
								return data
							}(),
						},
					},
				},
			},
			&user_model.UserSessionView{
				PasswordVerification:     testNow,
				SecondFactorVerification: testNow,
				ChangeDate:               testNow,
			},
			nil,
		},
		{
			"new user events (user deleted), session terminated",
			args{
				userProvider: &mockViewUserSession{
					PasswordVerification: testNow,
				},
				agentID: "agentID",
				userID:  "id",
				eventProvider: &mockEventUser{
					Events: []eventstore.Event{
						&es_models.Event{
							AggregateType: user_repo.AggregateType,
							Typ:           user_repo.UserRemovedType,
							CreationDate:  testNow,
						},
					},
				},
			},
			&user_model.UserSessionView{
				ChangeDate: testNow,
				State:      domain.UserSessionStateTerminated,
			},
			nil,
		},
		{
			"new user events (user deleted, readded and password checked)",
			args{
				userProvider: &mockViewUserSession{
					PasswordVerification: testNow,
				},
				agentID: "agentID",
				userID:  "id",
				eventProvider: &mockEventUser{
					Events: []eventstore.Event{
						&es_models.Event{
							AggregateType: user_repo.AggregateType,
							Typ:           user_repo.UserRemovedType,
						},
						&es_models.Event{
							AggregateType: user_repo.AggregateType,
							Typ:           user_repo.HumanAddedType,
						},
						&es_models.Event{
							AggregateType: user_repo.AggregateType,
							Typ:           user_repo.HumanPasswordCheckSucceededType,
							CreationDate:  testNow,
							Data: func() []byte {
								data, _ := json.Marshal(&user_es_model.AuthRequest{UserAgentID: "agentID"})
								return data
							}(),
						},
					},
				},
			},
			&user_model.UserSessionView{
				ChangeDate:           testNow,
				PasswordVerification: testNow,
				State:                domain.UserSessionStateActive,
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
			zerrors.IsNotFound,
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
			"new user events, new view model state",
			args{
				viewProvider: &mockViewUser{
					PasswordSet:            true,
					PasswordChangeRequired: true,
				},
				eventProvider: &mockEventUser{
					Events: []eventstore.Event{
						&es_models.Event{
							AggregateType: user_repo.AggregateType,
							Typ:           user_repo.UserV1PasswordChangedType,
							CreationDate:  testNow,
							Data: func() []byte {
								data, _ := json.Marshal(user_es_model.Password{ChangeRequired: false, Secret: &crypto.CryptoValue{}})
								return data
							}(),
						},
					},
				},
			},
			&user_model.UserView{
				ChangeDate: testNow,
				State:      user_model.UserStateActive,
				UserName:   "UserName",
				HumanView: &user_model.HumanView{
					PasswordSet:            true,
					PasswordChangeRequired: false,
					PasswordChanged:        testNow,
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

func TestAuthRequestRepo_VerifyPassword_IgnoreUnknownUsernames(t *testing.T) {
	authRequest := func(userID string) *domain.AuthRequest {
		a := &domain.AuthRequest{
			ID:      "authRequestID",
			AgentID: "userAgentID",
			UserID:  userID,
			LoginPolicy: &domain.LoginPolicy{
				ObjectRoot:            es_models.ObjectRoot{},
				Default:               true,
				AllowUsernamePassword: true,
				AllowRegister:         true,
				AllowExternalIDP:      true,
				IDPProviders: []*domain.IDPProvider{
					{
						ObjectRoot:  es_models.ObjectRoot{},
						Type:        domain.IdentityProviderTypeSystem,
						IDPConfigID: "idpConfig1",
						Name:        "IdP",
						IDPType:     domain.IDPTypeOIDC,
						IDPState:    domain.IDPConfigStateActive,
					},
				},
				IgnoreUnknownUsernames: true,
			},
			AllowedExternalIDPs: []*domain.IDPProvider{
				{
					ObjectRoot:  es_models.ObjectRoot{},
					Type:        domain.IdentityProviderTypeSystem,
					IDPConfigID: "idpConfig1",
					Name:        "IdP",
					IDPType:     domain.IDPTypeOIDC,
					IDPState:    domain.IDPConfigStateActive,
				},
			},
			LabelPolicy: &domain.LabelPolicy{
				ObjectRoot: es_models.ObjectRoot{},
				State:      domain.LabelPolicyStateActive,
				Default:    true,
			},
			PrivacyPolicy: &domain.PrivacyPolicy{
				ObjectRoot: es_models.ObjectRoot{},
				State:      domain.PolicyStateActive,
				Default:    true,
			},
			LockoutPolicy: &domain.LockoutPolicy{
				Default: true,
			},
			PasswordAgePolicy: &domain.PasswordAgePolicy{
				ObjectRoot:     es_models.ObjectRoot{},
				MaxAgeDays:     0,
				ExpireWarnDays: 0,
			},
			DefaultTranslations: []*domain.CustomText{{}},
			OrgTranslations:     []*domain.CustomText{{}},
			SAMLRequestID:       "",
		}
		a.SetPolicyOrgID("instance1")
		return a
	}

	type fields struct {
		AuthRequests      func(*testing.T, string) cache.AuthRequestCache
		UserViewProvider  userViewProvider
		UserEventProvider userEventProvider
		OrgViewProvider   orgViewProvider
		PasswordChecker   passwordChecker
	}
	type args struct {
		ctx           context.Context
		authReqID     string
		userID        string
		resourceOwner string
		password      string
		userAgentID   string
		info          *domain.BrowserInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "no user",
			fields: fields{
				AuthRequests: func(tt *testing.T, userID string) cache.AuthRequestCache {
					m := mock.NewMockAuthRequestCache(gomock.NewController(tt))
					a := authRequest(userID)
					m.EXPECT().GetAuthRequestByID(gomock.Any(), "authRequestID").Return(a, nil)
					m.EXPECT().CacheAuthRequest(gomock.Any(), a)
					return m
				},
				UserViewProvider:  &mockViewNoUser{},
				UserEventProvider: &mockEventUser{},
			},
			args: args{
				ctx:           authz.NewMockContext("instance1", "", ""),
				authReqID:     "authRequestID",
				userID:        unknownUserID,
				resourceOwner: "org1",
				password:      "password",
				userAgentID:   "userAgentID",
				info: &domain.BrowserInfo{
					UserAgent: "useragent",
				},
			},
		},
		{
			name: "invalid password",
			fields: fields{
				AuthRequests: func(tt *testing.T, userID string) cache.AuthRequestCache {
					m := mock.NewMockAuthRequestCache(gomock.NewController(tt))
					a := authRequest(userID)
					m.EXPECT().GetAuthRequestByID(gomock.Any(), "authRequestID").Return(a, nil)
					m.EXPECT().CacheAuthRequest(gomock.Any(), a)
					return m
				},
				UserViewProvider:  &mockViewUser{},
				UserEventProvider: &mockEventUser{},
				OrgViewProvider:   &mockViewOrg{State: domain.OrgStateActive},
				PasswordChecker: &mockPasswordChecker{
					err: command.ErrPasswordInvalid(nil),
				},
			},
			args: args{
				ctx:           authz.NewMockContext("instance1", "", ""),
				authReqID:     "authRequestID",
				userID:        "user1",
				resourceOwner: "org1",
				password:      "password",
				userAgentID:   "userAgentID",
				info: &domain.BrowserInfo{
					UserAgent: "useragent",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &AuthRequestRepo{
				AuthRequests:      tt.fields.AuthRequests(t, tt.args.userID),
				UserViewProvider:  tt.fields.UserViewProvider,
				UserEventProvider: tt.fields.UserEventProvider,
				OrgViewProvider:   tt.fields.OrgViewProvider,
				PasswordChecker:   tt.fields.PasswordChecker,
			}
			err := repo.VerifyPassword(tt.args.ctx, tt.args.authReqID, tt.args.userID, tt.args.resourceOwner, tt.args.password, tt.args.userAgentID, tt.args.info)
			assert.ErrorIs(t, err, zerrors.ThrowInvalidArgument(nil, "EVENT-SDe2f", "Errors.User.UsernameOrPassword.Invalid"))
		})
	}
}
