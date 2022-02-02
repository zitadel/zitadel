package admin

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	obj_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/notification/channels/smtp"
	"github.com/caos/zitadel/internal/query"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
	settings_pb "github.com/caos/zitadel/pkg/grpc/settings"
	"google.golang.org/protobuf/types/known/durationpb"
)

func listSecretGeneratorToModel(req *admin_pb.ListSecretGeneratorsRequest) (*query.SecretGeneratorSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := SecretGeneratorQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	return &query.SecretGeneratorSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}

func SecretGeneratorQueriesToModel(queries []*settings_pb.SecretGeneratorQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = SecretGeneratorQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func SecretGeneratorQueryToModel(apiQuery *settings_pb.SecretGeneratorQuery) (query.SearchQuery, error) {
	switch q := apiQuery.Query.(type) {
	case *settings_pb.SecretGeneratorQuery_TypeQuery:
		return query.NewProjectNameSearchQuery(object.TextMethodToQuery(q.TypeQuery.Method), q.TypeQuery.GeneratorType)
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ORG-fm9es", "List.Query.Invalid")
	}
}

func UpdateSecretGeneratorToConfig(req *admin_pb.UpdateSecretGeneratorRequest) *crypto.GeneratorConfig {
	return &crypto.GeneratorConfig{
		Length:              uint(req.Length),
		Expiry:              types.Duration{Duration: req.Expiry.AsDuration()},
		IncludeUpperLetters: req.IncludeUpperLetters,
		IncludeLowerLetters: req.IncludeLowerLetters,
		IncludeDigits:       req.IncludeDigits,
		IncludeSymbols:      req.IncludeSymbols,
	}
}

func SecretGeneratorToPb(generator *query.SecretGenerator) *settings_pb.SecretGenerator {
	mapped := &settings_pb.SecretGenerator{
		Length:              uint32(generator.Length),
		Expiry:              durationpb.New(generator.Expiry),
		IncludeUpperLetters: generator.IncludeUpperLetters,
		IncludeLowerLetters: generator.IncludeLowerLetters,
		IncludeDigits:       generator.IncludeDigits,
		IncludeSymbols:      generator.IncludeSymbols,
		Details:             obj_grpc.ToViewDetailsPb(generator.Sequence, generator.CreationDate, generator.ChangeDate, generator.ID),
	}
	return mapped
}

func UpdateSMTPToConfig(req *admin_pb.UpdateSMTPConfigRequest) *smtp.EmailConfig {
	return &smtp.EmailConfig{
		Tls:      req.Tls,
		From:     req.FromAddress,
		FromName: req.FromName,
		SMTP: smtp.SMTP{
			Host:     req.SmtpHost,
			User:     req.SmtpUser,
			Password: req.SmtpPassword,
		},
	}
}
