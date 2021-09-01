package backup

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io/ioutil"
	"strings"
)

func getS3Client(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
) *s3.Client {

	const defaultRegion = "us-east-1"
	staticResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               endpoint, // or where ever you ran minio
			SigningRegion:     defaultRegion,
			HostnameImmutable: true,
		}, nil
	})

	cfg := aws.Config{
		Region:           defaultRegion,
		Credentials:      credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
		EndpointResolver: staticResolver,
	}

	return s3.NewFromConfig(cfg)
}
func ListS3AssetBuckets(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
	prefix string,
) ([]string, error) {

	s3Client := getS3Client(endpoint, accessKeyID, secretAccessKey)

	input := &s3.ListBucketsInput{}

	result, err := s3Client.ListBuckets(context.Background(), input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, err
	}

	buckets := make([]string, 0)
	for _, bucket := range result.Buckets {
		name := *bucket.Name
		if strings.HasPrefix(name, prefix) {
			buckets = append(buckets, *bucket.Name)
		}
	}

	return buckets, nil
}

func ListS3Folders(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
	bucketName string,
	path string,
) ([]string, error) {
	s3Client := getS3Client(endpoint, accessKeyID, secretAccessKey)
	ctx := context.Background()
	input := &s3.ListObjectsV2Input{
		Bucket: &bucketName,
		Prefix: &path,
	}

	output, err := s3Client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, err
	}

	objects := make([]string, 0)
	for _, obj := range output.Contents {
		objects = append(objects, *obj.Key)
	}

	return objects, nil
}

func ListGCSFolders(
	saJSONPath string,
	bucketName string,
	path string,
) ([]string, error) {
	ctx := context.Background()
	data, err := ioutil.ReadFile(saJSONPath)
	if err != nil {
		return nil, err
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(data))
	if err != nil {
		return nil, err
	}

	iter := client.Bucket(bucketName).Objects(
		ctx,
		&storage.Query{
			Prefix:   path,
			Versions: false,
		},
	)
	objects := make([]string, 0)
	for {
		obj, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj.Name)
	}

	return objects, nil
}
