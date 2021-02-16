package bucket

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/database/kinds/backups/core"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"strings"
)

func BackupList() core.BackupListFunc {
	return func(monitor mntr.Monitor, name string, desired *tree.Tree) ([]string, error) {
		desiredKind, err := ParseDesiredV0(desired)
		if err != nil {
			return nil, errors.Wrap(err, "parsing desired state failed")
		}
		desired.Parsed = desiredKind

		if !monitor.IsVerbose() && desiredKind.Spec.Verbose {
			monitor.Verbose()
		}

		return listFilesWithFilter(desiredKind.Spec.ServiceAccountJSON.Value, desiredKind.Spec.Bucket, name)
	}
}

func listFilesWithFilter(serviceAccountJSON string, bucketName, name string) ([]string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(serviceAccountJSON)))

	if err != nil {
		return nil, err
	}
	bkt := client.Bucket(bucketName)

	names := make([]string, 0)
	it := bkt.Objects(ctx, &storage.Query{Prefix: name + "/"})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		parts := strings.Split(attrs.Name, "/")
		found := false
		for _, name := range names {
			if name == parts[1] {
				found = true
			}
		}
		if !found {
			names = append(names, parts[1])
		}
	}

	return names, nil
}
