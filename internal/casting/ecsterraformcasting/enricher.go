package ecsterraformcasting

import (
	"context"
	"fmt"

	"github.com/hanzoai/o11y-foundry/api/v1alpha1"
	"github.com/hanzoai/o11y-foundry/api/v1alpha1/installation"
	"github.com/hanzoai/o11y-foundry/internal/domain"
	"github.com/hanzoai/o11y-foundry/internal/errors"
	"github.com/hanzoai/o11y-foundry/internal/molding"
)

const (
	telemetryStorePort        = 9000
	telemetryKeeperClientPort = 9181
	telemetryKeeperRaftPort   = 9234
	metaStorePort             = 5432
	signozAPIPort             = 8080
	signozOpampPort           = 4320
)

var _ molding.MoldingEnricher = (*ecsMoldingEnricher)(nil)

type ecsMoldingEnricher struct {
	materials []domain.StructuredMaterial
}

func newEcsMoldingEnricher(config *installation.Casting) (*ecsMoldingEnricher, error) {
	materials, err := getMaterials(config)
	if err != nil {
		return nil, errors.Wrapf(err, errors.TypeInternal, "failed to get materials")
	}

	return &ecsMoldingEnricher{materials: materials}, nil
}

func (enricher *ecsMoldingEnricher) EnrichStatus(ctx context.Context, kind v1alpha1.MoldingKind, config *installation.Casting) error {
	namespaceBytes, err := enricher.materials[0].GetBytes("resource.aws_service_discovery_private_dns_namespace.main.name")
	if err != nil {
		return errors.Wrapf(err, errors.TypeInternal, "failed to get namespace")
	}
	namespace := string(namespaceBytes)

	switch kind {
	case v1alpha1.MoldingKindTelemetryStore:
		sdName, err := enricher.materials[1].GetBytes("resource.aws_service_discovery_service.telemetrystore.name")
		if err != nil {
			return errors.Wrapf(err, errors.TypeInternal, "failed to get telemetrystore service discovery name")
		}
		fqdn := fmt.Sprintf("%s.%s", string(sdName), namespace)
		config.Spec.TelemetryStore.Status.Addresses.TCP = []string{domain.MustNewAddress("tcp", fqdn, telemetryStorePort).String()}

	case v1alpha1.MoldingKindTelemetryKeeper:
		sdName, err := enricher.materials[2].GetBytes("resource.aws_service_discovery_service.telemetrykeeper.name")
		if err != nil {
			return errors.Wrapf(err, errors.TypeInternal, "failed to get telemetrykeeper service discovery name")
		}
		fqdn := fmt.Sprintf("%s.%s", string(sdName), namespace)
		config.Spec.TelemetryKeeper.Status.Addresses.Client = []string{domain.MustNewAddress("tcp", fqdn, telemetryKeeperClientPort).String()}
		config.Spec.TelemetryKeeper.Status.Addresses.Raft = []string{domain.MustNewAddress("tcp", fqdn, telemetryKeeperRaftPort).String()}

	case v1alpha1.MoldingKindMetaStore:
		sdName, err := enricher.materials[3].GetBytes("resource.aws_service_discovery_service.metastore.name")
		if err != nil {
			return errors.Wrapf(err, errors.TypeInternal, "failed to get metastore service discovery name")
		}
		fqdn := fmt.Sprintf("%s.%s", string(sdName), namespace)
		config.Spec.MetaStore.Status.Addresses.DSN = []string{domain.MustNewAddress("tcp", fqdn, metaStorePort).String()}

	case v1alpha1.MoldingKindO11y:
		sdName, err := enricher.materials[4].GetBytes("resource.aws_service_discovery_service.signoz.name")
		if err != nil {
			return errors.Wrapf(err, errors.TypeInternal, "failed to get signoz service discovery name")
		}
		fqdn := fmt.Sprintf("%s.%s", string(sdName), namespace)
		config.Spec.O11y.Status.Addresses.APIServer = []string{domain.MustNewAddress("tcp", fqdn, signozAPIPort).String()}
		config.Spec.O11y.Status.Addresses.Opamp = []string{domain.MustNewAddress("ws", fqdn, signozOpampPort).String()}
	}

	return nil
}
