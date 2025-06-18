target "docker-metadata-action" {}

target "login-pnpm" {
  cache-to = [
    { type: "registry", ref: "${IMAGE_REGISTRY}/login-pnpm-buildcache:${REF_TAG}", mode: "max", oci-mediatypes=true }
  ]
}

target "login-dev-base" {
  cache-to = [
    { type: "registry", ref: "${IMAGE_REGISTRY}/login-dev-base-buildcache:${REF_TAG}", mode: "max", oci-mediatypes=true }
  ]
}

target "login-lint" {
  cache-to = [
    { type: "registry", ref: "${IMAGE_REGISTRY}/login-lint-buildcache:${REF_TAG}", mode: "max", oci-mediatypes=true }
  ]
}

target "login-test-unit" {
  cache-to = [
    { type: "registry", ref: "${IMAGE_REGISTRY}/login-test-unit-buildcache:${REF_TAG}", mode: "max", oci-mediatypes=true }
  ]
}

target "login-test-integration" {
  cache-to = [
    { type: "registry", ref: "${IMAGE_REGISTRY}/login-test-integration-buildcache:${REF_TAG}", mode: "max", oci-mediatypes=true }
  ]
}

target "login-client" {
  cache-to = [
    { type: "registry", ref: "${IMAGE_REGISTRY}/login-client-buildcache:${REF_TAG}", mode: "max", oci-mediatypes=true }
  ]
}

target "login-test-acceptance" {
  cache-to = [
    { type: "registry", ref: "${IMAGE_REGISTRY}/login-test-acceptance-buildcache:${REF_TAG}", mode: "max", oci-mediatypes=true }
  ]
}

target "login-standalone" {
  cache-to = [
    { type: "registry", ref: "${IMAGE_REGISTRY}/login-buildcache:${REF_TAG}", mode: "max", oci-mediatypes=true }
  ]
}
