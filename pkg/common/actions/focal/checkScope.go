package focal

import (
	"context"
	"fmt"
	cloudresourcesv1beta1 "github.com/kyma-project/cloud-resources-control-plane/api/cloud-resources/v1beta1"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/composed"
)

func checkScope(ctx context.Context, st composed.State) (error, context.Context) {
	logger := composed.LoggerFromCtx(ctx)
	state := st.(*State)

	if state.Scope == nil {
		return nil, nil // whenNoScope will handle this, create the Scope and requeue
	}

	if state.CommonObj().ScopeRef() != nil &&
		state.CommonObj().ScopeRef().Name == state.Scope.Name &&
		state.CommonObj().GetNamespace() == state.Scope.Namespace {
		return nil, nil // all fine, Scope exists and object has reference to it
	}

	// set scope reference to the object, update it and requeue
	state.CommonObj().SetScopeRef(&cloudresourcesv1beta1.ScopeRef{Name: state.Scope.Name})
	err := state.UpdateObj(ctx)
	if err != nil {
		err = fmt.Errorf("error updating object with scope reference: %w", err)
		logger.Error(err, "Error checking scope")
		return composed.StopWithRequeue, nil
	}

	logger.Info("Scope reference is saved to the object")

	return nil, nil
}
