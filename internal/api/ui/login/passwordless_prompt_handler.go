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

type passwordlessPromptFormData struct{}

func (l *Login) handlePasswordlessPrompt(w http.ResponseWriter, r *http.Request) {
	data := new(passwordlessPromptFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.renderPasswordlessRegistration(w, r, authReq, "", "", "", "", 0, nil)
}

func (l *Login) renderPasswordlessPrompt(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, err error) {
	translator := l.getTranslator(r.Context(), authReq)
	data := &passwordlessPromptData{
		userData: l.getUserData(r, authReq, translator, "PasswordlessPrompt.Title", "PasswordlessPrompt.Description", err),
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplPasswordlessPrompt], data, nil)
}
