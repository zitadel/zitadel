//go:build ignore

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

// webhook HandleFunc to read the request body and then print out the contents
func webhook(signatureAlgorithms []jose.SignatureAlgorithm, keySet oidc.KeySet) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// read the body content
		sentBody, err := io.ReadAll(req.Body)
		if err != nil {
			// if there was an error while reading the body return an error
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
		defer req.Body.Close()
		// validate the JWT and extract the payload from it
		payload, err := validateJWT(req.Context(), string(sentBody), signatureAlgorithms, keySet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// print out the payload
		fmt.Println(string(payload))
	}
}

func validateJWT(ctx context.Context, jwtString string, algorithms []jose.SignatureAlgorithm, keySet oidc.KeySet) ([]byte, error) { //nolint:typecheck
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
	// set the issuer to the Zitadel instance URL
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

	// handle the HTTP call under "/webhook"
	http.HandleFunc("/webhook", webhook(signatureAlgorithms, keySet))

	// start an HTTP server with the before defined function to handle the endpoint under "http://localhost:8090"
	http.ListenAndServe(":8090", nil)
}
