package handler

import (
	"github.com/caos/zitadel/internal/errors"
	"net/http"
	"strings"
)

const (
	queryAuthSessionID = "authSessionID"
)

func (l *Login) getAuthSession(r *http.Request) (*model.AuthSession, error) {
	authSessionID := r.FormValue(queryAuthSessionID)
	if authSessionID == "" {
		return nil, nil
	}
	userAgent, err := l.userAgentHandler.GetUserAgent(r)
	if err != nil {
		return nil, err
	}
	ids := strings.Split(authSessionID, ":")
	if len(ids) != 2 {
		return nil, errors.ThrowInvalidArgument(nil, "APP-QfSSPm", "invalid id")
	}
	if ids[0] != userAgent.GetID() {
		return nil, errors.ThrowInvalidArgument(nil, "APP-x0UPKz", "invalid id")
	}
	return l.service.Auth.GetAuthSession(r.Context(), ids[1], userAgent.GetID(), &model.BrowserInformation{RemoteIP: &model.IP{}})
}

func (l *Login) getAuthSessionAndParseData(r *http.Request, data interface{}) (*model.AuthSession, error) {
	authSession, err := l.getAuthSession(r)
	if err != nil {
		return nil, err
	}
	err = l.parser.Parse(r, data)
	return authSession, err
}
