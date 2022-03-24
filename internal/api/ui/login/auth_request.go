package login

import (
	"net/http"

	"github.com/caos/zitadel/internal/domain"

	"github.com/caos/zitadel/internal/api/authz"
	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
)

const (
	QueryAuthRequestID = "authRequestID"
	queryUserAgentID   = "userAgentID"
)

func (l *Login) getAuthRequest(r *http.Request) (*domain.AuthRequest, error) {
	authRequestID := r.FormValue(QueryAuthRequestID)
	if authRequestID == "" {
		return nil, nil
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	instanceID := authz.GetInstance(r.Context()).ID
	return l.authRepo.AuthRequestByID(r.Context(), authRequestID, userAgentID, instanceID)
}

func (l *Login) getAuthRequestAndParseData(r *http.Request, data interface{}) (*domain.AuthRequest, error) {
	authReq, err := l.getAuthRequest(r)
	if err != nil {
		return authReq, err
	}
	err = l.parser.Parse(r, data)
	return authReq, err
}

func (l *Login) getParseData(r *http.Request, data interface{}) error {
	return l.parser.Parse(r, data)
}
