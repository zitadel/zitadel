package query

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	prepareUserSchemasStmt = `SELECT projections.user_schemas1.id,` +
		` projections.user_schemas1.creation_date,` +
		` projections.user_schemas1.change_date,` +
		` projections.user_schemas1.sequence,` +
		` projections.user_schemas1.instance_id,` +
		` projections.user_schemas1.state,` +
		` projections.user_schemas1.type,` +
		` projections.user_schemas1.revision,` +
		` projections.user_schemas1.schema,` +
		` projections.user_schemas1.possible_authenticators,` +
		` COUNT(*) OVER ()` +
		` FROM projections.user_schemas`
	prepareUserSchemasCols = []string{
		"id",
		"creation_date",
		"change_date",
		"sequence",
		"instance_id",
		"state",
		"type",
		"revision",
		"schema",
		"possible_authenticators",
		"count",
	}

	prepareUserSchemaStmt = `SELECT projections.user_schemas1.id,` +
		` projections.user_schemas1.creation_date,` +
		` projections.user_schemas1.change_date,` +
		` projections.user_schemas1.sequence,` +
		` projections.user_schemas1.instance_id,` +
		` projections.user_schemas1.state,` +
		` projections.user_schemas1.type,` +
		` projections.user_schemas1.revision,` +
		` projections.user_schemas1.schema,` +
		` projections.user_schemas1.possible_authenticators` +
		` FROM projections.user_schemas`
	prepareUserSchemaCols = []string{
		"id",
		"creation_date",
		"change_date",
		"sequence",
		"instance_id",
		"state",
		"type",
		"revision",
		"schema",
		"possible_authenticators",
	}
)

func Test_UserSchemaPrepares(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareUserSchemasQuery no result",
			prepare: prepareUserSchemasQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareUserSchemasStmt),
					nil,
					nil,
				),
			},
			object: &UserSchemas{UserSchemas: []*UserSchema{}},
		},
		{
			name:    "prepareUserSchemasQuery one result",
			prepare: prepareUserSchemasQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareUserSchemasStmt),
					prepareUserSchemasCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							uint64(20211109),
							"instance-id",
							domain.UserSchemaStateActive,
							"type",
							1,
							json.RawMessage(`{"$schema":"urn:zitadel:schema:v1","properties":{"name":{"type":"string","urn:zitadel:schema:permission":{"self":"rw"}}},"type":"object"}`),
							database.NumberArray[domain.AuthenticatorType]{domain.AuthenticatorTypeUsername, domain.AuthenticatorTypePassword},
						},
					},
				),
			},
			object: &UserSchemas{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				UserSchemas: []*UserSchema{
					{
						ObjectDetails: domain.ObjectDetails{
							ID:            "id",
							EventDate:     testNow,
							CreationDate:  testNow,
							Sequence:      20211109,
							ResourceOwner: "instance-id",
						},
						State:                  domain.UserSchemaStateActive,
						Type:                   "type",
						Revision:               1,
						Schema:                 json.RawMessage(`{"$schema":"urn:zitadel:schema:v1","properties":{"name":{"type":"string","urn:zitadel:schema:permission":{"self":"rw"}}},"type":"object"}`),
						PossibleAuthenticators: database.NumberArray[domain.AuthenticatorType]{domain.AuthenticatorTypeUsername, domain.AuthenticatorTypePassword},
					},
				},
			},
		},
		{
			name:    "prepareUserSchemasQuery multiple result",
			prepare: prepareUserSchemasQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareUserSchemasStmt),
					prepareUserSchemasCols,
					[][]driver.Value{
						{
							"id-1",
							testNow,
							testNow,
							uint64(20211109),
							"instance-id",
							domain.UserSchemaStateActive,
							"type1",
							1,
							json.RawMessage(`{"$schema":"urn:zitadel:schema:v1","properties":{"name":{"type":"string","urn:zitadel:schema:permission":{"self":"rw"}}},"type":"object"}`),
							database.NumberArray[domain.AuthenticatorType]{domain.AuthenticatorTypeUsername, domain.AuthenticatorTypePassword},
						},
						{
							"id-2",
							testNow,
							testNow,
							uint64(20211110),
							"instance-id",
							domain.UserSchemaStateInactive,
							"type2",
							2,
							json.RawMessage(`{"$schema":"urn:zitadel:schema:v1","properties":{"name":{"type":"string","urn:zitadel:schema:permission":{"self":"rw"}}},"type":"object"}`),
							database.NumberArray[domain.AuthenticatorType]{domain.AuthenticatorTypeUsername, domain.AuthenticatorTypePassword},
						},
					},
				),
			},
			object: &UserSchemas{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				UserSchemas: []*UserSchema{
					{
						ObjectDetails: domain.ObjectDetails{
							ID:            "id-1",
							EventDate:     testNow,
							CreationDate:  testNow,
							Sequence:      20211109,
							ResourceOwner: "instance-id",
						},
						State:                  domain.UserSchemaStateActive,
						Type:                   "type1",
						Revision:               1,
						Schema:                 json.RawMessage(`{"$schema":"urn:zitadel:schema:v1","properties":{"name":{"type":"string","urn:zitadel:schema:permission":{"self":"rw"}}},"type":"object"}`),
						PossibleAuthenticators: database.NumberArray[domain.AuthenticatorType]{domain.AuthenticatorTypeUsername, domain.AuthenticatorTypePassword},
					},
					{
						ObjectDetails: domain.ObjectDetails{
							ID:            "id-2",
							EventDate:     testNow,
							CreationDate:  testNow,
							Sequence:      20211110,
							ResourceOwner: "instance-id",
						},
						State:                  domain.UserSchemaStateInactive,
						Type:                   "type2",
						Revision:               2,
						Schema:                 json.RawMessage(`{"$schema":"urn:zitadel:schema:v1","properties":{"name":{"type":"string","urn:zitadel:schema:permission":{"self":"rw"}}},"type":"object"}`),
						PossibleAuthenticators: database.NumberArray[domain.AuthenticatorType]{domain.AuthenticatorTypeUsername, domain.AuthenticatorTypePassword},
					},
				},
			},
		},
		{
			name:    "prepareUserSchemasQuery sql err",
			prepare: prepareUserSchemasQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareUserSchemasStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*UserSchema)(nil),
		},
		{
			name:    "prepareUserSchemaQuery no result",
			prepare: prepareUserSchemaQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					regexp.QuoteMeta(prepareUserSchemaStmt),
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !zerrors.IsNotFound(err) {
						return fmt.Errorf("err should be zitadel.NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*UserSchema)(nil),
		},
		{
			name:    "prepareUserSchemaQuery found",
			prepare: prepareUserSchemaQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareUserSchemaStmt),
					prepareUserSchemaCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						uint64(20211109),
						"instance-id",
						domain.UserSchemaStateActive,
						"type",
						1,
						json.RawMessage(`{"$schema":"urn:zitadel:schema:v1","properties":{"name":{"type":"string","urn:zitadel:schema:permission":{"self":"rw"}}},"type":"object"}`),
						database.NumberArray[domain.AuthenticatorType]{domain.AuthenticatorTypeUsername, domain.AuthenticatorTypePassword},
					},
				),
			},
			object: &UserSchema{
				ObjectDetails: domain.ObjectDetails{
					ID:            "id",
					EventDate:     testNow,
					CreationDate:  testNow,
					Sequence:      20211109,
					ResourceOwner: "instance-id",
				},
				State:                  domain.UserSchemaStateActive,
				Type:                   "type",
				Revision:               1,
				Schema:                 json.RawMessage(`{"$schema":"urn:zitadel:schema:v1","properties":{"name":{"type":"string","urn:zitadel:schema:permission":{"self":"rw"}}},"type":"object"}`),
				PossibleAuthenticators: database.NumberArray[domain.AuthenticatorType]{domain.AuthenticatorTypeUsername, domain.AuthenticatorTypePassword},
			},
		},
		{
			name:    "prepareUserSchemaQuery sql err",
			prepare: prepareUserSchemaQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareUserSchemaStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*UserSchema)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
