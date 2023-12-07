package actions

import (
	"context"
	composed "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
)

func LoadKyma(ctx context.Context, state composed.State) error {
	// not sure if it's possible to use client.Client to get unstructured or
	// we must have a typed object, either as local mock or remote real dependency to LM
	state.(*State).ShootName = "something.read.from.kyma.cr"
	return nil
}
