package s3

import (
	"database/sql"
	"encoding/json"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/zitadel/zitadel/internal/static"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	SSL             bool
	Location        string
	BucketPrefix    string
	MultiDelete     bool
}

func (c *Config) NewStorage() (static.Storage, error) {
	minioClient, err := minio.New(c.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessKeyID, c.SecretAccessKey, ""),
		Secure: c.SSL,
		Region: c.Location,
	})
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "MINIO-2n9fs", "Errors.Assets.Store.NotInitialized")
	}
	return &Minio{
		Client:       minioClient,
		Location:     c.Location,
		BucketPrefix: c.BucketPrefix,
		MultiDelete:  c.MultiDelete,
	}, nil
}

func NewStorage(_ *sql.DB, rawConfig map[string]interface{}) (static.Storage, error) {
	configData, err := json.Marshal(rawConfig)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "MINIO-Ef2f2", "could not map config")
	}
	c := new(Config)
	if err := json.Unmarshal(configData, c); err != nil {
		return nil, zerrors.ThrowInternal(err, "MINIO-GB4nw", "could not map config")
	}
	return c.NewStorage()
}
