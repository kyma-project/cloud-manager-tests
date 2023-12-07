package actions

import (
	gardenerapiv1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	composedAction "github.com/kyma-project/cloud-resources-control-plane/pkg/common/composedAction"
)

type ProviderType string

const (
	ProviderGCP   = "gcp"
	ProviderAzure = "azure"
	ProviderAws   = "aws"
)

type State struct {
	composedAction.State

	ShootName      string
	Provider       ProviderType
	Shoot          *gardenerapiv1beta1.Shoot
	CredentialData map[string]string
}
