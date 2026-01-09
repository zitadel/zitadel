package repository

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

type userHuman struct {
	user
}

var existingHumanUser = userHuman{
	user: user{
		tableName: "existing_user",
	},
}

func HumanUserRepository() domain.HumanUserRepository {
	return &userHuman{
		user: user{},
	}
}

func (u userHuman) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if !condition.IsRestrictingColumn(u.TypeColumn()) {
		condition = database.And(condition, u.TypeCondition(domain.UserTypeHuman))
	}
	return u.user.Update(ctx, client, condition, changes...)
}

var _ domain.HumanUserRepository = (*userHuman)(nil)

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetAvatarKey implements [domain.HumanUserRepository.SetAvatarKey].
func (u userHuman) SetAvatarKey(avatarKey *string) database.Change {
	return database.NewChangePtr(u.avatarKeyColumn(), avatarKey)
}

// SetDisplayName implements [domain.HumanUserRepository.SetDisplayName].
func (u userHuman) SetDisplayName(displayName string) database.Change {
	return database.NewChange(u.DisplayNameColumn(), displayName)
}

// SetFirstName implements [domain.HumanUserRepository.SetFirstName].
func (u userHuman) SetFirstName(firstName string) database.Change {
	return database.NewChange(u.FirstNameColumn(), firstName)
}

// SetGender implements [domain.HumanUserRepository.SetGender].
func (u userHuman) SetGender(gender domain.HumanGender) database.Change {
	return database.NewChange(u.genderColumn(), gender)
}

// SetLastName implements [domain.HumanUserRepository.SetLastName].
func (u userHuman) SetLastName(lastName string) database.Change {
	return database.NewChange(u.LastNameColumn(), lastName)
}

// SkipMultifactorInitializationAt implements [domain.HumanUserRepository.SkipMultifactorInitializationAt].
func (u userHuman) SkipMultifactorInitializationAt(skippedAt time.Time) database.Change {
	if skippedAt.IsZero() {
		return u.SkipMultifactorInitialization()
	}
	return database.NewChange(u.multifactorInitializationSkippedAtColumn(), skippedAt)
}

// SkipMultifactorInitialization implements [domain.HumanUserRepository.SkipMultifactorInitialization].
func (u userHuman) SkipMultifactorInitialization() database.Change {
	return database.NewChange(u.multifactorInitializationSkippedAtColumn(), database.NowInstruction)
}

// SetNickname implements [domain.HumanUserRepository.SetNickname].
func (u userHuman) SetNickname(nickname string) database.Change {
	return database.NewChange(u.NicknameColumn(), nickname)
}

// SetPassword implements [domain.HumanUserRepository.SetPassword].
func (u userHuman) SetPassword(verification domain.VerificationType) database.Change {
	switch typ := verification.(type) {
	case *domain.VerificationTypeInit:
		var createdAt any = database.NowInstruction
		if !typ.CreatedAt.IsZero() {
			createdAt = typ.CreatedAt
		}
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.verifications (instance_id, user_id, value, code, created_at, expiry) SELECT")
				builder.WriteArgs(
					existingHumanUser.InstanceIDColumn(),
					existingHumanUser.idColumn(),
					typ.Value,
					typ.Code,
					createdAt,
					typ.Expiry,
				)
				builder.WriteString(" FROM ")
				builder.WriteString(existingHumanUser.unqualifiedTableName())
				builder.WriteString(" RETURNING verification.*")
			},
			func(name string) database.Change {
				return database.NewChangeToStatement(
					u.unverifiedPasswordIDColumn(),
					func(builder *database.StatementBuilder) {
						builder.WriteString(" SELECT ")
						existingHumanUser.verification.IDColumn().WriteQualified(builder)
						builder.WriteString(" FROM ")
						builder.WriteString(name)
						writeCondition(builder, database.And(
							database.NewColumnCondition(u.InstanceIDColumn(), database.NewColumn(name, "instance_id")),
							database.NewColumnCondition(u.idColumn(), database.NewColumn(name, "user_id")),
						))
					},
				)
			},
		)
	case *domain.VerificationTypeVerified:
		verifiedAtChange := database.NewChange(u.passwordVerifiedAtColumn(), database.NowInstruction)
		if !typ.VerifiedAt.IsZero() {
			verifiedAtChange = database.NewChange(u.passwordVerifiedAtColumn(), typ.VerifiedAt)
		}

		return database.NewChanges(
			verifiedAtChange,
			database.NewChangeToNull(u.unverifiedPasswordIDColumn()),
			database.NewChange(u.failedPasswordAttemptsColumn(), 0),
			database.NewChangeToStatement(u.passwordColumn(), func(builder *database.StatementBuilder) {
				builder.WriteString("DELETE FROM zitadel.verifications USING existing_user")
				writeCondition(builder, database.And(
					database.NewColumnCondition(u.verification.InstanceIDColumn(), existingHumanUser.InstanceIDColumn()),
					database.NewColumnCondition(u.verification.IDColumn(), existingHumanUser.unverifiedPasswordIDColumn()),
				))
				builder.WriteString(" RETURNING value")
			}),
		)
	case *domain.VerificationTypeUpdate:
		changes := make(database.Changes, 0, 3)
		if typ.Code != nil {
			changes = append(changes, database.NewChange(u.verification.CodeColumn(), typ.Code))
		}
		if typ.Expiry != nil {
			changes = append(changes, database.NewChange(u.verification.ExpiryColumn(), *typ.Expiry))
		}
		if typ.Value != nil {
			changes = append(changes, database.NewChange(u.verification.ValueColumn(), *typ.Value))
		}
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("UPDATE zitadel.verifications SET ")
				changes.Write(builder)
				builder.WriteString(" FROM existing_user")
				writeCondition(builder, database.And(
					database.NewColumnCondition(u.verification.InstanceIDColumn(), existingHumanUser.InstanceIDColumn()),
					database.NewColumnCondition(u.verification.IDColumn(), existingHumanUser.unverifiedPasswordIDColumn()),
				))
			},
			nil,
		)
	case *domain.VerificationTypeSkipped:
		verifiedAtChange := database.NewChange(u.passwordVerifiedAtColumn(), database.NowInstruction)
		if !typ.SkippedAt.IsZero() {
			verifiedAtChange = database.NewChange(u.passwordVerifiedAtColumn(), typ.SkippedAt)
		}
		return database.NewChanges(
			database.NewChangeToNull(u.unverifiedPasswordIDColumn()),
			verifiedAtChange,
			database.NewChange(u.passwordColumn(), *typ.Value),
		)
	case *domain.VerificationTypeFailed:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("UPDATE zitadel.verifications SET verifications.failed_attempts = verifications.failed_attempts + 1")
				builder.WriteString("FROM existing_user")
				writeCondition(builder, database.And(
					database.NewColumnCondition(u.verification.InstanceIDColumn(), existingHumanUser.InstanceIDColumn()),
					database.NewColumnCondition(u.verification.IDColumn(), existingHumanUser.unverifiedPasswordIDColumn()),
				))
			},
			nil,
		)
	}
	panic(fmt.Sprintf("undefined verification type %T", verification))
}

// CheckPassword implements [domain.HumanUserRepository.CheckPassword].
func (u userHuman) CheckPassword(check domain.PasswordCheckType) database.Change {
	switch check.(type) {
	case *domain.CheckTypeSucceeded:
		return database.NewChange(u.failedPasswordAttemptsColumn(), 0)
	case *domain.CheckTypeFailed:
		return database.NewIncrementColumnChange(u.failedPasswordAttemptsColumn())
	}
	panic(fmt.Sprintf("undefined password check type %T", check))
}

// SetPasswordChangeRequired implements [domain.HumanUserRepository.SetPasswordChangeRequired].
func (u userHuman) SetPasswordChangeRequired(required bool) database.Change {
	return database.NewChange(u.passwordChangeRequiredColumn(), true)
}

// SetPreferredLanguage implements [domain.HumanUserRepository.SetPreferredLanguage].
func (u userHuman) SetPreferredLanguage(preferredLanguage language.Tag) database.Change {
	if preferredLanguage == language.Und {
		return database.NewChangeToNull(u.preferredLanguageColumn())
	}
	return database.NewChange(u.preferredLanguageColumn(), preferredLanguage.String())
}

// SetVerification implements [domain.HumanUserRepository.SetVerification].
func (u userHuman) SetVerification(id string, verification domain.VerificationType) database.Change {
	panic("unimplemented")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// DisplayNameCondition implements [domain.HumanUserRepository.DisplayNameCondition].
func (u userHuman) DisplayNameCondition(op database.TextOperation, displayName string) database.Condition {
	return database.NewTextCondition(u.DisplayNameColumn(), op, displayName)
}

// FirstNameCondition implements [domain.HumanUserRepository.FirstNameCondition].
func (u userHuman) FirstNameCondition(op database.TextOperation, firstName string) database.Condition {
	return database.NewTextCondition(u.FirstNameColumn(), op, firstName)
}

// LastNameCondition implements [domain.HumanUserRepository.LastNameCondition].
func (u userHuman) LastNameCondition(op database.TextOperation, lastName string) database.Condition {
	return database.NewTextCondition(u.LastNameColumn(), op, lastName)
}

// NicknameCondition implements [domain.HumanUserRepository.NicknameCondition].
func (u userHuman) NicknameCondition(op database.TextOperation, nickname string) database.Condition {
	return database.NewTextCondition(u.NicknameColumn(), op, nickname)
}

// PreferredLanguageCondition implements [domain.HumanUserRepository.PreferredLanguageCondition].
func (u userHuman) PreferredLanguageCondition(lang language.Tag) database.Condition {
	return database.NewTextCondition(u.preferredLanguageColumn(), database.TextOperationEqual, lang.String())
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// DisplayNameColumn implements [domain.HumanUserRepository.DisplayNameColumn].
func (u userHuman) DisplayNameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "display_name")
}

func (u userHuman) genderColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "gender")
}

func (u userHuman) avatarKeyColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "avatar_key")
}

// FirstNameColumn implements [domain.HumanUserRepository.FirstNameColumn].
func (u userHuman) FirstNameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "first_name")
}

// LastNameColumn implements [domain.HumanUserRepository.LastNameColumn].
func (u userHuman) LastNameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "last_name")
}

// NicknameColumn implements [domain.HumanUserRepository.NicknameColumn].
func (u userHuman) NicknameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "nickname")
}

func (u userHuman) multifactorInitializationSkippedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "multifactor_initialization_skipped_at")
}

func (u userHuman) preferredLanguageColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "preferred_language")
}

func (u userHuman) passwordChangeRequiredColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_change_required")
}

func (u userHuman) unverifiedPasswordIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "unverified_password_id")
}

func (u userHuman) passwordColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password")
}

func (u userHuman) passwordVerifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_verified_at")
}

func (u userHuman) failedPasswordAttemptsColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "failed_password_attempts")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------

// -------------------------------------------------------------
// sub repositories
// -------------------------------------------------------------
