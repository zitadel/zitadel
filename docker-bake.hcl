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

target "console" {
  name     = "console-${tgt}"
  inherits = ["_console"]
  matrix = {
    tgt = ["build", "lint", "image", "generate"]
  }
  output = {
    "build"  = ["type=local,dest=.build/console"]
    "lint"   = ["type=cacheonly"]
    "image"   = ["type=image"]
    "generate" = ["type=local,dest=./"]
  }[tgt]
  tags = {
    "build"  = []
    "lint"   = []
    "image"   = ["${REGISTRY}/console:${GITHUB_SHA}"]
    "generate" = []
  }[tgt]
  cache-to = {
    "build"  =  ["type=gha,ignore-error=true,mode=max,scope=console-${tgt}"]
    "lint"   =  ["type=gha,ignore-error=true,mode=max,scope=console-${tgt}"]
    "image"   = ["type=gha,ignore-error=true,mode=max,scope=console-${tgt}"]
    "generate"   = ["type=gha,ignore-error=true,mode=max,scope=console-${tgt}"]
  }[tgt]
    cache-from = {
    "build"  =  ["type=gha,scope=console-${tgt}"]
    "lint"   =  ["type=gha,scope=console-${tgt}"]
    "image"   = ["type=gha,scope=console-${tgt}"]
    "generate" = ["type=gha,scope=console-${tgt}"]
  }[tgt]
  platforms = {
    "build"  =  []
    "lint"   =  []
    "unit"   =  []
    "image"   = ["linux/amd64", "linux/arm64"]
    "generate"   = []
  }[tgt]
  target = tgt
}

target "_core" {
  dockerfile = "Dockerfile.core"
  context = "."
  contexts = {
    golang = "docker-image://golang:1.24"
    console = "target:console-build"
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
    tgt = ["build", "lint", "image", "generate", "unit"]
  }
  output = {
    "build"  = ["type=local,dest=.build/core"]
    "lint"   = ["type=cacheonly"]
    "unit"   = ["type=local,dest=./.build/core"]
    "image"   = ["type=image"]
    "generate" = ["type=local,dest=./"]
  }[tgt]
  tags = {
    "build"  = []
    "lint"   = []
    "unit"   = []
    "image"   = ["${REGISTRY}/zitadel:${GITHUB_SHA}"]
    "generate"   = []
  }[tgt]
    cache-to = {
    "build"  =  ["type=gha,ignore-error=true,mode=max,scope=core-${tgt}"]
    "lint"   =  ["type=gha,ignore-error=true,mode=max,scope=core-${tgt}"]
    "unit"   =  ["type=gha,ignore-error=true,mode=max,scope=core-${tgt}"]    
    "image"   = ["type=gha,ignore-error=true,mode=max,scope=core-${tgt}"]
    "generate"   = ["type=gha,ignore-error=true,mode=max,scope=core-${tgt}"]
  }[tgt]
    cache-from = {
    "build"  =  ["type=gha,scope=core-${tgt}"]
    "lint"   =  ["type=gha,scope=core-${tgt}"]
    "unit"   =  ["type=gha,scope=core-${tgt}"]
    "image"   = ["type=gha,scope=core-${tgt}"]
    "generate"   = ["type=gha,scope=core-${tgt}"]
  }[tgt]
    platforms = {
    "build"  =  []
    "lint"   =  []
    "unit"   =  []
    "image"   = ["linux/amd64", "linux/arm64"]
    "generate"   = []
  }[tgt]
  target = tgt
}