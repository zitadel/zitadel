package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func Test_authzInstance_checkDomain(t *testing.T) {
	instance := &authzInstance{
		TrustedDomains: []string{"trusted.domain"},
	}
	tests := []struct {
		name    string
		domain  *http_util.DomainCtx
		wantErr bool
	}{
		{
			name:   "no public host",
			domain: &http_util.DomainCtx{InstanceHost: "instance.domain"},
		},
		{
			name:   "public host equals instance host",
			domain: &http_util.DomainCtx{InstanceHost: "instance.domain", PublicHost: "instance.domain"},
		},
		{
			name:   "trusted public host",
			domain: &http_util.DomainCtx{InstanceHost: "instance.domain", PublicHost: "trusted.domain"},
		},
		{
			name:   "trusted public host with port",
			domain: &http_util.DomainCtx{InstanceHost: "instance.domain", PublicHost: "trusted.domain:443"},
		},
		{
			name:    "untrusted public host",
			domain:  &http_util.DomainCtx{InstanceHost: "instance.domain", PublicHost: "evil.com"},
			wantErr: true,
		},
		{
			name:    "invalid public host (non-numeric port)",
			domain:  &http_util.DomainCtx{InstanceHost: "instance.domain", PublicHost: "evil.com:443/"},
			wantErr: true,
		},
		{
			name:    "invalid public host (userinfo)",
			domain:  &http_util.DomainCtx{InstanceHost: "instance.domain", PublicHost: "a:b@evil.com"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := instance.checkDomain(tt.domain.InstanceDomain(), tt.domain.RequestedDomain())
			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, zerrors.IsNotFound(err))
				return
			}
			assert.NoError(t, err)
		})
	}
}
