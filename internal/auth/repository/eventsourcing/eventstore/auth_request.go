package eventstore

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	cache "github.com/zitadel/zitadel/internal/auth_request/repository"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/query"
	user_repo "github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	user_model "github.com/zitadel/zitadel/internal/user/model"
	user_view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const unknownUserID = "UNKNOWN"

var (
	ErrUserNotFound = func(err error) error {
		return zerrors.ThrowNotFound(err, "EVENT-hodc6", "Errors.User.NotFound")
	}
)

type AuthRequestRepo struct {
	Command      *command.Commands
	Query        *query.Queries
	AuthRequests cache.AuthRequestCache
	View         *view.View
	UserCodeAlg  crypto.EncryptionAlgorithm

	LabelPolicyProvider       labelPolicyProvider
	UserSessionViewProvider   userSessionViewProvider
	UserViewProvider          userViewProvider
	UserCommandProvider       userCommandProvider
	UserEventProvider         userEventProvider
	OrgViewProvider           orgViewProvider
	LoginPolicyViewProvider   loginPolicyViewProvider
	LockoutPolicyViewProvider lockoutPolicyViewProvider
	PasswordAgePolicyProvider passwordAgePolicyProvider
	PrivacyPolicyProvider     privacyPolicyProvider
	IDPProviderViewProvider   idpProviderViewProvider
	IDPUserLinksProvider      idpUserLinksProvider
	UserGrantProvider         userGrantProvider
	ProjectProvider           projectProvider
	ApplicationProvider       applicationProvider
	CustomTextProvider        customTextProvider
	PasswordReset             passwordReset
	PasswordChecker           passwordChecker

	IdGenerator id.Generator
}

type labelPolicyProvider interface {
	ActiveLabelPolicyByOrg(context.Context, string, bool) (*query.LabelPolicy, error)
}

type privacyPolicyProvider interface {
	PrivacyPolicyByOrg(context.Context, bool, string, bool) (*query.PrivacyPolicy, error)
}

type userSessionViewProvider interface {
	UserSessionByIDs(context.Context, string, string, string) (*user_view_model.UserSessionView, error)
	UserSessionsByAgentID(context.Context, string, string) ([]*user_view_model.UserSessionView, error)
	GetLatestUserSessionSequence(ctx context.Context, instanceID string) (*query.CurrentState, error)
}

type userViewProvider interface {
	UserByID(context.Context, string, string) (*user_view_model.UserView, error)
}

type loginPolicyViewProvider interface {
	LoginPolicyByID(context.Context, bool, string, bool) (*query.LoginPolicy, error)
}

type lockoutPolicyViewProvider interface {
	LockoutPolicyByOrg(context.Context, bool, string) (*query.LockoutPolicy, error)
}

type passwordAgePolicyProvider interface {
	PasswordAgePolicyByOrg(context.Context, bool, string, bool) (*query.PasswordAgePolicy, error)
}

type idpProviderViewProvider interface {
	IDPLoginPolicyLinks(context.Context, string, *query.IDPLoginPolicyLinksSearchQuery, bool) (*query.IDPLoginPolicyLinks, error)
}

type idpUserLinksProvider interface {
	IDPUserLinks(ctx context.Context, queries *query.IDPUserLinksSearchQuery, permissionCheck domain.PermissionCheck) (*query.IDPUserLinks, error)
}

type userEventProvider interface {
	UserEventsByID(ctx context.Context, id string, changeDate time.Time, eventTypes []eventstore.EventType) ([]eventstore.Event, error)
	PasswordCodeExists(ctx context.Context, userID string) (exists bool, err error)
	InviteCodeExists(ctx context.Context, userID string) (exists bool, err error)
}

type userCommandProvider interface {
	BulkAddedUserIDPLinks(ctx context.Context, userID, resourceOwner string, externalIDPs []*domain.UserIDPLink) error
}

type orgViewProvider interface {
	OrgByID(context.Context, string) (*query.Org, error)
	OrgByPrimaryDomain(context.Context, string) (*query.Org, error)
}

type userGrantProvider interface {
	ProjectByClientID(context.Context, string) (*query.Project, error)
	UserGrantsByProjectAndUserID(context.Context, string, string) ([]*query.UserGrant, error)
}

type projectProvider interface {
	ProjectByClientID(context.Context, string) (*query.Project, error)
	SearchProjectGrants(ctx context.Context, queries *query.ProjectGrantSearchQueries, permissionCheck domain.PermissionCheck) (projects *query.ProjectGrants, err error)
}

type applicationProvider interface {
	AppByOIDCClientID(context.Context, string) (*query.App, error)
}

type customTextProvider interface {
	CustomTextListByTemplate(ctx context.Context, aggregateID string, text string, withOwnerRemoved bool) (texts *query.CustomTexts, err error)
}

type passwordReset interface {
	RequestSetPassword(ctx context.Context, userID, resourceOwner string, notifyType domain.NotificationType, authRequestID string) (objectDetails *domain.ObjectDetails, err error)
}

type passwordChecker interface {
	HumanCheckPassword(ctx context.Context, resourceOwner, userID, password string, authReq *domain.AuthRequest) error
}

func (repo *AuthRequestRepo) Health(ctx context.Context) error {
	return repo.AuthRequests.Health(ctx)
}

func (repo *AuthRequestRepo) CreateAuthRequest(ctx context.Context, request *domain.AuthRequest) (_ *domain.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	reqID, err := repo.IdGenerator.Next()
	if err != nil {
		return nil, err
	}
	request.ID = reqID
	project, err := repo.ProjectProvider.ProjectByClientID(ctx, request.ApplicationID)
	if err != nil {
		return nil, err
	}
	request.AppendAudIfNotExisting(project.ID)
	request.ApplicationResourceOwner = project.ResourceOwner
	request.PrivateLabelingSetting = project.PrivateLabelingSetting
	if err := setOrgID(ctx, repo.OrgViewProvider, request); err != nil {
		return nil, err
	}
	if request.LoginHint != "" {
		err = repo.checkLoginName(ctx, request, request.LoginHint)
		logging.WithFields("login name", request.LoginHint, "id", request.ID, "applicationID", request.ApplicationID, "traceID", tracing.TraceIDFromCtx(ctx)).OnError(err).Info("login hint invalid")
	}
	if request.UserID == "" && request.LoginHint == "" && domain.IsPrompt(request.Prompt, domain.PromptNone) {
		err = repo.tryUsingOnlyUserSession(ctx, request)
		logging.WithFields("id", request.ID, "applicationID", request.ApplicationID, "traceID", tracing.TraceIDFromCtx(ctx)).OnError(err).Debug("unable to select only user session")
	}

	err = repo.AuthRequests.SaveAuthRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (repo *AuthRequestRepo) AuthRequestByID(ctx context.Context, id, userAgentID string) (_ *domain.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return repo.getAuthRequestNextSteps(ctx, id, userAgentID, false)
}

func (repo *AuthRequestRepo) AuthRequestByIDCheckLoggedIn(ctx context.Context, id, userAgentID string) (_ *domain.AuthRequest, err error) {
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

func (repo *AuthRequestRepo) SaveSAMLRequestID(ctx context.Context, id, requestID, userAgentID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, id, userAgentID)
	if err != nil {
		return err
	}
	request.SAMLRequestID = requestID
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) AuthRequestByCode(ctx context.Context, code string) (_ *domain.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.AuthRequests.GetAuthRequestByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	err = repo.fillPolicies(ctx, request)
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

func (repo *AuthRequestRepo) SelectExternalIDP(ctx context.Context, authReqID, idpConfigID, userAgentID string, idpArguments map[string]any) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	err = repo.checkSelectedExternalIDP(request, idpConfigID, idpArguments)
	if err != nil {
		return err
	}
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) CheckExternalUserLogin(ctx context.Context, authReqID, userAgentID string, externalUser *domain.ExternalUser, info *domain.BrowserInfo, migrationCheck bool) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	err = repo.checkExternalUserLogin(ctx, request, externalUser.IDPConfigID, externalUser.ExternalUserID)
	if zerrors.IsNotFound(err) {
		// clear potential user information (e.g. when username was entered but another external user was returned)
		request.SetUserInfo("", "", "", "", "", request.UserOrgID)
		// in case the check was done with an ID, that was retrieved by a session that allows migration,
		// we do not need to set the linking user and return early
		if migrationCheck {
			return err
		}
		if err := repo.setLinkingUser(ctx, request, externalUser); err != nil {
			return err
		}
		return err
	}
	if err != nil {
		return err
	}

	request.IDPLoginChecked = true
	err = repo.Command.UserIDPLoginChecked(ctx, request.UserOrgID, request.UserID, request.WithCurrentInfo(info))
	if err != nil {
		return err
	}
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) SetExternalUserLogin(ctx context.Context, authReqID, userAgentID string, externalUser *domain.ExternalUser) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}

	err = repo.setLinkingUser(ctx, request, externalUser)
	if err != nil {
		return err
	}
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) SetLinkingUser(ctx context.Context, request *domain.AuthRequest, externalUser *domain.ExternalUser) error {
	for i, user := range request.LinkingUsers {
		if user.ExternalUserID == externalUser.ExternalUserID {
			request.LinkingUsers[i] = externalUser
			return repo.AuthRequests.UpdateAuthRequest(ctx, request)
		}
	}
	return nil
}

func (repo *AuthRequestRepo) setLinkingUser(ctx context.Context, request *domain.AuthRequest, externalUser *domain.ExternalUser) error {
	request.LinkingUsers = append(request.LinkingUsers, externalUser)
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) SelectUser(ctx context.Context, authReqID, userID, userAgentID string, enforceExistingSession bool) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	// Check if the session already exists in the same user agent, e.g. when selecting the user from the user selection page.
	// This is to prevent username enumeration attacks by checking if the user exists in the system.
	if enforceExistingSession {
		userSession, err := userSessionByIDs(ctx, repo.UserSessionViewProvider, repo.UserEventProvider, request.AgentID, userID)
		if err != nil {
			return err
		}
		if userSession.Sequence == 0 {
			return zerrors.ThrowNotFound(nil, "AUTH-2d3f4", "Errors.UserSession.NotFound")
		}
	}
	user, err := activeUserByID(ctx, repo.UserViewProvider, repo.UserEventProvider, repo.OrgViewProvider, repo.LockoutPolicyViewProvider, userID, false)
	if err != nil {
		return err
	}
	if request.RequestedOrgID != "" && request.RequestedOrgID != user.ResourceOwner {
		return zerrors.ThrowPreconditionFailed(nil, "EVENT-fJe2a", "Errors.User.NotAllowedOrg")
	}
	username := user.UserName
	if request.RequestedOrgID == "" {
		username = user.PreferredLoginName
	}
	request.SetUserInfo(user.ID, username, user.PreferredLoginName, user.DisplayName, user.AvatarKey, user.ResourceOwner)
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) VerifyPassword(ctx context.Context, authReqID, userID, resourceOwner, password, userAgentID string, info *domain.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequestEnsureUser(ctx, authReqID, userAgentID, userID)
	if err != nil {
		if isIgnoreUserNotFoundError(err, request) {
			// use the same errorID as below (otherwise it would expose the error reason)
			return zerrors.ThrowInvalidArgument(nil, "EVENT-SDe2f", "Errors.User.UsernameOrPassword.Invalid")
		}
		return err
	}
	err = repo.PasswordChecker.HumanCheckPassword(ctx, resourceOwner, userID, password, request.WithCurrentInfo(info))
	if isIgnoreUserInvalidPasswordError(err, request) {
		// use the same errorID as above (otherwise it would expose the error reason)
		return zerrors.ThrowInvalidArgument(nil, "EVENT-SDe2f", "Errors.User.UsernameOrPassword.Invalid")
	}
	return err
}

func isIgnoreUserNotFoundError(err error, request *domain.AuthRequest) bool {
	return request != nil && request.LoginPolicy != nil && request.LoginPolicy.IgnoreUnknownUsernames && errors.Is(err, ErrUserNotFound(nil))
}

func isIgnoreUserInvalidPasswordError(err error, request *domain.AuthRequest) bool {
	return request != nil && request.LoginPolicy != nil && request.LoginPolicy.IgnoreUnknownUsernames && errors.Is(err, command.ErrPasswordInvalid(nil))
}

func lockoutPolicyToDomain(policy *query.LockoutPolicy) *domain.LockoutPolicy {
	return &domain.LockoutPolicy{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:   policy.ID,
			Sequence:      policy.Sequence,
			ResourceOwner: policy.ResourceOwner,
			CreationDate:  policy.CreationDate,
			ChangeDate:    policy.ChangeDate,
		},
		Default:             policy.IsDefault,
		MaxPasswordAttempts: policy.MaxPasswordAttempts,
		MaxOTPAttempts:      policy.MaxOTPAttempts,
		ShowLockOutFailures: policy.ShowFailures,
	}
}

func (repo *AuthRequestRepo) VerifyMFAOTP(ctx context.Context, authRequestID, userID, resourceOwner, code, userAgentID string, info *domain.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequestEnsureUser(ctx, authRequestID, userAgentID, userID)
	if err != nil {
		return err
	}
	return repo.Command.HumanCheckMFATOTP(ctx, userID, code, resourceOwner, request.WithCurrentInfo(info))
}

func (repo *AuthRequestRepo) SendMFAOTPSMS(ctx context.Context, userID, resourceOwner, authRequestID, userAgentID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	request, err := repo.getAuthRequestEnsureUser(ctx, authRequestID, userAgentID, userID)
	if err != nil {
		return err
	}
	return repo.Command.HumanSendOTPSMS(ctx, userID, resourceOwner, request)
}

func (repo *AuthRequestRepo) VerifyMFAOTPSMS(ctx context.Context, userID, resourceOwner, code, authRequestID, userAgentID string, info *domain.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequestEnsureUser(ctx, authRequestID, userAgentID, userID)
	if err != nil {
		return err
	}
	return repo.Command.HumanCheckOTPSMS(ctx, userID, code, resourceOwner, request.WithCurrentInfo(info))
}

func (repo *AuthRequestRepo) SendMFAOTPEmail(ctx context.Context, userID, resourceOwner, authRequestID, userAgentID string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	request, err := repo.getAuthRequestEnsureUser(ctx, authRequestID, userAgentID, userID)
	if err != nil {
		return err
	}
	return repo.Command.HumanSendOTPEmail(ctx, userID, resourceOwner, request)
}

func (repo *AuthRequestRepo) VerifyMFAOTPEmail(ctx context.Context, userID, resourceOwner, code, authRequestID, userAgentID string, info *domain.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequestEnsureUser(ctx, authRequestID, userAgentID, userID)
	if err != nil {
		return err
	}
	return repo.Command.HumanCheckOTPEmail(ctx, userID, code, resourceOwner, request.WithCurrentInfo(info))
}

func (repo *AuthRequestRepo) BeginMFAU2FLogin(ctx context.Context, userID, resourceOwner, authRequestID, userAgentID string) (login *domain.WebAuthNLogin, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	request, err := repo.getAuthRequestEnsureUser(ctx, authRequestID, userAgentID, userID)
	if err != nil {
		return nil, err
	}
	return repo.Command.HumanBeginU2FLogin(ctx, userID, resourceOwner, request)
}

func (repo *AuthRequestRepo) VerifyMFAU2F(ctx context.Context, userID, resourceOwner, authRequestID, userAgentID string, credentialData []byte, info *domain.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequestEnsureUser(ctx, authRequestID, userAgentID, userID)
	if err != nil {
		return err
	}
	return repo.Command.HumanFinishU2FLogin(ctx, userID, resourceOwner, credentialData, request)
}

func (repo *AuthRequestRepo) BeginPasswordlessSetup(ctx context.Context, userID, resourceOwner string, authenticatorPlatform domain.AuthenticatorAttachment) (login *domain.WebAuthNToken, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	return repo.Command.HumanAddPasswordlessSetup(ctx, userID, resourceOwner, authenticatorPlatform)
}

func (repo *AuthRequestRepo) VerifyPasswordlessSetup(ctx context.Context, userID, resourceOwner, userAgentID, tokenName string, credentialData []byte) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	_, err = repo.Command.HumanHumanPasswordlessSetup(ctx, userID, resourceOwner, tokenName, userAgentID, credentialData)
	return err
}

func (repo *AuthRequestRepo) BeginPasswordlessInitCodeSetup(ctx context.Context, userID, resourceOwner, codeID, verificationCode string, preferredPlatformType domain.AuthenticatorAttachment) (login *domain.WebAuthNToken, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	passwordlessInitCode, err := repo.Query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypePasswordlessInitCode, repo.UserCodeAlg)
	if err != nil {
		return nil, err
	}
	return repo.Command.HumanAddPasswordlessSetupInitCode(ctx, userID, resourceOwner, codeID, verificationCode, preferredPlatformType, passwordlessInitCode)
}

func (repo *AuthRequestRepo) VerifyPasswordlessInitCodeSetup(ctx context.Context, userID, resourceOwner, userAgentID, tokenName, codeID, verificationCode string, credentialData []byte) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	passwordlessInitCode, err := repo.Query.InitEncryptionGenerator(ctx, domain.SecretGeneratorTypePasswordlessInitCode, repo.UserCodeAlg)
	if err != nil {
		return err
	}
	_, err = repo.Command.HumanPasswordlessSetupInitCode(ctx, userID, resourceOwner, tokenName, userAgentID, codeID, verificationCode, credentialData, passwordlessInitCode)
	return err
}

func (repo *AuthRequestRepo) BeginPasswordlessLogin(ctx context.Context, userID, resourceOwner, authRequestID, userAgentID string) (login *domain.WebAuthNLogin, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequestEnsureUser(ctx, authRequestID, userAgentID, userID)
	if err != nil {
		return nil, err
	}
	return repo.Command.HumanBeginPasswordlessLogin(ctx, userID, resourceOwner, request)
}

func (repo *AuthRequestRepo) VerifyPasswordless(ctx context.Context, userID, resourceOwner, authRequestID, userAgentID string, credentialData []byte, info *domain.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequestEnsureUser(ctx, authRequestID, userAgentID, userID)
	if err != nil {
		return err
	}
	return repo.Command.HumanFinishPasswordlessLogin(ctx, userID, resourceOwner, credentialData, request)
}

func (repo *AuthRequestRepo) LinkExternalUsers(ctx context.Context, authReqID, userAgentID string, info *domain.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	err = linkExternalIDPs(ctx, repo.UserCommandProvider, request)
	if err != nil {
		return err
	}
	err = repo.Command.UserIDPLoginChecked(ctx, request.UserOrgID, request.UserID, request.WithCurrentInfo(info))
	if err != nil {
		return err
	}
	request.LinkingUsers = nil
	request.IDPLoginChecked = true
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

func (repo *AuthRequestRepo) ResetSelectedIDP(ctx context.Context, authReqID, userAgentID string) error {
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	request.SelectedIDPConfigID = ""
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) RequestLocalAuth(ctx context.Context, authReqID, userAgentID string) error {
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	request.RequestLocalAuth = true
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) AutoRegisterExternalUser(ctx context.Context, registerUser *domain.Human, externalIDP *domain.UserIDPLink, orgMemberRoles []string, authReqID, userAgentID, resourceOwner string, metadatas []*domain.Metadata, info *domain.BrowserInfo) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	request, err := repo.getAuthRequest(ctx, authReqID, userAgentID)
	if err != nil {
		return err
	}
	addMetadata := make([]*command.AddMetadataEntry, len(metadatas))
	for i, metadata := range metadatas {
		addMetadata[i] = &command.AddMetadataEntry{
			Key:   metadata.Key,
			Value: metadata.Value,
		}
	}
	human := command.AddHumanFromDomain(registerUser, metadatas, request, externalIDP)
	err = repo.Command.AddUserHuman(ctx, resourceOwner, human, false, repo.UserCodeAlg)
	if err != nil {
		return err
	}
	request.SetUserInfo(human.ID, human.Username, human.Username, human.DisplayName, "", resourceOwner)
	request.SelectedIDPConfigID = externalIDP.IDPConfigID
	request.LinkingUsers = nil
	request.IDPLoginChecked = true
	err = repo.Command.UserIDPLoginChecked(ctx, request.UserOrgID, request.UserID, request.WithCurrentInfo(info))
	if err != nil {
		return err
	}
	if len(metadatas) > 0 {
		// user context necessary due to permission check in command
		userCtx := authz.SetCtxData(ctx, authz.CtxData{UserID: request.UserID, OrgID: request.UserOrgID})
		_, err := repo.Command.BulkSetUserMetadata(userCtx, request.UserID, request.UserOrgID, metadatas...)
		if err != nil {
			return err
		}
	}
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) getAuthRequestNextSteps(ctx context.Context, id, userAgentID string, checkLoggedIn bool) (*domain.AuthRequest, error) {
	request, err := repo.getAuthRequest(ctx, id, userAgentID)
	if err != nil {
		return nil, err
	}
	steps, err := repo.nextSteps(ctx, request, checkLoggedIn)
	if err != nil {
		return request, err
	}
	request.PossibleSteps = steps
	return request, nil
}

func (repo *AuthRequestRepo) getAuthRequestEnsureUser(ctx context.Context, authRequestID, userAgentID, userID string) (*domain.AuthRequest, error) {
	request, err := repo.getAuthRequest(ctx, authRequestID, userAgentID)
	if err != nil {
		return nil, err
	}
	// If there's no user, checks if the user could be reused (from the session).
	// (the nextStepsUser will update the userID in the request in that case)
	if request.UserID == "" {
		if _, err = repo.nextStepsUser(ctx, request); err != nil {
			return nil, err
		}
	}
	if request.UserID != userID {
		return nil, zerrors.ThrowPreconditionFailed(nil, "EVENT-GBH32", "Errors.User.NotMatchingUserID")
	}
	_, err = activeUserByID(ctx, repo.UserViewProvider, repo.UserEventProvider, repo.OrgViewProvider, repo.LockoutPolicyViewProvider, request.UserID, false)
	if err != nil {
		return request, err
	}
	return request, nil
}

func (repo *AuthRequestRepo) getAuthRequest(ctx context.Context, id, userAgentID string) (request *domain.AuthRequest, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	request, err = repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if request.AgentID != userAgentID {
		return nil, zerrors.ThrowPermissionDenied(nil, "EVENT-adk13", "Errors.AuthRequest.UserAgentNotCorresponding")
	}
	err = repo.fillPolicies(ctx, request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (repo *AuthRequestRepo) getLoginPolicyAndIDPProviders(ctx context.Context, orgID string) (*query.LoginPolicy, []*domain.IDPProvider, error) {
	policy, err := repo.LoginPolicyViewProvider.LoginPolicyByID(ctx, false, orgID, false)
	if err != nil {
		return nil, nil, err
	}
	if !policy.AllowExternalIDPs {
		return policy, nil, nil
	}
	idpProviders, err := getLoginPolicyIDPProviders(ctx, repo.IDPProviderViewProvider, authz.GetInstance(ctx).InstanceID(), orgID, policy.IsDefault)
	if err != nil {
		return nil, nil, err
	}
	return policy, idpProviders, nil
}

func (repo *AuthRequestRepo) fillPolicies(ctx context.Context, request *domain.AuthRequest) error {
	instance := authz.GetInstance(ctx)
	orgID := request.RequestedOrgID
	if orgID == "" {
		orgID = request.UserOrgID
	}
	if orgID == "" {
		orgID = authz.GetInstance(ctx).DefaultOrganisationID()
		if !instance.Features().LoginDefaultOrg {
			orgID = instance.InstanceID()
		}
	}

	if request.LoginPolicy == nil || len(request.AllowedExternalIDPs) == 0 || request.PolicyOrgID() != orgID {
		loginPolicy, idpProviders, err := repo.getLoginPolicyAndIDPProviders(ctx, orgID)
		if err != nil {
			return err
		}
		request.LoginPolicy = queryLoginPolicyToDomain(loginPolicy)
		if len(idpProviders) > 0 {
			request.AllowedExternalIDPs = idpProviders
		}
	}
	if request.LockoutPolicy == nil || request.PolicyOrgID() != orgID {
		lockoutPolicy, err := repo.getLockoutPolicy(ctx, orgID)
		if err != nil {
			return err
		}
		request.LockoutPolicy = lockoutPolicyToDomain(lockoutPolicy)
	}
	if request.PrivacyPolicy == nil || request.PolicyOrgID() != orgID {
		privacyPolicy, err := repo.GetPrivacyPolicy(ctx, orgID)
		if err != nil {
			return err
		}
		request.PrivacyPolicy = privacyPolicy
	}
	if request.LabelPolicy == nil || request.PolicyOrgID() != orgID {
		labelPolicy, err := repo.getLabelPolicy(ctx, request.PrivateLabelingOrgID(orgID))
		if err != nil {
			return err
		}
		request.LabelPolicy = labelPolicy
	}
	if request.PasswordAgePolicy == nil || request.PolicyOrgID() != orgID {
		passwordPolicy, err := repo.getPasswordAgePolicy(ctx, orgID)
		if err != nil {
			return err
		}
		request.PasswordAgePolicy = passwordPolicy
	}
	if len(request.DefaultTranslations) == 0 {
		defaultLoginTranslations, err := repo.getLoginTexts(ctx, instance.InstanceID())
		if err != nil {
			return err
		}
		request.DefaultTranslations = defaultLoginTranslations
	}
	if len(request.OrgTranslations) == 0 || request.PolicyOrgID() != orgID {
		orgLoginTranslations, err := repo.getLoginTexts(ctx, request.PrivateLabelingOrgID(orgID))
		if err != nil {
			return err
		}
		request.OrgTranslations = orgLoginTranslations
	}
	request.SetPolicyOrgID(orgID)
	repo.AuthRequests.CacheAuthRequest(ctx, request)
	return nil
}

func (repo *AuthRequestRepo) tryUsingOnlyUserSession(ctx context.Context, request *domain.AuthRequest) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	userSessions, err := userSessionsByUserAgentID(ctx, repo.UserSessionViewProvider, request.AgentID, request.InstanceID)
	if err != nil {
		return err
	}
	if len(userSessions) == 1 {
		user := userSessions[0]
		username := user.UserName
		if request.RequestedOrgID == "" {
			username = user.LoginName
		}
		request.SetUserInfo(user.UserID, username, user.LoginName, user.DisplayName, user.AvatarKey, user.ResourceOwner)
	}
	return nil
}

func (repo *AuthRequestRepo) checkLoginName(ctx context.Context, request *domain.AuthRequest, loginNameInput string) (err error) {
	var user *user_view_model.UserView
	loginNameInput = strings.TrimSpace(loginNameInput)
	preferredLoginName := loginNameInput
	if request.RequestedOrgID != "" {
		if request.RequestedOrgDomain {
			domainPolicy, err := repo.getDomainPolicy(ctx, request.RequestedOrgID)
			if err != nil {
				return err
			}
			if domainPolicy.UserLoginMustBeDomain {
				preferredLoginName += "@" + request.RequestedPrimaryDomain
			}
		}
		user, err = repo.checkLoginNameInputForResourceOwner(ctx, request, loginNameInput, preferredLoginName)
	} else {
		user, err = repo.checkLoginNameInput(ctx, request, loginNameInput, preferredLoginName)
	}
	// return any error apart from not found ones directly
	if err != nil && !zerrors.IsNotFound(err) {
		return err
	}
	// if there's an active (human) user, let's use it
	if user != nil && !user.HumanView.IsZero() && domain.UserState(user.State).IsEnabled() {
		request.SetUserInfo(user.ID, loginNameInput, preferredLoginName, "", "", user.ResourceOwner)
		return nil
	}
	// the user was either not found or not active
	// so check if the loginname suffix matches a verified org domain
	// but only if no org was requested (by id or domain)
	if request.RequestedOrgID == "" {
		ok, errDomainDiscovery := repo.checkDomainDiscovery(ctx, request, loginNameInput)
		if errDomainDiscovery != nil || ok {
			return errDomainDiscovery
		}
	}
	// let's once again check if the user was just inactive
	if user != nil && user.State == int32(domain.UserStateInactive) {
		return zerrors.ThrowPreconditionFailed(nil, "AUTH-2n8fs", "Errors.User.Inactive")
	}
	// or locked
	if user != nil && user.State == int32(domain.UserStateLocked) {
		return zerrors.ThrowPreconditionFailed(nil, "AUTH-SF3gb", "Errors.User.Locked")
	}
	// let's just check if unknown usernames are ignored
	if request.LoginPolicy != nil && request.LoginPolicy.IgnoreUnknownUsernames {
		if request.LabelPolicy != nil && request.LabelPolicy.HideLoginNameSuffix {
			preferredLoginName = loginNameInput
		}
		request.SetUserInfo(unknownUserID, preferredLoginName, preferredLoginName, preferredLoginName, "", request.RequestedOrgID)
		return nil
	}
	// there was no policy that allowed unknown loginnames in any case
	// so not found errors can now be returned
	if err != nil {
		return err
	}
	// let's check if it was a machine user
	if !user.MachineView.IsZero() {
		return zerrors.ThrowPreconditionFailed(nil, "AUTH-DGV4g", "Errors.User.NotHuman")
	}
	// everything should be handled by now
	logging.WithFields("authRequest", request.ID, "loginName", loginNameInput).Error("unhandled state for checkLoginName")
	return zerrors.ThrowInternal(nil, "AUTH-asf3df", "Errors.Internal")
}

func (repo *AuthRequestRepo) checkDomainDiscovery(ctx context.Context, request *domain.AuthRequest, loginName string) (bool, error) {
	// check if there's a suffix in the loginname
	loginName = strings.TrimSpace(strings.ToLower(loginName))
	index := strings.LastIndex(loginName, "@")
	if index < 0 {
		return false, nil
	}
	// check if the suffix matches a verified domain
	org, err := repo.Query.OrgByVerifiedDomain(ctx, loginName[index+1:])
	if err != nil {
		return false, nil
	}
	// and if the login policy allows domain discovery
	policy, err := repo.Query.LoginPolicyByID(ctx, true, org.ID, false)
	if err != nil || !policy.AllowDomainDiscovery {
		return false, nil
	}
	// discovery was allowed, so set the org as requested org
	// and clear all potentially existing user information and only set the loginname as hint (for registration)
	// also ensure that the policies are read from the org
	request.SetOrgInformation(org.ID, org.Name, org.Domain, false)
	request.SetUserInfo("", "", "", "", "", org.ID)
	if err = repo.fillPolicies(ctx, request); err != nil {
		return false, err
	}
	request.LoginHint = loginName
	request.Prompt = append(request.Prompt, domain.PromptCreate) // to trigger registration
	repo.AuthRequests.CacheAuthRequest(ctx, request)
	return true, nil
}

func (repo *AuthRequestRepo) checkLoginNameInput(ctx context.Context, request *domain.AuthRequest, loginNameInput, preferredLoginName string) (*user_view_model.UserView, error) {
	// always check the preferred / suffixed loginname first
	user, err := repo.View.UserByLoginName(ctx, preferredLoginName, request.InstanceID)
	if err == nil {
		// and take the user regardless if there would be a user with that email or phone
		return user, repo.checkLoginPolicyWithResourceOwner(ctx, request, user.ResourceOwner)
	}
	// for email and phone check we will use the loginname as provided by the user (without computed suffix)
	user, emailErr := repo.View.UserByEmail(ctx, loginNameInput, request.InstanceID)
	if emailErr == nil {
		// if there was a single user with the specified email
		// load and check the login policy
		if emailErr = repo.checkLoginPolicyWithResourceOwner(ctx, request, user.ResourceOwner); emailErr != nil {
			return nil, emailErr
		}
		// and in particular if the login with email is possible
		// if so take the user (and ignore possible phone matches)
		if !request.LoginPolicy.DisableLoginWithEmail {
			return user, nil
		}
	}
	user, phoneErr := repo.View.UserByPhone(ctx, loginNameInput, request.InstanceID)
	if phoneErr == nil {
		// if there was a single user with the specified phone
		// load and check the login policy
		if phoneErr = repo.checkLoginPolicyWithResourceOwner(ctx, request, user.ResourceOwner); phoneErr != nil {
			return nil, phoneErr
		}
		// and in particular if the login with phone is possible
		// if so take the user
		if !request.LoginPolicy.DisableLoginWithPhone {
			return user, nil
		}
	}
	// if we get here the user was not found by loginname
	// and either there was no match for email or phone as well, or they have been both disabled
	return nil, err
}

func (repo *AuthRequestRepo) checkLoginNameInputForResourceOwner(ctx context.Context, request *domain.AuthRequest, loginNameInput, preferredLoginName string) (*user_view_model.UserView, error) {
	// always check the preferred / suffixed loginname first
	user, err := repo.View.UserByLoginNameAndResourceOwner(ctx, preferredLoginName, request.RequestedOrgID, request.InstanceID)
	if err == nil {
		// and take the user regardless if there would be a user with that email or phone
		return user, nil
	}
	// for email and phone check we will use the loginname as provided by the user (without computed suffix)
	if request.LoginPolicy != nil && !request.LoginPolicy.DisableLoginWithEmail {
		// if login by email is allowed and there was a single user with the specified email
		// take that user (and ignore possible phone number matches)
		user, emailErr := repo.View.UserByEmailAndResourceOwner(ctx, loginNameInput, request.RequestedOrgID, request.InstanceID)
		if emailErr == nil {
			return user, nil
		}
	}
	if request.LoginPolicy != nil && !request.LoginPolicy.DisableLoginWithPhone {
		// if login by phone is allowed and there was a single user with the specified phone
		// take that user
		user, phoneErr := repo.View.UserByPhoneAndResourceOwner(ctx, loginNameInput, request.RequestedOrgID, request.InstanceID)
		if phoneErr == nil {
			return user, nil
		}
	}
	// if we get here the user was not found by loginname
	// and either there was no match for email or phone as well or they have been both disabled
	return nil, err
}

func (repo *AuthRequestRepo) checkLoginPolicyWithResourceOwner(ctx context.Context, request *domain.AuthRequest, resourceOwner string) (err error) {
	if request.LoginPolicy == nil || request.PolicyOrgID() != resourceOwner {
		loginPolicy, idps, err := repo.getLoginPolicyAndIDPProviders(ctx, resourceOwner)
		if err != nil {
			return err
		}
		request.LoginPolicy = queryLoginPolicyToDomain(loginPolicy)
		request.AllowedExternalIDPs = idps
	}
	if len(request.LinkingUsers) != 0 && !request.LoginPolicy.AllowExternalIDP {
		return zerrors.ThrowInvalidArgument(nil, "LOGIN-s9sio", "Errors.User.NotAllowedToLink")
	}
	if len(request.LinkingUsers) != 0 {
		exists := linkingIDPConfigExistingInAllowedIDPs(request.LinkingUsers, request.AllowedExternalIDPs)
		if !exists {
			return zerrors.ThrowInvalidArgument(nil, "LOGIN-Dj89o", "Errors.User.NotAllowedToLink")
		}
	}
	repo.AuthRequests.CacheAuthRequest(ctx, request)
	return nil
}

func queryLoginPolicyToDomain(policy *query.LoginPolicy) *domain.LoginPolicy {
	return &domain.LoginPolicy{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:   policy.OrgID,
			Sequence:      policy.Sequence,
			ResourceOwner: policy.OrgID,
			CreationDate:  policy.CreationDate,
			ChangeDate:    policy.ChangeDate,
		},
		Default:                    policy.IsDefault,
		AllowUsernamePassword:      policy.AllowUsernamePassword,
		AllowRegister:              policy.AllowRegister,
		AllowExternalIDP:           policy.AllowExternalIDPs,
		ForceMFA:                   policy.ForceMFA,
		ForceMFALocalOnly:          policy.ForceMFALocalOnly,
		SecondFactors:              policy.SecondFactors,
		MultiFactors:               policy.MultiFactors,
		PasswordlessType:           policy.PasswordlessType,
		HidePasswordReset:          policy.HidePasswordReset,
		IgnoreUnknownUsernames:     policy.IgnoreUnknownUsernames,
		AllowDomainDiscovery:       policy.AllowDomainDiscovery,
		DefaultRedirectURI:         policy.DefaultRedirectURI,
		PasswordCheckLifetime:      time.Duration(policy.PasswordCheckLifetime),
		ExternalLoginCheckLifetime: time.Duration(policy.ExternalLoginCheckLifetime),
		MFAInitSkipLifetime:        time.Duration(policy.MFAInitSkipLifetime),
		SecondFactorCheckLifetime:  time.Duration(policy.SecondFactorCheckLifetime),
		MultiFactorCheckLifetime:   time.Duration(policy.MultiFactorCheckLifetime),
		DisableLoginWithEmail:      policy.DisableLoginWithEmail,
		DisableLoginWithPhone:      policy.DisableLoginWithPhone,
	}
}

func (repo *AuthRequestRepo) checkSelectedExternalIDP(request *domain.AuthRequest, idpConfigID string, idpArguments map[string]any) error {
	for _, externalIDP := range request.AllowedExternalIDPs {
		if externalIDP.IDPConfigID == idpConfigID {
			request.SelectedIDPConfigID = idpConfigID
			request.SelectedIDPConfigArgs = idpArguments
			return nil
		}
	}
	return zerrors.ThrowNotFound(nil, "LOGIN-Nsm8r", "Errors.User.ExternalIDP.NotAllowed")
}

func (repo *AuthRequestRepo) checkExternalUserLogin(ctx context.Context, request *domain.AuthRequest, idpConfigID, externalUserID string) (err error) {
	idQuery, err := query.NewIDPUserLinkIDPIDSearchQuery(idpConfigID)
	if err != nil {
		return err
	}
	externalIDQuery, err := query.NewIDPUserLinksExternalIDSearchQuery(externalUserID)
	if err != nil {
		return err
	}
	queries := []query.SearchQuery{
		idQuery, externalIDQuery,
	}
	if request.RequestedOrgID != "" {
		orgIDQuery, err := query.NewIDPUserLinksResourceOwnerSearchQuery(request.RequestedOrgID)
		if err != nil {
			return err
		}
		queries = append(queries, orgIDQuery)
	}
	links, err := repo.Query.IDPUserLinks(ctx, &query.IDPUserLinksSearchQuery{Queries: queries}, nil)
	if err != nil {
		return err
	}
	if len(links.Links) != 1 {
		return zerrors.ThrowNotFound(nil, "AUTH-Sf8sd", "Errors.ExternalIDP.NotFound")
	}
	user, err := activeUserByID(ctx, repo.UserViewProvider, repo.UserEventProvider, repo.OrgViewProvider, repo.LockoutPolicyViewProvider, links.Links[0].UserID, false)
	if err != nil {
		return err
	}
	username := user.UserName
	if request.RequestedOrgID == "" {
		username = user.PreferredLoginName
	}
	request.SetUserInfo(user.ID, username, user.PreferredLoginName, user.DisplayName, user.AvatarKey, user.ResourceOwner)
	return nil
}

//nolint:gocognit
func (repo *AuthRequestRepo) nextSteps(ctx context.Context, request *domain.AuthRequest, checkLoggedIn bool) (steps []domain.NextStep, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if request == nil {
		return nil, zerrors.ThrowInvalidArgument(nil, "EVENT-ds27a", "Errors.Internal")
	}
	steps = make([]domain.NextStep, 0)
	if !checkLoggedIn && domain.IsPrompt(request.Prompt, domain.PromptNone) {
		return append(steps, &domain.RedirectToCallbackStep{}), nil
	}
	if request.UserID == "" {
		steps, err = repo.nextStepsUser(ctx, request)
		if err != nil || len(steps) > 0 {
			return steps, err
		}
	}
	user, err := activeUserByID(ctx, repo.UserViewProvider, repo.UserEventProvider, repo.OrgViewProvider, repo.LockoutPolicyViewProvider, request.UserID, request.LoginPolicy.IgnoreUnknownUsernames)
	if err != nil {
		return nil, err
	}
	// in case the user was set automatically, we might not have the org set
	if request.UserOrgID == "" {
		request.UserOrgID = user.ResourceOwner
	}
	userSession, err := userSessionByIDs(ctx, repo.UserSessionViewProvider, repo.UserEventProvider, request.AgentID, user.ID)
	if err != nil {
		return nil, err
	}
	request.SessionID = userSession.ID
	request.DisplayName = userSession.DisplayName
	request.AvatarKey = userSession.AvatarKey
	if user.HumanView != nil && user.HumanView.PreferredLanguage != "" {
		request.PreferredLanguage = gu.Ptr(language.Make(user.HumanView.PreferredLanguage))
	}

	isInternalLogin := (request.SelectedIDPConfigID == "" && userSession.SelectedIDPConfigID == "") || request.RequestLocalAuth
	idps, err := checkExternalIDPsOfUser(ctx, repo.IDPUserLinksProvider, user.ID)
	if err != nil {
		return nil, err
	}
	noLocalAuth := request.LoginPolicy != nil && !request.LoginPolicy.AllowUsernamePassword

	allowedLinkedIDPs := checkForAllowedIDPs(request.AllowedExternalIDPs, idps.Links)
	if (!isInternalLogin || len(allowedLinkedIDPs) > 0 || noLocalAuth) &&
		len(request.LinkingUsers) == 0 &&
		!request.RequestLocalAuth {
		step, err := repo.idpChecked(request, allowedLinkedIDPs, userSession)
		if err != nil {
			return nil, err
		}
		if step != nil {
			return append(steps, step), nil
		}
	}
	if isInternalLogin || (!isInternalLogin && len(request.LinkingUsers) > 0) {
		step := repo.firstFactorChecked(ctx, request, user, userSession)
		if step != nil {
			return append(steps, step), nil
		}
	}

	// If the user never had a verified email, we need to verify it.
	// This prevents situations, where OTP email is the only MFA method and no verified email is set.
	// If the user had a verified email, but change it and has not yet verified the new one, we'll verify it after we checked the MFA methods.
	if user.VerifiedEmail == "" && !user.IsEmailVerified {
		return append(steps, &domain.VerifyEMailStep{
			InitPassword: !user.PasswordSet && len(idps.Links) == 0,
		}), nil
	}

	step, ok, err := repo.mfaChecked(userSession, request, user, isInternalLogin && len(request.LinkingUsers) == 0)
	if err != nil {
		return nil, err
	}
	if !ok {
		return append(steps, step), nil
	}

	expired := passwordAgeChangeRequired(request.PasswordAgePolicy, user.PasswordChanged)
	if expired || user.PasswordChangeRequired {
		steps = append(steps, &domain.ChangePasswordStep{Expired: expired})
	}
	if !user.IsEmailVerified {
		steps = append(steps, &domain.VerifyEMailStep{
			InitPassword: !user.PasswordSet && len(idps.Links) == 0,
		})
	}
	if user.UsernameChangeRequired {
		steps = append(steps, &domain.ChangeUsernameStep{})
	}

	if expired || user.PasswordChangeRequired || !user.IsEmailVerified || user.UsernameChangeRequired {
		return steps, nil
	}

	if request.LinkingUsers != nil && len(request.LinkingUsers) != 0 {
		return append(steps, &domain.LinkUsersStep{}), nil
	}
	//PLANNED: consent step

	missing, err := projectRequired(ctx, request, repo.ProjectProvider)
	if err != nil {
		return nil, err
	}
	if missing {
		return append(steps, &domain.ProjectRequiredStep{}), nil
	}

	missing, err = userGrantRequired(ctx, request, user, repo.UserGrantProvider)
	if err != nil {
		return nil, err
	}
	if missing {
		return append(steps, &domain.GrantRequiredStep{}), nil
	}

	ok, err = repo.hasSucceededPage(ctx, request, repo.ApplicationProvider)
	if err != nil {
		return nil, err
	}
	if ok {
		steps = append(steps, &domain.LoginSucceededStep{})
	}
	return append(steps, &domain.RedirectToCallbackStep{}), nil
}

func checkForAllowedIDPs(allowedIDPs []*domain.IDPProvider, idps []*query.IDPUserLink) (_ []string) {
	allowedLinkedIDPs := make([]string, 0, len(idps))
	// only use allowed linked idps
	for _, idp := range idps {
		for _, allowedIdP := range allowedIDPs {
			if idp.IDPID == allowedIdP.IDPConfigID {
				allowedLinkedIDPs = append(allowedLinkedIDPs, allowedIdP.IDPConfigID)
			}
		}
	}
	return allowedLinkedIDPs
}

func passwordAgeChangeRequired(policy *domain.PasswordAgePolicy, changed time.Time) bool {
	if policy == nil || policy.MaxAgeDays == 0 {
		return false
	}
	maxDays := time.Duration(policy.MaxAgeDays) * 24 * time.Hour
	return time.Now().Add(-maxDays).After(changed)
}

func (repo *AuthRequestRepo) nextStepsUser(ctx context.Context, request *domain.AuthRequest) (_ []domain.NextStep, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	steps := make([]domain.NextStep, 0)
	if request.LinkingUsers != nil && len(request.LinkingUsers) > 0 {
		steps = append(steps, new(domain.ExternalNotFoundOptionStep))
		return steps, nil
	}
	if domain.IsPrompt(request.Prompt, domain.PromptCreate) {
		return append(steps, &domain.RegistrationStep{}), nil
	}
	// if there's a login or consent prompt, but not select account, just return the login step
	if len(request.Prompt) > 0 && !domain.IsPrompt(request.Prompt, domain.PromptSelectAccount) {
		return append(steps, new(domain.LoginStep)), nil
	} else {
		// if no user was specified, either select_account or no prompt was provided,
		// then check the active user sessions (of the user agent)
		users, err := repo.usersForUserSelection(ctx, request)
		if err != nil {
			return nil, err
		}
		// in case select_account was specified ignore it if there aren't any user sessions
		if domain.IsPrompt(request.Prompt, domain.PromptSelectAccount) && len(users) > 0 {
			steps = append(steps, &domain.SelectUserStep{Users: users})
		}
		// If we get here, either no sessions were found for select_account
		// or no prompt was provided.
		// In either case if there was a specific idp is selected (scope), directly redirect
		if request.SelectedIDPConfigID != "" {
			steps = append(steps, &domain.RedirectToExternalIDPStep{})
		}
		// or there aren't any sessions to use, present the login page (https://github.com/zitadel/zitadel/issues/7213)
		if len(users) == 0 {
			steps = append(steps, new(domain.LoginStep))
		}
		// if no prompt was provided, but there are multiple user sessions, then the user must decide which to use
		if len(request.Prompt) == 0 && len(users) > 1 {
			steps = append(steps, &domain.SelectUserStep{Users: users})
		}
		if len(steps) > 0 {
			return steps, nil
		}
		// the single user session was inactive
		if users[0].UserSessionState != domain.UserSessionStateActive {
			return append(steps, &domain.SelectUserStep{Users: users}), nil
		}
		// a single active user session was found, use that automatically
		request.SetUserInfo(users[0].UserID, users[0].UserName, users[0].LoginName, users[0].DisplayName, users[0].AvatarKey, users[0].ResourceOwner)
		if err = repo.fillPolicies(ctx, request); err != nil {
			return nil, err
		}
		if err = repo.AuthRequests.UpdateAuthRequest(ctx, request); err != nil {
			return nil, err
		}
	}
	return steps, nil
}

func checkExternalIDPsOfUser(ctx context.Context, idpUserLinksProvider idpUserLinksProvider, userID string) (*query.IDPUserLinks, error) {
	userIDQuery, err := query.NewIDPUserLinksUserIDSearchQuery(userID)
	if err != nil {
		return nil, err
	}
	return idpUserLinksProvider.IDPUserLinks(ctx, &query.IDPUserLinksSearchQuery{Queries: []query.SearchQuery{userIDQuery}}, nil)
}

func (repo *AuthRequestRepo) usersForUserSelection(ctx context.Context, request *domain.AuthRequest) ([]domain.UserSelection, error) {
	userSessions, err := userSessionsByUserAgentID(ctx, repo.UserSessionViewProvider, request.AgentID, request.InstanceID)
	if err != nil {
		return nil, err
	}
	users := make([]domain.UserSelection, 0)
	for _, session := range userSessions {
		if request.RequestedOrgID == "" || request.RequestedOrgID == session.ResourceOwner {
			users = append(users, domain.UserSelection{
				UserID:            session.UserID,
				DisplayName:       session.DisplayName,
				UserName:          session.UserName,
				LoginName:         session.LoginName,
				ResourceOwner:     session.ResourceOwner,
				AvatarKey:         session.AvatarKey,
				UserSessionState:  session.State,
				SelectionPossible: request.RequestedOrgID == "" || request.RequestedOrgID == session.ResourceOwner,
			})
		}
	}
	return users, nil
}

func (repo *AuthRequestRepo) firstFactorChecked(ctx context.Context, request *domain.AuthRequest, user *user_model.UserView, userSession *user_model.UserSessionView) domain.NextStep {
	if user.InitRequired {
		return &domain.InitUserStep{PasswordSet: user.PasswordSet}
	}

	var step domain.NextStep
	if request.LoginPolicy.PasswordlessType != domain.PasswordlessTypeNotAllowed && user.IsPasswordlessReady() {
		if checkVerificationTimeMaxAge(userSession.PasswordlessVerification, request.LoginPolicy.MultiFactorCheckLifetime, request) {
			request.MFAsVerified = append(request.MFAsVerified, domain.MFATypeU2FUserVerification)
			request.AuthTime = userSession.PasswordlessVerification
			return nil
		}
		step = &domain.PasswordlessStep{
			PasswordSet: user.PasswordSet,
		}
	}

	if user.PasswordlessInitRequired {
		return &domain.PasswordlessRegistrationPromptStep{}
	}

	if user.PasswordInitRequired {
		if !user.IsEmailVerified {
			// If the user was created through the user resource API,
			// they can either have an invite code...
			exists, err := repo.UserEventProvider.InviteCodeExists(ctx, user.ID)
			logging.WithFields("userID", user.ID).OnError(err).Error("unable to check if invite code exists")
			if err == nil && exists {
				return &domain.VerifyInviteStep{}
			}
			// or were created with an explicit email verification mail
			return &domain.VerifyEMailStep{InitPassword: true}
		}
		// If they were created with a verified mail, they might have never received mail to set their password,
		// e.g. when created through a user resource API. In this case we'll just create and send one now.
		exists, err := repo.UserEventProvider.PasswordCodeExists(ctx, user.ID)
		logging.WithFields("userID", user.ID).OnError(err).Error("unable to check if password code exists")
		if err == nil && !exists {
			_, err = repo.PasswordReset.RequestSetPassword(ctx, user.ID, user.ResourceOwner, domain.NotificationTypeEmail, request.ID)
			logging.WithFields("userID", user.ID).OnError(err).Error("unable to create password code")
		}
		return &domain.InitPasswordStep{}
	}

	if checkVerificationTimeMaxAge(userSession.PasswordVerification, request.LoginPolicy.PasswordCheckLifetime, request) {
		request.PasswordVerified = true
		request.AuthTime = userSession.PasswordVerification
		return nil
	}
	if step != nil {
		return step
	}
	return &domain.PasswordStep{}
}

func (repo *AuthRequestRepo) idpChecked(request *domain.AuthRequest, idps []string, userSession *user_model.UserSessionView) (domain.NextStep, error) {
	if checkVerificationTimeMaxAge(userSession.ExternalLoginVerification, request.LoginPolicy.ExternalLoginCheckLifetime, request) {
		request.IDPLoginChecked = true
		request.AuthTime = userSession.ExternalLoginVerification
		return nil, nil
	}
	// use the explicitly set IdP first
	if request.SelectedIDPConfigID != "" {
		// only use the explicitly set IdP if allowed
		for _, allowedIdP := range request.AllowedExternalIDPs {
			if request.SelectedIDPConfigID == allowedIdP.IDPConfigID {
				return &domain.ExternalLoginStep{SelectedIDPConfigID: request.SelectedIDPConfigID}, nil
			}
		}
		// error if the explicitly set IdP is not allowed, to avoid misinterpretation with usage of another IdP
		return nil, zerrors.ThrowPreconditionFailed(nil, "LOGIN-LWif2", "Errors.Org.IdpNotExisting")
	}
	// reuse the previously used IdP from the session
	if userSession.SelectedIDPConfigID != "" {
		// only use the previously used IdP if allowed
		for _, allowedIdP := range request.AllowedExternalIDPs {
			if userSession.SelectedIDPConfigID == allowedIdP.IDPConfigID {
				return &domain.ExternalLoginStep{SelectedIDPConfigID: userSession.SelectedIDPConfigID}, nil
			}
		}
	}
	// then use an existing linked and allowed IdP of the user
	if len(idps) > 0 {
		return &domain.ExternalLoginStep{SelectedIDPConfigID: idps[0]}, nil
	}
	// if the user did not link one, then just use one of the configured IdPs of the org
	if len(request.AllowedExternalIDPs) > 0 {
		return &domain.ExternalLoginStep{SelectedIDPConfigID: request.AllowedExternalIDPs[0].IDPConfigID}, nil
	}
	return nil, zerrors.ThrowPreconditionFailed(nil, "LOGIN-5Hm8s", "Errors.Org.IdpNotExisting")
}

func (repo *AuthRequestRepo) mfaChecked(userSession *user_model.UserSessionView, request *domain.AuthRequest, user *user_model.UserView, isInternalAuthentication bool) (domain.NextStep, bool, error) {
	mfaLevel := request.MFALevel()
	if slices.Contains(request.MFAsVerified, domain.MFATypeU2FUserVerification) {
		return nil, true, nil
	}
	allowedProviders, required := user.MFATypesAllowed(mfaLevel, request.LoginPolicy, isInternalAuthentication)
	promptRequired := (user.MFAMaxSetUp < mfaLevel) || (len(allowedProviders) == 0 && required)
	if promptRequired || !repo.mfaSkippedOrSetUp(user, request) {
		types := user.MFATypesSetupPossible(mfaLevel, request.LoginPolicy)
		if promptRequired && len(types) == 0 {
			return nil, false, zerrors.ThrowPreconditionFailed(nil, "LOGIN-5Hm8s", "Errors.Login.LoginPolicy.MFA.ForceAndNotConfigured")
		}
		if len(types) == 0 {
			return nil, true, nil
		}
		return &domain.MFAPromptStep{
			Required:     promptRequired,
			MFAProviders: types,
		}, false, nil
	}
	switch mfaLevel {
	default:
		fallthrough
	case domain.MFALevelNotSetUp:
		if len(allowedProviders) == 0 {
			return nil, true, nil
		}
		fallthrough
	case domain.MFALevelSecondFactor:
		if checkVerificationTimeMaxAge(userSession.SecondFactorVerification, request.LoginPolicy.SecondFactorCheckLifetime, request) {
			request.MFAsVerified = append(request.MFAsVerified, userSession.SecondFactorVerificationType)
			request.AuthTime = userSession.SecondFactorVerification
			return nil, true, nil
		}
		fallthrough
	case domain.MFALevelMultiFactor:
		if checkVerificationTimeMaxAge(userSession.MultiFactorVerification, request.LoginPolicy.MultiFactorCheckLifetime, request) {
			request.MFAsVerified = append(request.MFAsVerified, userSession.MultiFactorVerificationType)
			request.AuthTime = userSession.MultiFactorVerification
			return nil, true, nil
		}
	}
	return &domain.MFAVerificationStep{
		MFAProviders: allowedProviders,
	}, false, nil
}

func (repo *AuthRequestRepo) mfaSkippedOrSetUp(user *user_model.UserView, request *domain.AuthRequest) bool {
	if user.MFAMaxSetUp > domain.MFALevelNotSetUp {
		return true
	}
	if request.LoginPolicy.MFAInitSkipLifetime == 0 {
		return true
	}
	return checkVerificationTime(user.MFAInitSkipped, request.LoginPolicy.MFAInitSkipLifetime)
}

func (repo *AuthRequestRepo) GetPrivacyPolicy(ctx context.Context, orgID string) (*domain.PrivacyPolicy, error) {
	policy, err := repo.PrivacyPolicyProvider.PrivacyPolicyByOrg(ctx, false, orgID, false)
	if zerrors.IsNotFound(err) {
		return new(domain.PrivacyPolicy), nil
	}
	if err != nil {
		return nil, err
	}
	return privacyPolicyToDomain(policy), err
}

func privacyPolicyToDomain(p *query.PrivacyPolicy) *domain.PrivacyPolicy {
	return &domain.PrivacyPolicy{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:   p.ID,
			Sequence:      p.Sequence,
			ResourceOwner: p.ResourceOwner,
			CreationDate:  p.CreationDate,
			ChangeDate:    p.ChangeDate,
		},
		State:          p.State,
		Default:        p.IsDefault,
		TOSLink:        p.TOSLink,
		PrivacyLink:    p.PrivacyLink,
		HelpLink:       p.HelpLink,
		SupportEmail:   p.SupportEmail,
		DocsLink:       p.DocsLink,
		CustomLink:     p.CustomLink,
		CustomLinkText: p.CustomLinkText,
	}
}

func (repo *AuthRequestRepo) getLockoutPolicy(ctx context.Context, orgID string) (*query.LockoutPolicy, error) {
	policy, err := repo.LockoutPolicyViewProvider.LockoutPolicyByOrg(ctx, false, orgID)
	if err != nil {
		return nil, err
	}
	return policy, err
}

func (repo *AuthRequestRepo) getLabelPolicy(ctx context.Context, orgID string) (*domain.LabelPolicy, error) {
	policy, err := repo.LabelPolicyProvider.ActiveLabelPolicyByOrg(ctx, orgID, false)
	if err != nil {
		return nil, err
	}
	return labelPolicyToDomain(policy), nil
}

func labelPolicyToDomain(p *query.LabelPolicy) *domain.LabelPolicy {
	return &domain.LabelPolicy{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:   p.ID,
			Sequence:      p.Sequence,
			ResourceOwner: p.ResourceOwner,
			CreationDate:  p.CreationDate,
			ChangeDate:    p.ChangeDate,
		},
		State:               p.State,
		Default:             p.IsDefault,
		PrimaryColor:        p.Light.PrimaryColor,
		BackgroundColor:     p.Light.BackgroundColor,
		WarnColor:           p.Light.WarnColor,
		FontColor:           p.Light.FontColor,
		LogoURL:             p.Light.LogoURL,
		IconURL:             p.Light.IconURL,
		PrimaryColorDark:    p.Dark.PrimaryColor,
		BackgroundColorDark: p.Dark.BackgroundColor,
		WarnColorDark:       p.Dark.WarnColor,
		FontColorDark:       p.Dark.FontColor,
		LogoDarkURL:         p.Dark.LogoURL,
		IconDarkURL:         p.Dark.IconURL,
		Font:                p.FontURL,
		HideLoginNameSuffix: p.HideLoginNameSuffix,
		ErrorMsgPopup:       p.ShouldErrorPopup,
		DisableWatermark:    p.WatermarkDisabled,
		ThemeMode:           p.ThemeMode,
	}
}

func (repo *AuthRequestRepo) getPasswordAgePolicy(ctx context.Context, orgID string) (*domain.PasswordAgePolicy, error) {
	policy, err := repo.PasswordAgePolicyProvider.PasswordAgePolicyByOrg(ctx, false, orgID, false)
	if err != nil {
		return nil, err
	}
	return passwordAgePolicyToDomain(policy), nil
}

func passwordAgePolicyToDomain(policy *query.PasswordAgePolicy) *domain.PasswordAgePolicy {
	return &domain.PasswordAgePolicy{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:   policy.ID,
			Sequence:      policy.Sequence,
			ResourceOwner: policy.ResourceOwner,
			CreationDate:  policy.CreationDate,
			ChangeDate:    policy.ChangeDate,
		},
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
	}
}

func (repo *AuthRequestRepo) getLoginTexts(ctx context.Context, aggregateID string) ([]*domain.CustomText, error) {
	loginTexts, err := repo.CustomTextProvider.CustomTextListByTemplate(ctx, aggregateID, domain.LoginCustomText, false)
	if err != nil {
		return nil, err
	}
	return query.CustomTextsToDomain(loginTexts), err
}

func (repo *AuthRequestRepo) hasSucceededPage(ctx context.Context, request *domain.AuthRequest, provider applicationProvider) (bool, error) {
	if _, ok := request.Request.(*domain.AuthRequestOIDC); !ok {
		return false, nil
	}
	app, err := provider.AppByOIDCClientID(ctx, request.ApplicationID)
	if err != nil {
		return false, err
	}
	return app.OIDCConfig.AppType == domain.OIDCApplicationTypeNative && !app.OIDCConfig.SkipNativeAppSuccessPage, nil
}

func (repo *AuthRequestRepo) getDomainPolicy(ctx context.Context, orgID string) (*query.DomainPolicy, error) {
	return repo.Query.DomainPolicyByOrg(ctx, false, orgID, false)
}

func setOrgID(ctx context.Context, orgViewProvider orgViewProvider, request *domain.AuthRequest) error {
	orgID := request.GetScopeOrgID()
	if orgID != "" {
		org, err := orgViewProvider.OrgByID(ctx, orgID)
		if err != nil {
			return err
		}
		request.SetOrgInformation(org.ID, org.Name, org.Domain, false)
		return nil
	}

	primaryDomain := request.GetScopeOrgPrimaryDomain()
	if primaryDomain == "" {
		return nil
	}

	org, err := orgViewProvider.OrgByPrimaryDomain(ctx, primaryDomain)
	if err != nil {
		return err
	}
	request.SetOrgInformation(org.ID, org.Name, primaryDomain, true)
	return nil
}

func getLoginPolicyIDPProviders(ctx context.Context, provider idpProviderViewProvider, iamID, orgID string, defaultPolicy bool) ([]*domain.IDPProvider, error) {
	resourceOwner := iamID
	if !defaultPolicy {
		resourceOwner = orgID
	}
	links, err := provider.IDPLoginPolicyLinks(ctx, resourceOwner, &query.IDPLoginPolicyLinksSearchQuery{}, false)
	if err != nil {
		return nil, err
	}
	providers := make([]*domain.IDPProvider, len(links.Links))
	for i, link := range links.Links {
		providers[i] = &domain.IDPProvider{
			Type:        link.OwnerType,
			IDPConfigID: link.IDPID,
			Name:        link.IDPName,
			IDPType:     link.IDPType,
		}
	}
	return providers, nil
}

func checkVerificationTimeMaxAge(verificationTime time.Time, lifetime time.Duration, request *domain.AuthRequest) bool {
	if !checkVerificationTime(verificationTime, lifetime) {
		return false
	}
	if request.MaxAuthAge == nil {
		return true
	}
	return verificationTime.After(request.CreationDate.Add(-*request.MaxAuthAge))
}

func checkVerificationTime(verificationTime time.Time, lifetime time.Duration) bool {
	return verificationTime.Add(lifetime).After(time.Now().UTC())
}

func userSessionsByUserAgentID(ctx context.Context, provider userSessionViewProvider, agentID, instanceID string) (_ []*user_model.UserSessionView, err error) {
	//nolint
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	session, err := provider.UserSessionsByAgentID(ctx, agentID, instanceID)
	if err != nil {
		return nil, err
	}
	return user_view_model.UserSessionsToModel(session), nil
}

var (
	userSessionEventTypes = []eventstore.EventType{
		user_repo.UserV1PasswordCheckSucceededType,
		user_repo.UserV1PasswordCheckFailedType,
		user_repo.UserV1MFAOTPCheckSucceededType,
		user_repo.UserV1MFAOTPCheckFailedType,
		user_repo.UserV1SignedOutType,
		user_repo.UserLockedType,
		user_repo.UserDeactivatedType,
		user_repo.HumanPasswordCheckSucceededType,
		user_repo.HumanPasswordCheckFailedType,
		user_repo.UserIDPLoginCheckSucceededType,
		user_repo.HumanMFAOTPCheckSucceededType,
		user_repo.HumanMFAOTPCheckFailedType,
		user_repo.HumanSignedOutType,
		user_repo.HumanPasswordlessTokenCheckSucceededType,
		user_repo.HumanPasswordlessTokenCheckFailedType,
		user_repo.HumanU2FTokenCheckSucceededType,
		user_repo.HumanU2FTokenCheckFailedType,
		user_repo.UserRemovedType,
	}
)

func userSessionByIDs(ctx context.Context, provider userSessionViewProvider, eventProvider userEventProvider, agentID, userID string) (*user_model.UserSessionView, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()

	// always load the latest sequence first, so in case the session was not found by id,
	// the sequence will be equal or lower than the actual projection and no events are lost
	sequence, err := provider.GetLatestUserSessionSequence(ctx, instanceID)
	logging.WithFields("instanceID", instanceID, "userID", userID).
		OnError(err).
		Errorf("could not get current sequence for userSessionByIDs")

	session, err := provider.UserSessionByIDs(ctx, agentID, userID, instanceID)
	if err != nil {
		if !zerrors.IsNotFound(err) {
			return nil, err
		}
		session = &user_view_model.UserSessionView{UserAgentID: agentID, UserID: userID}
		if sequence != nil {
			session.ChangeDate = sequence.EventCreatedAt
		}
	}
	events, err := eventProvider.UserEventsByID(ctx, userID, session.ChangeDate, append(session.EventTypes(), userSessionEventTypes...))
	if err != nil {
		logging.WithFields("traceID", tracing.TraceIDFromCtx(ctx)).WithError(err).Debug("error retrieving new events")
		return user_view_model.UserSessionToModel(session), nil
	}
	sessionCopy := *session
	for _, event := range events {
		switch event.Type() {
		case user_repo.UserV1PasswordCheckSucceededType,
			user_repo.UserV1PasswordCheckFailedType,
			user_repo.UserV1MFAOTPCheckSucceededType,
			user_repo.UserV1MFAOTPCheckFailedType,
			user_repo.UserV1SignedOutType,
			user_repo.UserLockedType,
			user_repo.UserDeactivatedType,
			user_repo.HumanPasswordCheckSucceededType,
			user_repo.HumanPasswordCheckFailedType,
			user_repo.UserIDPLoginCheckSucceededType,
			user_repo.HumanMFAOTPCheckSucceededType,
			user_repo.HumanMFAOTPCheckFailedType,
			user_repo.HumanSignedOutType,
			user_repo.HumanPasswordlessTokenCheckSucceededType,
			user_repo.HumanPasswordlessTokenCheckFailedType,
			user_repo.HumanU2FTokenCheckSucceededType,
			user_repo.HumanU2FTokenCheckFailedType:
			userAgentID, err := user_view_model.UserAgentIDFromEvent(event)
			if err != nil {
				logging.WithFields("traceID", tracing.TraceIDFromCtx(ctx)).WithError(err).Debug("error getting event data")
				return user_view_model.UserSessionToModel(session), nil
			}
			if userAgentID != agentID {
				continue
			}
		}
		err := sessionCopy.AppendEvent(event)
		logging.WithFields("traceID", tracing.TraceIDFromCtx(ctx)).OnError(err).Warn("error appending event")
	}
	return user_view_model.UserSessionToModel(&sessionCopy), nil
}

func activeUserByID(ctx context.Context, userViewProvider userViewProvider, userEventProvider userEventProvider, queries orgViewProvider, lockoutPolicyProvider lockoutPolicyViewProvider, userID string, ignoreUnknownUsernames bool) (user *user_model.UserView, err error) {
	// PLANNED: Check LockoutPolicy
	user, err = userByID(ctx, userViewProvider, userEventProvider, userID)
	if err != nil {
		if ignoreUnknownUsernames && zerrors.IsNotFound(err) {
			return &user_model.UserView{
				ID:        userID,
				HumanView: &user_model.HumanView{},
			}, nil
		}
		return nil, err
	}

	if user.HumanView == nil {
		return nil, zerrors.ThrowPreconditionFailed(nil, "EVENT-Lm69x", "Errors.User.NotHuman")
	}
	if user.State == user_model.UserStateLocked || user.State == user_model.UserStateSuspend {
		return nil, zerrors.ThrowPreconditionFailed(nil, "EVENT-FJ262", "Errors.User.Locked")
	}
	if !(user.State == user_model.UserStateActive || user.State == user_model.UserStateInitial) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "EVENT-FJ262", "Errors.User.NotActive")
	}
	org, err := queries.OrgByID(ctx, user.ResourceOwner)
	if err != nil {
		return nil, err
	}
	if org.State != domain.OrgStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "EVENT-Zws3s", "Errors.User.NotActive")
	}
	return user, nil
}

func userByID(ctx context.Context, viewProvider userViewProvider, eventProvider userEventProvider, userID string) (_ *user_model.UserView, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	user, viewErr := viewProvider.UserByID(ctx, userID, authz.GetInstance(ctx).InstanceID())
	if viewErr != nil && !zerrors.IsNotFound(viewErr) {
		return nil, viewErr
	} else if user == nil {
		user = new(user_view_model.UserView)
	}
	events, err := eventProvider.UserEventsByID(ctx, userID, user.ChangeDate, user.EventTypes())
	if err != nil {
		logging.WithFields("traceID", tracing.TraceIDFromCtx(ctx)).WithError(err).Debug("error retrieving new events")
		return user_view_model.UserToModel(user), nil
	}
	if len(events) == 0 {
		if viewErr != nil {
			// We already returned all errors apart from not found, but need to make sure that can be checked in case IgnoreUnknownUsernames option is active.
			return nil, ErrUserNotFound(viewErr)
		}
		return user_view_model.UserToModel(user), nil
	}
	userCopy := *user
	for _, event := range events {
		if err := userCopy.AppendEvent(event); err != nil {
			return user_view_model.UserToModel(user), nil
		}
	}
	if userCopy.State == int32(user_model.UserStateDeleted) {
		return nil, zerrors.ThrowNotFound(nil, "EVENT-3F9so", "Errors.User.NotFound")
	}
	return user_view_model.UserToModel(&userCopy), nil
}

func linkExternalIDPs(ctx context.Context, userCommandProvider userCommandProvider, request *domain.AuthRequest) error {
	externalIDPs := make([]*domain.UserIDPLink, len(request.LinkingUsers))
	for i, linkingUser := range request.LinkingUsers {
		externalIDP := &domain.UserIDPLink{
			ObjectRoot:     es_models.ObjectRoot{AggregateID: request.UserID},
			IDPConfigID:    linkingUser.IDPConfigID,
			ExternalUserID: linkingUser.ExternalUserID,
			DisplayName:    linkingUser.PreferredUsername,
		}
		externalIDPs[i] = externalIDP
	}
	data := authz.CtxData{
		UserID: "LOGIN",
		OrgID:  request.UserOrgID,
	}
	return userCommandProvider.BulkAddedUserIDPLinks(authz.SetCtxData(ctx, data), request.UserID, request.UserOrgID, externalIDPs)
}

func linkingIDPConfigExistingInAllowedIDPs(linkingUsers []*domain.ExternalUser, idpProviders []*domain.IDPProvider) bool {
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

func userGrantRequired(ctx context.Context, request *domain.AuthRequest, user *user_model.UserView, userGrantProvider userGrantProvider) (_ bool, err error) {
	var project *query.Project
	switch request.Request.Type() {
	case domain.AuthRequestTypeOIDC, domain.AuthRequestTypeSAML, domain.AuthRequestTypeDevice:
		project, err = userGrantProvider.ProjectByClientID(ctx, request.ApplicationID)
		if err != nil {
			return false, err
		}
	default:
		return false, zerrors.ThrowPreconditionFailed(nil, "EVENT-dfrw2", "Errors.AuthRequest.RequestTypeNotSupported")
	}
	if !project.ProjectRoleCheck {
		return false, nil
	}
	grants, err := userGrantProvider.UserGrantsByProjectAndUserID(ctx, project.ID, user.ID)
	if err != nil {
		return false, err
	}
	return len(grants) == 0, nil
}

func projectRequired(ctx context.Context, request *domain.AuthRequest, projectProvider projectProvider) (missingGrant bool, err error) {
	var project *query.Project
	switch request.Request.Type() {
	case domain.AuthRequestTypeOIDC, domain.AuthRequestTypeSAML, domain.AuthRequestTypeDevice:
		project, err = projectProvider.ProjectByClientID(ctx, request.ApplicationID)
		if err != nil {
			return false, err
		}
	default:
		return false, zerrors.ThrowPreconditionFailed(nil, "EVENT-ku4He", "Errors.AuthRequest.RequestTypeNotSupported")
	}
	// if the user and project are part of the same organisation we do not need to check if the project exists on that org
	if !project.HasProjectCheck || project.ResourceOwner == request.UserOrgID {
		return false, nil
	}

	// else just check if there is a project grant for that org
	projectID, err := query.NewProjectGrantProjectIDSearchQuery(project.ID)
	if err != nil {
		return false, err
	}
	grantedOrg, err := query.NewProjectGrantGrantedOrgIDSearchQuery(request.UserOrgID)
	if err != nil {
		return false, err
	}
	grants, err := projectProvider.SearchProjectGrants(ctx, &query.ProjectGrantSearchQueries{Queries: []query.SearchQuery{projectID, grantedOrg}}, nil)
	if err != nil {
		return false, err
	}
	return len(grants.ProjectGrants) != 1, nil
}
