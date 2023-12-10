package common

import (
	cloudresourcesv1beta1 "github.com/kyma-project/cloud-resources-control-plane/api/cloud-resources/v1beta1"
)

type ProviderType string

const (
	ProviderGCP   = "gcp"
	ProviderAzure = "azure"
	ProviderAws   = "aws"
)

type CommonObject interface {
	Kyma() string

	Scope() *cloudresourcesv1beta1.Scope
	SetScope(scope *cloudresourcesv1beta1.Scope)
}
