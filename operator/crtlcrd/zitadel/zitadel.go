package zitadel

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/operator"
	"github.com/zitadel/zitadel/pkg/databases"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/zitadel/zitadel/operator/api/zitadel"
	v1 "github.com/zitadel/zitadel/operator/api/zitadel/v1"
	orbz "github.com/zitadel/zitadel/operator/zitadel/kinds/orb"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Reconciler struct {
	kubernetes.ClientInt
	Monitor mntr.Monitor
	Scheme  *runtime.Scheme
	Version string
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, err error) {
	internalMonitor := r.Monitor.WithFields(map[string]interface{}{
		"kind":      "zitadel",
		"namespace": req.NamespacedName,
	})

	defer func() {
		r.Monitor.Error(err)
	}()

	if req.Namespace != zitadel.Namespace || req.Name != zitadel.Name {
		return res, fmt.Errorf("resource must be named %s and namespaced in %s", zitadel.Name, zitadel.Namespace)
	}

	dbClient, err := databases.NewConnection(r.Monitor, r.ClientInt, false, nil)
	if err != nil {
		return res, err
	}

	if err := Takeoff(internalMonitor, r.ClientInt, orbz.AdaptFunc("ensure", &r.Version, false, []string{"operator", "iam", "dbconnection"}, dbClient)); err != nil {
		return res, err
	}

	return res, nil
}

func Takeoff(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	adaptFunc operator.AdaptFunc,
) error {
	desired, err := zitadel.ReadCrd(k8sClient)
	if err != nil {
		return err
	}

	query, _, _, _, _, _, err := adaptFunc(monitor, desired, &tree.Tree{})
	if err != nil {
		return err
	}

	ensure, err := query(k8sClient, map[string]interface{}{})
	if err != nil {
		return err
	}

	if err := ensure(k8sClient); err != nil {
		return err
	}

	return nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Zitadel{}).
		Complete(r)
}
