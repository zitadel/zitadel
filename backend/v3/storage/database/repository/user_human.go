package repository

import (
	"context"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.HumanUserRepository = (*userHuman)(nil)

type userHuman struct{}

func (h userHuman) unqualifiedTableName() string {
	return "human_users"
}

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

const insertHumanStmt = "INSERT INTO zitadel.human_users (" +
	"instance_id, organization_id, id, username, username_org_unique, state, created_at, updated_at" +
	", first_name, last_name, nickname, display_name, preferred_language, gender, avatar_key" +
	") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING created_at, updated_at"

// create inserts a new human user into the database.
// the type of the user must be checked before calling this method.
func (u userHuman) create(ctx context.Context, client database.QueryExecutor, user *domain.User) error {
	var createdAt, updatedAt any = database.DefaultInstruction, database.DefaultInstruction
	if !user.CreatedAt.IsZero() {
		createdAt = user.CreatedAt
	}
	if !user.UpdatedAt.IsZero() {
		updatedAt = user.UpdatedAt
	}

	return client.QueryRow(
		ctx, insertHumanStmt,
		user.InstanceID,
		user.OrgID,
		user.ID,
		user.Username,
		user.IsUsernameOrgUnique,
		user.State,
		createdAt,
		updatedAt,
		user.Human.FirstName,
		user.Human.LastName,
		user.Human.Nickname,
		user.Human.DisplayName,
		user.Human.PreferredLanguage,
		user.Human.Gender,
		user.Human.AvatarKey,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

// // Security implements [domain.HumanUserRepository].
// func (u userHuman) Security() domain.HumanSecurityRepository {
// 	panic("unimplemented")
// }

// Update implements [domain.HumanUserRepository].
func (u userHuman) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}
	if err := checkPKCondition(u, condition); err != nil {
		return 0, err
	}
	if !database.Changes(changes).IsOnColumn(u.UpdatedAtColumn()) {
		changes = append(changes, database.NewChange(u.UpdatedAtColumn(), database.NullInstruction))
	}
	builder := database.NewStatementBuilder(`UPDATE zitadel.human_users SET `)
	database.Changes(changes).Write(builder)
	writeCondition(builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetAvatarKey implements [domain.HumanUserRepository].
func (u userHuman) SetAvatarKey(key *string) database.Change {
	return database.NewChangePtr(u.AvatarKeyColumn(), key)
}

// SetDisplayName implements [domain.HumanUserRepository].
func (u userHuman) SetDisplayName(name string) database.Change {
	return database.NewChange(u.DisplayNameColumn(), name)
}

// SetFirstName implements [domain.HumanUserRepository].
func (u userHuman) SetFirstName(name string) database.Change {
	return database.NewChange(u.FirstNameColumn(), name)
}

// SetGender implements [domain.HumanUserRepository].
func (u userHuman) SetGender(gender *domain.Gender) database.Change {
	if gender == nil || *gender == domain.GenderUnspecified {
		return database.NewChangeToNull(u.GenderColumn())
	}
	return database.NewChange(u.GenderColumn(), *gender)
}

// SetLastName implements [domain.HumanUserRepository].
func (u userHuman) SetLastName(name string) database.Change {
	return database.NewChange(u.LastNameColumn(), name)
}

// SetNickname implements [domain.HumanUserRepository].
func (u userHuman) SetNickname(name string) database.Change {
	return database.NewChange(u.NicknameColumn(), name)
}

// SetPreferredLanguage implements [domain.HumanUserRepository].
func (u userHuman) SetPreferredLanguage(lang *language.Tag) database.Change {
	if lang == nil || *lang == language.Und {
		return database.NewChangeToNull(u.PreferredLanguageColumn())
	}
	return database.NewChange(u.PreferredLanguageColumn(), lang.String())
}

// SetState implements [domain.HumanUserRepository].
func (u userHuman) SetState(state domain.UserState) database.Change {
	return database.NewChange(u.StateColumn(), state)
}

// SetUpdatedAt implements [domain.HumanUserRepository].
func (u userHuman) SetUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(u.UpdatedAtColumn(), updatedAt)
}

// SetUsername implements [domain.HumanUserRepository].
func (u userHuman) SetUsername(username string) database.Change {
	return database.NewChange(u.UsernameColumn(), username)
}

// SetUsernameOrgUnique implements [domain.HumanUserRepository].
func (u userHuman) SetUsernameOrgUnique(usernameOrgUnique bool) database.Change {
	return database.NewChange(u.UsernameOrgUniqueColumn(), usernameOrgUnique)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// PrimaryKeyCondition implements domain.HumanUserRepository.
func (h userHuman) PrimaryKeyCondition(instanceID string, userID string) database.Condition {
	return database.And(
		h.InstanceIDCondition(instanceID),
		h.IDCondition(userID),
	)
}

// TypeCondition implements domain.HumanUserRepository.
func (h userHuman) TypeCondition(userType domain.UserType) database.Condition {
	// TODO(adlerhurst): it doesn't make sense to have this method on userHuman
	return user{}.TypeCondition(userType)
}

// CreatedAtCondition implements [domain.HumanUserRepository].
func (u userHuman) CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition {
	return database.NewNumberCondition(u.CreatedAtColumn(), op, createdAt)
}

// DisplayNameCondition implements [domain.HumanUserRepository].
func (u userHuman) DisplayNameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(u.DisplayNameColumn(), op, name)
}

// FirstNameCondition implements [domain.HumanUserRepository].
func (u userHuman) FirstNameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(u.FirstNameColumn(), op, name)
}

// IDCondition implements [domain.HumanUserRepository].
func (u userHuman) IDCondition(userID string) database.Condition {
	return database.NewTextCondition(u.IDColumn(), database.TextOperationEqual, userID)
}

// InstanceIDCondition implements [domain.HumanUserRepository].
func (u userHuman) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(u.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// LastNameCondition implements [domain.HumanUserRepository].
func (u userHuman) LastNameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(u.LastNameColumn(), op, name)
}

// NicknameCondition implements [domain.HumanUserRepository].
func (u userHuman) NicknameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(u.NicknameColumn(), op, name)
}

// OrgIDCondition implements [domain.HumanUserRepository].
func (u userHuman) OrgIDCondition(orgID string) database.Condition {
	return database.NewTextCondition(u.OrgIDColumn(), database.TextOperationEqual, orgID)
}

// StateCondition implements [domain.HumanUserRepository].
func (u userHuman) StateCondition(state domain.UserState) database.Condition {
	return database.NewNumberCondition(u.StateColumn(), database.NumberOperationEqual, state)
}

// UpdatedAtCondition implements [domain.HumanUserRepository].
func (u userHuman) UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition {
	return database.NewNumberCondition(u.UpdatedAtColumn(), op, updatedAt)
}

// UsernameCondition implements [domain.HumanUserRepository].
func (u userHuman) UsernameCondition(op database.TextOperation, username string) database.Condition {
	return database.NewTextCondition(u.UsernameColumn(), op, username)
}

// UsernameOrgUniqueCondition implements [domain.HumanUserRepository].
func (u userHuman) UsernameOrgUniqueCondition(condition bool) database.Condition {
	return database.NewBooleanCondition(u.UsernameOrgUniqueColumn(), condition)
}

// // FirstNameCondition implements [domain.humanConditions].
// func (h userHuman) FirstNameCondition(op database.TextOperation, firstName string) database.Condition {
// 	return database.NewTextCondition(h.FirstNameColumn(), op, firstName)
// }

// // LastNameCondition implements [domain.humanConditions].
// func (h userHuman) LastNameCondition(op database.TextOperation, lastName string) database.Condition {
// 	return database.NewTextCondition(h.LastNameColumn(), op, lastName)
// }

// // EmailAddressCondition implements [domain.humanConditions].
// func (h userHuman) EmailAddressCondition(op database.TextOperation, email string) database.Condition {
// 	return database.NewTextCondition(h.EmailAddressColumn(), op, email)
// }

// // EmailVerifiedCondition implements [domain.humanConditions].
// func (h userHuman) EmailVerifiedCondition(isVerified bool) database.Condition {
// 	if isVerified {
// 		return database.IsNotNull(h.EmailVerifiedAtColumn())
// 	}
// 	return database.IsNull(h.EmailVerifiedAtColumn())
// }

// // EmailVerifiedAtCondition implements [domain.humanConditions].
// func (h userHuman) EmailVerifiedAtCondition(op database.NumberOperation, verifiedAt time.Time) database.Condition {
// 	return database.NewNumberCondition(h.EmailVerifiedAtColumn(), op, verifiedAt)
// }

// // PhoneNumberCondition implements [domain.humanConditions].
// func (h userHuman) PhoneNumberCondition(op database.TextOperation, phoneNumber string) database.Condition {
// 	return database.NewTextCondition(h.PhoneNumberColumn(), op, phoneNumber)
// }

// // PhoneVerifiedCondition implements [domain.humanConditions].
// func (h userHuman) PhoneVerifiedCondition(isVerified bool) database.Condition {
// 	if isVerified {
// 		return database.IsNotNull(h.PhoneVerifiedAtColumn())
// 	}
// 	return database.IsNull(h.PhoneVerifiedAtColumn())
// }

// // PhoneVerifiedAtCondition implements [domain.humanConditions].
// func (h userHuman) PhoneVerifiedAtCondition(op database.NumberOperation, verifiedAt time.Time) database.Condition {
// 	return database.NewNumberCondition(h.PhoneVerifiedAtColumn(), op, verifiedAt)
// }

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// PrimaryKeyColumns implements domain.HumanUserRepository.
func (h userHuman) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		h.InstanceIDColumn(),
		h.IDColumn(),
	}
}

// AvatarKeyColumn implements [domain.HumanUserRepository].
func (u userHuman) AvatarKeyColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "avatar_key")
}

// CreatedAtColumn implements [domain.HumanUserRepository].
func (u userHuman) CreatedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "created_at")
}

// DisplayNameColumn implements [domain.HumanUserRepository].
func (u userHuman) DisplayNameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "display_name")
}

// FirstNameColumn implements [domain.HumanUserRepository].
func (u userHuman) FirstNameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "first_name")
}

// GenderColumn implements [domain.HumanUserRepository].
func (u userHuman) GenderColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "gender")
}

// IDColumn implements [domain.HumanUserRepository].
func (u userHuman) IDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "id")
}

// InstanceIDColumn implements [domain.HumanUserRepository].
func (u userHuman) InstanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "instance_id")
}

// LastNameColumn implements [domain.HumanUserRepository].
func (u userHuman) LastNameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "last_name")
}

// OrgIDColumn implements [domain.HumanUserRepository].
func (u userHuman) OrgIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "organization_id")
}

// PreferredLanguageColumn implements [domain.HumanUserRepository].
func (u userHuman) PreferredLanguageColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "preferred_language")
}

// StateColumn implements [domain.HumanUserRepository].
func (u userHuman) StateColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "state")
}

// UpdatedAtColumn implements [domain.HumanUserRepository].
func (u userHuman) UpdatedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "updated_at")
}

// UsernameColumn implements [domain.HumanUserRepository].
func (u userHuman) UsernameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "username")
}

// UsernameOrgUniqueColumn implements [domain.HumanUserRepository].
func (u userHuman) UsernameOrgUniqueColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "username_org_unique")
}

// NicknameColumn implements [domain.HumanUserRepository].
func (h userHuman) NicknameColumn() database.Column {
	return database.NewColumn("user_humans", "nick_name")
}

// // FirstNameColumn implements [domain.humanColumns].
// func (h userHuman) FirstNameColumn() database.Column {
// 	return database.NewColumn("user_humans", "first_name")
// }

// // LastNameColumn implements [domain.humanColumns].
// func (h userHuman) LastNameColumn() database.Column {
// 	return database.NewColumn("user_humans", "last_name")
// }

// // EmailAddressColumn implements [domain.humanColumns].
// func (h userHuman) EmailAddressColumn() database.Column {
// 	return database.NewColumn("user_humans", "email_address")
// }

// // EmailVerifiedAtColumn implements [domain.humanColumns].
// func (h userHuman) EmailVerifiedAtColumn() database.Column {
// 	return database.NewColumn("user_humans", "email_verified_at")
// }

// // PhoneNumberColumn implements [domain.humanColumns].
// func (h userHuman) PhoneNumberColumn() database.Column {
// 	return database.NewColumn("user_humans", "phone_number")
// }

// // PhoneVerifiedAtColumn implements [domain.humanColumns].
// func (h userHuman) PhoneVerifiedAtColumn() database.Column {
// 	return database.NewColumn("user_humans", "phone_verified_at")
// }

// // func (h userHuman) columns() database.Columns {
// // 	return append(h.user.columns(),
// // 		h.FirstNameColumn(),
// // 		h.LastNameColumn(),
// // 		h.EmailAddressColumn(),
// // 		h.EmailVerifiedAtColumn(),
// // 		h.PhoneNumberColumn(),
// // 		h.PhoneVerifiedAtColumn(),
// // 	)
// // }

// // func (h userHuman) writeReturning(builder *database.StatementBuilder) {
// // 	builder.WriteString(" RETURNING ")
// // 	h.columns().Write(builder)
// // }
