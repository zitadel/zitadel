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
	ErrMissingType       = errors.New("type must be set")
	ErrMissingOwnerType  = errors.New("owner must be set")
)

type settings struct{}

func (s settings) qualifiedTableName() string {
	return "zitadel." + s.unqualifiedTableName()
}

func (settings) unqualifiedTableName() string {
	return "settings"
}

var (
	ErrSettingObjectMustNotBeNil error = errors.New("setting object must not be nil")
	ErrLabelStateMustBeDefined   error = errors.New("label state must be defined")
)

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (s settings) IDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "id")
}

func (s settings) InstanceIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "instance_id")
}

func (s settings) OrganizationIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "organization_id")
}

func (s settings) TypeColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "type")
}

func (s settings) OwnerTypeColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "owner_type")
}

func (s settings) StateColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "state")
}

func (s settings) SettingsColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "settings")
}

func (s settings) CreatedAtColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "created_at")
}

func (s settings) UpdatedAtColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "updated_at")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (s settings) InstanceIDCondition(id string) database.Condition {
	return database.NewTextCondition(s.InstanceIDColumn(), database.TextOperationEqual, id)
}

func (s settings) OrganizationIDCondition(id *string) database.Condition {
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

func (s settings) StateCondition(typ domain.SettingState) database.Condition {
	return database.NewTextCondition(s.StateColumn(), database.TextOperationEqual, typ.String())
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

func (s settings) SetSettings(settings string) database.Change {
	return database.NewChange(s.SettingsColumn(), settings)
}

func (s settings) SetUpdatedAt(updatedAt *time.Time) database.Change {
	return database.NewChangePtr(s.UpdatedAtColumn(), updatedAt)
}

func (s settings) get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Settings, error) {
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

	return getOne[domain.Settings](ctx, client, &builder)
}

func (s settings) list(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Settings, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(querySettingStmt)

	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	options.Write(&builder)

	orderBy := database.OrderBy(s.CreatedAtColumn())
	orderBy.Write(&builder)

	return getMany[domain.Settings](ctx, client, &builder)
}

const setSettingStmt = `INSERT INTO zitadel.settings` +
	` (id, instance_id, organization_id, type, owner_type, state, settings)` +
	` VALUES ($1, $2, $3, $4, $5, $6, $7)` +
	` ON CONFLICT (id, instance_id, organization_id, type, owner_type) DO UPDATE SET`

func (s settings) set(ctx context.Context, client database.QueryExecutor, settings *domain.Settings, changes ...database.Change) error {
	builder := database.NewStatementBuilder(setSettingStmt,
		settings.ID,
		settings.InstanceID,
		settings.OrganizationID,
		settings.Type,
		settings.OwnerType,
		settings.State,
		settings.Settings,
	)
	err := database.Changes(changes).Write(builder)
	if err != nil {
		return err
	}
	builder.WriteString(` RETURNING created_at, updated_at`)

	return client.QueryRow(ctx, builder.String(), settings).Scan(&settings.CreatedAt, &settings.UpdatedAt)
}

func (s settings) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if err := s.CheckMandatoryCondtions(condition); err != nil {
		return 0, err
	}

	var builder database.StatementBuilder
	builder.WriteString(`DELETE FROM ` + s.qualifiedTableName())

	writeCondition(&builder, condition)
	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (s settings) CheckMandatoryCondtions(condition database.Condition) error {
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

// login
type loginSettings struct {
	settings
}

// -------------------------------------------------------------
// login changes
// -------------------------------------------------------------

func (l loginSettings) SetSettingFields(value *domain.LoginSettings) database.Change {
	changes := make([]db_json.JsonUpdate, 0)
	if value == nil {
		return nil
	}
	if value.AllowUserNamePassword != nil {
		changes = append(changes, l.SetAllowUserNamePassword(*value.AllowUserNamePassword))
	}
	if value.AllowRegister != nil {
		changes = append(changes, l.SetAllowRegister(*value.AllowRegister))
	}
	if value.AllowExternalIDP != nil {
		changes = append(changes, l.SetAllowExternalIDP(*value.AllowExternalIDP))
	}
	if value.ForceMFA != nil {
		changes = append(changes, l.SetForceMFA(*value.ForceMFA))
	}
	if value.ForceMFALocalOnly != nil {
		changes = append(changes, l.SetForceMFALocalOnly(*value.ForceMFALocalOnly))
	}
	if value.HidePasswordReset != nil {
		changes = append(changes, l.SetHidePasswordReset(*value.HidePasswordReset))
	}
	if value.IgnoreUnknownUsernames != nil {
		changes = append(changes, l.SetIgnoreUnknownUsernames(*value.IgnoreUnknownUsernames))
	}
	if value.AllowDomainDiscovery != nil {
		changes = append(changes, l.SetAllowDomainDiscovery(*value.AllowDomainDiscovery))
	}
	if value.DisableLoginWithEmail != nil {
		changes = append(changes, l.SetDisableLoginWithEmail(*value.DisableLoginWithEmail))
	}
	if value.DisableLoginWithPhone != nil {
		changes = append(changes, l.SetDisableLoginWithPhone(*value.DisableLoginWithPhone))
	}
	if value.PasswordlessType != nil {
		changes = append(changes, l.SetPasswordlessType(*value.PasswordlessType))
	}
	if value.DefaultRedirectURI != nil {
		changes = append(changes, l.SetDefaultRedirectURI(*value.DefaultRedirectURI))
	}
	if value.PasswordCheckLifetime != nil {
		changes = append(changes, l.SetPasswordCheckLifetime(*value.PasswordCheckLifetime))
	}
	if value.ExternalLoginCheckLifetime != nil {
		changes = append(changes, l.SetExternalLoginCheckLifetime(*value.ExternalLoginCheckLifetime))
	}
	if value.MFAInitSkipLifetime != nil {
		changes = append(changes, l.SetMFAInitSkipLifetime(*value.MFAInitSkipLifetime))
	}
	if value.SecondFactorCheckLifetime != nil {
		changes = append(changes, l.SetSecondFactorCheckLifetime(*value.SecondFactorCheckLifetime))
	}
	if value.MultiFactorCheckLifetime != nil {
		changes = append(changes, l.SetMultiFactorCheckLifetime(*value.MultiFactorCheckLifetime))
	}
	if value.DisableLoginWithEmail != nil {
		changes = append(changes, l.SetDisableLoginWithEmail(*value.DisableLoginWithEmail))
	}
	if value.MFAType != nil {
		// TODO(stebenz): determine what factors have to be removed and added
		//changes = append(changes, l.AddMFAType(0))
		//changes = append(changes, l.RemoveMFAType(0))
	}
	if value.SecondFactorTypes != nil {
		// TODO(stebenz): determine what factors have to be removed and added
		//changes = append(changes, l.AddSecondFactorTypes(0))
		//changes = append(changes, l.RemoveSecondFactorTypes(0))
	}
	return db_json.NewJsonChanges(l.SettingsColumn(), changes...)
}

func (loginSettings) SetAllowUserNamePassword(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"allowUsernamePassword"}, value)
}

func (loginSettings) SetAllowRegister(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"allowRegister"}, value)
}

func (loginSettings) SetAllowExternalIDP(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"allowExternalIdp"}, value)
}

func (loginSettings) SetForceMFA(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"forceMfa"}, value)
}

func (loginSettings) SetForceMFALocalOnly(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"forceMFALocalOnly"}, value)
}

func (loginSettings) SetHidePasswordReset(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"hidePasswordReset"}, value)
}

func (loginSettings) SetIgnoreUnknownUsernames(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"ignoreUnknownUsernames"}, value)
}

func (loginSettings) SetAllowDomainDiscovery(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"allowDomainDiscovery"}, value)
}

func (loginSettings) SetDisableLoginWithEmail(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"disableLoginWithEmail"}, value)
}

func (loginSettings) SetDisableLoginWithPhone(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"disableLoginWithPhone"}, value)
}

func (loginSettings) SetPasswordlessType(value domain.PasswordlessType) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"passwordlessType"}, value)
}

func (loginSettings) SetDefaultRedirectURI(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"defaultRedirectUri"}, value)
}

func (loginSettings) SetPasswordCheckLifetime(value time.Duration) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"passwordCheckLifetime"}, value)
}

func (loginSettings) SetExternalLoginCheckLifetime(value time.Duration) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"externalLoginCheckLifetime"}, value)
}

func (loginSettings) SetMFAInitSkipLifetime(value time.Duration) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"mfaInitSkipLifetime"}, value)
}

func (loginSettings) SetSecondFactorCheckLifetime(value time.Duration) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"secondFactorCheckLifetime"}, value)
}

func (loginSettings) SetMultiFactorCheckLifetime(value time.Duration) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"multiFactorCheckLifetime"}, value)
}

func (loginSettings) AddMFAType(value domain.MultiFactorType) db_json.JsonUpdate {
	return db_json.NewArrayChange([]string{"mfaType"}, value, false)
}

func (loginSettings) RemoveMFAType(value domain.MultiFactorType) db_json.JsonUpdate {
	return db_json.NewArrayChange([]string{"mfaType"}, value, true)
}

func (loginSettings) AddSecondFactorTypes(value domain.SecondFactorType) db_json.JsonUpdate {
	return db_json.NewArrayChange([]string{"secondFactors"}, value, false)
}

func (loginSettings) RemoveSecondFactorTypes(value domain.SecondFactorType) db_json.JsonUpdate {
	return db_json.NewArrayChange([]string{"secondFactors"}, value, true)
}

func LoginRepository() domain.LoginSettingsRepository {
	return &loginSettings{
		settings{},
	}
}

var _ domain.LoginSettingsRepository = (*loginSettings)(nil)

func (l loginSettings) Set(ctx context.Context, client database.QueryExecutor, settings *domain.LoginSettings) error {
	settings.Type = domain.SettingTypeLogin
	settings.State = domain.SettingStateActivated
	if settings.OrganizationID == nil {
		settings.OwnerType = domain.OwnerTypeOrganization
	} else {
		settings.OwnerType = domain.OwnerTypeInstance
	}

	var settingsJSON []byte
	if err := json.Unmarshal(settingsJSON, settings.LoginSettingsAttributes); err != nil {
		return err
	}
	settings.Settings.Settings = settingsJSON
	args, changes, err := setSettingChanges(settings.Settings, settingsJSON, l.SetSettingFields(settings))
	if err != nil {
		return err
	}
	return l.settings.set(ctx, client, settings.Settings, l.SetSettingFields(settings))
}

func (s *loginSettings) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.LoginSettings, error) {
	list := make([]*domain.LoginSettings, 0)

	settings, err := s.settings.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	for i := range settings {
		err = json.Unmarshal(settings[i].Settings, &list[i])
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}

func (s *loginSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.LoginSettings, error) {
	loginSettings := &domain.LoginSettings{}

	settings, err := s.settings.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(settings.Settings, &loginSettings)
	if err != nil {
		return nil, err
	}

	return loginSettings, nil
}

// label
type brandingSettings struct {
	*settings
}

// -------------------------------------------------------------
// label changes
// -------------------------------------------------------------

func (brandingSettings) SetPrimaryColor(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"primaryColorLight"}, value)
}

func (brandingSettings) SetBackgroundColor(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"backgroundColorLight"}, value)
}

func (brandingSettings) SetWarnColor(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"warnColorLight"}, value)
}

func (brandingSettings) SetFontColor(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"fontColorLight"}, value)
}

func (brandingSettings) SetPrimaryColorDark(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"primaryColorDark"}, value)
}

func (brandingSettings) SetBackgroundColorDark(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"backgroundColorDark"}, value)
}

func (brandingSettings) SetWarnColorDark(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"warnColorDark"}, value)
}

func (brandingSettings) SetFontColorDark(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"fontColorDark"}, value)
}

func (brandingSettings) SetHideLoginNameSuffix(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"hideLoginNameSuffix"}, value)
}

func (brandingSettings) SetErrorMsgPopup(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"errorMsgPopup"}, value)
}

func (brandingSettings) SetDisableWatermark(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"disableMsgPopup"}, value)
}

func (brandingSettings) SetThemeMode(value domain.BrandingPolicyThemeMode) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"themeMode"}, value)
}

func (brandingSettings) SetLogoURLLight(value *url.URL) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"logoURLLight"}, value)
}

func (brandingSettings) SetLogoURLDark(value *url.URL) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"logoUrlDark"}, value)
}

func (brandingSettings) SetIconURLLight(value *url.URL) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"iconUrlLight"}, value)
}

func (brandingSettings) SetIconURLDark(value *url.URL) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"iconUrlDark"}, value)
}

func (brandingSettings) SetFontURL(value *url.URL) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"fontUrl"}, value)
}

func BrandingSettingsRepository() domain.BrandingSettingsRepository {
	return &brandingSettings{
		settings{},
	}
}

var _ domain.BrandingSettingsRepository = (*brandingSettings)(nil)

const createLabelSettingStmt = `INSERT INTO zitadel.settings` +
	` (instance_id, organization_id, type, owner_type, label_state, settings)` +
	` VALUES ($1, $2, $3, $4, $5, $6)` +
	` ON CONFLICT (instance_id, organization_id, type, owner_type, label_state) WHERE type = 'label' DO UPDATE SET` +
	` settings = EXCLUDED.settings` +
	` RETURNING id, created_at, updated_at`

func (s *brandingSettings) Set(ctx context.Context, client database.QueryExecutor, setting *domain.BrandingSetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}

	if setting.InstanceID == "" {
		return ErrMissingInstanceID
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
		setting.BrandingState,
		string(settingJSON))

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)
}

func (s *brandingSettings) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.settings SET `)
	err := database.Changes(changes).Write(&builder)
	if err != nil {
		return 0, err
	}
	writeCondition(&builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (s *brandingSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.BrandingSetting, error) {
	labelSetting := &domain.BrandingSetting{}
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

func (s *brandingSettings) getLabel(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Setting, error) {
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

func (s *brandingSettings) Reset(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return s.Delete(ctx, client, condition)
}

func (s *settings) CreateLabel(ctx context.Context, client database.QueryExecutor, setting *domain.BrandingSetting) error {
	if setting == nil {
		return ErrSettingObjectMustNotBeNil
	}
	if setting.BrandingState == nil {
		return ErrLabelStateMustBeDefined
	}

	setting.Type = domain.SettingTypeBranding
	return setSetting(ctx, client, setting.Setting, setting)
}

func (s *brandingSettings) ActivateLabelSetting(ctx context.Context, client database.QueryExecutor, setting *domain.BrandingSetting) error {
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

func (s *brandingSettings) ActivateLabelSettingEvent(ctx context.Context, client database.QueryExecutor, condition database.Condition, uudateAt time.Time) (int64, error) {
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

func (s *settings) ActivateLabelSetting(ctx context.Context, client database.QueryExecutor, setting *domain.BrandingSetting) error {
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
	return setSetting(ctx, client, setting.Setting, setting)
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
	return setSetting(ctx, client, setting.Setting, setting)
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
	return setSetting(ctx, client, setting.Setting, setting)
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
	return setSetting(ctx, client, setting.Setting, setting)
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

	return setSetting(ctx, client, setting.Setting, setting)
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
	return setSetting(ctx, client, setting.Setting, setting)
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
	return setSetting(ctx, client, setting.Setting, setting)
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
		setting.BrandingState,
		string(setting.Settings),
		setting.CreatedAt,
		setting.UpdatedAt)

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&setting.ID, &setting.CreatedAt, &setting.UpdatedAt)
}
