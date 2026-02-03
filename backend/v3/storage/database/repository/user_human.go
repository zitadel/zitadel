package repository

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"github.com/muhlemmer/gu"
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

func (u userHuman) create(ctx context.Context, builder *database.StatementBuilder, client database.QueryExecutor, user *domain.User) error {
	var createdAt any = database.NowInstruction
	if !user.CreatedAt.IsZero() {
		createdAt = user.CreatedAt
	}
	columnValues := map[string]any{
		"instance_id":     user.InstanceID,
		"organization_id": user.OrganizationID,
		"id":              user.ID,
		"username":        user.Username,
		"state":           user.State,
		"type":            "human",
		"first_name":      user.Human.FirstName,
		"last_name":       user.Human.LastName,
		"created_at":      createdAt,
		"updated_at":      createdAt,
		"password":        user.Human.Password.Password,
		"email":           user.Human.Email.Address,
	}
	if !user.Human.MultifactorInitializationSkippedAt.IsZero() {
		columnValues["multifactor_initialization_skipped_at"] = user.Human.MultifactorInitializationSkippedAt
	}
	if !user.Human.PreferredLanguage.IsRoot() {
		columnValues["preferred_language"] = user.Human.PreferredLanguage
	}
	if user.Human.Password.IsChangeRequired {
		columnValues["password_change_required"] = true
	}
	if user.Human.Password.Unverified == nil {
		// TODO: do we need to create a verification if [user.Human.Password.Unverified] is not nil?
		columnValues["password_verified_at"] = createdAt
	}
	if user.Human.DisplayName != "" {
		columnValues["display_name"] = user.Human.DisplayName
	}
	if user.Human.Nickname != "" {
		columnValues["nickname"] = user.Human.Nickname
	}
	if user.Human.Gender != domain.HumanGenderUnspecified {
		columnValues["gender"] = user.Human.Gender
	}
	if user.Human.AvatarKey != "" {
		columnValues["avatar_key"] = user.Human.AvatarKey
	}

	ctes := make(map[string]database.CTEChange)

	// TODO: can we handle email the same as phone?
	// meaning that we either set it verified or unverified
	if user.Human.Email.Unverified != nil {
		verification := &domain.VerificationTypeInit{
			Value:     &user.Human.Email.Address,
			CreatedAt: user.CreatedAt,
			Code:      user.Human.Email.Unverified.Code,
		}
		if user.Human.Email.Unverified.Value != nil {
			verification.Value = user.Human.Email.Unverified.Value
		}
		if user.Human.Email.Unverified.ExpiresAt != nil && !user.Human.Email.Unverified.ExpiresAt.IsZero() {
			verification.Expiry = gu.Ptr(time.Since(*user.Human.Email.Unverified.ExpiresAt))
		}

		change := u.SetEmail(verification).(database.CTEChange)
		ctes["email_verification"] = change
		columnValues["email_verification_id"] = func(builder *database.StatementBuilder) {
			builder.WriteString(`(SELECT id FROM email_verification)`)
		}
	} else {
		columnValues["email"] = user.Human.Email.Address
		if !user.Human.Email.VerifiedAt.IsZero() {
			columnValues["email_verified_at"] = user.Human.Email.VerifiedAt
		} else {
			columnValues["email_verified_at"] = database.NowInstruction
		}
	}
	if user.Human.Phone != nil {
		if user.Human.Phone.Unverified != nil {
			verification := &domain.VerificationTypeInit{
				Value:     &user.Human.Phone.Number,
				CreatedAt: user.CreatedAt,
				Code:      user.Human.Phone.Unverified.Code,
			}
			if user.Human.Phone.Unverified.Value != nil {
				verification.Value = user.Human.Phone.Unverified.Value
			}
			if user.Human.Phone.Unverified.ExpiresAt != nil && !user.Human.Phone.Unverified.ExpiresAt.IsZero() {
				verification.Expiry = gu.Ptr(time.Since(*user.Human.Phone.Unverified.ExpiresAt))
			}

			change := u.SetPhone(verification).(database.CTEChange)
			ctes["phone_verification"] = change
			columnValues["phone_verification_id"] = func(builder *database.StatementBuilder) {
				builder.WriteString(`(SELECT id FROM phone_verification)`)
			}
		} else {
			columnValues["phone"] = user.Human.Phone.Number
			if !user.Human.Phone.VerifiedAt.IsZero() {
				columnValues["phone_verified_at"] = user.Human.Phone.VerifiedAt
			} else {
				columnValues["phone_verified_at"] = database.NowInstruction
			}
		}
	}

	for i, passkey := range user.Human.Passkeys {
		ctes[fmt.Sprintf("passkey_%d", i)] = u.AddPasskey(passkey).(database.CTEChange)
	}

	if user.Human.TOTP != nil {
		columnValues["totp_secret"] = user.Human.TOTP.Secret
		if !user.Human.TOTP.VerifiedAt.IsZero() {
			columnValues["totp_verified_at"] = user.Human.TOTP.VerifiedAt
		} else {
			columnValues["totp_verified_at"] = database.NowInstruction
		}
	}

	for i, link := range user.Human.IdentityProviderLinks {
		name := fmt.Sprintf("idp_link_%d", i)
		ctes[name] = u.AddIdentityProviderLink(link).(database.CTEChange)
	}

	for i, verification := range user.Human.Verifications {
		ctes[fmt.Sprintf("verification_%d", i)] = u.SetVerification(&domain.VerificationTypeInit{
			ID:        gu.Ptr(verification.ID),
			CreatedAt: user.CreatedAt,
			Code:      verification.Code,
			Value:     verification.Value,
		}).(database.CTEChange)
	}

	// write CTE changes
	for name, cte := range ctes {
		cte.SetName(name)
		builder.WriteString(", ")
		builder.WriteString(name)
		builder.WriteString(" AS (")
		cte.WriteCTE(builder)
		builder.WriteString(")")
	}
	columns := slices.Sorted(maps.Keys(columnValues))

	// write final insert
	builder.WriteString(" INSERT INTO zitadel.users (")
	builder.WriteString(strings.Join(columns, ", "))
	builder.WriteString(") VALUES (")
	for i, column := range columns {
		if i > 0 {
			builder.WriteString(", ")
		}
		value := columnValues[column]
		if fn, ok := value.(func(builder *database.StatementBuilder)); ok {
			fn(builder)
			continue
		}
		builder.WriteArg(value)
	}
	builder.WriteString(") RETURNING created_at, updated_at")

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&user.CreatedAt, &user.UpdatedAt)
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
		return u.verification.setInit(typ, existingUser.unqualifiedTableName(), existingHumanUser.passwordVerificationIDColumn())
	case *domain.VerificationTypeVerified:
		return u.verification.verified(typ, existingUser.unqualifiedTableName(), u.InstanceIDColumn(),
			u.passwordVerificationIDColumn(), u.passwordVerifiedAtColumn(), u.passwordColumn())
	case *domain.VerificationTypeUpdate:
		return u.verification.update(typ, existingHumanUser.unqualifiedTableName(),
			existingHumanUser.InstanceIDColumn(), existingHumanUser.passwordVerificationIDColumn(),
		)
	case *domain.VerificationTypeSkipped:
		return u.verification.skipped(typ, u.passwordVerifiedAtColumn(), u.passwordColumn())
	case *domain.VerificationTypeFailed:
		return u.verification.failed(existingHumanUser.unqualifiedTableName(), existingHumanUser.InstanceIDColumn(), existingHumanUser.passwordVerificationIDColumn())
	}
	panic(fmt.Sprintf("undefined verification type %T", verification))
}

// SetPasswordChangeRequired implements [domain.HumanUserRepository.SetPasswordChangeRequired].
func (u userHuman) SetPasswordChangeRequired(required bool) database.Change {
	return database.NewChange(u.passwordChangeRequiredColumn(), required)
}

func (u userHuman) SetLastSuccessfulPasswordCheck(checkedAt time.Time) database.Change {
	if checkedAt.IsZero() {
		return database.NewChange(u.lastSuccessfulPasswordCheckColumn(), database.NowInstruction)
	}
	return database.NewChange(u.lastSuccessfulPasswordCheckColumn(), checkedAt)
}

func (u userHuman) IncrementPasswordFailedAttempts() database.Change {
	return database.NewIncrementColumnChange(u.failedPasswordAttemptsColumn())
}

func (u userHuman) ResetPasswordFailedAttempts() database.Change {
	return database.NewChange(u.failedPasswordAttemptsColumn(), 0)
}

// SetPreferredLanguage implements [domain.HumanUserRepository.SetPreferredLanguage].
func (u userHuman) SetPreferredLanguage(preferredLanguage language.Tag) database.Change {
	if preferredLanguage == language.Und {
		return database.NewChangeToNull(u.preferredLanguageColumn())
	}
	return database.NewChange(u.preferredLanguageColumn(), preferredLanguage.String())
}

// SetVerification implements [domain.HumanUserRepository.SetVerification].
func (u userHuman) SetVerification(verification domain.VerificationType) database.Change {
	switch typ := verification.(type) {
	case *domain.VerificationTypeInit:
		return u.verification.init(typ, existingHumanUser.unqualifiedTableName())
	case *domain.VerificationTypeVerified:
		// return u.verification.verified()
		verifiedAtChange := database.NewChange(u.passwordVerifiedAtColumn(), database.NowInstruction)
		if !typ.VerifiedAt.IsZero() {
			verifiedAtChange = database.NewChange(u.passwordVerifiedAtColumn(), typ.VerifiedAt)
		}

		return database.NewChanges(
			verifiedAtChange,
			database.NewChangeToNull(u.passwordVerificationIDColumn()),
			database.NewChange(u.failedPasswordAttemptsColumn(), 0),
			database.NewChangeToStatement(u.passwordColumn(), func(builder *database.StatementBuilder) {
				builder.WriteString("DELETE FROM zitadel.verifications USING existing_user")
				writeCondition(builder, database.And(
					database.NewColumnCondition(u.verification.instanceIDColumn(), existingHumanUser.InstanceIDColumn()),
					database.NewTextCondition(u.verification.idColumn(), database.TextOperationEqual, *typ.ID),
				))
				builder.WriteString(" RETURNING value")
			}),
		)
	case *domain.VerificationTypeUpdate:
		changes := make(database.Changes, 0, 3)
		if typ.Code != nil {
			changes = append(changes, database.NewChange(u.verification.codeColumn(), typ.Code))
		}
		if typ.Expiry != nil {
			changes = append(changes, database.NewChange(u.verification.expiryColumn(), *typ.Expiry))
		}
		if typ.Value != nil {
			changes = append(changes, database.NewChange(u.verification.valueColumn(), *typ.Value))
		}
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("UPDATE zitadel.verifications SET ")
				changes.Write(builder)
				builder.WriteString(" FROM existing_user")
				writeCondition(builder, database.And(
					database.NewColumnCondition(u.verification.instanceIDColumn(), existingHumanUser.InstanceIDColumn()),
					database.NewTextCondition(u.verification.idColumn(), database.TextOperationEqual, *typ.ID),
				))
			},
			nil,
		)
	case *domain.VerificationTypeSkipped:
		panic("skip verification is not supported for verifications")
	case *domain.VerificationTypeFailed:
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("UPDATE zitadel.verifications SET verifications.failed_attempts = verifications.failed_attempts + 1")
				builder.WriteString("FROM existing_user")
				writeCondition(builder, database.And(
					database.NewColumnCondition(u.verification.instanceIDColumn(), existingHumanUser.InstanceIDColumn()),
					database.NewTextCondition(u.verification.idColumn(), database.TextOperationEqual, *typ.ID),
				))
			},
			nil,
		)
	}
	panic(fmt.Sprintf("undefined verification type %T", verification))
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

func (u userHuman) passwordColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password")
}

func (u userHuman) passwordChangeRequiredColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_change_required")
}

func (u userHuman) passwordVerifiedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_verified_at")
}

func (u userHuman) passwordVerificationIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_verification_id")
}

func (u userHuman) lastSuccessfulPasswordCheckColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_last_successful_check")
}

func (u userHuman) failedPasswordAttemptsColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "password_failed_attempts")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------

// -------------------------------------------------------------
// sub repositories
// -------------------------------------------------------------
