target "release" {}

target "login" {
    inherits = [ "release" ]
    context = "apps/login"
    tags = [ "zitadel-login:local"]
    annotations = [
        "org.opencontainers.image.description=ZITADEL Login - Identity infrastructure, simplified for you."
    ]
}
