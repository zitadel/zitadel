package query

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/text/language"
	"sigs.k8s.io/yaml"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type CustomTexts struct {
	SearchResponse
	CustomTexts []*CustomText
}

type CustomText struct {
	AggregateID  string
	Sequence     uint64
	CreationDate time.Time
	ChangeDate   time.Time

	Template string
	Language language.Tag
	Key      string
	Text     string
}

var (
	customTextTable = table{
		name:          projection.CustomTextTable,
		instanceIDCol: projection.CustomTextInstanceIDCol,
	}
	CustomTextColAggregateID = Column{
		name:  projection.CustomTextAggregateIDCol,
		table: customTextTable,
	}
	CustomTextColInstanceID = Column{
		name:  projection.CustomTextInstanceIDCol,
		table: customTextTable,
	}
	CustomTextColSequence = Column{
		name:  projection.CustomTextSequenceCol,
		table: customTextTable,
	}
	CustomTextColCreationDate = Column{
		name:  projection.CustomTextCreationDateCol,
		table: customTextTable,
	}
	CustomTextColChangeDate = Column{
		name:  projection.CustomTextChangeDateCol,
		table: customTextTable,
	}
	CustomTextColTemplate = Column{
		name:  projection.CustomTextTemplateCol,
		table: customTextTable,
	}
	CustomTextColLanguage = Column{
		name:  projection.CustomTextLanguageCol,
		table: customTextTable,
	}
	CustomTextColKey = Column{
		name:  projection.CustomTextKeyCol,
		table: customTextTable,
	}
	CustomTextColText = Column{
		name:  projection.CustomTextTextCol,
		table: customTextTable,
	}
	CustomTextOwnerRemoved = Column{
		name:  projection.CustomTextOwnerRemovedCol,
		table: customTextTable,
	}
)

func (q *Queries) CustomTextList(ctx context.Context, aggregateID, template, language string, withOwnerRemoved bool) (texts *CustomTexts, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareCustomTextsQuery(ctx, q.client)
	eq := sq.Eq{
		CustomTextColAggregateID.identifier(): aggregateID,
		CustomTextColTemplate.identifier():    template,
		CustomTextColLanguage.identifier():    language,
		CustomTextColInstanceID.identifier():  authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[CustomTextOwnerRemoved.identifier()] = false
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-M9gse", "Errors.Query.SQLStatement")
	}

	rows, err := q.client.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-2j00f", "Errors.Internal")
	}
	texts, err = scan(rows)
	if err != nil {
		return nil, err
	}
	texts.LatestSequence, err = q.latestSequence(ctx, projectsTable)
	return texts, err
}

func (q *Queries) CustomTextListByTemplate(ctx context.Context, aggregateID, template string, withOwnerRemoved bool) (texts *CustomTexts, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareCustomTextsQuery(ctx, q.client)
	eq := sq.Eq{
		CustomTextColAggregateID.identifier(): aggregateID,
		CustomTextColTemplate.identifier():    template,
		CustomTextColInstanceID.identifier():  authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[CustomTextOwnerRemoved.identifier()] = false
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-M49fs", "Errors.Query.SQLStatement")
	}

	rows, err := q.client.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-3n9ge", "Errors.Internal")
	}
	texts, err = scan(rows)
	if err != nil {
		return nil, err
	}
	texts.LatestSequence, err = q.latestSequence(ctx, projectsTable)
	return texts, err
}

func (q *Queries) GetDefaultLoginTexts(ctx context.Context, lang string) (_ *domain.CustomLoginText, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	contents, err := q.readLoginTranslationFile(ctx, lang)
	if err != nil {
		return nil, err
	}
	loginText := new(domain.CustomLoginText)
	if err := yaml.Unmarshal(contents, loginText); err != nil {
		return nil, errors.ThrowInternal(err, "TEXT-M0p4s", "Errors.TranslationFile.ReadError")
	}
	loginText.IsDefault = true
	loginText.AggregateID = authz.GetInstance(ctx).InstanceID()
	return loginText, nil
}

func (q *Queries) GetCustomLoginTexts(ctx context.Context, aggregateID, lang string) (_ *domain.CustomLoginText, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	texts, err := q.CustomTextList(ctx, aggregateID, domain.LoginCustomText, lang, false)
	if err != nil {
		return nil, err
	}
	return CustomTextsToLoginDomain(authz.GetInstance(ctx).InstanceID(), aggregateID, lang, texts), err
}

func (q *Queries) IAMLoginTexts(ctx context.Context, lang string) (_ *domain.CustomLoginText, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	contents, err := q.readLoginTranslationFile(ctx, lang)
	if err != nil {
		return nil, err
	}
	loginTextMap := make(map[string]interface{})
	if err := yaml.Unmarshal(contents, &loginTextMap); err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-m0Jf3", "Errors.TranslationFile.ReadError")
	}
	texts, err := q.CustomTextList(ctx, authz.GetInstance(ctx).InstanceID(), domain.LoginCustomText, lang, false)
	if err != nil {
		return nil, err
	}
	for _, text := range texts.CustomTexts {
		keys := strings.Split(text.Key, ".")
		screenTextMap, ok := loginTextMap[keys[0]].(map[string]interface{})
		if !ok {
			continue
		}
		screenTextMap[keys[1]] = text.Text
	}
	jsonbody, err := json.Marshal(loginTextMap)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-0nJ3f", "Errors.TranslationFile.MergeError")
	}
	loginText := new(domain.CustomLoginText)
	if err := json.Unmarshal(jsonbody, &loginText); err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-m93Jf", "Errors.TranslationFile.MergeError")
	}
	loginText.AggregateID = authz.GetInstance(ctx).InstanceID()
	loginText.IsDefault = true
	return loginText, nil
}

func (q *Queries) readLoginTranslationFile(ctx context.Context, lang string) ([]byte, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	contents, ok := q.LoginTranslationFileContents[lang]
	var err error
	if !ok {
		contents, err = q.readTranslationFile(q.LoginDir, fmt.Sprintf("/i18n/%s.yaml", lang))
		if errors.IsNotFound(err) {
			contents, err = q.readTranslationFile(q.LoginDir, fmt.Sprintf("/i18n/%s.yaml", authz.GetInstance(ctx).DefaultLanguage().String()))
		}
		if err != nil {
			return nil, err
		}
		q.LoginTranslationFileContents[lang] = contents
	}
	return contents, nil
}

func prepareCustomTextsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*CustomTexts, error)) {
	return sq.Select(
			CustomTextColAggregateID.identifier(),
			CustomTextColSequence.identifier(),
			CustomTextColCreationDate.identifier(),
			CustomTextColChangeDate.identifier(),
			CustomTextColLanguage.identifier(),
			CustomTextColTemplate.identifier(),
			CustomTextColKey.identifier(),
			CustomTextColText.identifier(),
			countColumn.identifier()).
			From(customTextTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*CustomTexts, error) {
			customTexts := make([]*CustomText, 0)
			var count uint64
			for rows.Next() {
				customText := new(CustomText)
				lang := ""
				err := rows.Scan(
					&customText.AggregateID,
					&customText.Sequence,
					&customText.CreationDate,
					&customText.ChangeDate,
					&lang,
					&customText.Template,
					&customText.Key,
					&customText.Text,
					&count,
				)
				if err != nil {
					return nil, err
				}
				customText.Language = language.Make(lang)
				customTexts = append(customTexts, customText)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-3n9fs", "Errors.Query.CloseRows")
			}

			return &CustomTexts{
				CustomTexts: customTexts,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func CustomTextsToDomain(texts *CustomTexts) []*domain.CustomText {
	result := make([]*domain.CustomText, len(texts.CustomTexts))
	for i, text := range texts.CustomTexts {
		result[i] = CustomTextToDomain(text)
	}
	return result
}

func CustomTextToDomain(text *CustomText) *domain.CustomText {
	return &domain.CustomText{
		ObjectRoot: models.ObjectRoot{
			AggregateID:  text.AggregateID,
			Sequence:     text.Sequence,
			CreationDate: text.CreationDate,
			ChangeDate:   text.ChangeDate,
		},
		Template: text.Template,
		Language: text.Language,
		Key:      text.Key,
		Text:     text.Text,
	}
}

func CustomTextsToLoginDomain(instanceID, aggregateID, lang string, texts *CustomTexts) *domain.CustomLoginText {
	langTag := language.Make(lang)
	result := &domain.CustomLoginText{
		ObjectRoot: models.ObjectRoot{
			AggregateID: aggregateID,
		},
		Language: langTag,
	}
	if len(texts.CustomTexts) == 0 {
		result.AggregateID = instanceID
		result.IsDefault = true
	}
	for _, text := range texts.CustomTexts {
		if text.CreationDate.Before(result.CreationDate) {
			result.CreationDate = text.CreationDate
		}
		if text.ChangeDate.After(result.ChangeDate) {
			result.ChangeDate = text.ChangeDate
		}
		if strings.HasPrefix(text.Key, domain.LoginKeySelectAccount) {
			selectAccountKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyLogin) {
			loginKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyPassword) {
			passwordKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyUsernameChange) {
			usernameChangeKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyUsernameChangeDone) {
			usernameChangeDoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitPassword) {
			initPasswordKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitPasswordDone) {
			initPasswordDoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyEmailVerification) {
			emailVerificationKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyEmailVerificationDone) {
			emailVerificationDoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitializeUser) {
			initializeUserKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitUserDone) {
			initializeUserDoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitMFAPrompt) {
			initMFAPromptKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitMFAOTP) {
			initMFAOTPKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitMFAU2F) {
			initMFAU2FKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyInitMFADone) {
			initMFADoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyMFAProviders) {
			mfaProvidersKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyVerifyMFAOTP) {
			verifyMFAOTPKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyVerifyMFAU2F) {
			verifyMFAU2FKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyPasswordless) {
			passwordlessKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyPasswordlessPrompt) {
			passwordlessPromptKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyPasswordlessRegistration) {
			passwordlessRegistrationKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyPasswordlessRegistrationDone) {
			passwordlessRegistrationDoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyPasswordChange) {
			passwordChangeKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyPasswordChangeDone) {
			passwordChangeDoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyPasswordResetDone) {
			passwordResetDoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyRegistrationOption) {
			registrationOptionKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyRegistrationUser) {
			registrationUserKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyExternalRegistrationUserOverview) {
			externalRegistrationUserKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyRegistrationOrg) {
			registrationOrgKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyLinkingUserDone) {
			linkingUserKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyExternalNotFound) {
			externalUserNotFoundKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeySuccessLogin) {
			successLoginKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyLogoutDone) {
			logoutDoneKeyToDomain(text, result)
		}
		if strings.HasPrefix(text.Key, domain.LoginKeyFooter) {
			footerKeyToDomain(text, result)
		}
	}
	return result
}

func selectAccountKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeySelectAccountTitle {
		result.SelectAccount.Title = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountDescription {
		result.SelectAccount.Description = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountTitleLinkingProcess {
		result.SelectAccount.TitleLinking = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountDescriptionLinkingProcess {
		result.SelectAccount.DescriptionLinking = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountOtherUser {
		result.SelectAccount.OtherUser = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountSessionStateActive {
		result.SelectAccount.SessionState0 = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountSessionStateInactive {
		result.SelectAccount.SessionState1 = text.Text
	}
	if text.Key == domain.LoginKeySelectAccountUserMustBeMemberOfOrg {
		result.SelectAccount.MustBeMemberOfOrg = text.Text
	}
}

func loginKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyLoginTitle {
		result.Login.Title = text.Text
	}
	if text.Key == domain.LoginKeyLoginDescription {
		result.Login.Description = text.Text
	}
	if text.Key == domain.LoginKeyLoginTitleLinkingProcess {
		result.Login.TitleLinking = text.Text
	}
	if text.Key == domain.LoginKeyLoginDescriptionLinkingProcess {
		result.Login.DescriptionLinking = text.Text
	}
	if text.Key == domain.LoginKeyLoginNameLabel {
		result.Login.LoginNameLabel = text.Text
	}
	if text.Key == domain.LoginKeyLoginUsernamePlaceHolder {
		result.Login.UsernamePlaceholder = text.Text
	}
	if text.Key == domain.LoginKeyLoginLoginnamePlaceHolder {
		result.Login.LoginnamePlaceholder = text.Text
	}
	if text.Key == domain.LoginKeyLoginExternalUserDescription {
		result.Login.ExternalUserDescription = text.Text
	}
	if text.Key == domain.LoginKeyLoginUserMustBeMemberOfOrg {
		result.Login.MustBeMemberOfOrg = text.Text
	}
	if text.Key == domain.LoginKeyLoginRegisterButtonText {
		result.Login.RegisterButtonText = text.Text
	}
	if text.Key == domain.LoginKeyLoginNextButtonText {
		result.Login.NextButtonText = text.Text
	}
}

func passwordKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordTitle {
		result.Password.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordDescription {
		result.Password.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordLabel {
		result.Password.PasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyPasswordResetLinkText {
		result.Password.ResetLinkText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordBackButtonText {
		result.Password.BackButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordNextButtonText {
		result.Password.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordMinLength {
		result.Password.MinLength = text.Text
	}
	if text.Key == domain.LoginKeyPasswordHasUppercase {
		result.Password.HasUppercase = text.Text
	}
	if text.Key == domain.LoginKeyPasswordHasLowercase {
		result.Password.HasLowercase = text.Text
	}
	if text.Key == domain.LoginKeyPasswordHasNumber {
		result.Password.HasNumber = text.Text
	}
	if text.Key == domain.LoginKeyPasswordHasSymbol {
		result.Password.HasSymbol = text.Text
	}
	if text.Key == domain.LoginKeyPasswordConfirmation {
		result.Password.Confirmation = text.Text
	}
}

func usernameChangeKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyUsernameChangeTitle {
		result.UsernameChange.Title = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeDescription {
		result.UsernameChange.Description = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeUsernameLabel {
		result.UsernameChange.UsernameLabel = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeCancelButtonText {
		result.UsernameChange.CancelButtonText = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeNextButtonText {
		result.UsernameChange.NextButtonText = text.Text
	}
}

func usernameChangeDoneKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyUsernameChangeDoneTitle {
		result.UsernameChangeDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeDoneDescription {
		result.UsernameChangeDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyUsernameChangeDoneNextButtonText {
		result.UsernameChangeDone.NextButtonText = text.Text
	}
}

func initPasswordKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitPasswordTitle {
		result.InitPassword.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordDescription {
		result.InitPassword.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordCodeLabel {
		result.InitPassword.CodeLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordNewPasswordLabel {
		result.InitPassword.NewPasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordNewPasswordConfirmLabel {
		result.InitPassword.NewPasswordConfirmLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordNextButtonText {
		result.InitPassword.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordResendButtonText {
		result.InitPassword.ResendButtonText = text.Text
	}
}

func initPasswordDoneKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitPasswordDoneTitle {
		result.InitPasswordDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordDoneDescription {
		result.InitPasswordDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordDoneNextButtonText {
		result.InitPasswordDone.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitPasswordDoneCancelButtonText {
		result.InitPasswordDone.CancelButtonText = text.Text
	}
}

func emailVerificationKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyEmailVerificationTitle {
		result.EmailVerification.Title = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationDescription {
		result.EmailVerification.Description = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationCodeLabel {
		result.EmailVerification.CodeLabel = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationNextButtonText {
		result.EmailVerification.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationResendButtonText {
		result.EmailVerification.ResendButtonText = text.Text
	}
}

func emailVerificationDoneKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyEmailVerificationDoneTitle {
		result.EmailVerificationDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationDoneDescription {
		result.EmailVerificationDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationDoneNextButtonText {
		result.EmailVerificationDone.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationDoneCancelButtonText {
		result.EmailVerificationDone.CancelButtonText = text.Text
	}
	if text.Key == domain.LoginKeyEmailVerificationDoneLoginButtonText {
		result.EmailVerificationDone.LoginButtonText = text.Text
	}
}

func initializeUserKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitializeUserTitle {
		result.InitUser.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserDescription {
		result.InitUser.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserCodeLabel {
		result.InitUser.CodeLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserNewPasswordLabel {
		result.InitUser.NewPasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserNewPasswordConfirmLabel {
		result.InitUser.NewPasswordConfirmLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserResendButtonText {
		result.InitUser.ResendButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitializeUserNextButtonText {
		result.InitUser.NextButtonText = text.Text
	}
}

func initializeUserDoneKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitUserDoneTitle {
		result.InitUserDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitUserDoneDescription {
		result.InitUserDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitUserDoneCancelButtonText {
		result.InitUserDone.CancelButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitUserDoneNextButtonText {
		result.InitUserDone.NextButtonText = text.Text
	}
}

func initMFAPromptKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitMFAPromptTitle {
		result.InitMFAPrompt.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAPromptDescription {
		result.InitMFAPrompt.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAPromptOTPOption {
		result.InitMFAPrompt.Provider0 = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAPromptU2FOption {
		result.InitMFAPrompt.Provider1 = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAPromptSkipButtonText {
		result.InitMFAPrompt.SkipButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAPromptNextButtonText {
		result.InitMFAPrompt.NextButtonText = text.Text
	}
}

func initMFAOTPKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitMFAOTPTitle {
		result.InitMFAOTP.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPDescription {
		result.InitMFAOTP.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPDescriptionOTP {
		result.InitMFAOTP.OTPDescription = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPCodeLabel {
		result.InitMFAOTP.CodeLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPSecretLabel {
		result.InitMFAOTP.SecretLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPNextButtonText {
		result.InitMFAOTP.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAOTPCancelButtonText {
		result.InitMFAOTP.CancelButtonText = text.Text
	}
}

func initMFAU2FKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitMFAU2FTitle {
		result.InitMFAU2F.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAU2FDescription {
		result.InitMFAU2F.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAU2FTokenNameLabel {
		result.InitMFAU2F.TokenNameLabel = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAU2FRegisterTokenButtonText {
		result.InitMFAU2F.RegisterTokenButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAU2FNotSupported {
		result.InitMFAU2F.NotSupported = text.Text
	}
	if text.Key == domain.LoginKeyInitMFAU2FErrorRetry {
		result.InitMFAU2F.ErrorRetry = text.Text
	}
}

func initMFADoneKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyInitMFADoneTitle {
		result.InitMFADone.Title = text.Text
	}
	if text.Key == domain.LoginKeyInitMFADoneDescription {
		result.InitMFADone.Description = text.Text
	}
	if text.Key == domain.LoginKeyInitMFADoneCancelButtonText {
		result.InitMFADone.CancelButtonText = text.Text
	}
	if text.Key == domain.LoginKeyInitMFADoneNextButtonText {
		result.InitMFADone.NextButtonText = text.Text
	}
}

func mfaProvidersKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyMFAProvidersChooseOther {
		result.MFAProvider.ChooseOther = text.Text
	}
	if text.Key == domain.LoginKeyMFAProvidersOTP {
		result.MFAProvider.Provider0 = text.Text
	}
	if text.Key == domain.LoginKeyMFAProvidersU2F {
		result.MFAProvider.Provider1 = text.Text
	}
}

func verifyMFAOTPKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyVerifyMFAOTPTitle {
		result.VerifyMFAOTP.Title = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAOTPDescription {
		result.VerifyMFAOTP.Description = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAOTPCodeLabel {
		result.VerifyMFAOTP.CodeLabel = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAOTPNextButtonText {
		result.VerifyMFAOTP.NextButtonText = text.Text
	}
}

func verifyMFAU2FKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyVerifyMFAU2FTitle {
		result.VerifyMFAU2F.Title = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAU2FDescription {
		result.VerifyMFAU2F.Description = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAU2FValidateTokenText {
		result.VerifyMFAU2F.ValidateTokenButtonText = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAU2FNotSupported {
		result.VerifyMFAU2F.NotSupported = text.Text
	}
	if text.Key == domain.LoginKeyVerifyMFAU2FErrorRetry {
		result.VerifyMFAU2F.ErrorRetry = text.Text
	}
}

func passwordlessKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordlessTitle {
		result.Passwordless.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessDescription {
		result.Passwordless.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessLoginWithPwButtonText {
		result.Passwordless.LoginWithPwButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessValidateTokenButtonText {
		result.Passwordless.ValidateTokenButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessNotSupported {
		result.Passwordless.NotSupported = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessErrorRetry {
		result.Passwordless.ErrorRetry = text.Text
	}
}

func passwordlessPromptKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordlessPromptTitle {
		result.PasswordlessPrompt.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessPromptDescription {
		result.PasswordlessPrompt.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessPromptDescriptionInit {
		result.PasswordlessPrompt.DescriptionInit = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessPromptPasswordlessButtonText {
		result.PasswordlessPrompt.PasswordlessButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessPromptNextButtonText {
		result.PasswordlessPrompt.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessPromptSkipButtonText {
		result.PasswordlessPrompt.SkipButtonText = text.Text
	}
}

func passwordlessRegistrationKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordlessRegistrationTitle {
		result.PasswordlessRegistration.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessRegistrationDescription {
		result.PasswordlessRegistration.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessRegistrationRegisterTokenButtonText {
		result.PasswordlessRegistration.RegisterTokenButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessRegistrationTokenNameLabel {
		result.PasswordlessRegistration.TokenNameLabel = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessRegistrationNotSupported {
		result.PasswordlessRegistration.NotSupported = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessRegistrationErrorRetry {
		result.PasswordlessRegistration.ErrorRetry = text.Text
	}
}

func passwordlessRegistrationDoneKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordlessRegistrationDoneTitle {
		result.PasswordlessRegistrationDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessRegistrationDoneDescription {
		result.PasswordlessRegistrationDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessRegistrationDoneDescriptionClose {
		result.PasswordlessRegistrationDone.DescriptionClose = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessRegistrationDoneNextButtonText {
		result.PasswordlessRegistrationDone.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordlessRegistrationDoneCancelButtonText {
		result.PasswordlessRegistrationDone.CancelButtonText = text.Text
	}
}

func passwordChangeKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordChangeTitle {
		result.PasswordChange.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeDescription {
		result.PasswordChange.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeOldPasswordLabel {
		result.PasswordChange.OldPasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeNewPasswordLabel {
		result.PasswordChange.NewPasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeNewPasswordConfirmLabel {
		result.PasswordChange.NewPasswordConfirmLabel = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeCancelButtonText {
		result.PasswordChange.CancelButtonText = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeNextButtonText {
		result.PasswordChange.NextButtonText = text.Text
	}
}

func passwordChangeDoneKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordChangeDoneTitle {
		result.PasswordChangeDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeDoneDescription {
		result.PasswordChangeDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordChangeDoneNextButtonText {
		result.PasswordChangeDone.NextButtonText = text.Text
	}
}

func passwordResetDoneKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyPasswordResetDoneTitle {
		result.PasswordResetDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyPasswordResetDoneDescription {
		result.PasswordResetDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyPasswordResetDoneNextButtonText {
		result.PasswordResetDone.NextButtonText = text.Text
	}
}

func registrationOptionKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyRegistrationOptionTitle {
		result.RegisterOption.Title = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationOptionDescription {
		result.RegisterOption.Description = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationOptionExternalLoginDescription {
		result.RegisterOption.ExternalLoginDescription = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationOptionUserNameButtonText {
		result.RegisterOption.RegisterUsernamePasswordButtonText = text.Text
	}
}

func externalRegistrationUserKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyExternalRegistrationUserOverviewBackButtonText {
		result.ExternalRegistrationUserOverview.BackButtonText = text.Text
	}
	if text.Key == domain.LoginKeyExternalRegistrationUserOverviewDescription {
		result.ExternalRegistrationUserOverview.Description = text.Text
	}
	if text.Key == domain.LoginKeyExternalRegistrationUserOverviewEmailLabel {
		result.ExternalRegistrationUserOverview.EmailLabel = text.Text
	}
	if text.Key == domain.LoginKeyExternalRegistrationUserOverviewFirstnameLabel {
		result.ExternalRegistrationUserOverview.FirstnameLabel = text.Text
	}
	if text.Key == domain.LoginKeyExternalRegistrationUserOverviewLanguageLabel {
		result.ExternalRegistrationUserOverview.LanguageLabel = text.Text
	}
	if text.Key == domain.LoginKeyExternalRegistrationUserOverviewLastnameLabel {
		result.ExternalRegistrationUserOverview.LastnameLabel = text.Text
	}
	if text.Key == domain.LoginKeyExternalRegistrationUserOverviewNextButtonText {
		result.ExternalRegistrationUserOverview.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyExternalRegistrationUserOverviewNicknameLabel {
		result.ExternalRegistrationUserOverview.NicknameLabel = text.Text
	}
	if text.Key == domain.LoginKeyExternalRegistrationUserOverviewPhoneLabel {
		result.ExternalRegistrationUserOverview.PhoneLabel = text.Text
	}
	if text.Key == domain.LoginKeyExternalRegistrationUserOverviewPrivacyLinkText {
		result.ExternalRegistrationUserOverview.PrivacyLinkText = text.Text
	}
	if text.Key == domain.LoginKeyExternalRegistrationUserOverviewTitle {
		result.ExternalRegistrationUserOverview.Title = text.Text
	}
	if text.Key == domain.LoginKeyExternalRegistrationUserOverviewTOSAndPrivacyLabel {
		result.ExternalRegistrationUserOverview.TOSAndPrivacyLabel = text.Text
	}
	if text.Key == domain.LoginKeyExternalRegistrationUserOverviewTOSConfirm {
		result.ExternalRegistrationUserOverview.TOSConfirm = text.Text
	}
	if text.Key == domain.LoginKeyExternalRegistrationUserOverviewPrivacyConfirm {
		result.ExternalRegistrationUserOverview.PrivacyConfirm = text.Text
	}
	if text.Key == domain.LoginKeyExternalRegistrationUserOverviewTOSLinkText {
		result.ExternalRegistrationUserOverview.TOSLinkText = text.Text
	}
	if text.Key == domain.LoginKeyExternalRegistrationUserOverviewUsernameLabel {
		result.ExternalRegistrationUserOverview.UsernameLabel = text.Text
	}
}

func registrationUserKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyRegistrationUserTitle {
		result.RegistrationUser.Title = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserDescription {
		result.RegistrationUser.Description = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserDescriptionOrgRegister {
		result.RegistrationUser.DescriptionOrgRegister = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserFirstnameLabel {
		result.RegistrationUser.FirstnameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserLastnameLabel {
		result.RegistrationUser.LastnameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserEmailLabel {
		result.RegistrationUser.EmailLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserUsernameLabel {
		result.RegistrationUser.UsernameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserLanguageLabel {
		result.RegistrationUser.LanguageLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserGenderLabel {
		result.RegistrationUser.GenderLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserPasswordLabel {
		result.RegistrationUser.PasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserPasswordConfirmLabel {
		result.RegistrationUser.PasswordConfirmLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserTOSAndPrivacyLabel {
		result.RegistrationUser.TOSAndPrivacyLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserTOSConfirm {
		result.RegistrationUser.TOSConfirm = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserTOSLinkText {
		result.RegistrationUser.TOSLinkText = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserPrivacyConfirm {
		result.RegistrationUser.PrivacyConfirm = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserPrivacyLinkText {
		result.RegistrationUser.PrivacyLinkText = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserNextButtonText {
		result.RegistrationUser.NextButtonText = text.Text
	}
	if text.Key == domain.LoginKeyRegistrationUserBackButtonText {
		result.RegistrationUser.BackButtonText = text.Text
	}
}

func registrationOrgKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyRegisterOrgTitle {
		result.RegistrationOrg.Title = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgDescription {
		result.RegistrationOrg.Description = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgOrgNameLabel {
		result.RegistrationOrg.OrgNameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgFirstnameLabel {
		result.RegistrationOrg.FirstnameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgLastnameLabel {
		result.RegistrationOrg.LastnameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgUsernameLabel {
		result.RegistrationOrg.UsernameLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgEmailLabel {
		result.RegistrationOrg.EmailLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgPasswordLabel {
		result.RegistrationOrg.PasswordLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgPasswordConfirmLabel {
		result.RegistrationOrg.PasswordConfirmLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgTOSAndPrivacyLabel {
		result.RegistrationOrg.TOSAndPrivacyLabel = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgTOSConfirm {
		result.RegistrationOrg.TOSConfirm = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgTOSLinkText {
		result.RegistrationOrg.TOSLinkText = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgPrivacyConfirm {
		result.RegistrationOrg.PrivacyConfirm = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgPrivacyLinkText {
		result.RegistrationOrg.PrivacyLinkText = text.Text
	}
	if text.Key == domain.LoginKeyRegisterOrgSaveButtonText {
		result.RegistrationOrg.SaveButtonText = text.Text
	}
}

func linkingUserKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyLinkingUserDoneTitle {
		result.LinkingUsersDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyLinkingUserDoneDescription {
		result.LinkingUsersDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyLinkingUserDoneCancelButtonText {
		result.LinkingUsersDone.CancelButtonText = text.Text
	}
	if text.Key == domain.LoginKeyLinkingUserDoneNextButtonText {
		result.LinkingUsersDone.NextButtonText = text.Text
	}
}

func externalUserNotFoundKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyExternalNotFoundTitle {
		result.ExternalNotFound.Title = text.Text
	}
	if text.Key == domain.LoginKeyExternalNotFoundDescription {
		result.ExternalNotFound.Description = text.Text
	}
	if text.Key == domain.LoginKeyExternalNotFoundLinkButtonText {
		result.ExternalNotFound.LinkButtonText = text.Text
	}
	if text.Key == domain.LoginKeyExternalNotFoundAutoRegisterButtonText {
		result.ExternalNotFound.AutoRegisterButtonText = text.Text
	}
	if text.Key == domain.LoginKeyExternalNotFoundTOSAndPrivacyLabel {
		result.ExternalNotFound.TOSAndPrivacyLabel = text.Text
	}
	if text.Key == domain.LoginKeyExternalNotFoundTOSConfirm {
		result.ExternalNotFound.TOSConfirm = text.Text
	}
	if text.Key == domain.LoginKeyExternalNotFoundTOSLinkText {
		result.ExternalNotFound.TOSLinkText = text.Text
	}
	if text.Key == domain.LoginKeyExternalNotFoundPrivacyConfirm {
		result.ExternalNotFound.PrivacyConfirm = text.Text
	}
	if text.Key == domain.LoginKeyExternalNotFoundPrivacyLinkText {
		result.ExternalNotFound.PrivacyLinkText = text.Text
	}
}

func successLoginKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeySuccessLoginTitle {
		result.LoginSuccess.Title = text.Text
	}
	if text.Key == domain.LoginKeySuccessLoginAutoRedirectDescription {
		result.LoginSuccess.AutoRedirectDescription = text.Text
	}
	if text.Key == domain.LoginKeySuccessLoginRedirectedDescription {
		result.LoginSuccess.RedirectedDescription = text.Text
	}
	if text.Key == domain.LoginKeySuccessLoginNextButtonText {
		result.LoginSuccess.NextButtonText = text.Text
	}
}

func logoutDoneKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyLogoutDoneTitle {
		result.LogoutDone.Title = text.Text
	}
	if text.Key == domain.LoginKeyLogoutDoneDescription {
		result.LogoutDone.Description = text.Text
	}
	if text.Key == domain.LoginKeyLogoutDoneLoginButtonText {
		result.LogoutDone.LoginButtonText = text.Text
	}
}

func footerKeyToDomain(text *CustomText, result *domain.CustomLoginText) {
	if text.Key == domain.LoginKeyFooterTOS {
		result.Footer.TOS = text.Text
	}
	if text.Key == domain.LoginKeyFooterPrivacyPolicy {
		result.Footer.PrivacyPolicy = text.Text
	}
	if text.Key == domain.LoginKeyFooterHelp {
		result.Footer.Help = text.Text
	}
}
