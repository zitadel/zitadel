package system

import (
	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/v2/internal/command"
	"github.com/zitadel/zitadel/v2/pkg/grpc/system"
)

func instanceLimitsPbToCommand(req *system.SetLimitsRequest) *command.SetLimits {
	var setLimits = new(command.SetLimits)
	if req.AuditLogRetention != nil {
		setLimits.AuditLogRetention = gu.Ptr(req.AuditLogRetention.AsDuration())
	}
	return setLimits
}
