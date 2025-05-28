package query

import (
	"context"
	"crypto/md5"
	"database/sql"
	_ "embed"
	"encoding/json"
	"time"

	"dario.cat/mergo"
	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	//go:embed v2-default.json
	defaultLoginTranslations []byte

	hostedLoginTranslationTable = table{
		name:          projection.HostedLoginTranslationTable,
		instanceIDCol: projection.HostedLoginTranslationInstaceIDCol,
	}

	hostedLoginTranslationColInstanceID = Column{
		name:  projection.HostedLoginTranslationInstaceIDCol,
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
	hostedLoginTranslationColCreationDate = Column{
		name:  projection.HostedLoginTranslationCreationDateCol,
		table: hostedLoginTranslationTable,
	}
	hostedLoginTranslationColChangeDate = Column{
		name:  projection.HostedLoginTranslationChangeDateCol,
		table: hostedLoginTranslationTable,
	}
	hostedLoginTranslationColSequence = Column{
		name:  projection.HostedLoginTranslationSequenceCol,
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
)

type HostedLoginTranslationLevelType uint

const (
	HostedLoginTranslationLevelInstance HostedLoginTranslationLevelType = iota + 1
	HostedLoginTranslationLevelOrg
)

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
	LevelType HostedLoginTranslationLevelType
	LevelID   string
}

func (q *Queries) GetHostedLoginTranslation(ctx context.Context, req *settings.GetHostedLoginTranslationRequest) (res *settings.GetHostedLoginTranslationResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	inst := authz.GetInstance(ctx)
	defaultInstLang := inst.DefaultLanguage()

	lang := language.BCP47.Make(req.GetLocale())

	sysTranslation, err := getSystemTranslation(lang.String(), defaultInstLang.String())

	stmt, scan := prepareHostedLoginTranslationQuery()

	eq := sq.Eq{
		hostedLoginTranslationColInstanceID.identifier(): inst.InstanceID(),
		hostedLoginTranslationColLocale.identifier():     lang.String(),
	}

	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-ZgCMux", "Errors.Query.SQLStatement")
	}

	trs := make([]*HostedLoginTranslation, 2)
	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		trs, err = scan(rows)
		return err
	}, query, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-6k1zjx", "Errors.Internal")
	}

	var requestedTranslation, otherTranslation *HostedLoginTranslation
	for _, tr := range trs {
		if tr == nil {
			continue
		}

		if tr.LevelType == HostedLoginTranslationLevelType(req.GetLevel()) {
			requestedTranslation = tr
			break
		}
		otherTranslation = tr
	}

	// Case where req.GetLeve() == ORGANIZATION -> Check if we have an instance level translation
	// If so, merge it with the translations we have
	if otherTranslation != nil && requestedTranslation.LevelType > otherTranslation.LevelType {
		if err := mergo.Merge(&requestedTranslation.File, otherTranslation.File); err != nil {
			return nil, zerrors.ThrowInternal(err, "QUERY-pdgEJd", "Errors.Query.MergeTranslations")
		}
	}

	// Merge the system translations
	if err := mergo.Merge(&requestedTranslation.File, sysTranslation); err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-HdprNF", "Errors.Query.MergeTranslations")
	}

	protoTranslation, err := structpb.NewStruct(requestedTranslation.File)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-70ppPp", "Errors.Protobuf.ConvertToStruct")
	}

	etag := md5.Sum([]byte(protoTranslation.String()))

	res.Etag = string(etag[:])
	res.Translations = protoTranslation

	return
}

func getSystemTranslation(lang, instanceDefaultLang string) (map[string]any, error) {
	defaultTranslations := map[string]any{}

	err := json.Unmarshal(defaultLoginTranslations, &defaultTranslations)
	if err != nil {
		zerrors.ThrowInternal(err, "QUERY-nvx88W", "Errors.Query.UnmarshalDefaultLoginTranslations")
	}

	translation, ok := defaultTranslations[lang]
	if !ok {
		translation, ok = defaultTranslations[instanceDefaultLang]
		if !ok {
			return nil, zerrors.ThrowNotFoundf(nil, "QUERY-6gb5QR", "Errors.Query.HostedLoginTranslationNotFound-%s", lang)
		}
	}

	castedTranslation, ok := translation.(map[string]any)
	if !ok {
		return nil, zerrors.ThrowInternal(nil, "QUERY-WrRn5e", "Errors.Query.HostedLoginCastError")
	}

	return castedTranslation, nil
}

func prepareHostedLoginTranslationQuery() (sq.SelectBuilder, func(*sql.Rows) ([]*HostedLoginTranslation, error)) {
	return sq.Select(
			hostedLoginTranslationColFile.identifier(),
			hostedLoginTranslationColResourceOwnerType.identifier(),
		).From(hostedLoginTranslationTable.identifier()).
			Limit(2).
			PlaceholderFormat(sq.Dollar),
		func(r *sql.Rows) ([]*HostedLoginTranslation, error) {
			translations := make([]*HostedLoginTranslation, 2)
			for r.Next() {
				translation := &HostedLoginTranslation{}
				err := r.Scan(
					&translation.File,
					&translation.LevelType,
				)
				if err != nil {
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
