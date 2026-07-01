package foundry

import (
	"context"
	"log/slog"

	"github.com/hanzoai/o11y-foundry/api/v1alpha1"
	"github.com/hanzoai/o11y-foundry/api/v1alpha1/collectionagent"
	"github.com/hanzoai/o11y-foundry/api/v1alpha1/installation"
	collectionagentcasting "github.com/hanzoai/o11y-foundry/internal/casting/collectionagent"
	installationcasting "github.com/hanzoai/o11y-foundry/internal/casting/installation"
	"github.com/hanzoai/o11y-foundry/internal/config"
	"github.com/hanzoai/o11y-foundry/internal/config/yamlconfig"
	foundryerrors "github.com/hanzoai/o11y-foundry/internal/errors"
	"github.com/hanzoai/o11y-foundry/internal/infrastructure"
	terraformgenerator "github.com/hanzoai/o11y-foundry/internal/infrastructure/terraform"
	"github.com/hanzoai/o11y-foundry/internal/patch"
	"github.com/hanzoai/o11y-foundry/internal/planner"
)

type plannerCtor func(ctx context.Context, m v1alpha1.Machinery, logger *slog.Logger) (planner.Planner, error)

type Foundry struct {
	// Config for loading the casting configuration.
	Config config.Config

	// Patchers for applying patches to generated materials, keyed by patch type.
	Patchers map[string]patch.Patch

	// Logger for logging.
	Logger *slog.Logger

	// Planners for the different casting kinds.
	Planners map[v1alpha1.Kind]plannerCtor

	// InfrastructureGenerator for generating infrastructure-as-code manifests.
	InfrastructureGenerator infrastructure.Generator
}

func New(logger *slog.Logger) (*Foundry, error) {
	return &Foundry{
		Config: yamlconfig.New(),
		Logger: logger,
		Planners: map[v1alpha1.Kind]plannerCtor{
			v1alpha1.KindInstallation: func(ctx context.Context, m v1alpha1.Machinery, logger *slog.Logger) (planner.Planner, error) {
				return installationcasting.NewPlanner(ctx, m.(*installation.Casting), logger)
			},
			v1alpha1.KindCollectionAgent: func(ctx context.Context, m v1alpha1.Machinery, logger *slog.Logger) (planner.Planner, error) {
				return collectionagentcasting.NewPlanner(ctx, m.(*collectionagent.Casting), logger)
			},
		},
		InfrastructureGenerator: terraformgenerator.New(logger),
	}, nil
}

func (foundry *Foundry) newPlanner(ctx context.Context, m v1alpha1.Machinery) (planner.Planner, error) {
	ctor, ok := foundry.Planners[m.Kind()]
	if !ok {
		return nil, foundryerrors.Newf(foundryerrors.TypeUnsupported, "unsupported casting kind %q", m.Kind())
	}
	return ctor(ctx, m, foundry.Logger)
}
