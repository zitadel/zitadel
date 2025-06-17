target "docker-metadata-action" {}

variable "IMAGE_REGISTRY" {
  default = "ghcr.io/zitadel"
}

variable "BUILD_CACHE_KEY" {
  default = "local"
}

target "login-pnpm" {
  cache-from = [
    { "type": "registry", "ref": "${IMAGE_REGISTRY}/login-pnpm-buildcache:${BUILD_CACHE_KEY}" },
    { "type": "registry", "ref": "${IMAGE_REGISTRY}/login-pnpm-buildcache:latest" },
  ]
  dockerfile = "dockerfiles/login-pnpm.Dockerfile"
 }

target "login-dev-base" {
  cache-from = [
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-dev-base-buildcache:${BUILD_CACHE_KEY}"},
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-dev-base-buildcache:latest"},
  ]
  dockerfile = "dockerfiles/login-dev-base.Dockerfile"
  contexts = {
    login-pnpm = "target:login-pnpm"
  }
}

target "login-lint" {
  cache-from = [
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-lint-buildcache:${BUILD_CACHE_KEY}"},
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-lint-buildcache:latest"},
  ]
  dockerfile = "dockerfiles/login-lint.Dockerfile"
  contexts = {
    login-dev-base = "target:login-dev-base"
  }
}

variable "LOGIN_TEST_UNIT_TAG" {
  default = "login-test-unit:local"
}

target "login-test-unit" {
  dockerfile = "dockerfiles/login-test-unit.Dockerfile"
  contexts = {
    login-client   = "target:login-client"
  }
  output = ["type=docker"]
  tags = ["${LOGIN_TEST_UNIT_TAG}"]
}

target "login-client" {
  dockerfile = "dockerfiles/login-client.Dockerfile"
  contexts = {
    login-pnpm              = "target:login-pnpm"
    typescript-proto-client = "target:typescript-proto-client"
  }
}

target "typescript-proto-client" {
  dockerfile = "dockerfiles/typescript-proto-client.Dockerfile"
  contexts = {
    # We directly generate and download the client server-side with buf, so we don't need the proto files
    login-pnpm = "target:login-pnpm"
  }
}

# proto-files is only used to build core-mock against which the integration tests run.
# To build the proto-client, we use buf to generate and download the client code directly.
target "proto-files" {
  dockerfile = "dockerfiles/proto-files.Dockerfile"
  contexts = {
    login-pnpm = "target:login-pnpm"
  }
}

variable "CORE_MOCK_TAG" {
  default = "core-mock:local"
}

target "core-mock" {
  context = "apps/core-mock"
  contexts = {
    protos = "target:proto-files"
  }
  tags = ["${CORE_MOCK_TAG}"]
}

variable "LOGIN_TEST_INTEGRATION_TAG" {
  default = "login-test-integration:local"
}

target "login-test-integration" {
  dockerfile = "dockerfiles/login-test-integration.Dockerfile"
  contexts = {
    login-pnpm = "target:login-pnpm"
  }
  tags = ["${LOGIN_TEST_INTEGRATION_TAG}"]
}

variable "LOGIN_TEST_ACCEPTANCE_TAG" {
  default = "login-test-acceptance:local"
}

target "login-test-acceptance" {
  dockerfile = "dockerfiles/login-test-acceptance.Dockerfile"
  contexts = {
    login-pnpm = "target:login-pnpm"
  }
  tags = ["${LOGIN_TEST_ACCEPTANCE_TAG}"]
}

variable "LOGIN_TAG" {
  default = "zitadel-login:local"
}

# We run integration and acceptance tests against the next standalone server for docker.
target "login-standalone" {
  dockerfile = "dockerfiles/login-standalone.Dockerfile"
  contexts = {
    login-client = "target:login-client"
  }
  tags = ["${LOGIN_TAG}"]
}
