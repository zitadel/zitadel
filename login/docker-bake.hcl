group "default" {
  targets = ["typescript-proto-client"]
}

target "login-platform" {
  dockerfile = "dockerfiles/login-platform.Dockerfile"
}

target "login-base" {
  dockerfile = "dockerfiles/login-base.Dockerfile"
  contexts = {
      login-platform = "target:login-platform"
  }
}

target "login-dependencies" {
  dockerfile = "dockerfiles/login-dependencies.Dockerfile"
  contexts = {
    login-base = "target:login-base"
  }
}

target "typescript-proto-client" {
  dockerfile = "dockerfiles/typescript-proto-client.Dockerfile"
  contexts = {
    # We directly generate and download the client server-side with buf, so we don't need the proto files
    login-base = "target:login-dependencies"
  }
}

# proto-files is only used to build core-mock against which the integration tests run.
# To build the proto-client, we use buf to generate and download the client code directly.
target "proto-files" {
  dockerfile = "dockerfiles/proto-files.Dockerfile"
  contexts = {
    login-base = "target:login-dependencies"
  }
}

target "core-mock" {
  context = "apps/core-mock"
  contexts = {
    protos = "target:proto-files"
  }
}

target "login-integration-testsuite" {
  dockerfile = "dockerfiles/login-integration-testsuite.Dockerfile"
  contexts = {
    login-base = "target:login-base"
  }
}

# We run integration and acceptance tests against the next standalone server for docker.
target "login-standalone" {
  dockerfile = "dockerfiles/login-standalone.Dockerfile"
  args = {
    NODE_ENV = "production"
  }
  contexts = {
      login-platform = "target:login-platform"
      login-base = "target:login-dependencies"
  }
}
