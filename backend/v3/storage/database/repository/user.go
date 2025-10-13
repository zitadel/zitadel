package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var (
	_ domain.UserRepository    = (*user)(nil)
	_ domain.HumanRepository   = (*human)(nil)
	_ domain.MachineRepository = (*machine)(nil)
)

type (
	contact struct{}
	human   struct {
		contact
	}
	machine struct{}
	user    struct {
		human
		machine
	}
)

func UserRepository() domain.UserRepository {
	return &user{
		// human:   human{},
		// machine: machine{},
	}
}

func (m machine) qualifiedTableName() string {
	return "zitadel.machine_users"
}

func (m machine) unqualifiedTableName() string {
	return "machine_users"
}

func (h human) qualifiedTableName() string {
	return "zitadel.human_users"
}

func (h human) unqualifiedTableName() string {
	return "human_users"
}

func (c contact) qualifiedTableName() string {
	return "zitadel.human_contacts"
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// machine
func (m machine) InstanceIDColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "instance_id")
}

func (m machine) OrgIDColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "org_id")
}

func (m machine) IDColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "id")
}

func (m machine) UsernameColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "username")
}

func (m machine) UsernameOrgUniqueColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "username_org_unique")
}

func (m machine) StateColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "state")
}

func (m machine) CreatedAtColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "created_at")
}

func (m machine) UpdatedAtColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "updated_at")
}

func (m machine) NameColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "name")
}

func (m machine) DescriptionColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "description")
}

func (m machine) SecretColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "secret")
}

func (m machine) AccessTokenTypeColumn() database.Column {
	return database.NewColumn(m.unqualifiedTableName(), "access_token_type")
}

// human
func (h human) InstanceIDColumn() database.Column {
	return database.NewColumn(h.unqualifiedTableName(), "instance_id")
}

func (h human) OrgIDColumn() database.Column {
	return database.NewColumn(h.unqualifiedTableName(), "org_id")
}

func (h human) IDColumn() database.Column {
	return database.NewColumn(h.unqualifiedTableName(), "id")
}

func (h human) UsernameColumn() database.Column {
	return database.NewColumn(h.unqualifiedTableName(), "username")
}

func (h human) UsernameOrgUniqueColumn() database.Column {
	return database.NewColumn(h.unqualifiedTableName(), "username_org_unique")
}

func (h human) StateColumn() database.Column {
	return database.NewColumn(h.unqualifiedTableName(), "state")
}

func (h human) CreatedAtColumn() database.Column {
	return database.NewColumn(h.unqualifiedTableName(), "created_at")
}

func (h human) UpdatedAtColumn() database.Column {
	return database.NewColumn(h.unqualifiedTableName(), "updated_at")
}

func (h human) FirstNameColumn() database.Column {
	return database.NewColumn(h.unqualifiedTableName(), "first_name")
}

func (h human) LastNameColumn() database.Column {
	return database.NewColumn(h.unqualifiedTableName(), "last_name")
}

func (h human) NickNameColumn() database.Column {
	return database.NewColumn(h.unqualifiedTableName(), "nick_name")
}

func (h human) DisplayNameColumn() database.Column {
	return database.NewColumn(h.unqualifiedTableName(), "display_name")
}

func (h human) PreferredLanguageColumn() database.Column {
	return database.NewColumn(h.unqualifiedTableName(), "preferred_language")
}

func (h human) GenderColumn() database.Column {
	return database.NewColumn(h.unqualifiedTableName(), "gender")
}

func (h human) AvatarKeyColumn() database.Column {
	return database.NewColumn(h.unqualifiedTableName(), "avatar_key")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// machine
func (m machine) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(m.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

func (m machine) OrgIDCondition(orgID string) database.Condition {
	return database.NewTextCondition(m.OrgIDColumn(), database.TextOperationEqual, orgID)
}

func (m machine) IDCondition(userID string) database.Condition {
	return database.NewTextCondition(m.IDColumn(), database.TextOperationEqual, userID)
}

func (m machine) UsernameCondition(op database.TextOperation, username string) database.Condition {
	return database.NewTextCondition(m.UsernameColumn(), op, username)
}

func (m machine) UsernameOrgUniqueCondition(condition bool) database.Condition {
	return database.NewBooleanCondition(m.UsernameOrgUniqueColumn(), condition)
}

func (m machine) StateCondition(state domain.UserState) database.Condition {
	return database.NewTextCondition(m.StateColumn(), database.TextOperationEqual, state.String())
}

func (m machine) CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition {
	return database.NewNumberCondition(m.CreatedAtColumn(), op, createdAt)
}

func (m machine) UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition {
	return database.NewNumberCondition(m.UpdatedAtColumn(), op, updatedAt)
}

func (m machine) NameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(m.DescriptionColumn(), op, name)
}

func (m machine) DescriptionCondition(op database.TextOperation, description string) database.Condition {
	return database.NewTextCondition(m.DescriptionColumn(), op, description)
}

func (m machine) AccessTokenTypeCondition(accessTokenType domain.AccessTokenType) database.Condition {
	return database.NewNumberCondition(m.AccessTokenTypeColumn(), database.NumberOperationEqual, uint8(accessTokenType))
}

// human
func (h human) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(h.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

func (h human) OrgIDCondition(orgID string) database.Condition {
	return database.NewTextCondition(h.OrgIDColumn(), database.TextOperationEqual, orgID)
}

func (h human) IDCondition(userID string) database.Condition {
	return database.NewTextCondition(h.IDColumn(), database.TextOperationEqual, userID)
}

func (h human) UsernameCondition(op database.TextOperation, username string) database.Condition {
	return database.NewTextCondition(h.UsernameColumn(), op, username)
}

func (h human) UsernameOrgUniqueCondition(condition bool) database.Condition {
	return database.NewBooleanCondition(h.UsernameOrgUniqueColumn(), condition)
}

func (h human) StateCondition(state domain.UserState) database.Condition {
	return database.NewTextCondition(h.StateColumn(), database.TextOperationEqual, state.String())
}

func (h human) CreatedAtCondition(op database.NumberOperation, createdAt time.Time) database.Condition {
	return database.NewNumberCondition(h.CreatedAtColumn(), op, createdAt)
}

func (h human) UpdatedAtCondition(op database.NumberOperation, updatedAt time.Time) database.Condition {
	return database.NewNumberCondition(h.UpdatedAtColumn(), op, updatedAt)
}

func (h human) FirstNameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(h.FirstNameColumn(), op, name)
}

func (h human) LastNameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(h.LastNameColumn(), op, name)
}

func (h human) NickNameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(h.NickNameColumn(), op, name)
}

func (h human) DisplayNameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(h.DisplayNameColumn(), op, name)
}

func (h human) PreferredLanguageCondition(language string) database.Condition {
	return database.NewTextCondition(h.PreferredLanguageColumn(), database.TextOperationEqual, language)
}

func (h human) GenderCondition(gender uint8) database.Condition {
	return database.NewNumberCondition(h.GenderColumn(), database.NumberOperationEqual, gender)
}

func (h human) AvatarKeyCondition(key string) database.Condition {
	return database.NewTextCondition(h.AvatarKeyColumn(), database.TextOperationEqual, key)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// machine
func (m machine) SetUsername(username string) database.Change {
	return database.NewChange(m.UsernameColumn(), username)
}

func (m machine) SetUsernameOrgUnique(op bool) database.Change {
	return database.NewChange(m.UsernameOrgUniqueColumn(), op)
}

func (m machine) SetState(state domain.UserState) database.Change {
	return database.NewChange(m.StateColumn(), state.String())
}

func (m machine) SetName(name string) database.Change {
	return database.NewChange(m.NameColumn(), name)
}

func (m machine) SetDescription(description string) database.Change {
	return database.NewChange(m.DescriptionColumn(), description)
}

func (m machine) SetUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(m.UpdatedAtColumn(), updatedAt)
}

func (m machine) SetSecret(secret *string) database.Change {
	return database.NewChangePtr(m.SecretColumn(), secret)
}

func (m machine) SetAccessTokenType(accessTokenType domain.AccessTokenType) database.Change {
	return database.NewChange(m.AccessTokenTypeColumn(), uint8(accessTokenType))
}

// human
func (h human) SetUsername(username string) database.Change {
	return database.NewChange(h.UsernameColumn(), username)
}

func (h human) SetUsernameOrgUnique(op bool) database.Change {
	return database.NewChange(h.UsernameOrgUniqueColumn(), op)
}

func (h human) SetState(state domain.UserState) database.Change {
	return database.NewChange(h.StateColumn(), state.String())
}

func (h human) SetFirstName(name string) database.Change {
	return database.NewChange(h.FirstNameColumn(), name)
}

func (h human) SetLastName(name string) database.Change {
	return database.NewChange(h.LastNameColumn(), name)
}

func (h human) SetNickName(name string) database.Change {
	return database.NewChange(h.NickNameColumn(), name)
}

func (h human) SetDisplayName(name string) database.Change {
	return database.NewChange(h.DisplayNameColumn(), name)
}

func (h human) SetPreferredLanguage(language string) database.Change {
	return database.NewChange(h.PreferredLanguageColumn(), language)
}

func (h human) SetGender(gender uint8) database.Change {
	return database.NewChange(h.GenderColumn(), gender)
}

func (h human) SetAvatarKey(key *string) database.Change {
	return database.NewChangePtr(h.AvatarKeyColumn(), key)
}

func (h human) SetUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(h.UpdatedAtColumn(), updatedAt)
}

func (u user) Human() domain.HumanRepository {
	return &human{}
}

func (u user) Machine() domain.MachineRepository {
	return &machine{}
}

// Create Human could have been done in one statement using CTE(s), but because a user may or may not email + phone or just email, this would require more code to handle
// place holder numbering, so I decided to use separate statements
func (u user) CreateHuman(ctx context.Context, client database.QueryExecutor, user *domain.Human) (*domain.Human, error) {
	// if user.HumanContact != nil {
	// }
	// if user.HumanSecurity != nil {
	// }

	user, err := u.createHuman(ctx, client, user)
	if err != nil {
		return nil, err
	}

	humanEmailContactCreateErrChann := make(chan error, 1)
	go func() {
		humanEmailContactCreateErrChann <- u.CreateHumanContact(ctx, client, user, &user.HumanEmailContact)
	}()

	var humanPhoneContactCreateErrChann chan error
	if user.HumanPhoneContact != nil {
		humanPhoneContactCreateErrChann = make(chan error, 1)
		go func() {
			humanPhoneContactCreateErrChann <- u.CreateHumanContact(ctx, client, user, user.HumanPhoneContact)
		}()
	}

	err = u.CreateHumanSecurity(ctx, client, user)
	if err != nil {
		return nil, err
	}

	if err := <-humanEmailContactCreateErrChann; err != nil {
		return nil, err
	}

	if humanPhoneContactCreateErrChann != nil {
		if err := <-humanPhoneContactCreateErrChann; err != nil {
			return nil, err
		}
	}

	return user, nil
}

const createHumaneStmt = `INSERT INTO zitadel.human_users (instance_id, org_id, id, username, username_org_unique, state,` +
	` first_name, last_name, nick_name, display_name, preferred_language, gender, avatar_key)` +
	` VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)` +
	` RETURNING created_at, updated_at`

func (u user) createHuman(ctx context.Context, client database.QueryExecutor, user *domain.Human) (*domain.Human, error) {
	builder := database.StatementBuilder{}
	fmt.Printf("[DEBUGPRINT] [users_test.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> user.State = %+v\n", user.State)
	builder.AppendArgs(user.User.InstanceID, user.User.OrgID, user.ID, user.Username, user.UsernameOrgUnique, user.State)
	builder.AppendArgs(user.FirstName, user.LastName, user.NickName, user.DisplayName, user.PreferredLanguage, user.Gender, user.AvatarKey)

	builder.WriteString(createHumaneStmt)

	err := client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&user.CreatedAt, &user.UpdatedAt)
	return user, err
}

const createHumaneContactStmt = `INSERT INTO zitadel.human_contacts (instance_id, org_id, user_id,` +
	` type, value, is_verified, unverified_value)` +
	` VALUES($1, $2, $3, $4, $5, $6, $7)`

func (u user) CreateHumanContact(ctx context.Context, client database.QueryExecutor, user *domain.Human, contact *domain.HumanContact) error {
	builder := database.StatementBuilder{}
	builder.AppendArgs(user.User.InstanceID, user.User.OrgID, user.ID)
	builder.AppendArgs(contact.Type, contact.Value, contact.IsVerified, contact.UnverifiedValue)

	builder.WriteString(createHumaneContactStmt)

	// return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&user.CreatedAt, &user.UpdatedAt)
	_, err := client.Exec(ctx, builder.String(), builder.Args()...)
	return err
}

const createHumaneSecuirtyStmt = `INSERT INTO zitadel.human_security (instance_id, org_id, user_id,` +
	` password_change_required, password_changed, mfa_init_skipped)` +
	` VALUES($1, $2, $3, $4, $5, $6)`

func (u user) CreateHumanSecurity(ctx context.Context, client database.QueryExecutor, user *domain.Human) error {
	builder := database.StatementBuilder{}
	builder.AppendArgs(user.User.InstanceID, user.User.OrgID, user.ID)
	builder.AppendArgs(user.HumanSecurity.PasswordChangeRequired, user.HumanSecurity.PasswordChange, user.HumanSecurity.MFAInitSkipped)

	builder.WriteString(createHumaneSecuirtyStmt)

	// return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&user.CreatedAt, &user.UpdatedAt)
	_, err := client.Exec(ctx, builder.String(), builder.Args()...)
	return err
}

func (u user) UpdateHuman(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}

	if !condition.IsRestrictingColumn(u.human.InstanceIDColumn()) {
		return 0, database.NewMissingConditionError(u.human.InstanceIDColumn())
	}

	// if !condition.IsRestrictingColumn(u.human.OrgIDColumn()) {
	// 	return 0, database.NewMissingConditionError(u.human.OrgIDColumn())
	// }

	if !condition.IsRestrictingColumn(u.human.IDColumn()) {
		return 0, database.NewMissingConditionError(u.human.IDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(`UPDATE zitadel.human_users SET `)
	database.Changes(changes).Write(&builder)
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// func (h human) GetHumanTableChanges(changes ...database.Change) []*database.Change {
// 	var humanChanges []*database.Change
// 	for _, change := range changes {
// 		if change.IsOnTable(h.contact.qualifiedTableName()) {
// 			if humanChanges == nil {
// 				humanChanges = make([]*database.Change, 0, len(changes))
// 			}
// 			humanChanges = append(humanChanges, &change)
// 		}
// 	}
// 	return humanChanges
// }

func (c contact) Set(ctx context.Context, client database.QueryExecutor, contact domain.HumanContact) (int64, error) {
	// 	for _, opt := range opts {
	// 		opt(options)
	return 0, nil
}

const queryHumanUserStmt = `SELECT zitadel.human_users.instance_id, zitadel.human_users.org_id, id, username, username_org_unique, state,` +
	` first_name, last_name, nick_name, display_name, preferred_language, gender, avatar_key,` +
	` created_at, updated_at,` +
	// email
	` email.type AS "email.type", email.value AS "email.value", email.is_verified AS "email.is_verified", email.unverified_value AS "email.unverified_value",` +
	// phone
	` phone.type AS "phone.type", phone.value AS "phone.value", phone.is_verified AS "phone.is_verified", phone.unverified_value AS "phone.unverified_value"` +
	` FROM zitadel.human_users` +
	// join email
	` LEFT JOIN zitadel.human_contacts AS email ON zitadel.human_users.id = email.user_id` +
	` AND zitadel.human_users.instance_id = email.instance_id` +
	` AND zitadel.human_users.org_id = email.org_id` +
	` AND email.type = 'email'` +
	// join phone
	` LEFT JOIN zitadel.human_contacts AS phone ON zitadel.human_users.id = phone.user_id` +
	` AND zitadel.human_users.instance_id = phone.instance_id` +
	` AND zitadel.human_users.org_id = phone.org_id` +
	` AND phone.type = 'phone'`

func (u user) GetHuman(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Human, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	if !options.Condition.IsRestrictingColumn(u.human.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(u.human.InstanceIDColumn())
	}

	// if !options.Condition.IsRestrictingColumn(u.human.OrgIDColumn()) {
	// 	return nil, database.NewMissingConditionError(u.human.OrgIDColumn())
	// }

	// if !options.Condition.IsRestrictingColumn(u.human.IDColumn()) {
	// 	return nil, database.NewMissingConditionError(u.human.IDColumn())
	// }

	var builder database.StatementBuilder
	builder.WriteString(queryHumanUserStmt)
	options.Write(&builder)

	user, err := scanHuman(ctx, client, &builder)
	if err != nil {
		return nil, err
	}

	if user.HumanPhoneContact != nil && user.HumanPhoneContact.Value == nil {
		user.HumanPhoneContact = nil
	}

	return user, nil
}

func (u user) ListHuman(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Human, error) {
	builder := database.StatementBuilder{}

	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	if !options.Condition.IsRestrictingColumn(u.Human().InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(u.Human().InstanceIDColumn())
	}

	builder.WriteString(queryHumanUserStmt)
	options.Write(&builder)

	orderBy := database.OrderBy(u.Human().CreatedAtColumn())
	orderBy.Write(&builder)

	return scanHumans(ctx, client, &builder)
}

func (u user) DeleteHuman(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if !condition.IsRestrictingColumn(u.Human().InstanceIDColumn()) {
		return 0, database.NewMissingConditionError(u.Human().InstanceIDColumn())
	}

	// if !condition.IsRestrictingColumn(u.Human().OrgIDColumn()) {
	// 	return 0, database.NewMissingConditionError(u.Human().OrgIDColumn())
	// }

	if !condition.IsRestrictingColumn(u.Human().IDColumn()) {
		return 0, database.NewMissingConditionError(u.Human().IDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(`DELETE FROM zitadel.human_users`)
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

const createMachineStmt = `INSERT INTO zitadel.machine_users (instance_id, org_id, id, username, username_org_unique, state,` +
	` name, description)` +
	` VALUES($1, $2, $3, $4, $5, $6, $7, $8)` +
	` RETURNING created_at, updated_at`

func (u user) CreateMachine(ctx context.Context, client database.QueryExecutor, user *domain.Machine) (*domain.Machine, error) {
	builder := database.StatementBuilder{}
	builder.AppendArgs(user.InstanceID, user.OrgID, user.ID, user.Username, user.UsernameOrgUnique, user.State)
	builder.AppendArgs(user.Name, user.Description)

	builder.WriteString(createMachineStmt)

	err := client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func (u user) UpdateMachine(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}

	if !condition.IsRestrictingColumn(u.machine.InstanceIDColumn()) {
		return 0, database.NewMissingConditionError(u.machine.InstanceIDColumn())
	}

	if !condition.IsRestrictingColumn(u.machine.OrgIDColumn()) {
		return 0, database.NewMissingConditionError(u.machine.OrgIDColumn())
	}

	if !condition.IsRestrictingColumn(u.machine.IDColumn()) {
		return 0, database.NewMissingConditionError(u.machine.IDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(`UPDATE zitadel.machine_users SET `)
	database.Changes(changes).Write(&builder)
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

const querMachineUserStmt = `SELECT instance_id, org_id, id, username, username_org_unique, state,` +
	` name, description,` +
	` created_at, updated_at` +
	` FROM zitadel.machine_users`

func (u user) GetMachine(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Machine, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	options = new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	if !options.Condition.IsRestrictingColumn(u.machine.InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(u.machine.InstanceIDColumn())
	}

	if !options.Condition.IsRestrictingColumn(u.machine.OrgIDColumn()) {
		return nil, database.NewMissingConditionError(u.machine.OrgIDColumn())
	}

	if !options.Condition.IsRestrictingColumn(u.machine.IDColumn()) {
		return nil, database.NewMissingConditionError(u.machine.IDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(querMachineUserStmt)
	options.Write(&builder)

	return scanMachine(ctx, client, &builder)
}

func (u user) ListMachine(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Machine, error) {
	builder := database.StatementBuilder{}

	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	if !options.Condition.IsRestrictingColumn(u.Machine().InstanceIDColumn()) {
		return nil, database.NewMissingConditionError(u.Machine().InstanceIDColumn())
	}

	builder.WriteString(querMachineUserStmt)
	options.Write(&builder)

	orderBy := database.OrderBy(u.Machine().CreatedAtColumn())
	orderBy.Write(&builder)

	return scanMachines(ctx, client, &builder)
}

func (u user) DeleteMachine(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if !condition.IsRestrictingColumn(u.Machine().InstanceIDColumn()) {
		return 0, database.NewMissingConditionError(u.Machine().InstanceIDColumn())
	}

	if !condition.IsRestrictingColumn(u.Machine().OrgIDColumn()) {
		return 0, database.NewMissingConditionError(u.Machine().OrgIDColumn())
	}

	if !condition.IsRestrictingColumn(u.Machine().IDColumn()) {
		return 0, database.NewMissingConditionError(u.Machine().IDColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(`DELETE FROM zitadel.machine_users`)
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func scanMachine(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.Machine, error) {
	user := &domain.Machine{}
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	err = rows.(database.CollectableRows).CollectExactlyOneRow(user)
	if err != nil {
		return nil, err
	}

	return user, err
}

func scanMachines(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.Machine, error) {
	users := []*domain.Machine{}

	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	err = rows.(database.CollectableRows).Collect(&users)
	if err != nil {
		return nil, err
	}

	return users, err
}

func scanHuman(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.Human, error) {
	user := &domain.Human{}
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	err = rows.(database.CollectableRows).CollectExactlyOneRow(user)
	if err != nil {
		return nil, err
	}

	return user, err
}

func scanHumans(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.Human, error) {
	users := []*domain.Human{}

	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	err = rows.(database.CollectableRows).Collect(&users)
	if err != nil {
		return nil, err
	}

	return users, err
}
