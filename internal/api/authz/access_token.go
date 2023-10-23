package authz

import (
	"context"
	"strings"

	zitadel_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

const (
	BearerPrefix = "Bearer "
)

type MembershipsResolver interface {
	SearchMyMemberships(ctx context.Context, orgID string, shouldTriggerBulk bool) ([]*Membership, error)
}

type authZRepo interface {
	MembershipsResolver
	VerifyAccessToken(ctx context.Context, token, verifierClientID, projectID string) (userID, agentID, clientID, prefLang, resourceOwner string, err error)
	VerifierClientID(ctx context.Context, name string) (clientID, projectID string, err error)
	ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (projectID string, origins []string, err error)
	ExistsOrg(ctx context.Context, id, domain string) (string, error)
}

var _ AccessTokenVerifier = (*AccessTokenVerifierFromRepo)(nil)

type AccessTokenVerifierFromRepo struct {
	authZRepo authZRepo
}

func StartAccessTokenVerifierFromRepo(authZRepo authZRepo) *AccessTokenVerifierFromRepo {
	return &AccessTokenVerifierFromRepo{authZRepo: authZRepo}
}

func (a *AccessTokenVerifierFromRepo) VerifyAccessToken(ctx context.Context, token string) (userID, clientID, agentID, prefLang, resourceOwner string, err error) {
	userID, agentID, clientID, prefLang, resourceOwner, err = a.authZRepo.VerifyAccessToken(ctx, token, "", GetInstance(ctx).ProjectID())
	return userID, clientID, agentID, prefLang, resourceOwner, err
}

type client struct {
	name string
}

func verifyAccessToken(ctx context.Context, token string, t AccessTokenVerifier) (userID, clientID, agentID, prefLan, resourceOwner string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	parts := strings.Split(token, BearerPrefix)
	if len(parts) != 2 {
		return "", "", "", "", "", zitadel_errors.ThrowUnauthenticated(nil, "AUTH-7fs1e", "invalid auth header")
	}
	return t.VerifyAccessToken(ctx, parts[1])
}
