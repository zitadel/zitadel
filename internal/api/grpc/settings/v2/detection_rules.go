package settings

import (
"context"
"time"

"connectrpc.com/connect"

"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
"github.com/zitadel/zitadel/internal/zerrors"
pb "github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

var zeroTime time.Time

func (s *Server) ListDetectionRules(ctx context.Context, _ *connect.Request[pb.ListDetectionRulesRequest]) (*connect.Response[pb.ListDetectionRulesResponse], error) {
persistedSettings, err := s.query.DetectionSettings(ctx)
if err != nil {
return nil, err
}
if persistedSettings == nil || !persistedSettings.RulesManaged {
defaults := s.command.DefaultDetectionRules()
rules := make([]*pb.DetectionRule, len(defaults))
for i, rule := range defaults {
rules[i] = detectionRuleToPb(rule, zeroTime, zeroTime)
}
return connect.NewResponse(&pb.ListDetectionRulesResponse{Rules: rules}), nil
}
rules, err := s.query.SearchDetectionRules(ctx)
if err != nil {
return nil, err
}
resp := make([]*pb.DetectionRule, len(rules))
for i, rule := range rules {
resp[i] = queryDetectionRuleToPb(rule)
}
return connect.NewResponse(&pb.ListDetectionRulesResponse{Rules: resp}), nil
}

func (s *Server) GetDetectionRule(ctx context.Context, req *connect.Request[pb.GetDetectionRuleRequest]) (*connect.Response[pb.GetDetectionRuleResponse], error) {
persistedSettings, err := s.query.DetectionSettings(ctx)
if err != nil {
return nil, err
}
if persistedSettings == nil || !persistedSettings.RulesManaged {
for _, rule := range s.command.DefaultDetectionRules() {
if rule.ID == req.Msg.GetRuleId() {
return connect.NewResponse(&pb.GetDetectionRuleResponse{Rule: detectionRuleToPb(rule, zeroTime, zeroTime)}), nil
}
}
return nil, zerrors.ThrowNotFound(nil, "SETT-uG7qP", "Errors.NotFound")
}
rule, err := s.query.DetectionRule(ctx, req.Msg.GetRuleId())
if err != nil {
return nil, err
}
if rule == nil {
return nil, zerrors.ThrowNotFound(nil, "SETT-Hx2w8", "Errors.NotFound")
}
return connect.NewResponse(&pb.GetDetectionRuleResponse{Rule: queryDetectionRuleToPb(rule)}), nil
}

func (s *Server) CreateDetectionRule(ctx context.Context, req *connect.Request[pb.CreateDetectionRuleRequest]) (*connect.Response[pb.CreateDetectionRuleResponse], error) {
rule, err := detectionRuleToDomain(req.Msg.GetRule())
if err != nil {
return nil, err
}
details, err := s.command.CreateDetectionRule(ctx, rule)
if err != nil {
return nil, err
}
return connect.NewResponse(&pb.CreateDetectionRuleResponse{Details: object.DomainToDetailsPb(details)}), nil
}

func (s *Server) UpdateDetectionRule(ctx context.Context, req *connect.Request[pb.UpdateDetectionRuleRequest]) (*connect.Response[pb.UpdateDetectionRuleResponse], error) {
rule, err := detectionRuleToDomain(req.Msg.GetRule())
if err != nil {
return nil, err
}
if rule.ID == "" {
rule.ID = req.Msg.GetRuleId()
}
if rule.ID != req.Msg.GetRuleId() {
return nil, zerrors.ThrowInvalidArgument(nil, "SETT-4A4wE", "Errors.Risk.Invalid")
}
details, err := s.command.ChangeDetectionRule(ctx, rule)
if err != nil {
return nil, err
}
return connect.NewResponse(&pb.UpdateDetectionRuleResponse{Details: object.DomainToDetailsPb(details)}), nil
}

func (s *Server) DeleteDetectionRule(ctx context.Context, req *connect.Request[pb.DeleteDetectionRuleRequest]) (*connect.Response[pb.DeleteDetectionRuleResponse], error) {
details, err := s.command.RemoveDetectionRule(ctx, req.Msg.GetRuleId())
if err != nil {
return nil, err
}
return connect.NewResponse(&pb.DeleteDetectionRuleResponse{Details: object.DomainToDetailsPb(details)}), nil
}
