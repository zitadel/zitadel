FROM login-pnpm AS typescript-proto-client
COPY packages/zitadel-proto/package.json ./packages/zitadel-proto/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter zitadel-proto
COPY packages/zitadel-proto ./packages/zitadel-proto
RUN pnpm generate

FROM scratch AS typescript-proto-client-out
COPY --from=typescript-proto-client /build/packages/zitadel-proto/zitadel /zitadel
COPY --from=typescript-proto-client /build/packages/zitadel-proto/google /google
COPY --from=typescript-proto-client /build/packages/zitadel-proto/protoc-gen-openapiv2 /protoc-gen-openapiv2
COPY --from=typescript-proto-client /build/packages/zitadel-proto/validate /validate

FROM typescript-proto-client
