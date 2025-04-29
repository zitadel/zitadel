package v4

import (
	"context"
	"time"
)

type Human struct {
	FirstName string
	LastName  string
	Email     *Email
	Phone     *Phone
}

const UserTypeHuman UserType = "human"

func (Human) userTrait() {}

func (h Human) Type() UserType {
	return UserTypeHuman
}

var _ userTrait = (*Human)(nil)

type Email struct {
	Address string
	Verification
}

type Phone struct {
	Number string
	Verification
}

type Verification struct {
	VerifiedAt time.Time
}

type userHuman struct {
	*user
}

func (u *user) Human() *userHuman {
	return &userHuman{user: u}
}

const userEmailQuery = `SELECT h.email_address, h.email_verified_at FROM user_humans h`

func (u *userHuman) GetEmail(ctx context.Context) (*Email, error) {
	var email Email

	u.builder.WriteString(userEmailQuery)
	u.writeCondition()

	err := u.client.QueryRow(ctx, u.builder.String(), u.builder.args...).Scan(
		&email.Address,
		&email.Verification.VerifiedAt,
	)

	if err != nil {
		return nil, err
	}
	return &email, nil
}

func (h userHuman) Update(ctx context.Context, changes ...Change) error {
	h.builder.WriteString(`UPDATE human_users h SET `)
	Changes(changes).writeTo(&h.builder)
	h.writeCondition()

	stmt := h.builder.String()

	return h.client.Exec(ctx, stmt, h.builder.args...)
}

func (h userHuman) SetFirstName(firstName string) Change {
	return newChange(h.FirstNameColumn(), firstName)
}

func (h userHuman) FirstNameColumn() Column {
	return column{"h.first_name"}
}

func (h userHuman) FirstNameCondition(op TextOperator, firstName string) Condition {
	return newTextCondition(h.FirstNameColumn(), op, firstName)
}

func (h userHuman) SetLastName(lastName string) Change {
	return newChange(h.LastNameColumn(), lastName)
}

func (h userHuman) LastNameColumn() Column {
	return column{"h.last_name"}
}

func (h userHuman) LastNameCondition(op TextOperator, lastName string) Condition {
	return newTextCondition(h.LastNameColumn(), op, lastName)
}

func (h userHuman) EmailAddressColumn() Column {
	return ignoreCaseCol{
		column: column{"h.email_address"},
		suffix: "_lower",
	}
}

func (h userHuman) EmailAddressCondition(op TextOperator, email string) Condition {
	return newTextCondition(h.EmailAddressColumn(), op, email)
}

func (h userHuman) EmailVerifiedAtColumn() Column {
	return column{"h.email_verified_at"}
}

func (h *userHuman) EmailAddressVerifiedCondition(isVerified bool) Condition {
	if isVerified {
		return IsNotNull(h.EmailVerifiedAtColumn())
	}
	return IsNull(h.EmailVerifiedAtColumn())
}

func (h userHuman) EmailVerifiedAtCondition(op TextOperator, emailVerifiedAt string) Condition {
	return newTextCondition(h.EmailVerifiedAtColumn(), op, emailVerifiedAt)
}

func (h userHuman) SetEmailAddress(address string) Change {
	return newChange(h.EmailAddressColumn(), address)
}

// SetEmailVerified sets the verified column of the email
// if at is zero the statement uses the database timestamp
func (h userHuman) SetEmailVerified(at time.Time) Change {
	if at.IsZero() {
		return newChange(h.EmailVerifiedAtColumn(), nowDBInstruction)
	}
	return newChange(h.EmailVerifiedAtColumn(), at)
}

func (h userHuman) SetEmail(address string, verified *time.Time) Change {
	return newChanges(
		h.SetEmailAddress(address),
		newUpdatePtrColumn(h.EmailVerifiedAtColumn(), verified),
	)
}

func (h userHuman) PhoneNumberColumn() Column {
	return column{"h.phone_number"}
}

func (h userHuman) SetPhoneNumber(number string) Change {
	return newChange(h.PhoneNumberColumn(), number)
}

func (h userHuman) PhoneNumberCondition(op TextOperator, phoneNumber string) Condition {
	return newTextCondition(h.PhoneNumberColumn(), op, phoneNumber)
}

func (h userHuman) PhoneVerifiedAtColumn() Column {
	return column{"h.phone_verified_at"}
}

func (h userHuman) PhoneNumberVerifiedCondition(isVerified bool) Condition {
	if isVerified {
		return IsNotNull(h.PhoneVerifiedAtColumn())
	}
	return IsNull(h.PhoneVerifiedAtColumn())
}

// SetPhoneVerified sets the verified column of the phone
// if at is zero the statement uses the database timestamp
func (h userHuman) SetPhoneVerified(at time.Time) Change {
	if at.IsZero() {
		return newChange(h.PhoneVerifiedAtColumn(), nowDBInstruction)
	}
	return newChange(h.PhoneVerifiedAtColumn(), at)
}

func (h userHuman) PhoneVerifiedAtCondition(op TextOperator, phoneVerifiedAt string) Condition {
	return newTextCondition(h.PhoneVerifiedAtColumn(), op, phoneVerifiedAt)
}

func (h userHuman) SetPhone(number string, verifiedAt *time.Time) Change {
	return newChanges(
		h.SetPhoneNumber(number),
		newUpdatePtrColumn(h.PhoneVerifiedAtColumn(), verifiedAt),
	)
}
