package login

import (
	"net/http"

	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/google"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	"github.com/zitadel/zitadel/internal/query"
)

func (l *Login) handleExternalLoginStep(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, selectedIDPID string) {
	for _, idp := range authReq.AllowedExternalIDPs {
		if idp.IDPConfigID == selectedIDPID {
			l.handleIDP(w, r, authReq, selectedIDPID)
			return
		}
	}
	l.renderLogin(w, r, authReq, errors.ThrowInvalidArgument(nil, "VIEW-Fsj7f", "Errors.User.ExternalIDP.NotAllowed"))
}

func (l *Login) handleExternalLogin(w http.ResponseWriter, r *http.Request) {
	data := new(externalIDPData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if authReq == nil {
		l.defaultRedirect(w, r)
		return
	}
	l.handleIDP(w, r, authReq, data.IDPConfigID)
}

func (l *Login) handleIDP(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, id string) {
	identityProvider, err := l.getIDPByID(r, id)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	err = l.authRepo.SelectExternalIDP(r.Context(), authReq.ID, identityProvider.ID, userAgentID)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	var provider idp.Provider
	switch identityProvider.Type {
	case domain.IDPTypeOIDC:
		//oidc.New(provider.Name, provider.)
	case domain.IDPTypeJWT:
	case domain.IDPTypeOAuth:
	case domain.IDPTypeLDAP:
	case domain.IDPTypeAzureAD:
	case domain.IDPTypeGitHub:
	case domain.IDPTypeGitHubEE:
	case domain.IDPTypeGitLab:
	case domain.IDPTypeGitLabSelfHosted:
	case domain.IDPTypeGoogle:
		secret, err := crypto.DecryptString(identityProvider.ClientSecret, l.idpConfigAlg)
		if err != nil {
			l.renderLogin(w, r, authReq, err)
			return
		}
		provider, err = google.New(identityProvider.ClientID, secret, EndpointExternalLoginCallback)
		if err != nil {

		}
	}
	session, err := provider.BeginAuth(r.Context(), authReq.ID)
	if err != nil {

	}
	http.Redirect(w, r, session.GetAuthURL(), http.StatusFound)
}

func (l *Login) handleExternalLoginCallback(w http.ResponseWriter, r *http.Request) {
	data := new(externalIDPCallbackData)
	err := l.getParseData(r, data)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	authReq, err := l.authRepo.AuthRequestByID(r.Context(), data.State, userAgentID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	identityProvider, err := l.getIDPByID(r, authReq.SelectedIDPConfigID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	var provider idp.Provider
	var session idp.Session
	switch identityProvider.Type {
	case domain.IDPTypeOIDC:
		//oidc.New(provider.Name, provider.)
	case domain.IDPTypeJWT:
	case domain.IDPTypeOAuth:
	case domain.IDPTypeLDAP:
	case domain.IDPTypeAzureAD:
	case domain.IDPTypeGitHub:
	case domain.IDPTypeGitHubEE:
	case domain.IDPTypeGitLab:
	case domain.IDPTypeGitLabSelfHosted:
	case domain.IDPTypeGoogle:
		secret, err := crypto.DecryptString(identityProvider.ClientSecret, l.idpConfigAlg)
		if err != nil {
			l.renderLogin(w, r, authReq, err)
			return
		}
		provider, err = google.New(identityProvider.ClientID, secret, EndpointExternalLoginCallback)
		if err != nil {

		}
		session = &openid.Session{Provider: provider, AuthURL: "", Code: data.Code}
	}

	user, err := session.FetchUser(r.Context())
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	l.handleExternalUserAuthenticated(w, r, authReq, provider, session, user)
	return
}

// tokens extracts the oidc.Tokens for backwards compatibility of PostExternalAuthenticationActions
func tokens(session idp.Session) *oidc.Tokens {
	if s, ok := session.(*openid.Session); ok {
		return s.Tokens
	}
	return nil
}

func (l *Login) handleExternalUserAuthenticated(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, provider *query.IDPTemplate, session idp.Session, user idp.User) {
	externalUser := mapIDPUserToExternalUser(user, provider.ID)
	externalUser, err := l.runPostExternalAuthenticationActions(externalUser, tokens(session), authReq, r, nil)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	err = l.authRepo.CheckExternalUserLogin(setContext(r.Context(), ""), authReq.ID, authReq.AgentID, externalUser, domain.BrowserInfoFromRequest(r))
	if err != nil {
		if errors.IsNotFound(err) { // TODO: handle error
			err = nil
		}
		l.externalUserNotExisting(w, r, authReq, provider, externalUser)
		return
	}
	if len(externalUser.Metadatas) > 0 {
		authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.ID)
		if err != nil {
			return
		}
		_, err = l.command.BulkSetUserMetadata(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, externalUser.Metadatas...)
		if err != nil {
			l.renderError(w, r, authReq, err)
			return
		}
	}
	l.renderNextStep(w, r, authReq)
}

func mapIDPUserToExternalUser(user idp.User, id string) *domain.ExternalUser {
	return &domain.ExternalUser{
		IDPConfigID:       id,
		ExternalUserID:    user.GetID(),
		PreferredUsername: user.GetPreferredUsername(),
		DisplayName:       user.GetDisplayName(),
		FirstName:         user.GetFirstName(),
		LastName:          user.GetLastName(),
		NickName:          user.GetNickname(),
		Email:             user.GetEmail(),
		IsEmailVerified:   user.IsEmailVerified(),
		PreferredLanguage: user.GetPreferredLanguage(),
		Phone:             user.GetPhone(),
		IsPhoneVerified:   user.IsPhoneVerified(),
	}
}

func mapExternalUserToLoginUser(externalUser *domain.ExternalUser, mustBeDomain bool) (*domain.Human, *domain.UserIDPLink, []*domain.Metadata) {
	human := &domain.Human{
		Username: externalUser.PreferredUsername,
		Profile: &domain.Profile{
			FirstName:         externalUser.FirstName,
			LastName:          externalUser.LastName,
			PreferredLanguage: externalUser.PreferredLanguage,
			NickName:          externalUser.NickName,
			DisplayName:       externalUser.DisplayName,
		},
		Email: &domain.Email{
			EmailAddress:    externalUser.Email,
			IsEmailVerified: externalUser.IsEmailVerified,
		},
	}
	if externalUser.Phone != "" {
		human.Phone = &domain.Phone{
			PhoneNumber:     externalUser.Phone,
			IsPhoneVerified: externalUser.IsPhoneVerified,
		}
	}
	externalIDP := &domain.UserIDPLink{
		IDPConfigID:    externalUser.IDPConfigID,
		ExternalUserID: externalUser.ExternalUserID,
		DisplayName:    externalUser.DisplayName,
	}
	return human, externalIDP, externalUser.Metadatas
}

func (l *Login) externalUserNotExisting(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, provider *query.IDPTemplate, externalUser *domain.ExternalUser) {
	resourceOwner := authz.GetInstance(r.Context()).DefaultOrganisationID()

	if authReq.RequestedOrgID != "" && authReq.RequestedOrgID != resourceOwner {
		resourceOwner = authReq.RequestedOrgID
	}

	orgIAMPolicy, err := l.getOrgDomainPolicy(r, resourceOwner)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, err)
		return
	}

	if !provider.IsAutoCreation {
		human, idpLinking, _ := mapExternalUserToLoginUser(externalUser, orgIAMPolicy.UserLoginMustBeDomain)
		l.renderExternalNotFoundOption(w, r, authReq, orgIAMPolicy, human, idpLinking, err)
		return
	}

	// reload auth request, to ensure current state (checked external login)
	authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
	if err != nil {
		l.renderExternalNotFoundOption(w, r, authReq, orgIAMPolicy, human, idpLinking, err)
		return
	}
	l.handleAutoRegister(w, r, authReq, false)
	return
}
