package foundry

import (
	"log/slog"

	"github.com/signoz/foundry/api/v1alpha1"
	"github.com/signoz/foundry/internal/casting"
	"github.com/signoz/foundry/internal/casting/dockercomposecasting"
	"github.com/signoz/foundry/internal/casting/systemdcasting"
	"github.com/signoz/foundry/internal/casting/terraformcasting"
	"github.com/signoz/foundry/internal/config"
	"github.com/signoz/foundry/internal/config/yamlconfig"
	"github.com/signoz/foundry/internal/molding"
	"github.com/signoz/foundry/internal/molding/ingestermolding"
	"github.com/signoz/foundry/internal/molding/metastoremolding"
	"github.com/signoz/foundry/internal/molding/signozmolding"
	"github.com/signoz/foundry/internal/molding/telemetrykeepermolding"
	"github.com/signoz/foundry/internal/molding/telemetrystoremolding"
	"github.com/signoz/foundry/internal/tooler"
	"github.com/signoz/foundry/internal/tooler/clickhousekeepertooler"
	"github.com/signoz/foundry/internal/tooler/clickhousetooler"
	"github.com/signoz/foundry/internal/tooler/dockercomposetooler"
	"github.com/signoz/foundry/internal/tooler/dockertooler"
	"github.com/signoz/foundry/internal/tooler/postgrestooler"
	"github.com/signoz/foundry/internal/tooler/systemdtooler"
	"github.com/signoz/foundry/internal/tooler/terraformtooler"
)

type Foundry struct {
	// Config for loading the casting configuration.
	Config config.Config

	// Logger for logging.
	Logger *slog.Logger

	// Registry for the different deployments.
	Registry *Registry

	// Moldings for the different molding kinds.
	Moldings map[v1alpha1.MoldingKind]molding.Molding

	// TerraformGenerator for generating infrastructure manifests.
	TerraformGenerator *terraformcasting.TerraformGenerator
}

func New(logger *slog.Logger) (*Foundry, error) {
	yamlConfig := yamlconfig.New()

	registry, err := NewRegistry(logger)
	if err != nil {
		return nil, err
	}

	return &Foundry{
		Config: yamlConfig,
		Logger: logger,
		Castings: map[string]casting.Casting{
			"docker":  dockercomposecasting.New(logger),
			"systemd": systemdcasting.New(logger),
		},
		Toolers: map[string][]tooler.Tooler{
			"terraform": {terraformtooler.New()},
			"docker":    {dockertooler.New(), dockercomposetooler.New()},
			"systemd": {
				systemdtooler.New(),
				clickhousekeepertooler.New(),
				clickhousetooler.New(),
				postgrestooler.New(),
			},
		},
		Moldings: map[v1alpha1.MoldingKind]molding.Molding{
			v1alpha1.MoldingKindTelemetryStore:  telemetrystoremolding.New(logger),
			v1alpha1.MoldingKindTelemetryKeeper: telemetrykeepermolding.New(logger),
			v1alpha1.MoldingKindMetaStore:       metastoremolding.New(logger),
			v1alpha1.MoldingKindSignoz:          signozmolding.New(logger),
			v1alpha1.MoldingKindIngester:        ingestermolding.New(logger),
		},
		TerraformGenerator: terraformcasting.NewGenerator(logger),
	}, nil
}
