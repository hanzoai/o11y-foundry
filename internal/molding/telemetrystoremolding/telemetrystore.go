package telemetrystoremolding

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/signoz/foundry/api/v1alpha1"
	foundryerrors "github.com/signoz/foundry/internal/errors"
	"github.com/signoz/foundry/internal/molding"
	"github.com/signoz/foundry/internal/types"
)

var _ molding.Molding = (*telemetrystore)(nil)

type telemetrystore struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *telemetrystore {
	return &telemetrystore{
		logger: logger,
	}
}

func (molding *telemetrystore) Kind() v1alpha1.MoldingKind {
	return v1alpha1.MoldingKindTelemetryStore
}

func (molding *telemetrystore) MoldV1Alpha1(ctx context.Context, config *v1alpha1.Casting) error {
	data, err := molding.getData(config)
	if err != nil {
		molding.logger.ErrorContext(ctx, "failed to get data", foundryerrors.LogAttr(err))
		return err
	}

	metaDataName := config.Metadata.Name
	storeKind := config.Spec.TelemetryStore.Kind.String()

	// Check if ports differ across addresses - if so, need per-instance config
	// Also generate per-instance config if we have multiple nodes (shards * replicas > 1)
	requiresPerInstanceConfig := hasDistinctPorts(data.StoreAddresses)

	configData := make(map[string]string)

	// Generate functions file (shared across all instances)
	functionBuf := bytes.NewBuffer(nil)
	if err := FunctionsClickhousev2556YAML.Execute(functionBuf, data); err != nil {
		return fmt.Errorf("failed to execute function template: %w", err)
	}
	configData[StoreFunctionsFileName(metaDataName, storeKind, 0)] = functionBuf.String()

	if requiresPerInstanceConfig {
		// Per-instance config: generate separate config for each shard/replica
		for shard := 0; shard < data.ShardCount; shard++ {
			for replica := 0; replica < data.ReplicaCount; replica++ {
				data.ServerID = shard*data.ReplicaCount + replica

				configBuf := bytes.NewBuffer(nil)
				if err := ConfigClickhousev2556YAML.Execute(configBuf, data); err != nil {
					return fmt.Errorf("failed to execute config template for shard %d replica %d: %w", shard, replica, err)
				}
				instanceID := shard*data.ReplicaCount + replica
				configData[StoreInstanceConfigFileName(metaDataName, storeKind, instanceID, data.ReplicaCount)] = configBuf.String()
			}
		}
	} else {
		// Shared config: generate one config file for all instances
		configBuf := bytes.NewBuffer(nil)
		if err := ConfigClickhousev2556YAML.Execute(configBuf, data); err != nil {
			return fmt.Errorf("failed to execute config template: %w", err)
		}
		configData[StoreInstanceConfigFileName(metaDataName, storeKind, 0, data.ReplicaCount)] = configBuf.String()
	}

	config.Spec.TelemetryStore.Spec.Config.Data = configData

	return nil
}

// hasDistinctPorts returns true if addresses have different ports.
func hasDistinctPorts(addresses []types.Address) bool {
	if len(addresses) <= 1 {
		return false
	}
	firstPort := addresses[0].Port()
	for _, addr := range addresses[1:] {
		if addr.Port() != firstPort {
			return true
		}
	}
	return false
}

func (molding *telemetrystore) getData(config *v1alpha1.Casting) (Data, error) {
	storeAddresses := config.Spec.TelemetryStore.Status.Addresses[v1alpha1.TelemetryStoreClusterAddresses]
	if len(storeAddresses) == 0 {
		return Data{}, fmt.Errorf("telemetry store addresses not set in status")
	}

	cluster := config.Spec.TelemetryStore.Spec.Cluster

	shardCount := max(*cluster.Shards, 1)
	replicaCount := max(*cluster.Replicas+1, 1)

	expectedNodes := shardCount * replicaCount
	if len(storeAddresses) < expectedNodes {
		return Data{}, fmt.Errorf(
			"insufficient addresses: have %d, need %d (shards=%d × replicas=%d)",
			len(storeAddresses), expectedNodes, shardCount, replicaCount,
		)
	}

	newStoreAddrs, err := types.NewAddresses(storeAddresses[:expectedNodes])
	if err != nil {
		return Data{}, fmt.Errorf("failed to parse addresses: %w", err)
	}

	keeperAddresses := config.Spec.TelemetryKeeper.Status.Addresses[v1alpha1.TelemetryKeeperClientAddresses]
	if len(keeperAddresses) == 0 {
		return Data{}, fmt.Errorf("telemetry keeper addresses not set in status")
	}

	newKeeperAddrs, err := types.NewAddresses(keeperAddresses)
	if err != nil {
		return Data{}, fmt.Errorf("failed to parse addresses: %w", err)
	}

	// Check if CreatePerInstance flag should be set (when key exists in Extras)
	createPerInstance := false
	if config.Spec.TelemetryStore.Status.Extras != nil {
		if val, ok := config.Spec.TelemetryStore.Status.Extras[v1alpha1.CreatePerInstanceKey]; ok && val == "true" {
			createPerInstance = true
		}
	}

	return Data{
		StoreAddresses:    newStoreAddrs,
		KeeperAddresses:   newKeeperAddrs,
		ShardCount:        shardCount,
		ReplicaCount:      replicaCount,
		CreatePerInstance: createPerInstance,
	}, nil
}
