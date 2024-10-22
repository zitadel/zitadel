package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 33.sql
	addTwilioVerifyServiceSID string
)

type SMSConfigs3TwilioAddVerifyServiceSid struct {
	dbClient *database.DB
}

func (mig *SMSConfigs3TwilioAddVerifyServiceSid) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addTwilioVerifyServiceSID)
	return err
}

func (mig *SMSConfigs3TwilioAddVerifyServiceSid) String() string {
	return "33_sms_configs3_twilio_add_verification_sid"
}
