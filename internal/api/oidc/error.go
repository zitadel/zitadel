package oidc

import (
	"context"
	"errors"
	"net/http"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// oidcError ensures [*oidc.Error] and [op.StatusError] types for err.
// It must be used when an error passes the package boundary towards oidc.
// When err is already of the correct type is passed as-is.
// If the err is a Zitadel error, it is transformed with a proper HTTP status code.
// Unknown errors are treated as internal server errors.
func oidcError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
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

	statusCode, _ := http_util.ZitadelErrorToHTTPStatusCode(ctx, err)
	newOidcErr := oidc.ErrServerError
	if statusCode < 500 {
		newOidcErr = oidc.ErrInvalidRequest
	}
	oidcErr := newOidcErr().WithParent(err)
	oidcErr.Description = zError.GetMessage()
	return op.NewStatusError(oidcErr, statusCode)
}

// asInvalidGrant maps errors which mean that the grant presented at the token
// endpoint - an authorization code or a refresh token - is unknown, malformed,
// expired or revoked, to an [*oidc.Error] with the invalid_grant error code.
// RFC 6749, section 5.2 defines invalid_grant for exactly those cases, and
// clients rely on it to tell a grant they have to renew from a request of their
// own that was malformed. Without this mapping [oidcError] reports them as
// invalid_request, because it maps every error below status 500 to that code.
// Any other error is returned unchanged, so that internal failures keep their
// status instead of telling the client its grant is gone.
func asInvalidGrant(ctx context.Context, err error, description string) error {
	if err == nil || !isInvalidGrant(err) {
		return err
	}
	oidcErr := oidc.ErrInvalidGrant().
		WithParent(err).
		WithReturnParentToClient(authz.GetFeatures(ctx).DebugOIDCParentError)
	oidcErr.Description = description
	return oidcErr
}

func isInvalidGrant(err error) bool {
	return errors.Is(err, op.ErrInvalidRefreshToken) ||
		errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "", "Errors.OIDCSession.RefreshTokenInvalid")) ||
		errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "", "Errors.OIDCSession.Token.Invalid")) ||
		errors.Is(err, zerrors.ThrowInvalidArgument(nil, "", "Errors.User.RefreshToken.Invalid")) ||
		errors.Is(err, zerrors.ThrowNotFound(nil, "", "Errors.User.RefreshToken.Invalid")) ||
		errors.Is(err, zerrors.ThrowInvalidArgument(nil, "", "Errors.User.Code.Invalid")) ||
		errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "", "Errors.AuthRequest.InvalidCode")) ||
		errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "", "Errors.AuthRequest.NoCode"))
}

func writeRecoverError(w http.ResponseWriter, r *http.Request, err error) {
	op.WriteError(w, r, err, logging.FromCtx(r.Context()))
}
