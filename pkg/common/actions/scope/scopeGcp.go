package scope

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	cloudresourcesv1beta1 "github.com/kyma-project/cloud-resources-control-plane/api/cloud-resources/v1beta1"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/composed"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func defineScopeGcp(ctx context.Context, st composed.State) (error, context.Context) {
	logger := composed.LoggerFromCtx(ctx)
	state := st.(*State)

	js, ok := state.CredentialData["serviceaccount.json"]
	if !ok {
		err := errors.New("gardener credential for gcp missing serviceaccount.json key")
		logger.Error(err, "error defining GCP scope")
		return composed.StopAndForget, nil // no requeue
	}

	var data map[string]string
	err := json.Unmarshal([]byte(js), &data)
	if err != nil {
		err := fmt.Errorf("error decoding serviceaccount.json: %w", err)
		logger.Error(err, "error defining GCP scope")
		return composed.StopAndForget, nil // no requeue
	}

	project, ok := data["project_id"]
	if !ok {
		err := errors.New("gardener gcp credentials missing project_id")
		logger.Error(err, "error defining GCP scope")
		return composed.StopAndForget, nil // no requeue
	}

	scope := &cloudresourcesv1beta1.Scope{
		ObjectMeta: metav1.ObjectMeta{
			Name:      state.Obj().GetName(),
			Namespace: state.Obj().GetNamespace(),
		},
		Spec: cloudresourcesv1beta1.ScopeSpec{
			Kyma:      "",
			ShootName: "",
			Scope: cloudresourcesv1beta1.ScopeInfo{
				Gcp: &cloudresourcesv1beta1.GcpScope{
					Project:    project,
					VpcNetwork: fmt.Sprintf("shoot--%s--%s", state.ShootNamespace, state.ShootName),
				},
			},
		},
	}

	err = state.Client().Create(ctx, scope)
	if err != nil {
		err = fmt.Errorf("error creating scope: %w", err)
		logger.Error(err, "error saving GCP scope")
		return composed.StopWithRequeue, nil // will requeue
	}

	state.CommonObj().SetScopeRef(&cloudresourcesv1beta1.ScopeRef{
		Name: scope.Name,
	})

	err = state.UpdateObj(ctx)
	if err != nil {
		err = fmt.Errorf("error updating object scope ref: %w", err)
		logger.Error(err, "error saving object with Gcp scope ref")
		return composed.StopWithRequeue, nil // will requeue
	}

	return composed.StopWithRequeue, nil
}
