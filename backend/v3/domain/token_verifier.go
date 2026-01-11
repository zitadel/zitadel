package domain

import "context"

type SessionTokenVerifier func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error)

func noopSessionTokenVerifier() SessionTokenVerifier {
	return func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
		return nil
	}
}
