target "docker-metadata-action" {}

variable "IMAGE_REGISTRY" {
  default = "ghcr.io/zitadel"
}

variable "REF_TAG" {
  default = "local"
}

target "login-pnpm" {
  cache-from = [
    { "type": "registry", "ref": "${IMAGE_REGISTRY}/login-pnpm-buildcache:${REF_TAG}" },
    { "type": "registry", "ref": "${IMAGE_REGISTRY}/login-pnpm-buildcache:latest" },
  ]
  dockerfile = "dockerfiles/login-pnpm.Dockerfile"
 }

target "login-dev-base" {
  cache-from = [
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-dev-base-buildcache:${REF_TAG}"},
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-dev-base-buildcache:latest"},
  ]
  dockerfile = "dockerfiles/login-dev-base.Dockerfile"
  contexts = {
    login-pnpm = "target:login-pnpm"
  }
}

target "login-lint" {
  cache-from = [
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-lint-buildcache:${REF_TAG}"},
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-lint-buildcache:latest"},
  ]
  dockerfile = "dockerfiles/login-lint.Dockerfile"
  contexts = {
    login-dev-base = "target:login-dev-base"
  }
}

target "login-test-unit" {
  cache-from = [
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-test-unit-buildcache:${REF_TAG}"},
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-test-unit-buildcache:latest"},
  ]
  dockerfile = "dockerfiles/login-test-unit.Dockerfile"
  contexts = {
    login-client   = "target:login-client"
  }
}

target "login-client" {
  cache-from = [
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-client-buildcache:${REF_TAG}"},
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-client-buildcache:latest"},
  ]
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
    output = ["type=docker"]
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
  output = ["type=docker"]
}

variable "LOGIN_TEST_INTEGRATION_TAG" {
  default = "login-test-integration:local"
}

target "login-test-integration" {
  cache-from = [
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-test-integration-buildcache:${REF_TAG}"},
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-test-integration-buildcache:latest"},
  ]
  dockerfile = "dockerfiles/login-test-integration.Dockerfile"
  contexts = {
    login-pnpm = "target:login-pnpm"
  }
  tags = ["${LOGIN_TEST_INTEGRATION_TAG}"]
  output = ["type=docker"]
}

variable "LOGIN_TEST_ACCEPTANCE_TAG" {
  default = "login-test-acceptance:local"
}

target "login-test-acceptance" {
  cache-from = [
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-test-acceptance-buildcache:${REF_TAG}"},
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-test-acceptance-buildcache:latest"},
  ]
  dockerfile = "dockerfiles/login-test-acceptance.Dockerfile"
  contexts = {
    login-pnpm = "target:login-pnpm"
  }
  tags = ["${LOGIN_TEST_ACCEPTANCE_TAG}"]
  output = ["type=docker"]
}

variable "LOGIN_TAG" {
  default = "zitadel-login:local"
}

# We run integration and acceptance tests against the next standalone server for docker.
target "login-standalone" {
  cache-from = [
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-buildcache:${REF_TAG}"},
    {"type": "registry", "ref": "${IMAGE_REGISTRY}/login-buildcache:latest"},
  ]
  dockerfile = "dockerfiles/login-standalone.Dockerfile"
  contexts = {
    login-client = "target:login-client"
  }
  tags = ["${LOGIN_TAG}"]
  output = ["type=docker"]
}
