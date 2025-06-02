package command

import (
	"context"
	"crypto/md5"
	"encoding/hex"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	level2AggregateMapper map[settings.TranslationLevelType]func(string) eventstore.Aggregate = map[settings.TranslationLevelType]func(string) eventstore.Aggregate{
		settings.TranslationLevelType_TRANSLATION_LEVEL_TYPE_INSTANCE: func(resourceID string) eventstore.Aggregate {
			return instance.NewAggregate(resourceID).Aggregate
		},
		settings.TranslationLevelType_TRANSLATION_LEVEL_TYPE_ORG: func(resourceID string) eventstore.Aggregate {
			return org.NewAggregate(resourceID).Aggregate
		},
	}
)

func (c *Commands) SetHostedLoginTranslation(ctx context.Context, req *settings.SetHostedLoginTranslationRequest) (res *settings.SetHostedLoginTranslationResponse, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	aggregateFunc, ok := level2AggregateMapper[req.GetLevel()]
	if !ok {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMA-YB6Sri", "Errors.Arguments.LevelType.Invalid")
	}
	agg := aggregateFunc(req.GetLevelId())

	lang, err := language.BCP47.Parse(req.GetLocale())
	if err != nil || lang.IsRoot() {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMA-xmjATA", "Errors.Arguments.Locale.Invalid")
	}
	baseLang, _ := lang.Base()

	commands, wm, err := c.setTranslation(ctx, agg, baseLang, req.GetTranslations().AsMap())

	pushedEvents, err := c.eventstore.Push(ctx, commands...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "COMMA-i8nqFl", "Errors.Internal")
	}

	err = AppendAndReduce(wm, pushedEvents...)
	if err != nil {
		return nil, err
	}

	protoTranslation, err := structpb.NewStruct(wm.Translation)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-70ppPp", "Errors.Protobuf.ConvertToStruct")
	}

	etag := md5.Sum([]byte(protoTranslation.String()))
	return &settings.SetHostedLoginTranslationResponse{
		Etag: hex.EncodeToString(etag[:]),
	}, nil
}

func (c *Commands) setTranslation(ctx context.Context, agg eventstore.Aggregate, lang language.Base, translations map[string]any) ([]eventstore.Command, *HostedLoginTranslationWriteModel, error) {
	wm := NewHostedLoginTranslationWriteModel(agg.ID)
	events := []eventstore.Command{}
	switch agg.Type {
	case "instance":
	case "org":
		events = append(events, org.NewHostedLoginTranslationSetEvent(ctx, &agg, translations, lang.String()))
	default:
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "COMMA-0aw7In", "Errors.Arguments.LevelType.Invalid")
	}

	return events, wm, nil
}
