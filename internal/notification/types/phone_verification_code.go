package types

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func (notify Notify) SendPhoneVerificationCode(user *query.NotifyUser, origin, code, requestedDomain string) error {
	args := make(map[string]interface{})
	args["Code"] = code
	args["Domain"] = requestedDomain
	return notify("", args, domain.VerifyPhoneMessageType, true)
}
