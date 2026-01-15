package molding

import (
	"context"
	"fmt"
	"strings"

	"github.com/signoz/foundry/api/v1alpha1"
)

type MoldingEnricher interface {
	// Enrich the molding status with the casting configuration.
	EnrichStatus(ctx context.Context, kind v1alpha1.MoldingKind, config *v1alpha1.Casting) error
}

type Molding interface {
	// Kind of the molding.
	Kind() v1alpha1.MoldingKind

	// Molds the v1alpha1 casting configuration. This function mutates the config in place. It is not safe for concurrent use.
	MoldV1Alpha1(ctx context.Context, config *v1alpha1.Casting) error
}

func MoldingsInOrder() []v1alpha1.MoldingKind {
	return []v1alpha1.MoldingKind{
		v1alpha1.MoldingKindTelemetryKeeper,
		v1alpha1.MoldingKindTelemetryStore,
		v1alpha1.MoldingKindMetaStore,
		v1alpha1.MoldingKindSignoz,
		v1alpha1.MoldingKindIngester,
	}
}

// FormatFileName generates a standardized filename from parts joined
// with the format as the extension. This follows the naming convention pattern:
//
//	{metadataName}-{moldingKind}-{subKind}-{suffix}-{instanceInfo}.{format}
//
// Example: signoz-ingester-config.yaml
//
// Each molding should implement naming functions in its template.go file that
// call this utility with the appropriate parts for their specific file types.
// This allows the core molding system to remain agnostic of specific naming
// patterns while ensuring consistency across all moldings.
func FormatFileName(parts []string, format string) string {
	if len(parts) == 0 {
		return format
	}
	return fmt.Sprintf("%s.%s", strings.Join(parts, "-"), format)
}
