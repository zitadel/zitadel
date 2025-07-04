variable "LOGIN_DIR" {
  default = "./"
}

variable "DOCKERFILES_DIR" {
  default = "dockerfiles/"
}

# The release target is overwritten in docker-bake-release.hcl
# It makes sure the image is built for multiple platforms.
# By default the platforms property is empty, so images are only built for the current bake runtime platform.
target "release" {}

# typescript-proto-client is used to generate the client code for the login service.
# It is not login-prefixed, so it is easily extendable.
# To extend this bake-file.hcl, set the context of all login-prefixed targets to a different directory.
# For example docker bake --file login/docker-bake.hcl --file docker-bake.hcl --set login-*.context=./login/
# The zitadel repository uses this to generate the client and the mock server from local proto files.
target "typescript-proto-client" {
  inherits   = ["release"]
  dockerfile = "${DOCKERFILES_DIR}typescript-proto-client.Dockerfile"
  contexts = {
    # We directly generate and download the client server-side with buf, so we don't need the proto files
    login-pnpm = "target:login-pnpm"
  }
}

# We prefix the target with login- so we can reuse the writing of protos if we overwrite the typescript-proto-client target.
target "login-typescript-proto-client-out" {
  dockerfile = "${DOCKERFILES_DIR}login-typescript-proto-client-out.Dockerfile"
  contexts = {
    typescript-proto-client = "target:typescript-proto-client"
  }
  output = [
    "type=local,dest=${LOGIN_DIR}packages/zitadel-proto"
  ]
}

# proto-files is only used to build core-mock against which the integration tests run.
# To build the proto-client, we use buf to generate and download the client code directly.
# It is not login-prefixed, so it is easily extendable.
# To extend this bake-file.hcl, set the context of all login-prefixed targets to a different directory.
# For example docker bake --file login/docker-bake.hcl --file docker-bake.hcl --set login-*.context=./login/
# The zitadel repository uses this to generate the client and the mock server from local proto files.
target "proto-files" {
  inherits   = ["release"]
  dockerfile = "${DOCKERFILES_DIR}proto-files.Dockerfile"
  contexts = {
    login-pnpm = "target:login-pnpm"
  }
}

variable "NODE_VERSION" {
  default = "20"
}

target "login-pnpm" {
  inherits   = ["release"]
  dockerfile = "${DOCKERFILES_DIR}login-pnpm.Dockerfile"
  args = {
    NODE_VERSION = "${NODE_VERSION}"
  }
}

target "login-dev-base" {
  dockerfile = "${DOCKERFILES_DIR}login-dev-base.Dockerfile"
  contexts = {
    login-pnpm = "target:login-pnpm"
  }
}

target "login-lint" {
  dockerfile = "${DOCKERFILES_DIR}login-lint.Dockerfile"
  contexts = {
    login-dev-base = "target:login-dev-base"
  }
}

target "login-test-unit" {
  dockerfile = "${DOCKERFILES_DIR}login-test-unit.Dockerfile"
  contexts = {
    login-client = "target:login-client"
  }
}

target "login-client" {
  inherits   = ["release"]
  dockerfile = "${DOCKERFILES_DIR}login-client.Dockerfile"
  contexts = {
    login-pnpm              = "target:login-pnpm"
    typescript-proto-client = "target:typescript-proto-client"
  }
}

variable "LOGIN_CORE_MOCK_TAG" {
  default = "login-core-mock:local"
}

# the core-mock context must not be overwritten, so we don't prefix it with login-.
target "core-mock" {
  context = "${LOGIN_DIR}apps/login-test-integration/core-mock"
  contexts = {
    protos = "target:proto-files"
  }
  tags = ["${LOGIN_CORE_MOCK_TAG}"]
}

variable "LOGIN_TEST_INTEGRATION_TAG" {
  default = "login-test-integration:local"
}

target "login-test-integration" {
  dockerfile = "${DOCKERFILES_DIR}login-test-integration.Dockerfile"
  contexts = {
    login-pnpm = "target:login-pnpm"
  }
  tags = ["${LOGIN_TEST_INTEGRATION_TAG}"]
}

variable "LOGIN_TEST_ACCEPTANCE_TAG" {
  default = "login-test-acceptance:local"
}

target "login-test-acceptance" {
  dockerfile = "${DOCKERFILES_DIR}login-test-acceptance.Dockerfile"
  contexts = {
    login-pnpm = "target:login-pnpm"
  }
  tags = ["${LOGIN_TEST_ACCEPTANCE_TAG}"]
}

variable "LOGIN_TAG" {
  default = "zitadel-login:local"
}

target "docker-metadata-action" {
  # In the pipeline, this target is overwritten by the docker metadata action.
  tags = ["${LOGIN_TAG}"]
}

# We run integration and acceptance tests against the next standalone server for docker.
target "login-standalone" {
  inherits = [
    "docker-metadata-action",
    "release",
  ]
  dockerfile = "${DOCKERFILES_DIR}login-standalone.Dockerfile"
  contexts = {
    login-client = "target:login-client"
  }
}

target "login-standalone-out" {
  inherits = ["login-standalone"]
  target   = "login-standalone-out"
  output = [
    "type=local,dest=${LOGIN_DIR}apps/login/standalone"
  ]
}
