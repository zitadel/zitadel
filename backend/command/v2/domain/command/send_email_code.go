package command

import (
	"context"

	"github.com/zitadel/zitadel/backend/command/v2/pattern"
)

var _ pattern.Command = (*sendEmailCode)(nil)

type sendEmailCode struct {
	UserID      string  `json:"userId"`
	Email       string  `json:"email"`
	URLTemplate *string `json:"urlTemplate"`
	code        string  `json:"-"`
}

func SendEmailCode(userID, email string, urlTemplate *string) pattern.Command {
	cmd := &sendEmailCode{
		UserID:      userID,
		Email:       email,
		URLTemplate: urlTemplate,
	}

	return pattern.Batch(GenerateCode(cmd.SetCode, generateCode))
}

// Name implements [pattern.Command].
func (c *sendEmailCode) Name() string {
	return "user.v2.email.send_code"
}

// Execute implements [pattern.Command].
func (c *sendEmailCode) Execute(ctx context.Context) error {
	// Implementation of the command execution
	return nil
}

func (c *sendEmailCode) SetCode(code string) {
	c.code = code
}
