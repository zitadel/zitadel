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

target "_console" {
  dockerfile = "Dockerfile.console"
  context = "."
  contexts = {
    node = "docker-image://node:22"
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
    golang = "docker-image://golang:1.24"
  }
  args = {
    SASS_VERSION      = "1.64.1"
    GOLANG_CI_VERSION = "1.64.5"
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
  output = ["type=local,dest=.build/core"]
  contexts = {
    console = "target:console-build"
  }
  target = "build"
  cache-to = ["type=gha,ignore-error=true,mode=max,scope=core-build"]
  cache-from = ["type=gha,scope=core-build"]
}