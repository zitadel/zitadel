package handler

import (
	"io"
	"net/http"
	"strconv"

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

	bucketName := domain.IAMID
	if data.OrgID != "" && !data.DefaultPolicy {
		bucketName = data.OrgID
	}
	reader, info, _ := l.staticStorage.GetObject(r.Context(), bucketName, data.FileName)
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return
	}
	w.Header().Set("content-length", strconv.Itoa(int(info.Size)))
	w.Header().Set("content-type", "text/css")
	w.Write(bytes)
}
