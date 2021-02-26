package ingress

import (
	"testing"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/labels/mocklabels"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/core"
)

func TestAdaptFuncCover(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	cover := func(controller string) {
		_, _, err := Adapt("")(core.PathArguments{
			Monitor:             mntr.Monitor{},
			Namespace:           "",
			ID:                  mocklabels.Name,
			GRPC:                false,
			OriginCASecretName:  "",
			Prefix:              "",
			Rewrite:             "",
			Service:             "",
			ServicePort:         0,
			TimeoutMS:           0,
			ConnectTimeoutMS:    0,
			CORS:                nil,
			ControllerSpecifics: nil,
		})
		checkErr(err)
	}
	cover("NGINX")
	cover("Ambassador")
}
