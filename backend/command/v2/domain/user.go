package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/command/v2/domain/command"
	"github.com/zitadel/zitadel/backend/command/v2/domain/query"
	"github.com/zitadel/zitadel/backend/command/v2/pattern"
	"github.com/zitadel/zitadel/backend/command/v2/storage/database"
	"github.com/zitadel/zitadel/backend/command/v2/telemetry/tracing"
)

type User struct {
	ID       string
	Username string
	Email    Email
}

type SetUserEmail struct {
	UserID string
	Email  string

	IsVerified *bool
	ReturnCode *ReturnCode
	SendCode   *SendCode

	code   string
	client database.QueryExecutor
}

func (e *SetUserEmail) SetCode(code string) {
	e.code = code
}

type ReturnCode struct {
	// Code is the code to be sent to the user
	Code string
}

type SendCode struct {
	// URLTemplate is the template for the URL that is rendered into the message
	URLTemplate *string
}

func (d *Domain) SetUserEmail(ctx context.Context, req *SetUserEmail) error {
	batch := pattern.Batch(
		tracing.Trace(d.tracer, command.SetEmail(req.UserID, req.Email)),
	)

	if req.IsVerified == nil {
		batch.Append(command.GenerateCode(
			req.SetCode,
			query.QueryEncryptionGenerator(
				database.Query(d.pool),
				d.userCodeAlg,
			),
		))
	} else {
		batch.Append(command.VerifyEmail(req.UserID, req.Email))
	}

	// if !req.GetVerification().GetIsVerified() {
	// 	batch.

	// switch req.GetVerification().(type) {
	// case *user.SetEmailRequest_IsVerified:
	// 	batch.Append(tracing.Trace(s.tracer, command.VerifyEmail(req.GetUserId(), req.GetEmail())))
	// case *user.SetEmailRequest_SendCode:
	// 	batch.Append(tracing.Trace(s.tracer, command.SendEmailCode(req.GetUserId(), req.GetEmail(), req.GetSendCode().UrlTemplate)))
	// case *user.SetEmailRequest_ReturnCode:
	// 	batch.Append(tracing.Trace(s.tracer, query.ReturnEmailCode(req.GetUserId(), req.GetEmail())))
	// }

	// if err := batch.Execute(ctx); err != nil {
	// 	return nil, err
	// }
}
