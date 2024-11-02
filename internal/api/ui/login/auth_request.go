package login

import (
	"net/http"

	"github.com/zitadel/logging"

	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	return l.authRepo.AuthRequestByID(r.Context(), authRequestID, userAgentID)
}

func (l *Login) ensureAuthRequest(r *http.Request) (*domain.AuthRequest, error) {
	authRequest, err := l.getAuthRequest(r)
	if authRequest != nil || err != nil {
		return authRequest, err
	}
	return nil, zerrors.ThrowInvalidArgument(nil, "LOGIN-OLah9", "invalid or missing auth request")
}

func (l *Login) getAuthRequestAndParseData(r *http.Request, data interface{}) (*domain.AuthRequest, error) {
	authReq, err := l.getAuthRequest(r)
	if err != nil {
		return authReq, err
	}
	err = l.parser.Parse(r, data)
	return authReq, err
}

func (l *Login) ensureAuthRequestAndParseData(r *http.Request, data interface{}) (*domain.AuthRequest, error) {
	authReq, err := l.ensureAuthRequest(r)
	if err != nil {
		return authReq, err
	}
	err = l.parser.Parse(r, data)
	return authReq, err
}

func (l *Login) getParseData(r *http.Request, data interface{}) error {
	return l.parser.Parse(r, data)
}

// checkOptionalAuthRequestOfEmailLinks tries to get the [domain.AuthRequest] from the request.
// In case any error occurs, e.g. if the user agent does not correspond, the `authRequestID` query parameter will be
// removed from the request URL and form to ensure subsequent functions and pages do not use it.
// This function is used for handling links in emails, which could possibly be opened on another device than the
// auth request was initiated.
func (l *Login) checkOptionalAuthRequestOfEmailLinks(r *http.Request) *domain.AuthRequest {
	authReq, err := l.getAuthRequest(r)
	if err == nil {
		return authReq
	}
	logging.WithError(err).Infof("authrequest could not be found for email link on path %s", r.URL.RequestURI())
	queries := r.URL.Query()
	queries.Del(QueryAuthRequestID)
	r.URL.RawQuery = queries.Encode()
	r.RequestURI = r.URL.RequestURI()
	r.Form.Del(QueryAuthRequestID)
	r.PostForm.Del(QueryAuthRequestID)
	return nil
}
