package login

import (
	"net/http"

	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
	"github.com/zitadel/zitadel/internal/query"
)

const (
	tmplExternalRegisterOverview = "externalregisteroverview"
)

type externalRegisterFormData struct {
	ExternalIDPConfigID    string `schema:"external-idp-config-id"`
	ExternalIDPExtUserID   string `schema:"external-idp-ext-user-id"`
	ExternalIDPDisplayName string `schema:"external-idp-display-name"`
	ExternalEmail          string `schema:"external-email"`
	ExternalEmailVerified  bool   `schema:"external-email-verified"`
	Email                  string `schema:"email"`
	Username               string `schema:"username"`
	Firstname              string `schema:"firstname"`
	Lastname               string `schema:"lastname"`
	Nickname               string `schema:"nickname"`
	ExternalPhone          string `schema:"external-phone"`
	ExternalPhoneVerified  bool   `schema:"external-phone-verified"`
	Phone                  string `schema:"phone"`
	Language               string `schema:"language"`
	TermsConfirm           bool   `schema:"terms-confirm"`
}

type externalRegisterData struct {
	baseData
	externalRegisterFormData
	ExternalIDPID              string
	ExternalIDPUserID          string
	ExternalIDPUserDisplayName string
	ShowUsername               bool
	ShowUsernameSuffix         bool
	OrgRegister                bool
	ExternalEmail              string
	ExternalEmailVerified      bool
	ExternalPhone              string
	ExternalPhoneVerified      bool
}

func (l *Login) handleExternalRegister(w http.ResponseWriter, r *http.Request) {
	data := new(externalIDPData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.handleIDP(w, r, authReq, data.IDPConfigID)
	//l.handleExternalRegisterByConfigID(w, r, authReq, data.IDPConfigID)
}

//
//func (l *Login) handleExternalRegisterByConfigID(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, configID string) {
//	if authReq == nil {
//		l.defaultRedirect(w, r)
//		return
//	}
//	idpConfig, err := l.getIDPConfigByID(r, configID)
//	if err != nil {
//		l.renderError(w, r, authReq, err)
//		return
//	}
//	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
//	err = l.authRepo.SelectExternalIDP(r.Context(), authReq.ID, idpConfig.IDPConfigID, userAgentID)
//	if err != nil {
//		l.renderLogin(w, r, authReq, err)
//		return
//	}
//	if !idpConfig.IsOIDC {
//		l.handleJWTAuthorize(w, r, authReq, idpConfig)
//		return
//	}
//	l.handleOIDCAuthorize(w, r, authReq, idpConfig, EndpointExternalRegisterCallback)
//}
//
//func (l *Login) handleExternalRegisterCallback(w http.ResponseWriter, r *http.Request) {
//	data := new(externalIDPCallbackData)
//	err := l.getParseData(r, data)
//	if err != nil {
//		l.renderError(w, r, nil, err)
//		return
//	}
//	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
//	authReq, err := l.authRepo.AuthRequestByID(r.Context(), data.State, userAgentID)
//	if err != nil {
//		l.renderError(w, r, authReq, err)
//		return
//	}
//	idpConfig, err := l.authRepo.GetIDPConfigByID(r.Context(), authReq.SelectedIDPConfigID)
//	if err != nil {
//		l.renderError(w, r, authReq, err)
//		return
//	}
//	provider, err := l.getRPConfig(r.Context(), idpConfig, EndpointExternalRegisterCallback)
//	if err != nil {
//		l.renderRegisterOption(w, r, authReq, err)
//		return
//	}
//	tokens, err := rp.CodeExchange(r.Context(), data.Code, provider)
//	if err != nil {
//		l.renderRegisterOption(w, r, authReq, err)
//		return
//	}
//	l.handleExternalUserRegister(w, r, authReq, idpConfig, userAgentID, tokens)
//}

func (l *Login) handleExternalUserRegister(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *iam_model.IDPConfigView, userAgentID string, tokens *oidc.Tokens) {
	resourceOwner := authz.GetInstance(r.Context()).DefaultOrganisationID()
	if authReq.RequestedOrgID != "" {
		resourceOwner = authReq.RequestedOrgID
	}
	externalUser, externalIDP := l.mapTokenToLoginHumanAndExternalIDP(tokens, idpConfig)
	externalUser, err := l.runPostExternalAuthenticationActions(externalUser, tokens, authReq, r, nil)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	if idpConfig.AutoRegister {
		l.registerExternalUser(w, r, authReq, externalUser)
		return
	}
	orgIamPolicy, err := l.getOrgDomainPolicy(r, resourceOwner)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	labelPolicy, err := l.getLabelPolicy(r, resourceOwner)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	l.renderExternalRegisterOverview(w, r, authReq, orgIamPolicy, externalUser, externalIDP, labelPolicy.HideLoginNameSuffix, nil)
}

//
//func (l *Login) registerExternalUser(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, externalUser *domain.ExternalUser) {
//	resourceOwner := authz.GetInstance(r.Context()).DefaultOrganisationID()
//
//	if authReq.RequestedOrgID != "" && authReq.RequestedOrgID != resourceOwner {
//		resourceOwner = authReq.RequestedOrgID
//	}
//	orgIamPolicy, err := l.getOrgDomainPolicy(r, resourceOwner)
//	if err != nil {
//		l.renderRegisterOption(w, r, authReq, err)
//		return
//	}
//
//	idpConfig, err := l.authRepo.GetIDPConfigByID(r.Context(), authReq.SelectedIDPConfigID)
//	if err != nil {
//		l.renderRegisterOption(w, r, authReq, err)
//		return
//	}
//	user, externalIDP, metadata := l.mapExternalUserToLoginUser(orgIamPolicy, externalUser, idpConfig)
//	user, metadata, err = l.runPreCreationActions(authReq, r, user, metadata, resourceOwner, domain.FlowTypeExternalAuthentication)
//	if err != nil {
//		l.renderRegisterOption(w, r, authReq, err)
//		return
//	}
//	err = l.authRepo.AutoRegisterExternalUser(setContext(r.Context(), resourceOwner), user, externalIDP, nil, authReq.ID, authReq.AgentID, resourceOwner, metadata, nil)
//	if err != nil {
//		l.renderRegisterOption(w, r, authReq, err)
//		return
//	}
//	// read auth request again to get current state including userID
//	authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
//	if err != nil {
//		l.renderError(w, r, authReq, err)
//		return
//	}
//	userGrants, err := l.runPostCreationActions(authReq.UserID, authReq, r, resourceOwner, domain.FlowTypeExternalAuthentication)
//	if err != nil {
//		l.renderError(w, r, authReq, err)
//		return
//	}
//	err = l.appendUserGrants(r.Context(), userGrants, resourceOwner)
//	if err != nil {
//		l.renderError(w, r, authReq, err)
//		return
//	}
//	l.renderNextStep(w, r, authReq)
//}

func (l *Login) renderExternalRegisterOverview(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, orgIAMPolicy *query.DomainPolicy, externalUser *domain.ExternalUser, idp *domain.UserIDPLink, hideLoginNameSuffix bool, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}

	translator := l.getTranslator(r.Context(), authReq)
	data := externalRegisterData{
		baseData: l.getBaseData(r, authReq, "ExternalRegistrationUserOverview.Title", "ExternalRegistrationUserOverview.Description", errID, errMessage),
		externalRegisterFormData: externalRegisterFormData{
			Email:     externalUser.Email,
			Username:  externalUser.PreferredUsername,
			Firstname: externalUser.FirstName,
			Lastname:  externalUser.LastName,
			Nickname:  externalUser.NickName,
			Language:  externalUser.PreferredLanguage.String(),
		},
		ExternalIDPID:              idp.IDPConfigID,
		ExternalIDPUserID:          idp.ExternalUserID,
		ExternalIDPUserDisplayName: idp.DisplayName,
		ExternalEmail:              externalUser.Email,
		ExternalEmailVerified:      externalUser.IsEmailVerified,
		ShowUsername:               orgIAMPolicy.UserLoginMustBeDomain,
		OrgRegister:                orgIAMPolicy.UserLoginMustBeDomain,
		ShowUsernameSuffix:         !hideLoginNameSuffix,
	}
	data.Phone = externalUser.Phone
	data.ExternalPhone = externalUser.Phone
	data.ExternalPhoneVerified = externalUser.IsPhoneVerified

	funcs := map[string]interface{}{
		"selectedLanguage": func(l string) bool {
			return data.Language == l
		},
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplExternalRegisterOverview], data, funcs)
}

func (l *Login) handleExternalRegisterCheck(w http.ResponseWriter, r *http.Request) {
	data := new(externalRegisterFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	resourceOwner := authz.GetInstance(r.Context()).DefaultOrganisationID()

	if authReq.RequestedOrgID != "" && authReq.RequestedOrgID != resourceOwner {
		resourceOwner = authReq.RequestedOrgID
	}

	user := l.mapExternalRegisterDataToUser(data)
	l.registerExternalUser(w, r, authReq, user)
}

func (l *Login) mapTokenToLoginHumanAndExternalIDP(tokens *oidc.Tokens, idpConfig *iam_model.IDPConfigView) (*domain.ExternalUser, *domain.UserIDPLink) {
	displayName := tokens.IDTokenClaims.GetPreferredUsername()
	switch idpConfig.OIDCIDPDisplayNameMapping {
	case iam_model.OIDCMappingFieldEmail:
		if tokens.IDTokenClaims.IsEmailVerified() && tokens.IDTokenClaims.GetEmail() != "" {
			displayName = tokens.IDTokenClaims.GetEmail()
		}
	}
	if displayName == "" {
		displayName = tokens.IDTokenClaims.GetEmail()
	}

	externalUser := &domain.ExternalUser{
		IDPConfigID:       idpConfig.IDPConfigID,
		ExternalUserID:    tokens.IDTokenClaims.GetSubject(),
		PreferredUsername: tokens.IDTokenClaims.GetPreferredUsername(),
		DisplayName:       displayName,
		FirstName:         tokens.IDTokenClaims.GetGivenName(),
		LastName:          tokens.IDTokenClaims.GetFamilyName(),
		NickName:          tokens.IDTokenClaims.GetNickname(),
		Email:             tokens.IDTokenClaims.GetEmail(),
		IsEmailVerified:   tokens.IDTokenClaims.IsEmailVerified(),
		PreferredLanguage: tokens.IDTokenClaims.GetLocale(),
	}

	if tokens.IDTokenClaims.GetPhoneNumber() != "" {
		externalUser.Phone = tokens.IDTokenClaims.GetPhoneNumber()
		externalUser.IsPhoneVerified = tokens.IDTokenClaims.IsPhoneNumberVerified()
	}

	externalIDP := &domain.UserIDPLink{
		IDPConfigID:    idpConfig.IDPConfigID,
		ExternalUserID: tokens.IDTokenClaims.GetSubject(),
		DisplayName:    displayName,
	}
	return externalUser, externalIDP
}

func (l *Login) mapExternalRegisterDataToUser(data *externalRegisterFormData) *domain.ExternalUser {
	isEmailVerified := data.ExternalEmailVerified && data.Email == data.ExternalEmail
	isPhoneVerified := data.ExternalPhoneVerified && data.Phone == data.ExternalPhone
	return &domain.ExternalUser{
		IDPConfigID:       data.ExternalIDPConfigID,
		ExternalUserID:    data.ExternalIDPExtUserID,
		PreferredUsername: data.Username,
		DisplayName:       data.Email,
		FirstName:         data.Firstname,
		LastName:          data.Lastname,
		NickName:          data.Nickname,
		PreferredLanguage: language.Make(data.Language),
		Email:             data.Email,
		IsEmailVerified:   isEmailVerified,
		Phone:             data.Phone,
		IsPhoneVerified:   isPhoneVerified,
	}
}
