package domain

// TODO(IAM-Marco): Implement when recovery codes repository is available

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

type RecoveryCodeCheck struct {
	instanceID string
	sessionID  string

	RecoveryCode *session_grpc.CheckRecoveryCode
}

func NewRecoveryCodeCheckCommand(sessionID, instanceID string, request *session_grpc.CheckRecoveryCode) *RecoveryCodeCheck {
	return &RecoveryCodeCheck{
		instanceID:   instanceID,
		sessionID:    sessionID,
		RecoveryCode: request,
	}
}

// RequiresTransaction implements [Transactional].
func (r *RecoveryCodeCheck) RequiresTransaction() {}

// Events implements [Commander].
func (r *RecoveryCodeCheck) Events(ctx context.Context, opts *InvokeOpts) ([]eventstore.Command, error) {
	return nil, nil
}

// Execute implements [Commander].
func (r *RecoveryCodeCheck) Execute(ctx context.Context, opts *InvokeOpts) (err error) {
	return nil
}

// String implements [Commander].
func (r *RecoveryCodeCheck) String() string {
	return "RecoveryCodeCheck"
}

// Validate implements [Commander].
func (r *RecoveryCodeCheck) Validate(ctx context.Context, opts *InvokeOpts) (err error) {
	return nil
}

var _ Commander = (*RecoveryCodeCheck)(nil)
var _ Transactional = (*RecoveryCodeCheck)(nil)
