package ingestermolding

import (
	"embed"

	"github.com/hanzoai/o11y-foundry/internal/types"
)

//go:embed templates/*.gotmpl
var templates embed.FS

var (
	ConfigV0129xTemplate *domain.Template = domain.MustNewTemplateFromFS(templates, "templates/config.v0129x.yaml.gotmpl", domain.FormatYAML)
	OpampV0129xTemplate  *domain.Template = domain.MustNewTemplateFromFS(templates, "templates/opamp.v0129x.yaml.gotmpl", domain.FormatYAML)
)

type Data struct {
	O11yOpampAddress           string
	TelemetryStoreTracesAddress  string
	TelemetryStoreMetricsAddress string
	TelemetryStoreLogsAddress    string
	TelemetryStoreMeterAddress   string
}
