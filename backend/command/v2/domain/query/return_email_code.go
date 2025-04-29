package query

import (
	"context"

	"github.com/zitadel/zitadel/backend/command/v2/pattern"
)

var _ pattern.Query[string] = (*returnEmailCode)(nil)

type returnEmailCode struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	code   string `json:"-"`
}

func ReturnEmailCode(userID, email string) *returnEmailCode {
	return &returnEmailCode{
		UserID: userID,
		Email:  email,
	}
}

// Name implements [pattern.Command].
func (c *returnEmailCode) Name() string {
	return "user.v2.email.return_code"
}

// Execute implements [pattern.Command].
func (c *returnEmailCode) Execute(ctx context.Context) error {
	// Implementation of the command execution
	return nil
}

// Result implements [pattern.Query].
func (c *returnEmailCode) Result() string {
	return c.code
}
