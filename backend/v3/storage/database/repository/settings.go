package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	db_json "github.com/zitadel/zitadel/backend/v3/storage/database/json"
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

func (settings) OwnerTypeColumn() database.Column {
	return database.NewColumn("settings", "owner_type")
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

func (s settings) OwnerTypeCondition(typ domain.OwnerType) database.Condition {
	return database.NewTextCondition(s.OwnerTypeColumn(), database.TextOperationEqual, typ.String())
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

func (s settings) SetLabelSettings(changes ...db_json.JSONFieldChange) database.Change {
	return db_json.NewJsonChange(s.SettingsColumn(), changes...)
}

const querySettingStmt = `SELECT instance_id, organization_id, id, type, owner_type, label_state, settings,` +
	` created_at, updated_at` +
	` FROM zitadel.settings`

func (s *settings) Get_(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Setting, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	if !options.Condition.IsRestrictingColumn(s.LabelStateColumn()) {
		return nil, database.NewMissingConditionError(s.LabelStateColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(querySettingStmt)
	options.Write(&builder)

	return getOne[domain.Setting](ctx, client, &builder)
}

func (s *settings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Setting, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	if err := s.checkMandatoryCondtions(options.Condition); err != nil {
		return nil, err
	}

	builder := database.StatementBuilder{}
	builder.WriteString(querySettingStmt)

	options.Write(&builder)

	return getOne[domain.Setting](ctx, client, &builder)
}

// func (s *settings) Get(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string, typ domain.SettingType, opts ...database.QueryOption) (*domain.Setting, error) {
// 	builder := database.StatementBuilder{}

// 	options := new(database.QueryOpts)
// 	for _, opt := range opts {
// 		opt(options)
// 	}

// 	s.checkMandatoryCondtions(options.Condition)

// 	builder.WriteString(querySettingStmt)

// 	options.Write(&builder)

// 	return getOne[domain.Setting](ctx, client, &builder)
// }

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

// -------------------------------------------------------------
// login changes
// -------------------------------------------------------------

func (l loginSettings) SetAllowUserNamePasswordField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{allowUsernamePassword}'", value)
}

func (l loginSettings) SetAllowRegisterField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{allowRegister}'", value)
}

func (l loginSettings) SetAllowExternalIDPField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{allowExternalIdp}'", value)
}

func (l loginSettings) SetForceMFAField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{forceMfa}'", value)
}

func (l loginSettings) SetForceMFALocalOnlyField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{forceMFALocalOnly}'", value)
}

func (l loginSettings) SetHidePasswordResetField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{hidePasswordReset}'", value)
}

func (l loginSettings) SetIgnoreUnknownUsernamesField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{ignoreUnknownUsernames}'", value)
}

func (l loginSettings) SetAllowDomainDiscoveryField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{allowDomainDiscovery}'", value)
}

func (l loginSettings) SetDisableLoginWithEmailField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{disableLoginWithEmail}'", value)
}

func (l loginSettings) SetDisableLoginWithPhoneField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{disableLoginWithPhone}'", value)
}

func (l loginSettings) SetPasswordlessTypeField(value domain.PasswordlessType) db_json.JSONFieldChange {
	return db_json.NewChange("'{passwordlessType}'", value)
}

func (l loginSettings) SetDefaultRedirectURIField(value string) db_json.JSONFieldChange {
	return db_json.NewChange("'{defaultRedirectUri}'", value)
}

func (l loginSettings) SetPasswordCheckLifetimeField(value time.Duration) db_json.JSONFieldChange {
	return db_json.NewChange("'{passwordCheckLifetime}'", value)
}

func (l loginSettings) SetExternalLoginCheckLifetimeField(value time.Duration) db_json.JSONFieldChange {
	return db_json.NewChange("'{externalLoginCheckLifetime}'", value)
}

func (l loginSettings) SetMFAInitSkipLifetimeField(value time.Duration) db_json.JSONFieldChange {
	return db_json.NewChange("'{mfaInitSkipLifetime}'", value)
}

func (l loginSettings) SetSecondFactorCheckLifetimeField(value time.Duration) db_json.JSONFieldChange {
	return db_json.NewChange("'{secondFactorCheckLifetime}'", value)
}

func (l loginSettings) SetMultiFactorCheckLifetimeField(value time.Duration) db_json.JSONFieldChange {
	return db_json.NewChange("'{multiFactorCheckLifetime}'", value)
}

func (l loginSettings) SetMFATypeField(value []domain.MultiFactorType) db_json.JSONFieldChange {
	return db_json.NewChange("'{mfaType}'", value)
}

func (l loginSettings) SetSecondFactorTypesField(value []domain.SecondFactorType) db_json.JSONFieldChange {
	return db_json.NewChange("'{secondFactors}'", value)
}

func LoginRepository() domain.LoginRepository {
	return &loginSettings{
		&settings{},
	}
}

var _ domain.LoginRepository = (*loginSettings)(nil)

func (s *loginSettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.LoginSetting, changes ...database.Change) error {
	// if setting == nil {
	// 	return ErrSettingObjectMustNotBeNil
	// }
	// setting.Type = domain.SettingTypeLogin
	// settingJSON, err := json.Marshal(setting.Settings)
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("[DEBUGPRINT] [:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> string(settingJSON = %+v\n", string(settingJSON))
	// return createSetting(ctx, client, setting.Setting, &setting.Settings)
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeLockout
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *loginSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.LoginSetting, error) {
	loginSetting := &domain.LoginSetting{}
	var err error

	loginSetting.Setting, err = s.SettingsRepository.Get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[DEBUGPRINT] [settings_relational.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> loginSetting.Setting = %+v\n", string(loginSetting.Setting.Settings))
	err = json.Unmarshal(loginSetting.Setting.Settings, &loginSetting.Settings)
	if err != nil {
		return nil, err
	}

	return loginSetting, nil
}

func (s *loginSettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return s.SettingsRepository.Update_(ctx, client, condition, changes...)
}

// label
type labelSettings struct {
	domain.SettingsRepository
}

// -------------------------------------------------------------
// label changes
// -------------------------------------------------------------

func (l labelSettings) SetPrimaryColorField(value string) db_json.JSONFieldChange {
	return db_json.NewChange("'{primaryColor}'", value)
}

func (l labelSettings) SetBackgroundColorField(value string) db_json.JSONFieldChange {
	return db_json.NewChange("'{backgroundColor}'", value)
}

func (l labelSettings) SetWarnColorField(value string) db_json.JSONFieldChange {
	return db_json.NewChange("'{warnColor}'", value)
}

func (l labelSettings) SetFontColorField(value string) db_json.JSONFieldChange {
	return db_json.NewChange("'{fontColor}'", value)
}

func (l labelSettings) SetPrimaryCcolorDarkField(value string) db_json.JSONFieldChange {
	return db_json.NewChange("'{primaryColorDark}'", value)
}

func (l labelSettings) SetBackgroundColorDarkField(value string) db_json.JSONFieldChange {
	return db_json.NewChange("'{backgroundColorDark}'", value)
}

func (l labelSettings) SetWarnColorDarkField(value string) db_json.JSONFieldChange {
	return db_json.NewChange("'{warnColorDark}'", value)
}

func (l labelSettings) SetFontColorDarkField(value string) db_json.JSONFieldChange {
	return db_json.NewChange("'{fontColorDark}'", value)
}

func (l labelSettings) SetHideLoginNameSuffixField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{hideLoginNameSuffix}'", value)
}

func (l labelSettings) SetErrorMsgPopupField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{errorMsgPopup}'", value)
}

func (l labelSettings) SetDisableWatermarkField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{disableMsgPopup}'", value)
}

func (l labelSettings) SetThemeModeField(value domain.LabelPolicyThemeMode) db_json.JSONFieldChange {
	return db_json.NewChange("'{themeMode}'", value)
}

func (l labelSettings) SetLabelPolicyLightLogoURL(value *url.URL) db_json.JSONFieldChange {
	return db_json.NewChange("'{labelPolicyLightLogoURL}'", value)
}

func (l labelSettings) SetLabelPolicyDarkLogoURL(value *url.URL) db_json.JSONFieldChange {
	return db_json.NewChange("'{labelPolicyDarkLogoURL}'", value)
}

func (l labelSettings) SetLabelPolicyLightIconURL(value *url.URL) db_json.JSONFieldChange {
	return db_json.NewChange("'{labelPolicyLightIconURL}'", value)
}

func (l labelSettings) SetLabelPolicyDarkIconURL(value *url.URL) db_json.JSONFieldChange {
	return db_json.NewChange("'{labelPolicyDarkIconURL}'", value)
}

func (l labelSettings) SetLabelPolicyFontURL(value *url.URL) db_json.JSONFieldChange {
	return db_json.NewChange("'{labelPolicyLightFontURL}'", value)
}

func LabelRepository() domain.LabelRepository {
	return &labelSettings{
		&settings{},
	}
}

var _ domain.LabelRepository = (*labelSettings)(nil)

const createLabelSettingStmt = `INSERT INTO zitadel.settings` +
	` (instance_id, organization_id, type, owner_type, label_state, settings)` +
	` VALUES ($1, $2, $3, $4, $5, $6)` +
	` ON CONFLICT (instance_id, organization_id, type, owner_type, label_state) WHERE type = 'label' DO UPDATE SET` +
	` settings = EXCLUDED.settings` +
	` RETURNING id, created_at, updated_at`

func (s *labelSettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.LabelSetting) error {
	settingJSON, err := json.Marshal(setting.Settings)
	if err != nil {
		return err
	}

	builder := database.NewStatementBuilder(
		createLabelSettingStmt,
		setting.InstanceID,
		setting.OrgID,
		setting.Type,
		setting.OwnerType,
		setting.LabelState,
		string(settingJSON))

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)
}

func (s *labelSettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.settings SET `)
	database.Changes(changes).Write(&builder)
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (s *labelSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.LabelSetting, error) {
	labelSetting := &domain.LabelSetting{}
	var err error

	labelSetting.Setting, err = s.SettingsRepository.Get_(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(labelSetting.Setting.Settings, &labelSetting.Settings)
	if err != nil {
		return nil, err
	}

	return labelSetting, nil
}

func (s *labelSettings) Reset(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return s.SettingsRepository.Delete(ctx, client, condition)
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

// func (s *settings) GetLabel(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string, state domain.LabelState) (*domain.LabelSetting, error) {
// 	labelSetting := &domain.LabelSetting{}
// 	var err error

// 	stateCond := database.WithCondition(s.LabelStateCondition(state))
// 	labelSetting.Setting, err = s.Get(ctx, client, instanceID, orgID, domain.SettingTypeLabel, stateCond)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = json.Unmarshal(labelSetting.Setting.Settings, &labelSetting.Settings)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return labelSetting, nil
// }

// func (s *settings) UpdateLabel(ctx context.Context, client database.QueryExecutor, setting *domain.LabelSetting, changes ...database.Change) (int64, error) {
// 	if setting.LabelState == nil {
// 		return 0, ErrLabelStateMustBeDefined
// 	}
// 	return s.updateSetting(ctx, client, setting.Setting, &setting.Settings, changes...)
// }

// INSERT INTO zitadel.settings (instance_id, org_id, type, label_state, settings, updated_at, created_at)
// SELECT instance_id, org_id, type, $1, settings, $2, $3 FROM zitadel.settings AS copy_table
// WHERE (copy_table.type = $4) AND (copy_table.instance_id = $5) AND (copy_table.org_id IS NULL) AND (copy_table.label_state = $6)
// ON CONFLICT (instance_id, org_id, type, label_state) WHERE (type = $7)
// DO UPDATE SET (instance_id, org_id, type, label_state, settings, updated_at, created_at) = (EXCLUDED.instance_id, EXCLUDED.org_id, EXCLUDED.type, $1, EXCLUDED.settings, $2, $3)

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

// func (s *settings) GetPasswordComplexity(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.PasswordComplexitySetting, error) {
// 	passwordComplexitySetting := &domain.PasswordComplexitySetting{}
// 	var err error

// 	passwordComplexitySetting.Setting, err = s.Get(ctx, client, instanceID, orgID, domain.SettingTypePasswordComplexity)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = json.Unmarshal(passwordComplexitySetting.Setting.Settings, &passwordComplexitySetting.Settings)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return passwordComplexitySetting, nil
// }

// passwordComplexity
type passwordComplexitySettings struct {
	domain.SettingsRepository
}

// -------------------------------------------------------------
// label changes
// -------------------------------------------------------------

func (l passwordComplexitySettings) SetMinLengthField(value uint64) db_json.JSONFieldChange {
	return db_json.NewChange("'{minLength}'", value)
}

func (l passwordComplexitySettings) SetHasLowercaseField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{hasLowercase}'", value)
}

func (l passwordComplexitySettings) SetHasUppercaseField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{hasUppercase}'", value)
}

func (l passwordComplexitySettings) SetHasNumberField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{hasNumber}'", value)
}

func (l passwordComplexitySettings) SetHasSymbolField(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{hasSymbol}'", value)
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

func (s *passwordComplexitySettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.PasswordComplexitySetting, error) {
	passwordComplexitySetting := &domain.PasswordComplexitySetting{}
	var err error

	passwordComplexitySetting.Setting, err = s.SettingsRepository.Get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[DEBUGPRINT] [settings_relational.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> string(passwordComplexitySetting.Setting.Settings = %+v\n", string(passwordComplexitySetting.Setting.Settings))
	err = json.Unmarshal(passwordComplexitySetting.Setting.Settings, &passwordComplexitySetting.Settings)
	if err != nil {
		return nil, err
	}

	return passwordComplexitySetting, nil
}

func (s *passwordComplexitySettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return s.SettingsRepository.Update_(ctx, client, condition, changes...)
}

// passwordExpiry
type passwordExpirySettings struct {
	domain.SettingsRepository
}

func (l passwordExpirySettings) SetExpireWarnDays(value uint64) db_json.JSONFieldChange {
	return db_json.NewChange("'{expireWarnDays}'", value)
}

func (l passwordExpirySettings) SetMaxAgeDays(value uint64) db_json.JSONFieldChange {
	return db_json.NewChange("'{maxAgeDays}'", value)
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

func (s *passwordExpirySettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.PasswordExpirySetting, error) {
	passwordExpirySetting := &domain.PasswordExpirySetting{}
	var err error

	passwordExpirySetting.Setting, err = s.SettingsRepository.Get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(passwordExpirySetting.Setting.Settings, &passwordExpirySetting.Settings)
	if err != nil {
		return nil, err
	}

	return passwordExpirySetting, nil
}

func (s *passwordExpirySettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return s.SettingsRepository.Update_(ctx, client, condition, changes...)
}

// lockout
type lockoutSettings struct {
	domain.SettingsRepository
}

func (l lockoutSettings) SetMaxPasswordAttempts(value uint64) db_json.JSONFieldChange {
	return db_json.NewChange("'{maxPasswordAttempts}'", value)
}

func (l lockoutSettings) SetMaxOTPAttempts(value uint64) db_json.JSONFieldChange {
	return db_json.NewChange("'{maxOtpAttempts}'", value)
}

func (l lockoutSettings) SetShowLockOutFailures(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{showLockOutFailures}'", value)
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

func (s *lockoutSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.LockoutSetting, error) {
	lockoutSetting := &domain.LockoutSetting{}
	var err error

	lockoutSetting.Setting, err = s.SettingsRepository.Get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(lockoutSetting.Setting.Settings, &lockoutSetting.Settings)
	if err != nil {
		return nil, err
	}

	return lockoutSetting, nil
}

func (s *lockoutSettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return s.SettingsRepository.Update_(ctx, client, condition, changes...)
}

// security
type securitySettings struct {
	domain.SettingsRepository
}

func (l securitySettings) SetEnabled(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{enabled}'", value)
}

func (l securitySettings) SetEnableIframeEmbedding(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{enableIframe_embedding}'", value)
}

func (l securitySettings) SetAllowedOrigins(value []string) db_json.JSONFieldChange {
	return db_json.NewChange("'{allowedOrigins}'", value)
}

func (l securitySettings) SetEnableImpersonation(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{enableImpersonation}'", value)
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
	return createSetting(ctx, client, setting.Setting, &setting.Settings, changes...)
}

func (s *securitySettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.SecuritySetting, error) {
	securitySetting := &domain.SecuritySetting{}
	var err error

	securitySetting.Setting, err = s.SettingsRepository.Get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(securitySetting.Setting.Settings, &securitySetting.Settings)
	if err != nil {
		return nil, err
	}

	return securitySetting, nil
}

func (s *securitySettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return s.SettingsRepository.Update_(ctx, client, condition, changes...)
}

// domain
type domainSettings struct {
	domain.SettingsRepository
}

func (l domainSettings) SetUserLoginMustBeDomain(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{userLoginMustBeDomain}'", value)
}

func (l domainSettings) SetValidateOrgDomains(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{validateOrgDomains}'", value)
}

func (l domainSettings) SetSMTPSenderAddressMatchesInstanceDomain(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{smtpSenderAddressMatchesInstanceDomain}'", value)
}

func DomainRepository() domain.DomainRepository {
	return &domainSettings{
		&settings{},
	}
}

var _ domain.DomainRepository = (*domainSettings)(nil)

func (s *domainSettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.DomainSetting, changes ...database.Change) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeDomain
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

func (s *domainSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.DomainSetting, error) {
	DomainSetting := &domain.DomainSetting{}
	var err error

	DomainSetting.Setting, err = s.SettingsRepository.Get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(DomainSetting.Setting.Settings, &DomainSetting.Settings)
	if err != nil {
		return nil, err
	}

	return DomainSetting, nil
}

func (s *domainSettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return s.SettingsRepository.Update_(ctx, client, condition, changes...)
}

// organization
type organizationSettings struct {
	domain.SettingsRepository
}

func (l organizationSettings) SetOrganizationScopedUsernames(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{organizationScopedUsernames}'", value)
}

func (l organizationSettings) SetOldOrganizationScopedUsernames(value bool) db_json.JSONFieldChange {
	return db_json.NewChange("'{oldOrganizationScopedUsernames}'", value)
}

func (l organizationSettings) SetUsernameChanges(value []string) db_json.JSONFieldChange {
	return db_json.NewChange("'{usernameChanges}'", value)
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

func (s *organizationSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.OrganizationSetting, error) {
	organizationSetting := &domain.OrganizationSetting{}
	var err error

	organizationSetting.Setting, err = s.SettingsRepository.Get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(organizationSetting.Setting.Settings, &organizationSetting.Settings)
	if err != nil {
		return nil, err
	}

	return organizationSetting, nil
}

func (s *organizationSettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return s.SettingsRepository.Update_(ctx, client, condition, changes...)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *settings) CreatePasswordExpiry(ctx context.Context, client database.QueryExecutor, setting *domain.PasswordExpirySetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypePasswordExpiry
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

// func (s *settings) GetPasswordExpiry(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.PasswordExpirySetting, error) {
// 	passwordPolicySetting := &domain.PasswordExpirySetting{}
// 	var err error

// 	passwordPolicySetting.Setting, err = s.Get(ctx, client, instanceID, orgID, domain.SettingTypePasswordExpiry)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = json.Unmarshal(passwordPolicySetting.Setting.Settings, &passwordPolicySetting.Settings)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return passwordPolicySetting, nil
// }

func (s *settings) CreateSecurity(ctx context.Context, client database.QueryExecutor, setting *domain.SecuritySetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeSecurity
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

// func (s *settings) GetSecurity(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.SecuritySetting, error) {
// 	securitySetting := &domain.SecuritySetting{}
// 	var err error

// 	securitySetting.Setting, err = s.Get(ctx, client, instanceID, orgID, domain.SettingTypeSecurity)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = json.Unmarshal(securitySetting.Setting.Settings, &securitySetting.Settings)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return securitySetting, nil
// }

func (s *settings) CreateLockout(ctx context.Context, client database.QueryExecutor, setting *domain.LockoutSetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeLockout
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

// func (s *settings) GetLockout(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.LockoutSetting, error) {
// 	lockoutSetting := &domain.LockoutSetting{}
// 	var err error

// 	lockoutSetting.Setting, err = s.Get(ctx, client, instanceID, orgID, domain.SettingTypeLockout)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = json.Unmarshal(lockoutSetting.Setting.Settings, &lockoutSetting.Settings)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return lockoutSetting, nil
// }

func (s *settings) CreateDomain(ctx context.Context, client database.QueryExecutor, setting *domain.DomainSetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeDomain
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

// func (s *settings) GetDomain(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.DomainSetting, error) {
// 	lockoutSetting := &domain.DomainSetting{}
// 	var err error

// 	lockoutSetting.Setting, err = s.Get(ctx, client, instanceID, orgID, domain.SettingTypeDomain)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = json.Unmarshal(lockoutSetting.Setting.Settings, &lockoutSetting.Settings)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return lockoutSetting, nil
// }

func (s *settings) CreateOrg(ctx context.Context, client database.QueryExecutor, setting *domain.OrganizationSetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeOrganization
	return createSetting(ctx, client, setting.Setting, &setting.Settings)
}

// func (s *settings) GetOrg(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string) (*domain.OrganizationSetting, error) {
// 	orgSetting := &domain.OrganizationSetting{}
// 	var err error

// 	orgSetting.Setting, err = s.Get(ctx, client, instanceID, orgID, domain.SettingTypeOrganization)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = json.Unmarshal(orgSetting.Setting.Settings, &orgSetting.Settings)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return orgSetting, nil
// }

func (s *settings) Update_(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if err := s.checkMandatoryCondtions(condition); err != nil {
		return 0, err
	}
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.settings SET `)
	err := database.Changes(changes).Write(&builder)
	if err != nil {
		return 0, err
	}
	writeCondition(&builder, condition)
	fmt.Printf("[DEBUGPRINT] [settings_relational.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> builder.String() = %+v\n", builder.String())
	fmt.Printf("[DEBUGPRINT] [settings_relational.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> builder.Args()... = %+v\n", builder.Args()...)

	return client.Exec(ctx, builder.String(), builder.Args()...)
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
	` settings =  EXCLUDED.settings::JSONB `

// const createSettingStmt = `INSERT INTO zitadel.settings` +
// 	` (instance_id, organization_id, type, owner_type, label_state, settings)` +
// 	` VALUES ($1, $2, $3, $4, $5, $6)j

func createSetting(ctx context.Context, client database.QueryExecutor, setting *domain.Setting, settings any, changes ...database.Change) error {
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

	database.Changes(changes).Write(builder)

	builder.WriteString(` RETURNING id, created_at, updated_at`)

	fmt.Printf("[DEBUGPRINT] [settings_relational.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> SETTINGS CREATE builder.Args()... = %+v\n", builder.String())

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)
}

// func (s *settings) Delete(ctx context.Context, client database.QueryExecutor, instanceID string, orgID *string, typ domain.SettingType) (int64, error) {
func (s *settings) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	// options := new(database.QueryOpts)
	// for _, opt := range opts {
	// 	opt(options)
	// }

	if condition.IsRestrictingColumn(s.LabelStateColumn()) {
		return 0, database.NewMissingConditionError(s.LabelStateColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(`DELETE FROM zitadel.settings`)
	// options.Write(&builder)

	writeCondition(&builder, condition)
	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (s *settings) checkMandatoryCondtions(condition database.Condition) error {
	if !condition.IsRestrictingColumn(s.InstanceIDColumn()) {
		return database.NewMissingConditionError(s.InstanceIDColumn())
	}
	if !condition.IsRestrictingColumn(s.OrgIDColumn()) {
		return database.NewMissingConditionError(s.OrgIDColumn())
	}
	if !condition.IsRestrictingColumn(s.TypeColumn()) {
		return database.NewMissingConditionError(s.TypeColumn())
	}
	if !condition.IsRestrictingColumn(s.OwnerTypeColumn()) {
		return database.NewMissingConditionError(s.OwnerTypeColumn())
	}
	return nil
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
