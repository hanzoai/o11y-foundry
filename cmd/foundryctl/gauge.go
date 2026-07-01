package main

import (
	"context"
	"log/slog"

	"github.com/hanzoai/o11y-foundry/internal/domain"
	foundryerrors "github.com/hanzoai/o11y-foundry/internal/errors"
	"github.com/hanzoai/o11y-foundry/internal/foundry"
	"github.com/spf13/cobra"
)

func registerGaugeCmd(rootCmd *cobra.Command) {
	gaugeCmd := &cobra.Command{
		Use:   "gauge",
		Short: "Gauge whether required tools are available.",
		RunE: recoverRunE(domain.EventGauge, func(cmd *cobra.Command, args []string) (domain.Properties, error) {
			return runGauge(cmd.Context(), rootLogger, commonCfg.File)
		}),
	}

	rootCmd.AddCommand(gaugeCmd)
}

func runGauge(ctx context.Context, logger *slog.Logger, path string) (domain.Properties, error) {
	foundry, err := foundry.New(logger)
	if err != nil {
		logger.ErrorContext(ctx, "failed to create foundry, please report this issues to developers at https://github.com/hanzoai/o11y-foundry/issues", foundryerrors.LogAttr(err))
		return domain.NewProperties(), err
	}

	casting, err := foundry.Config.GetV1Alpha1(ctx, path)
	if err != nil {
		return domain.NewProperties(), err
	}

	props := casting.TrackableProperties()

	err = foundry.Gauge(ctx, casting)
	return props, err
}
