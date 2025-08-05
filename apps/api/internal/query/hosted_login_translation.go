package query

import (
	"context"
	"crypto/md5"
	"database/sql"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"dario.cat/mergo"
	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/v2/org"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

var (
	//go:embed v2-default.json
	defaultLoginTranslations []byte

	defaultSystemTranslations map[language.Tag]map[string]any

	hostedLoginTranslationTable = table{
		name:          projection.HostedLoginTranslationTable,
		instanceIDCol: projection.HostedLoginTranslationInstanceIDCol,
	}

	hostedLoginTranslationColInstanceID = Column{
		name:  projection.HostedLoginTranslationInstanceIDCol,
		table: hostedLoginTranslationTable,
	}
	hostedLoginTranslationColResourceOwner = Column{
		name:  projection.HostedLoginTranslationAggregateIDCol,
		table: hostedLoginTranslationTable,
	}
	hostedLoginTranslationColResourceOwnerType = Column{
		name:  projection.HostedLoginTranslationAggregateTypeCol,
		table: hostedLoginTranslationTable,
	}
	hostedLoginTranslationColLocale = Column{
		name:  projection.HostedLoginTranslationLocaleCol,
		table: hostedLoginTranslationTable,
	}
	hostedLoginTranslationColFile = Column{
		name:  projection.HostedLoginTranslationFileCol,
		table: hostedLoginTranslationTable,
	}
	hostedLoginTranslationColEtag = Column{
		name:  projection.HostedLoginTranslationEtagCol,
		table: hostedLoginTranslationTable,
	}
)

func init() {
	err := json.Unmarshal(defaultLoginTranslations, &defaultSystemTranslations)
	if err != nil {
		panic(err)
	}
}

type HostedLoginTranslations struct {
	SearchResponse
	HostedLoginTranslations []*HostedLoginTranslation
}

type HostedLoginTranslation struct {
	AggregateID  string
	Sequence     uint64
	CreationDate time.Time
	ChangeDate   time.Time

	Locale    string
	File      map[string]any
	LevelType string
	LevelID   string
	Etag      string
}

func (q *Queries) GetHostedLoginTranslation(ctx context.Context, req *settings.GetHostedLoginTranslationRequest) (res *settings.GetHostedLoginTranslationResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	inst := authz.GetInstance(ctx)
	defaultInstLang := inst.DefaultLanguage()

	lang, err := language.BCP47.Parse(req.GetLocale())
	if err != nil || lang.IsRoot() {
		return nil, zerrors.ThrowInvalidArgument(nil, "QUERY-rZLAGi", "Errors.Arguments.Locale.Invalid")
	}
	parentLang := lang.Parent()
	if parentLang.IsRoot() {
		parentLang = lang
	}

	sysTranslation, systemEtag, err := getSystemTranslation(parentLang, defaultInstLang)
	if err != nil {
		return nil, err
	}

	var levelID, resourceOwner string
	switch t := req.GetLevel().(type) {
	case *settings.GetHostedLoginTranslationRequest_System:
		return getTranslationOutputMessage(sysTranslation, systemEtag)
	case *settings.GetHostedLoginTranslationRequest_Instance:
		levelID = authz.GetInstance(ctx).InstanceID()
		resourceOwner = instance.AggregateType
	case *settings.GetHostedLoginTranslationRequest_OrganizationId:
		levelID = t.OrganizationId
		resourceOwner = org.AggregateType
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "QUERY-YB6Sri", "Errors.Arguments.Level.Invalid")
	}

	stmt, scan := prepareHostedLoginTranslationQuery()

	langORBaseLang := sq.Or{
		sq.Eq{hostedLoginTranslationColLocale.identifier(): lang.String()},
		sq.Eq{hostedLoginTranslationColLocale.identifier(): parentLang.String()},
	}
	eq := sq.Eq{
		hostedLoginTranslationColInstanceID.identifier():        inst.InstanceID(),
		hostedLoginTranslationColResourceOwner.identifier():     levelID,
		hostedLoginTranslationColResourceOwnerType.identifier(): resourceOwner,
	}

	query, args, err := stmt.Where(eq).Where(langORBaseLang).ToSql()
	if err != nil {
		logging.WithError(err).Error("unable to generate sql statement")
		return nil, zerrors.ThrowInternal(err, "QUERY-ZgCMux", "Errors.Query.SQLStatement")
	}

	var trs []*HostedLoginTranslation
	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		trs, err = scan(rows)
		return err
	}, query, args...)
	if err != nil {
		logging.WithError(err).Error("failed to query translations")
		return nil, zerrors.ThrowInternal(err, "QUERY-6k1zjx", "Errors.Internal")
	}

	requestedTranslation, parentTranslation := &HostedLoginTranslation{}, &HostedLoginTranslation{}
	for _, tr := range trs {
		if tr == nil {
			continue
		}

		if tr.LevelType == resourceOwner {
			requestedTranslation = tr
		} else {
			parentTranslation = tr
		}
	}

	if !req.GetIgnoreInheritance() {

		// There is no record for the requested level, set the upper level etag
		if requestedTranslation.Etag == "" {
			requestedTranslation.Etag = parentTranslation.Etag
		}

		// Case where Level == ORGANIZATION -> Check if we have an instance level translation
		// If so, merge it with the translations we have
		if parentTranslation != nil && parentTranslation.LevelType == instance.AggregateType {
			if err := mergo.Merge(&requestedTranslation.File, parentTranslation.File); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-pdgEJd", "Errors.Query.MergeTranslations")
			}
		}

		// The DB query returned no results, we have to set the system translation etag
		if requestedTranslation.Etag == "" {
			requestedTranslation.Etag = systemEtag
		}

		// Merge the system translations
		if err := mergo.Merge(&requestedTranslation.File, sysTranslation); err != nil {
			return nil, zerrors.ThrowInternal(err, "QUERY-HdprNF", "Errors.Query.MergeTranslations")
		}
	}

	return getTranslationOutputMessage(requestedTranslation.File, requestedTranslation.Etag)
}

func getSystemTranslation(lang, instanceDefaultLang language.Tag) (map[string]any, string, error) {
	translation, ok := defaultSystemTranslations[lang]
	if !ok {
		translation, ok = defaultSystemTranslations[instanceDefaultLang]
		if !ok {
			return nil, "", zerrors.ThrowNotFoundf(nil, "QUERY-6gb5QR", "Errors.Query.HostedLoginTranslationNotFound-%s", lang)
		}
	}

	hash := md5.Sum(fmt.Append(nil, translation))

	return translation, hex.EncodeToString(hash[:]), nil
}

func prepareHostedLoginTranslationQuery() (sq.SelectBuilder, func(*sql.Rows) ([]*HostedLoginTranslation, error)) {
	return sq.Select(
			hostedLoginTranslationColFile.identifier(),
			hostedLoginTranslationColResourceOwnerType.identifier(),
			hostedLoginTranslationColEtag.identifier(),
		).From(hostedLoginTranslationTable.identifier()).
			Limit(2).
			PlaceholderFormat(sq.Dollar),
		func(r *sql.Rows) ([]*HostedLoginTranslation, error) {
			translations := make([]*HostedLoginTranslation, 0, 2)
			for r.Next() {
				var rawTranslation json.RawMessage
				translation := &HostedLoginTranslation{}
				err := r.Scan(
					&rawTranslation,
					&translation.LevelType,
					&translation.Etag,
				)
				if err != nil {
					return nil, err
				}

				if err := json.Unmarshal(rawTranslation, &translation.File); err != nil {
					return nil, err
				}

				translations = append(translations, translation)
			}

			if err := r.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-oc7r7i", "Errors.Query.CloseRows")
			}

			return translations, nil
		}
}

func getTranslationOutputMessage(translation map[string]any, etag string) (*settings.GetHostedLoginTranslationResponse, error) {
	protoTranslation, err := structpb.NewStruct(translation)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-70ppPp", "Errors.Protobuf.ConvertToStruct")
	}

	return &settings.GetHostedLoginTranslationResponse{
		Translations: protoTranslation,
		Etag:         etag,
	}, nil
}
