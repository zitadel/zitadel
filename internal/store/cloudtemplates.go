package store

import (
	"context"
	"io"

	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/zitadel/zitadel/internal/constants"
	"google.golang.org/api/option"
)

type Client struct {
	cl         *storage.Client
	bucketName string
	path       string
}

var transfer *Client

func Setup(path string) (*storage.Client, error) {
	var client *storage.Client
	var err error

	if len(constants.STORECREDENTIALSFILE) <= 0 {
		client, err = storage.NewClient(context.Background())
	} else {
		client, err = storage.NewClient(context.Background(), option.WithCredentialsFile(constants.STORECREDENTIALSFILE))
	}

	if err != nil {
		return nil, err
	}
	transfer = &Client{
		cl:         client,
		bucketName: constants.STOREBUCKETNAME,
		path:       path,
	}
	return client, nil
}

// Download file downloads an object
func (c *Client) DownloadFile(ctx io.Writer, fileName string) error {
	// Download an object with storage.Writer.
	reader, err := c.cl.Bucket(c.bucketName).Object(filepath.Join(c.path, fileName)).NewReader(context.Background())

	if err != nil {
		return err
	}
	defer reader.Close()

	if _, err = io.Copy(ctx, reader); err != nil {
		return err
	}

	return nil
}

func DownloadCaller(ctx io.Writer, filename string) error {
	file := strings.Join(strings.Split(filename, "%20"), " ")
	err := transfer.DownloadFile(ctx, file)
	return err
}
