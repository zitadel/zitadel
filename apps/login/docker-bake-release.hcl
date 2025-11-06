variable "ZITADEL_RELEASE_GITHUB_ORG" {
    type   = string
}

target "login" {
    tags = [
        "ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/zitadel-login:${ZITADEL_RELEASE_VERSION}",
        ZITADEL_RELEASE_IS_LATEST ? "ghcr.io/${ZITADEL_RELEASE_GITHUB_ORG}/zitadel-login:latest": "",
    ]
}
