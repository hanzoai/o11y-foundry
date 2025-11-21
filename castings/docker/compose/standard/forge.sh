#!/bin/bash
# forge.sh — Foundry Config Forge
#
# The forge transforms configuration molds into deployable artifacts(pours) for running casting.
set -euo pipefail

# Directories
BASE_CONFIG_DIR="../../../../moldings"
FLAVOUR="standard"
CONFIGS_OUT="./pours"
ENV_DIR="."

# Components to render
COMPONENTS=("clickhouse" "signoz" "zookeeper" "otel-collector")

log_info()  { echo "[INFO] $1"; }
log_error() { echo "[ERROR] $1"; }

# Ensure .env exists
if [ ! -f "${ENV_DIR}/.env" ]; then
    log_error ".env file not found in ${ENV_DIR}"
    exit 1
fi

log_info "Loading environment variables from .env"
set -a
source "${ENV_DIR}/.env"
set +a

# Required environment variables
REQUIRED_VARS=(
    "ZOOKEEPER_HOST"
    "ZOOKEEPER_PORT"
    "CLICKHOUSE_HOST"
    "CLICKHOUSE_PORT"
    "SIGNOZ_HOST"
)

missing_vars=()
for var in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!var:-}" ]; then
        missing_vars+=("$var")
    fi
done

if [ ${#missing_vars[@]} -gt 0 ]; then
    log_error "Missing required variables in .env:"
    for var in "${missing_vars[@]}"; do
        log_error "  $var"
    done
    exit 1
fi

log_info "Rendering configs for flavour: ${FLAVOUR}"

# Prepare output directory
rm -rf "${CONFIGS_OUT}"
mkdir -p "${CONFIGS_OUT}"

#
# Render each component folder
#
for component in "${COMPONENTS[@]}"; do
    SRC_DIR="${BASE_CONFIG_DIR}/${component}/${FLAVOUR}"

    if [ ! -d "${SRC_DIR}" ]; then
        log_info "Skipping ${component}: no configs found for flavour"
        continue
    fi

    OUT_DIR="${CONFIGS_OUT}/${component}"
    mkdir -p "${OUT_DIR}"

    log_info "Rendering ${component} configs"

    # Loop through all files
    find "${SRC_DIR}" -type f | while read -r file; do
        rel_path="${file#${SRC_DIR}/}"
        out_file="${OUT_DIR}/${rel_path}"
        out_dir_path="$(dirname "${out_file}")"

        mkdir -p "${out_dir_path}"

        # Render .env files too
        if [[ "${file}" == *.env ]]; then
            envsubst < "${file}" > "${out_file}"
            log_info "Rendered: ${component}/${rel_path}"
            continue
        fi

        # If file contains variable placeholders, render it
        if grep -q '\${' "${file}"; then
            envsubst < "${file}" > "${out_file}"
            log_info "Rendered: ${component}/${rel_path}"
        else
            cp "${file}" "${out_file}"
            log_info "Copied:   ${component}/${rel_path}"
        fi
    done
done

log_info "All configs rendered successfully"
log_info "Output directory: ${CONFIGS_OUT}"
