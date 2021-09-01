package backup

import (
	"github.com/caos/zitadel/pkg/backup/cockroachdb"
	"github.com/caos/zitadel/pkg/backup/rsync"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

func RsyncBackupS3ToS3(
	backupName string,
	backupNameEnv string,
	destinationName string,
	destinationEndpoint string,
	destinationAKID string,
	destinationSAK string,
	destinationBucket string,
	sourceName string,
	sourceEndpoint string,
	sourceAKID string,
	sourceSAK string,
	sourceBucketPrefix string,
	configFilePath string,
) error {

	assetBuckets, err := ListS3AssetBuckets(sourceEndpoint, sourceAKID, sourceSAK, sourceBucketPrefix)
	if err != nil {
		return err
	}

	sourcePart, err := rsync.GetConfigPartS3(sourceName, sourceEndpoint, sourceAKID, sourceSAK)
	if err != nil {
		return err
	}

	destPart, err := rsync.GetConfigPartS3(destinationName, destinationEndpoint, destinationAKID, destinationSAK)
	if err != nil {
		return err
	}

	config := strings.Join([]string{
		sourcePart,
		destPart,
	}, "\n")

	if err := ioutil.WriteFile(configFilePath, []byte(config), fs.ModePerm); err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, assetBucket := range assetBuckets {
		wg.Add(1)
		sourceBucket := assetBucket
		if err := runCommand(
			rsync.GetCommand(
				configFilePath,
				sourceName,
				sourceBucket,
				destinationName,
				getS3FullPath(destinationBucket, backupName, backupNameEnv),
			),
		); err != nil {
			return err
		}
	}
	wg.Wait()

	return nil
}

func RsyncBackupS3ToGCS(
	backupName string,
	backupNameEnv string,
	destinationName string,
	destinationSaJsonPath string,
	destinationBucket string,
	sourceName string,
	sourceEndpoint string,
	sourceAKID string,
	sourceSAK string,
	sourceBucketPrefix string,
	configFilePath string,
) error {

	assetBuckets, err := ListS3AssetBuckets(sourceEndpoint, sourceAKID, sourceSAK, sourceBucketPrefix)
	if err != nil {
		return err
	}

	sourcePart, err := rsync.GetConfigPartS3(sourceName, sourceEndpoint, sourceAKID, sourceSAK)
	if err != nil {
		return err
	}

	destPart, err := rsync.GetConfigPartGCS(destinationName, destinationSaJsonPath)
	if err != nil {
		return err
	}

	config := strings.Join([]string{
		sourcePart,
		destPart,
	}, "\n")

	if err := ioutil.WriteFile(configFilePath, []byte(config), fs.ModePerm); err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, assetBucket := range assetBuckets {
		wg.Add(1)
		sourceBucket := assetBucket
		if err := runCommand(
			rsync.GetCommand(
				configFilePath,
				sourceName,
				sourceBucket,
				destinationName,
				getS3FullPath(destinationBucket, backupName, backupNameEnv),
			),
		); err != nil {
			return err
		}
	}
	wg.Wait()

	return nil
}

func CockroachBackupToGCS(
	certsFolder string,
	bucketName string,
	backupName string,
	backupNameEnv string,
	host string,
	port string,
	serviceAccountPath string,
) error {
	return runCommand(
		cockroachdb.GetBackupToGCS(
			certsFolder,
			host,
			port,
			bucketName,
			getS3Path(backupName, backupNameEnv),
			serviceAccountPath,
		),
	)
}

func getS3FullPath(
	bucket string,
	backupName string,
	backupNameEnv string,
) string {
	return filepath.Join(bucket, getS3Path(backupName, backupNameEnv))
}
func getS3Path(
	backupName string,
	backupNameEnv string,
) string {
	return filepath.Join(backupName, "${"+backupNameEnv+"}")
}

func CockroachBackupToS3(
	certsFolder string,
	bucketName string,
	backupName string,
	backupNameEnv string,
	host string,
	port string,
	accessKeyIDPath string,
	secretAccessKeyPath string,
	sessionTokenPath string,
	endpoint string,
	region string,
) error {
	return runCommand(
		cockroachdb.GetBackupToS3(
			certsFolder,
			host,
			port,
			bucketName,
			getS3Path(backupName, backupNameEnv),
			accessKeyIDPath,
			secretAccessKeyPath,
			sessionTokenPath,
			endpoint,
			region,
		),
	)
}
