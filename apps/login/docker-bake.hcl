target "release" {}

target "login" {
    inherits = [ "release" ]
    context = "apps/login"
    tags = [ "zitadel-login:local"]
}
