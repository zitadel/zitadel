package object

import (
	"github.com/dop251/goja"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/domain"
)

// TokenActorField accepts a domain.TokenActor pointer and copies its content so scripts can't mutate the domain object.
func TokenActorField(actor *domain.TokenActor) func(c *actions.FieldConfig) interface{} {
	return func(c *actions.FieldConfig) interface{} {
		return TokenActorFromDomain(c, actor)
	}
}

// TokenActorFromDomain returns the actor of a token, in case it was obtained through
// token exchange or impersonation. It returns null when there is no actor.
func TokenActorFromDomain(c *actions.FieldConfig, actor *domain.TokenActor) goja.Value {
	if actor == nil {
		return c.Runtime.ToValue(nil)
	}
	return c.Runtime.ToValue(tokenActorFromDomain(actor))
}

// tokenActorFromDomain copies the delegation chain, so scripts can't mutate the domain object.
func tokenActorFromDomain(actor *domain.TokenActor) *tokenActor {
	if actor == nil {
		return nil
	}
	return &tokenActor{
		Actor:  tokenActorFromDomain(actor.Actor),
		UserId: actor.UserID,
		Issuer: actor.Issuer,
	}
}

type tokenActor struct {
	// Actor is the previous actor in the delegation chain, null if there is none.
	Actor  *tokenActor
	UserId string
	Issuer string
}
