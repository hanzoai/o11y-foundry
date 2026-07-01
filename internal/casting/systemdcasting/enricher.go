package systemdcasting

import (
	"context"
	"fmt"
	"path"

	"github.com/hanzoai/o11y-foundry/api/v1alpha1"
	"github.com/hanzoai/o11y-foundry/api/v1alpha1/installation"
	"github.com/hanzoai/o11y-foundry/internal/domain"
	"github.com/hanzoai/o11y-foundry/internal/molding"
)

var _ molding.MoldingEnricher = (*linuxMoldingEnricher)(nil)

const (
	baseTelemetryKeeperClientPort = 9181
	baseTelemetryKeeperRaftPort   = 9234
	baseTelemetryStoreClusterPort = 9000
	baseMetaStorePostgresPort     = 5432
)

type linuxMoldingEnricher struct {
	materials []domain.Material
}

func newLinuxMoldingEnricher(_ *installation.Casting) *linuxMoldingEnricher {
	return &linuxMoldingEnricher{materials: []domain.Material{}}
}

func (e *linuxMoldingEnricher) EnrichStatus(ctx context.Context, kind v1alpha1.MoldingKind, config *installation.Casting) error {
	switch kind {
	case v1alpha1.MoldingKindTelemetryStore:
		return e.enrichTelemetryStore(config)
	case v1alpha1.MoldingKindTelemetryKeeper:
		return e.enrichTelemetryKeeper(config)
	case v1alpha1.MoldingKindMetaStore:
		return e.enrichMetaStore(config)
	case v1alpha1.MoldingKindO11y:
		return e.enrichO11y(config)
	case v1alpha1.MoldingKindIngester:
		return e.enrichIngester(config)
	}
	return nil
}

func (e *linuxMoldingEnricher) enrichTelemetryStore(config *installation.Casting) error {
	spec := &config.Spec.TelemetryStore
	cluster := spec.Spec.Cluster

	replicas := 1
	shards := 1
	if cluster.Replicas != nil {
		replicas = max(*cluster.Replicas+1, 1)
	}
	if cluster.Shards != nil {
		shards = max(*cluster.Shards, 1)
	}

	if replicas > 1 || shards > 1 {
		return fmt.Errorf("deployment mode '%s' does not support Distributed Clickhouse Setup, raise an issue at https://github.com/hanzoai/o11y-foundry/issues", config.Spec.Deployment.Mode)
	}

	// Generate addresses for each shard/replica
	var addresses []string
	for shard := 0; shard < shards; shard++ {
		for replica := 0; replica < replicas; replica++ {
			port := baseTelemetryStoreClusterPort + (shard * replicas) + replica
			addresses = append(addresses, domain.MustNewAddress("tcp", "localhost", port).String())
		}
	}

	config.Spec.TelemetryStore.Status.Addresses.TCP = addresses
	return nil
}

func (e *linuxMoldingEnricher) enrichTelemetryKeeper(config *installation.Casting) error {
	spec := &config.Spec.TelemetryKeeper
	cluster := spec.Spec.Cluster

	replicas := 1
	if cluster.Replicas != nil {
		replicas = max(*cluster.Replicas, 1)
	}

	if replicas > 1 {
		return fmt.Errorf("deployment mode '%s' does not support Distributed Clickhouse Setup, raise an issue at https://github.com/hanzoai/o11y-foundry/issues", config.Spec.Deployment.Mode)
	}

	var clientAddresses, raftAddresses []string
	for r := 0; r < replicas; r++ {
		clientAddresses = append(clientAddresses, domain.MustNewAddress("tcp", "localhost", baseTelemetryKeeperClientPort+r).String())
		raftAddresses = append(raftAddresses, domain.MustNewAddress("tcp", "localhost", baseTelemetryKeeperRaftPort+r).String())
	}

	config.Spec.TelemetryKeeper.Status.Addresses.Client = clientAddresses
	config.Spec.TelemetryKeeper.Status.Addresses.Raft = raftAddresses
	return nil
}

func (e *linuxMoldingEnricher) enrichMetaStore(config *installation.Casting) error {
	switch config.Spec.MetaStore.Kind {
	case installation.MetaStoreKindSQLite:
		// SQLite — no addresses or binaries to enrich.
	case installation.MetaStoreKindPostgres:
		dsn := domain.MustNewAddress("postgres", "localhost", baseMetaStorePostgresPort).String()
		config.Spec.MetaStore.Status.Addresses.DSN = []string{dsn}

		// Get the annotation value
		metastoreBin := config.Metadata.Annotations["foundry.o11y.hanzo.ai/metastore-postgres-binary-path"]

		// If it's missing, apply the default and write it back
		if metastoreBin == "" {
			metastoreBin = "/usr/bin/postgres"

			if config.Metadata.Annotations == nil {
				config.Metadata.Annotations = make(map[string]string)
			}
			config.Metadata.Annotations["foundry.signoz.io/metastore-postgres-binary-path"] = metastoreBin
		}
		config.Metadata.Annotations["foundry.o11y.hanzo.ai/metastore-postgres-binary-path"] = metastoreBin
	}
	return nil
}

func (e *linuxMoldingEnricher) enrichO11y(config *installation.Casting) error {
	config.Spec.O11y.Status.Addresses.Opamp = []string{
		domain.MustNewAddress("ws", "localhost", 4320).String(),
	}
	config.Spec.O11y.Status.Addresses.APIServer = []string{
		domain.MustNewAddress("tcp", "localhost", 8080).String(),
	}

	// Get the annotation value
	o11yBin := config.Metadata.Annotations["foundry.o11y.hanzo.ai/o11y-binary-path"]

	// If it's missing, apply the default and write it back
	if o11yBin == "" {
		o11yBin = "/opt/o11y/bin/o11y"

		if config.Metadata.Annotations == nil {
			config.Metadata.Annotations = make(map[string]string)
		}
		config.Metadata.Annotations["foundry.o11y.hanzo.ai/o11y-binary-path"] = o11yBin
	}

	// The binary defaults these to its in-container paths, so point them at the
	// extracted tarball tree (binary lives at <root>/bin/signoz).
	root := path.Dir(path.Dir(o11yBin))
	if config.Spec.O11y.Status.Env == nil {
		config.Spec.O11y.Status.Env = make(map[string]string)
	}
	env := config.Spec.O11y.Status.Env
	env["SIGNOZ_WEB_DIRECTORY"] = path.Join(root, "web")
	env["SIGNOZ_EMAILING_TEMPLATES_DIRECTORY"] = path.Join(root, "templates", "email")
	env["SIGNOZ_ALERTMANAGER_SIGNOZ_TEMPLATES"] = path.Join(root, "templates", "alertmanager", "*.gotmpl")

	return nil
}

func (e *linuxMoldingEnricher) enrichIngester(config *installation.Casting) error {
	config.Spec.Ingester.Status.Addresses.OTLP = []string{
		domain.MustNewAddress("tcp", "localhost", 4317).String(),
	}

	// Get the annotation value
	ingesterBin := config.Metadata.Annotations["foundry.o11y.hanzo.ai/ingester-binary-path"]

	// If it's missing, apply the default and write it back
	if ingesterBin == "" {
		ingesterBin = "/opt/ingester/bin/o11y-otel-collector"

		if config.Metadata.Annotations == nil {
			config.Metadata.Annotations = make(map[string]string)
		}
		config.Metadata.Annotations["foundry.o11y.hanzo.ai/ingester-binary-path"] = ingesterBin
	}

	if config.Spec.Ingester.Status.Env == nil {
		config.Spec.Ingester.Status.Env = make(map[string]string)
	}
	config.Spec.Ingester.Status.Env["SIGNOZ_OTEL_COLLECTOR_TIMEOUT"] = "10m"

	return nil
}
