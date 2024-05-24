package login

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

const (
	tmplLinkingUserPrompt = "link_user_prompt"
)

type linkingUserPromptData struct {
	userData
	Username string
	Linking  domain.AutoLinkingOption
	UserID   string
}

type linkingUserPromptFormData struct {
	OtherUser bool   `schema:"other"`
	UserID    string `schema:"userID"`
}

func (l *Login) renderLinkingUserPrompt(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, user *query.NotifyUser, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	translator := l.getTranslator(r.Context(), authReq)
	identification := user.PreferredLoginName
	// hide the suffix in case the option is set and the auth request has been started with the primary domain scope
	if authReq.RequestedOrgDomain && authReq.LabelPolicy != nil && authReq.LabelPolicy.HideLoginNameSuffix {
		identification = user.Username
	}
	data := &linkingUserPromptData{
		Username: identification,
		UserID:   user.ID,
		userData: l.getUserData(r, authReq, translator, "LinkingUserPrompt.Title", "LinkingUserPrompt.Description", errID, errMessage),
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplLinkingUserPrompt], data, nil)
}

func (l *Login) handleLinkingUserPrompt(w http.ResponseWriter, r *http.Request) {
	data := new(linkingUserPromptFormData)
	authReq, err := l.ensureAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	if data.OtherUser {
		l.renderExternalNotFoundOption(w, r, authReq, nil, nil, nil, nil)
		return
	}
	err = l.authRepo.SelectUser(r.Context(), authReq.ID, data.UserID, authReq.AgentID)
	if err != nil {
		l.renderLogin(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}
