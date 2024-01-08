package system

import (
	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func setInstanceLimitsPbToCommand(req *system.SetLimitsRequest) *command.SetLimits {
	var setLimits = new(command.SetLimits)
	if req.AuditLogRetention != nil {
		setLimits.AuditLogRetention = gu.Ptr(req.AuditLogRetention.AsDuration())
	}
	if block := req.GetBlock(); block != nil {
		setLimits.Block = gu.Ptr(block.GetValue())
	}
	return setLimits
}

func bulkSetInstanceLimitsPbToCommand(req *system.BulkSetLimitsRequest) []*command.SetLimitsBulk {
	cmds := make([]*command.SetLimitsBulk, len(req.Limits))
	for i := range req.Limits {
		setLimitsReq := req.Limits[i]
		cmds[i] = &command.SetLimitsBulk{
			InstanceID:    setLimitsReq.GetInstanceId(),
			ResourceOwner: setLimitsReq.GetInstanceId(),
			SetLimits:     *setInstanceLimitsPbToCommand(req.Limits[i]),
		}
	}
	return cmds
}
