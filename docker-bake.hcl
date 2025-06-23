target "proto-files" {
  context = "./"
  dockerfile = "dockerfiles/proto-files.Dockerfile"
}

target "typescript-proto-client" {
  context = "./"
  dockerfile = "dockerfiles/typescript-proto-client.Dockerfile"
  contexts = {
    proto-files = "target:proto-files"
  }
}
