package scope

import (
	"context"
	"fmt"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/composed"
)

func handleScope(ctx context.Context, st composed.State) (error, context.Context) {
	state := st.(*State)
	switch state.Provider {
	case ProviderGCP:
		return defineScopeGcp(ctx, state)
	case ProviderAzure:
		return defineScopeAzure(ctx, state)
	case ProviderAws:
		return defineScopeAws(ctx, state)
	}

	err := fmt.Errorf("unable to handle unknown provider '%s'", state.Provider)
	logger := composed.LoggerFromCtx(ctx)
	logger.Error(err, "Error defining scope")
	return composed.StopAndForget, nil // no requeue

}
