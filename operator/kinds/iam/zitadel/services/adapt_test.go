package services

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes/mock"
	"github.com/golang/mock/gomock"
	"gotest.tools/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

func GetExpectedService(
	namespace string,
	labels map[string]string,
	grpcPortName string,
	grpcServiceName string,
	grpcPort int,
	httpPortName string,
	httpServiceName string,
	httpPort int,
	uiPortName string,
	uiServiceName string,
	uiPort int,
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

	return []*corev1.Service{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      grpcServiceName,
				Namespace: namespace,
				Labels:    labels,
			},
			Spec: corev1.ServiceSpec{
				Ports:                    grpcPorts,
				Selector:                 labels,
				Type:                     corev1.ServiceType(""),
				PublishNotReadyAddresses: false,
				ClusterIP:                "",
				ExternalName:             "",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      httpServiceName,
				Namespace: namespace,
				Labels:    labels,
			},
			Spec: corev1.ServiceSpec{
				Ports:                    httpPorts,
				Selector:                 labels,
				Type:                     corev1.ServiceType(""),
				PublishNotReadyAddresses: false,
				ClusterIP:                "",
				ExternalName:             "",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      uiServiceName,
				Namespace: namespace,
				Labels:    labels,
			},
			Spec: corev1.ServiceSpec{
				Ports:                    uiPorts,
				Selector:                 labels,
				Type:                     corev1.ServiceType(""),
				PublishNotReadyAddresses: false,
				ClusterIP:                "",
				ExternalName:             "",
			},
		},
	}

}

func TestServices_AdaptEnsure1(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	namespace := "test"
	labels := map[string]string{"test": "test"}
	grpcPortName := "grpc"
	grpcServiceName := "grpc"
	grpcPort := 1
	httpPortName := "http"
	httpServiceName := "http"
	httpPort := 2
	uiPortName := "ui"
	uiServiceName := "ui"
	uiPort := 3

	for _, rsc := range GetExpectedService(
		namespace,
		labels,
		grpcPortName,
		grpcServiceName,
		grpcPort,
		httpPortName,
		httpServiceName,
		httpPort,
		uiPortName,
		uiServiceName,
		uiPort) {

		client.EXPECT().ApplyService(rsc).Times(1)
	}

	query, _, err := AdaptFunc(
		mntr.Monitor{},
		namespace,
		labels,
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

func TestServices_AdaptEnsure2(t *testing.T) {
	client := kubernetesmock.NewMockClientInt(gomock.NewController(t))

	namespace := "test0"
	labels := map[string]string{"test0": "test0"}
	grpcPortName := "grpc"
	grpcServiceName := "grpc1"
	grpcPort := 11
	httpPortName := "http"
	httpServiceName := "http2"
	httpPort := 22
	uiPortName := "ui"
	uiServiceName := "ui3"
	uiPort := 33

	for _, rsc := range GetExpectedService(
		namespace,
		labels,
		grpcPortName,
		grpcServiceName,
		grpcPort,
		httpPortName,
		httpServiceName,
		httpPort,
		uiPortName,
		uiServiceName,
		uiPort) {

		client.EXPECT().ApplyService(rsc).Times(1)
	}

	query, _, err := AdaptFunc(
		mntr.Monitor{},
		namespace,
		labels,
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
	labels := map[string]string{"test00": "test00"}
	grpcPortName := "grpc"
	grpcServiceName := "grpc11"
	grpcPort := 111
	httpPortName := "http"
	httpServiceName := "http22"
	httpPort := 222
	uiPortName := "ui"
	uiServiceName := "ui33"
	uiPort := 333

	for _, rsc := range GetExpectedService(
		namespace,
		labels,
		grpcPortName,
		grpcServiceName,
		grpcPort,
		httpPortName,
		httpServiceName,
		httpPort,
		uiPortName,
		uiServiceName,
		uiPort) {

		client.EXPECT().ApplyService(rsc).Times(1)
	}

	query, _, err := AdaptFunc(
		mntr.Monitor{},
		namespace,
		labels,
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
	labels := map[string]string{"test": "test"}
	grpcPortName := "grpc"
	grpcServiceName := "grpc"
	grpcPort := 1
	httpPortName := "http"
	httpServiceName := "http"
	httpPort := 2
	uiPortName := "ui"
	uiServiceName := "ui"
	uiPort := 3

	for _, rsc := range GetExpectedService(
		namespace,
		labels,
		grpcPortName,
		grpcServiceName,
		grpcPort,
		httpPortName,
		httpServiceName,
		httpPort,
		uiPortName,
		uiServiceName,
		uiPort) {

		client.EXPECT().DeleteService(rsc.Namespace, rsc.Name).Times(1)
	}

	_, destroy, err := AdaptFunc(
		mntr.Monitor{},
		namespace,
		labels,
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
	labels := map[string]string{"test0": "test0"}
	grpcPortName := "grpc"
	grpcServiceName := "grpc1"
	grpcPort := 11
	httpPortName := "http"
	httpServiceName := "http2"
	httpPort := 22
	uiPortName := "ui"
	uiServiceName := "ui3"
	uiPort := 33

	for _, rsc := range GetExpectedService(
		namespace,
		labels,
		grpcPortName,
		grpcServiceName,
		grpcPort,
		httpPortName,
		httpServiceName,
		httpPort,
		uiPortName,
		uiServiceName,
		uiPort) {

		client.EXPECT().DeleteService(rsc.Namespace, rsc.Name).Times(1)
	}

	_, destroy, err := AdaptFunc(
		mntr.Monitor{},
		namespace,
		labels,
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
	labels := map[string]string{"test00": "test00"}
	grpcPortName := "grpc"
	grpcServiceName := "grpc11"
	grpcPort := 111
	httpPortName := "http"
	httpServiceName := "http22"
	httpPort := 222
	uiPortName := "ui"
	uiServiceName := "ui33"
	uiPort := 333

	for _, rsc := range GetExpectedService(
		namespace,
		labels,
		grpcPortName,
		grpcServiceName,
		grpcPort,
		httpPortName,
		httpServiceName,
		httpPort,
		uiPortName,
		uiServiceName,
		uiPort) {

		client.EXPECT().DeleteService(rsc.Namespace, rsc.Name).Times(1)
	}

	_, destroy, err := AdaptFunc(
		mntr.Monitor{},
		namespace,
		labels,
		grpcServiceName,
		grpcPort,
		httpServiceName,
		httpPort,
		uiServiceName,
		uiPort)

	assert.NilError(t, err)
	assert.NilError(t, destroy(client))
}
