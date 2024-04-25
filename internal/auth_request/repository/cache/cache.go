package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/golang-lru/v2"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AuthRequestCache struct {
	client    *database.DB
	idCache   *lru.Cache[string, *domain.AuthRequest]
	codeCache *lru.Cache[string, *domain.AuthRequest]
}

func Start(dbClient *database.DB, amountOfCachedAuthRequests uint16) *AuthRequestCache {
	cache := &AuthRequestCache{
		client: dbClient,
	}
	idCache, err := lru.New[string, *domain.AuthRequest](int(amountOfCachedAuthRequests))
	logging.OnError(err).Info("auth request cache disabled")
	if err == nil {
		cache.idCache = idCache
	}
	codeCache, err := lru.New[string, *domain.AuthRequest](int(amountOfCachedAuthRequests))
	logging.OnError(err).Info("auth request cache disabled")
	if err == nil {
		cache.codeCache = codeCache
	}
	return cache
}

func (c *AuthRequestCache) Health(ctx context.Context) error {
	return c.client.PingContext(ctx)
}

func (c *AuthRequestCache) GetAuthRequestByID(ctx context.Context, id string) (*domain.AuthRequest, error) {
	if authRequest, ok := c.getCachedByID(ctx, id); ok {
		return authRequest, nil
	}
	request, err := c.getAuthRequest(ctx, "id", id, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	c.CacheAuthRequest(ctx, request)
	return request, nil
}

func (c *AuthRequestCache) GetAuthRequestByCode(ctx context.Context, code string) (*domain.AuthRequest, error) {
	if authRequest, ok := c.getCachedByCode(ctx, code); ok {
		return authRequest, nil
	}
	request, err := c.getAuthRequest(ctx, "code", code, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	c.CacheAuthRequest(ctx, request)
	return request, nil
}

func (c *AuthRequestCache) SaveAuthRequest(ctx context.Context, request *domain.AuthRequest) error {
	return c.saveAuthRequest(ctx, request, "INSERT INTO auth.auth_requests (id, request, instance_id, creation_date, change_date, request_type) VALUES($1, $2, $3, $4, $4, $5)", request.CreationDate, request.Request.Type())
}

func (c *AuthRequestCache) UpdateAuthRequest(ctx context.Context, request *domain.AuthRequest) error {
	if request.ChangeDate.IsZero() {
		request.ChangeDate = time.Now()
	}
	return c.saveAuthRequest(ctx, request, "UPDATE auth.auth_requests SET request = $2, instance_id = $3, change_date = $4, code = $5 WHERE id = $1", request.ChangeDate, request.Code)
}

func (c *AuthRequestCache) DeleteAuthRequest(ctx context.Context, id string) error {
	_, err := c.client.Exec("DELETE FROM auth.auth_requests WHERE instance_id = $1 and id = $2", authz.GetInstance(ctx).InstanceID(), id)
	if err != nil {
		return zerrors.ThrowInternal(err, "CACHE-dsHw3", "unable to delete auth request")
	}
	c.deleteFromCache(ctx, id)
	return nil
}

func (c *AuthRequestCache) getAuthRequest(ctx context.Context, key, value, instanceID string) (*domain.AuthRequest, error) {
	var b []byte
	var requestType domain.AuthRequestType
	query := fmt.Sprintf("SELECT request, request_type FROM auth.auth_requests WHERE instance_id = $1 and %s = $2", key)
	err := c.client.QueryRowContext(
		ctx,
		func(row *sql.Row) error {
			return row.Scan(&b, &requestType)
		},
		query, instanceID, value)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, zerrors.ThrowNotFound(err, "CACHE-d24aD", "Errors.AuthRequest.NotFound")
		}
		return nil, zerrors.ThrowInternal(err, "CACHE-as3kj", "Errors.Internal")
	}
	request, err := domain.NewAuthRequestFromType(requestType)
	if err == nil {
		err = json.Unmarshal(b, request)
	}
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "CACHE-2wshg", "Errors.Internal")
	}
	return request, nil
}

func (c *AuthRequestCache) saveAuthRequest(ctx context.Context, request *domain.AuthRequest, query string, date time.Time, param interface{}) error {
	b, err := json.Marshal(request)
	if err != nil {
		return zerrors.ThrowInternal(err, "CACHE-os0GH", "Errors.Internal")
	}
	_, err = c.client.Exec(query, request.ID, b, request.InstanceID, date, param)
	if err != nil {
		return zerrors.ThrowInternal(err, "CACHE-su3GK", "Errors.Internal")
	}
	c.CacheAuthRequest(ctx, request)
	return nil
}

func (c *AuthRequestCache) getCachedByID(ctx context.Context, id string) (*domain.AuthRequest, bool) {
	if c.idCache == nil {
		return nil, false
	}
	authRequest, ok := c.idCache.Get(cacheKey(ctx, id))
	logging.WithFields("hit", ok, "type", "id").Info("get from auth request cache")
	return authRequest, ok
}

func (c *AuthRequestCache) getCachedByCode(ctx context.Context, code string) (*domain.AuthRequest, bool) {
	if c.codeCache == nil {
		return nil, false
	}
	authRequest, ok := c.codeCache.Get(cacheKey(ctx, code))
	logging.WithFields("hit", ok, "type", "code").Info("get from auth request cache")
	return authRequest, ok
}

func (c *AuthRequestCache) CacheAuthRequest(ctx context.Context, request *domain.AuthRequest) {
	if c.idCache == nil {
		return
	}
	c.idCache.Add(cacheKey(ctx, request.ID), request)
	if request.Code != "" {
		c.codeCache.Add(cacheKey(ctx, request.Code), request)
	}
}

func cacheKey(ctx context.Context, value string) string {
	return fmt.Sprintf("%s-%s", authz.GetInstance(ctx).InstanceID(), value)
}

func (c *AuthRequestCache) deleteFromCache(ctx context.Context, id string) {
	if c.idCache == nil {
		return
	}
	idKey := cacheKey(ctx, id)
	request, ok := c.idCache.Get(idKey)
	if !ok {
		return
	}
	c.idCache.Remove(idKey)
	c.codeCache.Remove(cacheKey(ctx, request.Code))
}
