package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

type AuthRequestCache struct {
	client *database.DB
}

func Start(dbClient *database.DB) *AuthRequestCache {
	return &AuthRequestCache{
		client: dbClient,
	}
}

func (c *AuthRequestCache) Health(ctx context.Context) error {
	return c.client.PingContext(ctx)
}

func (c *AuthRequestCache) GetAuthRequestByID(ctx context.Context, id string) (*domain.AuthRequest, error) {
	return c.getAuthRequest("id", id, authz.GetInstance(ctx).InstanceID())
}

func (c *AuthRequestCache) GetAuthRequestByCode(ctx context.Context, code string) (*domain.AuthRequest, error) {
	return c.getAuthRequest("code", code, authz.GetInstance(ctx).InstanceID())
}

func (c *AuthRequestCache) SaveAuthRequest(_ context.Context, request *domain.AuthRequest) error {
	return c.saveAuthRequest(request, "INSERT INTO auth.auth_requests (id, request, instance_id, creation_date, change_date, request_type) VALUES($1, $2, $3, $4, $4, $5)", request.CreationDate, request.Request.Type())
}

func (c *AuthRequestCache) UpdateAuthRequest(_ context.Context, request *domain.AuthRequest) error {
	if request.ChangeDate.IsZero() {
		request.ChangeDate = time.Now()
	}
	return c.saveAuthRequest(request, "UPDATE auth.auth_requests SET request = $2, instance_id = $3, change_date = $4, code = $5 WHERE id = $1", request.ChangeDate, request.Code)
}

func (c *AuthRequestCache) DeleteAuthRequest(ctx context.Context, id string) error {
	_, err := c.client.Exec("DELETE FROM auth.auth_requests WHERE instance_id = $1 and id = $2", authz.GetInstance(ctx).InstanceID(), id)
	if err != nil {
		return caos_errs.ThrowInternal(err, "CACHE-dsHw3", "unable to delete auth request")
	}
	return nil
}

func (c *AuthRequestCache) getAuthRequest(key, value, instanceID string) (*domain.AuthRequest, error) {
	var b []byte
	var requestType domain.AuthRequestType
	query := fmt.Sprintf("SELECT request, request_type FROM auth.auth_requests WHERE instance_id = $1 and %s = $2", key)
	err := c.client.QueryRow(query, instanceID, value).Scan(&b, &requestType)
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

func (c *AuthRequestCache) saveAuthRequest(request *domain.AuthRequest, query string, date time.Time, param interface{}) error {
	b, err := json.Marshal(request)
	if err != nil {
		return caos_errs.ThrowInternal(err, "CACHE-os0GH", "Errors.Internal")
	}
	_, err = c.client.Exec(query, request.ID, b, request.InstanceID, date, param)
	if err != nil {
		return caos_errs.ThrowInternal(err, "CACHE-su3GK", "Errors.Internal")
	}
	return nil
}
