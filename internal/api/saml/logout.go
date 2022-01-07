package saml

import (
	"fmt"
	"net/http"
)

func(p * IdentityProvider) logoutHandleFunc(w http.ResponseWriter, r *http.Request) {
	//TODO
	http.Error(w, fmt.Sprintf("not implemented yet"), http.StatusNotImplemented)
}
