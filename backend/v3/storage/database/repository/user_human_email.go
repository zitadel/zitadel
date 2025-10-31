package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userHumanEmail struct {
	verification
}

func (email userHumanEmail) unqualifiedTableName() string {
	return "human_users"
}

func (email userHumanEmail) SetEmail(ctx context.Context, client database.QueryExecutor, condition database.Condition, address string, verification domain.VerificationType) (int64, error) {
	switch v := verification.(type) {
	case *domain.VerificationTypeIsVerified:
		return email.setVerifiedEmail(ctx, client, condition, address, v)
	case *domain.VerificationTypeSet:
		return email.setEmailVerificationCode(ctx, client, condition, address, v)
	}
	panic("unknown verification type")
}

func (email userHumanEmail) UpdateEmailVerificationCode(ctx context.Context, client database.QueryExecutor, condition database.Condition, verification *domain.VerificationTypeUpdateVerification) (int64, error) {
	var builder database.StatementBuilder
	builder.WriteString("UPDATE zitadel.verifications SET ")
	if verification.Code != nil {
		email.SetCode(*verification.Code).Write(&builder)
	}
	if verification.Expiry != nil {
		email.setExpiresAt(*verification.Expiry)
	}
	writeCondition(&builder, database.Exists(
		email.unqualifiedTableName(),
		database.And(
			database.NewColumnCondition(email.InstanceIDColumn(), email.verification.InstanceIDColumn()),
			database.NewColumnCondition(email.emailVerificationIDColumn(), email.verification.IDColumn()),
			condition,
		),
	))
	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// GetEmailVerification implements [domain.HumanUserRepository].
func (email userHumanEmail) GetEmailVerification(ctx context.Context, client database.QueryExecutor, condition database.Condition) (*domain.Verification, error) {
	return email.verification.get(ctx, client, database.Exists(
		email.unqualifiedTableName(),
		database.And(
			database.NewColumnCondition(email.InstanceIDColumn(), email.verification.InstanceIDColumn()),
			database.NewColumnCondition(email.emailVerificationIDColumn(), email.verification.IDColumn()),
			condition,
		),
	))
}

func (email userHumanEmail) IncrementFailedEmailVerificationAttempts(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	var builder database.StatementBuilder
	builder.WriteString("UPDATE zitadel.verifications SET failed_attempts = failed_attempts + 1 WHERE (instance_id, id) = (SELECT instance_id, email_verification_id FROM zitadel.human_users")
	writeCondition(&builder, condition)
	builder.WriteRune(')')
	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

func (email userHumanEmail) SetEmailVerifiedAt(verifiedAt time.Time) database.Change {
	if verifiedAt.IsZero() {
		return email.SetEmailVerified()
	}
	changes := make(database.Changes, 0, 2)
	changes = append(changes, email.clearEmailVerificationID())
	return append(changes, database.NewChange(email.EmailVerifiedAtColumn(), verifiedAt))
}

// SetEmailVerified implements [domain.HumanUserRepository].
func (email userHumanEmail) SetEmailVerified() database.Change {
	changes := make(database.Changes, 0, 2)
	changes = append(changes, email.clearEmailVerificationID())
	return append(changes, database.NewChange(email.EmailVerifiedAtColumn(), database.NowInstruction))
}

func (email userHumanEmail) clearEmailVerificationID() database.Change {
	return database.NewChangeToNull(email.emailVerificationIDColumn())
}

func (email userHumanEmail) setAddress(address string) database.Change {
	return database.NewChange(email.EmailColumn(), address)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (email userHumanEmail) EmailCondition(op database.TextOperation, address string) database.Condition {
	return database.NewTextCondition(email.EmailColumn(), op, address)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// EmailColumn implements [domain.HumanUserRepository].
func (email userHumanEmail) EmailColumn() database.Column {
	return database.NewColumn(email.unqualifiedTableName(), "email")
}

// EmailVerifiedAtColumn implements [domain.HumanUserRepository].
func (email userHumanEmail) EmailVerifiedAtColumn() database.Column {
	return database.NewColumn(email.unqualifiedTableName(), "email_verified_at")
}

func (email userHumanEmail) emailVerificationIDColumn() database.Column {
	return database.NewColumn(email.unqualifiedTableName(), "email_verification_id")
}

func (email userHumanEmail) setEmailVerificationCode(ctx context.Context, client database.QueryExecutor, condition database.Condition, address string, verification *domain.VerificationTypeSet) (int64, error) {
	var builder database.StatementBuilder
	builder.WriteString("WITH found_user AS ( SELECT * FROM zitadel.human_users")
	writeCondition(&builder, condition)
	builder.WriteString(" FOR NO KEY UPDATE), verification AS (" +
		" INSERT INTO zitadel.verifications (instance_id, value, code, created_at, expires_at)" +
		" SELECT u.instance_id, ",
	)
	var createdAt any = database.NowInstruction
	if verification.CreatedAt != nil && !verification.CreatedAt.IsZero() {
		createdAt = verification.CreatedAt
	}
	builder.WriteArgs(
		address,
		verification.Code,
		createdAt,
		verification.Expiry,
	)
	builder.WriteString(
		" FROM found_user u" +
			" RETURNING id" +
			") " +
			"UPDATE zitadel.human_users SET password_verification_id = (SELECT id FROM verification) WHERE (instance_id, id) IN (SELECT instance_id, id FROM found_user)",
	)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (email userHumanEmail) setVerifiedEmail(ctx context.Context, client database.QueryExecutor, condition database.Condition, address string, verification *domain.VerificationTypeIsVerified) (int64, error) {
	var builder database.StatementBuilder
	builder.WriteString("UPDATE zitadel.human_users SET ")
	database.Changes{
		email.setAddress(address),
		email.SetEmailVerifiedAt(verification.VerifiedAt),
	}.Write(&builder)
	writeCondition(&builder, condition)
	return client.Exec(ctx, builder.String(), builder.Args()...)
}
