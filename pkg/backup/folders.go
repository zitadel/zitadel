package backup

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func ListS3Folders(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
	bucketName string,
	path string,
	region string,
) ([]string, error) {
	s3Client := getS3Client(endpoint, accessKeyID, secretAccessKey, region)
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

	return getAllAssetFoldersFromObjectList(objects, path), nil
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

	return getAllAssetFoldersFromObjectList(objects, path), nil
}

func getAllAssetFoldersFromObjectList(objects []string, path string) []string {
	folders := make([]string, 0)
	for _, object := range objects {
		if strings.HasPrefix(object, path) {
			assetFolder := getAssetFolderFromObject(path, object)

			alreadyListed := false
			for _, folder := range folders {
				if assetFolder == folder {
					alreadyListed = true
				}
			}
			if !alreadyListed {
				folders = append(folders, assetFolder)
			}
		}
	}
	return folders
}

func getAssetFolderFromObject(path, object string) string {
	folder := strings.TrimPrefix(object, path+"/")
	for {
		folderT, _ := filepath.Split(folder)
		folder = strings.TrimSuffix(folderT, string(filepath.Separator))
		if !strings.Contains(folder, string(filepath.Separator)) {
			break
		}
	}
	return folder
}
