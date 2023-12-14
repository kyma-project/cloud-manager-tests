package scope

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	cloudresourcesv1beta1 "github.com/kyma-project/cloud-resources-control-plane/api/cloud-resources/v1beta1"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
)

func defineScopeGcp(ctx context.Context, state *State) (error, context.Context) {
	logger := composedAction.LoggerFromCtx(ctx)
	js, ok := state.CredentialData["serviceaccount.json"]
	if !ok {
		err := errors.New("gardener credential for gcp missing serviceaccount.json key")
		logger.Error(err, "error defining GCP scope")
		return state.Stop(nil), nil // no requeue
	}

	var data map[string]string
	err := json.Unmarshal([]byte(js), &data)
	if err != nil {
		err := fmt.Errorf("error decoding serviceaccount.json: %w", err)
		logger.Error(err, "error defining GCP scope")
		return state.Stop(nil), nil // no requeue
	}

	project, ok := data["project_id"]
	if !ok {
		err := errors.New("gardener gcp credentials missing project_id")
		logger.Error(err, "error defining GCP scope")
		return state.Stop(nil), nil // no requeue
	}

	scope := &cloudresourcesv1beta1.ScopeX{
		Gcp: &cloudresourcesv1beta1.GcpScope{
			Project:    project,
			VpcNetwork: fmt.Sprintf("shoot--%s--%s", state.ShootNamespace, state.ShootName),
		},
	}

	state.Object().SetScope(scope)

	err = state.UpdateObjStatus(ctx)
	if err != nil {
		err = fmt.Errorf("error saving object status with scope: %w", err)
		logger.Error(err, "error saving GCP scope")
		return state.Stop(err), nil // will requeue
	}

	return nil, nil
}
