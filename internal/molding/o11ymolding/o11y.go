package o11ymolding

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/hanzoai/o11y-foundry/api/v1alpha1"
	"github.com/hanzoai/o11y-foundry/internal/molding"
	"github.com/hanzoai/o11y-foundry/internal/types"
)

var _ molding.Molding = (*o11y)(nil)

type o11y struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *o11y {
	return &o11y{
		logger: logger,
	}
}

func (molding *o11y) Kind() v1alpha1.MoldingKind {
	return v1alpha1.MoldingKindO11y
}

func (molding *o11y) MoldV1Alpha1(ctx context.Context, config *v1alpha1.Casting) error {
	if config.Spec.O11y.Status.Env == nil {
		config.Spec.O11y.Status.Env = make(map[string]string)
	}

	if config.Spec.O11y.Spec.Env == nil {
		config.Spec.O11y.Spec.Env = make(map[string]string)
	}

	// Add telemetry store addresses
	config.Spec.O11y.Status.Env["HANZO_TELEMETRYSTORE_PROVIDER"] = config.Spec.TelemetryStore.Kind.String()

	if val, ok := config.Spec.O11y.Spec.Env["HANZO_TELEMETRYSTORE_CLICKHOUSE_DSN"]; ok {
		molding.logger.WarnContext(ctx, "HANZO_TELEMETRYSTORE_CLICKHOUSE_DSN is going to be overridden", slog.String("value", val))
	}

	config.Spec.O11y.Status.Env["HANZO_TELEMETRYSTORE_CLICKHOUSE_DSN"] = strings.Join(config.Spec.TelemetryStore.Status.Addresses.TCP, ",")

	// Add metastore addresses
	config.Spec.O11y.Status.Env["HANZO_SQLSTORE_PROVIDER"] = config.Spec.MetaStore.Kind.String()

	if config.Spec.MetaStore.Status.Addresses.DSN != nil {
		if val, ok := config.Spec.O11y.Spec.Env["HANZO_SQLSTORE_POSTGRES_DSN"]; ok {
			molding.logger.WarnContext(ctx, "HANZO_SQLSTORE_POSTGRES_DSN is going to be overridden", slog.String("value", val))
		}
		// construct postgres dsn with user, password, host, port, and db
		addrs, err := types.NewAddresses(config.Spec.MetaStore.Status.Addresses.DSN)
		if err != nil {
			return fmt.Errorf("failed to parse addresses: %w", err)
		}
		var dsns []string
		user := config.Spec.MetaStore.Status.Env["POSTGRES_USER"]
		password := config.Spec.MetaStore.Status.Env["POSTGRES_PASSWORD"]
		db := config.Spec.MetaStore.Status.Env["POSTGRES_DB"]
		for _, addr := range addrs {
			dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, addr.Host(), addr.Port(), db)
			dsns = append(dsns, dsn)
		}
		config.Spec.O11y.Status.Env["HANZO_SQLSTORE_POSTGRES_DSN"] = strings.Join(dsns, ",")
	}
	return nil
}
