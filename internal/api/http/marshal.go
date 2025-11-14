package http

import (
	"encoding/json"
	"net/http"

	"github.com/zitadel/logging"
)

func MarshalJSON(w http.ResponseWriter, i interface{}, err error, statusCode int) {
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}
	b, err := json.Marshal(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(b)
	logging.WithFields("logID", "HTTP-sdgT2").OnError(err).Error("error writing response")
}
