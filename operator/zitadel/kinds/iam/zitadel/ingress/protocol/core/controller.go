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

type PathAdapter func(PathArguments) (operator.QueryFunc, operator.DestroyFunc, error)

type PathArguments struct {
	Monitor                                      mntr.Monitor
	Namespace                                    string
	ID                                           labels.IDLabels
	GRPC                                         bool
	OriginCASecretName, Prefix, Rewrite, Service string
	ServicePort                                  uint16
	TimeoutMS, ConnectTimeoutMS                  int
	CORS                                         *CORS
	ControllerSpecifics                          map[string]interface{}
}
