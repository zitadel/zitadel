package query

import (
	"context"
	"database/sql"
	errs "errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
	"github.com/ghodss/yaml"
	"golang.org/x/text/language"
)

type MessageTexts struct {
	InitCode                 MessageText
	PasswordReset            MessageText
	VerifyEmail              MessageText
	VerifyPhone              MessageText
	DomainClaimed            MessageText
	PasswordlessRegistration MessageText
}

type MessageText struct {
	AggregateID  string
	Sequence     uint64
	CreationDate time.Time
	ChangeDate   time.Time
	State        domain.PolicyState

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
		name: projection.MessageTextTable,
	}
	MessageTextColAggregateID = Column{
		name: projection.MessageTextAggregateIDCol,
	}
	MessageTextColSequence = Column{
		name: projection.MessageTextSequenceCol,
	}
	MessageTextColCreationDate = Column{
		name: projection.MessageTextCreationDateCol,
	}
	MessageTextColChangeDate = Column{
		name: projection.MessageTextChangeDateCol,
	}
	MessageTextColState = Column{
		name: projection.MessageTextStateCol,
	}
	MessageTextColType = Column{
		name: projection.MessageTextTypeCol,
	}
	MessageTextColLanguage = Column{
		name: projection.MessageTextLanguageCol,
	}
	MessageTextColTitle = Column{
		name: projection.MessageTextTitleCol,
	}
	MessageTextColPreHeader = Column{
		name: projection.MessageTextPreHeaderCol,
	}
	MessageTextColSubject = Column{
		name: projection.MessageTextSubjectCol,
	}
	MessageTextColGreeting = Column{
		name: projection.MessageTextGreetingCol,
	}
	MessageTextColText = Column{
		name: projection.MessageTextTextCol,
	}
	MessageTextColButtonText = Column{
		name: projection.MessageTextButtonTextCol,
	}
	MessageTextColFooter = Column{
		name: projection.MessageTextFooterCol,
	}
)

func (q *Queries) MessageTextByOrg(ctx context.Context, orgID string) (*MessageText, error) {
	stmt, scan := prepareMessageTextQuery()
	query, args, err := stmt.Where(
		sq.Or{
			sq.Eq{
				MessageTextColAggregateID.identifier(): orgID,
			},
			sq.Eq{
				MessageTextColAggregateID.identifier(): q.iamID,
			},
		}).
		//OrderBy(MessageTextColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-90n3N", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultMessageText(ctx context.Context) (*MessageText, error) {
	stmt, scan := prepareMessageTextQuery()
	query, args, err := stmt.Where(sq.Eq{
		MessageTextColAggregateID.identifier(): q.iamID,
	}).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-1b9mf", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultMessageTextByTypeAndLanguageFromFileSystem(messageType, language string) (*MessageText, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	var err error
	contents, ok := q.NotificationTranslationFileContents[language]
	if !ok {
		contents, err = q.readTranslationFile(q.NotificationDir, fmt.Sprintf("/i18n/%s.yaml", language))
		if errors.IsNotFound(err) {
			contents, err = q.readTranslationFile(q.NotificationDir, fmt.Sprintf("/i18n/%s.yaml", q.DefaultLanguage.String()))
		}
		if err != nil {
			return nil, err
		}
		q.NotificationTranslationFileContents[language] = contents
	}
	messageTexts := new(MessageTexts)
	if err := yaml.Unmarshal(contents, messageTexts); err != nil {
		return nil, errors.ThrowInternal(err, "TEXT-3N9fs", "Errors.TranslationFile.ReadError")
	}
	return messageTexts.GetMessageTextByType(messageType), nil
}

func (q *Queries) IAMMessageTextByTypeAndLanguage(ctx context.Context, messageType, language string) (*MessageText, error) {
	stmt, scan := prepareMessageTextQuery()
	query, args, err := stmt.Where(
		sq.Eq{
			MessageTextColAggregateID.identifier(): q.iamID,
		},
		sq.Eq{
			MessageTextColType.identifier(): messageType,
		},
		sq.Eq{
			MessageTextColLanguage.identifier(): language,
		},
	).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-1b9mf", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func prepareMessageTextQuery() (sq.SelectBuilder, func(*sql.Row) (*MessageText, error)) {
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
			From(messageTextTable.identifier()).PlaceholderFormat(sq.Dollar),
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
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-3nlrS", "Errors.MessageText.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-499gJ", "Errors.Internal")
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

func (q *Queries) readTranslationFile(dir http.FileSystem, filename string) ([]byte, error) {
	r, err := dir.Open(filename)
	if os.IsNotExist(err) {
		return nil, errors.ThrowNotFound(err, "QUERY-sN9wg", "Errors.TranslationFile.NotFound")
	}
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-93njw", "Errors.TranslationFile.ReadError")
	}
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-l0fse", "Errors.TranslationFile.ReadError")
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
	case domain.DomainClaimedMessageType:
		return &m.DomainClaimed
	case domain.PasswordlessRegistrationMessageType:
		return &m.PasswordlessRegistration
	}
	return nil
}
