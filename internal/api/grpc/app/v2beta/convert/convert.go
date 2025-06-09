package convert

import (
	"net/url"

	"github.com/zitadel/zitadel/internal/domain"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func LoginVersionToDomain(version *app.LoginVersion) (domain.LoginVersion, string, error) {
	switch v := version.GetVersion().(type) {
	case nil:
		return domain.LoginVersionUnspecified, "", nil
	case *app.LoginVersion_LoginV1:
		return domain.LoginVersion1, "", nil
	case *app.LoginVersion_LoginV2:
		_, err := url.Parse(v.LoginV2.GetBaseUri())
		return domain.LoginVersion2, v.LoginV2.GetBaseUri(), err
	default:
		return domain.LoginVersionUnspecified, "", nil
	}
}
