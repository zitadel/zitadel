FROM scratch AS login-zitadel-proto-out
COPY --from=login-zitadel-proto /build/packages/zitadel-proto/zitadel /zitadel
COPY --from=login-zitadel-proto /build/packages/zitadel-proto/google /google
COPY --from=login-zitadel-proto /build/packages/zitadel-proto/protoc-gen-openapiv2 /protoc-gen-openapiv2
COPY --from=login-zitadel-proto /build/packages/zitadel-proto/validate /validate
