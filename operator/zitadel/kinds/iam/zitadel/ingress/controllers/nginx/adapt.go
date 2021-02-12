package nginx

import (
	"fmt"
	"strconv"

	"github.com/caos/orbos/pkg/kubernetes"

	"github.com/caos/zitadel/operator"

	"github.com/caos/orbos/pkg/kubernetes/resources/ingress"
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
	return func(args core.PathArguments) (operator.QueryFunc, operator.DestroyFunc, error) {

		timeoutMSStr := fmt.Sprintf("%dms", args.TimeoutMS)
		connTimeoutMSStr := fmt.Sprintf("%dms", args.ConnectTimeoutMS)

		annotations := map[string]string{
			backendProtocolKey: "HTTP",
			rewriteKey:         args.Rewrite + "$1",
			readTimeoutKey:     timeoutMSStr,
			sendTimeoutKey:     timeoutMSStr,
			connectTimeoutKey:  connTimeoutMSStr,
		}

		if args.GRPC {
			annotations[backendProtocolKey] = "GRPC"
		}

		if args.CORS != nil {
			annotations[enableCorsKey] = "true"
			annotations[corsAllowOriginKey] = args.CORS.Origins
			annotations[corsAllowMethodsKey] = args.CORS.Methods
			annotations[corsAllowHeadersKey] = args.CORS.Headers
			annotations[corsAllowCredentialsKey] = strconv.FormatBool(args.CORS.Credentials)
			annotations[corsMaxAgeKey] = args.CORS.MaxAge
		}

		for k, v := range args.ControllerSpecifics {
			annotations[k] = fmt.Sprintf("%v", v)
		}

		query, err := ingress.AdaptFuncToEnsure(
			args.Namespace,
			args.ID,
			virtualHost,
			args.Prefix+"(.*)",
			args.Service,
			args.ServicePort,
			annotations,
		)
		if err != nil {
			return nil, nil, err
		}

		destroy, err := ingress.AdaptFuncToDestroy(args.Namespace, args.ID.Name())
		if err != nil {
			return nil, nil, err
		}

		return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
				return operator.QueriersToEnsureFunc(args.Monitor, false, []operator.QueryFunc{
					operator.ResourceQueryToZitadelQuery(query),
				}, k8sClient, queried)
			},
			operator.DestroyersToDestroyFunc(args.Monitor, []operator.DestroyFunc{
				operator.ResourceDestroyToZitadelDestroy(destroy)}),
			nil
	}
}
