package molding

import (
	"context"

	"github.com/hanzoai/o11y-foundry/api/v1alpha1"
)

// MoldingEnricher populates a molding's Status fields from the surrounding
// installation casting. The ordering of EnrichStatus calls is owned by the
// installation Planner, which iterates the kinds it knows about.
type MoldingEnricher interface {
	EnrichStatus(ctx context.Context, kind v1alpha1.MoldingKind, config *installation.Casting) error
}

// Molding generates materials for a single SigNoz component. Mutates the
// config in place; not safe for concurrent use.
type Molding interface {
	Kind() v1alpha1.MoldingKind

	// Molds the v1alpha1 casting configuration. This function mutates the config in place. It is not safe for concurrent use.
	MoldV1Alpha1(ctx context.Context, config *v1alpha1.Casting) error
}

func MoldingsInOrder() []v1alpha1.MoldingKind {
	return []v1alpha1.MoldingKind{
		v1alpha1.MoldingKindTelemetryKeeper,
		v1alpha1.MoldingKindTelemetryStore,
		v1alpha1.MoldingKindMetaStore,
		v1alpha1.MoldingKindO11y,
		v1alpha1.MoldingKindIngester,
	}
}
