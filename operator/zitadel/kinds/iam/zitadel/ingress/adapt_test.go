package ingress

import (
	"testing"

	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/labels/mocklabels"
)

func TestAdaptFunc(t *testing.T) {
	cover := func(controller string) {
		AdaptFunc(mntr.Monitor{}, mocklabels.Api, "", "", 0, "", 0, "", 0, &configuration.Ingress{
			Subdomains: &configuration.Subdomains{},
			Controller: controller,
		})
	}
	cover("NGINX")
	cover("Ambassador")
}
