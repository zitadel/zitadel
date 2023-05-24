package login

import (
	"net/http"

	"github.com/zitadel/logging"

	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp/providers/ldap"
)

const (
	tmplLDAPLogin = "ldap_login"
)

type ldapFormData struct {
	Username         string `schema:"ldapusername"`
	Password         string `schema:"ldappassword"`
	ResetExternalIDP bool   `schema:"resetexternalidp"`
}

func (l *Login) handleLDAP(w http.ResponseWriter, r *http.Request) {
	authReq, err := l.getAuthRequest(r)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.renderLDAPLogin(w, r, authReq, nil)
}

func (l *Login) renderLDAPLogin(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	temp := l.renderer.Templates[tmplLDAPLogin]
	data := l.getUserData(r, authReq, "Login.Title", "Login.Description", errID, errMessage)
	l.renderer.RenderTemplate(w, r, l.getTranslator(r.Context(), authReq), temp, data, nil)
}

func (l *Login) handleLDAPCallback(w http.ResponseWriter, r *http.Request) {
	data := new(ldapFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if data.ResetExternalIDP {
		userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
		err := l.authRepo.ResetSelectedIDP(r.Context(), authReq.ID, userAgentID)
		if err != nil {
			l.renderLDAPLogin(w, r, authReq, err)
			return
		}

		l.handleLoginName(w, r)
		return
	}

	identityProvider, err := l.getIDPByID(r, authReq.SelectedIDPConfigID)
	if err != nil {
		l.renderLDAPLogin(w, r, authReq, err)
		return
	}

	provider, err := l.ldapProvider(r.Context(), identityProvider)
	if err != nil {
		l.renderLDAPLogin(w, r, authReq, err)
		return
	}
	session := &ldap.Session{Provider: provider, User: data.Username, Password: data.Password}

	user, err := session.FetchUser(r.Context())
	if err != nil {
		if _, _, actionErr := l.runPostExternalAuthenticationActions(new(domain.ExternalUser), nil, authReq, r, nil, err); actionErr != nil {
			logging.WithError(err).Error("both external user authentication and action post authentication failed")
		}
		l.renderLDAPLogin(w, r, authReq, err)
		return
	}
	l.handleExternalUserAuthenticated(w, r, authReq, identityProvider, session, user, l.renderNextStep)
}
