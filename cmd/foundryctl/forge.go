package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/hanzoai/o11y-foundry/internal/domain"
	foundryerrors "github.com/hanzoai/o11y-foundry/internal/errors"
	"github.com/hanzoai/o11y-foundry/internal/foundry"
	"github.com/hanzoai/o11y-foundry/internal/writer"
	"github.com/spf13/cobra"
)

func registerForgeCmd(rootCmd *cobra.Command) {
	forgeCmd := &cobra.Command{
		Use:   "forge",
		Short: "Forge Configuration and Deployment Files",
		Long:  "Generate deployment configuration files from casting.yaml",
		RunE: recoverRunE(domain.EventForge, func(cmd *cobra.Command, args []string) (domain.Properties, error) {
			return runForge(cmd.Context(), rootLogger, commonCfg.File, poursCfg.Path)
		}),
	}

	rootCmd.AddCommand(forgeCmd)
}

func runForge(ctx context.Context, logger *slog.Logger, path string, poursPath string) (domain.Properties, error) {
	foundry, err := foundry.New(logger)
	if err != nil {
		logger.ErrorContext(ctx, "failed to create foundry, please report this issues to developers at https://github.com/hanzoai/o11y-foundry/issues", foundryerrors.LogAttr(err))
		return domain.NewProperties(), err
	}

	machinery, err := foundry.Config.GetV1Alpha1(ctx, path)
	if err != nil {
		return domain.NewProperties(), err
	}

	props := machinery.TrackableProperties()

	poursAbsPath, err := filepath.Abs(poursPath)
	if err != nil {
		return props, err
	}

	err = foundry.Forge(ctx, machinery, path, &writer.Options{Output: &os.File{}, TargetDirectory: poursAbsPath})
	return props, err
}
