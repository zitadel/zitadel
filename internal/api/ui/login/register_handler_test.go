package login

import (
	"context"
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
)

func Test_determineResourceOwner(t *testing.T) {
	// Note: the nil-authRequest and "both empty" cases fall back to
	// authz.GetInstance(ctx).DefaultOrganisationID() and are covered by the
	// integration/handler tests where an instance is present in the context.
	// Here we lock the resolution order for the cases that do not touch the
	// instance: RequestedOrgID takes precedence over UserOrgID, and UserOrgID
	// is used as the fallback before the default org.
	tests := []struct {
		name        string
		authRequest *domain.AuthRequest
		want        string
	}{
		{
			name:        "requested org wins over user org",
			authRequest: &domain.AuthRequest{RequestedOrgID: "requested", UserOrgID: "user"},
			want:        "requested",
		},
		{
			name:        "requested org only",
			authRequest: &domain.AuthRequest{RequestedOrgID: "requested"},
			want:        "requested",
		},
		{
			name:        "user org used when no requested org (matches fillPolicies resolution)",
			authRequest: &domain.AuthRequest{UserOrgID: "user"},
			want:        "user",
		},
		{
			// domain discovery sets BOTH RequestedOrgID and UserOrgID to the
			// discovered org, so it still resolves to the discovered org.
			name:        "domain discovery result",
			authRequest: &domain.AuthRequest{RequestedOrgID: "discovered", UserOrgID: "discovered"},
			want:        "discovered",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := determineResourceOwner(context.Background(), tt.authRequest); got != tt.want {
				t.Errorf("determineResourceOwner() = %q, want %q", got, tt.want)
			}
		})
	}
}
