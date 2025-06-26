FROM scratch AS typescript-proto-client-out
COPY --from=typescript-proto-client /build/packages/zitadel-proto/zitadel /zitadel
COPY --from=typescript-proto-client /build/packages/zitadel-proto/google /google
COPY --from=typescript-proto-client /build/packages/zitadel-proto/protoc-gen-openapiv2 /protoc-gen-openapiv2
COPY --from=typescript-proto-client /build/packages/zitadel-proto/validate /validate
