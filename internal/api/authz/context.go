//go:generate enumer -type MemberType -trimprefix MemberType

package authz

import (
	"context"
	"errors"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/grpc"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type key int

const (
	requestPermissionsKey key = 1
	dataKey               key = 2
	allPermissionsKey     key = 3
	instanceKey           key = 4
)

type CtxData struct {
	UserID            string
	OrgID             string
	ProjectID         string
	AgentID           string
	PreferredLanguage string
	ResourceOwner     string
	SystemMemberships Memberships
}

func (ctxData CtxData) IsZero() bool {
	return ctxData.UserID == "" || ctxData.OrgID == "" && ctxData.SystemMemberships == nil
}

type Grants []*Grant

type Grant struct {
	OrgID string
	Roles []string
}

type Memberships []*Membership

type Membership struct {
	MemberType  MemberType
	AggregateID string
	//ObjectID differs from aggregate id if object is sub of an aggregate
	ObjectID string

	Roles []string
}

type MemberType int32

const (
	MemberTypeUnspecified MemberType = iota
	MemberTypeOrganization
	MemberTypeProject
	MemberTypeProjectGrant
	MemberTypeIAM
	MemberTypeSystem
)

type TokenVerifier interface {
	ExistsOrg(ctx context.Context, id, domain string) (string, error)
	ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (projectID string, origins []string, err error)
	AccessTokenVerifier
	SystemTokenVerifier
}

type AccessTokenVerifier interface {
	VerifyAccessToken(ctx context.Context, token string) (userID, clientID, agentID, prefLan, resourceOwner string, err error)
}

// AccessTokenVerifierFunc implements the SystemTokenVerifier interface so that a function can be used as a AccessTokenVerifier.
type AccessTokenVerifierFunc func(context.Context, string) (string, string, string, string, string, error)

func (a AccessTokenVerifierFunc) VerifyAccessToken(ctx context.Context, token string) (string, string, string, string, string, error) {
	return a(ctx, token)
}

type SystemTokenVerifier interface {
	VerifySystemToken(ctx context.Context, token string, orgID string) (matchingMemberships Memberships, userID string, err error)
}

// SystemTokenVerifierFunc implements the SystemTokenVerifier interface so that a function can be used as a SystemTokenVerifier.
type SystemTokenVerifierFunc func(context.Context, string, string) (Memberships, string, error)

func (s SystemTokenVerifierFunc) VerifySystemToken(ctx context.Context, token string, orgID string) (Memberships, string, error) {
	return s(ctx, token, orgID)
}

func VerifyTokenAndCreateCtxData(ctx context.Context, token, orgID, orgDomain string, t APITokenVerifier) (_ CtxData, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	tokenWOBearer, err := extractBearerToken(token)
	if err != nil {
		return CtxData{}, err
	}
	userID, clientID, agentID, prefLang, resourceOwner, err := t.VerifyAccessToken(ctx, tokenWOBearer)
	var sysMemberships Memberships
	if err != nil && !zerrors.IsUnauthenticated(err) {
		return CtxData{}, err
	}
	if err != nil {
		logging.WithFields("org_id", orgID, "org_domain", orgDomain).WithError(err).Warn("authz: verify access token")
		var sysTokenErr error
		sysMemberships, userID, sysTokenErr = t.VerifySystemToken(ctx, tokenWOBearer, orgID)
		if sysTokenErr != nil || sysMemberships == nil {
			return CtxData{}, zerrors.ThrowUnauthenticated(errors.Join(err, sysTokenErr), "AUTH-7fs1e", "Errors.Token.Invalid")
		}
	}
	var projectID string
	var origins []string
	if clientID != "" {
		projectID, origins, err = t.ProjectIDAndOriginsByClientID(ctx, clientID)
		if err != nil {
			return CtxData{}, zerrors.ThrowPermissionDenied(err, "AUTH-GHpw2", "could not read projectid by clientid")
		}
		// We used to check origins for every token, but service users shouldn't be used publicly (native app / SPA).
		// Therefore, mostly won't send an origin and aren't able to configure them anyway.
		// For the current time we will only check origins for tokens issued to users through apps (code / implicit flow).
		if err := checkOrigin(ctx, origins); err != nil {
			return CtxData{}, err
		}
	}
	if orgID == "" && orgDomain == "" {
		orgID = resourceOwner
	}
	// System API calls don't have a resource owner
	if orgID != "" {
		orgID, err = t.ExistsOrg(ctx, orgID, orgDomain)
		if err != nil {
			return CtxData{}, zerrors.ThrowPermissionDenied(nil, "AUTH-Bs7Ds", "Organisation doesn't exist")
		}
	}
	return CtxData{
		UserID:            userID,
		OrgID:             orgID,
		ProjectID:         projectID,
		AgentID:           agentID,
		PreferredLanguage: prefLang,
		ResourceOwner:     resourceOwner,
		SystemMemberships: sysMemberships,
	}, nil
}

func SetCtxData(ctx context.Context, ctxData CtxData) context.Context {
	return context.WithValue(ctx, dataKey, ctxData)
}

func GetCtxData(ctx context.Context) CtxData {
	ctxData, _ := ctx.Value(dataKey).(CtxData)
	return ctxData
}

func GetRequestPermissionsFromCtx(ctx context.Context) []string {
	ctxPermission, _ := ctx.Value(requestPermissionsKey).([]string)
	return ctxPermission
}

func checkOrigin(ctx context.Context, origins []string) error {
	origin := grpc.GetGatewayHeader(ctx, http_util.Origin)
	if origin == "" {
		origin = http_util.OriginHeader(ctx)
		if origin == "" {
			return nil
		}
	}
	if http_util.IsOriginAllowed(origins, origin) {
		return nil
	}
	return zerrors.ThrowPermissionDenied(nil, "AUTH-DZG21", "Errors.OriginNotAllowed")
}

func extractBearerToken(token string) (part string, err error) {
	parts := strings.Split(token, BearerPrefix)
	if len(parts) != 2 {
		return "", zerrors.ThrowUnauthenticated(nil, "AUTH-7fs1e", "invalid auth header")
	}
	return parts[1], nil
}
