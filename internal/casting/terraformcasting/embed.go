package terraformcasting

import (
	"embed"

	"github.com/signoz/foundry/internal/types"
)

//go:embed templates/*.gotmpl templates/aws/*.gotmpl templates/gcp/*.gotmpl templates/azure/*.gotmpl
var templates embed.FS

// Common templates.
var (
	providersTFTemplate *types.Template = types.MustNewTemplateFromFS(templates, "templates/providers.tf.gotmpl", types.FormatHCL)
)

// AWS templates.
var (
	awsMainTFTemplate      *types.Template = types.MustNewTemplateFromFS(templates, "templates/aws/main.tf.gotmpl", types.FormatHCL)
	awsVariablesTFTemplate *types.Template = types.MustNewTemplateFromFS(templates, "templates/aws/variables.tf.gotmpl", types.FormatHCL)
	awsOutputsTFTemplate   *types.Template = types.MustNewTemplateFromFS(templates, "templates/aws/outputs.tf.gotmpl", types.FormatHCL)
)

// GCP templates.
var (
	gcpMainTFTemplate      *types.Template = types.MustNewTemplateFromFS(templates, "templates/gcp/main.tf.gotmpl", types.FormatHCL)
	gcpVariablesTFTemplate *types.Template = types.MustNewTemplateFromFS(templates, "templates/gcp/variables.tf.gotmpl", types.FormatHCL)
	gcpOutputsTFTemplate   *types.Template = types.MustNewTemplateFromFS(templates, "templates/gcp/outputs.tf.gotmpl", types.FormatHCL)
)

// Azure templates.
var (
	azureMainTFTemplate      *types.Template = types.MustNewTemplateFromFS(templates, "templates/azure/main.tf.gotmpl", types.FormatHCL)
	azureVariablesTFTemplate *types.Template = types.MustNewTemplateFromFS(templates, "templates/azure/variables.tf.gotmpl", types.FormatHCL)
	azureOutputsTFTemplate   *types.Template = types.MustNewTemplateFromFS(templates, "templates/azure/outputs.tf.gotmpl", types.FormatHCL)
)
