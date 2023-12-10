package actions

import (
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/abstractions"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/actions/focal"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
)

func New() composedAction.Action {
	return composedAction.ComposeActions(
		"main",
		focal.LoadObj,
		focal.WhenNoScope(abstractions.NewFileReader()),
	)
}
