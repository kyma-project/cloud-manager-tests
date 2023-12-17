package actions

import (
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/actions/focal"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/actions/scope"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/composed"
)

func New() composed.Action {
	return composed.ComposeActions(
		"main",
		focal.New(),
		scope.WhenNoScope(),
	)
}
