package handler

import (
	"github.com/caos/zitadel/internal/domain"
	"net/http"
)

const (
	tmplMFAInitDone = "mfainitdone"
)

type mfaInitDoneData struct {
}

func (l *Login) renderMFAInitDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, data *mfaDoneData) {
	var errType, errMessage string
	data.baseData = l.getBaseData(r, authReq, "MFA Init Done", errType, errMessage)
	data.profileData = l.getProfileData(authReq)
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMFAInitDone], data, nil)
}
