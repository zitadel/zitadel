package repository

import (
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetInviteVerification implements [domain.HumanUserRepository].
func (u userHuman) SetInviteVerification(verification domain.VerificationType) database.Change {
	switch typ := verification.(type) {
	case *domain.VerificationTypeInit:
		return u.verification.init(typ, existingUser.unqualifiedTableName(), existingHumanUser.inviteVerificationIDColumn())
	case *domain.VerificationTypeUpdate:
		return u.verification.update(typ, existingHumanUser.unqualifiedTableName(),
			existingHumanUser.InstanceIDColumn(), existingHumanUser.inviteVerificationIDColumn(),
		)
	case *domain.VerificationTypeSucceeded:
		return u.verification.verified(typ, existingHumanUser.unqualifiedTableName(), existingHumanUser.InstanceIDColumn(),
			existingHumanUser.inviteVerificationIDColumn(), u.inviteAcceptedAtColumn(), u.failedInviteAttemptsColumn(),
		)
	case *domain.VerificationTypeSkipped:
		return u.verification.skipped(typ, u.inviteAcceptedAtColumn(), u.inviteVerificationIDColumn(), u.failedInviteAttemptsColumn())
	case *domain.VerificationTypeFailed:
		return u.verification.failed(existingHumanUser.unqualifiedTableName(), existingHumanUser.InstanceIDColumn(), existingHumanUser.inviteVerificationIDColumn())
	}
	panic(fmt.Sprintf("undefined verification type %T", verification))
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userHuman) inviteVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "invite_verification_id")
}

func (u userHuman) inviteAcceptedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "invite_accepted_at")
}

func (u userHuman) failedInviteAttemptsColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "invite_failed_attempts")
}
