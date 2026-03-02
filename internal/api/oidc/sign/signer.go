package sign

import (
	"context"
	"sync"

	"github.com/go-jose/go-jose/v4"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// SignerFunc is a getter function that allows add-hoc retrieval of the instance's signer.
type SignerFunc func(ctx context.Context) (jose.Signer, jose.SignatureAlgorithm, error)

// GetSignerOnce returns a function which retrieves the instance's signer from the database once.
// Repeated calls of the returned function return the same results.
func GetSignerOnce(
	getActiveSigningWebKey func(ctx context.Context) (*jose.JSONWebKey, error),
) SignerFunc {
	var (
		once    sync.Once
		signer  jose.Signer
		signAlg jose.SignatureAlgorithm
		err     error
	)
	return func(ctx context.Context) (jose.Signer, jose.SignatureAlgorithm, error) {
		once.Do(func() {
			ctx, span := tracing.NewSpan(ctx)
			defer func() { span.EndWithError(err) }()

			var webKey *jose.JSONWebKey
			webKey, err = getActiveSigningWebKey(ctx)
			if err != nil {
				return
			}
			signer, signAlg, err = signerFromWebKey(webKey)
		})
		return signer, signAlg, err
	}
}

func signerFromWebKey(signingKey *jose.JSONWebKey) (jose.Signer, jose.SignatureAlgorithm, error) {
	signAlg := jose.SignatureAlgorithm(signingKey.Algorithm)
	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: signAlg,
			Key:       signingKey,
		},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return nil, "", zerrors.ThrowInternal(err, "OIDC-oaF0s", "Errors.Internal")
	}
	return signer, signAlg, nil
}
