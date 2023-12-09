package actions

import (
	gardenerTypes "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	gardenerClient "github.com/gardener/gardener/pkg/client/core/clientset/versioned/typed/core/v1beta1"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/abstractions"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
	kubernetesClient "k8s.io/client-go/kubernetes"
)

func NewState(base composedAction.State, fileReader abstractions.FileReader) *State {
	return &State{
		State:      base,
		FileReader: fileReader,
	}
}

type State struct {
	composedAction.State

	FileReader abstractions.FileReader

	ShootName      string
	ShootNamespace string

	GardenerClient  gardenerClient.CoreV1beta1Interface
	GardenK8sClient kubernetesClient.Interface

	Provider       ProviderType
	Shoot          *gardenerTypes.Shoot
	CredentialData map[string]string
}

func (s *State) Object() CommonObject {
	return s.Obj().(CommonObject)
}
