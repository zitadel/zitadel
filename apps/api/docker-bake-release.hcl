target "api" {
    tags = [
        "ghcr.io/eliobischof/api:${ZITADEL_RELEASE_VERSION}",
        ZITADEL_RELEASE_IS_LATEST ? "ghcr.io/eliobischof/api:latest": "",
    ]
}

target "api-debug" {
    tags = [
        "ghcr.io/eliobischof/api:${ZITADEL_RELEASE_VERSION}-debug",
        ZITADEL_RELEASE_IS_LATEST ? "ghcr.io/eliobischof/api:latest-debug": "",
    ]
}
