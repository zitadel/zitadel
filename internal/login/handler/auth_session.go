package handler

import (
	"github.com/caos/zitadel/internal/auth_request/model"
	"net/http"
)

const (
	queryAuthRequestID = "authRequestID"
)

func (l *Login) getAuthRequest(r *http.Request) (*model.AuthRequest, error) {
	authRequestID := r.FormValue(queryAuthRequestID)
	if authRequestID == "" {
		return nil, nil
	}
	return l.authRepo.AuthRequestByID(r.Context(), authRequestID)
}

func (l *Login) getAuthRequestAndParseData(r *http.Request, data interface{}) (*model.AuthRequest, error) {
	authSession, err := l.getAuthRequest(r)
	if err != nil {
		return nil, err
	}
	err = l.parser.Parse(r, data)
	return authSession, err
}
