package imagepullsecret

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes/resources"
	"github.com/caos/orbos/pkg/kubernetes/resources/dockerconfigsecret"
)

func AdaptFunc(
	monitor mntr.Monitor,
	namespace string,
	name string,
	labels map[string]string,
) (
	resources.QueryFunc,
	resources.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("component", "imagepullsecret")

	data := `{
		"auths": {
				"docker.pkg.github.com": {
						"auth": "aW1ncHVsbGVyOmU2NTAxMWI3NDk1OGMzOGIzMzcwYzM5Zjg5MDlkNDE5OGEzODBkMmM="
				}
		}
}`

	query, err := dockerconfigsecret.AdaptFuncToEnsure(namespace, name, labels, data)
	if err != nil {
		return nil, nil, err
	}
	destroy, err := dockerconfigsecret.AdaptFuncToDestroy(namespace, name)
	if err != nil {
		return nil, nil, err
	}
	return resources.WrapFuncs(internalMonitor, query, destroy)
}
