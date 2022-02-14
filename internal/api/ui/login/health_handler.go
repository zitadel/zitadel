package login

import (
	"net/http"
)

func (l *Login) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (l *Login) handleReadiness(w http.ResponseWriter, r *http.Request) {
	err := l.authRepo.Health(r.Context())
	if err != nil {
		http.Error(w, "not ready", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("OK"))
}
