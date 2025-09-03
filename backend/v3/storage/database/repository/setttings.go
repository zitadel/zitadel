package repository

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type settings struct {
	repository
}

func SettingsRepository(client database.QueryExecutor) domain.SettingsRepository {
	return &settings{
		repository: repository{
			client: client,
		},
	}
}

var _ domain.SettingsRepository = (*settings)(nil)

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (settings) IDColumn() database.Column {
	return database.NewColumn("settings", "id")
}

func (settings) InstanceIDColumn() database.Column {
	return database.NewColumn("settings", "instance_id")
}

func (settings) OrgIDColumn() database.Column {
	return database.NewColumn("settings", "org_id")
}

func (settings) TypeColumn() database.Column {
	return database.NewColumn("settings", "type")
}

func (settings) SettingsColumn() database.Column {
	return database.NewColumn("settings", "settings")
}

func (settings) CreatedAtColumn() database.Column {
	return database.NewColumn("settings", "created_at")
}

func (settings) UpdatedAtColumn() database.Column {
	return database.NewColumn("settings", "updated_at")
}

// // -------------------------------------------------------------
// // conditions
// // -------------------------------------------------------------

func (s settings) InstanceIDCondition(id string) database.Condition {
	return database.NewTextCondition(s.InstanceIDColumn(), database.TextOperationEqual, id)
}

func (s settings) OrgIDCondition(id *string) database.Condition {
	if id == nil {
		return database.IsNull(s.OrgIDColumn())
	}
	return database.NewTextCondition(s.OrgIDColumn(), database.TextOperationEqual, *id)
}

func (s settings) IDCondition(id string) database.Condition {
	return database.NewTextCondition(s.IDColumn(), database.TextOperationEqual, id)
}

func (s settings) TypeCondition(typ domain.SettingType) database.Condition {
	return database.NewTextCondition(s.TypeColumn(), database.TextOperationEqual, typ.String())
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

func (s settings) SetType(name domain.SettingType) database.Change {
	return database.NewChange(s.TypeColumn(), name.String())
}

func (s settings) SetSettings(settings string) database.Change {
	return database.NewChange(s.SettingsColumn(), settings)
}

const querySettingStmt = `SELECT instance_id, org_id, id, type, settings,` +
	` created_at, updated_at` +
	` FROM zitadel.settings`

func (s *settings) Get(ctx context.Context, id string, instanceID string, orgID *string) (*domain.Setting, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(querySettingStmt)

	conditions := []database.Condition{s.IDCondition(id), s.InstanceIDCondition(instanceID), s.OrgIDCondition(orgID)}

	writeCondition(&builder, database.And(conditions...))

	return scanSetting(ctx, s.client, &builder)
}

func (s *settings) List(ctx context.Context, conditions ...database.Condition) ([]*domain.Setting, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(querySettingStmt)

	if conditions != nil {
		writeCondition(&builder, database.And(conditions...))
	}

	orderBy := database.OrderBy(s.CreatedAtColumn())
	orderBy.Write(&builder)

	return scanSettings(ctx, s.client, &builder)
}

const createSettingStmt = `INSERT INTO zitadel.settings` +
	` (instance_id, org_id, id, type, settings)` +
	` VALUES ($1, $2, $3, $4, $5)` +
	` RETURNING created_at, updated_at`

func (s *settings) Create(ctx context.Context, setting *domain.Setting) error {
	builder := database.StatementBuilder{}
	builder.AppendArgs(
		setting.InstanceID,
		setting.OrgID,
		setting.ID,
		setting.Type,
		string(setting.Settings))
	builder.WriteString(createSettingStmt)

	return s.client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&setting.CreatedAt, &setting.UpdatedAt)
}

func (s *settings) Update(ctx context.Context, id string, instanceID string, orgID *string, changes ...database.Change) (int64, error) {
	if changes == nil {
		return 0, errors.New("Update must contain at least one change")
	}
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.settings SET `)

	conditions := []database.Condition{
		s.IDCondition(id),
		s.InstanceIDCondition(instanceID),
		s.OrgIDCondition(orgID),
	}
	database.Changes(changes).Write(&builder)
	writeCondition(&builder, database.And(conditions...))

	stmt := builder.String()

	return s.client.Exec(ctx, stmt, builder.Args()...)
}

func (s *settings) Delete(ctx context.Context, id string, instanceID string, orgID *string) (int64, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(`DELETE FROM zitadel.settings`)

	conditions := []database.Condition{
		s.IDCondition(id),
		s.InstanceIDCondition(instanceID),
		s.OrgIDCondition(orgID),
	}
	writeCondition(&builder, database.And(conditions...))

	return s.client.Exec(ctx, builder.String(), builder.Args()...)
}

func scanSetting(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.Setting, error) {
	setting := &domain.Setting{}

	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	err = rows.(database.CollectableRows).CollectExactlyOneRow(setting)
	if err != nil {
		return nil, err
	}

	return setting, err
}

func scanSettings(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.Setting, error) {
	settings := []*domain.Setting{}

	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	err = rows.(database.CollectableRows).Collect(&settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}
