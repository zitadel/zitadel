package domain

import "context"

type SessionTokenVerifier func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error)

type SessionTokenDecryptor func(ctx context.Context, sessionToken string) (sessionID, tokenID string, err error)
