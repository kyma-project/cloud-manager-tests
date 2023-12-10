package focal

import (
	"context"
	"fmt"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
)

func LoadObj(ctx context.Context, state composedAction.State) error {
	logger := composedAction.LoggerFromCtx(ctx)
	err := state.LoadObj(ctx)
	if err != nil {
		err = fmt.Errorf("error getting object: %w", err)
		logger.Error(err, "error")
		return state.RequeueIfError(err)
	}

	logger.Info("Object loaded")

	return nil
}
