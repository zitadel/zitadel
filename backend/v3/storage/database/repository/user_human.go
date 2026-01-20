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

func (u userHuman) create(ctx context.Context, client database.QueryExecutor, user *domain.User) error {
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
	if !user.Human.Password.VerifiedAt.IsZero() {
		columnValues["password_verified_at"] = user.Human.Password.VerifiedAt
	}
	if !user.Human.PreferredLanguage.IsRoot() {
		columnValues["preferred_language"] = user.Human.PreferredLanguage
	}
	if user.Human.Password.IsChangeRequired {
		columnValues["password_change_required"] = true
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

	// we need to cheat a bit here to be able to use CTEChanges
	// we pretend we have an existing user to be able to use the existing change mechanisms
	builder := database.NewStatementBuilder("WITH existing_user AS (SELECT $1 AS instance_id, $2 AS organization_id, $3 AS id)", user.InstanceID, user.OrganizationID, user.ID)
	ctes := make(map[string]database.CTEChange)

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
		columnValues["unverified_email_id"] = func(builder *database.StatementBuilder) {
			builder.WriteString(`(SELECT id FROM email_verification)`)
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
			columnValues["unverified_phone_id"] = func(builder *database.StatementBuilder) {
				builder.WriteString(`(SELECT id FROM phone_verification)`)
			}
		} else {
			columnValues["phone"] = user.Human.Phone.Number
			columnValues["phone_verified_at"] = user.Human.Phone.VerifiedAt
		}
	}

	for i, passkey := range user.Human.Passkeys {
		ctes[fmt.Sprintf("passkey_%d", i)] = u.AddPasskey(passkey).(database.CTEChange)
	}

	if user.Human.TOTP.Unverified != nil {
		ctes["totp"] = u.SetTOTP(&domain.VerificationTypeInit{
			CreatedAt: user.CreatedAt,
			Code:      user.Human.TOTP.Unverified.Code,
			Value:     user.Human.TOTP.Unverified.Value,
		}).(database.CTEChange)
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
				builder.WriteString(" RETURNING verifications.*")
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
func (u userHuman) SetVerification(verification domain.VerificationType) database.Change {
	// panic("not implemented")
	switch typ := verification.(type) {
	case *domain.VerificationTypeInit:
		var createdAt any = database.NowInstruction
		if !typ.CreatedAt.IsZero() {
			createdAt = typ.CreatedAt
		}
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("INSERT INTO zitadel.verifications (instance_id, user_id, id, value, code, created_at, expiry) SELECT existing_user.instance_id, existing_user.id, ")
				builder.WriteArgs(
					typ.ID,
					typ.Value,
					typ.Code,
					createdAt,
					typ.Expiry,
				)
				builder.WriteString(" FROM ")
				builder.WriteString(existingHumanUser.unqualifiedTableName())
				builder.WriteString(" RETURNING verifications.*")
			}, nil,
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
					database.NewTextCondition(u.verification.IDColumn(), database.TextOperationEqual, *typ.ID),
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
					database.NewTextCondition(u.verification.IDColumn(), database.TextOperationEqual, *typ.ID),
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
					database.NewColumnCondition(u.verification.InstanceIDColumn(), existingHumanUser.InstanceIDColumn()),
					database.NewTextCondition(u.verification.IDColumn(), database.TextOperationEqual, *typ.ID),
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
