package repository

import (
	"context"
	"encoding/json"
	"slices"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type setting[T domain.SettingTypes] struct{}

var unqualifiedSettingsTableName = "settings"

func (setting[T]) unqualifiedTableName() string {
	return unqualifiedSettingsTableName
}

// Delete implements [domain.SettingsRepository].
func (s setting[T]) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if !condition.IsRestrictingColumn(s.TypeColumn()) {
		var t T
		condition = database.And(
			s.typeCondition(t.Type()),
			condition,
		)
	}

	if err := checkRestrictingColumns(condition, s.UniqueColumns()...); err != nil {
		return 0, err
	}

	builder := database.NewStatementBuilder("DELETE FROM zitadel.settings")
	writeCondition(builder, condition)
	return client.Exec(ctx, builder.String(), builder.Args()...)
}

const (
	settingsEnsurePrefix = "WITH empty AS (SELECT) INSERT INTO zitadel.settings " +
		"SELECT COALESCE(settings.instance_id, $1), COALESCE(settings.organization_id, $2)" +
		", COALESCE(settings.id, gen_random_uuid()), COALESCE(settings.type, $3)" +
		", COALESCE(settings.created_at, now()), COALESCE(settings.updated_at, now())" +
		", COALESCE(settings.state, $4), "
	settingsEnsureSuffix = " FROM empty LEFT JOIN zitadel.settings ON" +
		" settings.instance_id = $1 AND" +
		" settings.organization_id = $2 AND" +
		" settings.type = $3 AND" +
		" settings.state = $4" +
		" ON CONFLICT (instance_id, organization_id, type, state) DO UPDATE SET" +
		" attributes = EXCLUDED.attributes, updated_at = EXCLUDED.updated_at"
)

func (s setting[T]) ensure(ctx context.Context, client database.QueryExecutor, instanceID string, organizationID *string, state domain.SettingState, changes ...database.Change) error {
	var t T
	builder := database.NewStatementBuilder(settingsEnsurePrefix,
		instanceID, organizationID, t.Type(), state,
	)

	var firstAttributeChange *attributeChange
	changes = slices.DeleteFunc(changes, func(change database.Change) bool {
		attrChange, ok := change.(*attributeChange)
		if !ok {
			return false
		}
		if firstAttributeChange == nil {
			firstAttributeChange = attrChange
			return true
		}
		firstAttributeChange.chain(attrChange)
		return true
	})
	firstAttributeChange.change(database.Coalesce(s.attributesColumn(), json.RawMessage("{}"))).
		Write(builder)
	builder.WriteString(settingsEnsureSuffix)
	writeCondition(builder, s.UniqueCondition(instanceID, organizationID, t.Type(), state))

	_, err := client.Exec(ctx, builder.String(), builder.Args()...)
	return err
}

const (
	settingsSetPrefix = "INSERT INTO zitadel.settings" +
		"(id, instance_id, organization_id, type, state, attributes, created_at, updated_at) VALUES ("
	settingsSetSuffix = ") ON CONFLICT (instance_id, organization_id, type, state) DO UPDATE SET (" +
		"attributes, updated_at) = (EXCLUDED.attributes, EXCLUDED.updated_at)"
)

// Set implements [domain.SettingsRepository].
func (s setting[T]) Set(ctx context.Context, client database.QueryExecutor, setting *T) error {
	var (
		id        any = database.DefaultInstruction
		createdAt any = database.NowInstruction
		updatedAt any = database.NowInstruction
		state         = domain.SettingStateActive
	)

	base := (*setting).Base()

	if base.ID != "" {
		id = base.ID
	}
	if !base.CreatedAt.IsZero() {
		createdAt = base.CreatedAt
	}
	if !base.UpdatedAt.IsZero() {
		updatedAt = base.UpdatedAt
	}

	if brandingSetting, ok := any(setting).(*domain.BrandingSetting); ok {
		state = brandingSetting.State
	}

	builder := database.NewStatementBuilder(settingsSetPrefix)
	builder.WriteArgs(id, base.InstanceID, base.OrganizationID, (*setting).Type(),
		state, setting, createdAt, updatedAt,
	)
	builder.WriteString(settingsSetSuffix)
	_, err := client.Exec(ctx, builder.String(), builder.Args()...)
	return err
}

const settingsQueryPrefix = "SELECT settings.instance_id, settings.organization_id, settings.id, " +
	"settings.created_at, settings.updated_at, settings.state, settings.attributes FROM zitadel.settings"

type rawSetting[T domain.SettingTypes] struct {
	domain.Setting
	State      domain.SettingState `db:"state"`
	Attributes *T                  `json:"attributes" db:"attributes"`
}

// get handles the [domain.SettingsRepository.Get].
// The query must be built by the caller to ensure the correct columns are used.
func (s setting[T]) get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*rawSetting[T], error) {
	var dt T
	opts = append(opts, database.EnsureConditionOnColumn(s.TypeColumn(), s.typeCondition(dt.Type())))

	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if err := checkRestrictingColumns(options.Condition, s.InstanceIDColumn()); err != nil {
		return nil, err
	}

	builder := database.NewStatementBuilder(settingsQueryPrefix)
	options.Write(builder)

	return getOne[rawSetting[T]](ctx, client, builder)
}

// list handles the [domain.SettingsRepository.List].
// The query must be built by the caller to ensure the correct columns are used.
func (s setting[T]) list(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*rawSetting[T], error) {
	var dt T
	opts = append(opts, database.EnsureConditionOnColumn(s.TypeColumn(), s.typeCondition(dt.Type())))

	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if err := checkRestrictingColumns(options.Condition, s.InstanceIDColumn()); err != nil {
		return nil, err
	}

	builder := database.NewStatementBuilder(settingsQueryPrefix)
	options.Write(builder)

	return getMany[rawSetting[T]](ctx, client, builder)
}

// func (s setting[T]) Ensure(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes database.Changes) error {

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// InstanceIDCondition implements [domain.SettingsRepository].
func (s setting[T]) InstanceIDCondition(id string) database.Condition {
	return database.NewTextCondition(s.InstanceIDColumn(), database.TextOperationEqual, id)
}

// OrganizationIDCondition implements [domain.SettingsRepository].
func (s setting[T]) OrganizationIDCondition(id *string) database.Condition {
	if id == nil {
		return database.IsNull(s.OrganizationIDColumn())
	}
	return database.NewTextCondition(s.OrganizationIDColumn(), database.TextOperationEqual, *id)
}

// UniqueCondition implements [domain.SettingsRepository].
func (s setting[T]) UniqueCondition(instanceID string, orgID *string, typ domain.SettingType, state domain.SettingState) database.Condition {
	conditions := make([]database.Condition, 0, 4)
	conditions = append(conditions, s.InstanceIDCondition(instanceID), s.typeCondition(typ), s.stateCondition(state))

	if orgID == nil {
		return database.And(append(conditions, database.ForceRestrictingColumn(database.IsNull(s.OrganizationIDColumn()), true))...)
	}
	return database.And(append(conditions, database.NewTextCondition(s.OrganizationIDColumn(), database.TextOperationEqual, *orgID))...)
}

func (s setting[T]) typeCondition(typ domain.SettingType) database.Condition {
	return database.NewNumberCondition(s.TypeColumn(), database.NumberOperationEqual, typ)
}

func (s setting[T]) stateCondition(state domain.SettingState) database.Condition {
	return database.NewNumberCondition(s.StateColumn(), database.NumberOperationEqual, state)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetUpdatedAt implements [domain.SettingsRepository].
func (s setting[T]) SetUpdatedAt(updatedAt time.Time) database.Change {
	if updatedAt.IsZero() {
		return database.NewChange(s.UpdatedAtColumn(), database.NowInstruction)
	}
	return database.NewChange(s.UpdatedAtColumn(), updatedAt)
}

func (s setting[T]) setAttributes(value *domain.Setting) database.Changes {
	var t T
	return database.NewChanges(
		s.SetUpdatedAt(value.UpdatedAt),
		s.setCreatedAt(value.CreatedAt),
		s.setInstanceID(value.InstanceID),
		s.setOrganizationID(value.OrganizationID),
		s.setState(domain.SettingStateActive),
		s.setType(t.Type()),
		s.setID(value.ID),
	)
}

func (s setting[T]) setCreatedAt(value time.Time) database.Change {
	if value.IsZero() {
		return database.NewChange(s.CreatedAtColumn(), database.NowInstruction)
	}
	return database.NewChange(s.CreatedAtColumn(), value)
}

func (s setting[T]) setInstanceID(value string) database.Change {
	return database.NewChange(s.InstanceIDColumn(), value)
}

func (s setting[T]) setOrganizationID(value *string) database.Change {
	return database.NewChangePtr(s.OrganizationIDColumn(), value)
}

func (s setting[T]) setState(value domain.SettingState) database.Change {
	return database.NewChange(s.StateColumn(), value)
}

func (s setting[T]) setType(value domain.SettingType) database.Change {
	return database.NewChange(s.TypeColumn(), value)
}

func (s setting[T]) setID(value string) database.Change {
	return database.NewChange(s.idColumn(), value)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// PrimaryKeyColumns implements [domain.SettingsRepository].
func (s setting[T]) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		s.InstanceIDColumn(),
		s.idColumn(),
	}
}

// UniqueColumns implements [domain.SettingsRepository].
func (s setting[T]) UniqueColumns() []database.Column {
	return database.Columns{
		s.InstanceIDColumn(),
		s.OrganizationIDColumn(),
		s.TypeColumn(),
		s.StateColumn(),
	}
}

// CreatedAtColumn implements [domain.SettingsRepository].
func (s setting[T]) CreatedAtColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "created_at")
}

// InstanceIDColumn implements [domain.SettingsRepository].
func (s setting[T]) InstanceIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "instance_id")
}

// OrganizationIDColumn implements [domain.SettingsRepository].
func (s setting[T]) OrganizationIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "organization_id")
}

// StateColumn implements [domain.SettingsRepository].
func (s setting[T]) StateColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "state")
}

// TypeColumn implements [domain.SettingsRepository].
func (s setting[T]) TypeColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "type")
}

// UpdatedAtColumn implements [domain.SettingsRepository].
func (s setting[T]) UpdatedAtColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "updated_at")
}

func (s setting[T]) idColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "id")
}

func (s setting[T]) attributesColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "attributes")
}

var settingAttributesColumn = database.NewColumn(unqualifiedSettingsTableName, "attributes")

type attributeChange struct {
	chained *attributeChange
	change  func(column database.Column) database.Change
}

// Equals implements [database.Column].
func (a *attributeChange) Equals(col database.Column) bool {
	return col.Equals(settingAttributesColumn)
}

// WriteQualified implements [database.Column].
func (a *attributeChange) WriteQualified(builder *database.StatementBuilder) {
	a.Write(builder)
}

// WriteUnqualified implements [database.Column].
func (a *attributeChange) WriteUnqualified(builder *database.StatementBuilder) {
	a.Write(builder)
}

// IsOnColumn implements [database.Change].
func (a *attributeChange) IsOnColumn(col database.Column) bool {
	return col.Equals(settingAttributesColumn)
}

// Matches implements [database.Change].
func (a *attributeChange) Matches(x any) bool {
	toMatch, ok := x.(*attributeChange)
	if !ok {
		return false
	}
	var wantBuilder, gotBuilder database.StatementBuilder
	a.change(settingAttributesColumn).WriteArg(&wantBuilder)
	toMatch.change(settingAttributesColumn).WriteArg(&gotBuilder)
	return wantBuilder.String() == gotBuilder.String()
}

// String implements [database.Change].
func (a *attributeChange) String() string {
	return "repository.setting.attributeChange"
}

// Write implements [database.Change].
func (a *attributeChange) Write(builder *database.StatementBuilder) {
	if a.chained == nil {
		a.change(settingAttributesColumn).Write(builder)
		return
	}
	a.change(a.chained).Write(builder)
}

// WriteArg implements [database.Change].
func (a *attributeChange) WriteArg(builder *database.StatementBuilder) {
	a.Write(builder)
}

func (a *attributeChange) chain(next *attributeChange) *attributeChange {
	if a.chained == nil {
		a.chained = next
		return a
	}
	return a.chained.chain(next)
}

var (
	_ database.Change = (*attributeChange)(nil)
	_ database.Column = (*attributeChange)(nil)
)
