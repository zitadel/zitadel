package command

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net/url"
	"strings"
	"testing"

	"github.com/zitadel/zitadel/internal/api/http"
	zitadel_crypto "github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
)

// mockFetcherForOrdering is a minimal mock for FederatedLogoutDataFetcher
type mockFetcherForOrdering struct{}

func (m *mockFetcherForOrdering) IDPUserLinks(ctx context.Context, queries *query.IDPUserLinksSearchQuery, permissionCheck domain.PermissionCheck) (*query.IDPUserLinks, error) {
	// Return a dummy link
	return &query.IDPUserLinks{
		Links: []*query.IDPUserLink{
			{
				IDPID:            "mock-idp-id",
				ProvidedUserID:   "mock-provider-user-id",
				UserID:           "mock-user-id",
				ProvidedUsername: "mock-user",
			},
		},
	}, nil
}

func (m *mockFetcherForOrdering) IDPTemplateByID(ctx context.Context, shouldTriggerBulk bool, id string, withOwnerRemoved bool, permissionCheck domain.PermissionCheck, queries ...query.SearchQuery) (*query.IDPTemplate, error) {
	// Generate a temporary key pair for signing
	key, cert, err := generateTestKeyAndCert()
	if err != nil {
		return nil, err
	}

	// Return a dummy IDP template with SAML configured
	return &query.IDPTemplate{
		ID:   "mock-idp-id",
		Name: "Mock IDP",
		Type: domain.IDPTypeSAML,
		SAMLIDPTemplate: &query.SAMLIDPTemplate{
			IDPID:                  "mock-idp-id",
			Metadata:               []byte(mockMetadata),
			Key:                    &zitadel_crypto.CryptoValue{Crypted: key}, // Store raw key in Crypted, no-op decrypt will return it
			Certificate:            cert,
			Binding:                "redirect",
			WithSignedRequest:      true,
			SignatureAlgorithm:     "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
			FederatedLogoutEnabled: true,
		},
	}, nil
}

type mockEventStoreForOrdering struct{}

func (m *mockEventStoreForOrdering) Filter(ctx context.Context, searchQuery *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
	return []eventstore.Event{}, nil
}
func (m *mockEventStoreForOrdering) Push(ctx context.Context, commands ...eventstore.Command) ([]eventstore.Event, error) {
	return nil, nil
}
func (m *mockEventStoreForOrdering) FilterToQueryReducer(ctx context.Context, reducer eventstore.QueryReducer) error {
	return nil
}

// TestStartFederatedLogout_ParameterOrdering verifies that the query parameters in the Redirect URL
// are strictly ordered as required by the SAML specification: SAMLRequest, RelayState, SigAlg.
func TestStartFederatedLogout_ParameterOrdering(t *testing.T) {
	// Setup
	// Commands struct is available in the package
	c := &Commands{
		idpConfigEncryption: &noOpEncryption{},
	}

	ctx := context.Background()
	ctx = http.WithDomainContext(ctx, http.NewDomainCtx("test.zitadel.ch", "test.zitadel.ch", "https"))

	// 1. Prepare inputs manually matching what IDPTemplateByID returns
	keyPEM, certPEM, err := generateTestKeyAndCert()
	if err != nil {
		t.Fatalf("failed to generate certs: %v", err)
	}

	idpTemplate := &query.IDPTemplate{
		ID:   "mock-idp-id",
		Name: "Mock IDP",
		SAMLIDPTemplate: &query.SAMLIDPTemplate{
			IDPID:                  "mock-idp-id",
			Metadata:               []byte(mockMetadata),
			Key:                    &zitadel_crypto.CryptoValue{KeyID: "no-op-key", Algorithm: "no-op", Crypted: keyPEM},
			Certificate:            certPEM,
			Binding:                "redirect",
			WithSignedRequest:      true,
			SignatureAlgorithm:     "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
			FederatedLogoutEnabled: true,
		},
	}

	// 2. Call the private method directly
	reqData, err := c.generateSAMLLogoutRequest(
		ctx,
		&mockEventStoreForOrdering{},
		idpTemplate,
		"user-id",
		"name-id",
		"relay-state-value",
		"instance-id",
	)

	if err != nil {
		t.Fatalf("generateSAMLLogoutRequest failed: %v", err)
	}

	// 3. Verify the URL parameters
	targetURL := reqData.RedirectURL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		t.Fatalf("failed to parse redirect URL: %v", err)
	}

	query := parsedURL.RawQuery

	// 4. Assert Strict Order
	// The query string MUST contain SAMLRequest, RelayState, SigAlg, Signature IN THAT ORDER

	t.Logf("Generated Query: %s", query)

	// Check SAMLRequest is present
	samlReqIdx := strings.Index(query, "SAMLRequest=")
	if samlReqIdx == -1 {
		t.Errorf("SAMLRequest parameter missing")
	}

	// Check RelayState is present
	relayStateIdx := strings.Index(query, "RelayState=")
	if relayStateIdx == -1 {
		t.Errorf("RelayState parameter missing")
	}

	// Check SigAlg is present
	sigAlgIdx := strings.Index(query, "SigAlg=")
	if sigAlgIdx == -1 {
		t.Errorf("SigAlg parameter missing")
	}

	// Check Signature is present
	signatureIdx := strings.Index(query, "Signature=")
	if signatureIdx == -1 {
		t.Errorf("Signature parameter missing")
	}

	// VERIFY ORDER
	if !(samlReqIdx < relayStateIdx && relayStateIdx < sigAlgIdx && sigAlgIdx < signatureIdx) {
		t.Errorf("Query parameters are NOT in the correct order!\nQuery: %s\nIndices: SAMLRequest=%d, RelayState=%d, SigAlg=%d, Signature=%d",
			query, samlReqIdx, relayStateIdx, sigAlgIdx, signatureIdx)
	} else {
		t.Logf("Success: Query parameters are strictly ordered.")
	}

	// 5. Verify Signature using public key
	// Reconstruct the signed string data
	// The part BEFORE &Signature=...
	signedData := query[:signatureIdx-1] // -1 for the '&'

	// Decode signature
	signatureVal := query[signatureIdx+len("Signature="):]
	decodedVal, err := url.QueryUnescape(signatureVal)
	if err != nil {
		t.Fatalf("failed to unescape signature: %v", err)
	}
	decodedSig, err := base64.StdEncoding.DecodeString(decodedVal)
	if err != nil {
		t.Fatalf("failed to decode signature: %v", err)
	}

	// Parse certificate to get public key
	block, _ := pem.Decode(certPEM)
	if block == nil {
		t.Fatal("failed to decode cert PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("failed to parse certificate: %v", err)
	}

	rsaPub, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		t.Fatal("not an RSA public key")
	}

	// Verify
	// SAML usually uses SHA256 by default in our config: http://www.w3.org/2001/04/xmldsig-more#rsa-sha256
	// Which maps to crypto.SHA256
	hashed := sha256.Sum256([]byte(signedData))
	err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hashed[:], decodedSig)
	if err != nil {
		t.Errorf("Signature verification failed: %v", err)
	} else {
		t.Logf("Success: Signature verified against public key.")
	}
}

const mockMetadata = `<?xml version="1.0" encoding="UTF-8"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" entityID="https://mock-idp.com">
    <md:IDPSSODescriptor WantAuthnRequestsSigned="true" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
        <md:SingleLogoutService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://mock-idp.com/slo"/>
        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>
    </md:IDPSSODescriptor>
</md:EntityDescriptor>`
