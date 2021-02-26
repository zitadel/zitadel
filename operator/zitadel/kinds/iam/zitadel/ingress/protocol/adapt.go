package protocol

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/core"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/grpc"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/http"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/ui"
)

func AdaptFunc(
	monitor mntr.Monitor,
	componentLabels *labels.Component,
	namespaceStr string,
	grpcService string,
	grpcPort uint16,
	httpService string,
	httpPort uint16,
	uiService string,
	uiPort uint16,
	dns *configuration.Ingress,
	controllerSpecifics map[string]string,
	hostAdapter core.HostAdapter,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	apiAdapter := hostAdapter(dns.Subdomains.API + "." + dns.Domain)
	accountsAdapter := hostAdapter(dns.Subdomains.Accounts + "." + dns.Domain)
	consoleAdapter := hostAdapter(dns.Subdomains.Console + "." + dns.Domain)
	issuerAdapter := hostAdapter(dns.Subdomains.Issuer + "." + dns.Domain)

	queryGRPC, destroyGRPC, err := grpc.AdaptFunc(
		monitor,
		componentLabels,
		namespaceStr,
		grpcService,
		grpcPort,
		controllerSpecifics,
		dns.TlsSecret,
		apiAdapter,
	)
	if err != nil {
		return nil, nil, err
	}

	queryUI, destroyHTTP, err := ui.AdaptFunc(
		monitor,
		componentLabels,
		namespaceStr,
		uiService,
		uiPort,
		controllerSpecifics,
		dns.TlsSecret,
		consoleAdapter,
		accountsAdapter,
	)
	if err != nil {
		return nil, nil, err
	}

	queryHTTP, destroyUI, err := http.AdaptFunc(
		monitor,
		componentLabels,
		namespaceStr,
		httpService,
		httpPort,
		controllerSpecifics,
		dns.TlsSecret,
		apiAdapter,
		accountsAdapter,
		issuerAdapter,
	)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		destroyGRPC,
		destroyUI,
		destroyHTTP,
	}

	queriers := []operator.QueryFunc{
		queryHTTP,
		queryUI,
		queryGRPC,
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(monitor, true, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(monitor, destroyers),
		nil
}
