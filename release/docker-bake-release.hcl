variable "ZITADEL_RELEASE_IS_LATEST" {
    type    = bool
    default = false
}

variable "ZITADEL_RELEASE_VERSION" {
    type    = string
}

variable "ZITADEL_RELEASE_REVISION" {
    type    = string
}

variable "ZITADEL_RELEASE_PUSH" {
    type    = bool
    default = false
}

target "release-common" {
    platforms = [ "linux/amd64", "linux/arm64" ]
    push = ZITADEL_RELEASE_PUSH
    labels = {
        "org.opencontainers.image.created" = timestamp()
        "org.opencontainers.image.version" = ZITADEL_RELEASE_VERSION
        "org.opencontainers.image.revision" = ZITADEL_RELEASE_REVISION
    }
}
