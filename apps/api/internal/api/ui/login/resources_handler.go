package login

import (
	"net/http"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/assets"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/i18n"
)

type dynamicResourceData struct {
	OrgID         string `schema:"orgId"`
	DefaultPolicy bool   `schema:"default-policy"`
	FileName      string `schema:"filename"`
}

func (l *Login) handleResources() http.Handler {
	return http.FileServer(i18n.LoadFilesystem(i18n.LOGIN))
}

func (l *Login) handleDynamicResources(w http.ResponseWriter, r *http.Request) {
	data := new(dynamicResourceData)
	err := l.getParseData(r, data)
	if err != nil {
		return
	}

	resourceOwner := authz.GetInstance(r.Context()).InstanceID()
	if data.OrgID != "" && !data.DefaultPolicy {
		resourceOwner = data.OrgID
	}

	err = assets.GetAsset(w, r, resourceOwner, data.FileName, l.staticStorage)
	logging.WithFields("file", data.FileName, "org", resourceOwner).OnError(err).Warn("asset in login could not be served")
}
