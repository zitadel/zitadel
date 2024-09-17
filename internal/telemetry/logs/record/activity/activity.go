package activity

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/pkg/streams"
	"google.golang.org/grpc/codes"
	"strconv"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/logs/record"
	"github.com/zitadel/zitadel/pkg/streams/activity"
)

var _ record.StreamRecord = (*activityRecord)(nil)

type activityRecord struct {
	*record.BaseStreamRecord
	OrgID, UserID, Domain, Method, Path, RequestMethod string
	Trigger                                            activity.TriggerMethod
	GRPCStatus                                         *codes.Code
	HTTPStatus                                         *int
	IsSystemUser                                       bool
}

func (ar *activityRecord) Base() *record.BaseStreamRecord {
	return ar.BaseStreamRecord
}

func (ar *activityRecord) Fields() logrus.Fields {
	fields := logrus.Fields{
		string(streams.LogFieldKeyStream):         ar.Stream,
		string(streams.LogFieldKeyVersion):        ar.Version,
		string(streams.LogFieldKeyInstanceID):     ar.ZITADELInstanceID,
		string(activity.LogFieldKeyOrgID):         ar.OrgID,
		string(activity.LogFieldKeyUserID):        ar.UserID,
		string(activity.LogFieldKeyDomain):        ar.Domain,
		string(activity.LogFieldKeyTrigger):       ar.Trigger.String(),
		string(activity.LogFieldKeyMethod):        ar.Method,
		string(activity.LogFieldKeyPath):          ar.Path,
		string(activity.LogFieldKeyRequestMethod): ar.RequestMethod,
		string(activity.LogFieldKeyIsSystemUser):  strconv.FormatBool(ar.IsSystemUser),
	}
	if ar.GRPCStatus != nil {
		fields[string(activity.LogFieldKeyGRPCStatus)] = ar.GRPCStatus.String()
	}
	if ar.HTTPStatus != nil {
		fields[string(activity.LogFieldKeyHTTPStatus)] = strconv.FormatInt(int64(*ar.HTTPStatus), 10)
	}
	for key, value := range fields {
		if value == "" {
			delete(fields, key)
		}
	}
	return fields
}

// Trigger is used to log a specific events for a user (e.g. session or oidc token creation)
func Trigger(ctx context.Context, orgID, userID string, trigger activity.TriggerMethod, reducer func(ctx context.Context, r eventstore.QueryReducer) error) {
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

func TriggerGRPCWithContext(ctx context.Context, trigger activity.TriggerMethod) {
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
	trigger activity.TriggerMethod,
	method, path, requestMethod string,
	grpcStatus *codes.Code,
	httpStatus *int,
	isSystemUser bool,
) {
	r := &activityRecord{
		BaseStreamRecord: &record.BaseStreamRecord{
			Version:           string(streams.LogFieldValueStreamVersion),
			Stream:            streams.LogFieldValueStreamActivity,
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
	logging.New().WithFields(r.Fields()).Info(streams.LogFieldValueStreamActivity)
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
