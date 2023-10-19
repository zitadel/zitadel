package system

import (
	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func instanceLimitsPbToCommand(req *system.SetLimitsRequest) *command.SetLimits {
	var setLimits = new(command.SetLimits)
	if req.AuditLogRetention != nil {
		setLimits.AuditLogRetention = gu.Ptr(req.AuditLogRetention.AsDuration())
	}
	return setLimits
}
