package handler

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
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

	etag := r.Header.Get("If-None-Match")
	asset, info, err := l.getStatic(r.Context(), bucketName, data.FileName)
	if info != nil && info.ETag == etag {
		w.WriteHeader(304)
		return
	}
	if err != nil {
		return
	}

	w.Header().Set("content-length", strconv.Itoa(int(info.Size)))
	w.Header().Set("content-type", info.ContentType)
	w.Header().Set("ETag", info.ETag)
	w.Write(asset)
}

func (l *Login) getStatic(ctx context.Context, bucketName, fileName string) ([]byte, *domain.AssetInfo, error) {
	s := new(staticAsset)
	key := bucketName + "-" + fileName
	err := l.staticCache.Get(key, s)
	if err == nil && s.Info != nil && (s.Info.Expiration.After(time.Now().Add(-1 * time.Minute))) { //TODO: config?
		return s.Data, s.Info, nil
	}

	info, err := l.staticStorage.GetObjectInfo(ctx, bucketName, fileName)
	if err != nil {
		if caos_errs.IsNotFound(err) {
			return nil, nil, err
		}
		return s.Data, s.Info, err
	}
	if s.Info != nil && s.Info.ETag == info.ETag {
		if info.Expiration.After(s.Info.Expiration) {
			s.Info = info
			l.cacheStatic(bucketName, fileName, s)
		}
		return s.Data, s.Info, nil
	}

	reader, _, err := l.staticStorage.GetObject(ctx, bucketName, fileName)
	if err != nil {
		return s.Data, s.Info, err
	}
	s.Data, err = ioutil.ReadAll(reader)
	if err != nil {
		return nil, nil, err
	}
	s.Info = info
	l.cacheStatic(bucketName, fileName, s)
	return s.Data, s.Info, nil
}

func (l *Login) cacheStatic(bucketName, fileName string, s *staticAsset) {
	key := bucketName + "-" + fileName
	err := l.staticCache.Set(key, &s)
	logging.Log("HANDLER-dfht2").OnError(err).Warnf("caching of asset %s: %s failed", bucketName, fileName)
}

type staticAsset struct {
	Data []byte
	Info *domain.AssetInfo
}
