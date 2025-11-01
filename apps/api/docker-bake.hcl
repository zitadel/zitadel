
target "release" {}

target "api" {
    inherits = [ "release" ]
    context = "."
    dockerfile = "apps/api/Dockerfile"
    tags = [ "zitadel-api:local"]
}

target "api-debug" {
    inherits = [ "release", "api" ]
    target = "builder"
    tags = [ "zitadel-api-debug:local" ]
}
