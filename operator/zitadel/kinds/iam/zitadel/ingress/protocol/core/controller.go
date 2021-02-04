package core

import (
	"github.com/caos/orbos/pkg/kubernetes/resources"
	"github.com/caos/orbos/pkg/kubernetes/resources/ambassador/mapping"
	"github.com/caos/orbos/pkg/labels"
)

type CORS mapping.CORS

func (c *CORS) ToAmassadorCORS() *mapping.CORS {
	if c == nil {
		return nil
	}

	ambassadorCORS := mapping.CORS(*c)
	return &ambassadorCORS
}

type IngressDefinitionDestroyFunc func(namespace, name string) (resources.DestroyFunc, error)

type IngressDefinitionQueryFunc func(
	namespace string,
	labels labels.IDLabels,
	grpc bool,
	host,
	prefix,
	rewrite,
	service string,
	servicePort uint16,
	timeoutMS,
	connectTimeoutMS int,
	cors *CORS,
	controllerSpecifics map[string]interface{},
) (resources.QueryFunc, error)
