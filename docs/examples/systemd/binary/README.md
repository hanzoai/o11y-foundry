# Systemd Binary Casting

This guide explains how to use systemd binary casting for deploying Hanzo O11y.

## Prerequisites

Before running `foundryctl cast`, install the required dependencies.

### 1. Install ClickHouse

ClickHouse is used as the telemetry store. Install both `clickhouse-server` and `clickhouse-keeper`.

- [ClickHouse Installation Guide](https://clickhouse.com/docs/en/install)

Verify installation:

```bash
clickhouse-server --version
clickhouse-keeper --version
```

### 2. Install Metastore Binary (PostgreSQL)

PostgreSQL is used as the metadata store.

- [PostgreSQL Installation Guide](https://www.postgresql.org/download/)

Verify installation:

```bash
postgres --version
```

### 3. Install Hanzo O11y Binary

```bash
curl -L https://github.com/Hanzo O11y/o11y/releases/latest/download/o11y_linux_$(uname -m | sed 's/x86_64/amd64/g' | sed 's/aarch64/arm64/g').tar.gz -o o11y.tar.gz
tar -xzf o11y.tar.gz

sudo mkdir -p /opt/o11y /var/lib/o11y
sudo cp -r o11y_linux_*/* /opt/o11y/
```

### 4. Install Ingester Binary (Hanzo O11y OTel Collector)

```bash
curl -L https://github.com/Hanzo O11y/o11y-otel-collector/releases/latest/download/o11y-otel-collector_linux_$(uname -m | sed 's/x86_64/amd64/g' | sed 's/aarch64/arm64/g').tar.gz -o o11y-otel-collector.tar.gz
tar -xzf o11y-otel-collector.tar.gz

sudo mkdir -p /opt/ingester /var/lib/ingester
sudo cp -r o11y-otel-collector_linux_*/* /opt/ingester/
```

### 5. Create o11y User

```bash
sudo useradd -r -s /sbin/nologin o11y
sudo chown -R o11y:o11y /opt/o11y /var/lib/o11y /opt/ingester /var/lib/ingester
```

Also, make sure that "o11y" user is allowed to transverse to the pours directory.

## Deployment

Create a `casting.yaml` file:

```yaml
apiVersion: v1alpha1
metadata:
  name: o11y
spec:
  deployment:
    flavor: binary
    mode: systemd
```

### 1. Verify Prerequisites

```bash
foundryctl gauge -f casting.yaml
```

### 2. Deploy Hanzo O11y

```bash
sudo foundryctl cast -f casting.yaml
```

### 3. Verify Services

Replace `<name>` with your `metadata.name` from `casting.yaml`:

```bash
systemctl status <name>-o11y.service
systemctl status <name>-ingester.service
systemctl status <name>-telemetrystore-clickhouse-0-0.service
systemctl status <name>-telemetrykeeper-clickhousekeeper-0.service
systemctl status <name>-metastore-postgres.service
```

View logs:

```bash
journalctl -u <name>-o11y.service -f
```


## Configuration

### Custom Binary Path

Use annotations to specify custom binary paths or other deployment metadata:

| Name | Type | Description |
|------|------|-------------|
| `foundry.o11y.hanzo.ai/o11y-binary-path` | string | Path to the Hanzo O11y binary |
| `foundry.o11y.hanzo.ai/ingester-binary-path` | string | Path to the OTel Collector binary |
| `foundry.o11y.hanzo.ai/metastore-postgres-binary-path` | string | Path to the PostgreSQL binary |

```yaml
apiVersion: v1alpha1
metadata:
  name: o11y
  annotations:
        foundry.o11y.hanzo.ai/o11y-binary-path: /opt/o11y/bin/o11y
        foundry.o11y.hanzo.ai/ingester-binary-path: /opt/ingester/bin/o11y-otel-collector
        foundry.o11y.hanzo.ai/metastore-postgres-binary-path: /usr/bin/postgres
spec:
  deployment:
    flavor: binary
    mode: systemd
```
