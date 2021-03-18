package zitadel

import (
	"context"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/zitadel"
	v1 "github.com/caos/zitadel/operator/api/zitadel/v1"
	orbz "github.com/caos/zitadel/operator/zitadel/kinds/orb"
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

	desired, err := zitadel.ReadCrd(r.ClientInt)
	if err != nil {
		return res, err
	}

	query, _, _, _, err := orbz.AdaptFunc(nil, "ensure", &r.Version, false, []string{"operator", "iam"})(internalMonitor, desired, &tree.Tree{})
	if err != nil {
		internalMonitor.Error(err)
		return res, err
	}

	ensure, err := query(r.ClientInt, map[string]interface{}{})
	if err != nil {
		internalMonitor.Error(err)
		return res, err
	}

	if err := ensure(r.ClientInt); err != nil {
		internalMonitor.Error(err)
		return res, err
	}

	return res, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Zitadel{}).
		Complete(r)
}
