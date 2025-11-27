package repository

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	db_json "github.com/zitadel/zitadel/backend/v3/storage/database/json"
)

type settings struct{}

func (s settings) qualifiedTableName() string {
	return "zitadel." + s.unqualifiedTableName()
}

func (settings) unqualifiedTableName() string {
	return "settings"
}

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
	builder, err := s.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return getOne[domain.Settings](ctx, client, builder)
}

func (s settings) list(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Settings, error) {
	builder, err := s.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return getMany[domain.Settings](ctx, client, builder)
}

func (s settings) set(ctx context.Context, client database.QueryExecutor, settings *domain.Settings, changes ...database.Change) error {
	var (
		id, createdAt, updatedAt any = database.DefaultInstruction, database.DefaultInstruction, database.DefaultInstruction
	)

	if settings.ID != "" {
		id = settings.ID
	}
	if !settings.CreatedAt.IsZero() {
		createdAt = settings.CreatedAt
	}
	if !settings.UpdatedAt.IsZero() {
		updatedAt = settings.UpdatedAt
		changes = append(changes, s.SetUpdatedAt(&settings.UpdatedAt))
	}

	builder := database.NewStatementBuilder(`INSERT INTO ` + s.qualifiedTableName())
	builder.WriteString(` (id, instance_id, organization_id, type, state, settings, created_at, updated_at) VALUES (`)
	builder.WriteArgs(
		id,
		settings.InstanceID,
		settings.OrganizationID,
		settings.Type.String(),
		settings.State.String(),
		settings.Settings,
		createdAt,
		updatedAt,
	)
	builder.WriteString(`) ON CONFLICT (instance_id, organization_id, type, state) DO UPDATE SET `)
	err := database.Changes(changes).Write(builder)
	if err != nil {
		return err
	}
	builder.WriteString(` RETURNING id, created_at, updated_at`)
	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&settings.ID, &settings.CreatedAt, &settings.UpdatedAt)
}

func (s settings) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if err := checkRestrictingColumns(condition, s.UniqueColumns()...); err != nil {
		return 0, err
	}

	var builder database.StatementBuilder
	builder.WriteString(`DELETE FROM ` + s.qualifiedTableName())

	writeCondition(&builder, condition)
	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// PrimaryKeyColumns implements [domain.Repository].
func (s settings) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		s.InstanceIDColumn(),
		s.IDColumn(),
	}
}

func (s settings) PrimaryKeyCondition(instanceID, id string) database.Condition {
	return database.And(
		s.InstanceIDCondition(instanceID),
		s.IDCondition(id),
	)
}

func (s settings) UniqueColumns() []database.Column {
	return []database.Column{
		s.InstanceIDColumn(),
		s.OrganizationIDColumn(),
		s.TypeColumn(),
		s.StateColumn(),
	}
}

func (s settings) UniqueCondition(instanceID string, orgID *string, typ domain.SettingType, state domain.SettingState) database.Condition {
	return database.And(
		s.InstanceIDCondition(instanceID),
		s.OrganizationIDCondition(orgID),
		s.TypeCondition(typ),
		s.StateCondition(state),
	)
}

// -------------------------------------------------------------
// login
// -------------------------------------------------------------
type loginSettings struct {
	settings
}

func (l loginSettings) SetSettingFields(value domain.LoginSettingsAttributes) database.Change {
	changes := make([]db_json.JsonUpdate, 0)
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
		changes = append(changes, l.SetMFAType(value.MFAType))
	}
	if value.SecondFactorTypes != nil {
		changes = append(changes, l.SetSecondFactorTypes(value.SecondFactorTypes))
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
	return db_json.NewFieldChange([]string{"forceMfaLocalOnly"}, value)
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

func (loginSettings) SetMFAType(values []domain.MultiFactorType) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"mfaType"}, values)
}

func (l loginSettings) AddMFAType(value domain.MultiFactorType) database.Change {
	return db_json.NewJsonChanges(l.SettingsColumn(),
		db_json.NewArrayChange([]string{"mfaType"}, value, false),
	)
}

func (l loginSettings) RemoveMFAType(value domain.MultiFactorType) database.Change {
	return db_json.NewJsonChanges(l.SettingsColumn(),
		db_json.NewArrayChange([]string{"mfaType"}, value, true),
	)
}

func (loginSettings) SetSecondFactorTypes(values []domain.SecondFactorType) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"secondFactors"}, values)
}

func (l loginSettings) AddSecondFactorTypes(value domain.SecondFactorType) database.Change {
	return db_json.NewJsonChanges(l.SettingsColumn(),
		db_json.NewArrayChange([]string{"secondFactors"}, value, false),
	)
}

func (l loginSettings) RemoveSecondFactorTypes(value domain.SecondFactorType) database.Change {
	return db_json.NewJsonChanges(l.SettingsColumn(),
		db_json.NewArrayChange([]string{"secondFactors"}, value, true),
	)
}

func LoginSettingsRepository() domain.LoginSettingsRepository {
	return &loginSettings{
		settings{},
	}
}

var _ domain.LoginSettingsRepository = (*loginSettings)(nil)

func (s loginSettings) Set(ctx context.Context, client database.QueryExecutor, settings *domain.LoginSettings) error {
	settings.Type = domain.SettingTypeLogin
	settings.State = domain.SettingStateActive

	settingsJSON, err := json.Marshal(settings.LoginSettingsAttributes)
	if err != nil {
		return err
	}
	settings.Settings.Settings = settingsJSON
	if err := s.set(ctx, client, &settings.Settings, s.SetSettingFields(settings.LoginSettingsAttributes)); err != nil {
		return err
	}
	settings.Settings.Settings = nil
	return nil
}

func (s loginSettings) SetColumns(ctx context.Context, client database.QueryExecutor, settings *domain.Settings, changes ...database.Change) error {
	settings.Type = domain.SettingTypeLogin
	settings.State = domain.SettingStateActive

	return s.set(ctx, client, settings, changes...)
}

func (s loginSettings) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.LoginSettings, error) {
	settings, err := s.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	list := make([]*domain.LoginSettings, len(settings))
	for i := range settings {
		list[i] = &domain.LoginSettings{Settings: *settings[i]}
		if err := json.Unmarshal(list[i].Settings.Settings, &list[i].LoginSettingsAttributes); err != nil {
			return nil, err
		}
		list[i].Settings.Settings = nil
	}
	return list, nil
}

func (s loginSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.LoginSettings, error) {
	settings, err := s.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	loginSettings := &domain.LoginSettings{Settings: *settings}
	if err = json.Unmarshal(loginSettings.Settings.Settings, &loginSettings.LoginSettingsAttributes); err != nil {
		return nil, err
	}
	loginSettings.Settings.Settings = nil
	return loginSettings, nil
}

// -------------------------------------------------------------
// branding
// -------------------------------------------------------------

type brandingSettings struct {
	settings
}

func (b brandingSettings) SetSettingFields(value domain.BrandingSettingsAttributes) database.Change {
	changes := make([]db_json.JsonUpdate, 0)
	if value.PrimaryColorLight != nil {
		changes = append(changes, b.SetPrimaryColorLight(*value.PrimaryColorLight))
	}
	if value.BackgroundColorLight != nil {
		changes = append(changes, b.SetBackgroundColorLight(*value.BackgroundColorLight))
	}
	if value.WarnColorLight != nil {
		changes = append(changes, b.SetWarnColorLight(*value.WarnColorLight))
	}
	if value.FontColorLight != nil {
		changes = append(changes, b.SetFontColorLight(*value.FontColorLight))
	}
	if value.PrimaryColorDark != nil {
		changes = append(changes, b.SetPrimaryColorDark(*value.PrimaryColorDark))
	}
	if value.BackgroundColorDark != nil {
		changes = append(changes, b.SetBackgroundColorDark(*value.BackgroundColorDark))
	}
	if value.WarnColorDark != nil {
		changes = append(changes, b.SetWarnColorDark(*value.WarnColorDark))
	}
	if value.FontColorDark != nil {
		changes = append(changes, b.SetFontColorDark(*value.FontColorDark))
	}
	if value.HideLoginNameSuffix != nil {
		changes = append(changes, b.SetHideLoginNameSuffix(*value.HideLoginNameSuffix))
	}
	if value.ErrorMsgPopup != nil {
		changes = append(changes, b.SetErrorMsgPopup(*value.ErrorMsgPopup))
	}
	if value.DisableWatermark != nil {
		changes = append(changes, b.SetDisableWatermark(*value.DisableWatermark))
	}
	if value.ThemeMode != nil {
		changes = append(changes, b.SetThemeMode(*value.ThemeMode))
	}
	if value.LogoURLLight != nil {
		changes = append(changes, b.SetLogoURLLight(*value.LogoURLLight))
	}
	if value.LogoURLDark != nil {
		changes = append(changes, b.SetLogoURLDark(*value.LogoURLDark))
	}
	if value.IconURLLight != nil {
		changes = append(changes, b.SetIconURLLight(*value.IconURLLight))
	}
	if value.IconURLDark != nil {
		changes = append(changes, b.SetIconURLDark(*value.IconURLDark))
	}
	if value.FontURL != nil {
		changes = append(changes, b.SetFontURL(*value.FontURL))
	}
	return db_json.NewJsonChanges(b.SettingsColumn(), changes...)
}

func (brandingSettings) SetPrimaryColorLight(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"primaryColorLight"}, value)
}

func (brandingSettings) SetBackgroundColorLight(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"backgroundColorLight"}, value)
}

func (brandingSettings) SetWarnColorLight(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"warnColorLight"}, value)
}

func (brandingSettings) SetFontColorLight(value string) db_json.JsonUpdate {
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

func (brandingSettings) SetLogoURLLight(value url.URL) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"logoURLLight"}, value)
}

func (brandingSettings) SetLogoURLDark(value url.URL) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"logoUrlDark"}, value)
}

func (brandingSettings) SetIconURLLight(value url.URL) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"iconUrlLight"}, value)
}

func (brandingSettings) SetIconURLDark(value url.URL) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"iconUrlDark"}, value)
}

func (brandingSettings) SetFontURL(value url.URL) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"fontUrl"}, value)
}

func BrandingSettingsRepository() domain.BrandingSettingsRepository {
	return &brandingSettings{
		settings{},
	}
}

var _ domain.BrandingSettingsRepository = (*brandingSettings)(nil)

func (s brandingSettings) Set(ctx context.Context, client database.QueryExecutor, settings *domain.BrandingSettings) (err error) {
	settings.Type = domain.SettingTypeBranding
	settings.State = domain.SettingStatePreview

	settingsJSON, err := json.Marshal(settings.BrandingSettingsAttributes)
	if err != nil {
		return err
	}
	settings.Settings.Settings = settingsJSON
	if err := s.set(ctx, client, &settings.Settings, s.SetSettingFields(settings.BrandingSettingsAttributes)); err != nil {
		return err
	}
	settings.Settings.Settings = nil
	return nil
}

func (s brandingSettings) SetColumns(ctx context.Context, client database.QueryExecutor, settings *domain.Settings, changes ...database.Change) error {
	settings.Type = domain.SettingTypeBranding
	settings.State = domain.SettingStatePreview

	return s.set(ctx, client, settings, changes...)
}

func (s brandingSettings) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.BrandingSettings, error) {
	settings, err := s.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	list := make([]*domain.BrandingSettings, len(settings))
	for i := range settings {
		list[i] = &domain.BrandingSettings{Settings: *settings[i]}
		if err := json.Unmarshal(list[i].Settings.Settings, &list[i].BrandingSettingsAttributes); err != nil {
			return nil, err
		}
		list[i].Settings.Settings = nil
	}
	return list, nil
}

func (s brandingSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.BrandingSettings, error) {
	settings, err := s.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	brandingSettings := &domain.BrandingSettings{Settings: *settings}
	if err = json.Unmarshal(brandingSettings.Settings.Settings, &brandingSettings.BrandingSettingsAttributes); err != nil {
		return nil, err
	}
	brandingSettings.Settings.Settings = nil
	return brandingSettings, nil
}

const queryActivateBrandingSettingsStmt = `SELECT instance_id, organization_id, id, type, created_at, updated_at, settings FROM `

func (s brandingSettings) Activate(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return s.activateAt(ctx, client, condition, time.Now())
}

func (s brandingSettings) ActivateAt(ctx context.Context, client database.QueryExecutor, condition database.Condition, updatedAt time.Time) (int64, error) {
	return s.activateAt(ctx, client, condition, updatedAt)
}

func (s brandingSettings) activateAt(ctx context.Context, client database.QueryExecutor, condition database.Condition, updatedAt time.Time) (int64, error) {
	condition = database.And(
		condition,
		s.TypeCondition(domain.SettingTypeBranding),
		s.StateCondition(domain.SettingStatePreview),
	)
	// if you want to update the roles you have to have a unique condition, otherwise multiple branding settings get activated
	if err := checkRestrictingColumns(condition, s.UniqueColumns()...); err != nil {
		return 0, err
	}

	// the statement begins with getting settings in preview state
	builder := database.NewStatementBuilder("WITH preview AS (")
	builder.WriteString(queryActivateBrandingSettingsStmt + s.qualifiedTableName())
	writeCondition(builder, condition)
	builder.WriteString(" ) ")

	// insert or update settings as active
	builder.WriteString(`INSERT INTO ` + s.qualifiedTableName())
	builder.WriteString(` (instance_id, organization_id, id, type, state, created_at, updated_at, settings)`)

	builder.WriteString(` SELECT preview.instance_id, preview.organization_id, preview.id, preview.type, `)
	builder.WriteArg(domain.SettingStateActive)
	builder.WriteString(`, preview.created_at, `)
	builder.WriteArg(updatedAt)
	builder.WriteString(`, preview.settings FROM preview`)

	builder.WriteString(` ON CONFLICT (instance_id, organization_id, type, state) DO UPDATE SET settings = EXCLUDED.settings, updated_at = `)
	builder.WriteArg(updatedAt)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// password complexity
// -------------------------------------------------------------
type passwordComplexitySettings struct {
	settings
}

func (s passwordComplexitySettings) SetSettingFields(value domain.PasswordComplexitySettingsAttributes) database.Change {
	changes := make([]db_json.JsonUpdate, 0)
	if value.MinLength != nil {
		changes = append(changes, s.SetMinLength(*value.MinLength))
	}
	if value.HasLowercase != nil {
		changes = append(changes, s.SetHasLowercase(*value.HasLowercase))
	}
	if value.HasUppercase != nil {
		changes = append(changes, s.SetHasUppercase(*value.HasUppercase))
	}
	if value.HasNumber != nil {
		changes = append(changes, s.SetHasNumber(*value.HasNumber))
	}
	if value.HasSymbol != nil {
		changes = append(changes, s.SetHasSymbol(*value.HasSymbol))
	}
	return db_json.NewJsonChanges(s.SettingsColumn(), changes...)
}

func (passwordComplexitySettings) SetMinLength(value uint64) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"minLength"}, value)
}

func (passwordComplexitySettings) SetHasLowercase(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"hasLowercase"}, value)
}

func (passwordComplexitySettings) SetHasUppercase(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"hasUppercase"}, value)
}

func (passwordComplexitySettings) SetHasNumber(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"hasNumber"}, value)
}

func (passwordComplexitySettings) SetHasSymbol(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"hasSymbol"}, value)
}

func PasswordComplexitySettingsRepository() domain.PasswordComplexitySettingsRepository {
	return &passwordComplexitySettings{
		settings{},
	}
}

var _ domain.PasswordComplexitySettingsRepository = (*passwordComplexitySettings)(nil)

func (s passwordComplexitySettings) Set(ctx context.Context, client database.QueryExecutor, settings *domain.PasswordComplexitySettings) (err error) {
	settings.Type = domain.SettingTypePasswordComplexity
	settings.State = domain.SettingStateActive

	settingsJSON, err := json.Marshal(settings.PasswordComplexitySettingsAttributes)
	if err != nil {
		return err
	}
	settings.Settings.Settings = settingsJSON
	if err := s.set(ctx, client, &settings.Settings, s.SetSettingFields(settings.PasswordComplexitySettingsAttributes)); err != nil {
		return err
	}
	settings.Settings.Settings = nil
	return nil
}

func (s passwordComplexitySettings) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.PasswordComplexitySettings, error) {
	settings, err := s.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	list := make([]*domain.PasswordComplexitySettings, len(settings))
	for i := range settings {
		list[i] = &domain.PasswordComplexitySettings{Settings: *settings[i]}
		if err := json.Unmarshal(list[i].Settings.Settings, &list[i].PasswordComplexitySettingsAttributes); err != nil {
			return nil, err
		}
		list[i].Settings.Settings = nil
	}
	return list, nil
}

func (s passwordComplexitySettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.PasswordComplexitySettings, error) {
	settings, err := s.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	item := &domain.PasswordComplexitySettings{Settings: *settings}
	if err := json.Unmarshal(item.Settings.Settings, &item.PasswordComplexitySettingsAttributes); err != nil {
		return nil, err
	}
	item.Settings.Settings = nil
	return item, nil
}

// -------------------------------------------------------------
// password expiry
// -------------------------------------------------------------

type passwordExpirySettings struct {
	settings
}

func (s passwordExpirySettings) SetSettingFields(value domain.PasswordExpirySettingsAttributes) database.Change {
	changes := make([]db_json.JsonUpdate, 0)
	if value.ExpireWarnDays != nil {
		changes = append(changes, s.SetExpireWarnDays(*value.ExpireWarnDays))
	}
	if value.MaxAgeDays != nil {
		changes = append(changes, s.SetMaxAgeDays(*value.MaxAgeDays))
	}
	return db_json.NewJsonChanges(s.SettingsColumn(), changes...)
}

func (passwordExpirySettings) SetExpireWarnDays(value uint64) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"expireWarnDays"}, value)
}

func (passwordExpirySettings) SetMaxAgeDays(value uint64) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"maxAgeDays"}, value)
}

var _ domain.PasswordExpirySettingsRepository = (*passwordExpirySettings)(nil)

func PasswordExpirySettingsRepository() domain.PasswordExpirySettingsRepository {
	return &passwordExpirySettings{
		settings{},
	}
}

func (s passwordExpirySettings) Set(ctx context.Context, client database.QueryExecutor, settings *domain.PasswordExpirySettings) (err error) {
	settings.Type = domain.SettingTypePasswordExpiry
	settings.State = domain.SettingStateActive

	settingsJSON, err := json.Marshal(settings.PasswordExpirySettingsAttributes)
	if err != nil {
		return err
	}
	settings.Settings.Settings = settingsJSON
	if err := s.set(ctx, client, &settings.Settings, s.SetSettingFields(settings.PasswordExpirySettingsAttributes)); err != nil {
		return err
	}
	settings.Settings.Settings = nil
	return nil
}

func (s passwordExpirySettings) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.PasswordExpirySettings, error) {
	settings, err := s.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	list := make([]*domain.PasswordExpirySettings, len(settings))
	for i := range settings {
		list[i] = &domain.PasswordExpirySettings{Settings: *settings[i]}
		if err := json.Unmarshal(list[i].Settings.Settings, &list[i].PasswordExpirySettingsAttributes); err != nil {
			return nil, err
		}
		list[i].Settings.Settings = nil
	}
	return list, nil
}

func (s passwordExpirySettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.PasswordExpirySettings, error) {
	settings, err := s.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	item := &domain.PasswordExpirySettings{Settings: *settings}
	if err = json.Unmarshal(item.Settings.Settings, &item.PasswordExpirySettingsAttributes); err != nil {
		return nil, err
	}
	item.Settings.Settings = nil
	return item, nil
}

// -------------------------------------------------------------
// lockout
// -------------------------------------------------------------
type lockoutSettings struct {
	settings
}

func (s lockoutSettings) SetSettingFields(value domain.LockoutSettingsAttributes) database.Change {
	changes := make([]db_json.JsonUpdate, 0)
	if value.MaxPasswordAttempts != nil {
		changes = append(changes, s.SetMaxPasswordAttempts(*value.MaxPasswordAttempts))
	}
	if value.MaxOTPAttempts != nil {
		changes = append(changes, s.SetMaxOTPAttempts(*value.MaxOTPAttempts))
	}
	if value.ShowLockOutFailures != nil {
		changes = append(changes, s.SetShowLockOutFailures(*value.ShowLockOutFailures))
	}
	return db_json.NewJsonChanges(s.SettingsColumn(), changes...)
}

func (lockoutSettings) SetMaxPasswordAttempts(value uint64) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"maxPasswordAttempts"}, value)
}

func (lockoutSettings) SetMaxOTPAttempts(value uint64) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"maxOtpAttempts"}, value)
}

func (lockoutSettings) SetShowLockOutFailures(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"showLockOutFailures"}, value)
}

func LockoutSettingsRepository() domain.LockoutSettingsRepository {
	return &lockoutSettings{
		settings{},
	}
}

var _ domain.LockoutSettingsRepository = (*lockoutSettings)(nil)

func (s lockoutSettings) Set(ctx context.Context, client database.QueryExecutor, settings *domain.LockoutSettings) error {
	settings.Type = domain.SettingTypeLockout
	settings.State = domain.SettingStateActive

	settingsJSON, err := json.Marshal(settings.LockoutSettingsAttributes)
	if err != nil {
		return err
	}
	settings.Settings.Settings = settingsJSON
	if err := s.set(ctx, client, &settings.Settings, s.SetSettingFields(settings.LockoutSettingsAttributes)); err != nil {
		return err
	}
	settings.Settings.Settings = nil
	return nil
}

func (s lockoutSettings) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.LockoutSettings, error) {
	settings, err := s.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	list := make([]*domain.LockoutSettings, len(settings))
	for i := range settings {
		list[i] = &domain.LockoutSettings{Settings: *settings[i]}
		if err := json.Unmarshal(list[i].Settings.Settings, &list[i].LockoutSettingsAttributes); err != nil {
			return nil, err
		}
		list[i].Settings.Settings = nil
	}
	return list, nil
}

func (s lockoutSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.LockoutSettings, error) {
	settings, err := s.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	item := &domain.LockoutSettings{Settings: *settings}
	if err = json.Unmarshal(item.Settings.Settings, &item.LockoutSettingsAttributes); err != nil {
		return nil, err
	}
	item.Settings.Settings = nil
	return item, nil
}

// -------------------------------------------------------------
// security
// -------------------------------------------------------------
type securitySettings struct {
	settings
}

func (s securitySettings) SetSettingFields(value domain.SecuritySettingsAttributes) database.Change {
	changes := make([]db_json.JsonUpdate, 0)
	if value.EnableIframeEmbedding != nil {
		changes = append(changes, s.SetEnableIframeEmbedding(*value.EnableIframeEmbedding))
	}
	if value.AllowedOrigins != nil {
		changes = append(changes, s.SetAllowedOrigins(value.AllowedOrigins))
	}
	if value.EnableImpersonation != nil {
		changes = append(changes, s.SetEnableImpersonation(*value.EnableImpersonation))
	}
	return db_json.NewJsonChanges(s.SettingsColumn(), changes...)
}

func (securitySettings) SetEnableIframeEmbedding(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"enableIframeEmbedding"}, value)
}

func (securitySettings) SetAllowedOrigins(value []string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"allowedOrigins"}, value)
}

func (securitySettings) SetEnableImpersonation(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"enableImpersonation"}, value)
}

func SecuritySettingsRepository() domain.SecuritySettingsRepository {
	return &securitySettings{
		settings{},
	}
}

var _ domain.SecuritySettingsRepository = (*securitySettings)(nil)

func (s securitySettings) Set(ctx context.Context, client database.QueryExecutor, settings *domain.SecuritySettings) (err error) {
	settings.Type = domain.SettingTypeSecurity
	settings.State = domain.SettingStateActive

	settingsJSON, err := json.Marshal(settings.SecuritySettingsAttributes)
	if err != nil {
		return err
	}
	settings.Settings.Settings = settingsJSON
	if err := s.set(ctx, client, &settings.Settings, s.SetSettingFields(settings.SecuritySettingsAttributes)); err != nil {
		return err
	}
	settings.Settings.Settings = nil
	return nil
}

func (s securitySettings) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.SecuritySettings, error) {
	settings, err := s.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	list := make([]*domain.SecuritySettings, len(settings))
	for i := range settings {
		list[i] = &domain.SecuritySettings{Settings: *settings[i]}
		if err := json.Unmarshal(list[i].Settings.Settings, &list[i].SecuritySettingsAttributes); err != nil {
			return nil, err
		}
		list[i].Settings.Settings = nil
	}
	return list, nil
}

func (s securitySettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.SecuritySettings, error) {
	settings, err := s.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	item := &domain.SecuritySettings{Settings: *settings}
	if err = json.Unmarshal(item.Settings.Settings, &item.SecuritySettingsAttributes); err != nil {
		return nil, err
	}
	item.Settings.Settings = nil
	return item, nil
}

// -------------------------------------------------------------
// domain
// -------------------------------------------------------------
type domainSettings struct {
	settings
}

func (s domainSettings) SetSettingFields(value domain.DomainSettingsAttributes) database.Change {
	changes := make([]db_json.JsonUpdate, 0)
	if value.LoginNameIncludesDomain != nil {
		changes = append(changes, s.SetLoginNameIncludesDomain(*value.LoginNameIncludesDomain))
	}
	if value.RequireOrgDomainVerification != nil {
		changes = append(changes, s.SetRequireOrgDomainVerification(*value.RequireOrgDomainVerification))
	}
	if value.SMTPSenderAddressMatchesInstanceDomain != nil {
		changes = append(changes, s.SetSMTPSenderAddressMatchesInstanceDomain(*value.SMTPSenderAddressMatchesInstanceDomain))
	}
	return db_json.NewJsonChanges(s.SettingsColumn(), changes...)
}

func (domainSettings) SetLoginNameIncludesDomain(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"loginNameIncludesDomain"}, value)
}

func (domainSettings) SetRequireOrgDomainVerification(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"requireOrgDomainVerification"}, value)
}

func (domainSettings) SetSMTPSenderAddressMatchesInstanceDomain(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"smtpSenderAddressMatchesInstanceDomain"}, value)
}

func DomainSettingsRepository() domain.DomainSettingsRepository {
	return &domainSettings{
		settings{},
	}
}

var _ domain.DomainSettingsRepository = (*domainSettings)(nil)

func (s domainSettings) Set(ctx context.Context, client database.QueryExecutor, settings *domain.DomainSettings) (err error) {
	settings.Type = domain.SettingTypeDomain
	settings.State = domain.SettingStateActive

	settingsJSON, err := json.Marshal(settings.DomainSettingsAttributes)
	if err != nil {
		return err
	}
	settings.Settings.Settings = settingsJSON
	if err := s.set(ctx, client, &settings.Settings, s.SetSettingFields(settings.DomainSettingsAttributes)); err != nil {
		return err
	}
	settings.Settings.Settings = nil
	return nil
}

func (s domainSettings) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.DomainSettings, error) {
	settings, err := s.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	list := make([]*domain.DomainSettings, len(settings))
	for i := range settings {
		list[i] = &domain.DomainSettings{Settings: *settings[i]}
		if err := json.Unmarshal(list[i].Settings.Settings, &list[i].DomainSettingsAttributes); err != nil {
			return nil, err
		}
		list[i].Settings.Settings = nil
	}
	return list, nil
}

func (s domainSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.DomainSettings, error) {
	settings, err := s.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	item := &domain.DomainSettings{Settings: *settings}
	if err = json.Unmarshal(item.Settings.Settings, &item.DomainSettingsAttributes); err != nil {
		return nil, err
	}
	item.Settings.Settings = nil
	return item, nil
}

// -------------------------------------------------------------
// organization
// -------------------------------------------------------------
type organizationSettings struct {
	settings
}

func (s organizationSettings) SetSettingFields(value domain.OrganizationSettingsAttributes) database.Change {
	changes := make([]db_json.JsonUpdate, 0)
	if value.OrganizationScopedUsernames != nil {
		changes = append(changes, s.SetOrganizationScopedUsernames(*value.OrganizationScopedUsernames))
	}
	return db_json.NewJsonChanges(s.SettingsColumn(), changes...)
}

func (l organizationSettings) SetOrganizationScopedUsernames(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"organizationScopedUsernames"}, value)
}

func OrganizationSettingsRepository() domain.OrganizationSettingsRepository {
	return &organizationSettings{
		settings{},
	}
}

var _ domain.OrganizationSettingsRepository = (*organizationSettings)(nil)

func (s organizationSettings) Set(ctx context.Context, client database.QueryExecutor, settings *domain.OrganizationSettings) (err error) {
	settings.Type = domain.SettingTypeOrganization
	settings.State = domain.SettingStateActive

	settingsJSON, err := json.Marshal(settings.OrganizationSettingsAttributes)
	if err != nil {
		return err
	}
	settings.Settings.Settings = settingsJSON
	if err := s.set(ctx, client, &settings.Settings, s.SetSettingFields(settings.OrganizationSettingsAttributes)); err != nil {
		return err
	}
	settings.Settings.Settings = nil
	return nil
}

func (s organizationSettings) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.OrganizationSettings, error) {
	settings, err := s.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	list := make([]*domain.OrganizationSettings, len(settings))
	for i := range settings {
		list[i] = &domain.OrganizationSettings{Settings: *settings[i]}
		if err := json.Unmarshal(list[i].Settings.Settings, &list[i].OrganizationSettingsAttributes); err != nil {
			return nil, err
		}
		list[i].Settings.Settings = nil
	}
	return list, nil
}

func (s organizationSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.OrganizationSettings, error) {
	settings, err := s.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	item := &domain.OrganizationSettings{Settings: *settings}
	if err = json.Unmarshal(item.Settings.Settings, &item.OrganizationSettingsAttributes); err != nil {
		return nil, err
	}
	item.Settings.Settings = nil
	return item, nil
}

// -------------------------------------------------------------
// notification
// -------------------------------------------------------------
type notificationSettings struct {
	settings
}

func (s notificationSettings) SetSettingFields(value domain.NotificationSettingsAttributes) database.Change {
	changes := make([]db_json.JsonUpdate, 0)
	if value.PasswordChange != nil {
		changes = append(changes, s.SetPasswordChange(*value.PasswordChange))
	}
	return db_json.NewJsonChanges(s.SettingsColumn(), changes...)
}

func (s notificationSettings) SetPasswordChange(value bool) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"passwordChange"}, value)
}

func NotificationSettingsRepository() domain.NotificationSettingsRepository {
	return &notificationSettings{
		settings{},
	}
}

var _ domain.NotificationSettingsRepository = (*notificationSettings)(nil)

func (s notificationSettings) Set(ctx context.Context, client database.QueryExecutor, settings *domain.NotificationSettings) (err error) {
	settings.Type = domain.SettingTypeNotification
	settings.State = domain.SettingStateActive

	settingsJSON, err := json.Marshal(settings.NotificationSettingsAttributes)
	if err != nil {
		return err
	}
	settings.Settings.Settings = settingsJSON
	if err := s.set(ctx, client, &settings.Settings, s.SetSettingFields(settings.NotificationSettingsAttributes)); err != nil {
		return err
	}
	settings.Settings.Settings = nil
	return nil
}

func (s notificationSettings) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.NotificationSettings, error) {
	settings, err := s.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	list := make([]*domain.NotificationSettings, len(settings))
	for i := range settings {
		list[i] = &domain.NotificationSettings{Settings: *settings[i]}
		if err := json.Unmarshal(list[i].Settings.Settings, &list[i].NotificationSettingsAttributes); err != nil {
			return nil, err
		}
		list[i].Settings.Settings = nil
	}
	return list, nil
}

func (s notificationSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.NotificationSettings, error) {
	settings, err := s.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	item := &domain.NotificationSettings{Settings: *settings}
	if err = json.Unmarshal(item.Settings.Settings, &item.NotificationSettingsAttributes); err != nil {
		return nil, err
	}
	item.Settings.Settings = nil
	return item, nil
}

// -------------------------------------------------------------
// legal and support
// -------------------------------------------------------------
type legalAndSupportSettings struct {
	settings
}

func (s legalAndSupportSettings) SetSettingFields(value domain.LegalAndSupportSettingsAttributes) database.Change {
	changes := make([]db_json.JsonUpdate, 0)
	if value.TOSLink != nil {
		changes = append(changes, s.SetTOSLink(*value.TOSLink))
	}
	if value.PrivacyPolicyLink != nil {
		changes = append(changes, s.SetPrivacyPolicyLink(*value.PrivacyPolicyLink))
	}
	if value.HelpLink != nil {
		changes = append(changes, s.SetHelpLink(*value.HelpLink))
	}
	if value.SupportEmail != nil {
		changes = append(changes, s.SetSupportEmail(*value.SupportEmail))
	}
	if value.DocsLink != nil {
		changes = append(changes, s.SetDocsLink(*value.DocsLink))
	}
	if value.CustomLink != nil {
		changes = append(changes, s.SetCustomLink(*value.CustomLink))
	}
	if value.CustomLinkText != nil {
		changes = append(changes, s.SetCustomLinkText(*value.CustomLinkText))
	}
	return db_json.NewJsonChanges(s.SettingsColumn(), changes...)
}

func (s legalAndSupportSettings) SetTOSLink(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"tosLink"}, value)
}

func (s legalAndSupportSettings) SetPrivacyPolicyLink(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"privacyPolicyLink"}, value)
}

func (s legalAndSupportSettings) SetHelpLink(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"helpLink"}, value)
}

func (s legalAndSupportSettings) SetSupportEmail(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"supportEmail"}, value)
}

func (s legalAndSupportSettings) SetDocsLink(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"docsLink"}, value)
}

func (s legalAndSupportSettings) SetCustomLink(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"customLink"}, value)
}

func (s legalAndSupportSettings) SetCustomLinkText(value string) db_json.JsonUpdate {
	return db_json.NewFieldChange([]string{"customLinkText"}, value)
}

func LegalAndSupportSettingsRepository() domain.LegalAndSupportSettingsRepository {
	return &legalAndSupportSettings{
		settings{},
	}
}

var _ domain.LegalAndSupportSettingsRepository = (*legalAndSupportSettings)(nil)

func (s legalAndSupportSettings) Set(ctx context.Context, client database.QueryExecutor, settings *domain.LegalAndSupportSettings) (err error) {
	settings.Type = domain.SettingTypeLegalAndSupport
	settings.State = domain.SettingStateActive

	settingsJSON, err := json.Marshal(settings.LegalAndSupportSettingsAttributes)
	if err != nil {
		return err
	}
	settings.Settings.Settings = settingsJSON
	if err := s.set(ctx, client, &settings.Settings, s.SetSettingFields(settings.LegalAndSupportSettingsAttributes)); err != nil {
		return err
	}
	settings.Settings.Settings = nil
	return nil
}

func (s legalAndSupportSettings) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.LegalAndSupportSettings, error) {
	settings, err := s.list(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	list := make([]*domain.LegalAndSupportSettings, len(settings))
	for i := range settings {
		list[i] = &domain.LegalAndSupportSettings{Settings: *settings[i]}
		if err := json.Unmarshal(list[i].Settings.Settings, &list[i].LegalAndSupportSettingsAttributes); err != nil {
			return nil, err
		}
		list[i].Settings.Settings = nil
	}
	return list, nil
}

func (s legalAndSupportSettings) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.LegalAndSupportSettings, error) {
	settings, err := s.get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	item := &domain.LegalAndSupportSettings{Settings: *settings}
	if err = json.Unmarshal(item.Settings.Settings, &item.LegalAndSupportSettingsAttributes); err != nil {
		return nil, err
	}
	item.Settings.Settings = nil
	return item, nil
}

// -------------------------------------------------------------
// helpers
// -------------------------------------------------------------

const querySettingsStmt = `SELECT 
	settings.instance_id, 
	settings.organization_id, 
	settings.id, 
	settings.type, 
	settings.state,
	settings.settings,
	settings.created_at, 
	settings.updated_at 
	FROM `

func (s settings) prepareQuery(opts []database.QueryOption) (*database.StatementBuilder, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if err := checkRestrictingColumns(options.Condition, s.InstanceIDColumn()); err != nil {
		return nil, err
	}
	builder := database.NewStatementBuilder(querySettingsStmt + s.qualifiedTableName())
	options.Write(builder)

	return builder, nil
}
