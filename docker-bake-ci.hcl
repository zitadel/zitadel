target "docker-metadata-action" {}

target "login-pnpm" {
  cache-to = [
    { type: "registry", ref: "${IMAGE_REGISTRY}/login-pnpm-buildcache:${BUILD_CACHE_KEY}", mode: "max", oci-mediatypes=true }
  ]
}

target "login-dev-base" {
  cache-to = [
    { type: "registry", ref: "${IMAGE_REGISTRY}/login-dev-base-buildcache:${BUILD_CACHE_KEY}", mode: "max", oci-mediatypes=true }
  ]
}

target "login-lint" {
  cache-to = [
    { type: "registry", ref: "${IMAGE_REGISTRY}/login-lint-buildcache:${BUILD_CACHE_KEY}", mode: "max", oci-mediatypes=true }
  ]
}

target "login-test-unit" {
  cache-to = [
    { type: "registry", ref: "${IMAGE_REGISTRY}/login-test-unit-buildcache:${BUILD_CACHE_KEY}", mode: "max", oci-mediatypes=true }
  ]
}

target "login-test-integration" {
  cache-to = [
    { type: "registry", ref: "${IMAGE_REGISTRY}/login-test-integration-buildcache:${BUILD_CACHE_KEY}", mode: "max", oci-mediatypes=true }
  ]
}

target "login-client" {
  cache-to = [
    { type: "registry", ref: "${IMAGE_REGISTRY}/login-client-buildcache:${BUILD_CACHE_KEY}", mode: "max", oci-mediatypes=true }
  ]
}

target "login-test-acceptance" {
  cache-to = [
    { type: "registry", ref: "${IMAGE_REGISTRY}/login-test-acceptance-buildcache:${BUILD_CACHE_KEY}", mode: "max", oci-mediatypes=true }
  ]
}

target "login-standalone" {
  cache-to = [
    { type: "registry", ref: "${IMAGE_REGISTRY}/login-buildcache:${BUILD_CACHE_KEY}", mode: "max", oci-mediatypes=true }
  ]
}
