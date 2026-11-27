#!/usr/bin/env bash
# Prints the integration test packages for go test, slowest first.
#
# go test starts packages in the order it gets them, as many at once as -p
# allows. In go list order the slowest packages come late and finish the suite
# late, so the packages below, each over 200 s in CI, go first, slowest first.
# A package not listed follows in go list order.
set -euo pipefail

slowest=(
  internal/api/grpc/user/v2/integration_test
  internal/api/grpc/action/v2/integration_test
  internal/api/scim/integration_test
  internal/api/grpc/session/v2/integration_test
  internal/api/oidc/integration_test
  internal/api/grpc/action/v2beta/integration_test
  internal/api/grpc/user/v2beta/integration_test
  internal/api/grpc/instance/v2/integration_test
)

module=$(go list -m)
packages=$(go list -tags integration ./... | grep -e integration_test)
if [ -z "$packages" ]; then
  echo "no integration test packages found" >&2
  exit 1
fi
first=$(printf "$module/%s\n" "${slowest[@]}" | grep -xF -f <(printf '%s\n' "$packages") || true)
printf '%s\n' "$first"
printf '%s\n' "$packages" | grep -vxF -f <(printf '%s\n' "$first") || true
