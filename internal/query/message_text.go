package query

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/text/language"
	"sigs.k8s.io/yaml"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type MessageTexts struct {
	InitCode                 MessageText
	PasswordReset            MessageText
	VerifyEmail              MessageText
	VerifyPhone              MessageText
	VerifySMSOTP             MessageText
	VerifyEmailOTP           MessageText
	DomainClaimed            MessageText
	PasswordlessRegistration MessageText
	PasswordChange           MessageText
	InviteUser               MessageText
}

type MessageText struct {
	AggregateID  string
	Sequence     uint64
	CreationDate time.Time
	ChangeDate   time.Time
	State        domain.PolicyState

	IsDefault bool

	Type       string
	Language   language.Tag
	Title      string
	PreHeader  string
	Subject    string
	Greeting   string
	Text       string
	ButtonText string
	Footer     string
}

var (
	messageTextTable = table{
		name:          projection.MessageTextTable,
		instanceIDCol: projection.MessageTextInstanceIDCol,
	}
	MessageTextColAggregateID = Column{
		name:  projection.MessageTextAggregateIDCol,
		table: messageTextTable,
	}
	MessageTextColInstanceID = Column{
		name:  projection.MessageTextInstanceIDCol,
		table: messageTextTable,
	}
	MessageTextColSequence = Column{
		name:  projection.MessageTextSequenceCol,
		table: messageTextTable,
	}
	MessageTextColCreationDate = Column{
		name:  projection.MessageTextCreationDateCol,
		table: messageTextTable,
	}
	MessageTextColChangeDate = Column{
		name:  projection.MessageTextChangeDateCol,
		table: messageTextTable,
	}
	MessageTextColState = Column{
		name:  projection.MessageTextStateCol,
		table: messageTextTable,
	}
	MessageTextColType = Column{
		name:  projection.MessageTextTypeCol,
		table: messageTextTable,
	}
	MessageTextColLanguage = Column{
		name:  projection.MessageTextLanguageCol,
		table: messageTextTable,
	}
	MessageTextColTitle = Column{
		name:  projection.MessageTextTitleCol,
		table: messageTextTable,
	}
	MessageTextColPreHeader = Column{
		name:  projection.MessageTextPreHeaderCol,
		table: messageTextTable,
	}
	MessageTextColSubject = Column{
		name:  projection.MessageTextSubjectCol,
		table: messageTextTable,
	}
	MessageTextColGreeting = Column{
		name:  projection.MessageTextGreetingCol,
		table: messageTextTable,
	}
	MessageTextColText = Column{
		name:  projection.MessageTextTextCol,
		table: messageTextTable,
	}
	MessageTextColButtonText = Column{
		name:  projection.MessageTextButtonTextCol,
		table: messageTextTable,
	}
	MessageTextColFooter = Column{
		name:  projection.MessageTextFooterCol,
		table: messageTextTable,
	}
	MessageTextColOwnerRemoved = Column{
		name:  projection.MessageTextOwnerRemovedCol,
		table: messageTextTable,
	}
)

func (q *Queries) DefaultMessageText(ctx context.Context) (text *MessageText, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareMessageTextQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		MessageTextColAggregateID.identifier(): authz.GetInstance(ctx).InstanceID(),
		MessageTextColInstanceID.identifier():  authz.GetInstance(ctx).InstanceID(),
	}).
		Limit(1).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-1b9mf", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		text, err = scan(row)
		return err
	}, query, args...)
	return text, err
}

func (q *Queries) DefaultMessageTextByTypeAndLanguageFromFileSystem(ctx context.Context, messageType, language string) (_ *MessageText, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	contents, err := q.readNotificationTextMessages(ctx, language)
	if err != nil {
		return nil, err
	}
	messageTexts := new(MessageTexts)
	if err := yaml.Unmarshal(contents, messageTexts); err != nil {
		return nil, zerrors.ThrowInternal(err, "TEXT-3N9fs", "Errors.TranslationFile.ReadError")
	}
	return messageTexts.GetMessageTextByType(messageType), nil
}

func (q *Queries) CustomMessageTextByTypeAndLanguage(ctx context.Context, aggregateID, messageType, language string, withOwnerRemoved bool) (msg *MessageText, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareMessageTextQuery(ctx, q.client)
	eq := sq.Eq{
		MessageTextColLanguage.identifier():    language,
		MessageTextColType.identifier():        messageType,
		MessageTextColAggregateID.identifier(): aggregateID,
		MessageTextColInstanceID.identifier():  authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[MessageTextColOwnerRemoved.identifier()] = false
	}

	query, args, err := stmt.Where(eq).OrderBy(MessageTextColAggregateID.identifier()).Limit(1).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-1b9mf", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		msg, err = scan(row)
		return err
	}, query, args...)
	if zerrors.IsNotFound(err) {
		return q.IAMMessageTextByTypeAndLanguage(ctx, messageType, language)
	}
	return msg, err
}

func (q *Queries) IAMMessageTextByTypeAndLanguage(ctx context.Context, messageType, language string) (_ *MessageText, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	contents, err := q.readNotificationTextMessages(ctx, language)
	if err != nil {
		return nil, err
	}
	notificationTextMap := make(map[string]interface{})
	if err := yaml.Unmarshal(contents, &notificationTextMap); err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-ekjFF", "Errors.TranslationFile.ReadError")
	}
	texts, err := q.CustomTextList(ctx, authz.GetInstance(ctx).InstanceID(), messageType, language, false)
	if err != nil {
		return nil, err
	}
	for _, text := range texts.CustomTexts {
		messageTextMap, ok := notificationTextMap[messageType].(map[string]interface{})
		if !ok {
			continue
		}
		messageTextMap[text.Key] = text.Text
	}
	jsonbody, err := json.Marshal(notificationTextMap)

	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-3m8fJ", "Errors.TranslationFile.MergeError")
	}
	notificationText := new(MessageTexts)
	if err := json.Unmarshal(jsonbody, &notificationText); err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-9MkfD", "Errors.TranslationFile.MergeError")
	}
	result := notificationText.GetMessageTextByType(messageType)
	result.IsDefault = true
	result.AggregateID = authz.GetInstance(ctx).InstanceID()
	return result, nil
}

func (q *Queries) readNotificationTextMessages(ctx context.Context, language string) ([]byte, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	var err error
	contents, ok := q.NotificationTranslationFileContents[language]
	if !ok {
		contents, err = q.readTranslationFile(i18n.NOTIFICATION, fmt.Sprintf("/i18n/%s.yaml", language))
		if zerrors.IsNotFound(err) {
			contents, err = q.readTranslationFile(i18n.NOTIFICATION, fmt.Sprintf("/i18n/%s.yaml", authz.GetInstance(ctx).DefaultLanguage().String()))
		}
		if err != nil {
			return nil, err
		}
		q.NotificationTranslationFileContents[language] = contents
	}
	return contents, nil
}

func prepareMessageTextQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*MessageText, error)) {
	return sq.Select(
			MessageTextColAggregateID.identifier(),
			MessageTextColSequence.identifier(),
			MessageTextColCreationDate.identifier(),
			MessageTextColChangeDate.identifier(),
			MessageTextColState.identifier(),
			MessageTextColType.identifier(),
			MessageTextColLanguage.identifier(),
			MessageTextColTitle.identifier(),
			MessageTextColPreHeader.identifier(),
			MessageTextColSubject.identifier(),
			MessageTextColGreeting.identifier(),
			MessageTextColText.identifier(),
			MessageTextColButtonText.identifier(),
			MessageTextColFooter.identifier(),
		).
			From(messageTextTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*MessageText, error) {
			msg := new(MessageText)
			lang := ""
			title := sql.NullString{}
			preHeader := sql.NullString{}
			subject := sql.NullString{}
			greeting := sql.NullString{}
			text := sql.NullString{}
			buttonText := sql.NullString{}
			footer := sql.NullString{}
			err := row.Scan(
				&msg.AggregateID,
				&msg.Sequence,
				&msg.CreationDate,
				&msg.ChangeDate,
				&msg.State,
				&msg.Type,
				&lang,
				&title,
				&preHeader,
				&subject,
				&greeting,
				&text,
				&buttonText,
				&footer,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-3nlrS", "Errors.MessageText.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-499gJ", "Errors.Internal")
			}
			msg.Language = language.Make(lang)
			msg.Title = title.String
			msg.PreHeader = preHeader.String
			msg.Subject = subject.String
			msg.Greeting = greeting.String
			msg.Text = text.String
			msg.ButtonText = buttonText.String
			msg.Footer = footer.String
			return msg, nil
		}
}

func (q *Queries) readTranslationFile(namespace i18n.Namespace, filename string) ([]byte, error) {
	r, err := i18n.LoadFilesystem(namespace).Open(filename)
	if os.IsNotExist(err) {
		return nil, zerrors.ThrowNotFound(err, "QUERY-sN9wg", "Errors.TranslationFile.NotFound")
	}
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-93njw", "Errors.TranslationFile.ReadError")
	}
	contents, err := io.ReadAll(r)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-l0fse", "Errors.TranslationFile.ReadError")
	}
	return contents, nil
}

func (m *MessageTexts) GetMessageTextByType(msgType string) *MessageText {
	switch msgType {
	case domain.InitCodeMessageType:
		return &m.InitCode
	case domain.PasswordResetMessageType:
		return &m.PasswordReset
	case domain.VerifyEmailMessageType:
		return &m.VerifyEmail
	case domain.VerifyPhoneMessageType:
		return &m.VerifyPhone
	case domain.VerifySMSOTPMessageType:
		return &m.VerifySMSOTP
	case domain.VerifyEmailOTPMessageType:
		return &m.VerifyEmailOTP
	case domain.DomainClaimedMessageType:
		return &m.DomainClaimed
	case domain.PasswordlessRegistrationMessageType:
		return &m.PasswordlessRegistration
	case domain.PasswordChangeMessageType:
		return &m.PasswordChange
	case domain.InviteUserMessageType:
		return &m.InviteUser
	}
	return nil
}
