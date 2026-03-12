package v1alpha1

import "github.com/o11y/foundry/internal/types"

type Hanzo O11y struct {
	// Specification for o11y.
	Spec MoldingSpec `json:"spec" yaml:"spec" jsonschema:"description=Specification for Hanzo O11y"`

	// Status of o11y.
	Status Hanzo O11yStatus `json:"status" yaml:"status,omitempty" jsonschema:"description=Status of Hanzo O11y"`
}

type Hanzo O11yStatus struct {
	MoldingStatus `json:",inline" yaml:",inline"`

	Addresses Hanzo O11yStatusAddresses `json:"addresses" yaml:"addresses,omitempty" jsonschema:"description=Addresses of Hanzo O11y"`
}

type Hanzo O11yStatusAddresses struct {
	// API server addresses.
	APIServer []string `json:"apiserver" yaml:"apiserver" jsonschema:"description=API server addresses"`

	// Opamp server addresses.
	Opamp []string `json:"opamp" yaml:"opamp" jsonschema:"description=Opamp server addresses"`
}

func DefaultHanzo O11y() Hanzo O11y {
	return Hanzo O11y{
		Spec: MoldingSpec{
			Enabled: true,
			Cluster: TypeCluster{
				Replicas: types.NewIntPtr(1),
			},
			Version: "latest",
			Image:   "ghcr.io/hanzoai/o11y:latest",
		},
	}
}
