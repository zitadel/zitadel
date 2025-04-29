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

const userQuery = `SELECT u.instance_id, u.org_id, u.id, u.username, u.type, u.created_at, u.updated_at, u.deleted_at,` +
	` h.first_name, h.last_name, h.email_address, h.email_verified_at, h.phone_number, h.phone_verified_at, m.description` +
	` FROM users u` +
	` LEFT JOIN user_humans h ON u.instance_id = h.instance_id AND u.org_id = h.org_id AND u.id = h.id` +
	` LEFT JOIN user_machines m ON u.instance_id = m.instance_id AND u.org_id = m.org_id AND u.id = m.id`

type user struct {
	builder statementBuilder
	client  database.QueryExecutor

	condition Condition
}

func UserRepository(client database.QueryExecutor) *user {
	return &user{
		client: client,
	}
}

func (u *user) WithCondition(condition Condition) *user {
	u.condition = condition
	return u
}

func (u *user) Get(ctx context.Context) (*User, error) {
	u.builder.WriteString(userQuery)
	u.writeCondition()
	return scanUser(u.client.QueryRow(ctx, u.builder.String(), u.builder.args...))
}

func (u *user) List(ctx context.Context) (users []*User, err error) {
	u.builder.WriteString(userQuery)
	u.writeCondition()

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

const (
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
	return u.client.QueryRow(ctx, u.builder.String(), u.builder.args...).Scan(user.CreatedAt, user.UpdatedAt)
}

func (u *user) InstanceIDColumn() Column {
	return column{name: "u.instance_id"}
}

func (u *user) InstanceIDCondition(instanceID string) Condition {
	return newTextCondition(u.InstanceIDColumn(), TextOperatorEqual, instanceID)
}

func (u *user) OrgIDColumn() Column {
	return column{name: "u.org_id"}
}

func (u *user) OrgIDCondition(orgID string) Condition {
	return newTextCondition(u.OrgIDColumn(), TextOperatorEqual, orgID)
}

func (u *user) IDColumn() Column {
	return column{name: "u.id"}
}

func (u *user) IDCondition(userID string) Condition {
	return newTextCondition(u.IDColumn(), TextOperatorEqual, userID)
}

func (u *user) UsernameColumn() Column {
	return ignoreCaseCol{
		column: column{name: "u.username"},
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
	return column{name: "u.created_at"}
}

func (u *user) CreatedAtCondition(op NumberOperator, createdAt time.Time) Condition {
	return newNumberCondition(u.CreatedAtColumn(), op, createdAt)
}

func (u *user) UpdatedAtColumn() Column {
	return column{name: "u.updated_at"}
}

func (u *user) UpdatedAtCondition(op NumberOperator, updatedAt time.Time) Condition {
	return newNumberCondition(u.UpdatedAtColumn(), op, updatedAt)
}

func (u *user) DeletedAtColumn() Column {
	return column{name: "u.deleted_at"}
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

func (u *user) writeCondition() {
	if u.condition == nil {
		return
	}
	u.builder.WriteString(" WHERE ")
	u.condition.writeTo(&u.builder)
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
