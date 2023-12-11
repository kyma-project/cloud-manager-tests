package scope

import (
	"context"
	"errors"
	"fmt"
	cloudresourcesv1beta1 "github.com/kyma-project/cloud-resources-control-plane/api/cloud-resources/v1beta1"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
)

func defineScopeAzure(ctx context.Context, state *State) (error, context.Context) {
	logger := composedAction.LoggerFromCtx(ctx)

	subscriptionID, ok := state.CredentialData["subscriptionID"]
	if !ok {
		err := errors.New("gardener credential for azure missing subscriptionID key")
		logger.Error(err, "error defining Azure scope")
		return state.Stop(nil), nil // no requeue
	}

	tenantID, ok := state.CredentialData["tenantID"]
	if !ok {
		err := errors.New("gardener credential for azure missing tenantID key")
		logger.Error(err, "error defining Azure scope")
		return state.Stop(nil), nil // no requeue
	}

	scope := &cloudresourcesv1beta1.Scope{
		Azure: &cloudresourcesv1beta1.AzureScope{
			TenantId:       tenantID,
			SubscriptionId: subscriptionID,
			VpcNetwork:     fmt.Sprintf("shoot--%s--%s", state.ShootNamespace, state.ShootName),
		},
	}

	state.Object().SetScope(scope)

	err := state.UpdateObjStatus(ctx)
	if err != nil {
		err = fmt.Errorf("error saving object status with scope: %w", err)
		logger.Error(err, "error saving Azure scope")
		return state.Stop(err), nil // will requeue
	}

	return nil, nil
}
