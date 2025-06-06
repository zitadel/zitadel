include = ["login/docker-bake.hcl"]

target "local-protos" {
  inherits = ["download-protos"]
  dockerfile = "dockerfiles/protos.Dockerfile"
}

target "login-generate" {
  inherits = ["login-generate"]
  dockerfile = "dockerfiles/login-generate.Dockerfile"
  contexts = {
    protos = "./proto
  }
}
