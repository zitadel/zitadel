package handler

import (
	"net/http"

	"github.com/caos/zitadel/internal/auth_request/model"
)

const (
	tmplMFAInitDone = "mfainitdone"
)

type mfaInitDoneData struct {
}

func (l *Login) renderMFAInitDone(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, data *mfaDoneData) {
	var errType, errMessage string
	data.baseData = l.getBaseData(r, authReq, "MFA Init Done", errType, errMessage)
	data.profileData = l.getProfileData(authReq)
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMFAInitDone], data, nil)
}
