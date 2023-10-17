package system

import (
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func instanceLimitsPbToCommand(req *system.SetLimitsRequest) *command.SetLimits {
	return &command.SetLimits{
		AuditLogRetention: req.GetAuditLogRetention().AsDuration(),
	}
}
