package signal

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	sig "github.com/zitadel/zitadel/internal/signals"
	objectpb "github.com/zitadel/zitadel/pkg/grpc/object/v2"
	signalpb "github.com/zitadel/zitadel/pkg/grpc/signal/v2"
)

func (s *Server) ListSignals(
	ctx context.Context,
	req *connect.Request[signalpb.ListSignalsRequest],
) (*connect.Response[signalpb.ListSignalsResponse], error) {
	offset := 0
	limit := 100
	if q := req.Msg.GetQuery(); q != nil {
		offset = int(q.GetOffset())
		if q.GetLimit() > 0 && int(q.GetLimit()) < 1000 {
			limit = int(q.GetLimit())
		}
	}

	filters := toSignalFilters(ctx, req.Msg.GetFilters())
	signals, total, err := s.store.SearchSignals(ctx, filters, offset, limit)
	if err != nil {
		return nil, err
	}

	resp := &signalpb.ListSignalsResponse{
		Details: &objectpb.ListDetails{
			TotalResult: uint64(total),
		},
		Signals: make([]*signalpb.Signal, 0, len(signals)),
	}
	for _, s := range signals {
		resp.Signals = append(resp.Signals, recordedSignalToProto(s))
	}
	return connect.NewResponse(resp), nil
}

func (s *Server) AggregateSignals(
	ctx context.Context,
	req *connect.Request[signalpb.AggregateSignalsRequest],
) (*connect.Response[signalpb.AggregateSignalsResponse], error) {
	filters := toSignalFilters(ctx, req.Msg.GetFilters())

	groupBy := sig.AggGroupByField
	if req.Msg.GetGroupBy() == "time_bucket" {
		groupBy = sig.AggGroupByTimeBucket
	}
	metric := sig.AggMetricCount
	if req.Msg.GetMetric() == "distinct_count" {
		metric = sig.AggMetricDistinctCount
	}

	aggReq := sig.AggregationRequest{
		GroupBy:            groupBy,
		FieldName:          req.Msg.GetGroupBy(),
		TimeBucketInterval: req.Msg.GetTimeBucket(),
		Metric:             metric,
	}

	buckets, err := s.store.AggregateSignals(ctx, filters, aggReq)
	if err != nil {
		return nil, err
	}

	resp := &signalpb.AggregateSignalsResponse{
		Buckets: make([]*signalpb.AggregationBucket, 0, len(buckets)),
	}
	for _, b := range buckets {
		resp.Buckets = append(resp.Buckets, &signalpb.AggregationBucket{
			Key:   b.Key,
			Count: b.Value,
		})
	}
	return connect.NewResponse(resp), nil
}

// toSignalFilters converts the proto filters to internal filters.
// The instance ID is always taken from the auth context — never from the request.
func toSignalFilters(ctx context.Context, f *signalpb.SignalFilters) sig.SignalFilters {
	sf := sig.SignalFilters{
		InstanceID: authz.GetInstance(ctx).InstanceID(),
	}
	if f == nil {
		return sf
	}
	sf.UserID = f.GetUserId()
	sf.SessionID = f.GetSessionId()
	sf.IP = f.GetIp()
	sf.Stream = f.GetStream()
	sf.Outcome = f.GetOutcome()
	sf.Operation = f.GetOperation()
	sf.Country = f.GetCountry()
	sf.Resource = f.GetResource()
	sf.Payload = f.GetPayload()
	sf.TraceID = f.GetTraceId()
	sf.SpanID = f.GetSpanId()
	if ts := f.GetAfter(); ts != nil {
		t := ts.AsTime()
		sf.After = &t
	}
	if ts := f.GetBefore(); ts != nil {
		t := ts.AsTime()
		sf.Before = &t
	}
	return sf
}

func (s *Server) ListFindings(
	ctx context.Context,
	req *connect.Request[signalpb.ListFindingsRequest],
) (*connect.Response[signalpb.ListFindingsResponse], error) {
	offset := 0
	limit := 100
	if q := req.Msg.GetQuery(); q != nil {
		offset = int(q.GetOffset())
		if q.GetLimit() > 0 && int(q.GetLimit()) < 1000 {
			limit = int(q.GetLimit())
		}
	}

	filters := toFindingFilters(ctx, req.Msg.GetFilters())
	results, total, err := s.store.SearchFindings(ctx, filters, offset, limit)
	if err != nil {
		return nil, err
	}

	resp := &signalpb.ListFindingsResponse{
		Details: &objectpb.ListDetails{
			TotalResult: uint64(total),
		},
		Findings: make([]*signalpb.FindingWithContext, 0, len(results)),
	}
	for _, r := range results {
		resp.Findings = append(resp.Findings, findingResultToProto(r))
	}
	return connect.NewResponse(resp), nil
}

func (s *Server) AggregateFindings(
	ctx context.Context,
	req *connect.Request[signalpb.AggregateFindingsRequest],
) (*connect.Response[signalpb.AggregateFindingsResponse], error) {
	filters := toFindingFilters(ctx, req.Msg.GetFilters())
	groupBy := req.Msg.GetGroupBy()
	if groupBy == "" {
		groupBy = "name"
	}
	topN := int(req.Msg.GetTopN())
	if topN <= 0 {
		topN = 10
	}

	buckets, err := s.store.AggregateFindings(ctx, filters, groupBy, topN)
	if err != nil {
		return nil, err
	}

	resp := &signalpb.AggregateFindingsResponse{
		Buckets: make([]*signalpb.AggregationBucket, 0, len(buckets)),
	}
	for _, b := range buckets {
		resp.Buckets = append(resp.Buckets, &signalpb.AggregationBucket{
			Key:   b.Key,
			Count: b.Value,
		})
	}
	return connect.NewResponse(resp), nil
}

func toFindingFilters(ctx context.Context, f *signalpb.FindingFilters) sig.FindingFilters {
	ff := sig.FindingFilters{
		SignalFilters: toSignalFilters(ctx, nil),
	}
	if f == nil {
		return ff
	}
	ff.SignalFilters = toSignalFilters(ctx, f.GetSignalFilters())
	ff.FindingName = f.GetFindingName()
	ff.FindingSource = f.GetFindingSource()
	ff.BlockOnly = f.GetBlockOnly()
	ff.ChallengeOnly = f.GetChallengeOnly()
	return ff
}

func findingResultToProto(r sig.FindingResult) *signalpb.FindingWithContext {
	return &signalpb.FindingWithContext{
		Finding: &signalpb.Finding{
			Name:          r.Name,
			Source:        r.Source,
			Message:       r.Message,
			Confidence:    r.Confidence,
			Block:         r.Block,
			Challenge:     r.Challenge,
			ChallengeType: r.ChallengeType,
		},
		SignalTimestamp: timestamppb.New(r.SignalTimestamp),
		UserId:         r.UserID,
		SessionId:      r.SessionID,
		Ip:             r.IP,
		Operation:      r.Operation,
		Stream:         string(r.Stream),
		Outcome:        string(r.Outcome),
		TraceId:        r.TraceID,
	}
}

func recordedSignalToProto(rs sig.RecordedSignal) *signalpb.Signal {
	findings := make([]*signalpb.Finding, 0, len(rs.Findings))
	for _, f := range rs.Findings {
		findings = append(findings, &signalpb.Finding{
			Name:          f.Name,
			Source:        f.Source,
			Message:       f.Message,
			Confidence:    f.Confidence,
			Block:         f.Block,
			Challenge:     f.Challenge,
			ChallengeType: f.ChallengeType,
		})
	}
	return &signalpb.Signal{
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
		ForwardedChain: strings.Join(rs.ForwardedChain, ","),
		Referer:        rs.Referer,
		SecFetchSite:   rs.SecFetchSite,
		IsHttps:        rs.IsHTTPS,
		Findings:       findings,
		Payload:        rs.Payload,
		TraceId:        rs.TraceID,
		SpanId:         rs.SpanID,
	}
}
