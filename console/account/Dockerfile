FROM node:22-bookworm AS build
WORKDIR /workspace
ENV NG_BUILD_MAX_WORKERS=2
RUN apt-get update && apt-get install -y --no-install-recommends protobuf-compiler unzip curl ca-certificates && rm -rf /var/lib/apt/lists/* \
    && corepack enable && corepack prepare pnpm@10.30.3 --activate
COPY . .
RUN pnpm install --frozen-lockfile
RUN bash console/scripts/install-proto-plugins.sh
RUN pnpm --dir packages/zitadel-proto exec buf generate ../../proto && pnpm --dir packages/zitadel-client run build
RUN arch="$(uname -m)"; case "$arch" in x86_64) arch=amd64;; aarch64) arch=arm64;; esac; \
    PATH="/workspace/.artifacts/bin/linux/$arch:$PATH" pnpm --dir console exec buf generate ../proto --include-imports --include-wkt
RUN node console/account/build.mjs production

FROM nginx:stable-alpine
RUN apk add --no-cache jq
COPY --from=build /workspace/console/account/dist /usr/share/nginx/html
COPY console/account/nginx.conf.template /etc/nginx/templates/default.conf.template
COPY console/account/10-account-config.sh /docker-entrypoint.d/10-account-config.sh
RUN chmod +x /docker-entrypoint.d/10-account-config.sh
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD wget -q -O /dev/null http://127.0.0.1:8080/healthz || exit 1
