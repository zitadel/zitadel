variable "tags" {
  default = ["zitadel-login:local"]
}

variable "login-context" {
  default = "."
}

group "default" {
  targets = ["login-docker-image"]
}

target "typescript-base" {
    context = "${login-context}"
    dockerfile = "bake/base.Dockerfile"
}

target "proto" {
  context = "${login-context}"
  dockerfile = "bake/proto.Dockerfile"
  output = ["type=local,dest=./packages/zitadel-proto"]
  contexts = {
    base = "target:typescript-base"
  }
}

target "login-docker-image" {
  context = "${login-context}"
  dockerfile = "bake/login-for-docker.Dockerfile"
  tags = "${tags}"
  args = {
    NODE_ENV = "production"
  }
  contexts = {
    proto = "target:proto"
  }
}
