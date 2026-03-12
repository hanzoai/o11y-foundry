package v1alpha1

import "github.com/hanzoai/o11y-foundry/internal/types"

type HanzoO11y struct {
	// Specification for o11y.
	Spec MoldingSpec `json:"spec" yaml:"spec" jsonschema:"description=Specification for HanzoO11y"`

	// Status of o11y.
	Status HanzoO11yStatus `json:"status" yaml:"status,omitempty" jsonschema:"description=Status of HanzoO11y"`
}

type HanzoO11yStatus struct {
	MoldingStatus `json:",inline" yaml:",inline"`

	Addresses HanzoO11yStatusAddresses `json:"addresses" yaml:"addresses,omitempty" jsonschema:"description=Addresses of HanzoO11y"`
}

type HanzoO11yStatusAddresses struct {
	// API server addresses.
	APIServer []string `json:"apiserver" yaml:"apiserver" jsonschema:"description=API server addresses"`

	// Opamp server addresses.
	Opamp []string `json:"opamp" yaml:"opamp" jsonschema:"description=Opamp server addresses"`
}

func DefaultHanzoO11y() HanzoO11y {
	return HanzoO11y{
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
