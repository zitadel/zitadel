package handler

import (
	"github.com/caos/zitadel/internal/auth_request/model"
	"net/http"
)

const (
	tmplUserSelection = "userselection"
)

type userSelectionFormData struct {
	UserSessionID string `schema:"userSessionID"`
}

func (l *Login) renderUserSelection(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, selectionData *model.SelectUserStep) {
	var errType, errMessage string
	data := userSelectionData{
		baseData: l.getBaseData(r, authReq, "Select User", errType, errMessage),
		Users:    selectionData.Users,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplUserSelection], data, nil)
}

func (l *Login) handleSelectUser(w http.ResponseWriter, r *http.Request) {
	data := new(userSelectionFormData)
	authSession, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	if data.UserSessionID == "0" {
		l.renderLogin(w, r, authSession, nil)
		return
	}
	//TODO: Choose User
	//authSession, err = l.authRepo.SelectUser(r.Context(), authSession, data.UserSessionID, browserInfo)
	//if err != nil {
	//	l.renderError(w, r, authSession, err)
	//	return
	//}
	l.renderNextStep(w, r, authSession)
}
