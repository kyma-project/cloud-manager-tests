package focal

import "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composed"

func New() composed.Action {
	return composed.ComposeActions(
		"focal",
		loadObj,
		loadScopeFromRef,
		fixInvalidScopeRef,
	)
}
