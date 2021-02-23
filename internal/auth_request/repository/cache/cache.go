package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/caos/zitadel/internal/domain"

	"github.com/caos/zitadel/internal/config/types"
	caos_errs "github.com/caos/zitadel/internal/errors"
)

type Config struct {
	Connection types.SQL
}

type AuthRequestCache struct {
	client *sql.DB
}

func Start(conf Config) (*AuthRequestCache, error) {
	client, err := conf.Connection.Start()
	if err != nil {
		return nil, caos_errs.ThrowPreconditionFailed(err, "SQL-9qBtr", "unable to open database connection")
	}
	return &AuthRequestCache{
		client: client,
	}, nil
}

func (c *AuthRequestCache) Health(ctx context.Context) error {
	return c.client.PingContext(ctx)
}

func (c *AuthRequestCache) GetAuthRequestByID(_ context.Context, id string) (*domain.AuthRequest, error) {
	return c.getAuthRequest("id", id)
}

func (c *AuthRequestCache) GetAuthRequestByCode(_ context.Context, code string) (*domain.AuthRequest, error) {
	return c.getAuthRequest("code", code)
}

func (c *AuthRequestCache) SaveAuthRequest(_ context.Context, request *domain.AuthRequest) error {
	return c.saveAuthRequest(request, "INSERT INTO auth.auth_requests (id, request, request_type) VALUES($1, $2, $3)", request.Request.Type())
}

func (c *AuthRequestCache) UpdateAuthRequest(_ context.Context, request *domain.AuthRequest) error {
	return c.saveAuthRequest(request, "UPDATE auth.auth_requests SET request = $2, code = $3 WHERE id = $1", request.Code)
}

func (c *AuthRequestCache) DeleteAuthRequest(_ context.Context, id string) error {
	_, err := c.client.Exec("DELETE FROM auth.auth_requests WHERE id = $1", id)
	if err != nil {
		return caos_errs.ThrowInternal(err, "CACHE-dsHw3", "unable to delete auth request")
	}
	return nil
}

func (c *AuthRequestCache) getAuthRequest(key, value string) (*domain.AuthRequest, error) {
	var b []byte
	var requestType domain.AuthRequestType
	query := fmt.Sprintf("SELECT request, request_type FROM auth.auth_requests WHERE %s = $1", key)
	err := c.client.QueryRow(query, value).Scan(&b, &requestType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, caos_errs.ThrowNotFound(err, "CACHE-d24aD", "Errors.AuthRequest.NotFound")
		}
		return nil, caos_errs.ThrowInternal(err, "CACHE-as3kj", "Errors.Internal")
	}
	request, err := domain.NewAuthRequestFromType(requestType)
	if err == nil {
		err = json.Unmarshal(b, request)
	}
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "CACHE-2wshg", "Errors.Internal")
	}
	return request, nil
}

func (c *AuthRequestCache) saveAuthRequest(request *domain.AuthRequest, query string, param interface{}) error {
	b, err := json.Marshal(request)
	if err != nil {
		return caos_errs.ThrowInternal(err, "CACHE-os0GH", "Errors.Internal")
	}
	stmt, err := c.client.Prepare(query)
	if err != nil {
		return caos_errs.ThrowInternal(err, "CACHE-su3GK", "Errors.Internal")
	}
	_, err = stmt.Exec(request.ID, b, param)
	if err != nil {
		return caos_errs.ThrowInternal(err, "CACHE-sj8iS", "Errors.Internal")
	}
	return nil
}
