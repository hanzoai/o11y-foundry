# Schema Definition
Definition and logic behind the schema for `casting.yml` - [reference](https://www.notion.so/signoz/Foundry-263fcc6bcd19804c89efcbb7bae9effb?source=copy_link#26bfcc6bcd198014b184de990757ff04)
- The schema is subject to change as it is in early development
- Schema is used by `foundryctl` to generate the moldings for each platform

The schema currently accepts the following:

## Platforms
- docker
- linux
- kubernetes
- aws
- gcp
- azure
- windows
- aws_marketplace

## Deployment Types
- standard
- highly-available

# Requirements
- Cue installed `brew install cue-lang/tap/cue`
- Text editor with cue lang highlights (VScode Extension available)


# Usage
- Navigate to the each platform directory
- Use cue to validate the schema is correct `cue vet -c schema/casting-schema.cue casting.yml`
- Cue will outline if your `casting.yml` matches the schema
 
### Example
 Given the following `casting.yaml`
 ```yaml
schemaVersion: v1
platform: docker
type: highly-available
components:
  - name: signoz
    replicas: 1
    version: "0.39.2"
    env:
      - key: SIGNOZ_HOST
        value: "signoz"
  - name: clickhouse
    replicas: 1
    version: "24.8"
    env:
      - key: CLICKHOUSE_HOST
        value: clickhouse
      - key: CLICKHOUSE_PORT
        value: "9000"
  - name: zookeeper
    replicas: 1
    version: "3.7"
    env:
      - key: ZOOKEEPER_HOST
        value: zookeeper-1
      - key: ZOOKEEPER_PORT
        value: "2181"
  - name: signoz-collector
    replicas: 1
    version: "0.39.2"
requirements:
  - docker
  - docker-compose

 ```
 Validate it with `cue vet -v casting-schema.cue casting-example.yaml`

If you replace a field or name it something different, for eg. `platform: test`. Will fail as it does not match the current supported platforms.

```bash
platform: 7 errors in empty disjunction:
platform: conflicting values "aws" and "test":
    ./casting-example.yaml:2:11
    ./casting-schema.cue:7:48
    ./casting-schema.cue:40:17
    ./casting-schema.cue:51:1
```
