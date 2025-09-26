#!/bin/sh

# This script is used to run Docker Compose commands with the appropriate compose project set.
# This makes sure the commands don't break on different workspace_root folder names.
export DEVCONTAINER_LOCAL_WORKSPACE_ROOT=${DEVCONTAINER_LOCAL_WORKSPACE_ROOT:-$NX_WORKSPACE_ROOT}
docker compose -p ${DEVCONTAINER_LOCAL_WORKSPACE_ROOT##*/}_devcontainer "$@"