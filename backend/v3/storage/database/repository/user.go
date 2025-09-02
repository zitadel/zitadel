package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

const queryUserStmt = `SELECT instance_id, org_id, id, username, type, created_at, updated_at, deleted_at,` +
	` first_name, last_name, email_address, email_verified_at, phone_number, phone_verified_at, description` +
	` FROM users_view users`

type user struct {
	repository
}

func UserRepository(client database.QueryExecutor) domain.UserRepository {
	return &user{
		repository: repository{
			client: client,
		},
	}
}

var _ domain.UserRepository = (*user)(nil)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

// Human implements [domain.UserRepository].
func (u *user) Human() domain.HumanRepository {
	return &userHuman{user: u}
}

// Machine implements [domain.UserRepository].
func (u *user) Machine() domain.MachineRepository {
	return &userMachine{user: u}
}

// List implements [domain.UserRepository].
func (u *user) List(ctx context.Context, opts ...database.QueryOption) (users []*domain.User, err error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	builder := database.StatementBuilder{}
	builder.WriteString(queryUserStmt)
	options.WriteCondition(&builder)
	options.WriteOrderBy(&builder)
	options.WriteLimit(&builder)
	options.WriteOffset(&builder)

	rows, err := u.client.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := rows.Close()
		if err != nil {
			return
		}
		err = closeErr
	}()
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// Get implements [domain.UserRepository].
func (u *user) Get(ctx context.Context, opts ...database.QueryOption) (*domain.User, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	builder := database.StatementBuilder{}
	builder.WriteString(queryUserStmt)
	options.WriteCondition(&builder)
	options.WriteOrderBy(&builder)
	options.WriteLimit(&builder)
	options.WriteOffset(&builder)

	return scanUser(u.client.QueryRow(ctx, builder.String(), builder.Args()...))
}

const (
	createHumanStmt = `INSERT INTO human_users (instance_id, org_id, user_id, username, first_name, last_name, email_address, email_verified_at, phone_number, phone_verified_at)` +
		` VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)` +
		` RETURNING created_at, updated_at`
	createMachineStmt = `INSERT INTO user_machines (instance_id, org_id, user_id, username, description)` +
		` VALUES ($1, $2, $3, $4, $5)` +
		` RETURNING created_at, updated_at`
)

// Create implements [domain.UserRepository].
func (u *user) Create(ctx context.Context, user *domain.User) error {
	builder := database.StatementBuilder{}
	builder.AppendArgs(user.InstanceID, user.OrgID, user.ID, user.Username, user.Traits.Type())
	switch trait := user.Traits.(type) {
	case *domain.Human:
		builder.WriteString(createHumanStmt)
		builder.AppendArgs(trait.FirstName, trait.LastName, trait.Email.Address, trait.Email.VerifiedAt, trait.Phone.Number, trait.Phone.VerifiedAt)
	case *domain.Machine:
		builder.WriteString(createMachineStmt)
		builder.AppendArgs(trait.Description)
	}
	return u.client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&user.CreatedAt, &user.UpdatedAt)
}

// Delete implements [domain.UserRepository].
func (u *user) Delete(ctx context.Context, condition database.Condition) error {
	builder := database.StatementBuilder{}
	builder.WriteString("DELETE FROM users")
	writeCondition(&builder, condition)
	_, err := u.client.Exec(ctx, builder.String(), builder.Args()...)
	return err
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetUsername implements [domain.userChanges].
func (u user) SetUsername(username string) database.Change {
	return database.NewChange(u.UsernameColumn(), username)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// InstanceIDCondition implements [domain.userConditions].
func (u user) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(u.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// OrgIDCondition implements [domain.userConditions].
func (u user) OrgIDCondition(orgID string) database.Condition {
	return database.NewTextCondition(u.OrgIDColumn(), database.TextOperationEqual, orgID)
}

// IDCondition implements [domain.userConditions].
func (u user) IDCondition(userID string) database.Condition {
	return database.NewTextCondition(u.IDColumn(), database.TextOperationEqual, userID)
}

// UsernameCondition implements [domain.userConditions].
func (u user) UsernameCondition(op database.TextOperation, username string) database.Condition {
	return database.NewTextCondition(u.UsernameColumn(), op, username)
}

// CreatedAtCondition implements [domain.userConditions].
func (u user) CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition {
	return database.NewNumberCondition(u.CreatedAtColumn(), op, createdAt)
}

// UpdatedAtCondition implements [domain.userConditions].
func (u user) UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition {
	return database.NewNumberCondition(u.UpdatedAtColumn(), op, updatedAt)
}

// DeletedCondition implements [domain.userConditions].
func (u user) DeletedCondition(isDeleted bool) database.Condition {
	if isDeleted {
		return database.IsNotNull(u.DeletedAtColumn())
	}
	return database.IsNull(u.DeletedAtColumn())
}

// DeletedAtCondition implements [domain.userConditions].
func (u user) DeletedAtCondition(op database.NumberOperation, deletedAt time.Time) database.Condition {
	return database.NewNumberCondition(u.DeletedAtColumn(), op, deletedAt)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// InstanceIDColumn implements [domain.userColumns].
func (user) InstanceIDColumn() database.Column {
	return database.NewColumn("users", "instance_id")
}

// OrgIDColumn implements [domain.userColumns].
func (user) OrgIDColumn() database.Column {
	return database.NewColumn("users", "org_id")
}

// IDColumn implements [domain.userColumns].
func (user) IDColumn() database.Column {
	return database.NewColumn("users", "id")
}

// UsernameColumn implements [domain.userColumns].
func (user) UsernameColumn() database.Column {
	return database.NewColumn("users", "username")
}

// FirstNameColumn implements [domain.userColumns].
func (user) CreatedAtColumn() database.Column {
	return database.NewColumn("users", "created_at")
}

// UpdatedAtColumn implements [domain.userColumns].
func (user) UpdatedAtColumn() database.Column {
	return database.NewColumn("users", "updated_at")
}

// DeletedAtColumn implements [domain.userColumns].
func (user) DeletedAtColumn() database.Column {
	return database.NewColumn("users", "deleted_at")
}

func (u user) columns() database.Columns {
	return database.Columns{
		u.InstanceIDColumn(),
		u.OrgIDColumn(),
		u.IDColumn(),
		u.UsernameColumn(),
		u.CreatedAtColumn(),
		u.UpdatedAtColumn(),
		u.DeletedAtColumn(),
	}
}

func scanUser(scanner database.Scanner) (*domain.User, error) {
	var (
		user    domain.User
		human   domain.Human
		email   domain.Email
		phone   domain.Phone
		machine domain.Machine
		typ     domain.UserType
	)
	err := scanner.Scan(
		&user.InstanceID,
		&user.OrgID,
		&user.ID,
		&user.Username,
		&typ,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&human.FirstName,
		&human.LastName,
		&email.Address,
		&email.VerifiedAt,
		&phone.Number,
		&phone.VerifiedAt,
		&machine.Description,
	)
	if err != nil {
		return nil, err
	}

	switch typ {
	case domain.UserTypeHuman:
		if email.Address != "" {
			human.Email = &email
		}
		if phone.Number != "" {
			human.Phone = &phone
		}
		user.Traits = &human
	case domain.UserTypeMachine:
		user.Traits = &machine
	}

	return &user, nil
}
