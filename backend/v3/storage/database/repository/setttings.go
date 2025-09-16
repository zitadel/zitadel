package repository

import (
	"context"
	"encoding/json"

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

func (s *settings) Get(ctx context.Context, instanceID string, orgID *string, typ domain.SettingType) (*domain.Setting, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(querySettingStmt)

	conditions := []database.Condition{s.TypeCondition(typ), s.InstanceIDCondition(instanceID), s.OrgIDCondition(orgID)}

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

func (s *settings) GetLogin(ctx context.Context, instanceID string, orgID *string) (*domain.LoginSetting, error) {
	loginSetting := &domain.LoginSetting{}
	var err error

	loginSetting.Setting, err = s.Get(ctx, instanceID, orgID, domain.SettingTypeLogin)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(loginSetting.Setting.Settings, &loginSetting.Settings)
	if err != nil {
		return nil, err
	}

	return loginSetting, nil
}

func (s *settings) UpdateLogin(ctx context.Context, setting *domain.LoginSetting) (int64, error) {
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.settings SET `)

	conditions := []database.Condition{
		s.IDCondition(*setting.ID),
		s.InstanceIDCondition(setting.InstanceID),
		s.OrgIDCondition(setting.OrgID),
		s.TypeCondition(setting.Type),
	}
	settingJSON, err := json.Marshal(setting.Settings)
	if err != nil {
		return 0, err
	}
	s.SetSettings(string(settingJSON)).Write(&builder)
	writeCondition(&builder, database.And(conditions...))

	stmt := builder.String()

	return s.client.Exec(ctx, stmt, builder.Args()...)
}

func (s *settings) GetLabel(ctx context.Context, instanceID string, orgID *string) (*domain.LabelSetting, error) {
	labelSetting := &domain.LabelSetting{}
	var err error

	labelSetting.Setting, err = s.Get(ctx, instanceID, orgID, domain.SettingTypeLabel)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(labelSetting.Setting.Settings, &labelSetting.Settings)
	if err != nil {
		return nil, err
	}

	return labelSetting, nil
}

// func (s *settings) UpdateLabel(ctx context.Context, setting *domain.LabelSetting) (int64, error) {
// 	builder := database.StatementBuilder{}
// 	builder.WriteString(`UPDATE zitadel.settings SET `)
// 	fmt.Printf("[DEBUGPRINT] [settings_instance_test.go:1] setting.InstanceID = %+v\n", setting.InstanceID)

// 	conditions := []database.Condition{
// 		s.IDCondition(*setting.ID),
// 		s.InstanceIDCondition(setting.InstanceID),
// 		s.OrgIDCondition(setting.OrgID),
// 		s.TypeCondition(setting.Type),
// 	}
// 	settingJSON, err := json.Marshal(setting.Settings)
// 	if err != nil {
// 		return 0, err
// 	}
// 	s.SetSettings(string(settingJSON)).Write(&builder)
// 	writeCondition(&builder, database.And(conditions...))

// 	stmt := builder.String()

//		return s.client.Exec(ctx, stmt, builder.Args()...)
//	}
func (s *settings) UpdateLabel(ctx context.Context, setting *domain.LabelSetting) (int64, error) {
	return s.updateSetting(ctx, setting.Setting, &setting.Settings)
}

func (s *settings) GetPasswordComplexity(ctx context.Context, instanceID string, orgID *string) (*domain.PasswordComplexitySetting, error) {
	passwordComplexitySetting := &domain.PasswordComplexitySetting{}
	var err error

	passwordComplexitySetting.Setting, err = s.Get(ctx, instanceID, orgID, domain.SettingTypePasswordComplexity)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(passwordComplexitySetting.Setting.Settings, &passwordComplexitySetting.Settings)
	if err != nil {
		return nil, err
	}

	return passwordComplexitySetting, nil
}

func (s *settings) UpdatePasswordComplexity(ctx context.Context, setting *domain.PasswordComplexitySetting) (int64, error) {
	return s.updateSetting(ctx, setting.Setting, &setting.Settings)
}

func (s *settings) GetPasswordExpiry(ctx context.Context, instanceID string, orgID *string) (*domain.PasswordExpirySetting, error) {
	passwordPolicySetting := &domain.PasswordExpirySetting{}
	var err error

	passwordPolicySetting.Setting, err = s.Get(ctx, instanceID, orgID, domain.SettingTypePasswordExpiry)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(passwordPolicySetting.Setting.Settings, &passwordPolicySetting.Settings)
	if err != nil {
		return nil, err
	}

	return passwordPolicySetting, nil
}

func (s *settings) UpdatePasswordExpiry(ctx context.Context, setting *domain.PasswordExpirySetting) (int64, error) {
	return s.updateSetting(ctx, setting.Setting, &setting.Settings)
}

func (s *settings) GetSecurity(ctx context.Context, instanceID string, orgID *string) (*domain.SecuritySetting, error) {
	securitySetting := &domain.SecuritySetting{}
	var err error

	securitySetting.Setting, err = s.Get(ctx, instanceID, orgID, domain.SettingTypeSecurity)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(securitySetting.Setting.Settings, &securitySetting.Settings)
	if err != nil {
		return nil, err
	}

	return securitySetting, nil
}

func (s *settings) UpdateSecurity(ctx context.Context, setting *domain.SecuritySetting) (int64, error) {
	return s.updateSetting(ctx, setting.Setting, &setting.Settings)
}

func (s *settings) DeleteLogin(ctx context.Context, instanceID string, orgID *string) (int64, error) {
	return s.Delete(ctx, instanceID, orgID, domain.SettingTypeLogin)
}

func (s *settings) GetLockout(ctx context.Context, instanceID string, orgID *string) (*domain.LockoutSetting, error) {
	lockoutSetting := &domain.LockoutSetting{}
	var err error

	lockoutSetting.Setting, err = s.Get(ctx, instanceID, orgID, domain.SettingTypeLockout)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(lockoutSetting.Setting.Settings, &lockoutSetting.Settings)
	if err != nil {
		return nil, err
	}

	return lockoutSetting, nil
}

func (s *settings) UpdateLockout(ctx context.Context, setting *domain.LockoutSetting) (int64, error) {
	return s.updateSetting(ctx, setting.Setting, &setting.Settings)
}

func (s *settings) GetDomain(ctx context.Context, instanceID string, orgID *string) (*domain.DomainSetting, error) {
	lockoutSetting := &domain.DomainSetting{}
	var err error

	lockoutSetting.Setting, err = s.Get(ctx, instanceID, orgID, domain.SettingTypeDomain)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(lockoutSetting.Setting.Settings, &lockoutSetting.Settings)
	if err != nil {
		return nil, err
	}

	return lockoutSetting, nil
}

func (s *settings) UpdateDomain(ctx context.Context, setting *domain.DomainSetting) (int64, error) {
	return s.updateSetting(ctx, setting.Setting, &setting.Settings)
}

// func (s *settings) updateSetting(ctx context.Context, setting *domain.Setting, settings any) (int64, error) {
// 	builder := database.StatementBuilder{}
// 	builder.WriteString(`UPDATE zitadel.settings SET `)
// 	conditions := []database.Condition{
// 		s.IDCondition(*setting.ID),
// 		s.InstanceIDCondition(setting.InstanceID),
// 		s.OrgIDCondition(setting.OrgID),
// 		s.TypeCondition(setting.Type),
// 	}
// 	settingJSON, err := json.Marshal(settings)
// 	if err != nil {
// 		return 0, err
// 	}
// 	s.SetSettings(string(settingJSON)).Write(&builder)
// 	writeCondition(&builder, database.And(conditions...))

// 	stmt := builder.String()

// 	return s.client.Exec(ctx, stmt, builder.Args()...)
// }

func (s *settings) GetOrg(ctx context.Context, instanceID string, orgID *string) (*domain.OrgSetting, error) {
	orgSetting := &domain.OrgSetting{}
	var err error

	orgSetting.Setting, err = s.Get(ctx, instanceID, orgID, domain.SettingTypeOrganization)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(orgSetting.Setting.Settings, &orgSetting.Settings)
	if err != nil {
		return nil, err
	}

	return orgSetting, nil
}

func (s *settings) UpdateOrg(ctx context.Context, setting *domain.OrgSetting) (int64, error) {
	return s.updateSetting(ctx, setting.Setting, &setting.Settings)
}

func (s *settings) updateSetting(ctx context.Context, setting *domain.Setting, settings any) (int64, error) {
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.settings SET `)
	conditions := []database.Condition{
		s.IDCondition(*setting.ID),
		s.InstanceIDCondition(setting.InstanceID),
		s.OrgIDCondition(setting.OrgID),
		s.TypeCondition(setting.Type),
	}
	settingJSON, err := json.Marshal(settings)
	if err != nil {
		return 0, err
	}
	s.SetSettings(string(settingJSON)).Write(&builder)
	writeCondition(&builder, database.And(conditions...))

	stmt := builder.String()

	return s.client.Exec(ctx, stmt, builder.Args()...)
}

const createSettingStmt = `INSERT INTO zitadel.settings` +
	// ` (instance_id, org_id, id, type, settings)` +
	` (instance_id, org_id, type, settings)` +
	// ` VALUES ($1, $2, $3, $4, $5)` +
	` VALUES ($1, $2, $3, $4)` +
	` RETURNING created_at, updated_at`

func (s *settings) Create(ctx context.Context, setting *domain.Setting) error {
	builder := database.StatementBuilder{}

	builder.AppendArgs(
		setting.InstanceID,
		setting.OrgID,
		// setting.ID,
		setting.Type,
		string(setting.Settings))
	builder.WriteString(createSettingStmt)

	return s.client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&setting.CreatedAt, &setting.UpdatedAt)
}

// func (s *settings) Update(ctx context.Context, id string, instanceID string, orgID *string, changes ...database.Change) (int64, error) {
// 	if changes == nil {
// 		return 0, errors.New("Update must contain at least one change")
// 	}
// 	builder := database.StatementBuilder{}
// 	builder.WriteString(`UPDATE zitadel.settings SET `)

// 	conditions := []database.Condition{
// 		s.IDCondition(id),
// 		s.InstanceIDCondition(instanceID),
// 		s.OrgIDCondition(orgID),
// 	}
// 	database.Changes(changes).Write(&builder)
// 	writeCondition(&builder, database.And(conditions...))

// 	stmt := builder.String()

// 	return s.client.Exec(ctx, stmt, builder.Args()...)
// }

func (s *settings) Delete(ctx context.Context, instanceID string, orgID *string, typ domain.SettingType) (int64, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(`DELETE FROM zitadel.settings`)

	conditions := []database.Condition{
		s.TypeCondition(typ),
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
