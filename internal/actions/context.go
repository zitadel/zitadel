package actions

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dop251/goja"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"github.com/zitadel/zitadel/internal/query"
)

type contextParam struct {
	runtime *goja.Runtime
	parameter
}

type ContextOption func(*contextParam)

func SetToken(t *oidc.Tokens) ContextOption {
	return func(c *contextParam) {
		if t == nil {
			return
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
	}
}

func SetUserID(id string) ContextOption {
	return func(c *contextParam) {
		c.setPath([]string{"v1", "user", "id"}, id)
	}
}

type userMetadataGetter interface {
	SearchUserMetadata(ctx context.Context, shouldTriggerBulk bool, userID string, queries *query.UserMetadataSearchQueries) (*query.UserMetadataList, error)
}

func SetUserMetadataGetter(ctx context.Context, getter userMetadataGetter, userID, resourceOwner string) ContextOption {
	return func(c *contextParam) {
		c.setPath([]string{"v1", "user", "setMetadata"}, func(call goja.FunctionCall) goja.Value {
			resourceOwnerQuery, err := query.NewUserMetadataResourceOwnerSearchQuery(resourceOwner)
			if err != nil {
				logging.WithError(err).Debug("unable to create search query")
				panic(err)
			}
			metadata, err := getter.SearchUserMetadata(
				ctx,
				true,
				userID,
				&query.UserMetadataSearchQueries{Queries: []query.SearchQuery{resourceOwnerQuery}},
			)
			if err != nil {
				logging.WithError(err).Info("unable to get md in action")
				panic(err)
			}
			return c.runtime.ToValue(c.userMetadataListFromQuery(metadata))
		})
	}
}

func (c *contextParam) userMetadataListFromQuery(metadata *query.UserMetadataList) *userMetadataList {
	result := &userMetadataList{
		Count:     metadata.Count,
		Sequence:  metadata.Sequence,
		Timestamp: metadata.Timestamp,
		Metadata:  make([]*userMetadata, len(metadata.Metadata)),
	}

	for i, md := range metadata.Metadata {
		var value interface{}
		err := json.Unmarshal(md.Value, &value)
		if err != nil {
			logging.WithError(err).Debug("unable to unmarshal into map")
			panic(err)
		}
		result.Metadata[i] = &userMetadata{
			CreationDate:  md.CreationDate,
			ChangeDate:    md.ChangeDate,
			ResourceOwner: md.ResourceOwner,
			Sequence:      md.Sequence,
			Key:           md.Key,
			Value:         c.runtime.ToValue(value),
		}
	}

	return result
}

type userMetadataList struct {
	Count     uint64
	Sequence  uint64
	Timestamp time.Time
	Metadata  []*userMetadata
}

type userMetadata struct {
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64
	Key           string
	Value         goja.Value
}
