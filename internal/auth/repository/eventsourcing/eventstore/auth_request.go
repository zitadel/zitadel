package eventstore

import (
	"context"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/auth_request/model"
	cache "github.com/caos/zitadel/internal/auth_request/repository"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/id"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_event "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	org_view_model "github.com/caos/zitadel/internal/org/repository/view/model"
	project_view_model "github.com/caos/zitadel/internal/project/repository/view/model"
	"github.com/caos/zitadel/internal/tracing"
	user_model "github.com/caos/zitadel/internal/user/model"
	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	user_view_model "github.com/caos/zitadel/internal/user/repository/view/model"
	grant_view_model "github.com/caos/zitadel/internal/usergrant/repository/view/model"
)

type AuthRequestRepo struct {
	UserEvents   *user_event.UserEventstore
	OrgEvents    *org_event.OrgEventstore
	AuthRequests cache.AuthRequestCache
	View         *view.View

	UserSessionViewProvider userSessionViewProvider
	UserViewProvider        userViewProvider
	UserEventProvider       userEventProvider
	OrgViewProvider         orgViewProvider
	LoginPolicyViewProvider loginPolicyViewProvider
	IDPProviderViewProvider idpProviderViewProvider
	UserGrantProvider       userGrantProvider

	IdGenerator id.Generator

	PasswordCheckLifeTime      time.Duration
	ExternalLoginCheckLifeTime time.Duration
	MfaInitSkippedLifeTime     time.Duration
	SecondFactorCheckLifeTime  time.Duration
	MultiFactorCheckLifeTime   time.Duration

	IAMID string
}

type userSessionViewProvider interface {
	UserSessionByIDs(string, string) (*user_view_model.UserSessionView, error)
	UserSessionsByAgentID(string) ([]*user_view_model.UserSessionView, error)
}
type userViewProvider interface {
	UserByID(string) (*user_view_model.UserView, error)
}

type loginPolicyViewProvider interface {
	LoginPolicyByAggregateID(string) (*iam_view_model.LoginPolicyView, error)
}

type idpProviderViewProvider interface {
	IDPProvidersByAggregateIDAndState(string, iam_model.IDPConfigState) ([]*iam_view_model.IDPProviderView, error)
}

type userEventProvider interface {
	UserEventsByID(ctx context.Context, id string, sequence uint64) ([]*es_models.Event, error)
	BulkAddExternalIDPs(ctx context.Context, userID string, externalIDPs []*user_model.ExternalIDP) error
}

type orgViewProvider interface {
	OrgByID(string) (*org_view_model.OrgView, error)
	OrgByPrimaryDomain(string) (*org_view_model.OrgView, error)
}

type userGrantProvider interface {
	ApplicationByClientID(context.Context, string) (*project_view_model.ApplicationView, error)
	UserGrantsByProjectAndUserID(string, string) ([]*grant_view_model.UserGrantView, error)
}

func (repo *AuthRequestRepo) Health(ctx context.Context) error {
	if err := repo.UserEvents.Health(ctx); err != nil {
		return err
	}
	return repo.AuthRequests.Health(ctx)
}

func (repo *AuthRequestRepo) CreateAuthRequest(ctx context.Context, request *model.AuthRequest) (_ *model.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	reqID, err := repo.IdGenerator.Next()
	if err != nil {
		return nil, err
	}
	request.ID = reqID
	app, err := repo.View.ApplicationByClientID(ctx, request.ApplicationID)
	if err != nil {
		return nil, err
	}
	appIDs, err := repo.View.AppIDsFromProjectID(ctx, app.ProjectID)
	if err != nil {
		return nil, err
	}
	request.Audience = appIDs
	request.AppendAudIfNotExisting(app.ProjectID)
	if request.LoginHint != "" {
		err = repo.checkLoginName(ctx, request, request.LoginHint)
		logging.LogWithFields("EVENT-aG311", "login name", request.LoginHint, "id", request.ID, "applicationID", request.ApplicationID, "traceID", tracing.TraceIDFromCtx(ctx)).OnError(err).Debug("login hint invalid")
	}
	err = repo.AuthRequests.SaveAuthRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (repo *AuthRequestRepo) AuthRequestByID(ctx context.Context, id, userAgentID string) (_ *model.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return repo.getAuthRequestNextSteps(ctx, id, userAgentID, false)
}

func (repo *AuthRequestRepo) AuthRequestByIDCheckLoggedIn(ctx context.Context, id, userAgentID string) (_ *model.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return repo.getAuthRequestNextSteps(ctx, id, userAgentID, true)
}

func (repo *AuthRequestRepo) SaveAuthCode(ctx context.Context, id, code, userAgentID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, id, userAgentID)
	if err != nil {
		return err
	}
	request.Code = code
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) AuthRequestByCode(ctx context.Context, code string) (_ *model.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.AuthRequests.GetAuthRequestByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	err = repo.fillLoginPolicy(ctx, request)
	if err != nil {
		return nil, err
	}
	steps, err := repo.nextSteps(ctx, request, true)
	if err != nil {
		return nil, err
	}
	request.PossibleSteps = steps
	return request, nil
}

func (repo *AuthRequestRepo) DeleteAuthRequest(ctx context.Context, id string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return repo.AuthRequests.DeleteAuthRequest(ctx, id)
}

func (repo *AuthRequestRepo) CheckLoginName(ctx context.Context, id, loginName, userAgentID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, id, userAgentID)
	if err != nil {
		return err
	}
	err = repo.checkLoginName(ctx, request, loginName)
	if err != nil {
		return err
	}
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) SelectExternalIDP(ctx context.Context, authReqID, idpConfigID, userAgentID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	err = repo.checkSelectedExternalIDP(request, idpConfigID)
	if err != nil {
		return err
	}
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) CheckExternalUserLogin(ctx context.Context, authReqID, userAgentID string, externalUser *model.ExternalUser, info *model.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	err = repo.checkExternalUserLogin(request, externalUser.IDPConfigID, externalUser.ExternalUserID)
	if errors.IsNotFound(err) {
		if err := repo.setLinkingUser(ctx, request, externalUser); err != nil {
			return err
		}
		return err
	}
	if err != nil {
		return err
	}

	err = repo.UserEvents.ExternalLoginChecked(ctx, request.UserID, request.WithCurrentInfo(info))
	if err != nil {
		return err
	}
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) setLinkingUser(ctx context.Context, request *model.AuthRequest, externalUser *model.ExternalUser) error {
	request.LinkingUsers = append(request.LinkingUsers, externalUser)
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) SelectUser(ctx context.Context, id, userID, userAgentID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, id, userAgentID)
	if err != nil {
		return err
	}
	user, err := activeUserByID(ctx, repo.UserViewProvider, repo.UserEventProvider, repo.OrgViewProvider, userID)
	if err != nil {
		return err
	}
	request.SetUserInfo(user.ID, user.PreferredLoginName, user.DisplayName, user.ResourceOwner)
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) VerifyPassword(ctx context.Context, id, userID, password, userAgentID string, info *model.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, id, userAgentID)
	if err != nil {
		return err
	}
	if request.UserID != userID {
		return errors.ThrowPreconditionFailed(nil, "EVENT-ds35D", "Errors.User.NotMatchingUserID")
	}
	return repo.UserEvents.CheckPassword(ctx, userID, password, request.WithCurrentInfo(info))
}

func (repo *AuthRequestRepo) VerifyMfaOTP(ctx context.Context, authRequestID, userID, code, userAgentID string, info *model.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, authRequestID, userAgentID)
	if err != nil {
		return err
	}
	if request.UserID != userID {
		return errors.ThrowPreconditionFailed(nil, "EVENT-ADJ26", "Errors.User.NotMatchingUserID")
	}
	return repo.UserEvents.CheckMfaOTP(ctx, userID, code, request.WithCurrentInfo(info))
}

func (repo *AuthRequestRepo) LinkExternalUsers(ctx context.Context, authReqID, userAgentID string, info *model.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	err = linkExternalIDPs(ctx, repo.UserEventProvider, request)
	if err != nil {
		return err
	}
	err = repo.UserEvents.ExternalLoginChecked(ctx, request.UserID, request.WithCurrentInfo(info))
	if err != nil {
		return err
	}
	request.LinkingUsers = nil
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) ResetLinkingUsers(ctx context.Context, authReqID, userAgentID string) error {
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	request.LinkingUsers = nil
	request.SelectedIDPConfigID = ""
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) AutoRegisterExternalUser(ctx context.Context, registerUser *user_model.User, externalIDP *user_model.ExternalIDP, orgMember *org_model.OrgMember, authReqID, userAgentID, resourceOwner string, info *model.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	policyResourceOwner := authz.GetCtxData(ctx).OrgID
	if resourceOwner != "" {
		policyResourceOwner = resourceOwner
	}
	pwPolicy, err := repo.View.PasswordComplexityPolicyByAggregateID(policyResourceOwner)
	if errors.IsNotFound(err) {
		pwPolicy, err = repo.View.PasswordComplexityPolicyByAggregateID(repo.IAMID)
	}
	if err != nil {
		return err
	}
	pwPolicyView := iam_es_model.PasswordComplexityViewToModel(pwPolicy)
	orgPolicy, err := repo.View.OrgIAMPolicyByAggregateID(policyResourceOwner)
	if errors.IsNotFound(err) {
		orgPolicy, err = repo.View.OrgIAMPolicyByAggregateID(repo.IAMID)
	}
	if err != nil {
		return err
	}
	orgPolicyView := iam_es_model.OrgIAMViewToModel(orgPolicy)
	user, aggregates, err := repo.UserEvents.PrepareRegisterUser(ctx, registerUser, externalIDP, pwPolicyView, orgPolicyView, resourceOwner)
	if err != nil {
		return err
	}
	if orgMember != nil {
		orgMember.UserID = user.AggregateID
		_, memberAggregate, err := repo.OrgEvents.PrepareAddOrgMember(ctx, orgMember, policyResourceOwner)
		if err != nil {
			return err
		}
		aggregates = append(aggregates, memberAggregate)
	}

	err = sdk.PushAggregates(ctx, repo.UserEvents.PushAggregates, user.AppendEvents, aggregates...)
	if err != nil {
		return err
	}
	request.UserID = user.AggregateID
	request.UserOrgID = user.ResourceOwner
	request.SelectedIDPConfigID = externalIDP.IDPConfigID
	request.LinkingUsers = nil
	err = repo.UserEvents.ExternalLoginChecked(ctx, request.UserID, request.WithCurrentInfo(info))
	if err != nil {
		return err
	}
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) getAuthRequestNextSteps(ctx context.Context, id, userAgentID string, checkLoggedIn bool) (*model.AuthRequest, error) {
	request, err := repo.getAuthRequest(ctx, id, userAgentID)
	if err != nil {
		return nil, err
	}
	steps, err := repo.nextSteps(ctx, request, checkLoggedIn)
	if err != nil {
		return nil, err
	}
	request.PossibleSteps = steps
	return request, nil
}

func (repo *AuthRequestRepo) getAuthRequest(ctx context.Context, id, userAgentID string) (*model.AuthRequest, error) {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if request.AgentID != userAgentID {
		return nil, errors.ThrowPermissionDenied(nil, "EVENT-adk13", "Errors.AuthRequest.UserAgentNotCorresponding")
	}
	err = repo.fillLoginPolicy(ctx, request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (repo *AuthRequestRepo) getLoginPolicyAndIDPProviders(ctx context.Context, orgID string) (*iam_model.LoginPolicyView, []*iam_model.IDPProviderView, error) {
	policy, err := repo.getLoginPolicy(ctx, orgID)
	if err != nil {
		return nil, nil, err
	}
	if !policy.AllowExternalIDP {
		return policy, nil, nil
	}
	idpProviders, err := getLoginPolicyIDPProviders(repo.IDPProviderViewProvider, repo.IAMID, orgID, policy.Default)
	if err != nil {
		return nil, nil, err
	}
	return policy, idpProviders, nil
}

func (repo *AuthRequestRepo) fillLoginPolicy(ctx context.Context, request *model.AuthRequest) error {
	orgID := request.UserOrgID
	if orgID == "" {
		primaryDomain := request.GetScopeOrgPrimaryDomain()
		if primaryDomain != "" {
			org, err := repo.GetOrgByPrimaryDomain(primaryDomain)
			if err != nil {
				return err
			}
			orgID = org.ID
		}
	}
	if orgID == "" {
		orgID = repo.IAMID
	}

	policy, idpProviders, err := repo.getLoginPolicyAndIDPProviders(ctx, orgID)
	if err != nil {
		return err
	}
	request.LoginPolicy = policy
	if idpProviders != nil {
		request.AllowedExternalIDPs = idpProviders
	}
	return nil
}

func (repo *AuthRequestRepo) checkLoginName(ctx context.Context, request *model.AuthRequest, loginName string) (err error) {
	primaryDomain := request.GetScopeOrgPrimaryDomain()
	orgID := ""
	if primaryDomain != "" {
		org, err := repo.GetOrgByPrimaryDomain(primaryDomain)
		if err != nil {
			return err
		}
		orgID = org.ID
	}

	user := new(user_view_model.UserView)
	if orgID != "" {
		user, err = repo.View.UserByLoginNameAndResourceOwner(loginName, orgID)
	} else {
		user, err = repo.View.UserByLoginName(loginName)
		if err == nil {
			err = repo.checkLoginPolicyWithResourceOwner(ctx, request, user)
			if err != nil {
				return err
			}
		}
	}
	if err != nil {
		return err
	}

	request.SetUserInfo(user.ID, loginName, "", user.ResourceOwner)
	return nil
}

func (repo AuthRequestRepo) GetOrgByPrimaryDomain(primaryDomain string) (*org_model.OrgView, error) {
	org, err := repo.OrgViewProvider.OrgByPrimaryDomain(primaryDomain)
	if err != nil {
		return nil, err
	}
	return org_view_model.OrgToModel(org), nil
}

func (repo AuthRequestRepo) checkLoginPolicyWithResourceOwner(ctx context.Context, request *model.AuthRequest, user *user_view_model.UserView) error {
	loginPolicy, idpProviders, err := repo.getLoginPolicyAndIDPProviders(ctx, user.ResourceOwner)
	if err != nil {
		return err
	}
	if len(request.LinkingUsers) != 0 && !loginPolicy.AllowExternalIDP {
		return errors.ThrowInvalidArgument(nil, "LOGIN-s9sio", "Errors.User.NotAllowedToLink")
	}
	if len(request.LinkingUsers) != 0 {
		exists := linkingIDPConfigExistingInAllowedIDPs(request.LinkingUsers, idpProviders)
		if !exists {
			return errors.ThrowInvalidArgument(nil, "LOGIN-Dj89o", "Errors.User.NotAllowedToLink")
		}
	}
	request.LoginPolicy = loginPolicy
	request.AllowedExternalIDPs = idpProviders
	return nil
}

func (repo *AuthRequestRepo) checkSelectedExternalIDP(request *model.AuthRequest, idpConfigID string) error {
	for _, externalIDP := range request.AllowedExternalIDPs {
		if externalIDP.IDPConfigID == idpConfigID {
			request.SelectedIDPConfigID = idpConfigID
			return nil
		}
	}
	return errors.ThrowNotFound(nil, "LOGIN-Nsm8r", "Errors.User.ExternalIDP.NotAllowed")
}

func (repo *AuthRequestRepo) checkExternalUserLogin(request *model.AuthRequest, idpConfigID, externalUserID string) (err error) {
	primaryDomain := request.GetScopeOrgPrimaryDomain()
	externalIDP := new(user_view_model.ExternalIDPView)
	org := new(org_model.OrgView)
	if primaryDomain != "" {
		org, err = repo.GetOrgByPrimaryDomain(primaryDomain)
		if err != nil {
			return err
		}
		externalIDP, err = repo.View.ExternalIDPByExternalUserIDAndIDPConfigIDAndResourceOwner(externalUserID, idpConfigID, org.ID)
	} else {
		externalIDP, err = repo.View.ExternalIDPByExternalUserIDAndIDPConfigID(externalUserID, idpConfigID)
	}
	if err != nil {
		return err
	}
	request.SetUserInfo(externalIDP.UserID, "", "", externalIDP.ResourceOwner)
	return nil
}

func (repo *AuthRequestRepo) nextSteps(ctx context.Context, request *model.AuthRequest, checkLoggedIn bool) ([]model.NextStep, error) {
	if request == nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-ds27a", "Errors.Internal")
	}
	steps := make([]model.NextStep, 0)
	if !checkLoggedIn && request.Prompt == model.PromptNone {
		return append(steps, &model.RedirectToCallbackStep{}), nil
	}
	if request.UserID == "" {
		if request.LinkingUsers != nil && len(request.LinkingUsers) > 0 {
			steps = append(steps, new(model.ExternalNotFoundOptionStep))
			return steps, nil
		}
		steps = append(steps, new(model.LoginStep))
		if request.Prompt == model.PromptSelectAccount || request.Prompt == model.PromptUnspecified {
			users, err := repo.usersForUserSelection(request)
			if err != nil {
				return nil, err
			}
			if len(users) > 0 || request.Prompt == model.PromptSelectAccount {
				steps = append(steps, &model.SelectUserStep{Users: users})
			}
		}
		return steps, nil
	}
	user, err := activeUserByID(ctx, repo.UserViewProvider, repo.UserEventProvider, repo.OrgViewProvider, request.UserID)
	if err != nil {
		return nil, err
	}
	request.LoginName = user.PreferredLoginName
	userSession, err := userSessionByIDs(ctx, repo.UserSessionViewProvider, repo.UserEventProvider, request.AgentID, user)
	if err != nil {
		return nil, err
	}

	if (request.SelectedIDPConfigID != "" || userSession.SelectedIDPConfigID != "") && (request.LinkingUsers == nil || len(request.LinkingUsers) == 0) {
		if !checkVerificationTime(userSession.ExternalLoginVerification, repo.ExternalLoginCheckLifeTime) {
			selectedIDPConfigID := request.SelectedIDPConfigID
			if selectedIDPConfigID == "" {
				selectedIDPConfigID = userSession.SelectedIDPConfigID
			}
			return append(steps, &model.ExternalLoginStep{SelectedIDPConfigID: selectedIDPConfigID}), nil
		}
	} else if (request.SelectedIDPConfigID == "" && userSession.SelectedIDPConfigID == "") || (request.SelectedIDPConfigID != "" && request.LinkingUsers != nil && len(request.LinkingUsers) > 0) {
		if user.InitRequired {
			return append(steps, &model.InitUserStep{PasswordSet: user.PasswordSet}), nil
		}
		if !user.PasswordSet {
			return append(steps, &model.InitPasswordStep{}), nil
		}

		if !checkVerificationTime(userSession.PasswordVerification, repo.PasswordCheckLifeTime) {
			return append(steps, &model.PasswordStep{}), nil
		}
		request.PasswordVerified = true
		request.AuthTime = userSession.PasswordVerification
	}

	step, ok, err := repo.mfaChecked(userSession, request, user)
	if err != nil {
		return nil, err
	}
	if !ok {
		return append(steps, step), nil
	}

	if user.PasswordChangeRequired {
		steps = append(steps, &model.ChangePasswordStep{})
	}
	if !user.IsEmailVerified {
		steps = append(steps, &model.VerifyEMailStep{})
	}
	if user.UsernameChangeRequired {
		steps = append(steps, &model.ChangeUsernameStep{})
	}

	if user.PasswordChangeRequired || !user.IsEmailVerified || user.UsernameChangeRequired {
		return steps, nil
	}

	if request.LinkingUsers != nil && len(request.LinkingUsers) != 0 {
		return append(steps, &model.LinkUsersStep{}), nil

	}
	//PLANNED: consent step

	missing, err := userGrantRequired(ctx, request, user, repo.UserGrantProvider)
	if err != nil {
		return nil, err
	}
	if missing {
		return append(steps, &model.GrantRequiredStep{}), nil
	}

	return append(steps, &model.RedirectToCallbackStep{}), nil
}

func (repo *AuthRequestRepo) usersForUserSelection(request *model.AuthRequest) ([]model.UserSelection, error) {
	userSessions, err := userSessionsByUserAgentID(repo.UserSessionViewProvider, request.AgentID)
	if err != nil {
		return nil, err
	}
	users := make([]model.UserSelection, len(userSessions))
	for i, session := range userSessions {
		users[i] = model.UserSelection{
			UserID:           session.UserID,
			DisplayName:      session.DisplayName,
			LoginName:        session.LoginName,
			UserSessionState: session.State,
		}
	}
	return users, nil
}

func (repo *AuthRequestRepo) mfaChecked(userSession *user_model.UserSessionView, request *model.AuthRequest, user *user_model.UserView) (model.NextStep, bool, error) {
	mfaLevel := request.MfaLevel()
	allowedProviders, required := user.MfaTypesAllowed(mfaLevel, request.LoginPolicy)
	promptRequired := (user.MfaMaxSetUp < mfaLevel) || (len(allowedProviders) == 0 && required)
	if promptRequired || !repo.mfaSkippedOrSetUp(user) {
		types := user.MfaTypesSetupPossible(mfaLevel, request.LoginPolicy)
		if promptRequired && len(types) == 0 {
			return nil, false, errors.ThrowPreconditionFailed(nil, "LOGIN-5Hm8s", "Errors.Login.LoginPolicy.MFA.ForceAndNotConfigured")
		}
		if len(types) == 0 {
			return nil, true, nil
		}
		return &model.MfaPromptStep{
			Required:     promptRequired,
			MfaProviders: types,
		}, false, nil
	}
	switch mfaLevel {
	default:
		fallthrough
	case model.MFALevelNotSetUp:
		if len(allowedProviders) == 0 {
			return nil, true, nil
		}
		fallthrough
	case model.MFALevelSecondFactor:
		if checkVerificationTime(userSession.SecondFactorVerification, repo.SecondFactorCheckLifeTime) {
			request.MfasVerified = append(request.MfasVerified, userSession.SecondFactorVerificationType)
			request.AuthTime = userSession.SecondFactorVerification
			return nil, true, nil
		}
		fallthrough
	case model.MFALevelMultiFactor:
		if checkVerificationTime(userSession.MultiFactorVerification, repo.MultiFactorCheckLifeTime) {
			request.MfasVerified = append(request.MfasVerified, userSession.MultiFactorVerificationType)
			request.AuthTime = userSession.MultiFactorVerification
			return nil, true, nil
		}
	}
	return &model.MfaVerificationStep{
		MfaProviders: allowedProviders,
	}, false, nil
}

func (repo *AuthRequestRepo) mfaSkippedOrSetUp(user *user_model.UserView) bool {
	if user.MfaMaxSetUp > model.MFALevelNotSetUp {
		return true
	}
	return checkVerificationTime(user.MfaInitSkipped, repo.MfaInitSkippedLifeTime)
}

func (repo *AuthRequestRepo) getLoginPolicy(ctx context.Context, orgID string) (*iam_model.LoginPolicyView, error) {
	policy, err := repo.View.LoginPolicyByAggregateID(orgID)
	if errors.IsNotFound(err) {
		policy, err = repo.View.LoginPolicyByAggregateID(repo.IAMID)
		if err != nil {
			return nil, err
		}
		policy.Default = true
	}
	if err != nil {
		return nil, err
	}
	return iam_es_model.LoginPolicyViewToModel(policy), err
}

func getLoginPolicyIDPProviders(provider idpProviderViewProvider, iamID, orgID string, defaultPolicy bool) ([]*iam_model.IDPProviderView, error) {
	if defaultPolicy {
		idpProviders, err := provider.IDPProvidersByAggregateIDAndState(iamID, iam_model.IDPConfigStateActive)
		if err != nil {
			return nil, err
		}
		return iam_es_model.IDPProviderViewsToModel(idpProviders), nil
	}
	idpProviders, err := provider.IDPProvidersByAggregateIDAndState(orgID, iam_model.IDPConfigStateActive)
	if err != nil {
		return nil, err
	}
	return iam_es_model.IDPProviderViewsToModel(idpProviders), nil
}

func checkVerificationTime(verificationTime time.Time, lifetime time.Duration) bool {
	return verificationTime.Add(lifetime).After(time.Now().UTC())
}

func userSessionsByUserAgentID(provider userSessionViewProvider, agentID string) ([]*user_model.UserSessionView, error) {
	session, err := provider.UserSessionsByAgentID(agentID)
	if err != nil {
		return nil, err
	}
	return user_view_model.UserSessionsToModel(session), nil
}

func userSessionByIDs(ctx context.Context, provider userSessionViewProvider, eventProvider userEventProvider, agentID string, user *user_model.UserView) (*user_model.UserSessionView, error) {
	session, err := provider.UserSessionByIDs(agentID, user.ID)
	if err != nil {
		if !errors.IsNotFound(err) {
			return nil, err
		}
		session = &user_view_model.UserSessionView{}
	}
	events, err := eventProvider.UserEventsByID(ctx, user.ID, session.Sequence)
	if err != nil {
		logging.Log("EVENT-Hse6s").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("error retrieving new events")
		return user_view_model.UserSessionToModel(session), nil
	}
	sessionCopy := *session
	for _, event := range events {
		switch event.Type {
		case es_model.UserPasswordCheckSucceeded,
			es_model.UserPasswordCheckFailed,
			es_model.MFAOTPCheckSucceeded,
			es_model.MFAOTPCheckFailed,
			es_model.SignedOut,
			es_model.UserLocked,
			es_model.UserDeactivated,
			es_model.HumanPasswordCheckSucceeded,
			es_model.HumanPasswordCheckFailed,
			es_model.HumanExternalLoginCheckSucceeded,
			es_model.HumanMFAOTPCheckSucceeded,
			es_model.HumanMFAOTPCheckFailed,
			es_model.HumanSignedOut:
			eventData, err := user_view_model.UserSessionFromEvent(event)
			if err != nil {
				logging.Log("EVENT-sdgT3").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("error getting event data")
				return user_view_model.UserSessionToModel(session), nil
			}
			if eventData.UserAgentID != agentID {
				continue
			}
		case es_model.UserRemoved:
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dG2fe", "Errors.User.NotActive")
		}
		sessionCopy.AppendEvent(event)
	}
	return user_view_model.UserSessionToModel(&sessionCopy), nil
}

func activeUserByID(ctx context.Context, userViewProvider userViewProvider, userEventProvider userEventProvider, orgViewProvider orgViewProvider, userID string) (*user_model.UserView, error) {
	user, err := userByID(ctx, userViewProvider, userEventProvider, userID)
	if err != nil {
		return nil, err
	}

	if user.HumanView == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Lm69x", "Errors.User.NotHuman")
	}

	if user.State == user_model.UserStateLocked || user.State == user_model.UserStateSuspend {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-FJ262", "Errors.User.Locked")
	}
	if !(user.State == user_model.UserStateActive || user.State == user_model.UserStateInitial) {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-FJ262", "Errors.User.NotActive")
	}
	org, err := orgViewProvider.OrgByID(user.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if org.State != int32(org_model.OrgStateActive) {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Zws3s", "Errors.User.NotActive")
	}
	return user, nil
}

func userByID(ctx context.Context, viewProvider userViewProvider, eventProvider userEventProvider, userID string) (*user_model.UserView, error) {
	user, viewErr := viewProvider.UserByID(userID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	} else if user == nil {
		user = new(user_view_model.UserView)
	}
	events, err := eventProvider.UserEventsByID(ctx, userID, user.Sequence)
	if err != nil {
		logging.Log("EVENT-dfg42").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("error retrieving new events")
		return user_view_model.UserToModel(user), nil
	}
	if len(events) == 0 {
		if viewErr != nil {
			return nil, viewErr
		}
		return user_view_model.UserToModel(user), viewErr
	}
	userCopy := *user
	for _, event := range events {
		if err := userCopy.AppendEvent(event); err != nil {
			return user_view_model.UserToModel(user), nil
		}
	}
	if userCopy.State == int32(user_model.UserStateDeleted) {
		return nil, errors.ThrowNotFound(nil, "EVENT-3F9so", "Errors.User.NotFound")
	}
	return user_view_model.UserToModel(&userCopy), nil
}

func linkExternalIDPs(ctx context.Context, userEventProvider userEventProvider, request *model.AuthRequest) error {
	externalIDPs := make([]*user_model.ExternalIDP, len(request.LinkingUsers))
	for i, linkingUser := range request.LinkingUsers {
		externalIDP := &user_model.ExternalIDP{
			ObjectRoot:  es_models.ObjectRoot{AggregateID: request.UserID},
			IDPConfigID: linkingUser.IDPConfigID,
			UserID:      linkingUser.ExternalUserID,
			DisplayName: linkingUser.DisplayName,
		}
		externalIDPs[i] = externalIDP
	}
	data := authz.CtxData{
		UserID: "LOGIN",
		OrgID:  request.UserOrgID,
	}
	return userEventProvider.BulkAddExternalIDPs(authz.SetCtxData(ctx, data), request.UserID, externalIDPs)
}

func linkingIDPConfigExistingInAllowedIDPs(linkingUsers []*model.ExternalUser, idpProviders []*iam_model.IDPProviderView) bool {
	for _, linkingUser := range linkingUsers {
		exists := false
		for _, idp := range idpProviders {
			if idp.IDPConfigID == linkingUser.IDPConfigID {
				exists = true
				continue
			}
		}
		if !exists {
			return false
		}
	}
	return true
}
func userGrantRequired(ctx context.Context, request *model.AuthRequest, user *user_model.UserView, userGrantProvider userGrantProvider) (_ bool, err error) {
	var app *project_view_model.ApplicationView
	switch request.Request.Type() {
	case model.AuthRequestTypeOIDC:
		app, err = userGrantProvider.ApplicationByClientID(ctx, request.ApplicationID)
		if err != nil {
			return false, err
		}
	default:
		return false, errors.ThrowPreconditionFailed(nil, "EVENT-dfrw2", "Errors.AuthRequest.RequestTypeNotSupported")
	}
	if !app.ProjectRoleCheck {
		return false, nil
	}
	grants, err := userGrantProvider.UserGrantsByProjectAndUserID(app.ProjectID, user.ID)
	if err != nil {
		return false, err
	}
	return len(grants) == 0, nil
}
