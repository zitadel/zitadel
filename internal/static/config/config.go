package config

import (
	"context"
	"encoding/json"
	"io"
	"net/url"
	"time"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/static"
	"github.com/caos/zitadel/internal/static/s3"
)

type AssetStorageConfig struct {
	Type   string
	Config static.Config
}

var storage = map[string]func() static.Config{
	"s3":   func() static.Config { return &s3.Config{} },
	"none": func() static.Config { return &NoStorage{} },
	"":     func() static.Config { return &NoStorage{} },
}

func (c *AssetStorageConfig) UnmarshalJSON(data []byte) error {
	var rc struct {
		Type   string
		Config json.RawMessage
	}

	if err := json.Unmarshal(data, &rc); err != nil {
		return errors.ThrowInternal(err, "STATIC-Bfn5r", "error parsing config")
	}

	c.Type = rc.Type

	var err error
	c.Config, err = newStorageConfig(c.Type, rc.Config)
	if err != nil {
		return err
	}

	return nil
}

func newStorageConfig(storageType string, configData []byte) (static.Config, error) {
	t, ok := storage[storageType]
	if !ok {
		return nil, errors.ThrowInternalf(nil, "STATIC-dsbjh", "config type %s not supported", storageType)
	}

	staticConfig := t()
	if len(configData) == 0 {
		return staticConfig, nil
	}

	if err := json.Unmarshal(configData, staticConfig); err != nil {
		return nil, errors.ThrowInternal(err, "STATIC-GB4nw", "Could not read config: %v")
	}

	return staticConfig, nil
}

var (
	errNoStorage = errors.ThrowInternal(nil, "STATIC-ashg4", "not configured")
)

type NoStorage struct{}

func (_ *NoStorage) NewStorage() (static.Storage, error) {
	return &NoStorage{}, nil
}

func (_ *NoStorage) CreateBucket(ctx context.Context, name, location string) error {
	return errNoStorage
}

func (_ *NoStorage) RemoveBucket(ctx context.Context, name string) error {
	return errNoStorage
}

func (_ *NoStorage) ListBuckets(ctx context.Context) ([]*domain.BucketInfo, error) {
	return nil, errNoStorage
}

func (_ *NoStorage) PutObject(ctx context.Context, bucketName, objectName, contentType string, object io.Reader, objectSize int64, createBucketIfNotExisting bool) (*domain.AssetInfo, error) {
	return nil, errNoStorage
}

func (_ *NoStorage) GetObjectInfo(ctx context.Context, bucketName, objectName string) (*domain.AssetInfo, error) {
	return nil, errNoStorage
}

func (_ *NoStorage) GetObject(ctx context.Context, bucketName, objectName string) (io.Reader, func() (*domain.AssetInfo, error), error) {
	return nil, nil, errNoStorage
}

func (_ *NoStorage) ListObjectInfos(ctx context.Context, bucketName, prefix string, recursive bool) ([]*domain.AssetInfo, error) {
	return nil, errNoStorage
}

func (_ *NoStorage) GetObjectPresignedURL(ctx context.Context, bucketName, objectName string, expiration time.Duration) (*url.URL, error) {
	return nil, errNoStorage
}

func (_ *NoStorage) RemoveObject(ctx context.Context, bucketName, objectName string) error {
	return errNoStorage
}

func (_ *NoStorage) RemoveObjects(ctx context.Context, bucketName, path string, recursive bool) error {
	return errNoStorage
}
