package handler

import (
	"net"
	"net/http"

	"github.com/caos/zitadel/internal/auth_request/model"
)

const (
	tmplPassword = "password"
)

type passwordData struct {
	Password string `schema:"password"`
}

func (l *Login) renderPassword(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, passwordStep *model.PasswordStep) {
	var errType, errMessage string
	if passwordStep != nil {
		errMessage = "Failure Count: " + string(passwordStep.FailureCount)
	}
	data := userData{
		baseData: l.getBaseData(r, authReq, "Password", errType, errMessage),
		//TODO: Add Username
		//UserName: authReq.UserName,
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplPassword], data, nil)
}

func (l *Login) handlePasswordCheck(w http.ResponseWriter, r *http.Request) {
	data := new(passwordData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	browserInfo := &model.BrowserInfo{RemoteIP: net.IP{}} //TODO: impl
	err = l.authRepo.VerifyPassword(setContext(r.Context(), authReq.UserOrgID), authReq.ID, authReq.UserID, data.Password, browserInfo)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}
