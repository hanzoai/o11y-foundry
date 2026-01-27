package types

var (
	FormatYAML Format = Format{s: "yaml"}
	FormatJSON Format = Format{s: "json"}
	FormatINI  Format = Format{s: "ini"}
	FormatHCL  Format = Format{s: "hcl"}
	FormatText Format = Format{s: "text"}
)

type Format struct{ s string }
