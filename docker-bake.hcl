variable "GITHUB_SHA" {
  default = "latest"
}

variable "REGISTRY" {
  default = "ghcr.io/zitadel"
}

group "all" {
  targets = ["build", "lint", "image", "unit"]
}

group "build" {
  targets = ["console-build", "core-build"]
}

group "generate" {
  targets = ["console-generate" , "core-generate"]
}

group "lint" {
  targets = ["console-lint", "core-lint"]
}

group "image" {
  targets = ["console-image", "core-image"]
}

group "unit" {
  targets = ["core-unit"]
}

target "devcontainer" {
  dockerfile = "Dockerfile.devcontainer"
  context = "."
  tags = ["${REGISTRY}/base:${GITHUB_SHA}"]
  push = false
}

target "_console" {
  dockerfile = "Dockerfile.console"
  context = "."
  contexts = {
    devcontainer = "target:devcontainer"
    nginx = "docker-image://nginx:1.27-alpine"
  }
}

target "console-generate" {
  inherits = ["_console"]
  output = ["type=local,dest=./"]
  target = "generate"
  cache-to = ["type=gha,ignore-error=true,mode=max,scope=console-generate"]
  cache-from = ["type=gha,scope=console-generate"]
}

target "console-build" {
  inherits = ["_console"]
  output = ["type=local,dest=.build/console"]
  target = "build"
  cache-to = ["type=gha,ignore-error=true,mode=max,scope=console-build"]
  cache-from = ["type=gha,scope=console-build"]
}

target "_core" {
  dockerfile = "Dockerfile.core"
  context = "."
  contexts = {
    devcontainer = "target:devcontainer"
  }
}

target "core-generate" {
  inherits = ["_core"]
  output = ["type=local,dest=./"]
  target = "generate"
  cache-to = ["type=gha,ignore-error=true,mode=max,scope=core-generate"]
  cache-from = ["type=gha,scope=core-generate"]
}

target "core-build" {
  inherits = ["_core"]
  name = "core-build-${os}-${arch}"
    matrix = {
    os = ["linux", "darwin", "windows"]
    arch = ["amd64", "arm64"]
  }
  args = {
    OS = os
    ARCH = arch
  }
  output = ["type=local,dest=.build/core"]
  contexts = {
    console = "target:console-build"
  }
  target = "build"
  cache-to = ["type=gha,ignore-error=true,mode=max,scope=core-build"]
  cache-from = ["type=gha,scope=core-build"]
}