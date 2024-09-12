package activity

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/telemetry/logs/record"
	"google.golang.org/grpc/codes"
	"strconv"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	PathKey          = "zitadel-activity-path"
	RequestMethodKey = "zitadel-activity-request-method"
)

type TriggerMethod int

const (
	Unspecified TriggerMethod = iota
	ResourceAPI
	OIDCAccessToken
	OIDCRefreshToken
	SessionAPI
	SAMLResponse
)

func (t TriggerMethod) String() string {
	switch t {
	case Unspecified:
		return "unspecified"
	case ResourceAPI:
		return "resourceAPI"
	case OIDCRefreshToken:
		return "refreshToken"
	case OIDCAccessToken:
		return "accessToken"
	case SessionAPI:
		return "sessionAPI"
	case SAMLResponse:
		return "samlResponse"
	default:
		return "unknown"
	}
}

var _ record.StreamRecord = (*activityRecord)(nil)

type activityRecord struct {
	*record.BaseStreamRecord
	OrgID, UserID, Domain, Method, Path, RequestMethod string
	Trigger                                            TriggerMethod
	GRPCStatus                                         *codes.Code
	HTTPStatus                                         *int
	IsSystemUser                                       bool
}

func (ar *activityRecord) Base() *record.BaseStreamRecord {
	return ar.BaseStreamRecord
}

func (ar *activityRecord) Fields() logrus.Fields {
	fields := logrus.Fields{
		"stream":        ar.Stream,
		"version":       ar.Version,
		"instance":      ar.ZITADELInstanceID,
		"org":           ar.OrgID,
		"user":          ar.UserID,
		"domain":        ar.Domain,
		"trigger":       ar.Trigger.String(),
		"method":        ar.Method,
		"path":          ar.Path,
		"requestMethod": ar.RequestMethod,
		"isSystemUser":  strconv.FormatBool(ar.IsSystemUser),
	}
	if ar.GRPCStatus != nil {
		fields["grpcStatus"] = ar.GRPCStatus.String()
	}
	if ar.HTTPStatus != nil {
		fields["httpStatus"] = strconv.FormatInt(int64(*ar.HTTPStatus), 10)
	}
	for key, value := range fields {
		if value == "" {
			delete(fields, key)
		}
	}
	return fields
}

// Trigger is used to log a specific events for a user (e.g. session or oidc token creation)
func Trigger(ctx context.Context, orgID, userID string, trigger TriggerMethod, reducer func(ctx context.Context, r eventstore.QueryReducer) error) {
	if orgID == "" && userID != "" {
		orgID = getOrgOfUser(ctx, userID, reducer)
	}
	ai := ActivityInfoFromContext(ctx)
	triggerLog(
		authz.GetInstance(ctx).InstanceID(),
		orgID,
		userID,
		http_utils.DomainContext(ctx).Origin(), // TODO: origin?
		trigger,
		ai.Method,
		ai.Path,
		ai.RequestMethod,
		nil,
		nil,
		authz.GetCtxData(ctx).SystemMemberships != nil,
	)
}

func TriggerGRPCWithContext(ctx context.Context, trigger TriggerMethod) {
	ai := ActivityInfoFromContext(ctx)
	httpStatus := runtime.HTTPStatusFromCode(ai.GRPCStatus)
	triggerLog(
		authz.GetInstance(ctx).InstanceID(),
		authz.GetCtxData(ctx).OrgID,
		authz.GetCtxData(ctx).UserID,
		http_utils.DomainContext(ctx).Origin(), // TODO: origin?
		trigger,
		ai.Method,
		ai.Path,
		ai.RequestMethod,
		&ai.GRPCStatus,
		&httpStatus,
		authz.GetCtxData(ctx).SystemMemberships != nil,
	)
}

func triggerLog(
	instanceID, orgID, userID, domain string,
	trigger TriggerMethod,
	method, path, requestMethod string,
	grpcStatus *codes.Code,
	httpStatus *int,
	isSystemUser bool,
) {
	r := &activityRecord{
		BaseStreamRecord: &record.BaseStreamRecord{
			Version:           "v1",
			Stream:            record.StreamActivity,
			ZITADELInstanceID: instanceID,
		},
		OrgID:         orgID,
		UserID:        userID,
		Domain:        domain,
		Method:        method,
		Path:          path,
		RequestMethod: requestMethod,
		Trigger:       trigger,
		GRPCStatus:    grpcStatus,
		HTTPStatus:    httpStatus,
		IsSystemUser:  isSystemUser,
	}
	logging.New().WithFields(r.Fields()).Info(record.StreamActivity)
}

func getOrgOfUser(ctx context.Context, userID string, reducer func(ctx context.Context, r eventstore.QueryReducer) error) string {
	org := &orgIDOfUser{userID: userID}
	err := reducer(ctx, org)
	if err != nil {
		logging.WithError(err).Error("could not get org id of user for trigger log")
		return ""
	}
	return org.orgID
}

type orgIDOfUser struct {
	eventstore.WriteModel

	userID string
	orgID  string
}

func (u *orgIDOfUser) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		OrderDesc().
		Limit(1).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(u.userID).
		Builder()
}

func (u *orgIDOfUser) Reduce() error {
	if len(u.Events) == 0 {
		return nil
	}
	u.orgID = u.Events[0].Aggregate().ResourceOwner
	return nil
}
