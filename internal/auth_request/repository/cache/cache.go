package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/caos/zitadel/internal/auth_request/model"
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
	client, err := sql.Open("postgres", conf.Connection.ConnectionString())
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

func (c *AuthRequestCache) GetAuthRequestByID(_ context.Context, id string) (*model.AuthRequest, error) {
	var b []byte
	err := c.client.QueryRow("SELECT request FROM auth.auth_requests WHERE id = $1", id).Scan(&b)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, caos_errs.ThrowNotFound(err, "CACHE-d24aD", "auth request not found")
		}
		return nil, caos_errs.ThrowInternal(err, "CACHE-as3kj", "unable to get auth request from database")
	}
	request := &model.AuthRequest{Request: &model.AuthRequestOIDC{}} //TODO: ?
	err = json.Unmarshal(b, request)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "CACHE-2wshg", "unable to unmarshal auth request")
	}
	return request, nil
}

func (c *AuthRequestCache) SaveAuthRequest(_ context.Context, request *model.AuthRequest) error {
	return c.saveAuthRequest(request, "INSERT INTO auth.auth_requests (id, request) VALUES($1, $2)")
}

func (c *AuthRequestCache) UpdateAuthRequest(_ context.Context, request *model.AuthRequest) error {
	return c.saveAuthRequest(request, "UPDATE auth.auth_requests SET request = $2 WHERE id = $1")
}

func (c *AuthRequestCache) DeleteAuthRequest(_ context.Context, id string) error {
	_, err := c.client.Exec("DELETE FROM auth.auth_requests WHERE id = $1", id)
	if err != nil {
		return caos_errs.ThrowInternal(err, "CACHE-dsHw3", "unable to delete auth request")
	}
	return nil
}

func (c *AuthRequestCache) saveAuthRequest(request *model.AuthRequest, query string) error {
	b, err := json.Marshal(request)
	if err != nil {
		return caos_errs.ThrowInternal(err, "CACHE-os0GH", "unable to marshal auth request")
	}
	stmt, err := c.client.Prepare(query)
	if err != nil {
		return caos_errs.ThrowInternal(err, "CACHE-su3GK", "sql prepare failed")
	}
	_, err = stmt.Exec(request.ID, b)
	if err != nil {
		return caos_errs.ThrowInternal(err, "CACHE-sj8iS", "unable to save auth request")
	}
	return nil
}
