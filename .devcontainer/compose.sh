#!/bin/sh

# This script is used to run Docker Compose commands with the appropriate compose project set.
# This makes sure the commands don't break on different workspace_root folder names.
export DEVCONTAINER_LOCAL_WORKSPACE_ROOT=${DEVCONTAINER_LOCAL_WORKSPACE_ROOT:-${NX_WORKSPACE_ROOT:-zitadel}}

# Forward SIGINT/SIGTERM to the docker compose child process and exit 0.
# This prevents Nx from treating signal-terminated continuous tasks as failures.
_stop() {
  kill -TERM "$_child" 2>/dev/null
  wait "$_child"
  exit 0
}
trap _stop INT TERM

docker compose -p ${DEVCONTAINER_LOCAL_WORKSPACE_ROOT##*/}_devcontainer "$@" &
_child=$!
wait "$_child"
