package terraformcasting

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/signoz/foundry/api/v1alpha1"
	"github.com/signoz/foundry/internal/types"
)

// TerraformGenerator generates Terraform manifests for infrastructure deployment.
type TerraformGenerator struct {
	logger *slog.Logger
}

// NewGenerator creates a new TerraformGenerator.
func NewGenerator(logger *slog.Logger) *TerraformGenerator {
	return &TerraformGenerator{
		logger: logger,
	}
}

// Generate creates Terraform manifests based on the casting configuration and infrastructure provider.
func (g *TerraformGenerator) Generate(ctx context.Context, config v1alpha1.Casting) ([]types.Material, error) {
	if !config.Spec.Infrastructure.Enabled {
		return nil, nil
	}

	provider := config.Spec.Infrastructure.Provider
	g.logger.InfoContext(ctx, "Generating Terraform manifests", slog.String("provider", string(provider)))

	var materials []types.Material

	// Get provider-specific templates
	mainTemplate, varsTemplate, outputsTemplate, err := g.getTemplatesForProvider(provider)
	if err != nil {
		return nil, err
	}

	// Generate main.tf
	mainBuf := bytes.NewBuffer(nil)
	if err := mainTemplate.Execute(mainBuf, config); err != nil {
		return nil, fmt.Errorf("failed to execute main.tf template: %w", err)
	}
	mainMaterial, err := types.NewHCLMaterial(mainBuf.Bytes(), "terraform/main.tf")
	if err != nil {
		return nil, fmt.Errorf("failed to create main.tf material: %w", err)
	}
	materials = append(materials, mainMaterial)

	// Generate variables.tf
	varsBuf := bytes.NewBuffer(nil)
	if err := varsTemplate.Execute(varsBuf, config); err != nil {
		return nil, fmt.Errorf("failed to execute variables.tf template: %w", err)
	}
	varsMaterial, err := types.NewHCLMaterial(varsBuf.Bytes(), "terraform/variables.tf")
	if err != nil {
		return nil, fmt.Errorf("failed to create variables.tf material: %w", err)
	}
	materials = append(materials, varsMaterial)

	// Generate providers.tf (common template)
	providersBuf := bytes.NewBuffer(nil)
	if err := providersTFTemplate.Execute(providersBuf, config); err != nil {
		return nil, fmt.Errorf("failed to execute providers.tf template: %w", err)
	}
	providersMaterial, err := types.NewHCLMaterial(providersBuf.Bytes(), "terraform/providers.tf")
	if err != nil {
		return nil, fmt.Errorf("failed to create providers.tf material: %w", err)
	}
	materials = append(materials, providersMaterial)

	// Generate outputs.tf
	outputsBuf := bytes.NewBuffer(nil)
	if err := outputsTemplate.Execute(outputsBuf, config); err != nil {
		return nil, fmt.Errorf("failed to execute outputs.tf template: %w", err)
	}
	outputsMaterial, err := types.NewHCLMaterial(outputsBuf.Bytes(), "terraform/outputs.tf")
	if err != nil {
		return nil, fmt.Errorf("failed to create outputs.tf material: %w", err)
	}
	materials = append(materials, outputsMaterial)

	return materials, nil
}

// getTemplatesForProvider returns the appropriate templates for the given infrastructure provider.
func (g *TerraformGenerator) getTemplatesForProvider(provider v1alpha1.InfrastructureProvider) (main, vars, outputs *types.Template, err error) {
	switch provider {
	case v1alpha1.InfrastructureProviderAWS:
		return awsMainTFTemplate, awsVariablesTFTemplate, awsOutputsTFTemplate, nil
	case v1alpha1.InfrastructureProviderGCP:
		return gcpMainTFTemplate, gcpVariablesTFTemplate, gcpOutputsTFTemplate, nil
	case v1alpha1.InfrastructureProviderAzure:
		return azureMainTFTemplate, azureVariablesTFTemplate, azureOutputsTFTemplate, nil
	default:
		return nil, nil, nil, fmt.Errorf("unsupported infrastructure provider: %s", provider)
	}
}
