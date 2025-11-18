variable "ZITADEL_RELEASE_GITHUB_ORG" {
    type   = string
}

target "api" {
    tags = [
        "ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/zitadel:${ZITADEL_RELEASE_VERSION}",
        ZITADEL_RELEASE_IS_LATEST ? "ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/zitadel:latest": "",
    ]
}

target "api-debug" {
    tags = [
        "ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/zitadel-debug:${ZITADEL_RELEASE_VERSION}",
        ZITADEL_RELEASE_IS_LATEST ? "ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/zitadel-debug:latest": "",
    ]
}
