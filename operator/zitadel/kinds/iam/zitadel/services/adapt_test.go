package services

import (
	"testing"

	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/labels/mocklabels"

	"github.com/caos/orbos/mntr"
	kubernetesmock "github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/golang/mock/gomock"
	"gotest.tools/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func GetExpectedService(
	namespace string,
	zitadelPodSelector *labels.Selector,
	grpcPortName string,
	grpcServiceName *labels.Name,
	grpcPort uint16,
	httpPortName string,
	httpServiceName *labels.Name,
	httpPort uint16,
	uiPortName string,
	uiServiceName *labels.Name,
	uiPort uint16,
) []*corev1.Service {

	grpcPorts := []corev1.ServicePort{{
		Name:       grpcPortName,
		Protocol:   "",
		Port:       int32(grpcPort),
		TargetPort: intstr.Parse(grpcPortName),
		NodePort:   int32(0),
	},
	}

	httpPorts := []corev1.ServicePort{{
		Name:       httpPortName,
		Protocol:   "",
		Port:       int32(httpPort),
		TargetPort: intstr.Parse(httpPortName),
		NodePort:   int32(0),
	},
	}

	uiPorts := []corev1.ServicePort{{
		Name:       uiPortName,
		Protocol:   "",
		Port:       int32(uiPort),
		TargetPort: intstr.Parse(uiPortName),
		NodePort:   int32(0),
	},
	}

	zitadelPodSelectorMap := labels.MustK8sMap(zitadelPodSelector)

	return []*corev1.Service{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      grpcServiceName.Name(),
				Namespace: namespace,
				Labels:    labels.MustK8sMap(grpcServiceName),
			},
			Spec: corev1.ServiceSpec{
				Ports:                    grpcPorts,
				Selector:                 zitadelPodSelectorMap,
				Type:                     "",
				PublishNotReadyAddresses: false,
				ClusterIP:                "",
				ExternalName:             "",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      httpServiceName.Name(),
				Namespace: namespace,
				Labels:    labels.MustK8sMap(labels.AsSelectable(httpServiceName)),
			},
			Spec: corev1.ServiceSpec{
				Ports:                    httpPorts,
				Selector:                 zitadelPodSelectorMap,
				Type:                     "",
				PublishNotReadyAddresses: false,
				ClusterIP:                "",
				ExternalName:             "",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uiServiceName.Name(),
				Namespace: namespace,
				Labels:    labels.MustK8sMap(uiServiceName),
			},
			Spec: corev1.ServiceSpec{
				Ports:                    uiPorts,
				Selector:                 zitadelPodSelectorMap,
				Type:                     "",
				PublishNotReadyAddresses: false,
				ClusterIP:                "",
				ExternalName:             "",
			},
		},
	}
}

func serviceLabels(name ...string) (*labels.Component, *labels.Selector, []*labels.Name) {
	componentLabels := mocklabels.Component
	podSelectorLabels := labels.DeriveNameSelector(labels.MustForName(componentLabels, "zitadel"), false)
	nameLabels := make([]*labels.Name, len(name))
	for idx := range name {
		nameLabels[idx] = labels.MustForName(componentLabels, name[idx])
	}
	return componentLabels, podSelectorLabels, nameLabels
}

func TestServices_AdaptEnsure1(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	namespace := "test"
	grpcPortName := "grpc"
	grpcServiceName := "grpc"
	var grpcPort uint16 = 1
	httpPortName := "http"
	httpServiceName := "http"
	var httpPort uint16 = 2
	uiPortName := "ui"
	uiServiceName := "ui"
	var uiPort uint16 = 3

	componentLabels, podSelectorLabels, nameLabels := serviceLabels(grpcServiceName, httpServiceName, uiServiceName)

	for _, rsc := range GetExpectedService(
		namespace,
		podSelectorLabels,
		grpcPortName,
		nameLabels[0],
		grpcPort,
		httpPortName,
		nameLabels[1],
		httpPort,
		uiPortName,
		nameLabels[2],
		uiPort,
	) {
		client.EXPECT().ApplyService(rsc).Times(1)
	}

	query, _, err := AdaptFunc(
		mntr.Monitor{},
		componentLabels,
		podSelectorLabels,
		namespace,
		grpcServiceName,
		grpcPort,
		httpServiceName,
		httpPort,
		uiServiceName,
		uiPort,
	)

	assert.NilError(t, err)
	ensure, err := query(client, nil)
	assert.NilError(t, err)
	assert.NilError(t, ensure(client))
}

func TestServices_AdaptEnsure2(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	namespace := "test0"
	grpcPortName := "grpc"
	grpcServiceName := "grpc1"
	var grpcPort uint16 = 11
	httpPortName := "http"
	httpServiceName := "http2"
	var httpPort uint16 = 22
	uiPortName := "ui"
	uiServiceName := "ui3"
	var uiPort uint16 = 33

	componentLabels, podSelectorLabels, nameLabels := serviceLabels(grpcServiceName, httpServiceName, uiServiceName)

	for _, rsc := range GetExpectedService(
		namespace,
		podSelectorLabels,
		grpcPortName,
		nameLabels[0],
		grpcPort,
		httpPortName,
		nameLabels[1],
		httpPort,
		uiPortName,
		nameLabels[2],
		uiPort,
	) {

		client.EXPECT().ApplyService(rsc).Times(1)
	}

	query, _, err := AdaptFunc(
		mntr.Monitor{},
		componentLabels,
		podSelectorLabels,
		namespace,
		grpcServiceName,
		grpcPort,
		httpServiceName,
		httpPort,
		uiServiceName,
		uiPort)

	assert.NilError(t, err)
	ensure, err := query(client, nil)
	assert.NilError(t, err)
	assert.NilError(t, ensure(client))
}

func TestServices_AdaptEnsure3(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	namespace := "test00"
	grpcPortName := "grpc"
	grpcServiceName := "grpc11"
	var grpcPort uint16 = 111
	httpPortName := "http"
	httpServiceName := "http22"
	var httpPort uint16 = 222
	uiPortName := "ui"
	uiServiceName := "ui33"
	var uiPort uint16 = 333

	componentLabels, podSelectorLabels, nameLabels := serviceLabels(grpcServiceName, httpServiceName, uiServiceName)

	for _, rsc := range GetExpectedService(
		namespace,
		podSelectorLabels,
		grpcPortName,
		nameLabels[0],
		grpcPort,
		httpPortName,
		nameLabels[1],
		httpPort,
		uiPortName,
		nameLabels[2],
		uiPort,
	) {

		client.EXPECT().ApplyService(rsc).Times(1)
	}

	query, _, err := AdaptFunc(
		mntr.Monitor{},
		componentLabels,
		podSelectorLabels,
		namespace,
		grpcServiceName,
		grpcPort,
		httpServiceName,
		httpPort,
		uiServiceName,
		uiPort)

	assert.NilError(t, err)
	ensure, err := query(client, nil)
	assert.NilError(t, err)
	assert.NilError(t, ensure(client))
}

func TestServices_AdaptDestroy1(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	namespace := "test"
	grpcPortName := "grpc"
	grpcServiceName := "grpc"
	var grpcPort uint16 = 1
	httpPortName := "http"
	httpServiceName := "http"
	var httpPort uint16 = 2
	uiPortName := "ui"
	uiServiceName := "ui"
	var uiPort uint16 = 3

	componentLabels, podSelectorLabels, nameLabels := serviceLabels(grpcServiceName, httpServiceName, uiServiceName)

	for _, rsc := range GetExpectedService(
		namespace,
		podSelectorLabels,
		grpcPortName,
		nameLabels[0],
		grpcPort,
		httpPortName,
		nameLabels[1],
		httpPort,
		uiPortName,
		nameLabels[2],
		uiPort,
	) {

		client.EXPECT().DeleteService(rsc.Namespace, rsc.Name).Times(1)
	}

	_, destroy, err := AdaptFunc(
		mntr.Monitor{},
		componentLabels,
		podSelectorLabels,
		namespace,
		grpcServiceName,
		grpcPort,
		httpServiceName,
		httpPort,
		uiServiceName,
		uiPort)

	assert.NilError(t, err)
	assert.NilError(t, destroy(client))
}

func TestServices_AdaptDestroy2(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	namespace := "test0"
	grpcPortName := "grpc"
	grpcServiceName := "grpc1"
	var grpcPort uint16 = 11
	httpPortName := "http"
	httpServiceName := "http2"
	var httpPort uint16 = 22
	uiPortName := "ui"
	uiServiceName := "ui3"
	var uiPort uint16 = 33

	componentLabels, podSelectorLabels, nameLabels := serviceLabels(grpcServiceName, httpServiceName, uiServiceName)

	for _, rsc := range GetExpectedService(
		namespace,
		podSelectorLabels,
		grpcPortName,
		nameLabels[0],
		grpcPort,
		httpPortName,
		nameLabels[1],
		httpPort,
		uiPortName,
		nameLabels[2],
		uiPort,
	) {

		client.EXPECT().DeleteService(rsc.Namespace, rsc.Name).Times(1)
	}

	_, destroy, err := AdaptFunc(
		mntr.Monitor{},
		componentLabels,
		podSelectorLabels,
		namespace,
		grpcServiceName,
		grpcPort,
		httpServiceName,
		httpPort,
		uiServiceName,
		uiPort)

	assert.NilError(t, err)
	assert.NilError(t, destroy(client))
}

func TestServices_AdaptDestroy3(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	namespace := "test00"
	grpcPortName := "grpc"
	grpcServiceName := "grpc11"
	var grpcPort uint16 = 111
	httpPortName := "http"
	httpServiceName := "http22"
	var httpPort uint16 = 222
	uiPortName := "ui"
	uiServiceName := "ui33"
	var uiPort uint16 = 333

	componentLabels, podSelectorLabels, nameLabels := serviceLabels(grpcServiceName, httpServiceName, uiServiceName)

	for _, rsc := range GetExpectedService(
		namespace,
		podSelectorLabels,
		grpcPortName,
		nameLabels[0],
		grpcPort,
		httpPortName,
		nameLabels[1],
		httpPort,
		uiPortName,
		nameLabels[2],
		uiPort,
	) {

		client.EXPECT().DeleteService(rsc.Namespace, rsc.Name).Times(1)
	}

	_, destroy, err := AdaptFunc(
		mntr.Monitor{},
		componentLabels,
		podSelectorLabels,
		namespace,
		grpcServiceName,
		grpcPort,
		httpServiceName,
		httpPort,
		uiServiceName,
		uiPort)

	assert.NilError(t, err)
	assert.NilError(t, destroy(client))
}
