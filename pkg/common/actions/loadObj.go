package actions

import (
	"context"
	"fmt"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
)

func LoadObj(ctx context.Context, state composedAction.State) error {
	err := state.LoadObj(ctx)
	return state.RequeueIfError(err, fmt.Sprintf("error getting object %s", state.Name()))
}
