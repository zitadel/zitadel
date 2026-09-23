package setup

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query/projection"
)

var (
	//go:embed 79/*.sql
	backfillUniqueConstraintOwnersFS embed.FS
)

type BackfillUniqueConstraintOwners struct {
	dbClient *database.DB

	Version       string `json:"version"`
	Finalized     bool   `json:"finalized"`
	ForceFinalize bool   `json:"-"` // YAML/env only; not persisted on lastRun
	BatchSize     int    `json:"-"` // YAML/env only; not persisted on lastRun

	lastVersion   string `json:"-"`
	lastFinalized bool   `json:"-"`
}

const defaultUniqueConstraintOwnerBatchSize = 5000

func (mig *BackfillUniqueConstraintOwners) batchSize() int {
	if mig.BatchSize > 0 {
		return mig.BatchSize
	}
	return defaultUniqueConstraintOwnerBatchSize
}

func (mig *BackfillUniqueConstraintOwners) Check(lastRun map[string]interface{}) bool {
	if lastRun == nil {
		lastRun = map[string]interface{}{}
	}
	mig.lastVersion, _ = lastRun["version"].(string)
	mig.lastFinalized, _ = lastRun["finalized"].(bool)

	versionChanged := mig.lastVersion != mig.Version
	if mig.lastVersion == "" {
		return true
	}
	return versionChanged || (mig.ForceFinalize && !mig.lastFinalized)
}

func (mig *BackfillUniqueConstraintOwners) finalizeAfterSuccess() {
	mig.Finalized = mig.lastFinalized || mig.ForceFinalize || (mig.lastVersion != "" && mig.lastVersion != mig.Version)
}

func (mig *BackfillUniqueConstraintOwners) Execute(ctx context.Context, _ eventstore.Event) error {
	statements, err := backfillUniqueConstraintOwnersStatements()
	if err != nil {
		return err
	}
	batchSize := mig.batchSize()
	for _, stmt := range statements {
		var cursor database.TextArray[string]
		for {
			var next database.TextArray[string]
			var updated int64
			pageCursor := []string(cursor)
			if pageCursor == nil {
				pageCursor = []string{}
			}
			cursorArg, err := json.Marshal(pageCursor)
			if err != nil {
				return fmt.Errorf("backfill unique constraint owners %s: %w", stmt.file, err)
			}
			err = mig.dbClient.QueryRowContext(ctx, func(row *sql.Row) error {
				return row.Scan(&next, &updated)
			}, stmt.query, string(cursorArg), batchSize)
			if err != nil {
				if isUndefinedTable(err) {
					logging.Info(ctx, "skip unique constraint owners backfill, relation missing", "file", stmt.file, "migration", mig.String())
					break
				}
				return fmt.Errorf("backfill unique constraint owners %s: %w", stmt.file, err)
			}
			if len(next) == 0 {
				break
			}
			if sameTextArray(cursor, next) {
				return fmt.Errorf("backfill unique constraint owners %s: cursor did not advance", stmt.file)
			}
			logging.Info(ctx, "backfill unique constraint owners page", "file", stmt.file, "updated", updated, "cursor", []string(next), "migration", mig.String())
			cursor = next
		}
	}

	var unmatched int64
	err = mig.dbClient.QueryRowContext(ctx, func(row *sql.Row) error {
		return row.Scan(&unmatched)
	}, `SELECT COUNT(*) FROM eventstore.unique_constraints WHERE unique_type = ANY($1) AND owners = '{}'`,
		database.TextArray[string](eventstore.UniqueTypesWithOwners))
	if err != nil {
		return err
	}
	mig.finalizeAfterSuccess()
	logging.Info(ctx, "unique constraint owners backfill complete", "unmatched", unmatched, "finalized", mig.Finalized, "migration", mig.String())
	return nil
}

func (mig *BackfillUniqueConstraintOwners) String() string {
	return eventstore.UniqueConstraintOwnersBackfillStep
}

func sameTextArray(a, b database.TextArray[string]) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func isUndefinedTable(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "42P01"
}

func interpolateBackfillSQL(query string) string {
	return strings.NewReplacer(
		"{{.users}}", projection.UserTable,
		"{{.organization_settings}}", projection.OrganizationSettingsTable,
		"{{.domain_policies}}", projection.DomainPolicyTable,
		"{{.idp_user_links}}", projection.IDPUserLinkTable,
		"{{.orgs}}", projection.OrgProjectionTable,
		"{{.org_domains}}", projection.OrgDomainTable,
		"{{.projects}}", projection.ProjectProjectionTable,
		"{{.apps}}", projection.AppProjectionTable,
		"{{.apps_saml}}", projection.AppSAMLTable,
		"{{.project_roles}}", projection.ProjectRoleProjectionTable,
		"{{.project_grants}}", projection.ProjectGrantProjectionTable,
		"{{.project_grant_members}}", projection.ProjectGrantMemberProjectionTable,
		"{{.user_grants}}", projection.UserGrantProjectionTable,
		"{{.org_members}}", projection.OrgMemberProjectionTable,
		"{{.project_members}}", projection.ProjectMemberProjectionTable,
		"{{.instance_members}}", projection.InstanceMemberProjectionTable,
		"{{.groups}}", projection.GroupProjectionTable,
		"{{.actions}}", projection.ActionTable,
		"{{.idps}}", projection.IDPTable,
		"{{.idp_templates}}", projection.IDPTemplateTable,
		"{{.message_texts}}", projection.MessageTextTable,
	).Replace(query)
}

func backfillUniqueConstraintOwnersStatements() ([]statement, error) {
	statements, err := readStatements(backfillUniqueConstraintOwnersFS, "79")
	if err != nil {
		return nil, err
	}
	for i := range statements {
		statements[i].query = interpolateBackfillSQL(statements[i].query)
	}
	return statements, nil
}

// BackfillUniqueConstraintOwnersQuery returns the interpolated SQL for a setup 79 file.
// The statement reads one projection page: $1 is a JSON array primary-key cursor ("[]" starts at the beginning) and $2 is the page size.
// It returns the next cursor and the number of unique-constraint rows updated.
func BackfillUniqueConstraintOwnersQuery(file string) (string, error) {
	statements, err := backfillUniqueConstraintOwnersStatements()
	if err != nil {
		return "", err
	}
	for _, stmt := range statements {
		if stmt.file == file {
			return stmt.query, nil
		}
	}
	return "", errors.New("unknown backfill statement: " + file)
}
