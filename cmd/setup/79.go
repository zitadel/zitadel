package setup

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type BackfillUniqueConstraintOwners struct {
	dbClient *database.DB
}

func (mig *BackfillUniqueConstraintOwners) Execute(ctx context.Context, _ eventstore.Event) error {
	for i, stmt := range backfillUniqueConstraintOwnersStmts {
		logging.Info(ctx, "backfill unique constraint owners", "index", i, "migration", mig.String())
		if _, err := mig.dbClient.ExecContext(ctx, stmt); err != nil {
			if isUndefinedTable(err) {
				logging.Info(ctx, "skip unique constraint owners backfill, relation missing", "index", i, "migration", mig.String())
				continue
			}
			return err
		}
	}

	var unmatched int64
	err := mig.dbClient.QueryRowContext(ctx, func(row *sql.Row) error {
		return row.Scan(&unmatched)
	}, `SELECT COUNT(*) FROM eventstore.unique_constraints WHERE owners = '{}'`)
	if err != nil {
		return err
	}
	logging.Info(ctx, "unique constraint owners backfill complete", "unmatched", unmatched, "migration", mig.String())
	return nil
}

func (mig *BackfillUniqueConstraintOwners) String() string {
	return "79_backfill_unique_constraint_owners"
}

func isUndefinedTable(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "42P01" {
		return true
	}
	return strings.Contains(err.Error(), "does not exist")
}

var backfillUniqueConstraintOwnersStmts = []string{
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || u.resource_owner, 'user:' || u.id]
FROM projections.users14 u
WHERE uc.unique_type = 'usernames'
	AND uc.instance_id = u.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(u.username || u.resource_owner)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || u.resource_owner, 'user:' || u.id]
FROM projections.users14 u
WHERE uc.unique_type = 'usernames'
	AND uc.instance_id = u.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(u.username)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || l.resource_owner, 'idp:' || l.idp_id, 'user:' || l.user_id]
FROM projections.idp_user_links3 l
WHERE uc.unique_type = 'external_idps'
	AND uc.instance_id = l.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(l.idp_id || l.external_user_id)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || o.id]
FROM projections.orgs1 o
WHERE uc.unique_type = 'org_name'
	AND uc.instance_id = o.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(o.name)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || d.org_id]
FROM projections.org_domains2 d
WHERE uc.unique_type = 'org_domain'
	AND uc.instance_id = d.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(d.domain)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || p.resource_owner, 'project:' || p.id]
FROM projections.projects4 p
WHERE uc.unique_type = 'project_names'
	AND uc.instance_id = p.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(p.name || p.resource_owner)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || a.resource_owner, 'project:' || a.project_id]
FROM projections.apps7 a
WHERE uc.unique_type = 'appname'
	AND uc.instance_id = a.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(a.name || ':' || a.project_id)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || r.resource_owner, 'project:' || r.project_id]
FROM projections.project_roles4 r
WHERE uc.unique_type = 'project_role'
	AND uc.instance_id = r.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(r.role_key || ':' || r.project_id)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || a.resource_owner, 'project:' || a.project_id]
FROM projections.apps7 a
JOIN projections.apps7_saml_configs s ON s.app_id = a.id AND s.instance_id = a.instance_id
WHERE uc.unique_type = 'entity_ids'
	AND uc.instance_id = a.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(s.entity_id)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || g.resource_owner, 'org:' || g.granted_org_id, 'project:' || g.project_id]
FROM projections.project_grants4 g
WHERE uc.unique_type = 'project_grant'
	AND uc.instance_id = g.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(g.granted_org_id || ':' || g.project_id)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || m.resource_owner, 'user:' || m.user_id, 'project:' || m.project_id] || CASE WHEN COALESCE(m.grant_id, '') <> '' THEN ARRAY['grant:' || m.grant_id] ELSE ARRAY[]::text[] END
FROM projections.project_grant_members4 m
WHERE uc.unique_type = 'project_grant_member'
	AND uc.instance_id = m.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(m.project_id || ':' || m.user_id || ':' || m.grant_id)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || g.resource_owner, 'user:' || g.user_id, 'project:' || g.project_id] || CASE WHEN COALESCE(g.grant_id, '') <> '' THEN ARRAY['grant:' || g.grant_id] ELSE ARRAY[]::text[] END
FROM projections.user_grants5 g
WHERE uc.unique_type = 'user_grant'
	AND uc.instance_id = g.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(g.resource_owner || ':' || g.user_id || ':' || g.project_id || ':' || g.grant_id)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || m.org_id, 'user:' || m.user_id]
FROM projections.org_members4 m
WHERE uc.unique_type = 'member'
	AND uc.instance_id = m.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(m.org_id || ':' || m.user_id)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || m.resource_owner, 'user:' || m.user_id, 'project:' || m.project_id]
FROM projections.project_members4 m
WHERE uc.unique_type = 'member'
	AND uc.instance_id = m.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(m.project_id || ':' || m.user_id)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['user:' || m.user_id]
FROM projections.instance_members4 m
WHERE uc.unique_type = 'member'
	AND uc.instance_id = m.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(m.id || ':' || m.user_id)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || g.resource_owner]
FROM projections.groups1 g
WHERE uc.unique_type = 'group_name'
	AND uc.instance_id = g.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(g.name || ':' || g.resource_owner)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || a.resource_owner]
FROM projections.actions3 a
WHERE uc.unique_type = 'action_names'
	AND uc.instance_id = a.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(a.name || ':' || a.resource_owner)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || i.resource_owner, 'idp:' || i.id]
FROM projections.idps3 i
WHERE uc.unique_type = 'idp_config_names'
	AND uc.instance_id = i.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(i.name || i.resource_owner)`,
	`UPDATE eventstore.unique_constraints uc
SET owners = ARRAY['org:' || t.resource_owner, 'idp:' || t.id]
FROM projections.idp_templates6 t
WHERE uc.unique_type = 'idp_config_names'
	AND uc.instance_id = t.instance_id
	AND uc.owners = '{}'
	AND lower(uc.unique_field) = lower(t.name || t.resource_owner)`,
}
