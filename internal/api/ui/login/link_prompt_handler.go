package login

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

const (
	tmplLinkingPrompt = "linking_prompt"
)

type LinkingPromptData struct {
	userData
	Identification string
	Linking        domain.AutoLinkingOption
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
	data := &LinkingPromptData{
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
