package command

import (
	"context"

	"github.com/zitadel/zitadel/backend/command/v2/pattern"
)

var _ pattern.Command = (*verifyEmail)(nil)

type verifyEmail struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
}

func VerifyEmail(userID, email string) *verifyEmail {
	return &verifyEmail{
		UserID: userID,
		Email:  email,
	}
}

// Name implements [pattern.Command].
func (c *verifyEmail) Name() string {
	return "user.v2.verify_email"
}

// Execute implements [pattern.Command].
func (c *verifyEmail) Execute(ctx context.Context) error {
	// Implementation of the command execution
	return nil
}
