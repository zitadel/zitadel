target "docker-metadata-action" {}

target "login-pnpm" {
  cache-to = [
    { "type": "registry", "ref": "${IMAGE_REGISTRY}/login-pnpm-buildcache:${BUILD_CACHE_KEY}", "mode": "max" }
  ]
#  output = [
#    { "type" : "image", "name": "${IMAGE_REGISTRY}/login-pnpm:${BUILD_CACHE_KEY}", push: true },
#  ]
}

target "login-dev-base" {
  cache-to = [
    { "type": "registry", "ref": "${IMAGE_REGISTRY}/login-dev-base-buildcache:${BUILD_CACHE_KEY}", "mode": "max" }
  ]
#  output = [
#    { "type" : "image", "name": "${IMAGE_REGISTRY}/login-dev-base:${BUILD_CACHE_KEY}", push: true },
#  ]
}

target "login-lint" {
  cache-to = [
    { "type": "registry", "ref": "${IMAGE_REGISTRY}/login-lint-buildcache:${BUILD_CACHE_KEY}", "mode": "max" }
  ]
#  output = [
#    { "type" : "image", "name": "${IMAGE_REGISTRY}/login-lint:${BUILD_CACHE_KEY}", push: true },
#  ]
}
