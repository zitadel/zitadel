package setup

import (
	"context"
	_ "embed"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed 81.sql
	addLabelPolicyLinkPreview string
)

type LabelPolicies3AddLinkPreview struct {
	dbClient *database.DB
}

func (mig *LabelPolicies3AddLinkPreview) Execute(ctx context.Context, _ eventstore.Event) error {
	_, err := mig.dbClient.ExecContext(ctx, addLabelPolicyLinkPreview)
	return err
}

func (mig *LabelPolicies3AddLinkPreview) String() string {
	return "81_label_policies3_add_link_preview"
}
