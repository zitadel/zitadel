package static

import (
	"context"
	"database/sql"
	"io"
	"time"
)

type CreateStorage func(client *sql.DB, rawConfig map[string]interface{}) (Storage, error)

type Storage interface {
	PutObject(ctx context.Context, tenantID, location, resourceOwner, name, contentType string, objectType ObjectType, object io.Reader, objectSize int64) (*Asset, error)
	GetObject(ctx context.Context, tenantID, resourceOwner, name string) ([]byte, func() (*Asset, error), error)
	GetObjectInfo(ctx context.Context, tenantID, resourceOwner, name string) (*Asset, error)
	RemoveObject(ctx context.Context, tenantID, resourceOwner, name string) error
	RemoveObjects(ctx context.Context, tenantID, resourceOwner string, objectType ObjectType) error
	//TODO: add functionality to move asset location
}

type ObjectType int32

const (
	ObjectTypeUserAvatar = iota
	ObjectTypeStyling
)

type Asset struct {
	TenantID      string
	ResourceOwner string
	Name          string
	Hash          string
	Size          int64
	LastModified  time.Time
	Location      string
	ContentType   string
}
