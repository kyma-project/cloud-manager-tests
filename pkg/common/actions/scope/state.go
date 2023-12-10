package scope

import (
	gardenerTypes "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	gardenerClient "github.com/gardener/gardener/pkg/client/core/clientset/versioned/typed/core/v1beta1"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/abstractions"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/actions/focal"
	kubernetesClient "k8s.io/client-go/kubernetes"
)

func NewState(base *focal.State, fileReader abstractions.FileReader) *State {
	return &State{
		State:      base,
		FileReader: fileReader,
	}
}

type State struct {
	*focal.State

	FileReader abstractions.FileReader

	ShootName      string
	ShootNamespace string

	GardenerClient  gardenerClient.CoreV1beta1Interface
	GardenK8sClient kubernetesClient.Interface

	Provider       common.ProviderType
	Shoot          *gardenerTypes.Shoot
	CredentialData map[string]string
}
