package telemetrystoremolding

import (
	"embed"

	"github.com/signoz/foundry/internal/types"
)

//go:embed templates/*.gotmpl
var templates embed.FS

var (
	ConfigClickhousev2556YAML    *types.Template = types.MustNewTemplateFromFS(templates, "templates/config.clickhouse.v2556.yaml.gotmpl", types.FormatYAML)
	FunctionsClickhousev2556YAML *types.Template = types.MustNewTemplateFromFS(templates, "templates/functions.clickhouse.v2556.yaml.gotmpl", types.FormatYAML)
)

// Data is the template data for rendering ClickHouse telemetry store configs.
type Data struct {
	StoreAddresses  []types.Address
	KeeperAddresses []types.Address
	ShardCount      int
	ReplicaCount    int

	// ServerID is the index into StoreAddresses (for per-instance config)
	ServerID int
	// CreatePerInstance indicates if per-instance resources should be created (e.g., numbered paths, instance-specific configs)
	// This is set by the casting's MoldingEnricher when needed by the template
	CreatePerInstance bool
}
