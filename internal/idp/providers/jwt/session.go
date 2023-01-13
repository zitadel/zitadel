package jwt

import (
	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Session = (*Session)(nil)

type Session struct {
	AuthURL string
	Tokens  *oidc.Tokens
}

func (s *Session) GetAuthURL() string {
	return s.AuthURL
}
