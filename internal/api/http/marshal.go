package http

import (
	"encoding/json"
	"net/http"

	"github.com/caos/logging"
)

func MarshalJSON(w http.ResponseWriter, i interface{}) {
	b, err := json.Marshal(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(b)
	logging.Log("HTTP-sdgT2").OnError(err).Error("error writing response")
}
