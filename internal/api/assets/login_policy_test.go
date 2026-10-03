package assets

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
)

func TestFontUploadSizeLimit(t *testing.T) {
	for _, factory := range []struct {
		name string
		new  func(*Handler) Uploader
	}{
		{"instance", (*Handler).UploadDefaultLabelPolicyFont},
		{"organization", (*Handler).UploadOrgLabelPolicyFont},
	} {
		for _, tt := range []struct {
			name   string
			limit  int64
			size   int
			status int
		}{
			{"default boundary", 524288, 524288, http.StatusOK},
			{"default exceeded", 524288, 524289, http.StatusBadRequest},
			{"increased limit", 3145728, 524289, http.StatusOK},
			{"increased boundary", 3145728, 3145728, http.StatusOK},
			{"increased exceeded", 3145728, 3145729, http.StatusBadRequest},
			{"decreased exceeded", 262144, 262145, http.StatusBadRequest},
		} {
			t.Run(factory.name+"/"+tt.name, func(t *testing.T) {
				t.Parallel()
				h := &Handler{
					maxFontSize: tt.limit,
					errorHandler: func(w http.ResponseWriter, r *http.Request, err error, code int) {
						http.Error(w, err.Error(), code)
					},
				}
				// Keep the real font type and size validation, replacing only the
				// persistence operations that would require a database.
				uploader := &fontUploadRecorder{Uploader: factory.new(h)}
				font := make([]byte, tt.size)
				copy(font, "OTTO") // OpenType signature for content-type detection.
				var body bytes.Buffer
				writer := multipart.NewWriter(&body)
				part, err := writer.CreateFormFile("file", "font.otf")
				require.NoError(t, err)
				_, err = part.Write(font)
				require.NoError(t, err)
				require.NoError(t, writer.Close())
				req := httptest.NewRequest(http.MethodPost, "/", &body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				t.Cleanup(func() {
					if req.MultipartForm != nil {
						require.NoError(t, req.MultipartForm.RemoveAll())
					}
				})
				response := httptest.NewRecorder()
				UploadHandleFunc(h, uploader)(response, req)
				require.Equal(t, tt.status, response.Code, response.Body.String())
				if tt.status == http.StatusOK {
					require.Equal(t, font, uploader.data)
				} else {
					require.Nil(t, uploader.data)
					require.Contains(t, response.Body.String(), "file too big")
				}
			})
		}
	}
}

type fontUploadRecorder struct {
	Uploader
	data []byte
}

func (u *fontUploadRecorder) ResourceOwner(authz.Instance, authz.CtxData) string {
	return "owner"
}

func (u *fontUploadRecorder) ObjectName(authz.CtxData) (string, error) {
	return "font", nil
}

func (u *fontUploadRecorder) UploadAsset(_ context.Context, _ string, asset *command.AssetUpload, _ *command.Commands) error {
	var err error
	u.data, err = io.ReadAll(asset.File)
	return err
}
