variable "LOGIN_DIR" {
  default = "./"
}

variable "DOCKERFILES_DIR" {
  default = "dockerfiles/"
}

target "login-pnpm" {
  context = "${LOGIN_DIR}"
  dockerfile = "${DOCKERFILES_DIR}login-pnpm.Dockerfile"
}

target "login-dev-base" {
  dockerfile = "${DOCKERFILES_DIR}login-dev-base.Dockerfile"
  context = "${LOGIN_DIR}"
  contexts = {
    login-pnpm = "target:login-pnpm"
  }
}

target "login-lint" {
  dockerfile = "${DOCKERFILES_DIR}login-lint.Dockerfile"
  context = "${LOGIN_DIR}"
  contexts = {
    login-dev-base = "target:login-dev-base"
  }
}

target "login-test-unit" {
  dockerfile = "${DOCKERFILES_DIR}login-test-unit.Dockerfile"
  context = "${LOGIN_DIR}"
  contexts = {
    login-client = "target:login-client"
  }
}

target "login-client" {
  dockerfile = "${DOCKERFILES_DIR}login-client.Dockerfile"
  context = "${LOGIN_DIR}"
  contexts = {
    login-pnpm              = "target:login-pnpm"
    typescript-proto-client = "target:typescript-proto-client"
  }
}

target "typescript-proto-client" {
  dockerfile = "${DOCKERFILES_DIR}typescript-proto-client.Dockerfile"
  context = "${LOGIN_DIR}"
  contexts = {
    # We directly generate and download the client server-side with buf, so we don't need the proto files
    login-pnpm = "target:login-pnpm"
  }
  output = ["type=docker"]
}

# proto-files is only used to build core-mock against which the integration tests run.
# To build the proto-client, we use buf to generate and download the client code directly.
target "proto-files" {
  dockerfile = "${DOCKERFILES_DIR}proto-files.Dockerfile"
  context = "${LOGIN_DIR}"
  contexts = {
    login-pnpm = "target:login-pnpm"
  }
}

variable "CORE_MOCK_TAG" {
  default = "core-mock:local"
}

target "core-mock" {
  context = "${LOGIN_DIR}apps/login-test-integration/core-mock"
  contexts = {
    protos = "target:proto-files"
  }
  tags   = ["${CORE_MOCK_TAG}"]
  output = ["type=docker"]
}

variable "LOGIN_TEST_INTEGRATION_TAG" {
  default = "login-test-integration:local"
}

target "login-test-integration" {
  dockerfile = "${DOCKERFILES_DIR}login-test-integration.Dockerfile"
  context = "${LOGIN_DIR}"
  contexts = {
    login-pnpm = "target:login-pnpm"
  }
  tags   = ["${LOGIN_TEST_INTEGRATION_TAG}"]
  output = ["type=docker"]
}

variable "LOGIN_TEST_ACCEPTANCE_TAG" {
  default = "login-test-acceptance:local"
}

target "login-test-acceptance" {
  dockerfile = "${DOCKERFILES_DIR}login-test-acceptance.Dockerfile"
  context = "${LOGIN_DIR}"
  contexts = {
    login-pnpm = "target:login-pnpm"
  }
  tags   = ["${LOGIN_TEST_ACCEPTANCE_TAG}"]
  output = ["type=docker"]
}

variable "LOGIN_TAG" {
  default = "zitadel-login:local"
}

target "docker-metadata-action" {}

# We run integration and acceptance tests against the next standalone server for docker.
target "login-standalone" {
  inherits   = ["docker-metadata-action"]
  dockerfile = "${DOCKERFILES_DIR}login-standalone.Dockerfile"
  context = "${LOGIN_DIR}"
  contexts = {
    login-client = "target:login-client"
  }
  tags   = ["${LOGIN_TAG}"]
  output = ["type=docker"]
}
