package admin

import (
	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

func instanceLimitsPbToCommand(req *admin.SetInstanceLimitsRequest) *command.SetLimits {
	return &command.SetLimits{DisallowPublicOrgRegistration: gu.Ptr(req.GetDisallowPublicOrgRegistration())}
}
