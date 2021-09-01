package restore

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/zitadel/cmd/backupctl/cmds/helpers"
	"github.com/caos/zitadel/pkg/backup"
	"github.com/spf13/cobra"
)

func S3Command(monitor mntr.Monitor) *cobra.Command {
	var (
		backupName     string
		backupNameEnv  string
		assetEndpoint  string
		assetAKID      string
		assetSAK       string
		sourceEndpoint string
		sourceBucket   string
		sourceAKID     string
		sourceSAK      string
		configPath     string
		certsDir       string
		host           string
		port           string
		cmd            = &cobra.Command{
			Use:   "s3",
			Short: "Backup to S3 storage",
			Long:  "Backup to S3 storage",
		}
	)

	flags := cmd.Flags()
	flags.StringVar(&backupName, "backupname", "", "Backupname used in destination file path")
	flags.StringVar(&backupNameEnv, "backupnameenv", "", "Backupnameenv used in destination file path")
	flags.StringVar(&assetEndpoint, "asset-endpoint", "", "Endpoint for the asset S3 storage")
	flags.StringVar(&assetAKID, "asset-akid", "", "AccessKeyID for the asset S3 storage")
	flags.StringVar(&assetSAK, "asset-sak", "", "SecretAccessKey for the asset S3 storage")
	flags.StringVar(&sourceEndpoint, "source-endpoint", "", "Endpoint for the source S3 storage")
	flags.StringVar(&sourceAKID, "source-akid", "", "AccessKeyID for the source S3 storage")
	flags.StringVar(&sourceSAK, "source-sak", "", "SecretAccessKey for the source S3 storage")
	flags.StringVar(&sourceBucket, "source-bucket", "", "Bucketname in the source S3 storage")
	flags.StringVar(&configPath, "configpath", "", "Path used to save rsync configuration")
	flags.StringVar(&certsDir, "certs-dir", "", "Folder with certificates used to connect to cockroachdb")
	flags.StringVar(&host, "host", "", "Host used to connect to cockroachdb")
	flags.StringVar(&port, "port", "", "Port used to connect to cockroachdb")

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {

		if err := helpers.ValidateBackupFlags(
			backupName,
			backupNameEnv,
		); err != nil {
			return err
		}

		if err := helpers.ValidateDestinationS3Flags(
			sourceEndpoint,
			sourceAKID,
			sourceSAK,
			sourceBucket,
		); err != nil {
			return err
		}

		if err := helpers.ValidateSourceS3Flags(
			assetEndpoint,
			assetAKID,
			assetSAK,
			"notempty",
		); err != nil {
			return err
		}

		if err := helpers.ValidateCockroachFlags(
			certsDir,
			host,
			port,
		); err != nil {
			return err
		}

		if err := backup.RsyncRestoreS3ToS3(
			backupName,
			backupNameEnv,
			"destination",
			assetEndpoint,
			assetAKID,
			assetSAK,
			"source",
			sourceEndpoint,
			sourceAKID,
			sourceSAK,
			sourceBucket,
			configPath,
		); err != nil {
			return err
		}

		if err := backup.CockroachRestoreFromS3(
			certsDir,
			sourceBucket,
			backupName,
			backupNameEnv,
			host,
			port,
			sourceAKID,
			sourceSAK,
			"",
			sourceEndpoint,
			"",
		); err != nil {
			return err
		}

		return nil
	}
	return cmd
}
