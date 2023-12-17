package scope

import (
	"context"
	"errors"
	"fmt"
	cloudresourcesv1beta1 "github.com/kyma-project/cloud-resources-control-plane/api/cloud-resources/v1beta1"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/composed"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func defineScopeAzure(ctx context.Context, st composed.State) (error, context.Context) {
	logger := composed.LoggerFromCtx(ctx)
	state := st.(*State)

	subscriptionID, ok := state.CredentialData["subscriptionID"]
	if !ok {
		err := errors.New("gardener credential for azure missing subscriptionID key")
		logger.Error(err, "error defining Azure scope")
		return composed.StopAndForget, nil // no requeue
	}

	tenantID, ok := state.CredentialData["tenantID"]
	if !ok {
		err := errors.New("gardener credential for azure missing tenantID key")
		logger.Error(err, "error defining Azure scope")
		return composed.StopAndForget, nil // no requeue
	}

	scope := &cloudresourcesv1beta1.Scope{
		ObjectMeta: metav1.ObjectMeta{
			Name:      state.Obj().GetName(),
			Namespace: state.Obj().GetNamespace(),
			Labels: map[string]string{
				cloudresourcesv1beta1.ScopeKymaLabel: state.CommonObj().KymaName(),
			},
		},
		Spec: cloudresourcesv1beta1.ScopeSpec{
			Kyma:      "",
			ShootName: "",
			Scope: cloudresourcesv1beta1.ScopeInfo{
				Azure: &cloudresourcesv1beta1.AzureScope{
					TenantId:       tenantID,
					SubscriptionId: subscriptionID,
					VpcNetwork:     fmt.Sprintf("shoot--%s--%s", state.ShootNamespace, state.ShootName),
				},
			},
		},
	}

	err := state.Client().Create(ctx, scope)
	if err != nil {
		err = fmt.Errorf("error creating scope: %w", err)
		logger.Error(err, "error saving Azure scope")
		return composed.StopWithRequeue, nil // will requeue
	}

	state.CommonObj().SetScopeRef(&cloudresourcesv1beta1.ScopeRef{
		Name: scope.Name,
	})

	err = state.UpdateObj(ctx)
	if err != nil {
		err = fmt.Errorf("error updating object scope ref: %w", err)
		logger.Error(err, "error saving object with Azure scope ref")
		return composed.StopWithRequeue, nil // will requeue
	}

	// scope ref is set, can retry new
	return composed.StopWithRequeue, nil
}
