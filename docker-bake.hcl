variable "GITHUB_SHA" {
  default = "latest"
}

variable "REGISTRY" {
  default = "ghcr.io/zitadel"
}

group "build" {
  targets = ["console-build", "core-build"]
}

group "output" {
  targets = ["console-output", "core-output"]
}

group "lint" {
  targets = ["console-lint", "core-lint"]
}

group "image" {
  targets = ["console-image", "core-image"]
}

target "_console" {
  dockerfile = "Dockerfile.console"
  context = "."
  contexts = {
    node = "docker-image://node:22"
    nginx = "docker-image://nginx:1.27-alpine"
  }
}

target "console" {
  name     = "console-${tgt}"
  inherits = ["_console"]
  matrix = {
    tgt = ["build", "output", "lint", "image"]
  }
  output = {
    "build"  = ["type=cacheonly"]
    "output" = ["type=local,dest=.build/console"]
    "lint"   = ["type=cacheonly"]
    "image"   = ["type=docker"]
  }[tgt]
  tags = {
    "build"  = []
    "output" = []
    "lint"   = []
    "image"   = ["${REGISTRY}/console:${GITHUB_SHA}"]
  }[tgt]
  cache-to = {
    "build"  =  ["type=gha,ignore-error=true,mode=max,scope=${tgt}"]
    "output" =  ["type=gha,ignore-error=true,mode=max,scope=${tgt}"]
    "lint"   =  ["type=gha,ignore-error=true,mode=max,scope=${tgt}"]
    "image"   = ["type=gha,ignore-error=true,mode=max,scope=${tgt}"]
  }[tgt]
    cache-from = {
    "build"  =  ["type=gha,ignore-error=true,mode=max,scope=${tgt}"]
    "output" =  ["type=gha,ignore-error=true,mode=max,scope=${tgt}"]
    "lint"   =  ["type=gha,ignore-error=true,mode=max,scope=${tgt}"]
    "image"   = ["type=gha,ignore-error=true,mode=max,scope=${tgt}"]
  }[tgt]
  target = tgt
}

target "_core" {
  dockerfile = "Dockerfile.core"
  context = "."
  contexts = {
    node = "docker-image://golang:1.23"
    console = "target:console-output"
  }
  args = {
    SASS_VERSION      = "1.64.1"
    GOLANG_CI_VERSION = "1.64.5"
  }
}

target "core" {
  name     = "core-${tgt}"
  inherits = ["_core"]
  matrix = {
    tgt = ["build", "output", "lint", "image"]
  }
  output = {
    "build"  = ["type=cacheonly"]
    "output" = ["type=local,dest=.build/core"]
    "lint"   = ["type=cacheonly"]
    "image"   = ["type=docker"]
  }[tgt]
  tags = {
    "build"  = []
    "output" = []
    "lint"   = []
    "image"   = ["${REGISTRY}/zitadel:${GITHUB_SHA}"]
  }[tgt]
    cache-to = {
    "build"  =  ["type=gha,ignore-error=true,mode=max,scope=${tgt}"]
    "output" =  ["type=gha,ignore-error=true,mode=max,scope=${tgt}"]
    "lint"   =  ["type=gha,ignore-error=true,mode=max,scope=${tgt}"]
    "image"   = ["type=gha,ignore-error=true,mode=max,scope=${tgt}"]
  }[tgt]
    cache-from = {
    "build"  =  ["type=gha,ignore-error=true,mode=max,scope=${tgt}"]
    "output" =  ["type=gha,ignore-error=true,mode=max,scope=${tgt}"]
    "lint"   =  ["type=gha,ignore-error=true,mode=max,scope=${tgt}"]
    "image"   = ["type=gha,ignore-error=true,mode=max,scope=${tgt}"]
  }[tgt]
  target = tgt
}