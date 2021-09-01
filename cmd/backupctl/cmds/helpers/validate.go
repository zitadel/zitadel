package helpers

import "errors"

func ValidateBackupFlags(
	backupName string,
	backupNameEnv string,
) error {
	if backupName == "" {
		return errors.New("missing or empty backup name parameter")
	}
	if backupNameEnv == "" {
		return errors.New("missing or empty backup name environment variable parameter")
	}
	return nil
}

func ValidateGCSFlags(
	saJSONPath string,
	bucket string,
) error {
	if saJSONPath == "" {
		return errors.New("missing or empty service account json path parameter")
	}
	if bucket == "" {
		return errors.New("missing or empty GCS bucket name parameter")
	}
	return nil
}

func ValidateDestinationS3Flags(
	endpoint string,
	akid string,
	sak string,
	bucket string,
) error {
	if endpoint == "" {
		return errors.New("missing or empty destination endpoint parameter")
	}
	if akid == "" {
		return errors.New("missing or empty destination access key ID parameter")
	}
	if sak == "" {
		return errors.New("missing or empty destination secret access key parameter")
	}
	if bucket == "" {
		return errors.New("missing or empty destination bucket parameter")
	}
	return nil
}

func ValidateSourceS3Flags(
	endpoint string,
	akid string,
	sak string,
	prefix string,
) error {
	if endpoint == "" {
		return errors.New("missing or empty source endpoint parameter")
	}
	if akid == "" {
		return errors.New("missing or empty source access key ID parameter")
	}
	if sak == "" {
		return errors.New("missing or empty source secret access key parameter")
	}
	if prefix == "" {
		return errors.New("missing or empty source bucket prefix parameter")
	}
	return nil
}

func ValidateCockroachFlags(
	certsFolder string,
	host string,
	port string,
) error {
	if certsFolder == "" {
		return errors.New("missing or empty cockroachdb used certs-folder parameter")
	}
	if host == "" {
		return errors.New("missing or empty cockroachdb host parameter")
	}
	if port == "" {
		return errors.New("missing or empty cockroachdb port parameter")
	}
	return nil
}
