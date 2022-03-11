package crtlgitops

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/zitadel"
	"github.com/caos/zitadel/operator/zitadel/kinds/orb"
	"github.com/caos/zitadel/pkg/databases/db"
	macherrs "k8s.io/apimachinery/pkg/api/errors"
)

func ScaleDown(
	monitor mntr.Monitor,
	gitClient *git.Client,
	k8sClient *kubernetes.Client,
	version *string,
	gitops bool,
	dbClient db.Client,
) (bool, error) {
	noZitadel := false
	if err := zitadel.Takeoff(monitor, gitClient, orb.AdaptFunc("scaledown", version, gitops, []string{"scaledown"}, dbClient), k8sClient)(); err != nil {
		if macherrs.IsNotFound(err) {
			noZitadel = true
		} else {
			return noZitadel, err
		}
	}
	return noZitadel, nil
}

func ScaleUp(
	monitor mntr.Monitor,
	gitClient *git.Client,
	k8sClient *kubernetes.Client,
	version *string,
	gitops bool,
	dbClient db.Client,
) error {
	return zitadel.Takeoff(monitor, gitClient, orb.AdaptFunc("scaleup", version, gitops, []string{"scaleup"}, dbClient), k8sClient)()
}
