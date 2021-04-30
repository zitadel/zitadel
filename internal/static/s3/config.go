package s3

type AssetStorage struct {
	Type   string
	Config S3Config
}

type S3Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	SSL             bool
	Location        string
}
