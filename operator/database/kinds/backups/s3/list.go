package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/secret/read"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/database/kinds/backups/core"
	"github.com/pkg/errors"
	"strings"
)

func BackupList() core.BackupListFunc {
	return func(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, name string, desired *tree.Tree) ([]string, error) {
		desiredKind, err := ParseDesiredV0(desired)
		if err != nil {
			return nil, errors.Wrap(err, "parsing desired state failed")
		}
		desired.Parsed = desiredKind

		if !monitor.IsVerbose() && desiredKind.Spec.Verbose {
			monitor.Verbose()
		}

		valueAKI, err := read.GetSecretValue(k8sClient, desiredKind.Spec.AccessKeyID, desiredKind.Spec.ExistingAccessKeyID)
		if err != nil {
			return nil, err
		}

		valueSAK, err := read.GetSecretValue(k8sClient, desiredKind.Spec.SecretAccessKey, desiredKind.Spec.ExistingSecretAccessKey)
		if err != nil {
			return nil, err
		}

		valueST, err := read.GetSecretValue(k8sClient, desiredKind.Spec.SessionToken, desiredKind.Spec.ExistingSessionToken)
		if err != nil {
			return nil, err
		}

		return listFilesWithFilter(valueAKI, valueSAK, valueST, desiredKind.Spec.Region, desiredKind.Spec.Endpoint, desiredKind.Spec.Bucket, name)
	}
}

func listFilesWithFilter(akid, secretkey, sessionToken string, region string, endpoint string, bucketName, name string) ([]string, error) {
	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           endpoint,
			SigningRegion: region,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(akid, secretkey, sessionToken)),
		config.WithEndpointResolver(customResolver),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	prefix := name + "/"
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	}
	paginator := s3.NewListObjectsV2Paginator(client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
		o.Limit = 10
	})

	names := make([]string, 0)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		for _, value := range output.Contents {
			if strings.HasPrefix(*value.Key, prefix) {
				names = append(names, *value.Key)
			}
		}
	}
	return names, nil
}
