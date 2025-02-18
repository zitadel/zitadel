variable "GITHUB_SHA" {
  default = "latest"
}

variable "REGISTRY" {
  default = "ghcr.io/zitadel"
}

group "generate" {
  targets = ["console-base"]
}

group "unit-test" {
  targets = ["core-unit-test"]
}

target "console-base" {
  target = "console-base"
  cache-from = ["type=gha,scope=console-base"]
  cache-to = ["type=gha,mode=max,scope=console-base"]
}

target "console-builder" {
  target = "console-builder"
  cache-from = ["type=gha,scope=console-builder"]
  cache-to = ["type=gha,mode=max,scope=console-builder"]
}

target "console" {
  target = "console"
  tags = [
    "${REGISTRY}/console:${GITHUB_SHA}",
  ]
  cache-from = ["type=gha,scope=console"]
  cache-to = ["type=gha,mode=max,scope=console"]
}

target "core-base" {
  target = "core-base"
  cache-from = ["type=gha,scope=core-base"]
  cache-to = ["type=gha,mode=max,scope=core-base"]
}

target "core-unit-test" {
  target = "core-unit-test"
  cache-from = ["type=gha,scope=core-unit-test"]
  cache-to = ["type=gha,mode=max,scope=core-unit-test"]
}