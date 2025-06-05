FROM typescript-base

COPY --from=proto . /proto

RUN cd /app/packages/zitadel-proto && pnpm exec buf generate /proto
