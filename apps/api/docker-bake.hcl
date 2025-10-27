
target "release" {}

target "api" {
    inherits = [ "release" ]
    context = "."
    dockerfile = "apps/api/Dockerfile"
    tags = [ "ghcr.io/eliobischof/api:local"]
}

target "api-debug" {
    inherits = [ "release", "api" ]
    target = "builder"
    tags = [ "ghcr.io/eliobischof/api:local-debug" ]
}
