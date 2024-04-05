package login

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

const (
	tmplLinkingPrompt = "linking_prompt"
)

type linkingPromptData struct {
	userData
	Username string
	Linking  domain.AutoLinkingOption
	UserID   string
}

type linkingPromptFormData struct {
	OtherUser bool   `schema:"other"`
	UserID    string `schema:"userID"`
}

func (l *Login) renderLinkingPrompt(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, user *query.NotifyUser, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	translator := l.getTranslator(r.Context(), authReq)
	identification := user.PreferredLoginName
	if authReq.LabelPolicy != nil && authReq.LabelPolicy.HideLoginNameSuffix {
		identification = user.Username
	}
	data := &linkingPromptData{
		Username: identification,
		UserID:   user.ID,
		userData: l.getUserData(r, authReq, translator, "Login.Title", "Login.Description", errID, errMessage),
	}
	funcs := map[string]interface{}{
		"hasUsernamePasswordLogin": func() bool {
			return authReq != nil && authReq.LoginPolicy != nil && authReq.LoginPolicy.AllowUsernamePassword
		},
		"hasExternalLogin": func() bool {
			return authReq != nil && authReq.LoginPolicy != nil && authReq.LoginPolicy.AllowExternalIDP && authReq.AllowedExternalIDPs != nil && len(authReq.AllowedExternalIDPs) > 0
		},
		"hasRegistration": func() bool {
			return authReq != nil && authReq.LoginPolicy != nil && authReq.LoginPolicy.AllowRegister
		},
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplLinkingPrompt], data, funcs)
}

func (l *Login) handleLinkingPrompt(w http.ResponseWriter, r *http.Request) {
	data := new(linkingPromptFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	if data.OtherUser {
		l.renderLogin(w, r, authReq, nil)
		return
	}
	err = l.authRepo.SelectUser(r.Context(), authReq.ID, data.UserID, authReq.AgentID)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}
