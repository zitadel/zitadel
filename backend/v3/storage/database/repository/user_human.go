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
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
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
	}
	if !user.Human.MultifactorInitializationSkippedAt.IsZero() {
		columnValues["multifactor_initialization_skipped_at"] = user.Human.MultifactorInitializationSkippedAt
	}
	columnValues = u.setOptionalProfileColumns(columnValues, user.Human)

	ctes := make(map[string]database.CTEChange)

	columnValues, ctes = u.setPasswordColumns(columnValues, ctes, user, createdAt)
	columnValues, ctes = u.setEmailColumns(columnValues, ctes, user, createdAt)
	columnValues, ctes = u.setPhoneColumns(columnValues, ctes, user, createdAt)
	columnValues, ctes = u.setInviteColumns(columnValues, ctes, user, createdAt)

	for i, passkey := range user.Human.Passkeys {
		ctes[fmt.Sprintf("passkey_%d", i)] = u.AddPasskey(passkey).(database.CTEChange)
	}

	if user.Human.TOTP != nil {
		columnValues["totp_secret"] = user.Human.TOTP.Secret
		if !user.Human.TOTP.VerifiedAt.IsZero() {
			columnValues["totp_verified_at"] = user.Human.TOTP.VerifiedAt
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
		}).(database.CTEChange)
	}

	for i, metadata := range user.Metadata {
		ctes[fmt.Sprintf("metadata_%d", i)] = u.SetMetadata(metadata).(database.CTEChange)
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

func (u userHuman) setPasswordColumns(columnValues map[string]any, ctes map[string]database.CTEChange, user *domain.User, createdAt any) (map[string]any, map[string]database.CTEChange) {
	if user.Human.Password.PendingVerification != nil {
		verification := &domain.VerificationTypeInit{
			CreatedAt: user.CreatedAt,
			Code:      user.Human.Password.PendingVerification.Code,
		}
		if user.Human.Password.PendingVerification.ID != "" {
			verification.ID = &user.Human.Password.PendingVerification.ID
		}
		if user.Human.Password.PendingVerification.ExpiresAt != nil && !user.Human.Password.PendingVerification.ExpiresAt.IsZero() {
			verification.Expiry = gu.Ptr(time.Since(*user.Human.Password.PendingVerification.ExpiresAt))
		}

		ctes["password_verification"] = u.SetResetPasswordVerification(verification).(database.CTEChange)
		columnValues["password_verification_id"] = func(builder *database.StatementBuilder) {
			builder.WriteString(`(SELECT id FROM password_verification)`)
		}
	} else if user.Human.Password.Hash != "" {
		columnValues["password_hash"] = user.Human.Password.Hash
		columnValues["password_changed_at"] = createdAt
		if !user.Human.Password.ChangedAt.IsZero() {
			columnValues["password_changed_at"] = user.Human.Password.ChangedAt
		}
	}
	return columnValues, ctes
}

func (u userHuman) setEmailColumns(columnValues map[string]any, ctes map[string]database.CTEChange, user *domain.User, createdAt any) (map[string]any, map[string]database.CTEChange) {
	if user.Human.Email.PendingVerification != nil {
		verification := &domain.VerificationTypeInit{
			CreatedAt: user.CreatedAt,
			Code:      user.Human.Email.PendingVerification.Code,
		}
		if user.Human.Email.PendingVerification.ID != "" {
			verification.ID = &user.Human.Email.PendingVerification.ID
		}
		if user.Human.Email.PendingVerification.ExpiresAt != nil && !user.Human.Email.PendingVerification.ExpiresAt.IsZero() {
			verification.Expiry = gu.Ptr(time.Since(*user.Human.Email.PendingVerification.ExpiresAt))
		}

		ctes["email_verification"] = u.SetEmailVerification(verification).(database.CTEChange)
		columnValues["email_verification_id"] = func(builder *database.StatementBuilder) {
			builder.WriteString(`(SELECT id FROM email_verification)`)
		}
	}
	if user.Human.Email.Address != "" {
		columnValues["email"] = user.Human.Email.Address
		columnValues["unverified_email"] = user.Human.Email.Address
		columnValues["email_verified_at"] = createdAt
		if !user.Human.Email.VerifiedAt.IsZero() {
			columnValues["email_verified_at"] = user.Human.Email.VerifiedAt
		}
	}
	if user.Human.Email.UnverifiedAddress != "" {
		columnValues["unverified_email"] = user.Human.Email.UnverifiedAddress
	}
	return columnValues, ctes
}

func (u userHuman) setPhoneColumns(columnValues map[string]any, ctes map[string]database.CTEChange, user *domain.User, createdAt any) (map[string]any, map[string]database.CTEChange) {
	if user.Human.Phone == nil {
		return columnValues, ctes
	}
	if user.Human.Phone.PendingVerification != nil {
		verification := &domain.VerificationTypeInit{
			CreatedAt: user.CreatedAt,
			Code:      user.Human.Phone.PendingVerification.Code,
		}
		if user.Human.Phone.PendingVerification.ID != "" {
			verification.ID = &user.Human.Phone.PendingVerification.ID
		}
		if user.Human.Phone.PendingVerification.ExpiresAt != nil && !user.Human.Phone.PendingVerification.ExpiresAt.IsZero() {
			verification.Expiry = gu.Ptr(time.Since(*user.Human.Phone.PendingVerification.ExpiresAt))
		}

		ctes["phone_verification"] = u.SetPhoneVerification(verification).(database.CTEChange)
		columnValues["phone_verification_id"] = func(builder *database.StatementBuilder) {
			builder.WriteString(`(SELECT id FROM phone_verification)`)
		}
	}
	if user.Human.Phone.Number != "" {
		columnValues["phone"] = user.Human.Phone.Number
		columnValues["unverified_phone"] = user.Human.Phone.Number
		columnValues["phone_verified_at"] = createdAt
		if !user.Human.Phone.VerifiedAt.IsZero() {
			columnValues["phone_verified_at"] = user.Human.Phone.VerifiedAt
		}
	}
	if user.Human.Phone.UnverifiedNumber != "" {
		columnValues["unverified_phone"] = user.Human.Phone.UnverifiedNumber
	}
	return columnValues, ctes
}

func (u userHuman) setInviteColumns(columnValues map[string]any, ctes map[string]database.CTEChange, user *domain.User, createdAt any) (map[string]any, map[string]database.CTEChange) {
	if user.Human.Invite == nil {
		return columnValues, ctes
	}
	if user.Human.Invite.PendingVerification != nil {
		verification := &domain.VerificationTypeInit{
			CreatedAt: user.CreatedAt,
			Code:      user.Human.Invite.PendingVerification.Code,
		}
		if user.Human.Invite.PendingVerification.ID != "" {
			verification.ID = &user.Human.Invite.PendingVerification.ID
		}
		if user.Human.Invite.PendingVerification.ExpiresAt != nil && !user.Human.Invite.PendingVerification.ExpiresAt.IsZero() {
			verification.Expiry = gu.Ptr(time.Since(*user.Human.Invite.PendingVerification.ExpiresAt))
		}
		ctes["invite_verification"] = u.SetInviteVerification(verification).(database.CTEChange)
	}
	if !user.Human.Invite.AcceptedAt.IsZero() {
		columnValues["invite_accepted_at"] = createdAt
	}
	return columnValues, ctes
}

func (u userHuman) setOptionalProfileColumns(columnValues map[string]any, human *domain.HumanUser) map[string]any {
	if !human.PreferredLanguage.IsRoot() {
		columnValues["preferred_language"] = human.PreferredLanguage
	}
	if human.DisplayName != "" {
		columnValues["display_name"] = human.DisplayName
	}
	if human.Nickname != "" {
		columnValues["nickname"] = human.Nickname
	}
	if human.Gender != domain.HumanGenderUnspecified {
		columnValues["gender"] = human.Gender
	}
	if human.AvatarKey != "" {
		columnValues["avatar_key"] = human.AvatarKey
	}
	if human.Password.IsChangeRequired {
		columnValues["password_change_required"] = true
	}
	return columnValues
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

// SetAvatarKey implements [domain.HumanUserRepository].
func (u userHuman) SetAvatarKey(avatarKey *string) database.Change {
	return database.NewChangePtr(u.avatarKeyColumn(), avatarKey)
}

// SetDisplayName implements [domain.HumanUserRepository].
func (u userHuman) SetDisplayName(displayName string) database.Change {
	return database.NewChange(u.DisplayNameColumn(), displayName)
}

// SetFirstName implements [domain.HumanUserRepository].
func (u userHuman) SetFirstName(firstName string) database.Change {
	return database.NewChange(u.FirstNameColumn(), firstName)
}

// SetGender implements [domain.HumanUserRepository].
func (u userHuman) SetGender(gender domain.HumanGender) database.Change {
	return database.NewChange(u.genderColumn(), gender)
}

// SetLastName implements [domain.HumanUserRepository].
func (u userHuman) SetLastName(lastName string) database.Change {
	return database.NewChange(u.LastNameColumn(), lastName)
}

// SkipMultifactorInitializationAt implements [domain.HumanUserRepository].
func (u userHuman) SkipMultifactorInitializationAt(skippedAt time.Time) database.Change {
	if skippedAt.IsZero() {
		return u.SkipMultifactorInitialization()
	}
	return database.NewChange(u.multifactorInitializationSkippedAtColumn(), skippedAt)
}

// SkipMultifactorInitialization implements [domain.HumanUserRepository].
func (u userHuman) SkipMultifactorInitialization() database.Change {
	return database.NewChange(u.multifactorInitializationSkippedAtColumn(), database.NowInstruction)
}

// SetNickname implements [domain.HumanUserRepository].
func (u userHuman) SetNickname(nickname string) database.Change {
	return database.NewChange(u.NicknameColumn(), nickname)
}

// SetPreferredLanguage implements [domain.HumanUserRepository].
func (u userHuman) SetPreferredLanguage(preferredLanguage language.Tag) database.Change {
	if preferredLanguage == language.Und {
		return database.NewChangeToNull(u.preferredLanguageColumn())
	}
	return database.NewChange(u.preferredLanguageColumn(), preferredLanguage.String())
}

// SetVerification implements [domain.HumanUserRepository].
func (u userHuman) SetVerification(verification domain.VerificationType) database.Change {
	switch typ := verification.(type) {
	case *domain.VerificationTypeInit:
		return database.NewCTEChange(func(builder *database.StatementBuilder) {
			var (
				createdAt any = database.NowInstruction
				expiry    any = database.NullInstruction
				id        any = database.DefaultInstruction
			)
			if !typ.CreatedAt.IsZero() {
				createdAt = typ.CreatedAt
			}
			if typ.Expiry != nil {
				expiry = *typ.Expiry
			}
			if typ.ID != nil {
				id = *typ.ID
			}
			builder.WriteString("INSERT INTO zitadel.verifications(instance_id, user_id, id, code, created_at, expiry) SELECT instance_id, id, ")
			builder.WriteArgs(id, typ.Code, createdAt, expiry)
			builder.WriteString(" FROM ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			builder.WriteString(" RETURNING id")
		}, nil)
	case *domain.VerificationTypeSucceeded:
		return database.NewCTEChange(func(builder *database.StatementBuilder) {
			builder.WriteString("DELETE FROM zitadel.verifications USING existing_user")

			writeCondition(builder, database.And(
				database.NewColumnCondition(u.verification.instanceIDColumn(), existingHumanUser.InstanceIDColumn()),
				database.NewColumnCondition(u.verification.userIDColumn(), existingHumanUser.IDColumn()),
				database.NewTextCondition(u.verification.idColumn(), database.TextOperationEqual, *typ.ID),
			))
		}, nil)

	case *domain.VerificationTypeUpdate:
		changes := make(database.Changes, 0, 3)
		if typ.Code != nil {
			changes = append(changes, database.NewChange(u.verification.codeColumn(), typ.Code))
		}
		if typ.Expiry != nil {
			changes = append(changes, database.NewChange(u.verification.expiryColumn(), *typ.Expiry))
		}
		return database.NewCTEChange(
			func(builder *database.StatementBuilder) {
				builder.WriteString("UPDATE zitadel.verifications SET ")
				err := changes.Write(builder)
				logging.New(logging.StreamRuntime).Debug("write changes in cte failed", "error", err)
				builder.WriteString(" FROM existing_user")
				writeCondition(builder, database.And(
					database.NewColumnCondition(u.verification.instanceIDColumn(), existingHumanUser.InstanceIDColumn()),
					database.NewColumnCondition(u.verification.userIDColumn(), existingHumanUser.IDColumn()),
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
				builder.WriteString("UPDATE ")
				builder.WriteString(u.verification.qualifiedTableName())
				builder.WriteString(" SET ")
				err := database.NewIncrementColumnChange(u.verification.failedAttemptsColumn(), database.Coalesce(u.verification.failedAttemptsColumn(), 0)).Write(builder)
				logging.New(logging.StreamRuntime).Debug("write cte failed", "error", err)
				builder.WriteString(" FROM existing_user")
				writeCondition(builder, database.And(
					database.NewColumnCondition(u.verification.instanceIDColumn(), existingHumanUser.InstanceIDColumn()),
					database.NewColumnCondition(u.verification.userIDColumn(), existingHumanUser.IDColumn()),
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

// DisplayNameCondition implements [domain.HumanUserRepository].
func (u userHuman) DisplayNameCondition(op database.TextOperation, displayName string) database.Condition {
	return database.NewTextCondition(u.DisplayNameColumn(), op, displayName)
}

// FirstNameCondition implements [domain.HumanUserRepository].
func (u userHuman) FirstNameCondition(op database.TextOperation, firstName string) database.Condition {
	return database.NewTextCondition(u.FirstNameColumn(), op, firstName)
}

// LastNameCondition implements [domain.HumanUserRepository].
func (u userHuman) LastNameCondition(op database.TextOperation, lastName string) database.Condition {
	return database.NewTextCondition(u.LastNameColumn(), op, lastName)
}

// NicknameCondition implements [domain.HumanUserRepository].
func (u userHuman) NicknameCondition(op database.TextOperation, nickname string) database.Condition {
	return database.NewTextCondition(u.NicknameColumn(), op, nickname)
}

// PreferredLanguageCondition implements [domain.HumanUserRepository].
func (u userHuman) PreferredLanguageCondition(lang language.Tag) database.Condition {
	return database.NewTextCondition(u.preferredLanguageColumn(), database.TextOperationEqual, lang.String())
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// DisplayNameColumn implements [domain.HumanUserRepository].
func (u userHuman) DisplayNameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "display_name")
}

func (u userHuman) genderColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "gender")
}

func (u userHuman) avatarKeyColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "avatar_key")
}

// FirstNameColumn implements [domain.HumanUserRepository].
func (u userHuman) FirstNameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "first_name")
}

// LastNameColumn implements [domain.HumanUserRepository].
func (u userHuman) LastNameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "last_name")
}

// NicknameColumn implements [domain.HumanUserRepository].
func (u userHuman) NicknameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "nickname")
}

func (u userHuman) multifactorInitializationSkippedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "multifactor_initialization_skipped_at")
}

func (u userHuman) preferredLanguageColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "preferred_language")
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------

// -------------------------------------------------------------
// sub repositories
// -------------------------------------------------------------
