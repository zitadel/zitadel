package setup

import (
	"context"
	"embed"
	"fmt"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type CleanupOrgDomainRemoved struct {
	eventstoreClient *database.DB
}

var (
	//go:embed 47/*.sql
	cleanupOrgDomainRemoved embed.FS
)

func (mig *CleanupOrgDomainRemoved) Execute(ctx context.Context, _ eventstore.Event) error {
	statements, err := readStatements(cleanupOrgDomainRemoved, "47", "")
	if err != nil {
		return err
	}
	for _, stmt := range statements {
		logging.WithFields("file", stmt.file, "migration", mig.String()).Info("execute statement")
		if _, err := mig.eventstoreClient.ExecContext(ctx, stmt.query); err != nil {
			return fmt.Errorf("%s %s: %w", mig.String(), stmt.file, err)
		}
	}
	return nil
}

func (*CleanupOrgDomainRemoved) String() string {
	return "47_cleanup_org_domain_removed"
}
