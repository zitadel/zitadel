package handler

import (
	"net/http"
)

func (l *Login) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (l *Login) handleReadiness(w http.ResponseWriter, r *http.Request) {
	errs := l.service.Hellth(r.Context())
	for _, err := range errs {
		if err != nil {
			http.Error(w, "not ready", http.StatusInternalServerError)
			return
		}
	}
	w.Write([]byte("OK"))
}
