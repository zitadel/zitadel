group "default" {
  targets = ["typescript-proto-client"]
}

target "login-platform" {
  dockerfile = "dockerfiles/login-platform.Dockerfile"
}

target "login-pnpm" {
  dockerfile = "dockerfiles/login-pnpm.Dockerfile"
  contexts = {
    login-platform = "target:login-platform"
  }
}

target "login-dev-base" {
  dockerfile = "dockerfiles/login-dev-base.Dockerfile"
  contexts = {
      login-pnpm = "target:login-pnpm"
  }
}

target "login-lint" {
  dockerfile = "dockerfiles/login-lint.Dockerfile"
  contexts = {
    login-dev-base = "target:login-dev-base"
  }
}

target "login-test-unit" {
  dockerfile = "dockerfiles/login-test-unit.Dockerfile"
  contexts = {
    login-pnpm = "target:login-pnpm"
    login-dev-base = "target:login-dev-base"
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

target "core-mock" {
  context = "apps/core-mock"
  contexts = {
    protos = "target:proto-files"
  }
}

target "login-test-integration" {
  dockerfile = "dockerfiles/login-test-integration.Dockerfile"
  contexts = {
    login-pnpm = "target:login-pnpm"
  }
}

target "login-test-acceptance" {
  context = "apps/login-test-acceptance"
  contexts = {
    login-pnpm = "target:login-pnpm"
    login-test-acceptance-setup = "login-test-acceptance-setup:latest"
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
      login-pnpm = "target:login-pnpm"
  }
}
