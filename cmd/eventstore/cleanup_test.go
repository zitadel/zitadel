package eventstore

import (
	"strings"
	"testing"
)

func TestCleanupSQLTargetsOnlySafeTerminalAggregates(t *testing.T) {
	for _, sqlText := range []string{cleanupEvents2CountSQL, cleanupEvents2DeleteSQL} {
		for _, forbidden := range []string{
			"'session'",
			"'session_logout'",
			"'auth_request.code.exchanged'",
		} {
			if strings.Contains(sqlText, forbidden) {
				t.Fatalf("cleanup SQL must not target %s", forbidden)
			}
		}

		for _, required := range []string{
			"'auth_request.failed'",
			"'auth_request.succeeded'",
			"'device.authorization.canceled'",
			"'device.authorization.done'",
			"'idpintent.consumed'",
			"'idpintent.failed'",
			"'oidc_session'",
			"'oidc_session.access_token.added'",
			"'oidc_session.refresh_token.added'",
			"'oidc_session.refresh_token.renewed'",
			"'oidc_session.refresh_token.revoked'",
			"'saml_request.failed'",
			"'saml_request.succeeded'",
		} {
			if !strings.Contains(sqlText, required) {
				t.Fatalf("cleanup SQL must include %s", required)
			}
		}
	}
}
