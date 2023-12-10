package scope

import (
	"context"
	"fmt"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
)

func defineScope(ctx context.Context, st composedAction.State) error {
	state := st.(*State)

	switch state.Provider {
	case common.ProviderGCP:
		return defineScopeGcp(ctx, state)
	case common.ProviderAzure:
		return defineScopeAzure(ctx, state)
	case common.ProviderAws:
		return defineScopeAws(ctx, state)
	}

	err := fmt.Errorf("unable to handle unknown provider '%s'", state.Provider)
	logger := composedAction.LoggerFromCtx(ctx)
	logger.Error(err, "Error defining scope")
	return state.Stop(nil) // no requeue
}
