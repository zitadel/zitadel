package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

// webhook HandleFunc to read the request body and then print out the contents
func webhook(signatureAlgorithms []jose.SignatureAlgorithm, keySet oidc.KeySet, privateKey any) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// read the body content
		sentBody, err := io.ReadAll(req.Body)
		if err != nil {
			// if there was an error while reading the body return an error
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
		defer req.Body.Close()
		// decrypt the JWE, validate the JWT and extract the payload from it
		payload, err := validateJWE(req.Context(), string(sentBody), signatureAlgorithms, keySet, privateKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// print out the payload
		fmt.Println(string(payload))
	}
}

func validateJWE(ctx context.Context, jweString string, signatureAlgorithms []jose.SignatureAlgorithm, keySet oidc.KeySet, privateKey any) ([]byte, error) {
	// For parsing the JWE, we need to specify which algorithms we expect.
	// You can specify either `RSA_OAEP_256` (rsa.PrivateKey) or `ECDH_ES_A256KW` (ecdsa.PrivateKey) or both, depending on which private key(s) you have.
	// The content encryption algorithm will always be `A256GCM` as used by Zitadel.
	parsedJWE, err := jose.ParseEncrypted(jweString, []jose.KeyAlgorithm{jose.RSA_OAEP_256, jose.ECDH_ES_A256KW}, []jose.ContentEncryption{jose.A256GCM})
	if err != nil {
		return nil, err
	}

	// In this example we only use a single key and loaded it on start, but you can manage and rotate the key used for the encryption through the API.
	// You might also have a key management system in your application.
	// Using the KeyID from the JWE header (parsedJWE.Header.KeyID), we could load the correct private key for decryption.

	// Decrypt the JWE using the private key to get the inner JWT
	decryptedJWT, err := parsedJWE.Decrypt(privateKey)
	if err != nil {
		return nil, err
	}
	// Now validate the inner JWT and return the payload
	return validateJWT(ctx, string(decryptedJWT), signatureAlgorithms, keySet)
}

func validateJWT(ctx context.Context, jwtString string, algorithms []jose.SignatureAlgorithm, keySet oidc.KeySet) ([]byte, error) {
	// Parse the signed JWT to get the JWS object (which contains the signatures and unverified payload)
	parsedJWS, err := jose.ParseSigned(jwtString, algorithms)
	if err != nil {
		return nil, err
	}

	// Verify the signature using the retrieved key and return the payload
	return keySet.VerifySignature(ctx, parsedJWS)
}

func main() {
	ctx := context.Background()
	// set the issuer to the ZITADEL instance URL
	issuer := "http://localhost:8080"

	// the oidc client library will call the discovery endpoint to get the JWKS URI and supported signing algorithms
	discover, err := client.Discover(ctx, issuer, http.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}
	signatureAlgorithms := make([]jose.SignatureAlgorithm, len(discover.IDTokenSigningAlgValuesSupported))
	for i, alg := range discover.IDTokenSigningAlgValuesSupported {
		signatureAlgorithms[i] = jose.SignatureAlgorithm(alg)
	}
	keySet := rp.NewRemoteKeySet(http.DefaultClient, discover.JwksURI)

	// this example uses the private key from the PEM file to decrypt the JWE
	// make sure to load your private key accordingly
	privateKeyPEM, err := os.ReadFile("./private-key.pem")
	if err != nil {
		log.Fatal(err)
	}
	block, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	// handle the HTTP call under "/webhook"
	http.HandleFunc("/webhook", webhook(signatureAlgorithms, keySet, privateKey))

	// start an HTTP server with the before defined function to handle the endpoint under "http://localhost:8090"
	http.ListenAndServe(":8090", nil)
}
