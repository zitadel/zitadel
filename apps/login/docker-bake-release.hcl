target "login" {
    tags = [
        "ghcr.io/eliobischof/login:${ZITADEL_RELEASE_VERSION}",
        ZITADEL_RELEASE_IS_LATEST ? "ghcr.io/eliobischof/login:latest": "",
    ]
}
