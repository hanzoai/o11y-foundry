package dockercomposecasting

import (
	"embed"

	"github.com/hanzoai/o11y-foundry/internal/types"
)

//go:embed templates/*.gotmpl
var templates embed.FS

var (
	composeYAMLTemplate *domain.Template = domain.MustNewTemplateFromFS(templates, "templates/compose.yaml.gotmpl", domain.FormatYAML)
)
