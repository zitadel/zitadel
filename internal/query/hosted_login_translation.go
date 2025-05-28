package query

import (
	"context"
	"crypto/md5"
	_ "embed"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/structpb"
)

//go:embed v2-default.json
var defaultLoginTranslations []byte

type HostedLoginTranslationLevelType uint

const (
	HostedLoginTranslationLevelInstance HostedLoginTranslationLevelType = iota + 1
	HostedLoginTranslationLevelOrg
)

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

	translation, err := getSystemTranslation(req.GetLocale())

	protoTranslation, err := structpb.NewStruct(translation)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-70ppPp", "Errors.Protobuf.ConvertToStruct")
	}

	etag := md5.Sum([]byte(protoTranslation.String()))

	res.Etag = string(etag[:])
	res.Translations = protoTranslation

	return
}

func getSystemTranslation(lang string) (map[string]any, error) {
	defaultTranslations := map[string]any{}

	language := language.BCP47.Make(lang)

	err := json.Unmarshal(defaultLoginTranslations, &defaultTranslations)
	if err != nil {
		zerrors.ThrowInternal(err, "QUERY-nvx88W", "Errors.Query.UnmarshalDefaultLoginTranslations")
	}

	translation, ok := defaultTranslations[language.String()]
	if !ok {
		return nil, zerrors.ThrowNotFoundf(nil, "QUERY-6gb5QR", "Errors.Query.HostedLoginTranslationNotFound-%s", language.String())
	}

	castedTranslation, ok := translation.(map[string]any)
	if !ok {
		return nil, zerrors.ThrowInternal(nil, "QUERY-WrRn5e", "Errors.Query.HostedLoginCastError")
	}

	return castedTranslation, nil
}
