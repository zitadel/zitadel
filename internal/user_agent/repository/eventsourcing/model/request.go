package model

import "github.com/caos/zitadel/internal/user_agent/model"

type Request interface {
	Type() int32
}

type AuthSessionOIDC struct {
	Scopes        []string
	ResponseType  int32
	Nonce         string
	CodeChallenge *OIDCCodeChallenge
}

func (a *AuthSessionOIDC) Type() int32 {
	return 0
}

func RequestFromModel(request model.Request) Request {
	switch req := request.(type) {
	case *model.AuthSessionOIDC:
		return OIDCRequestFromModel(req)
	}
	return nil
}
func RequestToModel(request Request) model.Request {
	switch req := request.(type) {
	case *AuthSessionOIDC:
		return OIDCRequestToModel(req)
	}
	return nil
}

func OIDCRequestFromModel(request *model.AuthSessionOIDC) *AuthSessionOIDC {
	return &AuthSessionOIDC{
		Scopes:        request.Scopes,
		ResponseType:  int32(request.ResponseType),
		Nonce:         request.Nonce,
		CodeChallenge: OIDCCodeChallengeFromModel(request.CodeChallenge),
	}
}
func OIDCRequestToModel(request *AuthSessionOIDC) *model.AuthSessionOIDC {
	return &model.AuthSessionOIDC{
		Scopes:        request.Scopes,
		ResponseType:  model.OIDCResponseType(request.ResponseType),
		Nonce:         request.Nonce,
		CodeChallenge: OIDCCodeChallengeToModel(request.CodeChallenge),
	}
}
