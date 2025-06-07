group "default" {
  targets = ["typescript-proto-client"]
}

target "login-platform" {
  dockerfile = "dockerfiles/login-platform.Dockerfile"
}

target "login-dev-base" {
  dockerfile = "dockerfiles/login-dev-base.Dockerfile"
  contexts = {
      login-platform = "target:login-platform"
  }
}

target "login-dev-dependencies" {
  dockerfile = "dockerfiles/login-dev-dependencies.Dockerfile"
  contexts = {
    login-dev-base = "target:login-dev-base"
  }
}

# proto-files is only used to build core-mock against which the integration tests run.
# To build the proto-client, we use buf to generate and download the client code directly.
target "proto-files" {
  dockerfile = "dockerfiles/proto-files.Dockerfile"
  contexts = {
    login-dev-base = "target:login-dev-dependencies"
  }
}

target "core-mock" {
  context = "apps/login/mock"
  dockerfile = "Dockerfile"
  contexts = {
    protos = "target:proto-files"
  }
}

target "login-integration-testsuite" {
  context = "apps/login/cypress"
  contexts = {
      login-dev-dependencies = "target:login-dev-dependencies"
  }
}

target "typescript-proto-client" {
  dockerfile = "dockerfiles/typescript-proto-client.Dockerfile"
  contexts = {
    # We directly generate and download the client server-side with buf, so we don't need the proto files
    login-dev-base = "target:login-dev-dependencies"
  }
}

# We run integration and acceptance tests against the next standalone server for docker.
target "login-image" {
  dockerfile = "dockerfiles/login-image.Dockerfile"
  args = {
    NODE_ENV = "production"
  }
  contexts = {
      login-platform = "target:login-platform"
      login-dev-base = "target:login-dev-dependencies"
  }
}
