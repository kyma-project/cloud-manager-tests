package scope

import (
	"context"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/composed"
)

func WhenNoScope() composed.Action {
	return composed.ComposeActions(
		"whenNoScope",
		loadKyma,
		createGardenerClient,
		loadShoot,
		loadGardenerCredentials,
		handleScope,
		// just in case actions before didn't stopped it
		// scope is created, requeue now
		func(_ context.Context, state composed.State) (error, context.Context) {
			return composed.StopWithRequeue, nil
		},
	)
}
