variable "login-context" {
  default = "./login"
}

group "default" {
  targets = ["login-docker-image-local-protos"]
}

target "typescript-base" {
  context = "./login"
}

target "typescript-proto" {
  dockerfile = "bake/typescript-proto.Dockerfile"
  output = ["type=local,dest=./login/packages/zitadel-proto"]
  contexts = {
    typescript-base    = "target:typescript-base"
    proto = "./proto"
  }
}

target "login-docker-image-local-protos" {
  inherits = ["login-docker-image"]
  contexts = {
    proto = "target:typescript-proto"
  }
}
