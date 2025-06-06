variable "release_tags" {
  default = ["zitadel-login:local"]
}

group "default" {
  targets = ["login-generate"]
}

target "login-base" {
  context = "."
  dockerfile = "dockerfiles/login-base.Dockerfile"
}

target "download-protos" {
    dockerfile = "dockerfiles/download-protos.Dockerfile"
    contexts = {
      base = "target:login-base"
    }
}

target "core-mock" {
  dockerfile = "dockerfiles/core-mock.Dockerfile"
  contexts = {
    protos = "target:download-protos"
  }
}

target "login-generate" {
  dockerfile = "dockerfiles/login-generate.Dockerfile"
  contexts = {
    base = "target:login-base"
  }
}

target "login-image" {
  dockerfile = "dockerfiles/login-image.Dockerfile"
  tags = "${release_tags}"
  args = {
    NODE_ENV = "production"
  }
  contexts = {
      generated = "target:login-generate"
  }
}
