package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
)

type humanOperation struct {
	userOperation
}

// GetEmail implements domain.HumanOperation.
func (h *humanOperation) GetEmail(ctx context.Context) (*domain.Email, error) {
	var email domain.Email
	err := h.QueryExecutor.QueryRow(ctx, `SELECT email, is_email_verified FROM human_users WHERE id = $1`, h.clauses).Scan(
		&email.Address,
		&email.IsVerified,
	)
	if err != nil {
		return nil, err
	}
	return &email, nil
}

// SetEmail implements domain.HumanOperation.
func (h *humanOperation) SetEmail(ctx context.Context, email string) error {
	return h.QueryExecutor.Exec(ctx, `UPDATE human_users SET email = $1 WHERE id = $2`, email, h.clauses)
}

// SetEmailVerified implements domain.HumanOperation.
func (h *humanOperation) SetEmailVerified(ctx context.Context, email string) error {
	return h.QueryExecutor.Exec(ctx, `UPDATE human_users SET is_email_verified = $1 WHERE id = $2 AND email = $3`, true, h.clauses, email)
}

var _ domain.HumanOperation = (*humanOperation)(nil)
