package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	return database.NewColumn("settings", "organization_id")
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

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

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

const querySettingStmt = `SELECT instance_id, organization_id, id, type, owner_type, label_state, settings,` +
	` created_at, updated_at` +
	` FROM zitadel.settings`

func (s *settings) Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string, typ domain.SettingType, opts ...database.QueryOption) (*domain.Setting, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(querySettingStmt)

	// opts = append(opts, database.WithCondition(s.TypeCondition(typ)), database.WithCondition(s.InstanceIDCondition(instanceID)), database.WithCondition(s.OrgIDCondition(orgID)))
	opts = append(opts, database.WithCondition(database.And(s.TypeCondition(typ), s.InstanceIDCondition(instanceID), s.OrgIDCondition(orgID))))
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	options.Write(&builder)
	// writeCondition(&builder, database.And(cond...))

	return getOne[domain.Setting](ctx, client, &builder)
}

func (s *settings) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Setting, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(querySettingStmt)

	// if conditions != nil {
	// 	writeCondition(&builder, database.And(conditions...))
	// }

	// if opts != nil {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	options.Write(&builder)
	// }

	orderBy := database.OrderBy(s.CreatedAtColumn())
	orderBy.Write(&builder)

	return getMany[domain.Setting](ctx, client, &builder)
}

// login
type loginSettings struct {
	domain.SettingsRepository
}

func LoginRepository() domain.LoginRepository {
	return &loginSettings{
		&settings{},
	}
}

var _ domain.LoginRepository = (*loginSettings)(nil)

func (s *loginSettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.LoginSetting, changes ...database.Change) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeLogin
	settingJSON, err := json.Marshal(setting.Settings)
	if err != nil {
		return err
	}
	fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> string(settingJSON = %+v\n", string(settingJSON))
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *loginSettings) Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.LoginSetting, error) {
	loginSetting := &domain.LoginSetting{}
	var err error

	loginSetting.Setting, err = s.SettingsRepository.Get(ctx, client, instanceID, orgID, domain.SettingTypeLogin)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> loginSetting.Setting.Settings = %+v\n", string(loginSetting.Setting.Settings))
	err = json.Unmarshal(loginSetting.Setting.Settings, &loginSetting.Settings)
	if err != nil {
		return nil, err
	}
	fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> loginSetting.Settings.ForceMFA = %+v\n", loginSetting.Settings.ForceMFA)

	return loginSetting, nil
}

// label
type labelSettings struct {
	domain.SettingsRepository
}

func LabelRepository() domain.LabelRepository {
	return &labelSettings{
		&settings{},
	}
}

var _ domain.LabelRepository = (*labelSettings)(nil)

func (s *labelSettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.LabelSetting, changes ...database.Change) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeLabel
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *labelSettings) Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.LabelSetting, error) {
	labelSetting := &domain.LabelSetting{}
	var err error

	labelSetting.Setting, err = s.SettingsRepository.Get(ctx, client, instanceID, orgID, domain.SettingTypeLabel)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> labelSetting.Setting.Settings = %+v\n", labelSetting.Setting.Settings)
	err = json.Unmarshal(labelSetting.Setting.Settings, &labelSetting.Settings)
	if err != nil {
		return nil, err
	}

	return labelSetting, nil
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
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *settings) GetLabel(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string, state domain.LabelState) (*domain.LabelSetting, error) {
	labelSetting := &domain.LabelSetting{}
	var err error

	stateCond := database.WithCondition(s.LabelStateCondition(state))
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

const activatedLabelSettingStmt = `INSERT INTO zitadel.settings` +
	` (instance_id, organization_id, type, owner_type, label_state, settings)` +
	` VALUES ($1, $2, 'label', $3, 'activated', $4)` +
	` ON CONFLICT (instance_id, organization_id, type, owner_type, label_state) WHERE type = 'label' DO UPDATE SET` +
	` settings = EXCLUDED.settings` +
	` RETURNING id, created_at, updated_at`

func (s *settings) ActivateLabelSetting(ctx context.Context, client database.QueryExecutor, setting *domain.LabelSetting) error {
	builder := database.StatementBuilder{}

	settingJSON, err := json.Marshal(setting.Settings)
	if err != nil {
		return err
	}

	builder.AppendArgs(
		setting.InstanceID,
		setting.OrgID,
		setting.OwnerType,
		string(settingJSON))

	builder.WriteString(activatedLabelSettingStmt)

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)
}

func (s *settings) CreatePasswordComplexity(ctx context.Context, client database.QueryExecutor, setting *domain.PasswordComplexitySetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypePasswordComplexity
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
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

// passwordComplexity
type passwordComplexitySettings struct {
	domain.SettingsRepository
}

func PasswordComplexityRepository() domain.PasswordComplexityRepository {
	return &passwordComplexitySettings{
		&settings{},
	}
}

var _ domain.PasswordComplexityRepository = (*passwordComplexitySettings)(nil)

func (s *passwordComplexitySettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.PasswordComplexitySetting, changes ...database.Change) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypePasswordComplexity
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *passwordComplexitySettings) Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.PasswordComplexitySetting, error) {
	passwordComplexitySetting := &domain.PasswordComplexitySetting{}
	var err error

	passwordComplexitySetting.Setting, err = s.SettingsRepository.Get(ctx, client, instanceID, orgID, domain.SettingTypePasswordComplexity)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> passwordComplexitySetting.Setting.Settings = %+v\n", passwordComplexitySetting.Setting.Settings)
	err = json.Unmarshal(passwordComplexitySetting.Setting.Settings, &passwordComplexitySetting.Settings)
	if err != nil {
		return nil, err
	}

	return passwordComplexitySetting, nil
}

// passwordExpiry
type passwordExpirySettings struct {
	domain.SettingsRepository
}

func PasswordExpiryRepository() domain.PasswordExpiryRepository {
	return &passwordExpirySettings{
		&settings{},
	}
}

var _ domain.PasswordExpiryRepository = (*passwordExpirySettings)(nil)

func (s *passwordExpirySettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.PasswordExpirySetting, changes ...database.Change) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypePasswordExpiry
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *passwordExpirySettings) Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.PasswordExpirySetting, error) {
	passwordExpirySetting := &domain.PasswordExpirySetting{}
	var err error

	passwordExpirySetting.Setting, err = s.SettingsRepository.Get(ctx, client, instanceID, orgID, domain.SettingTypePasswordExpiry)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(passwordExpirySetting.Setting.Settings, &passwordExpirySetting.Settings)
	if err != nil {
		return nil, err
	}

	return passwordExpirySetting, nil
}

// lockout
type lockoutSettings struct {
	domain.SettingsRepository
}

func LockoutRepository() domain.LockoutRepository {
	return &lockoutSettings{
		&settings{},
	}
}

var _ domain.LockoutRepository = (*lockoutSettings)(nil)

func (s *lockoutSettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.LockoutSetting, changes ...database.Change) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeLockout
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *lockoutSettings) Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.LockoutSetting, error) {
	lockoutSetting := &domain.LockoutSetting{}
	var err error

	lockoutSetting.Setting, err = s.SettingsRepository.Get(ctx, client, instanceID, orgID, domain.SettingTypeLockout)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> lockoutSetting.Setting.Settings = %+v\n", lockoutSetting.Setting.Settings)
	err = json.Unmarshal(lockoutSetting.Setting.Settings, &lockoutSetting.Settings)
	if err != nil {
		return nil, err
	}

	return lockoutSetting, nil
}

// security
type securitySettings struct {
	domain.SettingsRepository
}

func SecurityRepository() domain.SecurityRepository {
	return &securitySettings{
		&settings{},
	}
}

var _ domain.SecurityRepository = (*securitySettings)(nil)

func (s *securitySettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.SecuritySetting, changes ...database.Change) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeSecurity
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *securitySettings) Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.SecuritySetting, error) {
	securitySetting := &domain.SecuritySetting{}
	var err error

	securitySetting.Setting, err = s.SettingsRepository.Get(ctx, client, instanceID, orgID, domain.SettingTypeSecurity)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> securitySetting.Setting.Settings = %+v\n", securitySetting.Setting.Settings)
	err = json.Unmarshal(securitySetting.Setting.Settings, &securitySetting.Settings)
	if err != nil {
		return nil, err
	}

	return securitySetting, nil
}

// Domain
type DomainSettings struct {
	domain.SettingsRepository
}

func DomainRepository() domain.DomainRepository {
	return &DomainSettings{
		&settings{},
	}
}

var _ domain.DomainRepository = (*DomainSettings)(nil)

func (s *DomainSettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.DomainSetting, changes ...database.Change) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeDomain
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *DomainSettings) Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.DomainSetting, error) {
	DomainSetting := &domain.DomainSetting{}
	var err error

	DomainSetting.Setting, err = s.SettingsRepository.Get(ctx, client, instanceID, orgID, domain.SettingTypeDomain)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DomainSetting.Setting.Settings = %+v\n", DomainSetting.Setting.Settings)
	err = json.Unmarshal(DomainSetting.Setting.Settings, &DomainSetting.Settings)
	if err != nil {
		return nil, err
	}

	return DomainSetting, nil
}

// organization
type organizationSettings struct {
	domain.SettingsRepository
}

func OrganizationSettingRepository() domain.OrganizationSettingRepository {
	return &organizationSettings{
		&settings{},
	}
}

var _ domain.OrganizationSettingRepository = (*organizationSettings)(nil)

func (s *organizationSettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.OrganizationSetting, changes ...database.Change) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeOrganization
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *organizationSettings) Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.OrganizationSetting, error) {
	organizationSetting := &domain.OrganizationSetting{}
	var err error

	organizationSetting.Setting, err = s.SettingsRepository.Get(ctx, client, instanceID, orgID, domain.SettingTypeOrganization)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> organizationSetting.Setting.Settings = %+v\n", organizationSetting.Setting.Settings)
	err = json.Unmarshal(organizationSetting.Setting.Settings, &organizationSetting.Settings)
	if err != nil {
		return nil, err
	}

	return organizationSetting, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *settings) UpdatePasswordComplexity(ctx context.Context, client database.QueryExecutor, setting *domain.PasswordComplexitySetting, changes ...database.Change) (int64, error) {
	return s.updateSetting(ctx, client, setting.Setting, &setting.Settings, changes...)
}

func (s *settings) CreatePasswordExpiry(ctx context.Context, client database.QueryExecutor, setting *domain.PasswordExpirySetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypePasswordExpiry
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
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
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
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
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
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
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
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

func (s *settings) CreateOrg(ctx context.Context, client database.QueryExecutor, setting *domain.OrganizationSetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeOrganization
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *settings) GetOrg(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.OrganizationSetting, error) {
	orgSetting := &domain.OrganizationSetting{}
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

func (s *settings) UpdateOrg(ctx context.Context, client database.QueryExecutor, setting *domain.OrganizationSetting, changes ...database.Change) (int64, error) {
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
	` (instance_id, organization_id, type, owner_type, label_state, settings, created_at, updated_at)` +
	` VALUES ($1, $2, $3, $4, $5, $6, $7, $8)` +
	` RETURNING id, created_at, updated_at`

func (s *settings) Create(ctx context.Context, client database.QueryExecutor, setting *domain.Setting) error {
	builder := database.NewStatementBuilder(
		eventCreateSettingStmt,
		setting.InstanceID,
		setting.OrgID,
		setting.Type,
		setting.OwnerType,
		setting.LabelState,
		string(setting.Settings),
		setting.CreatedAt,
		setting.UpdatedAt)

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)
}

// const createSettingStmt = `INSERT INTO zitadel.settings` +
// 	` (instance_id, organization_id, type, owner_type, label_state, settings)` +
// 	` VALUES ($1, $2, $3, $4, $5, $6)` +
// 	` RETURNING id, created_at, updated_at`

const createSettingStmt = `INSERT INTO zitadel.settings` +
	` (instance_id, organization_id, type, owner_type, label_state, settings)` +
	` VALUES ($1, $2, $3, $4, $5, $6)` +
	` ON CONFLICT (instance_id, organization_id, type, owner_type) WHERE type != 'label' DO UPDATE SET` +
	` settings = zitadel.settings.settings || EXCLUDED.settings::JSONB` +
	// ` settings =  EXCLUDED.settings::JSONB` +
	` RETURNING id, created_at, updated_at`

// const createSettingStmt = `INSERT INTO zitadel.settings` +
// 	` (instance_id, organization_id, type, owner_type, label_state, settings)` +
// 	` VALUES ($1, $2, $3, $4, $5, $6)`

func createSetting(ctx context.Context, client database.QueryExecutor, setting *domain.Setting, settings any) error {
	settingJSON, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	builder := database.NewStatementBuilder(
		createSettingStmt,
		setting.InstanceID,
		setting.OrgID,
		setting.Type,
		setting.OwnerType,
		setting.LabelState,
		string(settingJSON))

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
