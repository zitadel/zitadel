package nginx

import (
	"fmt"
	"strconv"

	"github.com/caos/orbos/pkg/kubernetes/resources"
	"github.com/caos/orbos/pkg/kubernetes/resources/ingress"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/ingress/protocol/core"
)

var (
	DestroyIngress core.IngressDefinitionDestroyFunc = ingress.AdaptFuncToDestroy
	_              core.IngressDefinitionQueryFunc   = QueryIngress
)

const (
	backendPoolKey          = "nginx.ingress.kubernetes.io/backend-protocol"
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

func QueryIngress(
	namespace string,
	id labels.IDLabels,
	grpc bool,
	host,
	prefix,
	rewrite,
	service string,
	servicePort uint16,
	timeoutMS,
	connectTimeoutMS int,
	cors *core.CORS,
	controllerSpecifics map[string]interface{},
) (resources.QueryFunc, error) {

	timeoutMSStr := fmt.Sprintf("%dms", timeoutMS)
	connTimeoutMSStr := fmt.Sprintf("%dms", connectTimeoutMS)

	annotations := map[string]string{
		backendPoolKey:    "HTTP",
		rewriteKey:        rewrite,
		readTimeoutKey:    timeoutMSStr,
		sendTimeoutKey:    timeoutMSStr,
		connectTimeoutKey: connTimeoutMSStr,
	}

	if grpc {
		annotations[connectTimeoutKey] = "GRPC"
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

	return ingress.AdaptFuncToEnsure(
		namespace,
		id,
		host,
		prefix,
		service,
		servicePort,
		annotations,
	)
}
