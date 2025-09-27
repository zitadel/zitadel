package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type settings struct{}

func SettingsRepository() domain.SettingsRepository {
	return new(settings)
}

var _ domain.SettingsRepository = (*settings)(nil)

var (
	ErrSettingObjectMustNotBeNil error = errors.New("setting object must not be nill")
	ErrLabelStateMustBeDefined   error = errors.New("label state must be defined")
)

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

func (settings) LabelStateColumn() database.Column {
	return database.NewColumn("settings", "label_state")
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

func (s settings) LabelStateCondition(typ domain.LabelState) database.Condition {
	return database.NewTextCondition(s.LabelStateColumn(), database.TextOperationEqual, typ.String())
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

func (s settings) SetUpdatedAt(updatedAt *time.Time) database.Change {
	return database.NewChangePtr(s.UpdatedAtColumn(), updatedAt)
}

const querySettingStmt = `SELECT instance_id, org_id, id, type, label_state, settings,` +
	` created_at, updated_at` +
	` FROM zitadel.settings`

func (s *settings) Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string, typ domain.SettingType, cond ...database.Condition) (*domain.Setting, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(querySettingStmt)

	cond = append(cond, s.TypeCondition(typ), s.InstanceIDCondition(instanceID), s.OrgIDCondition(orgID))

	writeCondition(&builder, database.And(cond...))

	return scanSetting(ctx, client, &builder)
}

func (s *settings) List(ctx context.Context, client database.QueryExecutor, conditions ...database.Condition) ([]*domain.Setting, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(querySettingStmt)

	if conditions != nil {
		writeCondition(&builder, database.And(conditions...))
	}

	orderBy := database.OrderBy(s.CreatedAtColumn())
	orderBy.Write(&builder)

	return scanSettings(ctx, client, &builder)
}

func (s *settings) CreateLogin(ctx context.Context, client database.QueryExecutor, setting *domain.LoginSetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeLogin
	return s.createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *settings) GetLogin(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.LoginSetting, error) {
	loginSetting := &domain.LoginSetting{}
	var err error

	loginSetting.Setting, err = s.Get(ctx, client, instanceID, orgID, domain.SettingTypeLogin)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(loginSetting.Setting.Settings, &loginSetting.Settings)
	if err != nil {
		return nil, err
	}

	return loginSetting, nil
}

func (s *settings) UpdateLogin(ctx context.Context, client database.QueryExecutor, setting *domain.LoginSetting, changes ...database.Change) (int64, error) {
	return s.updateSetting(ctx, client, setting.Setting, &setting.Settings, changes...)
}

func (s *settings) CreateLabel(ctx context.Context, client database.QueryExecutor, setting *domain.LabelSetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	if setting.LabelState == nil {
		return ErrLabelStateMustBeDefined
	}

	setting.Type = domain.SettingTypeLabel
	return s.createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *settings) GetLabel(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string, state domain.LabelState) (*domain.LabelSetting, error) {
	labelSetting := &domain.LabelSetting{}
	var err error

	stateCond := s.LabelStateCondition(state)
	labelSetting.Setting, err = s.Get(ctx, client, instanceID, orgID, domain.SettingTypeLabel, stateCond)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(labelSetting.Setting.Settings, &labelSetting.Settings)
	if err != nil {
		return nil, err
	}

	return labelSetting, nil
}

func (s *settings) UpdateLabel(ctx context.Context, client database.QueryExecutor, setting *domain.LabelSetting, changes ...database.Change) (int64, error) {
	if setting.LabelState == nil {
		return 0, ErrLabelStateMustBeDefined
	}
	return s.updateSetting(ctx, client, setting.Setting, &setting.Settings, changes...)
}

func (s *settings) CreatePasswordComplexity(ctx context.Context, client database.QueryExecutor, setting *domain.PasswordComplexitySetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypePasswordComplexity
	return s.createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *settings) GetPasswordComplexity(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.PasswordComplexitySetting, error) {
	passwordComplexitySetting := &domain.PasswordComplexitySetting{}
	var err error

	passwordComplexitySetting.Setting, err = s.Get(ctx, client, instanceID, orgID, domain.SettingTypePasswordComplexity)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(passwordComplexitySetting.Setting.Settings, &passwordComplexitySetting.Settings)
	if err != nil {
		return nil, err
	}

	return passwordComplexitySetting, nil
}

func (s *settings) UpdatePasswordComplexity(ctx context.Context, client database.QueryExecutor, setting *domain.PasswordComplexitySetting, changes ...database.Change) (int64, error) {
	return s.updateSetting(ctx, client, setting.Setting, &setting.Settings, changes...)
}

func (s *settings) CreatePasswordExpiry(ctx context.Context, client database.QueryExecutor, setting *domain.PasswordExpirySetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypePasswordExpiry
	return s.createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *settings) GetPasswordExpiry(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.PasswordExpirySetting, error) {
	passwordPolicySetting := &domain.PasswordExpirySetting{}
	var err error

	passwordPolicySetting.Setting, err = s.Get(ctx, client, instanceID, orgID, domain.SettingTypePasswordExpiry)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(passwordPolicySetting.Setting.Settings, &passwordPolicySetting.Settings)
	if err != nil {
		return nil, err
	}

	return passwordPolicySetting, nil
}

func (s *settings) UpdatePasswordExpiry(ctx context.Context, client database.QueryExecutor, setting *domain.PasswordExpirySetting, changes ...database.Change) (int64, error) {
	return s.updateSetting(ctx, client, setting.Setting, &setting.Settings, changes...)
}

func (s *settings) CreateSecurity(ctx context.Context, client database.QueryExecutor, setting *domain.SecuritySetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeSecurity
	return s.createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *settings) GetSecurity(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.SecuritySetting, error) {
	securitySetting := &domain.SecuritySetting{}
	var err error

	securitySetting.Setting, err = s.Get(ctx, client, instanceID, orgID, domain.SettingTypeSecurity)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(securitySetting.Setting.Settings, &securitySetting.Settings)
	if err != nil {
		return nil, err
	}

	return securitySetting, nil
}

func (s *settings) UpdateSecurity(ctx context.Context, client database.QueryExecutor, setting *domain.SecuritySetting, changes ...database.Change) (int64, error) {
	return s.updateSetting(ctx, client, setting.Setting, &setting.Settings, changes...)
}

func (s *settings) CreateLockout(ctx context.Context, client database.QueryExecutor, setting *domain.LockoutSetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeLockout
	return s.createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *settings) GetLockout(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.LockoutSetting, error) {
	lockoutSetting := &domain.LockoutSetting{}
	var err error

	lockoutSetting.Setting, err = s.Get(ctx, client, instanceID, orgID, domain.SettingTypeLockout)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(lockoutSetting.Setting.Settings, &lockoutSetting.Settings)
	if err != nil {
		return nil, err
	}

	return lockoutSetting, nil
}

func (s *settings) UpdateLockout(ctx context.Context, client database.QueryExecutor, setting *domain.LockoutSetting, changes ...database.Change) (int64, error) {
	return s.updateSetting(ctx, client, setting.Setting, &setting.Settings, changes...)
}

func (s *settings) CreateDomain(ctx context.Context, client database.QueryExecutor, setting *domain.DomainSetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeDomain
	return s.createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *settings) GetDomain(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.DomainSetting, error) {
	lockoutSetting := &domain.DomainSetting{}
	var err error

	lockoutSetting.Setting, err = s.Get(ctx, client, instanceID, orgID, domain.SettingTypeDomain)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(lockoutSetting.Setting.Settings, &lockoutSetting.Settings)
	if err != nil {
		return nil, err
	}

	return lockoutSetting, nil
}

func (s *settings) UpdateDomain(ctx context.Context, client database.QueryExecutor, setting *domain.DomainSetting, changes ...database.Change) (int64, error) {
	return s.updateSetting(ctx, client, setting.Setting, &setting.Settings, changes...)
}

func (s *settings) CreateOrg(ctx context.Context, client database.QueryExecutor, setting *domain.OrgSetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeOrganization
	return s.createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *settings) GetOrg(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.OrgSetting, error) {
	orgSetting := &domain.OrgSetting{}
	var err error

	orgSetting.Setting, err = s.Get(ctx, client, instanceID, orgID, domain.SettingTypeOrganization)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(orgSetting.Setting.Settings, &orgSetting.Settings)
	if err != nil {
		return nil, err
	}

	return orgSetting, nil
}

func (s *settings) UpdateOrg(ctx context.Context, client database.QueryExecutor, setting *domain.OrgSetting, changes ...database.Change) (int64, error) {
	return s.updateSetting(ctx, client, setting.Setting, &setting.Settings, changes...)
}

func (s *settings) updateSetting(ctx context.Context, client database.QueryExecutor, setting *domain.Setting, settings any, changes ...database.Change) (int64, error) {
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.settings SET `)
	conditions := []database.Condition{
		s.IDCondition(setting.ID),
		s.InstanceIDCondition(setting.InstanceID),
		s.OrgIDCondition(setting.OrgID),
		s.TypeCondition(setting.Type),
	}
	settingJSON, err := json.Marshal(settings)
	if err != nil {
		return 0, err
	}

	if changes == nil || !database.Changes(changes).IsOnColumn(s.UpdatedAtColumn()) {
		// if updated_at is not set, then explicitly set it to NULL so that the db trigger sets it to NOW()
		changes = append(changes, s.SetSettings(string(settingJSON)), s.SetUpdatedAt(nil))
	} else {
		changes = append(changes, s.SetSettings(string(settingJSON)))
	}

	database.Changes(changes).Write(&builder)

	writeCondition(&builder, database.And(conditions...))

	stmt := builder.String()

	return client.Exec(ctx, stmt, builder.Args()...)
}

const eventCreateSettingStmt = `INSERT INTO zitadel.settings` +
	` (instance_id, org_id, type, label_state, settings, created_at, updated_at)` +
	` VALUES ($1, $2, $3, $4, $5, $6, $7)` +
	` RETURNING id, created_at, updated_at`

func (s *settings) Create(ctx context.Context, client database.QueryExecutor, setting *domain.Setting) error {
	builder := database.StatementBuilder{}

	builder.AppendArgs(
		setting.InstanceID,
		setting.OrgID,
		setting.Type,
		setting.LabelState,
		string(setting.Settings),
		setting.CreatedAt,
		setting.UpdatedAt)
	builder.WriteString(eventCreateSettingStmt)

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)
}

const createSettingStmt = `INSERT INTO zitadel.settings` +
	` (instance_id, org_id, type, label_state, settings)` +
	` VALUES ($1, $2, $3, $4, $5)` +
	` RETURNING id, created_at, updated_at`

func (s *settings) createSetting(ctx context.Context, client database.QueryExecutor, setting *domain.Setting, settings any) error {
	builder := database.StatementBuilder{}

	settingJSON, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	builder.AppendArgs(
		setting.InstanceID,
		setting.OrgID,
		setting.Type,
		setting.LabelState,
		string(settingJSON))

	builder.WriteString(createSettingStmt)

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)
}

func (s *settings) Delete(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string, typ domain.SettingType) (int64, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(`DELETE FROM zitadel.settings`)

	conditions := []database.Condition{
		s.TypeCondition(typ),
		s.InstanceIDCondition(instanceID),
		s.OrgIDCondition(orgID),
	}
	writeCondition(&builder, database.And(conditions...))

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (s *settings) DeleteSettingsForInstance(ctx context.Context, client database.QueryExecutor, instanceID string) (int64, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(`DELETE FROM zitadel.settings`)

	conditions := []database.Condition{
		s.InstanceIDCondition(instanceID),
	}
	writeCondition(&builder, database.And(conditions...))

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (s *settings) DeleteSettingsForOrg(ctx context.Context, client database.QueryExecutor, orgID string) (int64, error) {
	if orgID == "" {
		return 0, domain.ErrNoOrgIdSpecified
	}
	builder := database.StatementBuilder{}

	builder.WriteString(`DELETE FROM zitadel.settings`)

	conditions := []database.Condition{
		s.OrgIDCondition(&orgID),
	}
	writeCondition(&builder, database.And(conditions...))

	return client.Exec(ctx, builder.String(), builder.Args()...)
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
