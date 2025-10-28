variable "ZITADEL_RELEASE_GITHUB_ORG" {
    type   = string
}

target "api" {
    tags = [
        "ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/api:${ZITADEL_RELEASE_VERSION}",
        ZITADEL_RELEASE_IS_LATEST ? "ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/api:latest": "",
    ]
}

target "api-debug" {
    tags = [
        "ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/api:${ZITADEL_RELEASE_VERSION}-debug",
        ZITADEL_RELEASE_IS_LATEST ? "ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/api:latest-debug": "",
    ]
}
