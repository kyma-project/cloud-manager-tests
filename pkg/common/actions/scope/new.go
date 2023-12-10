package scope

import composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"

func New() composedAction.Action {
	return composedAction.ComposeActions(
		"whenNoScope",
		loadKyma,
		createGardenerClient,
		loadShoot,
		loadGardenerCredentials,
		defineScope,
	)
}
