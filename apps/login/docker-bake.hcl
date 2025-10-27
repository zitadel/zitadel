target "release" {}

target "login" {
    inherits = [ "release" ]
    context = "apps/login"
    tags = [ "ghcr.io/eliobischof/login:local"]
}
