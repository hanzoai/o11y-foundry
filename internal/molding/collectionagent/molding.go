package collectionagent

import (
	"context"

	"github.com/hanzoai/o11y-foundry/api/v1alpha1"
	"github.com/hanzoai/o11y-foundry/api/v1alpha1/collectionagent"
)

type MoldingEnricher interface {
	EnrichStatus(ctx context.Context, kind v1alpha1.MoldingKind, config *collectionagent.Casting) error
}

type Molding interface {
	Kind() v1alpha1.MoldingKind
	MoldV1Alpha1(ctx context.Context, config *collectionagent.Casting) error
}
