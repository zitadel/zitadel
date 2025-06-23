include = ["login/docker-bake.hcl"]

target "proto-files" {
  dockerfile = "dockerfiles/proto-files.Dockerfile"
}

target "typescript-proto-client" {
  dockerfile = "dockerfiles/typescript-proto-client.Dockerfile"
  contexts = {
    proto-files = "target:proto-files"
  }
}
