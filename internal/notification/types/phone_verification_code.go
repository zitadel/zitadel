package types

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
)

func (notify Notify) SendPhoneVerificationCode(ctx context.Context, code string) error {
	args := make(map[string]interface{})
	args["Code"] = code
	args["Domain"] = authz.GetInstance(ctx).RequestedDomain()
	return notify("", args, domain.VerifyPhoneMessageType, true)
}
