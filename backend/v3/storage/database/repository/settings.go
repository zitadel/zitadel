package repository

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	db_json "github.com/zitadel/zitadel/backend/v3/storage/database/json"
)

var (
	ErrMissingInstanceID = errors.New("instance id must be set")
	ErrMissingOrgID      = errors.New("org id must be set")
	ErrMissingType       = errors.New("type must be set")
	ErrMissingOwnerType  = errors.New("owner must be set")
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

func (settings) OrganizationIDColumn() database.Column {
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
		return database.IsNull(s.OrganizationIDColumn())
	}
	return database.NewTextCondition(s.OrganizationIDColumn(), database.TextOperationEqual, *id)
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

func (s settings) SetLabelSettings(changes ...db_json.JsonUpdate) database.Change {
	return db_json.NewJsonChanges(s.SettingsColumn(), changes...)
}

func (s *settings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Setting, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	if err := s.CheckMandatoryCondtions(options.Condition); err != nil {
		return nil, err
	}

	builder := database.StatementBuilder{}
	builder.WriteString(querySettingStmt)

	options.Write(&builder)

	return getOne[domain.Setting](ctx, client, &builder)
}

func (s *settings) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Setting, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(querySettingStmt)

	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	options.Write(&builder)

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

func (l loginSettings) SetAllowUserNamePasswordField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"allowUsernamePassword"}, value)
}

func (l loginSettings) SetAllowRegisterField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"allowRegister"}, value)
}

func (l loginSettings) SetAllowExternalIDPField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"allowExternalIdp"}, value)
}

func (l loginSettings) SetForceMFAField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"forceMfa"}, value)
}

func (l loginSettings) SetForceMFALocalOnlyField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"forceMFALocalOnly"}, value)
}

func (l loginSettings) SetHidePasswordResetField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"hidePasswordReset"}, value)
}

func (l loginSettings) SetIgnoreUnknownUsernamesField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"ignoreUnknownUsernames"}, value)
}

func (l loginSettings) SetAllowDomainDiscoveryField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"allowDomainDiscovery"}, value)
}

func (l loginSettings) SetDisableLoginWithEmailField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"disableLoginWithEmail"}, value)
}

func (l loginSettings) SetDisableLoginWithPhoneField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"disableLoginWithPhone"}, value)
}

func (l loginSettings) SetPasswordlessTypeField(value domain.PasswordlessType) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"passwordlessType"}, value)
}

func (l loginSettings) SetDefaultRedirectURIField(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"defaultRedirectUri"}, value)
}

func (l loginSettings) SetPasswordCheckLifetimeField(value time.Duration) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"passwordCheckLifetime"}, value)
}

func (l loginSettings) SetExternalLoginCheckLifetimeField(value time.Duration) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"externalLoginCheckLifetime"}, value)
}

func (l loginSettings) SetMFAInitSkipLifetimeField(value time.Duration) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"mfaInitSkipLifetime"}, value)
}

func (l loginSettings) SetSecondFactorCheckLifetimeField(value time.Duration) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"secondFactorCheckLifetime"}, value)
}

func (l loginSettings) SetMultiFactorCheckLifetimeField(value time.Duration) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"multiFactorCheckLifetime"}, value)
}

func (l loginSettings) AddMFAType(value domain.MultiFactorType) db_json.JsonUpdate {
	return db_json.NewArrayChange([]string{"mfaType"}, value, false)
}

func (l loginSettings) RemoveMFAType(value domain.MultiFactorType) db_json.JsonUpdate {
	return db_json.NewArrayChange([]string{"mfaType"}, value, true)
}

func (l loginSettings) AddSecondFactorTypesField(value domain.SecondFactorType) db_json.JsonUpdate {
	return db_json.NewArrayChange([]string{"secondFactors"}, value, false)
}

func (l loginSettings) RemoveSecondFactorTypesField(value domain.SecondFactorType) db_json.JsonUpdate {
	return db_json.NewArrayChange([]string{"secondFactors"}, value, true)
}

func LoginRepository() domain.LoginRepository {
	return &loginSettings{
		&settings{},
	}
}

var _ domain.LoginRepository = (*loginSettings)(nil)

func (s *loginSettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.LoginSetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeLockout
	return createSetting(ctx, client, setting.Setting, setting)
}

func (s *loginSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.LoginSetting, error) {
	loginSetting := &domain.LoginSetting{}
	var err error

	loginSetting.Setting, err = s.SettingsRepository.Get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(loginSetting.Settings, &loginSetting)
	if err != nil {
		return nil, err
	}

	return loginSetting, nil
}

func (s *loginSettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return s.SettingsRepository.Update(ctx, client, condition, changes...)
}

// label
type labelSettings struct {
	*settings
}

// -------------------------------------------------------------
// label changes
// -------------------------------------------------------------

func (l labelSettings) SetPrimaryColorField(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"primaryColor"}, value)
}

func (l labelSettings) SetBackgroundColorField(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"backgroundColor"}, value)
}

func (l labelSettings) SetWarnColorField(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"warnColor"}, value)
}

func (l labelSettings) SetFontColorField(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"fontColor"}, value)
}

func (l labelSettings) SetPrimaryColorDarkField(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"primaryColorDark"}, value)
}

func (l labelSettings) SetBackgroundColorDarkField(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"backgroundColorDark"}, value)
}

func (l labelSettings) SetWarnColorDarkField(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"warnColorDark"}, value)
}

func (l labelSettings) SetFontColorDarkField(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"fontColorDark"}, value)
}

func (l labelSettings) SetHideLoginNameSuffixField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"hideLoginNameSuffix"}, value)
}

func (l labelSettings) SetErrorMsgPopupField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"errorMsgPopup"}, value)
}

func (l labelSettings) SetDisableWatermarkField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"disableMsgPopup"}, value)
}

func (l labelSettings) SetThemeModeField(value domain.LabelPolicyThemeMode) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"themeMode"}, value)
}

func (l labelSettings) SetLabelPolicyLightLogoURL(value *url.URL) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"labelPolicyLightLogoURL"}, value)
}

func (l labelSettings) SetLabelPolicyDarkLogoURL(value *url.URL) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"labelPolicyDarkLogoURL"}, value)
}

func (l labelSettings) SetLabelPolicyLightIconURL(value *url.URL) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"labelPolicyLightIconURL"}, value)
}

func (l labelSettings) SetLabelPolicyDarkIconURL(value *url.URL) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"labelPolicyDarkIconURL"}, value)
}

func (l labelSettings) SetLabelPolicyFontURL(value *url.URL) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"labelPolicyLightFontURL"}, value)
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
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}

	if setting.InstanceID == "" {
		return ErrMissingInstanceID
	}

	if setting.OrganizationID == nil {
		return ErrMissingOrgID
	}

	// type is always set from the calling function to that functions respective type
	// the check is left as a failsafe
	if setting.Type == domain.SettingTypeUnspecified {
		return ErrMissingType
	}

	if setting.OwnerType == domain.OwnerTypeUnspecified {
		return ErrMissingOwnerType
	}

	settingJSON, err := json.Marshal(setting)
	if err != nil {
		return err
	}

	builder := database.NewStatementBuilder(
		createLabelSettingStmt,
		setting.InstanceID,
		setting.OrganizationID,
		setting.Type,
		setting.OwnerType,
		setting.LabelState,
		string(settingJSON))

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)
}

func (s *labelSettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.settings SET `)
	err := database.Changes(changes).Write(&builder)
	if err != nil {
		return 0, err
	}
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (s *labelSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.LabelSetting, error) {
	labelSetting := &domain.LabelSetting{}
	var err error

	labelSetting.Setting, err = s.getLabel(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(labelSetting.Settings, &labelSetting)
	if err != nil {
		return nil, err
	}

	return labelSetting, nil
}

const querySettingStmt = `SELECT instance_id, organization_id, id, type, owner_type, label_state, settings,` +
	` created_at, updated_at` +
	` FROM zitadel.settings`

func (s *labelSettings) getLabel(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Setting, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	if err := s.CheckMandatoryCondtions(options.Condition); err != nil {
		return nil, err
	}

	if !options.Condition.IsRestrictingColumn(s.LabelStateColumn()) {
		return nil, database.NewMissingConditionError(s.LabelStateColumn())
	}

	var builder database.StatementBuilder
	builder.WriteString(querySettingStmt)
	options.Write(&builder)

	return getOne[domain.Setting](ctx, client, &builder)
}

func (s *labelSettings) Reset(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return s.Delete(ctx, client, condition)
}

func (s *settings) CreateLabel(ctx context.Context, client database.QueryExecutor, setting *domain.LabelSetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	if setting.LabelState == nil {
		return ErrLabelStateMustBeDefined
	}

	setting.Type = domain.SettingTypeLabel
	return createSetting(ctx, client, setting.Setting, setting)
}

func (s *labelSettings) ActivateLabelSetting(ctx context.Context, client database.QueryExecutor, setting *domain.LabelSetting) error {
	builder := database.StatementBuilder{}

	settingJSON, err := json.Marshal(setting)
	if err != nil {
		return err
	}

	builder.AppendArgs(
		setting.InstanceID,
		setting.OrganizationID,
		setting.OwnerType,
		string(settingJSON))

	builder.WriteString(activatedLabelSettingStmt)

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)
}

const activatedLabelSettingEventStmtStart = `INSERT INTO zitadel.settings` +
	` (instance_id, organization_id, type, owner_type, label_state, settings, created_at)` +
	` SELECT instance_id, organization_id, type, owner_type, 'activated', settings, $1 FROM zitadel.settings`

const activatedLabelSettingEventStmtEnd = ` ON CONFLICT (instance_id, organization_id, type, owner_type, label_state) WHERE type = 'label' DO UPDATE SET` +
	` settings = EXCLUDED.settings, updated_at = $1` +
	` RETURNING id, created_at, updated_at`

func (s *labelSettings) ActivateLabelSettingEvent(ctx context.Context, client database.QueryExecutor, condition database.Condition, uudateAt time.Time) (int64, error) {
	if !condition.IsRestrictingColumn(s.InstanceIDColumn()) {
		return 0, database.NewMissingConditionError(s.InstanceIDColumn())
	}
	if !condition.IsRestrictingColumn(s.OrganizationIDColumn()) {
		return 0, database.NewMissingConditionError(s.OrganizationIDColumn())
	}
	if !condition.IsRestrictingColumn(s.TypeColumn()) {
		return 0, database.NewMissingConditionError(s.TypeColumn())
	}
	if !condition.IsRestrictingColumn(s.OwnerTypeColumn()) {
		return 0, database.NewMissingConditionError(s.OwnerTypeColumn())
	}
	builder := database.StatementBuilder{}

	builder.WriteString(activatedLabelSettingEventStmtStart)

	builder.AppendArgs(
		uudateAt)
	writeCondition(&builder, condition)

	builder.WriteString(activatedLabelSettingEventStmtEnd)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

const activatedLabelSettingStmt = `INSERT INTO zitadel.settings` +
	` (instance_id, organization_id, type, owner_type, label_state, settings)` +
	` VALUES ($1, $2, 'label', $3, 'activated', $4)` +
	` ON CONFLICT (instance_id, organization_id, type, owner_type, label_state) WHERE type = 'label' DO UPDATE SET` +
	` settings = EXCLUDED.settings` +
	` RETURNING id, created_at, updated_at`

func (s *settings) ActivateLabelSetting(ctx context.Context, client database.QueryExecutor, setting *domain.LabelSetting) error {
	builder := database.StatementBuilder{}

	settingJSON, err := json.Marshal(setting)
	if err != nil {
		return err
	}

	builder.AppendArgs(
		setting.InstanceID,
		setting.OrganizationID,
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
	return createSetting(ctx, client, setting.Setting, setting)
}

// passwordComplexity
type passwordComplexitySettings struct {
	domain.SettingsRepository
}

// -------------------------------------------------------------
// label changes
// -------------------------------------------------------------

func (l passwordComplexitySettings) SetMinLengthField(value uint64) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"minLength"}, value)
}

func (l passwordComplexitySettings) SetHasLowercaseField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"hasLowercase"}, value)
}

func (l passwordComplexitySettings) SetHasUppercaseField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"hasUppercase"}, value)
}

func (l passwordComplexitySettings) SetHasNumberField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"hasNumber"}, value)
}

func (l passwordComplexitySettings) SetHasSymbolField(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"hasSymbol"}, value)
}

func PasswordComplexityRepository() domain.PasswordComplexityRepository {
	return &passwordComplexitySettings{
		&settings{},
	}
}

var _ domain.PasswordComplexityRepository = (*passwordComplexitySettings)(nil)

func (s *passwordComplexitySettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.PasswordComplexitySetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypePasswordComplexity
	return createSetting(ctx, client, setting.Setting, setting)
}

func (s *passwordComplexitySettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.PasswordComplexitySetting, error) {
	passwordComplexitySetting := &domain.PasswordComplexitySetting{}
	var err error

	passwordComplexitySetting.Setting, err = s.SettingsRepository.Get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(passwordComplexitySetting.Settings, &passwordComplexitySetting)
	if err != nil {
		return nil, err
	}

	return passwordComplexitySetting, nil
}

func (s *passwordComplexitySettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return s.SettingsRepository.Update(ctx, client, condition, changes...)
}

func (s *passwordComplexitySettings) Reset(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return s.Delete(ctx, client, condition)
}

// passwordExpiry
type passwordExpirySettings struct {
	domain.SettingsRepository
}

func (l passwordExpirySettings) SetExpireWarnDays(value uint64) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"expireWarnDays"}, value)
}

func (l passwordExpirySettings) SetMaxAgeDays(value uint64) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"maxAgeDays"}, value)
}

func PasswordExpiryRepository() domain.PasswordExpiryRepository {
	return &passwordExpirySettings{
		&settings{},
	}
}

var _ domain.PasswordExpiryRepository = (*passwordExpirySettings)(nil)

func (s *passwordExpirySettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.PasswordExpirySetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypePasswordExpiry
	return createSetting(ctx, client, setting.Setting, setting)
}

func (s *passwordExpirySettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.PasswordExpirySetting, error) {
	passwordExpirySetting := &domain.PasswordExpirySetting{}
	var err error

	passwordExpirySetting.Setting, err = s.SettingsRepository.Get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(passwordExpirySetting.Settings, &passwordExpirySetting)
	if err != nil {
		return nil, err
	}

	return passwordExpirySetting, nil
}

func (s *passwordExpirySettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return s.SettingsRepository.Update(ctx, client, condition, changes...)
}

func (s *passwordExpirySettings) Reset(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return s.Delete(ctx, client, condition)
}

// lockout
type lockoutSettings struct {
	domain.SettingsRepository
}

func (l lockoutSettings) SetMaxPasswordAttempts(value uint64) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"maxPasswordAttempts"}, value)
}

func (l lockoutSettings) SetMaxOTPAttempts(value uint64) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"maxOtpAttempts"}, value)
}

func (l lockoutSettings) SetShowLockOutFailures(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"showLockOutFailures"}, value)
}

func LockoutRepository() domain.LockoutRepository {
	return &lockoutSettings{
		&settings{},
	}
}

var _ domain.LockoutRepository = (*lockoutSettings)(nil)

func (s *lockoutSettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.LockoutSetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeLockout
	return createSetting(ctx, client, setting.Setting, setting)
}

func (s *lockoutSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.LockoutSetting, error) {
	lockoutSetting := &domain.LockoutSetting{}
	var err error

	lockoutSetting.Setting, err = s.SettingsRepository.Get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(lockoutSetting.Settings, &lockoutSetting)
	if err != nil {
		return nil, err
	}

	return lockoutSetting, nil
}

func (s *lockoutSettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return s.SettingsRepository.Update(ctx, client, condition, changes...)
}

func (s *lockoutSettings) Reset(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return s.Delete(ctx, client, condition)
}

// security
type securitySettings struct {
	domain.SettingsRepository
}

func (l securitySettings) SetEnabled(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"enabled"}, value)
}

func (l securitySettings) SetEnableIframeEmbedding(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"enableIframe_embedding"}, value)
}

func (l securitySettings) AddAllowedOrigins(value string) db_json.JsonUpdate {
	return db_json.NewArrayChange([]string{"allowedOrigins"}, value, false)
}

func (l securitySettings) RemoveAllowedOrigins(value string) db_json.JsonUpdate {
	return db_json.NewArrayChange([]string{"allowedOrigins"}, value, true)
}

func (l securitySettings) SetEnableImpersonation(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"enableImpersonation"}, value)
}

func SecurityRepository() domain.SecurityRepository {
	return &securitySettings{
		&settings{},
	}
}

var _ domain.SecurityRepository = (*securitySettings)(nil)

func (s *securitySettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.SecuritySetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeSecurity

	// return createSetting(ctx, client, setting.Setting, &setting.Settings)
	return createSetting(ctx, client, setting.Setting, setting)
}

func (s *securitySettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.SecuritySetting, error) {
	securitySetting := &domain.SecuritySetting{}
	var err error

	securitySetting.Setting, err = s.SettingsRepository.Get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(securitySetting.Settings, &securitySetting)
	if err != nil {
		return nil, err
	}

	return securitySetting, nil
}

func (s *securitySettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return s.SettingsRepository.Update(ctx, client, condition, changes...)
}

func (s *securitySettings) Reset(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return s.Delete(ctx, client, condition)
}

const setSecuritySettingEventStmt = `INSERT INTO zitadel.settings` +
	` (instance_id, organization_id, type, owner_type, settings, created_at)` +
	` VALUES ($1, $2, 'security', $3, $4, $5)` +
	` ON CONFLICT (instance_id, organization_id, type, owner_type) WHERE type != 'label' DO UPDATE Set `

func (s *securitySettings) SetEvent(ctx context.Context, client database.QueryExecutor, setting *domain.SecuritySetting, changes ...database.Change) (int64, error) {
	if setting == nil {
		return 0, ErrSettingObjectMustNotBeNil
	}

	settingJSON, err := json.Marshal(setting)
	if err != nil {
		return 0, err
	}

	builder := database.NewStatementBuilder(
		setSecuritySettingEventStmt,
		setting.InstanceID,
		setting.OrganizationID,
		setting.OwnerType,
		string(settingJSON),
		setting.CreatedAt)

	err = database.Changes(changes).Write(builder)
	if err != nil {
		return 0, err
	}

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// domain
type domainSettings struct {
	domain.SettingsRepository
}

func (l domainSettings) SetUserLoginMustBeDomain(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"userLoginMustBeDomain"}, value)
}

func (l domainSettings) SetValidateOrgDomains(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"validateOrgDomains"}, value)
}

func (l domainSettings) SetSMTPSenderAddressMatchesInstanceDomain(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"smtpSenderAddressMatchesInstanceDomain"}, value)
}

func DomainRepository() domain.DomainRepository {
	return &domainSettings{
		&settings{},
	}
}

var _ domain.DomainRepository = (*domainSettings)(nil)

func (s *domainSettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.DomainSetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeDomain
	return createSetting(ctx, client, setting.Setting, setting)
}

func (s *domainSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.DomainSetting, error) {
	DomainSetting := &domain.DomainSetting{}
	var err error

	DomainSetting.Setting, err = s.SettingsRepository.Get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(DomainSetting.Settings, &DomainSetting)
	if err != nil {
		return nil, err
	}

	return DomainSetting, nil
}

func (s *domainSettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return s.SettingsRepository.Update(ctx, client, condition, changes...)
}

func (s *domainSettings) Reset(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return s.Delete(ctx, client, condition)
}

// organization
type organizationSettings struct {
	domain.SettingsRepository
}

func (l organizationSettings) SetOrganizationScopedUsernames(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"organizationScopedUsernames"}, value)
}

func OrganizationSettingRepository() domain.OrganizationSettingRepository {
	return &organizationSettings{
		&settings{},
	}
}

var _ domain.OrganizationSettingRepository = (*organizationSettings)(nil)

func (s *organizationSettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.OrganizationSetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	setting.Type = domain.SettingTypeOrganization
	return createSetting(ctx, client, setting.Setting, setting)
}

func (s *organizationSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.OrganizationSetting, error) {
	organizationSetting := &domain.OrganizationSetting{}
	var err error

	organizationSetting.Setting, err = s.SettingsRepository.Get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(organizationSetting.Settings, &organizationSetting)
	if err != nil {
		return nil, err
	}

	return organizationSetting, nil
}

func (s *organizationSettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return s.SettingsRepository.Update(ctx, client, condition, changes...)
}

func (s *organizationSettings) Reset(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return s.Delete(ctx, client, condition)
}

const setOrganizationSettingEventStmt = `INSERT INTO zitadel.settings` +
	` (instance_id, organization_id, type, owner_type, settings, created_at)` +
	` VALUES ($1, $2, 'organization', $3, $4, $5)` +
	` ON CONFLICT (instance_id, organization_id, type, owner_type) WHERE type != 'label'` +
	` DO UPDATE SET settings = EXCLUDED.settings, updated_at = $5`

func (s *organizationSettings) SetEvent(ctx context.Context, client database.QueryExecutor, setting *domain.OrganizationSetting) (int64, error) {
	if setting == nil {
		return 0, ErrSettingObjectMustNotBeNil
	}

	settingJSON, err := json.Marshal(setting)
	if err != nil {
		return 0, err
	}

	builder := database.NewStatementBuilder(
		setOrganizationSettingEventStmt,
		setting.InstanceID,
		setting.OrganizationID,
		setting.OwnerType,
		string(settingJSON),
		setting.CreatedAt)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (s *settings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if err := s.CheckMandatoryCondtions(condition); err != nil {
		return 0, err
	}
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.settings SET `)
	err := database.Changes(changes).Write(&builder)
	if err != nil {
		return 0, err
	}
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// func (s *settings) updateSetting(ctx context.Context, client database.QueryExecutor, setting *domain.Setting, settings any, changes ...database.Change) (int64, error) {
// 	builder := database.StatementBuilder{}
// 	builder.WriteString(`UPDATE zitadel.settings SET `)
// 	conditions := []database.Condition{
// 		s.IDCondition(setting.ID),
// 		s.InstanceIDCondition(setting.InstanceID),
// 		s.OrgIDCondition(setting.OrganizationID),
// 		s.TypeCondition(setting.Type),
// 	}
// 	settingJSON, err := json.Marshal(settings)
// 	if err != nil {
// 		return 0, err
// 	}

// 	if changes == nil || !database.Changes(changes).IsOnColumn(s.UpdatedAtColumn()) {
// 		// if updated_at is not set, then explicitly set it to NULL so that the db trigger sets it to NOW()
// 		changes = append(changes, s.SetSettings(string(settingJSON)), s.SetUpdatedAt(nil))
// 	} else {
// 		changes = append(changes, s.SetSettings(string(settingJSON)))
// 	}

// 	err = database.Changes(changes).Write(&builder)
// 	if err != nil {
// 		return 0, err
// 	}

// 	writeCondition(&builder, database.And(conditions...))

// 	stmt := builder.String()

// 	return client.Exec(ctx, stmt, builder.Args()...)
// / }

const eventCreateSettingStmt = `INSERT INTO zitadel.settings` +
	` (instance_id, organization_id, type, owner_type, label_state, settings, created_at, updated_at)` +
	` VALUES ($1, $2, $3, $4, $5, $6, $7, $8)` +
	` RETURNING id, created_at, updated_at`

func (s *settings) Create(ctx context.Context, client database.QueryExecutor, setting *domain.Setting) error {
	builder := database.NewStatementBuilder(
		eventCreateSettingStmt,
		setting.InstanceID,
		setting.OrganizationID,
		setting.Type,
		setting.OwnerType,
		setting.LabelState,
		string(setting.Settings),
		setting.CreatedAt,
		setting.UpdatedAt)

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)
}

const createSettingStmt = `INSERT INTO zitadel.settings` +
	` (instance_id, organization_id, type, owner_type, label_state, settings)` +
	` VALUES ($1, $2, $3, $4, $5, $6)` +
	` ON CONFLICT (instance_id, organization_id, type, owner_type) WHERE type != 'label' DO UPDATE SET` +
	` settings =  EXCLUDED.settings::JSONB `

func createSetting(ctx context.Context, client database.QueryExecutor, setting *domain.Setting, settings any, changes ...database.Change) error {
	if setting.InstanceID == "" {
		return ErrMissingInstanceID
	}

	if setting.OrganizationID == nil {
		return ErrMissingOrgID
	}

	// type is always set from the calling function to that functions respective type
	// the check is left as a failsafe
	if setting.Type == domain.SettingTypeUnspecified {
		return ErrMissingType
	}

	if setting.OwnerType == domain.OwnerTypeUnspecified {
		return ErrMissingOwnerType
	}

	settingJSON, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	builder := database.NewStatementBuilder(
		createSettingStmt,
		setting.InstanceID,
		setting.OrganizationID,
		setting.Type,
		setting.OwnerType,
		// TODO
		setting.LabelState,
		string(settingJSON))

	err = database.Changes(changes).Write(builder)
	if err != nil {
		return err
	}

	builder.WriteString(` RETURNING id, created_at, updated_at`)

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)
}

func (s *settings) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	// options := new(database.QueryOpts)
	// for _, opt := range opts {
	// 	opt(options)
	// }

	// if condition.IsRestrictingColumn(s.LabelStateColumn()) {
	// 	return 0, database.NewMissingConditionError(s.LabelStateColumn())
	// }
	// TODO label
	if err := s.CheckMandatoryCondtions(condition); err != nil {
		return 0, err
	}

	var builder database.StatementBuilder
	builder.WriteString(`DELETE FROM zitadel.settings`)
	// options.Write(&builder)

	writeCondition(&builder, condition)
	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (s *settings) CheckMandatoryCondtions(condition database.Condition) error {
	if !condition.IsRestrictingColumn(s.InstanceIDColumn()) {
		return database.NewMissingConditionError(s.InstanceIDColumn())
	}
	if !condition.IsRestrictingColumn(s.OrganizationIDColumn()) {
		return database.NewMissingConditionError(s.OrganizationIDColumn())
	}
	if !condition.IsRestrictingColumn(s.TypeColumn()) {
		return database.NewMissingConditionError(s.TypeColumn())
	}
	if !condition.IsRestrictingColumn(s.OwnerTypeColumn()) {
		return database.NewMissingConditionError(s.OwnerTypeColumn())
	}
	return nil
}
