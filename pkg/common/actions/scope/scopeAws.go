package scope

import (
	"context"
	"errors"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
)

func defineScopeAws(ctx context.Context, state *State) error {
	logger := composedAction.LoggerFromCtx(ctx)
	err := errors.New("aws scope definition not implemented")
	logger.Error(err, "error defining AWS scope")

	return state.Stop(nil) // no requeue
}
