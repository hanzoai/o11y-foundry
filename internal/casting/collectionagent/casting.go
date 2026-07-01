package collectionagent

import (
	"context"

	"github.com/hanzoai/o11y-foundry/api/v1alpha1/collectionagent"
	"github.com/hanzoai/o11y-foundry/internal/domain"
	collectionagentmolding "github.com/hanzoai/o11y-foundry/internal/molding/collectionagent"
)

type Casting interface {
	Enricher(ctx context.Context, config *collectionagent.Casting) (collectionagentmolding.MoldingEnricher, error)
	Forge(ctx context.Context, config collectionagent.Casting, poursPath string) ([]domain.Material, error)
	Cast(ctx context.Context, config collectionagent.Casting, poursPath string) error
}
