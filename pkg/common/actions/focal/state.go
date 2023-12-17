package focal

import (
	cloudresourcesv1beta1 "github.com/kyma-project/cloud-resources-control-plane/api/cloud-resources/v1beta1"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/composed"
)

func NewState(base composed.State) *State {
	return &State{
		State: base,
	}
}

type State struct {
	composed.State

	Scope *cloudresourcesv1beta1.Scope
}

func (s *State) CommonObj() CommonObject {
	return s.Obj().(CommonObject)
}
