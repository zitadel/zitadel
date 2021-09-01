package backup

import (
	"github.com/caos/zitadel/pkg/backup/cockroachdb"
	"github.com/caos/zitadel/pkg/backup/rsync"
	"io/fs"
	"io/ioutil"
	"strings"
	"sync"
)

func RsyncRestoreS3ToS3(
	backupName string,
	backupNameEnv string,
	destinationName string,
	destinationEndpoint string,
	destinationAKID string,
	destinationSAK string,
	sourceName string,
	sourceEndpoint string,
	sourceAKID string,
	sourceSAK string,
	sourceBucket string,
	configFilePath string,
) error {

	assetBuckets, err := ListS3Folders(sourceEndpoint, sourceAKID, sourceSAK, sourceBucket, getS3Path(backupName, backupNameEnv))
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
				getS3FullPath(sourceBucket, backupName, backupNameEnv),
				destinationName,
				assetBucket,
			),
		); err != nil {
			return err
		}
	}
	wg.Wait()

	return nil
}

func RsyncRestoreGCSToS3(
	backupName string,
	backupNameEnv string,
	destinationName string,
	destinationEndpoint string,
	destinationAKID string,
	destinationSAK string,
	sourceName string,
	sourceSaJsonPath string,
	sourceBucket string,
	configFilePath string,
) error {
	assetBuckets, err := ListGCSFolders(sourceSaJsonPath, sourceBucket, getS3Path(backupName, backupNameEnv))
	if err != nil {
		return err
	}

	sourcePart, err := rsync.GetConfigPartGCS(sourceName, sourceSaJsonPath)
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
				getS3FullPath(sourceBucket, backupName, backupNameEnv),
				destinationName,
				assetBucket,
			),
		); err != nil {
			return err
		}
	}
	wg.Wait()

	return nil
}

func CockroachRestoreFromGCS(
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

func CockroachRestoreFromS3(
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
