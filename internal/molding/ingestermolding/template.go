package ingestermolding

import (
	"embed"

	"github.com/signoz/foundry/api/v1alpha1"
	"github.com/signoz/foundry/internal/molding"
	"github.com/signoz/foundry/internal/types"
)

//go:embed templates/*.gotmpl
var templates embed.FS

var (
	ConfigV0129xTemplate *types.Template = types.MustNewTemplateFromFS(templates, "templates/config.v0129x.yaml.gotmpl", types.FormatYAML)
	OpampV0129xTemplate  *types.Template = types.MustNewTemplateFromFS(templates, "templates/opamp.v0129x.yaml.gotmpl", types.FormatYAML)
)

type Data struct {
	SignozOpampAddress           string
	TelemetryStoreTracesAddress  string
	TelemetryStoreMetricsAddress string
	TelemetryStoreLogsAddress    string
	TelemetryStoreMeterAddress   string
}

// IngesterConfigFileName generates the filename for the ingester config file.
// Pattern: {moldingKind}-{metaName}-config.yaml.
// Example: ingester-signoz-config.yaml.
func IngesterConfigFileName(metaName, kind string, instance int) string {
	return molding.FormatFileName([]string{v1alpha1.MoldingKindIngester.String(), metaName, "config"}, "yaml")
}

// IngesterOpampFileName generates the filename for the ingester opamp file.
// Pattern: {moldingKind}-{metaName}-opamp.yaml.
// Example: ingester-signoz-opamp.yaml.
func IngesterOpampFileName(metaName, kind string, instance int) string {
	return molding.FormatFileName([]string{v1alpha1.MoldingKindIngester.String(), metaName, "opamp"}, "yaml")
}
