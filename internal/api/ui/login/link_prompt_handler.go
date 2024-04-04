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
	Identification string
	Linking        domain.AutoLinkingOption
}

type linkingPromptFormData struct {
	OtherUser bool   `schema:"other"`
	UserID    string `schema:"userId"`
}

func (l *Login) renderLinkingPrompt(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, user *query.User, linking domain.AutoLinkingOption, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	//if singleIDPAllowed(authReq) {
	//	l.handleIDP(w, r, authReq, authReq.AllowedExternalIDPs[0].IDPConfigID)
	//	return
	//}
	translator := l.getTranslator(r.Context(), authReq)
	data := &linkingPromptData{
		Identification: user.PreferredLoginName,
		Linking:        linking,
		userData:       l.getUserData(r, authReq, translator, "Login.Title", "Login.Description", errID, errMessage),
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
		l.renderLinkingPrompt(w, r, authReq, nil, 0, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}
