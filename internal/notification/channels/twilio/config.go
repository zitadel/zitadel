package twilio

import (
	newTwilio "github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type Config struct {
	SID              string
	Token            string
	SenderNumber     string
	VerifyServiceSID string
}

func (t *Config) IsValid() bool {
	return t.SID != "" && t.Token != "" && t.SenderNumber != ""
}

func (t *Config) VerifyCode(verificationID, code string) error {
	client := newTwilio.NewRestClientWithParams(newTwilio.ClientParams{Username: t.SID, Password: t.Token})
	checkParams := &verify.CreateVerificationCheckParams{}
	checkParams.SetVerificationSid(verificationID)
	checkParams.SetCode(code)
	resp, err := client.VerifyV2.CreateVerificationCheck(t.VerifyServiceSID, checkParams)
	if err != nil || resp.Status == nil {
		return zerrors.ThrowInvalidArgument(err, "TWILI-JK3ta", "Errors.User.Code.NotFound")
	}
	switch *resp.Status {
	case "approved":
		return nil
	case "expired":
		return zerrors.ThrowInvalidArgument(nil, "TWILI-SF3ba", "Errors.User.Code.Expired")
	case "max_attempts_reached":
		return zerrors.ThrowInvalidArgument(nil, "TWILI-Ok39a", "Errors.User.Code.NotFound")
	default:
		return zerrors.ThrowInvalidArgument(nil, "TWILI-Skwe4", "Errors.User.Code.Invalid")
	}
}
