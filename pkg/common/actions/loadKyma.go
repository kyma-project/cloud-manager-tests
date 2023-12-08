package actions

import (
	"context"
	composed "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
	apimachineryapi "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func LoadKyma(ctx context.Context, state composed.State) error {
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
		return err
	}

	state.(*State).ShootName = u.GetLabels()["kyma-project.io/shoot-name"]

	return nil
}
