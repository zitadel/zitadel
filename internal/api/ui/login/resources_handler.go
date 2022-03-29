package login

import (
	"net/http"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/assets"
	"github.com/caos/zitadel/internal/domain"
)

type dynamicResourceData struct {
	OrgID         string `schema:"orgId"`
	DefaultPolicy bool   `schema:"default-policy"`
	FileName      string `schema:"filename"`
}

func (l *Login) handleResources(staticDir http.FileSystem) http.Handler {
	return http.FileServer(staticDir)
}

func (l *Login) handleDynamicResources(w http.ResponseWriter, r *http.Request) {
	data := new(dynamicResourceData)
	err := l.getParseData(r, data)
	if err != nil {
		return
	}

	bucketName := domain.IAMID
	if data.OrgID != "" && !data.DefaultPolicy {
		bucketName = data.OrgID
	}

	err = assets.GetAsset(w, r, bucketName, data.FileName, l.staticStorage)
	logging.WithFields("file", data.FileName, "org", bucketName).OnError(err).Warn("asset in login could not be served")
}
