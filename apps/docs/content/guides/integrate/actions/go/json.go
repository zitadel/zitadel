//go:build ignore

package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/zitadel/zitadel-go/v3/pkg/actions"
)

const signingKey = "somekey" // signing key received after creating the target

// webhook HandleFunc to read the request body and then print out the contents
func webhook(w http.ResponseWriter, req *http.Request) {
	// read the body content
	sentBody, err := io.ReadAll(req.Body)
	if err != nil {
		// if there was an error while reading the body return an error
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()
	// validate signature
	if err := actions.ValidateRequestPayload(sentBody, &req.Header, signingKey); err != nil {
		// if the signed content is not equal the sent content return an error
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	// print out the read content
	fmt.Println(string(sentBody))
}

func main() {
	// handle the HTTP call under "/webhook"
	http.HandleFunc("/webhook", webhook)

	// start an HTTP server with the before defined function to handle the endpoint under "http://localhost:8090"
	http.ListenAndServe(":8090", nil)
}
