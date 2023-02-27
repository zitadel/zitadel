package assets

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/static"
)

func (h *Handler) UploadApplicationLightIcon(appID string) Uploader {
	return &applicationIconUploader{h.idGenerator, appID, false, []string{"image/"}, 1 << 19}
}

func (h *Handler) UploadApplicationDarkIcon(appID string) Uploader {
	return &applicationIconUploader{h.idGenerator, appID, true, []string{"image/"}, 1 << 19}
}

type applicationIconUploader struct {
	idGenerator   id.Generator
	appID string
	darkMode      bool
	contentTypes  []string
	maxSize       int64
}

func (a *applicationIconUploader) ContentTypeAllowed(contentType string) bool {
	for _, ct := range a.contentTypes {
		if strings.HasPrefix(contentType, ct) {
			return true
		}
	}
	return false
}

func (a *applicationIconUploader) ObjectType() static.ObjectType {
	return static.ObjectTypeStyling
}

func (a *applicationIconUploader) MaxFileSize() int64 {
	return a.maxSize
}

func (a *applicationIconUploader) ObjectName(ctx authz.CtxData) (string, error) {
	suffixID, err := a.idGenerator.Next()
	if err != nil {
		return "", err
	}

	prefix := domain.GetApplicationIconAssetPath(ctx.ProjectID, a.appID)
	if a.darkMode {
		return prefix + "-" + domain.Dark + "-" + suffixID, nil
	}
	return prefix + "-" + suffixID, nil
}

func (a *applicationIconUploader) ResourceOwner(instance authz.Instance, ctxData authz.CtxData) string {
	return ctxData.ResourceOwner
}

func (a *applicationIconUploader) UploadAsset(ctx context.Context, appID string, upload *command.AssetUpload, commands *command.Commands) error {
	_, err := commands.AddApplicationIcon(ctx, authz.GetCtxData(ctx).ProjectID, appID, a.darkMode, upload)
	return err
	
}

func (h *Handler) GetApplicationIcon() Downloader {
	return &applicationIconDownloader{query: h.query, darkMode: false}
}

func (h *Handler) GetApplicationIconDark() Downloader {
	return &applicationIconDownloader{query: h.query, darkMode: true}
}

type applicationIconDownloader struct {
	query         *query.Queries
	darkMode      bool
}

func (a *applicationIconDownloader) ObjectName(ctx context.Context, appID string) (string, error) {
	app, err := getApplication(ctx, appID, a.query)
	if err != nil {
		return "", nil
	}
	if a.darkMode {
		return app.DarkIconURL, nil
	}
	return app.LightIconURL, nil
}

func (a *applicationIconDownloader) ResourceOwner(ctx context.Context, _ string) string {
	return authz.GetCtxData(ctx).ResourceOwner
}

func getApplication(ctx context.Context, appID string, queries *query.Queries) (*query.App, error) {
	return queries.AppByID(ctx, appID, false)
}