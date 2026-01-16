package telemetrystoremolding

import (
	"embed"
	"fmt"

	"github.com/signoz/foundry/api/v1alpha1"
	"github.com/signoz/foundry/internal/molding"
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

// StoreFunctionsFileName generates the filename for the telemetry store functions file.
// Pattern: {metaName}-{moldingKind}-{storeKind}-functions.yaml.
// Example: signoz-telemetrystore-clickhouse-functions.yaml.
func StoreFunctionsFileName(metaName, kind string) string {
	return molding.FormatFileName([]string{metaName, v1alpha1.MoldingKindTelemetryStore.String(), kind, "functions"}, "yaml")
}

// StoreInstanceConfigFileName generates the filename for a per-instance telemetry store config file.
// Pattern: {metaName}-{moldingKind}-{storeKind}-cluster-{shard}-{replica}.yaml.
// Example: signoz-telemetrystore-clickhouse-cluster-0-1.yaml.
func StoreInstanceConfigFileName(metaName, kind string, shard int, replica int) string {
	return molding.FormatFileName([]string{metaName, v1alpha1.MoldingKindTelemetryStore.String(), kind, fmt.Sprintf("cluster-%d-%d", shard, replica)}, "yaml")
}
