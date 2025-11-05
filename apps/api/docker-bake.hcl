
target "release" {}

target "api" {
    inherits = [ "release" ]
    context = "."
    dockerfile = "apps/api/Dockerfile"
    tags = [ "zitadel-api:local"]
    annotations = [
        "org.opencontainers.image.description=ZITADEL API - Identity infrastructure, simplified for you."
    ]
}

target "api-debug" {
    inherits = [ "release", "api" ]
    target = "builder"
    tags = [ "zitadel-api-debug:local" ]
}
