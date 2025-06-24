target "zitadel" {
  dockerfile = "build/Dockerfile"
}

target "typescript-proto-client" {
  contexts = {
    proto-files = "target:proto-files"
  }
  output = [
    "type=local,dest=login/packages/zitadel-proto"
  ]
}

target "typescript-proto-client-out" {
  output = [
    "type=local,dest=login/packages/zitadel-proto"
  ]
}

