package ambassador

import (
	"testing"

	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels/mocklabels"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/controllers/ambassador/hosts"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	apixv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func SetReturnResourceVersion(
	k8sClient *kubernetesmock.MockClientInt,
	group,
	version,
	kind,
	namespace,
	name string,
	resourceVersion string,
) {
	ret := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"resourceVersion": resourceVersion,
			},
		},
	}
	k8sClient.EXPECT().GetNamespacedCRDResource(group, version, kind, namespace, name).Return(ret, nil)
}

func SetHosts(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
) {
	group := "getambassador.io"
	version := "v2"
	kind := "Host"
	k8sClient.EXPECT().CheckCRD("hosts.getambassador.io").Times(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)

	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, hosts.AccountsHostName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, hosts.AccountsHostName, gomock.Any()).Times(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, hosts.ApiHostName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, hosts.ApiHostName, gomock.Any()).Times(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, hosts.ConsoleHostName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, hosts.ConsoleHostName, gomock.Any()).Times(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, hosts.IssuerHostName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, hosts.IssuerHostName, gomock.Any()).Times(1)

}

func TestAmbassador_Adapt(t *testing.T) {

	monitor := mntr.Monitor{}
	namespace := "test"
	dns := &configuration.Ingress{
		Domain:    "",
		TlsSecret: "",
		Subdomains: &configuration.Subdomains{
			Accounts: "",
			API:      "",
			Console:  "",
			Issuer:   "",
		},
	}
	k8sClient := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	SetHosts(k8sClient, namespace)

	query, _, err := AdaptFunc(monitor, mocklabels.Component, namespace, dns)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}
