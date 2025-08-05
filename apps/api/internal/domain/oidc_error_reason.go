package domain

import (
	"errors"

	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type OIDCErrorReason int32

const (
	OIDCErrorReasonUnspecified OIDCErrorReason = iota
	OIDCErrorReasonInvalidRequest
	OIDCErrorReasonUnauthorizedClient
	OIDCErrorReasonAccessDenied
	OIDCErrorReasonUnsupportedResponseType
	OIDCErrorReasonInvalidScope
	OIDCErrorReasonServerError
	OIDCErrorReasonTemporaryUnavailable
	OIDCErrorReasonInteractionRequired
	OIDCErrorReasonLoginRequired
	OIDCErrorReasonAccountSelectionRequired
	OIDCErrorReasonConsentRequired
	OIDCErrorReasonInvalidRequestURI
	OIDCErrorReasonInvalidRequestObject
	OIDCErrorReasonRequestNotSupported
	OIDCErrorReasonRequestURINotSupported
	OIDCErrorReasonRegistrationNotSupported
	OIDCErrorReasonInvalidGrant
)

func OIDCErrorReasonFromError(err error) OIDCErrorReason {
	if errors.Is(err, oidc.ErrInvalidRequest()) {
		return OIDCErrorReasonInvalidRequest
	}
	if errors.Is(err, oidc.ErrInvalidGrant()) {
		return OIDCErrorReasonInvalidGrant
	}
	if zerrors.IsPreconditionFailed(err) {
		return OIDCErrorReasonAccessDenied
	}
	if zerrors.IsInternal(err) {
		return OIDCErrorReasonServerError
	}
	return OIDCErrorReasonUnspecified
}
