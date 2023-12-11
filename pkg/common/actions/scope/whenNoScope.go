package scope

import (
	"context"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/abstractions"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/actions/focal"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
)

func WhenNoScope(fileReader abstractions.FileReader) composedAction.Action {
	return func(ctx context.Context, st composedAction.State) (error, context.Context) {
		state := st.(*focal.State)
		if state.Object().Scope() != nil {
			logger := composedAction.LoggerFromCtx(ctx)
			logger.Info("Object has scope")
			return nil, nil // continue
		}

		logger := composedAction.LoggerFromCtx(ctx)
		logger.Info("Object has no scope, running define scope branch")

		scopeState := NewState(state, fileReader)
		action := New()
		return action(ctx, scopeState)
	}
}
