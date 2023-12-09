package actions

type ProviderType string

const (
	ProviderGCP   = "gcp"
	ProviderAzure = "azure"
	ProviderAws   = "aws"
)

type CommonObject interface {
	Kyma() string
}
