//go:generate enumer -type MemberType -trimprefix MemberType

package authz

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
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

func VerifyTokenAndCreateCtxData(ctx context.Context, token, orgID, orgDomain string, t *TokenVerifier) (_ CtxData, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	userID, clientID, agentID, prefLang, resourceOwner, isSystemUser, err := verifyAccessToken(ctx, token, t)
	if err != nil {
		return CtxData{}, err
	}
	var systemMemberships Memberships
	if isSystemUser {
		systemMemberships = make(Memberships, 0, len(t.systemUsers[userID]))
		for _, membership := range t.systemUsers[userID] {
			if membership.MemberType == MemberTypeSystem ||
				membership.MemberType == MemberTypeIAM && GetInstance(ctx).InstanceID() == membership.AggregateID ||
				membership.MemberType == MemberTypeOrganization && orgID == membership.AggregateID {
				systemMemberships = append(systemMemberships, membership)
			}
		}
	}
	var projectID string
	var origins []string
	if clientID != "" {
		projectID, origins, err = t.ProjectIDAndOriginsByClientID(ctx, clientID)
		if err != nil {
			return CtxData{}, errors.ThrowPermissionDenied(err, "AUTH-GHpw2", "could not read projectid by clientid")
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
	// System API calls dont't have a resource owner
	if orgID != "" {
		orgID, err = t.ExistsOrg(ctx, orgID, orgDomain)
		if err != nil {
			return CtxData{}, errors.ThrowPermissionDenied(nil, "AUTH-Bs7Ds", "Organisation doesn't exist")
		}
	}
	return CtxData{
		UserID:            userID,
		OrgID:             orgID,
		ProjectID:         projectID,
		AgentID:           agentID,
		PreferredLanguage: prefLang,
		ResourceOwner:     resourceOwner,
		SystemMemberships: systemMemberships,
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

func GetAllPermissionsFromCtx(ctx context.Context) []string {
	ctxPermission, _ := ctx.Value(allPermissionsKey).([]string)
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
	return errors.ThrowPermissionDenied(nil, "AUTH-DZG21", "Errors.OriginNotAllowed")
}
