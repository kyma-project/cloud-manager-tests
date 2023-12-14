package v1beta1

type ScopeX struct {
	// +optional
	Gcp *GcpScope `json:"gcp,omitempty"`

	// +optional
	Azure *AzureScope `json:"azure,omitempty"`

	// +optional
	Aws *AwsScope `json:"aws,omitempty"`
}

type GcpScope struct {
	// +kubebuilder:validation:Required
	Project string `json:"project"`

	// +kubebuilder:validation:Required
	VpcNetwork string `json:"vpcNetwork"`
}

type AzureScope struct {
	// +kubebuilder:validation:Required
	TenantId string `json:"tenantId"`

	// +kubebuilder:validation:Required
	SubscriptionId string `json:"subscriptionId"`

	// +kubebuilder:validation:Required
	VpcNetwork string `json:"vpcNetwork"`
}

type AwsScope struct {
	// +kubebuilder:validation:Required
	Foo string `json:"foo"`
}
