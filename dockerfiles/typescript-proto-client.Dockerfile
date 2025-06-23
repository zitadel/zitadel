FROM login-pnpm AS typescript-proto-client
COPY ./login/packages/zitadel-proto/package.json ./packages/zitadel-proto/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter zitadel-proto
COPY --from=proto-files / ./packages/zitadel-proto/proto
COPY ./login/packages/zitadel-proto/buf.gen.yaml ./packages/zitadel-proto/
RUN ls -la packages/zitadel-proto
RUN cd packages/zitadel-proto && pnpm exec buf generate . --path ./proto/zitadel
