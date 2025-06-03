//go:build integration

package setup_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test58(t *testing.T) {
	const query = `INSERT INTO projections.hosted_login_translations
	(instance_id, aggregate_id, aggregate_type, creation_date, change_date, sequence, locale) VALUES
	($1, $2, $3, $4, $5, $6, $7)`
	now := time.Now().UTC()

	type params struct {
		instanceID    string
		aggregateID   string
		aggregateType string
		creationDate  time.Time
		changeDate    *time.Time
		sequence      uint64
		locale        string
	}

	tt := []struct {
		testName     string
		input        params
		expectsError bool
		errMsg       string
	}{
		{
			testName:     "when aggregate_type is neither 'instance' nor 'org' should return error",
			input:        params{instanceID: "instID", aggregateID: "instID", aggregateType: "organization", creationDate: now, changeDate: nil, sequence: 1, locale: "en"},
			expectsError: true,
			errMsg:       `violates check constraint "hosted_login_translations_aggregate_type_check"`,
		},
		{
			testName:     "when locale is less than 2 chars should return error",
			input:        params{instanceID: "instID", aggregateID: "instID", aggregateType: "instance", creationDate: now, changeDate: nil, sequence: 1, locale: " e "},
			expectsError: true,
			errMsg:       `violates check constraint "hosted_login_translations_locale_check"`,
		},
		{
			testName:     "when record is valid should not return error",
			input:        params{instanceID: "instID", aggregateID: "instID", aggregateType: "org", creationDate: now, changeDate: nil, sequence: 1, locale: " en "},
			expectsError: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			tx, err := dbPool.Begin(CTX)
			require.Nil(t, err)

			t.Cleanup(func() {
				tx.Rollback(CTX)
			})

			_, err = tx.Exec(CTX, query, tc.input.instanceID, tc.input.aggregateID, tc.input.aggregateType, tc.input.creationDate, tc.input.changeDate, tc.input.sequence, tc.input.locale)

			require.Equal(t, tc.expectsError, err != nil, err)
			if tc.expectsError {
				assert.Contains(t, err.Error(), tc.errMsg)
			}
		})
	}
}
