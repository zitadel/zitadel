package signal

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/signals"
	signal "github.com/zitadel/zitadel/pkg/grpc/signal/v2"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2"
)

func (s *Server) ListSignals(ctx context.Context, req *connect.Request[signal.ListSignalsRequest]) (*connect.Response[signal.ListSignalsResponse], error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	filters := filtersFromProto(instanceID, req.Msg.GetFilters())

	if err := authorizeSignalAccess(ctx, filters); err != nil {
		return nil, err
	}

	offset, limit := listQueryToOffsetLimit(req.Msg.GetQuery())

	results, total, err := s.reader.SearchSignals(ctx, filters, offset, limit)
	if err != nil {
		return nil, err
	}

	pbSignals := make([]*signal.Signal, len(results))
	for i := range results {
		pbSignals[i] = signalToProto(&results[i])
	}

	return connect.NewResponse(&signal.ListSignalsResponse{
		Details: &object.ListDetails{
			TotalResult: uint64(total),
		},
		Signals: pbSignals,
	}), nil
}

func (s *Server) AggregateSignals(ctx context.Context, req *connect.Request[signal.AggregateSignalsRequest]) (*connect.Response[signal.AggregateSignalsResponse], error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	filters := filtersFromProto(instanceID, req.Msg.GetFilters())

	if err := authorizeSignalAccess(ctx, filters); err != nil {
		return nil, err
	}

	aggReq := signals.AggregateRequest{
		GroupBy:          req.Msg.GetGroupBy(),
		TimeBucket:       req.Msg.GetTimeBucket(),
		Metric:           req.Msg.GetMetric(),
		SecondaryGroupBy: req.Msg.GetSecondaryGroupBy(),
		Limit:            int(req.Msg.GetLimit()),
	}

	buckets, err := s.reader.AggregateSignals(ctx, filters, aggReq)
	if err != nil {
		return nil, err
	}

	pbBuckets := make([]*signal.AggregationBucket, len(buckets))
	for i, b := range buckets {
		pbBuckets[i] = &signal.AggregationBucket{
			Key:    b.Key,
			Count:  b.Count,
			Series: b.Series,
			Value:  b.Value,
		}
	}

	return connect.NewResponse(&signal.AggregateSignalsResponse{
		Buckets: pbBuckets,
	}), nil
}

func listQueryToOffsetLimit(q *object.ListQuery) (offset, limit int) {
	if q == nil {
		return 0, 20
	}
	offset = int(q.GetOffset())
	limit = int(q.GetLimit())
	if limit <= 0 {
		limit = 20
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}
	if offset > 10000 {
		offset = 10000
	}
	return offset, limit
}

func filtersFromProto(instanceID string, pf *signal.SignalFilters) signals.SignalFilters {
	f := signals.SignalFilters{
		InstanceID: instanceID,
		Fields:     make(map[string]string),
	}
	if pf == nil {
		return f
	}

	// Map all proto filter fields into the generic Fields map.
	// The field registry in fields.go determines how each is applied.
	set := func(col, val string) {
		if val != "" {
			f.Fields[col] = val
		}
	}
	set("user_id", pf.GetUserId())
	set("session_id", pf.GetSessionId())
	set("ip", pf.GetIp())
	set("stream", pf.GetStream())
	set("outcome", pf.GetOutcome())
	set("operation", pf.GetOperation())
	set("country", pf.GetCountry())
	set("resource", pf.GetResource())
	set("org_id", pf.GetOrgId())
	set("project_id", pf.GetProjectId())
	set("client_id", pf.GetClientId())
	set("payload", pf.GetPayload())
	set("trace_id", pf.GetTraceId())
	set("span_id", pf.GetSpanId())
	set("user_agent", pf.GetUserAgent())
	set("fingerprint_id", pf.GetFingerprintId())
	set("caller_id", pf.GetCallerId())
	set("referer", pf.GetReferer())
	set("accept_language", pf.GetAcceptLanguage())
	set("forwarded_chain", pf.GetForwardedChain())
	set("sec_fetch_site", pf.GetSecFetchSite())

	if pf.GetAfter() != nil {
		t := pf.GetAfter().AsTime()
		f.After = &t
	}
	if pf.GetBefore() != nil {
		t := pf.GetBefore().AsTime()
		f.Before = &t
	}
	return f
}

func signalToProto(rs *signals.RecordedSignal) *signal.Signal {
	findings := make([]*signal.Finding, len(rs.Findings))
	for i := range rs.Findings {
		findings[i] = findingToProto(&rs.Findings[i])
	}

	return &signal.Signal{
		InstanceId:     rs.InstanceID,
		UserId:         rs.UserID,
		CallerId:       rs.CallerID,
		SessionId:      rs.SessionID,
		FingerprintId:  rs.FingerprintID,
		Operation:      rs.Operation,
		Stream:         string(rs.Stream),
		Resource:       rs.Resource,
		Outcome:        string(rs.Outcome),
		CreatedAt:      timestamppb.New(rs.Timestamp),
		Ip:             rs.IP,
		UserAgent:      rs.UserAgent,
		AcceptLanguage: rs.AcceptLanguage,
		Country:        rs.Country,
		ForwardedChain: strings.Join(rs.ForwardedChain, ", "),
		Referer:        rs.Referer,
		SecFetchSite:   rs.SecFetchSite,
		IsHttps:        rs.IsHTTPS,
		Findings:       findings,
		Payload:        rs.Payload,
		TraceId:        rs.TraceID,
		SpanId:         rs.SpanID,
		OrgId:          rs.OrgID,
		ProjectId:      rs.ProjectID,
		ClientId:       rs.ClientID,
		DurationMs:     rs.DurationMs,
	}
}

func findingToProto(f *signals.RecordedFinding) *signal.Finding {
	return &signal.Finding{
		Name:          f.Name,
		Source:        f.Source,
		Message:       f.Message,
		Confidence:    f.Confidence,
		Block:         f.Block,
		Challenge:     f.Challenge,
		ChallengeType: f.ChallengeType,
	}
}

// authorizeSignalAccess is a placeholder for future fine-grained
// authorization (e.g. org-scoped queries). Currently the proto
// auth_option enforces iam.read which limits access to instance
// administrators, so no additional restriction is applied here.
func authorizeSignalAccess(_ context.Context, _ signals.SignalFilters) error {
	return nil
}
