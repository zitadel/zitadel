package ambassador

import (
	"testing"

	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/caos/orbos/pkg/labels/mocklabels"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ambassador/grpc"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ambassador/hosts"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ambassador/http"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ambassador/ui"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
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
				"annotations":     map[string]string{},
				"labels":          map[string]string{},
				"resourceVersion": resourceVersion,
			},
		},
	}
	k8sClient.EXPECT().GetNamespacedCRDResource(group, version, kind, namespace, name).MinTimes(1).MaxTimes(1).Return(ret, nil)
}

func SetMappingsUI(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
) {
	group := "getambassador.io"
	version := "v2"
	kind := "Mapping"
	k8sClient.EXPECT().CheckCRD("mappings.getambassador.io").MinTimes(1).MaxTimes(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)

	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, ui.AccountsName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, ui.AccountsName, gomock.Any()).MinTimes(1).MaxTimes(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, ui.ConsoleName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, ui.ConsoleName, gomock.Any()).MinTimes(1).MaxTimes(1)
}

func SetMappingsHTTP(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
) {
	group := "getambassador.io"
	version := "v2"
	kind := "Mapping"
	k8sClient.EXPECT().CheckCRD("mappings.getambassador.io").MinTimes(1).MaxTimes(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)

	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, http.AdminRName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, http.AdminRName, gomock.Any()).MinTimes(1).MaxTimes(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, http.AuthorizeName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, http.AuthorizeName, gomock.Any()).MinTimes(1).MaxTimes(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, http.AuthRName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, http.AuthRName, gomock.Any()).MinTimes(1).MaxTimes(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, http.EndsessionName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, http.EndsessionName, gomock.Any()).MinTimes(1).MaxTimes(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, http.IssuerName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, http.IssuerName, gomock.Any()).MinTimes(1).MaxTimes(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, http.MgmtName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, http.MgmtName, gomock.Any()).MinTimes(1).MaxTimes(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, http.OauthName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, http.OauthName, gomock.Any()).MinTimes(1).MaxTimes(1)
}

func SetMappingsGRPC(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
) {
	group := "getambassador.io"
	version := "v2"
	kind := "Mapping"
	k8sClient.EXPECT().CheckCRD("mappings.getambassador.io").MinTimes(1).MaxTimes(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)

	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, grpc.AdminMName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, grpc.AdminMName, gomock.Any()).MinTimes(1).MaxTimes(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, grpc.AuthMName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, grpc.AuthMName, gomock.Any()).MinTimes(1).MaxTimes(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, grpc.MgmtMName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, grpc.MgmtMName, gomock.Any()).MinTimes(1).MaxTimes(1)
}

func SetHosts(
	k8sClient *kubernetesmock.MockClientInt,
	namespace string,
) {
	group := "getambassador.io"
	version := "v2"
	kind := "Host"
	k8sClient.EXPECT().CheckCRD("hosts.getambassador.io").MinTimes(1).MaxTimes(1).Return(&apixv1beta1.CustomResourceDefinition{}, nil)

	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, hosts.AccountsHostName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, hosts.AccountsHostName, gomock.Any()).MinTimes(1).MaxTimes(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, hosts.ApiHostName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, hosts.ApiHostName, gomock.Any()).MinTimes(1).MaxTimes(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, hosts.ConsoleHostName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, hosts.ConsoleHostName, gomock.Any()).MinTimes(1).MaxTimes(1)
	SetReturnResourceVersion(k8sClient, group, version, kind, namespace, hosts.IssuerHostName, "")
	k8sClient.EXPECT().ApplyNamespacedCRDResource(group, version, kind, namespace, hosts.IssuerHostName, gomock.Any()).MinTimes(1).MaxTimes(1)

}

func TestAmbassador_Adapt(t *testing.T) {

	monitor := mntr.Monitor{}
	namespace := "test"
	grpcURL := "grpc"
	httpURL := "http"
	uiURL := "ui"
	dns := &configuration.DNS{
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

	SetMappingsUI(k8sClient, namespace)
	SetMappingsHTTP(k8sClient, namespace)
	SetMappingsGRPC(k8sClient, namespace)
	SetHosts(k8sClient, namespace)

	query, _, err := AdaptFunc(monitor, mocklabels.Component, namespace, grpcURL, httpURL, uiURL, dns)
	assert.NoError(t, err)
	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	assert.NoError(t, err)
	assert.NoError(t, ensure(k8sClient))
}
