package kuberneteshelmcasting

import (
	"embed"

	"github.com/hanzoai/o11y-foundry/internal/domain"
)

//go:embed templates/values.yaml.gotmpl
var templates embed.FS

var valuesYAMLTemplate *domain.Template = domain.MustNewTemplateFromFS(templates, "templates/values.yaml.gotmpl", domain.FormatYAML)
