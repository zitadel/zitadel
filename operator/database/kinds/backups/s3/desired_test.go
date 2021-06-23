package s3

import (
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

const (
	masterkey  = "testMk"
	cron       = "testCron"
	bucketName = "testBucket"
	region     = "testRegion"
	endpoint   = "testEndpoint"
	akid       = "testAKID"
	sak        = "testSAK"
	st         = "testST"
	yamlFile   = `kind: databases.caos.ch/BucketBackup
version: v0
spec:
    verbose: true
    cron: testCron
    bucket: testBucket
    region: testRegion
    endpoint: testEndpoint
    accessKeyID:
        encryption: AES256
        encoding: Base64
        value: l7GEXvmCT8hBXereT4FIG4j5vKQIycjS
    secretAccessKey:
        encryption: AES256
        encoding: Base64
        value: NWYnOpFpME-9FESqWi0bFQ3M6e0iNQw=
    sessionToken:
        encryption: AES256
        encoding: Base64
        value: xVY9pEXuh0Wbf2P2X_yThXwqRX08sA==
`

	yamlFileWithoutSecret = `kind: databases.caos.ch/BucketBackup
version: v0
spec:
    verbose: true
    cron: testCron
    bucket: testBucket
    endpoint: testEndpoint
    region: testRegion
`
	yamlEmpty = `kind: databases.caos.ch/BucketBackup
version: v0`
)

var (
	desired = DesiredV0{
		Common: &tree.Common{
			Kind:    "databases.caos.ch/BucketBackup",
			Version: "v0",
		},
		Spec: &Spec{
			Verbose:  true,
			Cron:     cron,
			Bucket:   bucketName,
			Endpoint: endpoint,
			Region:   region,
			AccessKeyID: &secret.Secret{
				Value:      akid,
				Encryption: "AES256",
				Encoding:   "Base64",
			},
			SecretAccessKey: &secret.Secret{
				Value:      sak,
				Encryption: "AES256",
				Encoding:   "Base64",
			},
			SessionToken: &secret.Secret{
				Value:      st,
				Encryption: "AES256",
				Encoding:   "Base64",
			},
		},
	}
	desiredWithoutSecret = DesiredV0{
		Common: &tree.Common{
			Kind:    "databases.caos.ch/BucketBackup",
			Version: "v0",
		},
		Spec: &Spec{
			Verbose:  true,
			Cron:     cron,
			Bucket:   bucketName,
			Region:   region,
			Endpoint: endpoint,
		},
	}
	desiredEmpty = DesiredV0{
		Common: &tree.Common{
			Kind:    "databases.caos.ch/BucketBackup",
			Version: "v0",
		},
		Spec: &Spec{
			Verbose:  false,
			Cron:     "",
			Bucket:   "",
			Endpoint: "",
			Region:   "",
			AccessKeyID: &secret.Secret{
				Value: "",
			},
			SecretAccessKey: &secret.Secret{
				Value: "",
			},
			SessionToken: &secret.Secret{
				Value: "",
			},
		},
	}

	desiredNil = DesiredV0{
		Common: &tree.Common{
			Kind:    "databases.caos.ch/BucketBackup",
			Version: "v0",
		},
	}
)

func marshalYaml(t *testing.T, masterkey string, struc *DesiredV0) []byte {
	secret.Masterkey = masterkey
	data, err := yaml.Marshal(struc)
	assert.NoError(t, err)
	return data
}

func unmarshalYaml(t *testing.T, masterkey string, yamlFile []byte) *tree.Tree {
	secret.Masterkey = masterkey
	desiredTree := &tree.Tree{}
	assert.NoError(t, yaml.Unmarshal(yamlFile, desiredTree))
	return desiredTree
}

func getDesiredTree(t *testing.T, masterkey string, desired *DesiredV0) *tree.Tree {
	return unmarshalYaml(t, masterkey, marshalYaml(t, masterkey, desired))
}

func TestBucket_DesiredParse(t *testing.T) {
	assert.Equal(t, yamlFileWithoutSecret, string(marshalYaml(t, masterkey, &desiredWithoutSecret)))

	desiredTree := unmarshalYaml(t, masterkey, []byte(yamlFile))
	desiredKind, err := ParseDesiredV0(desiredTree)
	assert.NoError(t, err)
	assert.Equal(t, &desired, desiredKind)
}

func TestBucket_DesiredNotZero(t *testing.T) {
	desiredTree := unmarshalYaml(t, masterkey, []byte(yamlFile))
	desiredKind, err := ParseDesiredV0(desiredTree)
	assert.NoError(t, err)
	assert.False(t, desiredKind.Spec.IsZero())
}

func TestBucket_DesiredZero(t *testing.T) {
	desiredTree := unmarshalYaml(t, masterkey, []byte(yamlEmpty))
	desiredKind, err := ParseDesiredV0(desiredTree)
	assert.NoError(t, err)
	assert.True(t, desiredKind.Spec.IsZero())
}
