variable "ZITADEL_RELEASE_VERSION" {
    type    = string
}

variable "ZITADEL_RELEASE_REVISION" {
    type    = string
}

variable "ZITADEL_RELEASE_IS_LATEST" {
    type    = bool
    default = false
}

variable "ZITADEL_RELEASE_GITHUB_REPO" {
    type    = string
}

target "release" {
    platforms = [ "linux/amd64", "linux/arm64" ]
    labels = {
        "org.opencontainers.image.created" = timestamp()
        "org.opencontainers.image.version" = ZITADEL_RELEASE_VERSION
        "org.opencontainers.image.revision" = ZITADEL_RELEASE_REVISION
        "org.opencontainers.image.source" = "https://github.com/${ZITADEL_RELEASE_GITHUB_REPO}"
    }
}
