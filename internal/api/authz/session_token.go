package authz

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	SessionTokenPrefix = "sess_"
	SessionTokenFormat = SessionTokenPrefix + "%s:%s"
)

func SessionTokenVerifier(algorithm crypto.EncryptionAlgorithm) func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
	return func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
		decodedToken, err := base64.RawURLEncoding.DecodeString(sessionToken)
		if err != nil {
			return err
		}
		_, spanPasswordComparison := tracing.NewNamedSpan(ctx, "crypto.CompareHash")
		token, err := algorithm.DecryptString(decodedToken, algorithm.EncryptionKeyID())
		spanPasswordComparison.EndWithError(err)
		if err != nil || token != fmt.Sprintf(SessionTokenFormat, sessionID, tokenID) {
			return zerrors.ThrowPermissionDenied(err, "COMMAND-sGr42", "Errors.Session.Token.Invalid")
		}
		return nil
	}
}
