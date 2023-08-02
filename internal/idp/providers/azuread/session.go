package azuread

import (
	"errors"

	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
)

var ErrIDTokenMissing = errors.New("no id_token provided")

type Session struct {
	*oauth.Session
}

func (s *Session) RetrieveOldID() (string, error) {
	idToken, ok := s.Tokens.Token.Extra("id_token").(string)
	if !ok {
		return "", ErrIDTokenMissing
	}
	claims := new(oidc.IDTokenClaims)
	_, err := oidc.ParseToken(idToken, claims)
	if err != nil {
		return "", err
	}
	return claims.Subject, nil
}
