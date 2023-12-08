package actions

import (
	gardenerapiv1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cloudresourcesv1beta1 "github.com/kyma-project/cloud-resources-control-plane/api/cloud-resources/v1beta1"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
)

func NewState(base composedAction.State) *State {
	return &State{
		State: base,
	}
}

type State struct {
	composedAction.State

	ShootName      string
	Provider       cloudresourcesv1beta1.ProviderType
	Shoot          *gardenerapiv1beta1.Shoot
	CredentialData map[string]string
}

func (s *State) Object() CommonObject {
	return s.Obj().(CommonObject)
}
