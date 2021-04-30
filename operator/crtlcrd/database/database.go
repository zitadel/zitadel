package database

import (
	"context"
	"fmt"

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

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, err error) {
	internalMonitor := r.Monitor.WithFields(map[string]interface{}{
		"kind":      "database",
		"namespace": req.NamespacedName,
	})

	defer func() {
		r.Monitor.Error(err)
	}()

	if req.Namespace != database.Namespace || req.Name != database.Name {
		return res, fmt.Errorf("resource must be named %s and namespaced in %s", database.Name, database.Namespace)
	}

	desired, err := database.ReadCrd(r.ClientInt)
	if err != nil {
		internalMonitor.Error(err)
		return res, err
	}

	query, _, _, _, _, _, err := orbdb.AdaptFunc("", &r.Version, false, "database")(internalMonitor, desired, &tree.Tree{})
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
		For(&v1.Database{}).
		Complete(r)
}
