package settings

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/detection"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	pb "github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

func detectionSettingsToPb(current command.DetectionSettings) *pb.DetectionSettings {
	return &pb.DetectionSettings{
		Enabled:               current.Enabled,
		FailOpen:              current.FailOpen,
		FailureBurstThreshold: uint32(current.FailureBurstThreshold),
		HistoryWindow:         durationpb.New(current.HistoryWindow),
		ContextChangeWindow:   durationpb.New(current.ContextChangeWindow),
		MaxSignalsPerUser:     uint32(current.MaxSignalsPerUser),
		MaxSignalsPerSession:  uint32(current.MaxSignalsPerSession),
	}
}

func detectionSettingsToCommand(current *pb.DetectionSettings) *command.DetectionSettings {
	return &command.DetectionSettings{
		Enabled:               current.GetEnabled(),
		FailOpen:              current.GetFailOpen(),
		FailureBurstThreshold: int(current.GetFailureBurstThreshold()),
		HistoryWindow:         durationOrZero(current.GetHistoryWindow()),
		ContextChangeWindow:   durationOrZero(current.GetContextChangeWindow()),
		MaxSignalsPerUser:     int(current.GetMaxSignalsPerUser()),
		MaxSignalsPerSession:  int(current.GetMaxSignalsPerSession()),
	}
}

func detectionRuleToPb(rule detection.Rule, creationDate, changeDate time.Time) *pb.DetectionRule {
	resp := &pb.DetectionRule{
		Id:          rule.ID,
		Description: rule.Description,
		Expr:        rule.Expr,
		Engine:      detectionRuleActionToPb(rule.Action),
		Finding: &pb.DetectionRuleFinding{
			Name:    rule.FindingCfg.Name,
			Message: rule.FindingCfg.Message,
			Block:   rule.FindingCfg.Block,
		},
		ContextTemplate: rule.ContextTemplate,
	}
	if rule.Action == detection.ActionRateLimit {
		resp.RateLimit = &pb.DetectionRuleRateLimit{
			Key:    rule.RateLimitCfg.KeyTemplate,
			Window: durationpb.New(rule.RateLimitCfg.Window),
			Max:    uint32(rule.RateLimitCfg.Max),
		}
	}
	if !creationDate.IsZero() {
		resp.CreationDate = timestamppb.New(creationDate)
	}
	if !changeDate.IsZero() {
		resp.ChangeDate = timestamppb.New(changeDate)
	}
	return resp
}

func queryDetectionRuleToPb(rule *query.DetectionRule) *pb.DetectionRule {
	if rule == nil {
		return nil
	}
	return detectionRuleToPb(detection.Rule{
		ID:          rule.ID,
		Description: rule.Description,
		Expr:        rule.Expr,
		Action:      rule.Action,
		Priority:    int(rule.Priority),
		StopOnMatch: rule.StopOnMatch,
		FindingCfg: detection.RuleFinding{
			Name:    rule.FindingName,
			Message: rule.FindingMessage,
			Block:   rule.FindingBlock,
		},
		ContextTemplate: rule.ContextTemplate,
		RateLimitCfg: detection.RuleRateLimit{
			KeyTemplate: rule.RateLimitKey,
			Window:      time.Duration(rule.RateLimitWindow),
			Max:         int(rule.RateLimitMax),
		},
	}, rule.CreationDate, rule.ChangeDate)
}

func detectionRuleToDomain(rule *pb.DetectionRule) (detection.Rule, error) {
	action, err := detectionRuleActionToDomain(rule.GetEngine())
	if err != nil {
		return detection.Rule{}, err
	}
	domainRule := detection.Rule{
		ID:              rule.GetId(),
		Description:     rule.GetDescription(),
		Expr:            rule.GetExpr(),
		Action:          action,
		ContextTemplate: rule.GetContextTemplate(),
	}
	if finding := rule.GetFinding(); finding != nil {
		domainRule.FindingCfg = detection.RuleFinding{
			Name:    finding.GetName(),
			Message: finding.GetMessage(),
			Block:   finding.GetBlock(),
		}
	}
	if rateLimit := rule.GetRateLimit(); rateLimit != nil {
		domainRule.RateLimitCfg = detection.RuleRateLimit{
			KeyTemplate: rateLimit.GetKey(),
			Window:      durationOrZero(rateLimit.GetWindow()),
			Max:         int(rateLimit.GetMax()),
		}
	}
	return domainRule, nil
}

func detectionRuleActionToPb(action detection.ActionType) pb.DetectionRuleEngine {
	switch action {
	case detection.ActionBlock:
		return pb.DetectionRuleEngine_DETECTION_RULE_ENGINE_BLOCK
	case detection.ActionRateLimit:
		return pb.DetectionRuleEngine_DETECTION_RULE_ENGINE_RATE_LIMIT
	case detection.ActionLLM:
		return pb.DetectionRuleEngine_DETECTION_RULE_ENGINE_LLM
	case detection.ActionLog:
		return pb.DetectionRuleEngine_DETECTION_RULE_ENGINE_LOG
	case detection.ActionCaptcha:
		return pb.DetectionRuleEngine_DETECTION_RULE_ENGINE_CAPTCHA
	default:
		return pb.DetectionRuleEngine_DETECTION_RULE_ENGINE_UNSPECIFIED
	}
}

func detectionRuleActionToDomain(action pb.DetectionRuleEngine) (detection.ActionType, error) {
	switch action {
	case pb.DetectionRuleEngine_DETECTION_RULE_ENGINE_BLOCK:
		return detection.ActionBlock, nil
	case pb.DetectionRuleEngine_DETECTION_RULE_ENGINE_RATE_LIMIT:
		return detection.ActionRateLimit, nil
	case pb.DetectionRuleEngine_DETECTION_RULE_ENGINE_LLM:
		return detection.ActionLLM, nil
	case pb.DetectionRuleEngine_DETECTION_RULE_ENGINE_LOG:
		return detection.ActionLog, nil
	case pb.DetectionRuleEngine_DETECTION_RULE_ENGINE_CAPTCHA:
		return detection.ActionCaptcha, nil
	default:
		return "", zerrors.ThrowInvalidArgument(nil, "SETT-sC7k1", "Errors.Risk.Invalid")
	}
}

func durationOrZero(d *durationpb.Duration) time.Duration {
	if d == nil {
		return 0
	}
	return d.AsDuration()
}
