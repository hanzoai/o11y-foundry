package kuberneteskustomizecasting

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/signoz/foundry/api/v1alpha1"
	rootcasting "github.com/signoz/foundry/internal/casting"
	"github.com/signoz/foundry/internal/molding"
	"github.com/signoz/foundry/internal/types"
)

var _ rootcasting.Casting = (*kustomizeCasting)(nil)

type kustomizeCasting struct {
	logger   *slog.Logger
	castings []*types.Template
}

// KustomizeKnobs defines the supported knobs for the kustomize casting.
type KustomizeKnobs struct {
	// Resources defines CPU and memory requests/limits for the container.
	Resources map[string]any `json:"resources,omitempty" yaml:"resources,omitempty"`

	// Tolerations defines pod tolerations for scheduling.
	Tolerations []map[string]any `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`

	// NodeSelector defines node selection constraints.
	NodeSelector map[string]string `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`

	// Affinity defines pod affinity and anti-affinity rules.
	Affinity map[string]any `json:"affinity,omitempty" yaml:"affinity,omitempty"`

	// TopologySpreadConstraints defines how pods are spread across topology domains.
	TopologySpreadConstraints []map[string]any `json:"topologySpreadConstraints,omitempty" yaml:"topologySpreadConstraints,omitempty"`

	// PodSecurityContext defines the pod-level security context (e.g. runAsUser, fsGroup).
	PodSecurityContext map[string]any `json:"podSecurityContext,omitempty" yaml:"podSecurityContext,omitempty"`

	// SecurityContext defines the container-level security context (e.g. runAsNonRoot, readOnlyRootFilesystem).
	SecurityContext map[string]any `json:"securityContext,omitempty" yaml:"securityContext,omitempty"`

	// ImagePullSecrets lists secret names for pulling container images.
	ImagePullSecrets []map[string]any `json:"imagePullSecrets,omitempty" yaml:"imagePullSecrets,omitempty"`

	// MinReadySeconds is the minimum seconds a pod should be ready before considered available (Deployment/StatefulSet).
	MinReadySeconds *int `json:"minReadySeconds,omitempty" yaml:"minReadySeconds,omitempty"`

	// StorageSize defines the PVC storage request size (e.g. "10Gi").
	StorageSize string `json:"storageSize,omitempty" yaml:"storageSize,omitempty"`

	// StorageClass defines the PVC storage class name.
	StorageClass string `json:"storageClass,omitempty" yaml:"storageClass,omitempty"`

	// ServiceType defines the Kubernetes service type (e.g. "ClusterIP", "LoadBalancer").
	ServiceType string `json:"serviceType,omitempty" yaml:"serviceType,omitempty"`

	// ServiceAnnotations defines annotations to add to the service.
	ServiceAnnotations map[string]string `json:"serviceAnnotations,omitempty" yaml:"serviceAnnotations,omitempty"`

	// ServiceLabels defines labels to add to the service.
	ServiceLabels map[string]string `json:"serviceLabels,omitempty" yaml:"serviceLabels,omitempty"`

	// PodAnnotations defines annotations to add to the pod template.
	PodAnnotations map[string]string `json:"podAnnotations,omitempty" yaml:"podAnnotations,omitempty"`

	// PodLabels defines extra labels to add to the pod template (podExtraLabels in Kubedeploy, podLabels in Altinity).
	PodLabels map[string]string `json:"podLabels,omitempty" yaml:"podLabels,omitempty"`
}

func New(logger *slog.Logger) *kustomizeCasting {
	return &kustomizeCasting{
		logger: logger,
		castings: []*types.Template{
			clickhouseOperatorClusterrole,
			clickhouseOperatorClusterrolebinding,
			clickhouseOperatorConfigmap,
			clickhouseOperatorDeployment,
			clickhouseOperatorService,
			clickhouseOperatorServiceaccount,
			clickhouseInstanceInstallation,
			clickhouseInstanceConfigmap,
			clickhouseKeeperInstallation,
			signozService,
			signozServiceaccount,
			signozStatefulset,
			ingesterConfigmap,
			ingesterDeployment,
			ingesterService,
			ingesterServiceaccount,
			telemetrystoreMigratorJob,
			clickhouseOperatorKustomization,
			clickhouseInstallationKustomization,
			clickhouseKeeperKustomization,
			signozKustomization,
			ingesterKustomization,
			telemetrystoreMigratorKustomization,
			deploymentKustomization,
		},
	}
}

func (c *kustomizeCasting) Enricher(ctx context.Context, config *v1alpha1.Casting) (molding.MoldingEnricher, error) {
	return newKustomizeMoldingEnricher(config), nil
}

func (c *kustomizeCasting) Forge(ctx context.Context, cfg v1alpha1.Casting, poursPath string) ([]types.Material, error) {
	var materials []types.Material
	for _, tmpl := range c.castings {
		m, err := c.forgeCasting(tmpl, &cfg, poursPath)
		if err != nil {
			return nil, fmt.Errorf("failed to forge: %w", err)
		}
		materials = append(materials, m...)
	}
	return materials, nil
}

func (c *kustomizeCasting) Cast(ctx context.Context, config v1alpha1.Casting, poursPath string) error {
	c.logger.InfoContext(ctx, "Please run 'forge' first to generate the Kubernetes manifests",
		slog.String("pours_path", poursPath))
	c.logger.InfoContext(ctx, "After forging, apply with: kubectl apply -k pours/deployment",
		slog.String("Docs", "https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/"))
	return nil
}

func (c *kustomizeCasting) forgeCasting(tmpl *types.Template, cfg *v1alpha1.Casting, poursPath string) ([]types.Material, error) {
	templatePath := tmpl.GetPath()
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cfg); err != nil {
		return nil, fmt.Errorf("execute template %s: %w", templatePath, err)
	}
	relPath := strings.TrimPrefix(templatePath, "templates/")
	relPath = strings.TrimSuffix(relPath, filepath.Ext(relPath))
	path := filepath.Join(rootcasting.DeploymentDir, relPath)
	material, err := types.NewYAMLMaterial(buf.Bytes(), path)
	if err != nil {
		return nil, fmt.Errorf("create material %s: %w", templatePath, err)
	}
	return []types.Material{material}, nil
}
