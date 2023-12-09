package actions

import (
	"context"
	"fmt"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func LoadShoot(ctx context.Context, st composedAction.State) error {
	logger := composedAction.LoggerFromCtx(ctx)
	state := st.(*State)

	shoot, err := state.GardenerClient.Shoots(state.ShootNamespace).Get(ctx, state.ShootName, metav1.GetOptions{})
	if err != nil {
		err = fmt.Errorf("error getting shoot: %w", err)
		logger.Error(err, "Error loading shoot")
		return err
	}

	state.Shoot = shoot

	logger.Info("Shoot loaded")

	return nil
}
