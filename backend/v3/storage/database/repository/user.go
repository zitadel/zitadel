package repository

import (
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.UserRepository = (*user)(nil)

type user struct{}

func UserRepository() domain.UserRepository {
	return new(user)
}

// TODO rename
func (o user) qualifiedTableName() string {
	return "zitadel.users"
}

// TODO rename
func (o user) unqualifiedTableName() string {
	return "users"
}

func (o user) qualifiedMachineUsersTableName() string {
	return "zitadel.machine_users"
}

func (o user) unqualifiedMachineUsersTableName() string {
	return "machine_users"
}

func (o user) qualifiedHumanUsersTableName() string {
	return "zitadel.human_users"
}

func (o user) unqualifiedHumanUsersTableName() string {
	return "human_users"
}

// const queryUserStmt = `SELECT instance_id, org_id, id, username, type, created_at, updated_at, deleted_at,` +
// 	` first_name, last_name, email_address, email_verified_at, phone_number, phone_verified_at, description` +
// 	` FROM users_view users`

// type user struct{}

// func UserRepository() domain.UserRepository {
// 	return new(user)
// }

// var _ domain.UserRepository = (*user)(nil)

// // -------------------------------------------------------------
// // repository
// // -------------------------------------------------------------

// // Human implements [domain.UserRepository].
// func (u user) Human() domain.HumanRepository {
// 	return &userHuman{user: u}
// }

// // Machine implements [domain.UserRepository].
// func (u user) Machine() domain.MachineRepository {
// 	return &userMachine{user: u}
// }

// // List implements [domain.UserRepository].
// func (u user) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (users []*domain.User, err error) {
// 	options := new(database.QueryOpts)
// 	for _, opt := range opts {
// 		opt(options)
// 	}

// 	builder := database.StatementBuilder{}
// 	builder.WriteString(queryUserStmt)
// 	options.WriteCondition(&builder)
// 	options.WriteOrderBy(&builder)
// 	options.WriteLimit(&builder)
// 	options.WriteOffset(&builder)

// 	rows, err := client.Query(ctx, builder.String(), builder.Args()...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		closeErr := rows.Close()
// 		if err != nil {
// 			return
// 		}
// 		err = closeErr
// 	}()
// 	for rows.Next() {
// 		user, err := scanUser(rows)
// 		if err != nil {
// 			return nil, err
// 		}
// 		users = append(users, user)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return users, nil
// }

// // Get implements [domain.UserRepository].
// func (u user) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.User, error) {
// 	options := new(database.QueryOpts)
// 	for _, opt := range opts {
// 		opt(options)
// 	}

// 	builder := database.StatementBuilder{}
// 	builder.WriteString(queryUserStmt)
// 	options.WriteCondition(&builder)
// 	options.WriteOrderBy(&builder)
// 	options.WriteLimit(&builder)
// 	options.WriteOffset(&builder)

// 	return scanUser(client.QueryRow(ctx, builder.String(), builder.Args()...))
// }

// const (
// 	createHumanStmt = `INSERT INTO human_users (instance_id, org_id, user_id, username, first_name, last_name, email_address, email_verified_at, phone_number, phone_verified_at)` +
// 		` VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)` +
// 		` RETURNING created_at, updated_at`
// 	createMachineStmt = `INSERT INTO user_machines (instance_id, org_id, user_id, username, description)` +
// 		` VALUES ($1, $2, $3, $4, $5)` +
// 		` RETURNING created_at, updated_at`
// )

// // Create implements [domain.UserRepository].
// func (u user) Create(ctx context.Context, client database.QueryExecutor, user *domain.User) error {
// 	builder := database.StatementBuilder{}
// 	builder.AppendArgs(user.InstanceID, user.OrgID, user.ID, user.Username, user.Traits.Type())
// 	switch trait := user.Traits.(type) {
// 	case *domain.Human:
// 		builder.WriteString(createHumanStmt)
// 		builder.AppendArgs(trait.FirstName, trait.LastName, trait.Email.Address, trait.Email.VerifiedAt, trait.Phone.Number, trait.Phone.VerifiedAt)
// 	case *domain.Machine:
// 		builder.WriteString(createMachineStmt)
// 		builder.AppendArgs(trait.Description)
// 	}
// 	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&user.CreatedAt, &user.UpdatedAt)
// }

// // Delete implements [domain.UserRepository].
// func (u user) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) error {
// 	builder := database.StatementBuilder{}
// 	builder.WriteString("DELETE FROM users")
// 	writeCondition(&builder, condition)
// 	_, err := client.Exec(ctx, builder.String(), builder.Args()...)
// 	return err
// }

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// user
func (u user) InstanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "instance_id")
}

func (u user) OrgIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "org_id")
}

func (u user) IDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "id")
}

func (u user) UsernameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "username")
}

func (u user) UsernameOrgUniqueColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "username_org_unique")
}

func (u user) StateColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "state")
}

func (u user) CreatedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "created_at")
}

func (u user) UpdatedAtColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "updated_at")
}

// machine
func (u user) NameColumn() database.Column {
	return database.NewColumn(u.unqualifiedMachineUsersTableName(), "name")
}

func (u user) DescriptionColumn() database.Column {
	return database.NewColumn(u.unqualifiedMachineUsersTableName(), "description")
}

// human
func (u user) FirstNameColumn() database.Column {
	return database.NewColumn(u.unqualifiedHumanUsersTableName(), "first_name")
}

func (u user) LastNameColumn() database.Column {
	return database.NewColumn(u.unqualifiedHumanUsersTableName(), "last_name")
}

func (u user) NickNameColumn() database.Column {
	return database.NewColumn(u.unqualifiedHumanUsersTableName(), "nick_name")
}

func (u user) DisplayNameColumn() database.Column {
	return database.NewColumn(u.unqualifiedHumanUsersTableName(), "display_name")
}

func (u user) PreferredLanguageColumn() database.Column {
	return database.NewColumn(u.unqualifiedHumanUsersTableName(), "preferred_language")
}

func (u user) GenderColumn() database.Column {
	return database.NewColumn(u.unqualifiedHumanUsersTableName(), "gender")
}

func (u user) AvatarKeyColumn() database.Column {
	return database.NewColumn(u.unqualifiedHumanUsersTableName(), "avatar_key")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (u user) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(u.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

func (u user) OrgIDCondition(orgID string) database.Condition {
	return database.NewTextCondition(u.OrgIDColumn(), database.TextOperationEqual, orgID)
}

func (u user) IDCondition(userID string) database.Condition {
	return database.NewTextCondition(u.IDColumn(), database.TextOperationEqual, userID)
}

func (u user) UsernameCondition(op database.TextOperation, username string) database.Condition {
	return database.NewTextCondition(u.UsernameColumn(), op, username)
}

func (u user) UsernameOrgUniqueCondition(condition bool) database.Condition {
	return database.NewBooleanCondition(u.UsernameColumn(), condition)
}

func (u user) StateCondition(state domain.UserState) database.Condition {
	return database.NewTextCondition(u.UsernameColumn(), database.TextOperationEqual, state.String())
}

func (u user) CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition {
	return database.NewNumberCondition(u.CreatedAtColumn(), op, createdAt)
}

func (u user) UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition {
	return database.NewNumberCondition(u.UpdatedAtColumn(), op, updatedAt)
}

// machine
func (u user) NameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(u.DescriptionColumn(), op, name)
}

func (u user) DescriptionCondition(op database.TextOperation, description string) database.Condition {
	return database.NewTextCondition(u.DescriptionColumn(), op, description)
}

// human
func (u user) FirstNameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(u.FirstNameColumn(), op, name)
}

func (u user) LastNameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(u.LastNameColumn(), op, name)
}

func (u user) NickNameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(u.NickNameColumn(), op, name)
}

func (u user) DisplayNameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(u.DisplayNameColumn(), op, name)
}

func (u user) PreferredLanguageCondition(language string) database.Condition {
	return database.NewTextCondition(u.PreferredLanguageColumn(), database.TextOperationEqual, language)
}

func (u user) GenderCondition(gender uint8) database.Condition {
	return database.NewNumberCondition(u.GenderColumn(), database.NumberOperationEqual, gender)
}

func (u user) AvatarKeyCondition(key string) database.Condition {
	return database.NewTextCondition(u.AvatarKeyColumn(), database.TextOperationEqual, key)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

func (u user) SetUsername(username string) database.Change {
	return database.NewChange(u.UsernameColumn(), username)
}

func (u user) SetUsernameOrgUnique(op bool) database.Change {
	return database.NewChange(u.UsernameOrgUniqueColumn(), op)
}

func (u user) SetState(state domain.UserState) database.Change {
	return database.NewChange(u.StateColumn(), state.String())
}

// machine
func (u user) SetName(name string) database.Change {
	return database.NewChange(u.NameColumn(), name)
}

func (u user) SetDescription(description string) database.Change {
	return database.NewChange(u.DescriptionColumn(), description)
}

// human
func (u user) SetFirstName(name string) database.Change {
	return database.NewChange(u.FirstNameColumn(), name)
}

func (u user) SetLastName(name string) database.Change {
	return database.NewChange(u.LastNameColumn(), name)
}

func (u user) SetNickName(name string) database.Change {
	return database.NewChange(u.NickNameColumn(), name)
}

func (u user) SetDisplayName(name string) database.Change {
	return database.NewChange(u.DisplayNameColumn(), name)
}

func (u user) SetPreferredLanguage(language string) database.Change {
	return database.NewChange(u.PreferredLanguageColumn(), language)
}

func (u user) SetGender(gender uint8) database.Change {
	return database.NewChange(u.GenderColumn(), gender)
}

func (u user) SetAvatarKey(key string) database.Change {
	return database.NewChange(u.AvatarKeyColumn(), key)
}
