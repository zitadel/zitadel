target "typescript-proto-client" {
  contexts = {
    proto-files = "target:proto-files"
  }
}

target "typescript-proto-client-out" {
  contexts = {
    proto-files = "target:proto-files"
  }
  output = [
    "type=local,dest=login/packages/zitadel-proto"
  ]
}
