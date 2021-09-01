package backup

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/zitadel/cmd/backupctl/cmds/helpers"
	"github.com/caos/zitadel/pkg/backup"
	"github.com/spf13/cobra"
)

func GCSCommand(monitor mntr.Monitor) *cobra.Command {
	var (
		backupName     string
		backupNameEnv  string
		assetEndpoint  string
		assetAKID      string
		assetSAK       string
		assetPrefix    string
		destBucket     string
		destSAJSONPath string
		configPath     string
		certsDir       string
		host           string
		port           string
		cmd            = &cobra.Command{
			Use:   "gcs",
			Short: "Backup to GCS Bucket",
			Long:  "Backup to GCS Bucket",
		}
	)

	flags := cmd.Flags()
	flags.StringVar(&backupName, "backupname", "", "Backupname used in destination file path")
	flags.StringVar(&backupNameEnv, "backupnameenv", "", "Backupnameenv used in destination file path")
	flags.StringVar(&assetEndpoint, "asset-endpoint", "", "Endpoint for the asset S3 storage")
	flags.StringVar(&assetAKID, "asset-akid", "", "AccessKeyID for the asset S3 storage")
	flags.StringVar(&assetSAK, "asset-sak", "", "SecretAccessKey for the asset S3 storage")
	flags.StringVar(&assetPrefix, "asset-prefix", "", "Bucket-Prefix in the asset S3 storage")
	flags.StringVar(&destSAJSONPath, "destination-sajsonpath", "~/sa.json", "Path to where ServiceAccount-json will be written for the destination GCS")
	flags.StringVar(&destBucket, "destination-bucket", "", "Bucketname in the destination GCS")
	flags.StringVar(&configPath, "configpath", "~/rsync.conf", "Path used to save rsync configuration")
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

		if err := helpers.ValidateGCSFlags(
			destSAJSONPath,
			destBucket,
		); err != nil {
			return err
		}

		if err := helpers.ValidateSourceS3Flags(
			assetEndpoint,
			assetAKID,
			assetSAK,
			assetPrefix,
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

		if err := backup.RsyncBackupS3ToGCS(
			backupName,
			backupNameEnv,
			"destination",
			destSAJSONPath,
			destBucket,
			"source",
			assetEndpoint,
			assetAKID,
			assetSAK,
			assetPrefix,
			configPath,
		); err != nil {
			return err
		}

		if err := backup.CockroachBackupToGCS(
			certsDir,
			destBucket,
			backupName,
			backupNameEnv,
			host,
			port,
			destSAJSONPath,
		); err != nil {
			return err
		}

		return nil
	}
	return cmd
}
