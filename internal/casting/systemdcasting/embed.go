package systemdcasting

import (
	"embed"

	"github.com/hanzoai/o11y-foundry/internal/types"
)

//go:embed templates/*.gotmpl
var templates embed.FS

var (
	o11yServiceTemplate                 *types.Template = types.MustNewTemplateFromFS(templates, "templates/o11y.service.gotmpl", types.FormatINI)
	ingesterServiceTemplate               *types.Template = types.MustNewTemplateFromFS(templates, "templates/ingester.service.gotmpl", types.FormatINI)
	telemetryStoreServiceTemplate         *types.Template = types.MustNewTemplateFromFS(templates, "templates/clickhouse.telemetrystore.v2556.service.gotmpl", types.FormatINI)
	telemetryKeeperServiceTemplate        *types.Template = types.MustNewTemplateFromFS(templates, "templates/clickhousekeeper.telemetrykeeper.v2556.service.gotmpl", types.FormatINI)
	metaStoreServiceTemplate              *types.Template = types.MustNewTemplateFromFS(templates, "templates/postgres.metastore.service.gotmpl", types.FormatINI)
	telemetryStoreMigratorServiceTemplate *types.Template = types.MustNewTemplateFromFS(templates, "templates/migrator.telemetrystore.service.gotmpl", types.FormatINI)
)
