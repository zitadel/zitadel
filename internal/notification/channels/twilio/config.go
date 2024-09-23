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
	if err != nil || resp.Status == nil || *resp.Status != "approved" {
		return zerrors.ThrowInvalidArgument(err, "TWILI-JK3ta", "Errors.Code.Invalid")
	}
	return nil
}
