package types

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

func (notify Notify) SendOTPSMSCode(user *query.NotifyUser, requestedDomain, origin, code string, expiry time.Duration) error {
	var url string
	url = "test"
	//if urlTmpl == "" {
	//	url = login.MailVerificationLink(origin, user.ID, code, user.ResourceOwner)
	//} else {
	//	var buf strings.Builder
	//	if err := domain.RenderConfirmURLTemplate(&buf, urlTmpl, user.ID, code, user.ResourceOwner); err != nil {
	//		return err
	//	}
	//	url = buf.String()
	//}

	args := make(map[string]interface{})
	args["OTP"] = code
	args["URL"] = url
	args["Domain"] = requestedDomain
	args["Expiry"] = expiry
	return notify(url, args, domain.VerifySMSOTPMessageType, false)
}
