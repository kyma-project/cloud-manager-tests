/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloudresources

import (
	"context"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/abstractions"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/actions"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
	apimachineryapi "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cloudresourcesv1beta1 "github.com/kyma-project/cloud-resources-control-plane/api/cloud-resources/v1beta1"
)

// NfsShareReconciler reconciles a NfsShare object
type NfsShareReconciler struct {
	client.Client
	record.EventRecorder
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cloud-resources.kyma-project.io,resources=nfsshares,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cloud-resources.kyma-project.io,resources=nfsshares/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cloud-resources.kyma-project.io,resources=nfsshares/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NfsShare object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *NfsShareReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO: this should be moved into separate reconciler package
	err := composedAction.ComposeActions(
		"vpcPeering",
		actions.LoadObj,
		actions.LoadKyma,
	)(ctx, actions.NewState(
		composedAction.NewState(r.Client, r.EventRecorder, req.NamespacedName, &cloudresourcesv1beta1.VpcPeering{}),
		abstractions.NewFileReader(),
	))

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *NfsShareReconciler) SetupWithManager(mgr ctrl.Manager) error {
	u := &apimachineryapi.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "operator.kyma-project.io",
		Version: "v1beta1",
		Kind:    "Kyma",
	})
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudresourcesv1beta1.NfsShare{}).
		// Kyma CR should be watched on one place only so it gets into the cache
		// we're using empty handler since we're not interested into starting
		// reconciliation when Kyma CR changes, we just want them cached
		Watches(
			u,
			handler.Funcs{},
		).
		Complete(r)
}
