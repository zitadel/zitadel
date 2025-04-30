package v4

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Dates struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type User struct {
	InstanceID string
	OrgID      string
	ID         string
	Username   string
	Traits     userTrait
	Dates
}

type UserType string

type userTrait interface {
	userTrait()
	Type() UserType
}

const queryUserStmt = `SELECT instance_id, org_id, id, username, type, created_at, updated_at, deleted_at,` +
	` first_name, last_name, email_address, email_verified_at, phone_number, phone_verified_at, description` +
	` FROM users_view`

type user struct {
	builder statementBuilder
	client  database.QueryExecutor
}

func UserRepository(client database.QueryExecutor) *user {
	return &user{
		client: client,
	}
}

func (u *user) List(ctx context.Context, opts ...QueryOption) (users []*User, err error) {
	options := new(queryOpts)
	for _, opt := range opts {
		opt(options)
	}

	u.builder.WriteString(queryUserStmt)
	options.writeCondition(&u.builder)
	options.writeOrderBy(&u.builder)
	options.writeLimit(&u.builder)
	options.writeOffset(&u.builder)

	rows, err := u.client.Query(ctx, u.builder.String(), u.builder.args...)
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

func (u *user) Get(ctx context.Context, opts ...QueryOption) (*User, error) {
	options := new(queryOpts)
	for _, opt := range opts {
		opt(options)
	}

	u.builder.WriteString(queryUserStmt)
	options.writeCondition(&u.builder)
	options.writeOrderBy(&u.builder)
	options.writeLimit(&u.builder)
	options.writeOffset(&u.builder)

	return scanUser(u.client.QueryRow(ctx, u.builder.String(), u.builder.args...))
}

const (
	// TODO: change to separate statements and tables
	createUserCte = `WITH user AS (` +
		`INSERT INTO users (instance_id, org_id, id, username, type) VALUES ($1, $2, $3, $4, $5)` +
		` RETURNING *)`
	createHumanStmt = createUserCte + ` INSERT INTO user_humans h (instance_id, org_id, user_id, first_name, last_name, email_address, email_verified_at, phone_number, phone_verified_at)` +
		` SELECT u.instance_id, u.org_id, u.id, $6, $7, $8, $9, $10, $11` +
		` FROM user u` +
		` RETURNING u.created_at, u.updated_at, u.deleted_at`
	createMachineStmt = createUserCte + ` INSERT INTO user_machines (instance_id, org_id, user_id, description)` +
		` SELECT u.instance_id, u.org_id, u.id, $6` +
		` FROM user u` +
		` RETURNING u.created_at, u.updated_at`
)

func (u *user) Create(ctx context.Context, user *User) error {
	u.builder.appendArgs(user.InstanceID, user.OrgID, user.ID, user.Username, user.Traits.Type())
	switch trait := user.Traits.(type) {
	case *Human:
		u.builder.WriteString(createHumanStmt)
		u.builder.appendArgs(trait.FirstName, trait.LastName, trait.Email.Address, trait.Email.VerifiedAt, trait.Phone.Number, trait.Phone.VerifiedAt)
	case *Machine:
		u.builder.WriteString(createMachineStmt)
		u.builder.appendArgs(trait.Description)
	}
	return u.client.QueryRow(ctx, u.builder.String(), u.builder.args...).Scan(&user.Dates.CreatedAt, &user.Dates.UpdatedAt)
}

func (u *user) Update(ctx context.Context, condition Condition, changes ...Change) error {
	u.builder.WriteString("UPDATE users SET ")
	Changes(changes).writeTo(&u.builder)
	u.writeCondition(condition)
	return u.client.Exec(ctx, u.builder.String(), u.builder.args...)
}

func (u *user) Delete(ctx context.Context, condition Condition) error {
	u.builder.WriteString("DELETE FROM users")
	u.writeCondition(condition)
	return u.client.Exec(ctx, u.builder.String(), u.builder.args...)
}

func (u *user) InstanceIDColumn() Column {
	return column{name: "instance_id"}
}

func (u *user) InstanceIDCondition(instanceID string) Condition {
	return newTextCondition(u.InstanceIDColumn(), TextOperatorEqual, instanceID)
}

func (u *user) OrgIDColumn() Column {
	return column{name: "org_id"}
}

func (u *user) OrgIDCondition(orgID string) Condition {
	return newTextCondition(u.OrgIDColumn(), TextOperatorEqual, orgID)
}

func (u *user) IDColumn() Column {
	return column{name: "id"}
}

func (u *user) IDCondition(userID string) Condition {
	return newTextCondition(u.IDColumn(), TextOperatorEqual, userID)
}

func (u *user) UsernameColumn() Column {
	return ignoreCaseCol{
		column: column{name: "username"},
		suffix: "_lower",
	}
}

func (u user) SetUsername(username string) Change {
	return newChange(u.UsernameColumn(), username)
}

func (u *user) UsernameCondition(op TextOperator, username string) Condition {
	return newTextCondition(u.UsernameColumn(), op, username)
}

func (u *user) CreatedAtColumn() Column {
	return column{name: "created_at"}
}

func (u *user) CreatedAtCondition(op NumberOperator, createdAt time.Time) Condition {
	return newNumberCondition(u.CreatedAtColumn(), op, createdAt)
}

func (u *user) UpdatedAtColumn() Column {
	return column{name: "updated_at"}
}

func (u *user) UpdatedAtCondition(op NumberOperator, updatedAt time.Time) Condition {
	return newNumberCondition(u.UpdatedAtColumn(), op, updatedAt)
}

func (u *user) DeletedAtColumn() Column {
	return column{name: "deleted_at"}
}

func (u *user) DeletedCondition(isDeleted bool) Condition {
	if isDeleted {
		return IsNotNull(u.DeletedAtColumn())
	}
	return IsNull(u.DeletedAtColumn())
}

func (u *user) DeletedAtCondition(op NumberOperator, deletedAt time.Time) Condition {
	return newNumberCondition(u.DeletedAtColumn(), op, deletedAt)
}

func (u *user) writeCondition(condition Condition) {
	if condition == nil {
		return
	}
	u.builder.WriteString(" WHERE ")
	condition.writeTo(&u.builder)
}

func (u user) columns() Columns {
	return Columns{
		u.InstanceIDColumn(),
		u.OrgIDColumn(),
		u.IDColumn(),
		u.UsernameColumn(),
		u.CreatedAtColumn(),
		u.UpdatedAtColumn(),
		u.DeletedAtColumn(),
	}
}

func scanUser(scanner database.Scanner) (*User, error) {
	var (
		user    User
		human   Human
		email   Email
		phone   Phone
		machine Machine
		typ     UserType
	)
	err := scanner.Scan(
		&user.InstanceID,
		&user.OrgID,
		&user.ID,
		&user.Username,
		&typ,
		&user.Dates.CreatedAt,
		&user.Dates.UpdatedAt,
		&user.Dates.DeletedAt,
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
	case UserTypeHuman:
		if email.Address != "" {
			human.Email = &email
		}
		if phone.Number != "" {
			human.Phone = &phone
		}
		user.Traits = &human
	case UserTypeMachine:
		user.Traits = &machine
	}

	return &user, nil
}
