package login

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
)

const (
	tmplPasswordlessPrompt = "passwordlessprompt"
)

type passwordlessPromptData struct {
	userData
}

// renderPasswordlessPrompt shows the informational passwordless prompt page for the
// PasswordlessRegistrationPromptStep, instructing the user to complete setup via the link
// they were emailed. The interactive auth-request based setup that this page used to POST to
// was a never-active leftover and has been removed (GHSA-45f2-5q3r-xgg6).
func (l *Login) renderPasswordlessPrompt(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	translator := l.getTranslator(r.Context(), authReq)
	data := &passwordlessPromptData{
		userData: l.getUserData(r, authReq, translator, "PasswordlessPrompt.Title", "PasswordlessPrompt.Description", err),
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplPasswordlessPrompt], data, nil)
}
