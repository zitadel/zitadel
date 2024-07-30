package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 32.sql
	addTwilioVerifyServiceSID string
)

type SMSConfigs2TwilioAddVerifyServiceSid struct {
	dbClient *database.DB
}

func (mig *SMSConfigs2TwilioAddVerifyServiceSid) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addTwilioVerifyServiceSID)
	return err
}

func (mig *SMSConfigs2TwilioAddVerifyServiceSid) String() string {
	return "32_sms_configs2_twilio_add_verification_sid"
}
