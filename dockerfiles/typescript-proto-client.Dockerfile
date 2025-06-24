FROM login-pnpm AS typescript-proto-client
COPY ./login/packages/zitadel-proto/package.json ./packages/zitadel-proto/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter zitadel-proto
COPY --from=proto-files /buf.yaml /buf.lock /proto-files/
COPY --from=proto-files /zitadel /proto-files/zitadel
COPY ./login/packages/zitadel-proto/buf.gen.yaml ./packages/zitadel-proto/
RUN ls /proto-files && cd packages/zitadel-proto && pnpm exec buf generate /proto-files

FROM scratch AS typescript-proto-client-out
COPY --from=typescript-proto-client /build/packages/zitadel-proto/zitadel /zitadel
COPY --from=typescript-proto-client /build/packages/zitadel-proto/google /google
COPY --from=typescript-proto-client /build/packages/zitadel-proto/protoc-gen-openapiv2 /protoc-gen-openapiv2
COPY --from=typescript-proto-client /build/packages/zitadel-proto/validate /validate
