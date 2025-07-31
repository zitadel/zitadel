variable "LOGIN_TAG" {
  default = "zitadel-login:local"
}

group "default" {
  targets = ["login-standalone"]
}

target "docker-metadata-action" {
  # In the pipeline, this target is overwritten by the docker metadata action.
  tags = ["${LOGIN_TAG}"]
}

# We run integration and acceptance tests against the next standalone server for docker.
target "login-standalone" {
  inherits = [
    "docker-metadata-action",
  ]
}
