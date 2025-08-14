variable "LOGIN_TAG" {
  default = "zitadel-login:local"
}

group "default" {
  targets = ["login-standalone"]
}

# The release target is overwritten in docker-bake-release.hcl
# It makes sure the image is built for multiple platforms.
# By default the platforms property is empty, so images are only built for the current bake runtime platform.
target "release" {}

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
}
