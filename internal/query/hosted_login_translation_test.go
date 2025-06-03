package query

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/structpb"
)

func jsonSyntaxErrorGenerator(notAMap []byte) error {
	d := map[string]any{}

	return json.Unmarshal(notAMap, &d)
}

func TestGetSystemTranslation(t *testing.T) {
	okTranslation := defaultLoginTranslations

	parsedOKTranslation := map[string]map[string]any{}
	require.Nil(t, json.Unmarshal(okTranslation, &parsedOKTranslation))

	malformedTranslation := []byte{1, 2}

	tt := []struct {
		testName string

		inputLanguage          string
		inputInstanceLanguage  string
		systemTranslationToSet []byte

		expectedLanguage map[string]any
		expectedError    error
	}{
		{
			testName:               "when unmarshalling default translation fails should return internal error",
			systemTranslationToSet: malformedTranslation,

			expectedError: zerrors.ThrowInternal(jsonSyntaxErrorGenerator(malformedTranslation), "QUERY-nvx88W", "Errors.Query.UnmarshalDefaultLoginTranslations"),
		},
		{
			testName:               "when neither input language nor system default language have translation should return not found error",
			systemTranslationToSet: okTranslation,
			inputLanguage:          "ro",
			inputInstanceLanguage:  "fr",

			expectedError: zerrors.ThrowNotFoundf(nil, "QUERY-6gb5QR", "Errors.Query.HostedLoginTranslationNotFound-%s", "ro"),
		},
		{
			testName:               "when input language has no translation should fallback onto instance default",
			systemTranslationToSet: okTranslation,
			inputLanguage:          "ro",
			inputInstanceLanguage:  "de",

			expectedLanguage: parsedOKTranslation["de"],
		},
		{
			testName:               "when input language has translation should return it",
			systemTranslationToSet: okTranslation,
			inputLanguage:          "de",
			inputInstanceLanguage:  "en",

			expectedLanguage: parsedOKTranslation["de"],
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// Given
			defaultLoginTranslations = tc.systemTranslationToSet

			// When
			translation, err := getSystemTranslation(tc.inputLanguage, tc.inputInstanceLanguage)

			// Verify
			require.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedLanguage, translation)
		})
	}
}

func TestGetTranslationOutput(t *testing.T) {
	t.Parallel()

	validMap := map[string]any{"loginHeader": "A login header"}
	protoMap, err := structpb.NewStruct(validMap)
	require.Nil(t, err)

	hash := md5.Sum([]byte(protoMap.String()))

	tt := []struct {
		testName         string
		inputTranslation map[string]any
		expectedError    error
		expectedResponse *settings.GetHostedLoginTranslationResponse
	}{
		{
			testName:         "when unparseable map should return internal error",
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
			res, err := getTranslationOutputMessage(tc.inputTranslation)

			// Verify
			require.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResponse, res)
		})
	}
}
