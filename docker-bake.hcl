variable "GITHUB_SHA" {
  default = "latest"
}

variable "REGISTRY" {
  default = "ghcr.io/fforootd"
}

group "generate" {
  targets = ["console-base", "core-base"]
}

group "build" {
  targets = ["console-build", "core-build"]
}

group "output" {
  targets = ["console-output", "core-output"]
}

group "unit-test" {
  targets = ["core-unit-test"]
}

group "lint" {
  targets = ["console-lint", "core-lint"]
}

target "console-base" {
  target = "console-base"
}

target "console-build" {
  target = "console-build"
}

target "console-lint" {
  target = "console-lint"
}

target "console-image" {
  target = "console-image"
  tags = [
    "${REGISTRY}/console:${GITHUB_SHA}",
  ]
}

target "console-output" {
  target = "console-output"
  output = ["type=local,dest=.build/console"]
}

target "core-base" {
  target = "core-base"
}

target "core-build" {
  target = "core-build"
}

target "core-lint" {
  target = "core-lint"
}

target "core-image" {
  target = "core-image"
  tags = [
    "${REGISTRY}/zitadel:${GITHUB_SHA}",
  ]
}

target "core-output" {
  target = "core-output"
  output = ["type=local,dest=.build/core"]
}

target "core-unit-test" {
  target = "core-unit-test"
}