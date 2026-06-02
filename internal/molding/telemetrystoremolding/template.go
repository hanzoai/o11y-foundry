package telemetrystoremolding

import (
	"embed"

	"github.com/hanzoai/o11y-foundry/internal/types"
)

//go:embed templates/*.gotmpl
var templates embed.FS

var (
	ConfigClickhousev2556YAML    *domain.Template = domain.MustNewTemplateFromFS(templates, "templates/config.clickhouse.v2556.yaml.gotmpl", domain.FormatYAML)
	FunctionsClickhousev2556YAML *domain.Template = domain.MustNewTemplateFromFS(templates, "templates/functions.clickhouse.v2556.yaml.gotmpl", domain.FormatYAML)
)

// Data is the template data for rendering ClickHouse telemetry store configs.
type Data struct {
	StoreAddresses  []domain.Address
	KeeperAddresses []domain.Address
	ShardCount      int
	ReplicaCount    int
	ShardID         int // 0-indexed, used to render per-node macros.shard
	ReplicaID       int // 0-indexed, used to render per-node macros.replica
}
