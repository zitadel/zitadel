//go:build integration

package setup_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/cmd/setup"
	"github.com/zitadel/zitadel/internal/query/projection"
)

func TestBackfillUniqueConstraintOwners_OrgDomainVerifiedOnly(t *testing.T) {
	instanceID := "uc-owners-domain-" + time.Now().Format("150405.000")
	defer cleanupInstance(t, instanceID)

	_, err := dbPool.Exec(CTX, `
		INSERT INTO eventstore.unique_constraints (instance_id, unique_type, unique_field, owners)
		VALUES ($1, 'org_domain', 'example.com', '{}')`, instanceID)
	require.NoError(t, err)

	now := time.Now()
	_, err = dbPool.Exec(CTX, `
		INSERT INTO `+projection.OrgDomainTable+` (
			creation_date, change_date, sequence, domain, org_id, instance_id, is_verified, is_primary, validation_type
		) VALUES
			($2, $2, 1, 'example.com', 'org-a', $1, true, true, 0),
			($2, $2, 1, 'example.com', 'org-b', $1, false, false, 0)`, instanceID, now)
	require.NoError(t, err)

	query, err := setup.BackfillUniqueConstraintOwnersQuery("05_org_domain.sql")
	require.NoError(t, err)
	_, err = dbPool.Exec(CTX, query, nil, 5000)
	require.NoError(t, err)

	var owners []string
	err = dbPool.QueryRow(CTX, `
		SELECT owners FROM eventstore.unique_constraints
		WHERE instance_id = $1 AND unique_type = 'org_domain' AND unique_field = 'example.com'`, instanceID).Scan(&owners)
	require.NoError(t, err)
	assert.Equal(t, []string{"org:org-a"}, owners)
}

func TestBackfillUniqueConstraintOwners_UsernameInstanceScopedOrg(t *testing.T) {
	instanceID := "uc-owners-user-" + time.Now().Format("150405.000")
	defer cleanupInstance(t, instanceID)

	now := time.Now()
	_, err := dbPool.Exec(CTX, `
		INSERT INTO eventstore.unique_constraints (instance_id, unique_type, unique_field, owners)
		VALUES
			($1, 'usernames', 'bob', '{}'),
			($1, 'usernames', 'boborg-y', '{}')`, instanceID)
	require.NoError(t, err)

	_, err = dbPool.Exec(CTX, `
		INSERT INTO `+projection.UserTable+` (
			id, creation_date, change_date, resource_owner, instance_id, state, sequence, username, type
		) VALUES
			('user-x', $2, $2, 'org-x', $1, 1, 1, 'bob', 1),
			('user-y', $2, $2, 'org-y', $1, 1, 1, 'bob', 1)`, instanceID, now)
	require.NoError(t, err)

	_, err = dbPool.Exec(CTX, `
		INSERT INTO `+projection.DomainPolicyTable+` (
			creation_date, change_date, sequence, id, state, user_login_must_be_domain,
			validate_org_domains, smtp_sender_address_matches_instance_domain, is_default, resource_owner, instance_id
		) VALUES ($2, $2, 1, $1, 1, false, false, false, true, $1, $1)`, instanceID, now)
	require.NoError(t, err)

	// org-y scopes usernames to the organization through its own domain policy.
	_, err = dbPool.Exec(CTX, `
		INSERT INTO `+projection.DomainPolicyTable+` (
			creation_date, change_date, sequence, id, state, user_login_must_be_domain,
			validate_org_domains, smtp_sender_address_matches_instance_domain, is_default, resource_owner, instance_id
		) VALUES ($2, $2, 1, 'org-y', 1, true, false, false, false, 'org-y', $1)`, instanceID, now)
	require.NoError(t, err)

	query, err := setup.BackfillUniqueConstraintOwnersQuery("02_usernames.sql")
	require.NoError(t, err)
	_, err = dbPool.Exec(CTX, query, nil, 5000)
	require.NoError(t, err)

	var instanceOwners []string
	err = dbPool.QueryRow(CTX, `
		SELECT owners FROM eventstore.unique_constraints
		WHERE instance_id = $1 AND unique_type = 'usernames' AND unique_field = 'bob'`, instanceID).Scan(&instanceOwners)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"org:org-x", "user:user-x"}, instanceOwners)

	var orgScopedOwners []string
	err = dbPool.QueryRow(CTX, `
		SELECT owners FROM eventstore.unique_constraints
		WHERE instance_id = $1 AND unique_type = 'usernames' AND unique_field = 'boborg-y'`, instanceID).Scan(&orgScopedOwners)
	require.NoError(t, err)
	assert.Empty(t, orgScopedOwners)
}

func cleanupInstance(t *testing.T, instanceID string) {
	t.Helper()
	_, _ = dbPool.Exec(CTX, `DELETE FROM eventstore.unique_constraints WHERE instance_id = $1`, instanceID)
	_, _ = dbPool.Exec(CTX, `DELETE FROM `+projection.OrgDomainTable+` WHERE instance_id = $1`, instanceID)
	_, _ = dbPool.Exec(CTX, `DELETE FROM `+projection.UserTable+` WHERE instance_id = $1`, instanceID)
	_, _ = dbPool.Exec(CTX, `DELETE FROM `+projection.DomainPolicyTable+` WHERE instance_id = $1`, instanceID)
}
