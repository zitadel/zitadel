package oidc

import (
	"errors"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// oidcError ensures [*oidc.Error] and [op.StatusError] types for err.
// It must be used when an error passes the package boundary towards oidc.
// When err is already of the correct type is passed as-is.
// If the err is a Zitadel error, it is transformed with a proper HTTP status code.
// Unknown errors are treated as internal server errors.
func oidcError(err error) error {
	if err == nil {
		return nil
	}
	logging.WithError(err).Warn("OIDC error")
	if errors.Is(err, op.ErrInvalidRefreshToken) {
		err = zerrors.ThrowInvalidArgument(err, "OIDCS-ef2Gi", "Errors.User.RefreshToken.Invalid")
	}
	var (
		sError op.StatusError
		oError *oidc.Error
		zError *zerrors.ZitadelError
	)
	if errors.As(err, &sError) || errors.As(err, &oError) {
		return err
	}

	// here we are encountering an error type that is completely unknown to us.
	if !errors.As(err, &zError) {
		err = zerrors.ThrowInternal(err, "OIDC-AhX2u", "Errors.Internal")
		errors.As(err, &zError)
	}

	statusCode, _ := http_util.ZitadelErrorToHTTPStatusCode(err)
	newOidcErr := oidc.ErrServerError
	if statusCode < 500 {
		newOidcErr = oidc.ErrInvalidRequest
	}

	entry := logging.WithError(err).WithField("status_code", statusCode)
	if statusCode >= 500 {
		entry.Error("OIDC error")
	} else {
		entry.Warn("OIDC error")
	}

	return op.NewStatusError(
		newOidcErr().
			WithParent(err).
			WithDescription(zError.GetMessage()),
		statusCode,
	)
}
