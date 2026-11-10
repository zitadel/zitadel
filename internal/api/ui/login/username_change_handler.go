package login

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
)

const (
	tmplChangeUsername     = "changeusername"
	tmplChangeUsernameDone = "changeusernamedone"
)

type changeUsernameData struct {
	Username string `schema:"username"`
}

func (l *Login) renderChangeUsername(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	translator := l.getTranslator(r.Context(), authReq)
	data := l.getUserData(r, authReq, translator, "UsernameChange.Title", "UsernameChange.Description", err)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplChangeUsername], data, nil)
}

func (l *Login) handleChangeUsername(w http.ResponseWriter, r *http.Request) {
	data := new(changeUsernameData)
	authReq, err := l.ensureAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	// The rename may only proceed once the flow has reached the change-username step;
	// isChangeUsernameStep documents why that also proves the caller is authenticated.
	if !isChangeUsernameStep(authReq) {
		l.renderError(w, r, authReq, nil)
		return
	}
	_, err = l.command.ChangeUsername(setContext(r.Context(), authReq.UserOrgID), authReq.UserOrgID, authReq.UserID, data.Username)
	if err != nil {
		l.renderChangeUsername(w, r, authReq, err)
		return
	}
	l.renderChangeUsernameDone(w, r, authReq)
}

func (l *Login) renderChangeUsernameDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	translator := l.getTranslator(r.Context(), authReq)
	data := l.getUserData(r, authReq, translator, "UsernameChangeDone.Title", "UsernameChangeDone.Description", nil)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplChangeUsernameDone], data, nil)
}

// isChangeUsernameStep reports whether the login flow's current step is the change-username
// step. ChangeUsernameStep is appended in nextSteps only after the first factor and any
// required MFA have succeeded, so gating on it guarantees the caller is the authenticated
// owner of authReq.UserID. Without this gate an unauthenticated caller could bind
// authReq.UserID to any account by submitting its login name and then rename it.
func isChangeUsernameStep(authReq *domain.AuthRequest) bool {
	if authReq == nil || len(authReq.PossibleSteps) == 0 {
		return false
	}
	_, ok := authReq.PossibleSteps[0].(*domain.ChangeUsernameStep)
	return ok
}
