package admin

import (
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/api/grpc/object"
	obj_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	settings_pb "github.com/zitadel/zitadel/pkg/grpc/settings"
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
		domainType := SecretGeneratorTypeToDomain(q.TypeQuery.GeneratorType)
		return query.NewSecretGeneratorTypeSearchQuery(int32(domainType))
	default:
		return nil, errors.ThrowInvalidArgument(nil, "ORG-fm9es", "List.Query.Invalid")
	}
}

func UpdateSecretGeneratorToConfig(req *admin_pb.UpdateSecretGeneratorRequest) *crypto.GeneratorConfig {
	return &crypto.GeneratorConfig{
		Length:              uint(req.Length),
		Expiry:              req.Expiry.AsDuration(),
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
		Details:             obj_grpc.ToViewDetailsPb(generator.Sequence, generator.CreationDate, generator.ChangeDate, generator.AggregateID),
	}
	return mapped
}

func SecretGeneratorTypeToPb(generatorType domain.SecretGeneratorType) settings_pb.SecretGeneratorType {
	switch generatorType {
	case domain.SecretGeneratorTypeInitCode:
		return settings_pb.SecretGeneratorType_SECRET_GENERATOR_TYPE_INIT_CODE
	case domain.SecretGeneratorTypeVerifyEmailCode:
		return settings_pb.SecretGeneratorType_SECRET_GENERATOR_TYPE_VERIFY_EMAIL_CODE
	case domain.SecretGeneratorTypeVerifyPhoneCode:
		return settings_pb.SecretGeneratorType_SECRET_GENERATOR_TYPE_VERIFY_PHONE_CODE
	case domain.SecretGeneratorTypePasswordResetCode:
		return settings_pb.SecretGeneratorType_SECRET_GENERATOR_TYPE_PASSWORD_RESET_CODE
	case domain.SecretGeneratorTypePasswordlessInitCode:
		return settings_pb.SecretGeneratorType_SECRET_GENERATOR_TYPE_PASSWORDLESS_INIT_CODE
	case domain.SecretGeneratorTypeAppSecret:
		return settings_pb.SecretGeneratorType_SECRET_GENERATOR_TYPE_APP_SECRET
	default:
		return settings_pb.SecretGeneratorType_SECRET_GENERATOR_TYPE_UNSPECIFIED
	}
}

func SecretGeneratorTypeToDomain(generatorType settings_pb.SecretGeneratorType) domain.SecretGeneratorType {
	switch generatorType {
	case settings_pb.SecretGeneratorType_SECRET_GENERATOR_TYPE_INIT_CODE:
		return domain.SecretGeneratorTypeInitCode
	case settings_pb.SecretGeneratorType_SECRET_GENERATOR_TYPE_VERIFY_EMAIL_CODE:
		return domain.SecretGeneratorTypeVerifyEmailCode
	case settings_pb.SecretGeneratorType_SECRET_GENERATOR_TYPE_VERIFY_PHONE_CODE:
		return domain.SecretGeneratorTypeVerifyPhoneCode
	case settings_pb.SecretGeneratorType_SECRET_GENERATOR_TYPE_PASSWORD_RESET_CODE:
		return domain.SecretGeneratorTypePasswordResetCode
	case settings_pb.SecretGeneratorType_SECRET_GENERATOR_TYPE_PASSWORDLESS_INIT_CODE:
		return domain.SecretGeneratorTypePasswordlessInitCode
	case settings_pb.SecretGeneratorType_SECRET_GENERATOR_TYPE_APP_SECRET:
		return domain.SecretGeneratorTypeAppSecret
	default:
		return domain.SecretGeneratorTypeUnspecified
	}
}

func AddSMTPToConfig(req *admin_pb.AddSMTPConfigRequest) *smtp.Config {
	return &smtp.Config{
		Tls:      req.Tls,
		From:     req.SenderAddress,
		FromName: req.SenderName,
		SMTP: smtp.SMTP{
			Host:     req.Host,
			User:     req.User,
			Password: req.Password,
		},
	}
}

func UpdateSMTPToConfig(req *admin_pb.UpdateSMTPConfigRequest) *smtp.Config {
	return &smtp.Config{
		Tls:      req.Tls,
		From:     req.SenderAddress,
		FromName: req.SenderName,
		SMTP: smtp.SMTP{
			Host: req.Host,
			User: req.User,
		},
	}
}

func SMTPConfigToPb(smtp *query.SMTPConfig) *settings_pb.SMTPConfig {
	mapped := &settings_pb.SMTPConfig{
		Tls:           smtp.TLS,
		SenderAddress: smtp.SenderAddress,
		SenderName:    smtp.SenderName,
		Host:          smtp.Host,
		User:          smtp.User,
		Details:       obj_grpc.ToViewDetailsPb(smtp.Sequence, smtp.CreationDate, smtp.ChangeDate, smtp.AggregateID),
	}
	return mapped
}

func SecurityPolicyToPb(policy *query.SecurityPolicy) *settings_pb.SecurityPolicy {
	return &settings_pb.SecurityPolicy{
		Details:               obj_grpc.ToViewDetailsPb(policy.Sequence, policy.CreationDate, policy.ChangeDate, policy.AggregateID),
		EnableIframeEmbedding: policy.Enabled,
		AllowedOrigins:        policy.AllowedOrigins,
	}
}
