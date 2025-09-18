package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

type userHuman struct {
	*user
}

var _ domain.HumanRepository = (*userHuman)(nil)

const userEmailQuery = `SELECT h.email_address, h.email_verified_at FROM user_humans h`

// GetEmail implements [domain.HumanRepository].
func (u *userHuman) GetEmail(ctx context.Context, condition database.Condition) (*domain.Email, error) {
	var email domain.Email

	builder := database.StatementBuilder{}
	builder.WriteString(userEmailQuery)
	writeCondition(&builder, condition)

	err := u.client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(
		&email.Address,
		&email.VerifiedAt,
	)
	if err != nil {
		return nil, err
	}
	return &email, nil
}

// Update implements [domain.HumanRepository].
func (h userHuman) Update(ctx context.Context, condition database.Condition, changes ...database.Change) error {
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE human_users SET `)
	database.Changes(changes).Write(&builder)
	writeCondition(&builder, condition)

	stmt := builder.String()

	_, err := h.client.Exec(ctx, stmt, builder.Args()...)
	return err
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetFirstName implements [domain.humanChanges].
func (h userHuman) SetFirstName(firstName string) database.Change {
	return database.NewChange(h.FirstNameColumn(), firstName)
}

// SetLastName implements [domain.humanChanges].
func (h userHuman) SetLastName(lastName string) database.Change {
	return database.NewChange(h.LastNameColumn(), lastName)
}

// SetEmail implements [domain.humanChanges].
func (h userHuman) SetEmail(address string, verified *time.Time) database.Change {
	return database.NewChanges(
		h.SetEmailAddress(address),
		database.NewChangePtr(h.EmailVerifiedAtColumn(), verified),
	)
}

// SetEmailAddress implements [domain.humanChanges].
func (h userHuman) SetEmailAddress(address string) database.Change {
	return database.NewChange(h.EmailAddressColumn(), address)
}

// SetEmailVerifiedAt implements [domain.humanChanges].
func (h userHuman) SetEmailVerifiedAt(at time.Time) database.Change {
	if at.IsZero() {
		return database.NewChange(h.EmailVerifiedAtColumn(), database.NowInstruction)
	}
	return database.NewChange(h.EmailVerifiedAtColumn(), at)
}

// SetPhone implements [domain.humanChanges].
func (h userHuman) SetPhone(number string, verifiedAt *time.Time) database.Change {
	return database.NewChanges(
		h.SetPhoneNumber(number),
		database.NewChangePtr(h.PhoneVerifiedAtColumn(), verifiedAt),
	)
}

// SetPhoneNumber implements [domain.humanChanges].
func (h userHuman) SetPhoneNumber(number string) database.Change {
	return database.NewChange(h.PhoneNumberColumn(), number)
}

// SetPhoneVerifiedAt implements [domain.humanChanges].
func (h userHuman) SetPhoneVerifiedAt(at time.Time) database.Change {
	if at.IsZero() {
		return database.NewChange(h.PhoneVerifiedAtColumn(), database.NowInstruction)
	}
	return database.NewChange(h.PhoneVerifiedAtColumn(), at)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// FirstNameCondition implements [domain.humanConditions].
func (h userHuman) FirstNameCondition(op database.TextOperation, firstName string) database.Condition {
	return database.NewTextCondition(h.FirstNameColumn(), op, firstName)
}

// LastNameCondition implements [domain.humanConditions].
func (h userHuman) LastNameCondition(op database.TextOperation, lastName string) database.Condition {
	return database.NewTextCondition(h.LastNameColumn(), op, lastName)
}

// EmailAddressCondition implements [domain.humanConditions].
func (h userHuman) EmailAddressCondition(op database.TextOperation, email string) database.Condition {
	return database.NewTextCondition(h.EmailAddressColumn(), op, email)
}

// EmailVerifiedCondition implements [domain.humanConditions].
func (h *userHuman) EmailVerifiedCondition(isVerified bool) database.Condition {
	if isVerified {
		return database.IsNotNull(h.EmailVerifiedAtColumn())
	}
	return database.IsNull(h.EmailVerifiedAtColumn())
}

// EmailVerifiedAtCondition implements [domain.humanConditions].
func (h userHuman) EmailVerifiedAtCondition(op database.NumberOperation, verifiedAt time.Time) database.Condition {
	return database.NewNumberCondition(h.EmailVerifiedAtColumn(), op, verifiedAt)
}

// PhoneNumberCondition implements [domain.humanConditions].
func (h userHuman) PhoneNumberCondition(op database.TextOperation, phoneNumber string) database.Condition {
	return database.NewTextCondition(h.PhoneNumberColumn(), op, phoneNumber)
}

// PhoneVerifiedCondition implements [domain.humanConditions].
func (h userHuman) PhoneVerifiedCondition(isVerified bool) database.Condition {
	if isVerified {
		return database.IsNotNull(h.PhoneVerifiedAtColumn())
	}
	return database.IsNull(h.PhoneVerifiedAtColumn())
}

// PhoneVerifiedAtCondition implements [domain.humanConditions].
func (h userHuman) PhoneVerifiedAtCondition(op database.NumberOperation, verifiedAt time.Time) database.Condition {
	return database.NewNumberCondition(h.PhoneVerifiedAtColumn(), op, verifiedAt)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// FirstNameColumn implements [domain.humanColumns].
func (h userHuman) FirstNameColumn() database.Column {
	return database.NewColumn("user_humans", "first_name")
}

// LastNameColumn implements [domain.humanColumns].
func (h userHuman) LastNameColumn() database.Column {
	return database.NewColumn("user_humans", "last_name")
}

// EmailAddressColumn implements [domain.humanColumns].
func (h userHuman) EmailAddressColumn() database.Column {
	return database.NewColumn("user_humans", "email_address")
}

// EmailVerifiedAtColumn implements [domain.humanColumns].
func (h userHuman) EmailVerifiedAtColumn() database.Column {
	return database.NewColumn("user_humans", "email_verified_at")
}

// PhoneNumberColumn implements [domain.humanColumns].
func (h userHuman) PhoneNumberColumn() database.Column {
	return database.NewColumn("user_humans", "phone_number")
}

// PhoneVerifiedAtColumn implements [domain.humanColumns].
func (h userHuman) PhoneVerifiedAtColumn() database.Column {
	return database.NewColumn("user_humans", "phone_verified_at")
}

// func (h userHuman) columns() database.Columns {
// 	return append(h.user.columns(),
// 		h.FirstNameColumn(),
// 		h.LastNameColumn(),
// 		h.EmailAddressColumn(),
// 		h.EmailVerifiedAtColumn(),
// 		h.PhoneNumberColumn(),
// 		h.PhoneVerifiedAtColumn(),
// 	)
// }

// func (h userHuman) writeReturning(builder *database.StatementBuilder) {
// 	builder.WriteString(" RETURNING ")
// 	h.columns().Write(builder)
// }
