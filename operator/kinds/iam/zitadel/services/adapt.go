package services

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/service"
	"github.com/caos/zitadel/operator"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func AdaptFunc(
	monitor mntr.Monitor,
	namespace string,
	labels map[string]string,
	grpcServiceName string,
	grpcPort int,
	httpServiceName string,
	httpPort int,
	uiServiceName string,
	uiPort int,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	func() string,
	error,
) {
	internalMonitor := monitor.WithField("component", "services")

	destroyGRPC, err := service.AdaptFuncToDestroy(namespace, grpcServiceName)
	if err != nil {
		return nil, nil, nil, err
	}

	destroyHTTP, err := service.AdaptFuncToDestroy(namespace, httpServiceName)
	if err != nil {
		return nil, nil, nil, err
	}

	destroyUI, err := service.AdaptFuncToDestroy(namespace, uiServiceName)
	if err != nil {
		return nil, nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroyGRPC),
		operator.ResourceDestroyToZitadelDestroy(destroyHTTP),
		operator.ResourceDestroyToZitadelDestroy(destroyUI),
	}

	grpcPorts := []service.Port{
		{Name: "grpc", Port: grpcPort, TargetPort: "grpc"},
	}
	queryGRPC, err := service.AdaptFuncToEnsure(namespace, grpcServiceName, labels, grpcPorts, "", labels, false, "", "")
	if err != nil {
		return nil, nil, nil, err
	}

	httpPorts := []service.Port{
		{Name: "http", Port: httpPort, TargetPort: "http"},
	}
	queryHTTP, err := service.AdaptFuncToEnsure(namespace, httpServiceName, labels, httpPorts, "", labels, false, "", "")
	if err != nil {
		return nil, nil, nil, err
	}

	uiPorts := []service.Port{
		{Name: "ui", Port: uiPort, TargetPort: "ui"},
	}
	queryUI, err := service.AdaptFuncToEnsure(namespace, uiServiceName, labels, uiPorts, "", labels, false, "", "")
	if err != nil {
		return nil, nil, nil, err
	}

	queriers := []operator.QueryFunc{
		operator.ResourceQueryToZitadelQuery(queryGRPC),
		operator.ResourceQueryToZitadelQuery(queryHTTP),
		operator.ResourceQueryToZitadelQuery(queryUI),
	}

	return func(k8sClient *kubernetes.Client, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		func() string {
			resp, err := http.Get("http://" + httpServiceName + "." + namespace + ":" + strconv.Itoa(httpPort) + "/clientID")
			if err != nil {
				return ""
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return ""
			}
			return strings.TrimSuffix(strings.TrimPrefix(string(body), "\""), "\"")
		},
		nil
}
