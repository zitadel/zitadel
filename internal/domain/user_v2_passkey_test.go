package domain

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

func TestRenderPasskeyURLTemplate(t *testing.T) {
	type args struct {
		tmpl          string
		userID        string
		resourceOwner string
		codeID        string
		code          string
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr error
	}{
		{
			name: "parse error",
			args: args{
				tmpl: "{{",
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "DOMAIN-oGh5e", "Errors.User.InvalidURLTemplate"),
		},
		{
			name: "success",
			args: args{
				tmpl:          "https://example.com/passkey/register?userID={{.UserID}}&orgID={{.ResourceOwner}}&codeID={{.CodeID}}&code={{.Code}}",
				userID:        "user1",
				resourceOwner: "org1",
				codeID:        "99",
				code:          "123",
			},
			wantW: "https://example.com/passkey/register?userID=user1&orgID=org1&codeID=99&code=123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := RenderPasskeyURLTemplate(w, tt.args.tmpl, tt.args.userID, tt.args.resourceOwner, tt.args.codeID, tt.args.code)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantW, w.String())
		})
	}
}

func TestWebAuthNToken_PasskeyRegistrationDetails(t *testing.T) {
	webAuthN := &WebAuthNToken{
		WebAuthNTokenID:        "1",
		CredentialCreationData: []byte{1, 2, 3},
		State:                  MFAStateReady,
		Challenge:              "challenge",
		AllowedCredentialIDs:   [][]byte{{4, 5, 6}, {7, 8, 9}},
		UserVerification:       UserVerificationRequirementRequired,
		KeyID:                  []byte{10, 11, 12},
		PublicKey:              []byte{13, 14, 15},
		AttestationType:        "attestation_type",
		AAGUID:                 []byte{16, 17, 18},
		SignCount:              999,
		WebAuthNTokenName:      "awesome",
	}
	details := &ObjectDetails{
		Sequence:      77,
		EventDate:     time.Now(),
		ResourceOwner: "memememe",
	}
	want := &PasskeyRegistrationDetails{
		ObjectDetails:          details,
		WebAuthNTokenID:        "1",
		CredentialCreationData: []byte{1, 2, 3},
		State:                  MFAStateReady,
		Challenge:              "challenge",
		AllowedCredentialIDs:   [][]byte{{4, 5, 6}, {7, 8, 9}},
		UserVerification:       UserVerificationRequirementRequired,
		KeyID:                  []byte{10, 11, 12},
		PublicKey:              []byte{13, 14, 15},
		AttestationType:        "attestation_type",
		AAGUID:                 []byte{16, 17, 18},
		SignCount:              999,
		WebAuthNTokenName:      "awesome",
	}
	got := webAuthN.PasskeyRegistrationDetails(details)
	assert.Equal(t, want, got)
}
