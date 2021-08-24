package common

type image string

func (i image) String() string { return string(i) }

type dockerhubImage image

type zitadelImage image

const (
	CockroachImage       dockerhubImage = "cockroachdb/cockroach:v21.1.7"
	PostgresImage        dockerhubImage = "postgres:9.6.17"
	FlywayImage          dockerhubImage = "flyway/flyway:7.12.1"
	AlpineImage          dockerhubImage = "alpine:3.11"
	ZITADELImage         zitadelImage   = "caos/zitadel"
	BackupImage          zitadelImage   = "caos/zitadel-crbackup"
	ZITADELOperatorImage zitadelImage   = "caos/zitadel-operator"
)

func (z zitadelImage) Reference(customImageRegistry, version string) string {

	reg := "ghcr.io"
	if customImageRegistry != "" {
		reg = customImageRegistry
	}

	return concat(image(z), reg, version)
}

func (d dockerhubImage) Reference(customImageRegistry string) string {
	return concat(image(d), customImageRegistry, "")
}

func concat(img image, customImageRegistry, version string) string {
	str := img.String()

	if customImageRegistry != "" {
		str = customImageRegistry + "/" + str
	}

	if version != "" {
		str = str + ":" + version
	}
	return str
}
