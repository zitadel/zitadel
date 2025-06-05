FROM typescript-base

COPY --from=proto . /proto

RUN ls -la /proto
WORKDIR /app/packages/zitadel-proto

RUN pnpm exec buf generate /proto
