package ingress

import (
	"testing"

	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/labels/mocklabels"
)

func TestAdaptFuncCover(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	cover := func(controller string) {
		_, _, err := AdaptFunc(mntr.Monitor{}, mocklabels.Api, "", "", 0, "", 0, "", 0, &configuration.Ingress{
			Subdomains: &configuration.Subdomains{},
			Controller: controller,
		})
		checkErr(err)
	}
	cover("NGINX")
	cover("Ambassador")
}
