package actions

import (
	"context"
	"errors"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
)

var gardenerProviderExtensionsMap map[ProviderType]string

func init() {
	gardenerProviderExtensionsMap = map[ProviderType]string{
		ProviderGCP:   "provider.extensions.gardener.cloud/gcp",
		ProviderAzure: "provider.extensions.gardener.cloud/azure",
		ProviderAws:   "provider.extensions.gardener.cloud/aws",
	}
}

func DetectProvider(ctx context.Context, st composedAction.State) error {
	logger := composedAction.LoggerFromCtx(ctx)
	state := st.(*State)

	for providerType, label := range gardenerProviderExtensionsMap {
		val, exists := state.Shoot.Labels[label]
		if exists && (val == "true" || val == "True") {
			logger = logger.WithValues("provider", providerType)
			logger.Info("Detected provider from shoot annotations")
			state.Provider = providerType
			composedAction.LoggerIntoCtx(ctx, logger, state)
			return nil
		}
	}

	logger.Error(errors.New("unable to detect provider"), "no known provider found in shoot annotations")

	return state.Stop(nil) // no requeue
}
