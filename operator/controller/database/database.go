package database

import (
	"context"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/database"
	v1 "github.com/caos/zitadel/operator/api/database/v1"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Reconciler struct {
	kubernetes.ClientInt
	Monitor mntr.Monitor
	Scheme  *runtime.Scheme
	Version string
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	internalMonitor := r.Monitor.WithFields(map[string]interface{}{
		"kind":      "database",
		"namespace": req.NamespacedName,
	})

	desired, err := database.ReadCrd(r.ClientInt, req.Namespace, req.Name)
	if err != nil {
		return ctrl.Result{}, err
	}

	query, _, _, err := orbdb.AdaptFunc("", &r.Version, false, "database")(internalMonitor, desired, &tree.Tree{})
	if err != nil {
		return ctrl.Result{}, err
	}

	ensure, err := query(r.ClientInt, map[string]interface{}{})
	if err != nil {
		return ctrl.Result{}, err
	}

	if err := ensure(r.ClientInt); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Database{}).
		Complete(r)
}
