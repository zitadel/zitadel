package actions

import (
	"encoding/json"

	"github.com/caos/oidc/pkg/oidc"
)

type Context map[string]interface{}

func (c Context) set(name string, value interface{}) {
	map[string]interface{}(c)[name] = value
}

func (c *Context) SetToken(t *oidc.Tokens) *Context {
	if t == nil {
		return c
	}
	if t.Token != nil && t.Token.AccessToken != "" {
		c.set("accessToken", t.AccessToken)
	}
	if t.IDToken != "" {
		c.set("idToken", t.IDToken)
	}
	if t.IDTokenClaims != nil {
		c.set("getClaim", func(claim string) interface{} { return t.IDTokenClaims.GetClaim(claim) })
		c.set("claimsJSON", func() (string, error) {
			c, err := json.Marshal(t.IDTokenClaims)
			if err != nil {
				return "", err
			}
			return string(c), nil
		})
	}
	return c
}
