package v1alpha1

type Casting struct {
	TypeVersion `json:",inline" yaml:",inline"`

	// Metadata of the casting configuration.
	Metadata TypeMetadata `json:"metadata" yaml:"metadata" description:"Metadata of the casting configuration"`

	// Specification for the casting.
	Spec CastingSpec `json:"spec" yaml:"spec" description:"Specification for the casting"`

	// Status of the casting.
	Status CastingStatus `json:"status,omitzero" yaml:"status,omitempty" description:"Status of the casting"`
}

type CastingSpec struct {
	// Mode platform in which the platform will run.
	Deployment TypeDeployment `json:"deployment" yaml:"deployment" description:"Deployment configuration for the platform"`

	// The configuration for the o11y molding.
	O11y HanzoO11y `json:"o11y,omitzero" yaml:"o11y,omitempty" description:"The configuration for the HanzoO11y molding"`

	// The configuration for the telemetry store molding.
	TelemetryStore TelemetryStore `json:"telemetrystore,omitzero" yaml:"telemetrystore,omitempty" description:"The configuration for the telemetry store molding"`

	// The configuration for the telemetry keeper molding.
	TelemetryKeeper TelemetryKeeper `json:"telemetrykeeper,omitzero" yaml:"telemetrykeeper,omitempty" description:"The configuration for the telemetry keeper molding"`

	// The configuration for the meta store molding.
	MetaStore MetaStore `json:"metastore,omitzero" yaml:"metastore,omitempty" description:"The configuration for the meta store molding"`

	// The configuration for the ingester molding.
	Ingester Ingester `json:"ingester,omitzero" yaml:"ingester,omitempty" description:"The configuration for the ingester molding"`
}

type CastingStatus struct {
	// Checksum of the casting file.
	Checksum string `json:"checksum" yaml:"checksum" description:"Checksum of the casting file"`
}

func MergeCastingSpecAndStatus(base *Casting) error {
	if err := base.Spec.O11y.Spec.MergeStatus(base.Spec.O11y.Status.MoldingStatus); err != nil {
		return err
	}

	if err := base.Spec.TelemetryStore.Spec.MergeStatus(base.Spec.TelemetryStore.Status.MoldingStatus); err != nil {
		return err
	}

	if err := base.Spec.TelemetryKeeper.Spec.MergeStatus(base.Spec.TelemetryKeeper.Status.MoldingStatus); err != nil {
		return err
	}

	if err := base.Spec.MetaStore.Spec.MergeStatus(base.Spec.MetaStore.Status.MoldingStatus); err != nil {
		return err
	}

	if err := base.Spec.Ingester.Spec.MergeStatus(base.Spec.Ingester.Status.MoldingStatus); err != nil {
		return err
	}

	return nil
}

func DefaultCasting() Casting {
	return Casting{
		TypeVersion: TypeVersion{
			APIVersion: "v1alpha1",
		},
		Metadata: TypeMetadata{
			Name: "o11y",
		},
		Spec: CastingSpec{
			O11y:          DefaultHanzoO11y(),
			TelemetryStore:  DefaultTelemetryStore(),
			TelemetryKeeper: DefaultTelemetryKeeper(),
			MetaStore:       DefaultMetaStore(),
			Ingester:        DefaultIngester(),
		},
	}
}

func ExampleCasting() Casting {
	return Casting{
		TypeVersion: TypeVersion{
			APIVersion: "v1alpha1",
		},
		Metadata: TypeMetadata{
			Name: "o11y",
		},
	}
}
