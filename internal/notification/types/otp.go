package types

import (
	"context"
	"time"

	http_utils "github.com/zitadel/zitadel/internal/api/http"

	"github.com/zitadel/zitadel/internal/api/authz"

	"github.com/zitadel/zitadel/internal/domain"
)

func (notify Notify) SendOTPSMSCode(ctx context.Context, code string, expiry time.Duration) error {
	args := otpArgs(ctx, code, expiry)
	return notify("", args, domain.VerifySMSOTPMessageType, false)
}

func (notify Notify) SendOTPEmailCode(ctx context.Context, url, code string, expiry time.Duration) error {
	args := otpArgs(ctx, code, expiry)
	return notify(url, args, domain.VerifyEmailOTPMessageType, false)
}

func otpArgs(ctx context.Context, code string, expiry time.Duration) map[string]interface{} {
	args := make(map[string]interface{})
	args["OTP"] = code
	args["Origin"] = http_utils.ComposedOrigin(ctx)
	args["Domain"] = authz.GetInstance(ctx).RequestedDomain()
	args["Expiry"] = expiry
	return args
}
