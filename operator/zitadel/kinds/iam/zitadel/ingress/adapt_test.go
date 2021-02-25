package ingress

import (
	"testing"

	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"

	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/golang/mock/gomock"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/labels/mocklabels"
)

func TestAdaptFunc(t *testing.T) {
	q, _, err := AdaptFunc(mntr.Monitor{}, mocklabels.Api, "", "", 0, "", 0, "", 0, &configuration.Ingress{
		Subdomains: &configuration.Subdomains{},
	})
	if err != nil {
		t.Fatal(err)
	}

	q(kubernetesmock.NewMockClientInt(gomock.NewController(t)), make(map[string]interface{}))
}
