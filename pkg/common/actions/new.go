package actions

import (
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/abstractions"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/actions/focal"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/actions/scope"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
)

func New() composedAction.Action {
	return composedAction.ComposeActions(
		"main",
		focal.LoadObj,
		scope.WhenNoScope(abstractions.NewFileReader()),
	)
}
