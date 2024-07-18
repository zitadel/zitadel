package action

import (
	"context"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	action "github.com/zitadel/zitadel/pkg/grpc/resources/action/v3alpha"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
)

func (s *Server) GetTarget(ctx context.Context, req *action.GetTargetRequest) (*action.GetTargetResponse, error) {
	if err := checkExecutionEnabled(ctx); err != nil {
		return nil, err
	}

	resp, err := s.query.GetTargetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &action.GetTargetResponse{
		Target: targetToPb(resp),
	}, nil
}

func (s *Server) SearchTargets(ctx context.Context, req *action.SearchTargetsRequest) (*action.SearchTargetsResponse, error) {
	if err := checkExecutionEnabled(ctx); err != nil {
		return nil, err
	}

	queries, err := listTargetsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchTargets(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &action.SearchTargetsResponse{
		Result:  targetsToPb(resp.Targets),
		Details: object.ToListDetails(resp.SearchResponse),
	}, nil
}

func (s *Server) SearchExecutions(ctx context.Context, req *action.SearchExecutionsRequest) (*action.SearchExecutionsResponse, error) {
	if err := checkExecutionEnabled(ctx); err != nil {
		return nil, err
	}

	queries, err := listExecutionsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchExecutions(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &action.SearchExecutionsResponse{
		Result:  executionsToPb(resp.Executions),
		Details: object.ToListDetails(resp.SearchResponse),
	}, nil
}

func targetToPb(t *query.Target) *action.GetTarget {
	target := &action.GetTarget{
		Details: object.DomainToDetailsPb(&t.ObjectDetails),
		Target: &action.Target{
			Name:     t.Name,
			Timeout:  durationpb.New(t.Timeout),
			Endpoint: t.Endpoint,
		},
	}
	switch t.TargetType {
	case domain.TargetTypeWebhook:
		target.Target.TargetType = &action.Target_RestWebhook{RestWebhook: &action.SetRESTWebhook{InterruptOnError: t.InterruptOnError}}
	case domain.TargetTypeCall:
		target.Target.TargetType = &action.Target_RestCall{RestCall: &action.SetRESTCall{InterruptOnError: t.InterruptOnError}}
	case domain.TargetTypeAsync:
		target.Target.TargetType = &action.Target_RestAsync{RestAsync: &action.SetRESTAsync{}}
	default:
		target.Target.TargetType = nil
	}
	return target
}
