package handler

import (
	"net/http"

)

const (
	tmplUserSelection = "userselection"
)

type userSelectionFormData struct {
	UserSessionID string `schema:"userSessionID"`
}

func (l *Login) renderUserSelection(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession, selectionData *model.UserSelectionData) {
	var errType, errMessage string
	data := userSelectionData{
		baseData: l.getBaseData(r, authSession, "Select User", errType, errMessage),
		Users:    selectionData.Users,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplUserSelection], data, nil)
}

func (l *Login) handleSelectUser(w http.ResponseWriter, r *http.Request) {
	data := new(userSelectionFormData)
	authSession, err := l.getAuthSessionAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	browserInfo := &model.BrowserInformation{RemoteIP: &model.IP{}} //TODO: impl
	if data.UserSessionID == "0" {
		l.renderLogin(w, r, authSession, nil)
		return
	}
	authSession, err = l.service.Auth.SelectUser(r.Context(), authSession, data.UserSessionID, browserInfo)
	if err != nil {
		l.renderError(w, r, authSession, err)
		return
	}
	l.renderNextStep(w, r, authSession)
}
