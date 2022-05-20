package handler

import (
	"net/http"
	"strings"

	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/oidc"
	"golang.org/x/text/language"

	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
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
	if authReq == nil {
		l.defaultRedirect(w, r)
		return
	}
	idpConfig, err := l.getIDPConfigByID(r, data.IDPConfigID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	err = l.authRepo.SelectExternalIDP(r.Context(), authReq.ID, idpConfig.IDPConfigID, userAgentID)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	if !idpConfig.IsOIDC {
		l.handleJWTAuthorize(w, r, authReq, idpConfig)
		return
	}
	l.handleOIDCAuthorize(w, r, authReq, idpConfig, EndpointExternalRegisterCallback)
}

func (l *Login) handleExternalRegisterCallback(w http.ResponseWriter, r *http.Request) {
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
	idpConfig, err := l.authRepo.GetIDPConfigByID(r.Context(), authReq.SelectedIDPConfigID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	provider, err := l.getRPConfig(idpConfig, EndpointExternalRegisterCallback)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	tokens, err := rp.CodeExchange(r.Context(), data.Code, provider)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	l.handleExternalUserRegister(w, r, authReq, idpConfig, userAgentID, tokens)
}

func (l *Login) handleExternalUserRegister(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *iam_model.IDPConfigView, userAgentID string, tokens *oidc.Tokens) {
	iam, err := l.query.IAMByID(r.Context(), domain.IAMID)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	resourceOwner := iam.GlobalOrgID
	if authReq.RequestedOrgID != "" {
		resourceOwner = authReq.RequestedOrgID
	}
	orgIamPolicy, err := l.getOrgIamPolicy(r, resourceOwner)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	user, externalIDP := l.mapTokenToLoginHumanAndExternalIDP(orgIamPolicy, tokens, idpConfig)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	if !idpConfig.AutoRegister {
		l.renderExternalRegisterOverview(w, r, authReq, orgIamPolicy, user, externalIDP, nil)
		return
	}
	l.registerExternalUser(w, r, authReq, iam, user, externalIDP)
}

func (l *Login) registerExternalUser(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, iam *query.IAM, user *domain.Human, externalIDP *domain.UserIDPLink) {
	resourceOwner := iam.GlobalOrgID
	memberRoles := []string{domain.RoleSelfManagementGlobal}

	if authReq.RequestedOrgID != "" && authReq.RequestedOrgID != resourceOwner {
		memberRoles = nil
		resourceOwner = authReq.RequestedOrgID
	}
	_, err := l.command.RegisterHuman(setContext(r.Context(), resourceOwner), resourceOwner, user, externalIDP, memberRoles)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}

func (l *Login) renderExternalRegisterOverview(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, orgIAMPolicy *query.OrgIAMPolicy, human *domain.Human, idp *domain.UserIDPLink, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}

	data := externalRegisterData{
		baseData: l.getBaseData(r, authReq, "ExternalRegisterOverview", errID, errMessage),
		externalRegisterFormData: externalRegisterFormData{
			Email:     human.EmailAddress,
			Username:  human.Username,
			Firstname: human.FirstName,
			Lastname:  human.LastName,
			Nickname:  human.NickName,
			Language:  human.PreferredLanguage.String(),
		},
		ExternalIDPID:              idp.IDPConfigID,
		ExternalIDPUserID:          idp.ExternalUserID,
		ExternalIDPUserDisplayName: idp.DisplayName,
		ExternalEmail:              human.EmailAddress,
		ExternalEmailVerified:      human.IsEmailVerified,
		ShowUsername:               orgIAMPolicy.UserLoginMustBeDomain,
		OrgRegister:                orgIAMPolicy.UserLoginMustBeDomain,
	}
	if human.Phone != nil {
		data.Phone = human.PhoneNumber
		data.ExternalPhone = human.PhoneNumber
		data.ExternalPhoneVerified = human.IsPhoneVerified
	}
	translator := l.getTranslator(authReq)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplExternalRegisterOverview], data, nil)
}

func (l *Login) handleExternalRegisterCheck(w http.ResponseWriter, r *http.Request) {
	data := new(externalRegisterFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	iam, err := l.query.IAMByID(r.Context(), domain.IAMID)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	resourceOwner := iam.GlobalOrgID
	memberRoles := []string{domain.RoleSelfManagementGlobal}

	if authReq.RequestedOrgID != "" && authReq.RequestedOrgID != iam.GlobalOrgID {
		memberRoles = nil
		resourceOwner = authReq.RequestedOrgID
	}
	externalIDP, err := l.getExternalIDP(data)
	if externalIDP == nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	user, err := l.mapExternalRegisterDataToUser(r, data)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	_, err = l.command.RegisterHuman(setContext(r.Context(), resourceOwner), resourceOwner, user, externalIDP, memberRoles)
	if err != nil {
		l.renderRegisterOption(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}

func (l *Login) mapTokenToLoginHumanAndExternalIDP(orgIamPolicy *query.OrgIAMPolicy, tokens *oidc.Tokens, idpConfig *iam_model.IDPConfigView) (*domain.Human, *domain.UserIDPLink) {
	username := tokens.IDTokenClaims.GetPreferredUsername()
	switch idpConfig.OIDCUsernameMapping {
	case iam_model.OIDCMappingFieldEmail:
		if tokens.IDTokenClaims.IsEmailVerified() && tokens.IDTokenClaims.GetEmail() != "" {
			username = tokens.IDTokenClaims.GetEmail()
		}
	}
	if username == "" {
		username = tokens.IDTokenClaims.GetEmail()
	}

	if orgIamPolicy.UserLoginMustBeDomain {
		splittedUsername := strings.Split(username, "@")
		if len(splittedUsername) > 1 {
			username = splittedUsername[0]
		}
	}

	human := &domain.Human{
		Username: username,
		Profile: &domain.Profile{
			FirstName:         tokens.IDTokenClaims.GetGivenName(),
			LastName:          tokens.IDTokenClaims.GetFamilyName(),
			PreferredLanguage: tokens.IDTokenClaims.GetLocale(),
			NickName:          tokens.IDTokenClaims.GetNickname(),
		},
		Email: &domain.Email{
			EmailAddress:    tokens.IDTokenClaims.GetEmail(),
			IsEmailVerified: tokens.IDTokenClaims.IsEmailVerified(),
		},
	}

	if tokens.IDTokenClaims.GetPhoneNumber() != "" {
		human.Phone = &domain.Phone{
			PhoneNumber:     tokens.IDTokenClaims.GetPhoneNumber(),
			IsPhoneVerified: tokens.IDTokenClaims.IsPhoneNumberVerified(),
		}
	}

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

	externalIDP := &domain.UserIDPLink{
		IDPConfigID:    idpConfig.IDPConfigID,
		ExternalUserID: tokens.IDTokenClaims.GetSubject(),
		DisplayName:    displayName,
	}
	return human, externalIDP
}

func (l *Login) mapExternalRegisterDataToUser(r *http.Request, data *externalRegisterFormData) (*domain.Human, error) {
	human := &domain.Human{
		Username: data.Username,
		Profile: &domain.Profile{
			FirstName:         data.Firstname,
			LastName:          data.Lastname,
			PreferredLanguage: language.Make(data.Language),
			NickName:          data.Nickname,
		},
		Email: &domain.Email{
			EmailAddress: data.Email,
		},
	}
	if data.ExternalEmail != data.Email {
		human.IsEmailVerified = false
	} else {
		human.IsEmailVerified = data.ExternalEmailVerified
	}
	if data.ExternalPhone == "" {
		return human, nil
	}
	human.Phone = &domain.Phone{
		PhoneNumber: data.Phone,
	}
	if data.ExternalPhone != data.Phone {
		human.IsPhoneVerified = false
	} else {
		human.IsPhoneVerified = data.ExternalPhoneVerified
	}
	return human, nil
}

func (l *Login) getExternalIDP(data *externalRegisterFormData) (*domain.UserIDPLink, error) {
	return &domain.UserIDPLink{
		IDPConfigID:    data.ExternalIDPConfigID,
		ExternalUserID: data.ExternalIDPExtUserID,
		DisplayName:    data.ExternalIDPDisplayName,
	}, nil
}
