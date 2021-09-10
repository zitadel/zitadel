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
	c.set("accessToken", t.AccessToken)
	c.set("idToken", t.IDToken)
	c.set("getClaim", func(claim string) interface{} { return t.IDTokenClaims.GetClaim(claim) })
	c.set("claimsJSON", func() (string, error) {
		c, err := json.Marshal(t.IDTokenClaims)
		if err != nil {
			return "", err
		}
		return string(c), nil
	})
	return c
}
