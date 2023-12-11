package focal

import (
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
)

func NewState(base composedAction.State) *State {
	return &State{
		State: base,
	}
}

type State struct {
	composedAction.State
}

func (s *State) Object() common.CommonObject {
	return s.Obj().(common.CommonObject)
}
