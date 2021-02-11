package core

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes/resources/ambassador/mapping"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
)

type CORS mapping.CORS

func (c *CORS) ToAmassadorCORS() *mapping.CORS {
	if c == nil {
		return nil
	}

	ambassadorCORS := mapping.CORS(*c)
	return &ambassadorCORS
}

type HostAdapter func(virtualHost string) PathAdapter

type PathAdapter func(
	monitor mntr.Monitor,
	namespace string,
	labels labels.IDLabels,
	grpc bool,
	originCASecretName,
	prefix,
	rewrite,
	service string,
	servicePort uint16,
	timeoutMS,
	connectTimeoutMS int,
	cors *CORS,
	controllerSpecifics map[string]interface{}) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
)
