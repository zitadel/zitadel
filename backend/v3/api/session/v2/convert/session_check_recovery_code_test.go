package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/domain"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestCheckRecoveryCodeGRPCToDomain(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		check *session_grpc.CheckRecoveryCode
		want  *domain.CheckTypeRecoveryCode
	}{
		{
			name: "nil recovery code check",
		},
		{
			name: "recovery code check",
			check: &session_grpc.CheckRecoveryCode{
				Code: "code1",
			},
			want: &domain.CheckTypeRecoveryCode{
				RecoveryCode: "code1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := CheckRecoveryCodeGRPCToDomain(tt.check)
			assert.Equal(t, tt.want, got)
		})
	}
}
