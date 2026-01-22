package v1alpha1

// InfrastructureProvider represents the cloud provider for infrastructure deployment.
type InfrastructureProvider string

const (
	InfrastructureProviderAWS   InfrastructureProvider = "aws"
	InfrastructureProviderGCP   InfrastructureProvider = "gcp"
	InfrastructureProviderAzure InfrastructureProvider = "azure"
)

// Infrastructure holds the configuration for infrastructure manifest generation (e.g., Terraform).
type Infrastructure struct {
	// Whether infrastructure manifest generation is enabled
	Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`

	// The cloud provider to generate infrastructure manifests for (aws, gcp, azure)
	Provider InfrastructureProvider `json:"provider,omitempty" yaml:"provider,omitempty"`
}

// DefaultInfrastructure returns the default Infrastructure configuration.
func DefaultInfrastructure() Infrastructure {
	return Infrastructure{
		Enabled:  false,
		Provider: InfrastructureProviderAWS,
	}
}
