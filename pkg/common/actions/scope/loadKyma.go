package scope

import (
	"context"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
	apimachineryapi "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func loadKyma(ctx context.Context, state composedAction.State) (error, context.Context) {
	logger := composedAction.LoggerFromCtx(ctx)

	u := &apimachineryapi.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "operator.kyma-project.io",
		Version: "v1beta1",
		Kind:    "Kyma",
	})
	err := state.Client().Get(ctx, types.NamespacedName{
		Namespace: state.Obj().GetNamespace(),
		Name:      state.(*State).Object().Kyma(),
	}, u)
	if err != nil {
		logger.Error(err, "error loading Kyma CR")
		return err, nil
	}

	state.(*State).ShootName = u.GetLabels()["kyma-project.io/shoot-name"]

	logger = logger.WithValues("shootName", state.(*State).ShootName)
	logger.Info("Shoot name found")

	return nil, composedAction.LoggerIntoCtx(ctx, logger)
}
