package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 68.sql
	targetAddPayloadTypeAndKeyFingerprintColumn string
)

type TargetAddPayloadTypeColumn struct {
	dbClient *database.DB
}

func (mig *TargetAddPayloadTypeColumn) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, targetAddPayloadTypeAndKeyFingerprintColumn)
	return err
}

func (mig *TargetAddPayloadTypeColumn) String() string {
	return "68_target2_add_payload_type_authn_key2_fingerprint_expiration"
}
