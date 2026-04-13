package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"

	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestCheckPasskeyGRPCToDomain(t *testing.T) {
	t.Parallel()
	tt := []struct {
		testName      string
		input         *session_grpc.CheckWebAuthN
		expectedBytes []byte
		expectedError error
	}{
		{
			testName:      "when input is nil should return nil bytes and nil error",
			input:         nil,
			expectedBytes: nil,
			expectedError: nil,
		},
		{
			testName: "when credential assertion data is nil should return nil bytes and nil error",
			input: &session_grpc.CheckWebAuthN{
				CredentialAssertionData: nil,
			},
			expectedBytes: nil,
			expectedError: nil,
		},
		{
			testName: "when credential assertion data is empty should return empty JSON bytes",
			input: &session_grpc.CheckWebAuthN{
				CredentialAssertionData: &structpb.Struct{},
			},
			expectedBytes: []byte("{}"),
			expectedError: nil,
		},
		{
			testName: "when credential assertion data is populated should return marshaled JSON bytes",
			input: &session_grpc.CheckWebAuthN{
				CredentialAssertionData: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"publicKeyCredentialRequestOptions": {
							Kind: &structpb.Value_StructValue{
								StructValue: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"challenge": structpb.NewStringValue("Y2hhbGxlbmdl"),
										"rpId":      structpb.NewStringValue("example.com"),
										"timeout":   structpb.NewNumberValue(5000),
									},
								},
							},
						},
					},
				},
			},
			expectedBytes: []byte(`{"publicKeyCredentialRequestOptions":{"challenge":"Y2hhbGxlbmdl","rpId":"example.com","timeout":5000}}`),
			expectedError: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Test
			result, err := CheckPasskeyGRPCToDomain(tc.input)

			// Verify
			assert.Equal(t, tc.expectedBytes, result)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
