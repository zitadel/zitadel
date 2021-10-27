package common

type image struct {
	image     string
	runAsUser int64
}

func (i image) String() string { return i.image }

type dockerhubImage image

type zitadelImage image

var (
	CockroachImage        = dockerhubImage{"cockroachdb/cockroach:v21.1.8", 0}
	FlywayImage           = dockerhubImage{"flyway/flyway:7.12.1", 101}
	ZITADELImage          = zitadelImage{"caos/zitadel", 1000}
	ZITADELCockroachImage = zitadelImage{"caos/zitadel-cockroach", 1000}
	ZITADELOperatorImage  = zitadelImage{"caos/zitadel-operator", 1000}
)

func (z zitadelImage) RunAsUser() int64 {
	return z.runAsUser
}

func (z zitadelImage) Reference(customImageRegistry, version string) string {
	reg := "ghcr.io"
	if customImageRegistry != "" {
		reg = customImageRegistry
	}

	return concat(image(z), reg, version)
}

func (d dockerhubImage) RunAsUser() int64 {
	return d.runAsUser
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
