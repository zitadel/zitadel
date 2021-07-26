package handler

import (
	"net/http"

	"github.com/caos/zitadel/internal/domain"
)

const (
	tmplPasswordlessPrompt = "passwordlessprompt"
)

type passwordlessPromptData struct {
	userData
	SetupEnabled bool
	Required     bool
}

type passwordlessPromptFormData struct {
	Skip bool `schema:"skip"`
}

func (l *Login) handlePasswordlessPrompt(w http.ResponseWriter, r *http.Request) {
	data := new(passwordlessPromptFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if !data.Skip {
		l.renderPasswordlessRegistration(w, r, authReq, "", "", "", "", nil)
		return
	}
	err = l.command.HumanSkipMFAInit(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.handleLogin(w, r)
}

func (l *Login) renderPasswordlessPrompt(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, step *domain.PasswordlessRegistrationPromptStep, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	data := &passwordlessPromptData{
		userData:     l.getUserData(r, authReq, "Passwordless Prompt", errID, errMessage),
		SetupEnabled: step.SetupEnabled,
		Required:     step.Required,
	}

	translator := l.getTranslator(authReq)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplPasswordlessPrompt], data, nil)
}
