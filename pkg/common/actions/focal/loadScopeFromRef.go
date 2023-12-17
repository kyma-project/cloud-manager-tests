package focal

import (
	"context"
	"fmt"
	cloudresourcesv1beta1 "github.com/kyma-project/cloud-resources-control-plane/api/cloud-resources/v1beta1"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/composed"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func loadScopeFromRef(ctx context.Context, state composed.State) (error, context.Context) {
	logger := composed.LoggerFromCtx(ctx)
	logger.Info("Loading scope from reference")

	scope := &cloudresourcesv1beta1.Scope{}
	err := state.Client().Get(ctx, state.Name(), scope)
	if client.IgnoreNotFound(err) != nil {
		err = fmt.Errorf("error getting Scope from reference: %w", err)
		logger.Error(err, "Error loading scope from ref")
		return composed.StopWithRequeue, nil
	}

	logger.Info("Loaded Scope from reference")

	state.(*State).Scope = scope

	return nil, nil
}
