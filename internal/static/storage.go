package static

import (
	"context"
	"database/sql"
	"io"
	"time"
)

type CreateStorage func(client *sql.DB, rawConfig map[string]interface{}) (Storage, error)

type Storage interface {
	PutObject(ctx context.Context, instanceID, location, resourceOwner, name, contentType string, objectType ObjectType, object io.Reader, objectSize int64) (*Asset, error)
	GetObject(ctx context.Context, instanceID, resourceOwner, name string) ([]byte, func() (*Asset, error), error)
	GetObjectInfo(ctx context.Context, instanceID, resourceOwner, name string) (*Asset, error)
	RemoveObject(ctx context.Context, instanceID, resourceOwner, name string) error
	RemoveObjects(ctx context.Context, instanceID, resourceOwner string, objectType ObjectType) error
	RemoveInstanceObjects(ctx context.Context, instanceID string) error
	//TODO: add functionality to move asset location
}

type ObjectType int32

const (
	ObjectTypeUserAvatar ObjectType = iota
	ObjectTypeStyling
)

func (o ObjectType) String() string {
	switch o {
	case ObjectTypeUserAvatar:
		return "0"
	case ObjectTypeStyling:
		return "1"
	default:
		return ""
	}
}

type Asset struct {
	InstanceID    string
	ResourceOwner string
	Name          string
	Hash          string
	Size          int64
	LastModified  time.Time
	Location      string
	ContentType   string
}

func (a *Asset) VersionedName() string {
	return a.Name + "?v=" + a.Hash
}
