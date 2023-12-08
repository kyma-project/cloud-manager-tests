package v1beta1

type ProviderType string

const (
	ProviderGCP   = "gcp"
	ProviderAzure = "azure"
	ProviderAws   = "aws"
)

type StatusState string

const (
	UnknownState StatusState = "Unknown"
	ReadyState   StatusState = "Ready"
	ErrorState   StatusState = "Error"
)
