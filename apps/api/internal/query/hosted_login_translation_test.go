package query

import (
	"crypto/md5"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/database/mock"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

func TestGetSystemTranslation(t *testing.T) {
	okTranslation := defaultLoginTranslations

	parsedOKTranslation := map[string]map[string]any{}
	require.Nil(t, json.Unmarshal(okTranslation, &parsedOKTranslation))

	hashOK := md5.Sum(fmt.Append(nil, parsedOKTranslation["de"]))

	tt := []struct {
		testName string

		inputLanguage          language.Tag
		inputInstanceLanguage  language.Tag
		systemTranslationToSet []byte

		expectedLanguage map[string]any
		expectedEtag     string
		expectedError    error
	}{
		{
			testName:               "when neither input language nor system default language have translation should return not found error",
			systemTranslationToSet: okTranslation,
			inputLanguage:          language.MustParse("ro"),
			inputInstanceLanguage:  language.MustParse("fr"),

			expectedError: zerrors.ThrowNotFoundf(nil, "QUERY-6gb5QR", "Errors.Query.HostedLoginTranslationNotFound-%s", "ro"),
		},
		{
			testName:               "when input language has no translation should fallback onto instance default",
			systemTranslationToSet: okTranslation,
			inputLanguage:          language.MustParse("ro"),
			inputInstanceLanguage:  language.MustParse("de"),

			expectedLanguage: parsedOKTranslation["de"],
			expectedEtag:     hex.EncodeToString(hashOK[:]),
		},
		{
			testName:               "when input language has translation should return it",
			systemTranslationToSet: okTranslation,
			inputLanguage:          language.MustParse("de"),
			inputInstanceLanguage:  language.MustParse("en"),

			expectedLanguage: parsedOKTranslation["de"],
			expectedEtag:     hex.EncodeToString(hashOK[:]),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Given
			defaultLoginTranslations = tc.systemTranslationToSet

			// When
			translation, etag, err := getSystemTranslation(tc.inputLanguage, tc.inputInstanceLanguage)

			// Verify
			require.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedLanguage, translation)
			assert.Equal(t, tc.expectedEtag, etag)
		})
	}
}

func TestGetTranslationOutput(t *testing.T) {
	t.Parallel()

	validMap := map[string]any{"loginHeader": "A login header"}
	protoMap, err := structpb.NewStruct(validMap)
	require.NoError(t, err)

	hash := md5.Sum(fmt.Append(nil, validMap))
	encodedHash := hex.EncodeToString(hash[:])

	tt := []struct {
		testName         string
		inputTranslation map[string]any
		expectedError    error
		expectedResponse *settings.GetHostedLoginTranslationResponse
	}{
		{
			testName:         "when unparsable map should return internal error",
			inputTranslation: map[string]any{"\xc5z": "something"},
			expectedError:    zerrors.ThrowInternal(protoimpl.X.NewError("invalid UTF-8 in string: %q", "\xc5z"), "QUERY-70ppPp", "Errors.Protobuf.ConvertToStruct"),
		},
		{
			testName:         "when input translation is valid should return expected response message",
			inputTranslation: validMap,
			expectedResponse: &settings.GetHostedLoginTranslationResponse{
				Translations: protoMap,
				Etag:         hex.EncodeToString(hash[:]),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := getTranslationOutputMessage(tc.inputTranslation, encodedHash)

			// Verify
			require.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResponse, res)
		})
	}
}

func TestGetHostedLoginTranslation(t *testing.T) {
	query := `SELECT projections.hosted_login_translations.file, projections.hosted_login_translations.aggregate_type, projections.hosted_login_translations.etag
	FROM projections.hosted_login_translations
	WHERE projections.hosted_login_translations.aggregate_id = $1
	AND projections.hosted_login_translations.aggregate_type = $2
	AND projections.hosted_login_translations.instance_id = $3
	AND (projections.hosted_login_translations.locale = $4 OR projections.hosted_login_translations.locale = $5)
	LIMIT 2`
	okTranslation := defaultLoginTranslations

	parsedOKTranslation := map[string]map[string]any{}
	require.NoError(t, json.Unmarshal(okTranslation, &parsedOKTranslation))

	protoDefaultTranslation, err := structpb.NewStruct(parsedOKTranslation["en"])
	require.Nil(t, err)

	defaultWithDBTranslations := maps.Clone(parsedOKTranslation["en"])
	defaultWithDBTranslations["test"] = "translation"
	defaultWithDBTranslations["test2"] = "translation2"
	protoDefaultWithDBTranslation, err := structpb.NewStruct(defaultWithDBTranslations)
	require.NoError(t, err)

	nilProtoDefaultMap, err := structpb.NewStruct(nil)
	require.NoError(t, err)

	hashDefaultTranslations := md5.Sum(fmt.Append(nil, parsedOKTranslation["en"]))

	tt := []struct {
		testName string

		defaultInstanceLanguage language.Tag
		sqlExpectations         []mock.Expectation

		inputRequest *settings.GetHostedLoginTranslationRequest

		expectedError  error
		expectedResult *settings.GetHostedLoginTranslationResponse
	}{
		{
			testName: "when input language is invalid should return invalid argument error",

			inputRequest: &settings.GetHostedLoginTranslationRequest{},

			expectedError: zerrors.ThrowInvalidArgument(nil, "QUERY-rZLAGi", "Errors.Arguments.Locale.Invalid"),
		},
		{
			testName: "when input language is root should return invalid argument error",

			defaultInstanceLanguage: language.English,
			inputRequest: &settings.GetHostedLoginTranslationRequest{
				Locale: "root",
			},

			expectedError: zerrors.ThrowInvalidArgument(nil, "QUERY-rZLAGi", "Errors.Arguments.Locale.Invalid"),
		},
		{
			testName: "when no system translation is available should return not found error",

			defaultInstanceLanguage: language.Romanian,
			inputRequest: &settings.GetHostedLoginTranslationRequest{
				Locale: "ro-RO",
			},

			expectedError: zerrors.ThrowNotFoundf(nil, "QUERY-6gb5QR", "Errors.Query.HostedLoginTranslationNotFound-%s", "ro"),
		},
		{
			testName: "when requesting system translation should return it",

			defaultInstanceLanguage: language.English,
			inputRequest: &settings.GetHostedLoginTranslationRequest{
				Locale: "en-US",
				Level:  &settings.GetHostedLoginTranslationRequest_System{},
			},

			expectedResult: &settings.GetHostedLoginTranslationResponse{
				Translations: protoDefaultTranslation,
				Etag:         hex.EncodeToString(hashDefaultTranslations[:]),
			},
		},
		{
			testName: "when querying DB fails should return internal error",

			defaultInstanceLanguage: language.English,
			sqlExpectations: []mock.Expectation{
				mock.ExpectQuery(
					query,
					mock.WithQueryArgs("123", "org", "instance-id", "en-US", "en"),
					mock.WithQueryErr(sql.ErrConnDone),
				),
			},
			inputRequest: &settings.GetHostedLoginTranslationRequest{
				Locale: "en-US",
				Level: &settings.GetHostedLoginTranslationRequest_OrganizationId{
					OrganizationId: "123",
				},
			},

			expectedError: zerrors.ThrowInternal(sql.ErrConnDone, "QUERY-6k1zjx", "Errors.Internal"),
		},
		{
			testName: "when querying DB returns no result should return system translations",

			defaultInstanceLanguage: language.English,
			sqlExpectations: []mock.Expectation{
				mock.ExpectQuery(
					query,
					mock.WithQueryArgs("123", "org", "instance-id", "en-US", "en"),
					mock.WithQueryResult(
						[]string{"file", "aggregate_type", "etag"},
						[][]driver.Value{},
					),
				),
			},
			inputRequest: &settings.GetHostedLoginTranslationRequest{
				Locale: "en-US",
				Level: &settings.GetHostedLoginTranslationRequest_OrganizationId{
					OrganizationId: "123",
				},
			},

			expectedResult: &settings.GetHostedLoginTranslationResponse{
				Translations: protoDefaultTranslation,
				Etag:         hex.EncodeToString(hashDefaultTranslations[:]),
			},
		},
		{
			testName: "when querying DB returns no result and inheritance disabled should return empty result",

			defaultInstanceLanguage: language.English,
			sqlExpectations: []mock.Expectation{
				mock.ExpectQuery(
					query,
					mock.WithQueryArgs("123", "org", "instance-id", "en-US", "en"),
					mock.WithQueryResult(
						[]string{"file", "aggregate_type", "etag"},
						[][]driver.Value{},
					),
				),
			},
			inputRequest: &settings.GetHostedLoginTranslationRequest{
				Locale: "en-US",
				Level: &settings.GetHostedLoginTranslationRequest_OrganizationId{
					OrganizationId: "123",
				},
				IgnoreInheritance: true,
			},

			expectedResult: &settings.GetHostedLoginTranslationResponse{
				Etag:         "",
				Translations: nilProtoDefaultMap,
			},
		},
		{
			testName: "when querying DB returns records should return merged result",

			defaultInstanceLanguage: language.English,
			sqlExpectations: []mock.Expectation{
				mock.ExpectQuery(
					query,
					mock.WithQueryArgs("123", "org", "instance-id", "en-US", "en"),
					mock.WithQueryResult(
						[]string{"file", "aggregate_type", "etag"},
						[][]driver.Value{
							{[]byte(`{"test": "translation"}`), "org", "etag-org"},
							{[]byte(`{"test2": "translation2"}`), "instance", "etag-instance"},
						},
					),
				),
			},
			inputRequest: &settings.GetHostedLoginTranslationRequest{
				Locale: "en-US",
				Level: &settings.GetHostedLoginTranslationRequest_OrganizationId{
					OrganizationId: "123",
				},
			},

			expectedResult: &settings.GetHostedLoginTranslationResponse{
				Etag:         "etag-org",
				Translations: protoDefaultWithDBTranslation,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Given
			db := &database.DB{DB: mock.NewSQLMock(t, tc.sqlExpectations...).DB}
			querier := Queries{client: db}

			ctx := authz.NewMockContext("instance-id", "org-id", "user-id", authz.WithMockDefaultLanguage(tc.defaultInstanceLanguage))

			// When
			res, err := querier.GetHostedLoginTranslation(ctx, tc.inputRequest)

			// Verify
			require.Equal(t, tc.expectedError, err)

			if tc.expectedError == nil {
				assert.Equal(t, tc.expectedResult.GetEtag(), res.GetEtag())
				assert.Equal(t, tc.expectedResult.GetTranslations().GetFields(), res.GetTranslations().GetFields())
			}
		})
	}
}
