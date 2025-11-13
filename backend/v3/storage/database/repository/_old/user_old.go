package repository

// import (
// 	"context"
// 	"time"

// 	"github.com/zitadel/zitadel/backend/v3/domain"
// 	"github.com/zitadel/zitadel/backend/v3/storage/database"
// )

// var (
// 	_ domain.UserRepository          = (*user)(nil)
// 	_ domain.HumanSecurityRepository = (*security)(nil)
// )

// type (
// 	contact  struct{}
// 	security struct{}
// 	// human    struct {
// 	// 	contact
// 	// 	security
// 	// }
// 	user struct {
// 		human
// 		machine
// 	}
// )

// func UserRepository() domain.UserRepository {
// 	return &user{
// 		// human:   human{},
// 		// machine: machine{},
// 	}
// }

// func (s security) qualifiedTableName() string {
// 	return "zitadel.human_security"
// }

// func (s security) unqualifiedTableName() string {
// 	return "human_security"
// }

// func (c contact) qualifiedTableName() string {
// 	return "zitadel.human_contacts"
// }

// // -------------------------------------------------------------
// // columns
// // -------------------------------------------------------------

// // security
// func (s security) InstanceIDColumn() database.Column {
// 	return database.NewColumn(s.unqualifiedTableName(), "instance_id")
// }

// func (s security) OrgIDColumn() database.Column {
// 	return database.NewColumn(s.unqualifiedTableName(), "org_id")
// }

// func (s security) UserIDColumn() database.Column {
// 	return database.NewColumn(s.unqualifiedTableName(), "id")
// }

// func (s security) PasswordChangeRequiredColumn() database.Column {
// 	return database.NewColumn(s.unqualifiedTableName(), "password_change_required")
// }

// func (s security) PasswordChangedColumn() database.Column {
// 	return database.NewColumn(s.unqualifiedTableName(), "password_changed")
// }

// func (s security) MFAInitSkippedColumn() database.Column {
// 	return database.NewColumn(s.unqualifiedTableName(), "mfa_init_skipped")
// }

// // -------------------------------------------------------------
// // conditions
// // -------------------------------------------------------------

// // security
// func (s security) InstanceIDCondition(instanceID string) database.Condition {
// 	return database.NewTextCondition(s.InstanceIDColumn(), database.TextOperationEqual, instanceID)
// }

// func (s security) OrgIDCondition(orgID string) database.Condition {
// 	return database.NewTextCondition(s.OrgIDColumn(), database.TextOperationEqual, orgID)
// }

// func (s security) UserIDCondition(userID string) database.Condition {
// 	return database.NewTextCondition(s.UserIDColumn(), database.TextOperationEqual, userID)
// }

// func (s security) PassswordChangeRequiredCondition(required bool) database.Condition {
// 	return database.NewBooleanCondition(s.PasswordChangeRequiredColumn(), required)
// }

// func (s security) PasswordChangeCondition(op database.NumberOperation, time time.Time) database.Condition {
// 	return database.NewNumberCondition(s.PasswordChangedColumn(), op, time)
// }

// func (s security) MFAInitSkippedCondition(skipped bool) database.Condition {
// 	return database.NewBooleanCondition(s.PasswordChangeRequiredColumn(), skipped)
// }

// // -------------------------------------------------------------
// // changes
// // -------------------------------------------------------------

// // security
// func (s security) SetPasswordChangeRequired(required bool) database.Change {
// 	return database.NewChange(s.PasswordChangeRequiredColumn(), required)
// }

// func (s security) SetPasswordChanged(time time.Time) database.Change {
// 	return database.NewChange(s.PasswordChangedColumn(), time)
// }

// func (s security) SetMFAInitSkipped(skipped bool) database.Change {
// 	return database.NewChange(s.MFAInitSkippedColumn(), skipped)
// }

// // Create Human could have been done in one statement using CTE(s), but because a user may or may not email + phone or just email, this would require more code to handle
// // place holder numbering, so I decided to use separate statements
// func (u user) CreateHuman(ctx context.Context, client database.QueryExecutor, user *domain.Human) (*domain.Human, error) {
// 	// if user.HumanContact != nil {
// 	// }
// 	// if user.HumanSecurity != nil {
// 	// }

// 	user, err := u.createHuman(ctx, client, user)
// 	if err != nil {
// 		return nil, err
// 	}

// 	humanEmailContactCreateErrChann := make(chan error, 1)
// 	go func() {
// 		humanEmailContactCreateErrChann <- u.CreateHumanContact(ctx, client, user, &user.HumanEmailContact)
// 	}()

// 	var humanPhoneContactCreateErrChann chan error
// 	if user.HumanPhoneContact != nil {
// 		humanPhoneContactCreateErrChann = make(chan error, 1)
// 		go func() {
// 			humanPhoneContactCreateErrChann <- u.CreateHumanContact(ctx, client, user, user.HumanPhoneContact)
// 		}()
// 	}

// 	humanSecurityCreateErrChann := make(chan error, 1)
// 	go func() {
// 		humanSecurityCreateErrChann <- u.Human().Security().Create(ctx, client, user)
// 	}()

// 	if err := <-humanEmailContactCreateErrChann; err != nil {
// 		return nil, err
// 	}

// 	if humanPhoneContactCreateErrChann != nil {
// 		if err := <-humanPhoneContactCreateErrChann; err != nil {
// 			return nil, err
// 		}
// 	}

// 	if err := <-humanSecurityCreateErrChann; err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// const createHumaneStmt = `INSERT INTO zitadel.human_users (instance_id, org_id, id, username, username_org_unique, state,` +
// 	` first_name, last_name, nick_name, display_name, preferred_language, gender, avatar_key)` +
// 	` VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)` +
// 	` RETURNING created_at, updated_at`

// func (u user) createHuman(ctx context.Context, client database.QueryExecutor, user *domain.Human) (*domain.Human, error) {
// 	builder := database.StatementBuilder{}
// 	builder.AppendArgs(user.User.InstanceID, user.User.OrgID, user.ID, user.Username, user.IsUsernameOrgUnique, user.State)
// 	builder.AppendArgs(user.FirstName, user.LastName, user.NickName, user.DisplayName, user.PreferredLanguage, user.Gender, user.AvatarKey)

// 	builder.WriteString(createHumaneStmt)

// 	err := client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&user.CreatedAt, &user.UpdatedAt)
// 	return user, err
// }

// const createHumaneContactStmt = `INSERT INTO zitadel.human_contacts (instance_id, org_id, user_id,` +
// 	` type, value, is_verified, unverified_value)` +
// 	` VALUES($1, $2, $3, $4, $5, $6, $7)`

// func (u user) CreateHumanContact(ctx context.Context, client database.QueryExecutor, user *domain.Human, contact *domain.HumanContact) error {
// 	builder := database.StatementBuilder{}
// 	builder.AppendArgs(user.User.InstanceID, user.User.OrgID, user.ID)
// 	builder.AppendArgs(contact.Type, contact.Value, contact.IsVerified, contact.UnverifiedValue)

// 	builder.WriteString(createHumaneContactStmt)

// 	_, err := client.Exec(ctx, builder.String(), builder.Args()...)
// 	return err
// }

// func (u user) UpdateHuman(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
// 	if len(changes) == 0 {
// 		return 0, database.ErrNoChanges
// 	}

// 	if !condition.IsRestrictingColumn(u.human.InstanceIDColumn()) {
// 		return 0, database.NewMissingConditionError(u.human.InstanceIDColumn())
// 	}

// 	// if !condition.IsRestrictingColumn(u.human.OrgIDColumn()) {
// 	// 	return 0, database.NewMissingConditionError(u.human.OrgIDColumn())
// 	// }

// 	if !condition.IsRestrictingColumn(u.human.IDColumn()) {
// 		return 0, database.NewMissingConditionError(u.human.IDColumn())
// 	}

// 	var builder database.StatementBuilder
// 	builder.WriteString(`UPDATE zitadel.human_users SET `)
// 	database.Changes(changes).Write(&builder)
// 	writeCondition(&builder, condition)

// 	return client.Exec(ctx, builder.String(), builder.Args()...)
// }

// func (c contact) Set(ctx context.Context, client database.QueryExecutor, contact domain.HumanContact) (int64, error) {
// 	// 	for _, opt := range opts {
// 	// 		opt(options)
// 	return 0, nil
// }

// // TODO use LeftJoin()
// const queryHumanUserStmt = `SELECT zitadel.human_users.instance_id, zitadel.human_users.org_id, id, username, username_org_unique, state,` +
// 	` first_name, last_name, nick_name, display_name, preferred_language, gender, avatar_key,` +
// 	` created_at, updated_at,` +
// 	// email
// 	` email.type AS "email.type", email.value AS "email.value", email.is_verified AS "email.is_verified", email.unverified_value AS "email.unverified_value",` +
// 	// phone
// 	` phone.type AS "phone.type", phone.value AS "phone.value", phone.is_verified AS "phone.is_verified", phone.unverified_value AS "phone.unverified_value"` +
// 	` FROM zitadel.human_users` +
// 	// join email
// 	` LEFT JOIN zitadel.human_contacts AS email ON zitadel.human_users.id = email.user_id` +
// 	` AND zitadel.human_users.instance_id = email.instance_id` +
// 	` AND zitadel.human_users.org_id = email.org_id` +
// 	` AND email.type = 'email'` +
// 	// join phone
// 	` LEFT JOIN zitadel.human_contacts AS phone ON zitadel.human_users.id = phone.user_id` +
// 	` AND zitadel.human_users.instance_id = phone.instance_id` +
// 	` AND zitadel.human_users.org_id = phone.org_id` +
// 	` AND phone.type = 'phone'`

// func (u user) GetHuman(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Human, error) {
// 	options := new(database.QueryOpts)
// 	for _, opt := range opts {
// 		opt(options)
// 	}

// 	if !options.Condition.IsRestrictingColumn(u.human.InstanceIDColumn()) {
// 		return nil, database.NewMissingConditionError(u.human.InstanceIDColumn())
// 	}

// 	// if !options.Condition.IsRestrictingColumn(u.human.OrgIDColumn()) {
// 	// 	return nil, database.NewMissingConditionError(u.human.OrgIDColumn())
// 	// }

// 	// if !options.Condition.IsRestrictingColumn(u.human.IDColumn()) {
// 	// 	return nil, database.NewMissingConditionError(u.human.IDColumn())
// 	// }

// 	var builder database.StatementBuilder
// 	builder.WriteString(queryHumanUserStmt)
// 	options.Write(&builder)

// 	user, err := scanHuman(ctx, client, &builder)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if user.HumanPhoneContact != nil && user.HumanPhoneContact.Value == nil {
// 		user.HumanPhoneContact = nil
// 	}

// 	return user, nil
// }

// func (u user) ListHuman(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Human, error) {
// 	builder := database.StatementBuilder{}

// 	options := new(database.QueryOpts)
// 	for _, opt := range opts {
// 		opt(options)
// 	}

// 	if !options.Condition.IsRestrictingColumn(u.Human().InstanceIDColumn()) {
// 		return nil, database.NewMissingConditionError(u.Human().InstanceIDColumn())
// 	}

// 	builder.WriteString(queryHumanUserStmt)
// 	options.Write(&builder)

// 	orderBy := database.OrderBy(u.Human().CreatedAtColumn())
// 	orderBy.Write(&builder)

// 	return scanHumans(ctx, client, &builder)
// }

// func (u user) DeleteHuman(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
// 	if !condition.IsRestrictingColumn(u.Human().InstanceIDColumn()) {
// 		return 0, database.NewMissingConditionError(u.Human().InstanceIDColumn())
// 	}

// 	// if !condition.IsRestrictingColumn(u.Human().OrgIDColumn()) {
// 	// 	return 0, database.NewMissingConditionError(u.Human().OrgIDColumn())
// 	// }

// 	if !condition.IsRestrictingColumn(u.Human().IDColumn()) {
// 		return 0, database.NewMissingConditionError(u.Human().IDColumn())
// 	}

// 	var builder database.StatementBuilder
// 	builder.WriteString(`DELETE FROM zitadel.human_users`)
// 	writeCondition(&builder, condition)

// 	return client.Exec(ctx, builder.String(), builder.Args()...)
// }

// const createMachineStmt = `INSERT INTO zitadel.machine_users (instance_id, org_id, id, username, username_org_unique, state,` +
// 	` name, description)` +
// 	` VALUES($1, $2, $3, $4, $5, $6, $7, $8)` +
// 	` RETURNING created_at, updated_at`

// func (u user) CreateMachine(ctx context.Context, client database.QueryExecutor, user *domain.Machine) (*domain.Machine, error) {
// 	builder := database.StatementBuilder{}
// 	builder.AppendArgs(user.InstanceID, user.OrgID, user.ID, user.Username, user.IsUsernameOrgUnique, user.State)
// 	builder.AppendArgs(user.Name, user.Description)

// 	builder.WriteString(createMachineStmt)

// 	err := client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&user.CreatedAt, &user.UpdatedAt)
// 	return user, err
// }

// func (u user) UpdateMachine(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
// 	if len(changes) == 0 {
// 		return 0, database.ErrNoChanges
// 	}

// 	if !condition.IsRestrictingColumn(u.machine.InstanceIDColumn()) {
// 		return 0, database.NewMissingConditionError(u.machine.InstanceIDColumn())
// 	}

// 	if !condition.IsRestrictingColumn(u.machine.OrgIDColumn()) {
// 		return 0, database.NewMissingConditionError(u.machine.OrgIDColumn())
// 	}

// 	if !condition.IsRestrictingColumn(u.machine.IDColumn()) {
// 		return 0, database.NewMissingConditionError(u.machine.IDColumn())
// 	}

// 	var builder database.StatementBuilder
// 	builder.WriteString(`UPDATE zitadel.machine_users SET `)
// 	database.Changes(changes).Write(&builder)
// 	writeCondition(&builder, condition)

// 	return client.Exec(ctx, builder.String(), builder.Args()...)
// }

// const querMachineUserStmt = `SELECT instance_id, org_id, id, username, username_org_unique, state,` +
// 	` name, description,` +
// 	` created_at, updated_at` +
// 	` FROM zitadel.machine_users`

// func (u user) GetMachine(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Machine, error) {
// 	options := new(database.QueryOpts)
// 	for _, opt := range opts {
// 		opt(options)
// 	}

// 	options = new(database.QueryOpts)
// 	for _, opt := range opts {
// 		opt(options)
// 	}

// 	if !options.Condition.IsRestrictingColumn(u.machine.InstanceIDColumn()) {
// 		return nil, database.NewMissingConditionError(u.machine.InstanceIDColumn())
// 	}

// 	if !options.Condition.IsRestrictingColumn(u.machine.OrgIDColumn()) {
// 		return nil, database.NewMissingConditionError(u.machine.OrgIDColumn())
// 	}

// 	if !options.Condition.IsRestrictingColumn(u.machine.IDColumn()) {
// 		return nil, database.NewMissingConditionError(u.machine.IDColumn())
// 	}

// 	var builder database.StatementBuilder
// 	builder.WriteString(querMachineUserStmt)
// 	options.Write(&builder)

// 	return scanMachine(ctx, client, &builder)
// }

// func (u user) ListMachine(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Machine, error) {
// 	builder := database.StatementBuilder{}

// 	options := new(database.QueryOpts)
// 	for _, opt := range opts {
// 		opt(options)
// 	}

// 	if !options.Condition.IsRestrictingColumn(u.Machine().InstanceIDColumn()) {
// 		return nil, database.NewMissingConditionError(u.Machine().InstanceIDColumn())
// 	}

// 	builder.WriteString(querMachineUserStmt)
// 	options.Write(&builder)

// 	orderBy := database.OrderBy(u.Machine().CreatedAtColumn())
// 	orderBy.Write(&builder)

// 	return scanMachines(ctx, client, &builder)
// }

// func (u user) DeleteMachine(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
// 	if !condition.IsRestrictingColumn(u.Machine().InstanceIDColumn()) {
// 		return 0, database.NewMissingConditionError(u.Machine().InstanceIDColumn())
// 	}

// 	if !condition.IsRestrictingColumn(u.Machine().OrgIDColumn()) {
// 		return 0, database.NewMissingConditionError(u.Machine().OrgIDColumn())
// 	}

// 	if !condition.IsRestrictingColumn(u.Machine().IDColumn()) {
// 		return 0, database.NewMissingConditionError(u.Machine().IDColumn())
// 	}

// 	var builder database.StatementBuilder
// 	builder.WriteString(`DELETE FROM zitadel.machine_users`)
// 	writeCondition(&builder, condition)

// 	return client.Exec(ctx, builder.String(), builder.Args()...)
// }

// const createHumanSecurityStmt = `INSERT INTO zitadel.human_security (instance_id, org_id, user_id, password_change_required, password_changed, mfa_init_skipped)` +
// 	` VALUES($1, $2, $3, $4, $5, $6)`

// func (s security) Create(ctx context.Context, client database.QueryExecutor, user *domain.Human) error {
// 	builder := database.StatementBuilder{}
// 	builder.AppendArgs(user.User.InstanceID, user.User.OrgID, user.ID, user.HumanSecurity.PasswordChangeRequired,
// 		user.HumanSecurity.PasswordChange, user.HumanSecurity.MFAInitSkipped)

// 	builder.WriteString(createHumanSecurityStmt)

// 	_, err := client.Exec(ctx, builder.String(), builder.Args()...)
// 	return err
// }

// const querySecuirtyStmt = `SELECT zitadel.human_security.instance_id, zitadel.human_security.org_id, zitadel.human_security.user_id,` +
// 	` password_change_required, password_changed, mfa_init_skipped FROM zitadel.human_security`

// func (s security) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.HumanSecurity, error) {
// 	builder := database.StatementBuilder{}
// 	builder.WriteString(querMachineUserStmt)

// 	options := new(database.QueryOpts)
// 	for _, opt := range opts {
// 		opt(options)
// 	}

// 	human := human{}
// 	database.WithLeftJoin(
// 		human.qualifiedTableName(),
// 		database.And(
// 			database.NewColumnCondition(s.InstanceIDColumn(), human.InstanceIDColumn()),
// 			database.NewColumnCondition(s.OrgIDColumn(), human.OrgIDColumn()),
// 			database.NewColumnCondition(s.UserIDColumn(), human.OrgIDColumn()),
// 		),
// 	)(options)

// 	options.Write(&builder)

// 	return getOne[domain.HumanSecurity](ctx, client, &builder)
// }

// func (s security) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
// 	if len(changes) == 0 {
// 		return 0, database.ErrNoChanges
// 	}

// 	if !condition.IsRestrictingColumn(s.InstanceIDColumn()) {
// 		return 0, database.NewMissingConditionError(s.InstanceIDColumn())
// 	}

// 	// if !condition.IsRestrictingColumn(human.OrgIDColumn()) {
// 	// 	return 0, database.NewMissingConditionError(human.OrgIDColumn())
// 	// }

// 	if !condition.IsRestrictingColumn(s.UserIDColumn()) {
// 		return 0, database.NewMissingConditionError(s.UserIDColumn())
// 	}

// 	var builder database.StatementBuilder
// 	builder.WriteString(`UPDATE zitadel.human_security SET `)
// 	database.Changes(changes).Write(&builder)
// 	writeCondition(&builder, condition)

// 	return client.Exec(ctx, builder.String(), builder.Args()...)
// }

// func scanMachine(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.Machine, error) {
// 	user := &domain.Machine{}
// 	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = rows.(database.CollectableRows).CollectExactlyOneRow(user)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return user, err
// }

// func scanMachines(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.Machine, error) {
// 	users := []*domain.Machine{}

// 	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = rows.(database.CollectableRows).Collect(&users)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return users, err
// }

// func scanHuman(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.Human, error) {
// 	user := &domain.Human{}
// 	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = rows.(database.CollectableRows).CollectExactlyOneRow(user)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return user, err
// }

// func scanHumans(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.Human, error) {
// 	users := []*domain.Human{}

// 	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = rows.(database.CollectableRows).Collect(&users)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return users, err
// }
