package authz

import (
	"context"

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
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	userID, agentID, clientID, prefLang, resourceOwner, err = a.authZRepo.VerifyAccessToken(ctx, token, "", GetInstance(ctx).ProjectID())
	return userID, clientID, agentID, prefLang, resourceOwner, err
}

type client struct {
	name string
}
