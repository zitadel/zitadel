package login

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp/providers/ldap"
)

const (
	tmplLDAPLogin = "ldap_login"
)

type ldapFormData struct {
	Username string `schema:"username"`
	Password string `schema:"password"`
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

	idpTemplate, err := l.getIDPTemplateByID(r, "199408199776862496", "202307103371559120")
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	password, err := crypto.DecryptString(idpTemplate.Password, l.idpConfigAlg)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	provider := ldap.New(
		idpTemplate.Name,
		idpTemplate.Host,
		idpTemplate.BaseDN,
		idpTemplate.UserObjectClass,
		idpTemplate.UserUniqueAttribute,
		idpTemplate.Admin,
		password,
		"not used",
		ldap.Insecure(),
	)

	session, err := provider.BeginAuth(r.Context(), authReq.ID, data.Username, data.Password)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}

	_, err = session.FetchUser(r.Context())
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
}
