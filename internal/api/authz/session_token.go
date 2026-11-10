package authz

import (
	"context"
	"fmt"
	"strings"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	SessionTokenPrefix = "sess_"
	SessionTokenFormat = SessionTokenPrefix + "%s:%s"
)

func SessionTokenVerifier(algorithm crypto.AuthAlgorithm) func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
	return func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
		_, spanPasswordComparison := tracing.NewNamedSpan(ctx, "crypto.CompareHash")
		token, err := algorithm.DecryptToken(sessionToken)
		spanPasswordComparison.EndWithError(err)
		if err != nil || token != fmt.Sprintf(SessionTokenFormat, sessionID, tokenID) {
			return zerrors.ThrowPermissionDenied(err, "COMMAND-sGr42", "Errors.Session.Token.Invalid")
		}
		return nil
	}
}

func SessionTokenDecryptor(algorithm crypto.AuthAlgorithm) func(ctx context.Context, sessionToken string) (sessionID, tokenID string, err error) {
	return func(ctx context.Context, sessionToken string) (sessionID, tokenID string, err error) {
		_, spanPasswordComparison := tracing.NewNamedSpan(ctx, "crypto.CompareHash")
		token, err := algorithm.DecryptToken(sessionToken)
		spanPasswordComparison.EndWithError(err)
		if err != nil {
			return "", "", zerrors.ThrowPermissionDenied(err, "AUTHZ-sGr42", "Errors.Session.Token.Invalid")
		}

		n, err := fmt.Sscanf(strings.ReplaceAll(token, ":", " "), strings.ReplaceAll(SessionTokenFormat, ":", " "), &sessionID, &tokenID)
		if err != nil || n != 2 {
			return "", "", zerrors.ThrowInvalidArgument(err, "AUTHZ-3uike", "Errors.Session.Token.Invalid")
		}
		return sessionID, tokenID, nil
	}
}
