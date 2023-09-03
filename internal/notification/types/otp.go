package types

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func (notify Notify) SendOTPSMSCode(requestedDomain, origin, code string, expiry time.Duration) error {
	args := otpArgs(code, origin, requestedDomain, expiry)
	return notify("", args, domain.VerifySMSOTPMessageType, false)
}

func (notify Notify) SendOTPEmailCode(user *query.NotifyUser, url, requestedDomain, origin, code string, expiry time.Duration) error {
	args := otpArgs(code, origin, requestedDomain, expiry)
	return notify(url, args, domain.VerifyEmailOTPMessageType, false)
}

func otpArgs(code, origin, requestedDomain string, expiry time.Duration) map[string]interface{} {
	args := make(map[string]interface{})
	args["OTP"] = code
	args["Origin"] = origin
	args["Domain"] = requestedDomain
	args["Expiry"] = expiry
	return args
}
