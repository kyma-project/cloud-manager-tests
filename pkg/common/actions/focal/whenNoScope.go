package focal

import (
	"context"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/abstractions"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/actions/scope"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
)

func WhenNoScope(fileReader abstractions.FileReader) composedAction.Action {
	return func(ctx context.Context, st composedAction.State) error {
		state := st.(*State)
		if state.Object().Scope() != nil {
			return nil // continue
		}

		logger := composedAction.LoggerFromCtx(ctx)
		logger.Info("Object has no scope, running define scope branch")

		scopeState := scope.NewState(state, fileReader)
		action := scope.New()
		return action(ctx, scopeState)
	}
}
