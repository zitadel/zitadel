package nginx

import (
	"fmt"
	"strconv"

	"github.com/caos/orbos/pkg/kubernetes"

	"github.com/caos/zitadel/operator"

	"github.com/caos/orbos/mntr"

	"github.com/caos/orbos/pkg/kubernetes/resources/ingress"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/core"
)

const (
	backendProtocolKey      = "nginx.ingress.kubernetes.io/backend-protocol"
	rewriteKey              = "nginx.ingress.kubernetes.io/rewrite-target"
	connectTimeoutKey       = "nginx.ingress.kubernetes.io/proxy-connect-timeout"
	sendTimeoutKey          = "nginx.ingress.kubernetes.io/proxy-send-timeout"
	readTimeoutKey          = "nginx.ingress.kubernetes.io/proxy-read-timeout"
	enableCorsKey           = "nginx.ingress.kubernetes.io/enable-cors"
	corsAllowOriginKey      = "nginx.ingress.kubernetes.io/cors-allow-origin"
	corsAllowMethodsKey     = "nginx.ingress.kubernetes.io/cors-allow-methods"
	corsAllowHeadersKey     = "nginx.ingress.kubernetes.io/cors-expose-headers"
	corsAllowCredentialsKey = "nginx.ingress.kubernetes.io/cors-allow-credentials"
	corsMaxAgeKey           = "nginx.ingress.kubernetes.io/cors-max-age"
)

var _ core.HostAdapter = Adapt

func Adapt(virtualHost string) core.PathAdapter {
	return func(
		monitor mntr.Monitor,
		namespace string,
		id labels.IDLabels,
		grpc bool,
		originCASecretName,
		prefix,
		rewrite,
		service string,
		servicePort uint16,
		timeoutMS,
		connectTimeoutMS int,
		cors *core.CORS,
		controllerSpecifics map[string]interface{},
	) (operator.QueryFunc, operator.DestroyFunc, error) {

		timeoutMSStr := fmt.Sprintf("%dms", timeoutMS)
		connTimeoutMSStr := fmt.Sprintf("%dms", connectTimeoutMS)

		annotations := map[string]string{
			backendProtocolKey: "HTTP",
			rewriteKey:         rewrite + "$1",
			readTimeoutKey:     timeoutMSStr,
			sendTimeoutKey:     timeoutMSStr,
			connectTimeoutKey:  connTimeoutMSStr,
		}

		if grpc {
			annotations[backendProtocolKey] = "GRPC"
		}

		if cors != nil {
			annotations[enableCorsKey] = "true"
			annotations[corsAllowOriginKey] = cors.Origins
			annotations[corsAllowMethodsKey] = cors.Methods
			annotations[corsAllowHeadersKey] = cors.Headers
			annotations[corsAllowCredentialsKey] = strconv.FormatBool(cors.Credentials)
			annotations[corsMaxAgeKey] = cors.MaxAge
		}

		for k, v := range controllerSpecifics {
			annotations[k] = fmt.Sprintf("%v", v)
		}

		query, err := ingress.AdaptFuncToEnsure(
			namespace,
			id,
			virtualHost,
			prefix+"(.*)",
			service,
			servicePort,
			annotations,
		)
		if err != nil {
			return nil, nil, err
		}

		destroy, err := ingress.AdaptFuncToDestroy(namespace, id.Name())
		if err != nil {
			return nil, nil, err
		}

		return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
				return operator.QueriersToEnsureFunc(monitor, false, []operator.QueryFunc{
					operator.ResourceQueryToZitadelQuery(query),
				}, k8sClient, queried)
			},
			operator.DestroyersToDestroyFunc(monitor, []operator.DestroyFunc{
				operator.ResourceDestroyToZitadelDestroy(destroy)}),
			nil
	}
}
