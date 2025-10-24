#!/bin/bash

set -ex

for os in linux darwin windows; do
    for arch in amd64 arm64; do
        echo "Releasing for $os-$arch..."
        GOOS=$os GOARCH=$arch pnpm nx run @zitadel/api:pack
  done
done

gh release upload ${ZITADEL_VERSION} ./.artifacts/pack/zitadel-*-*.tar.gz
docker buildx build --push --platform linux/amd64,linux/arm64 -t ghcr.io/zitadel/zitadel:v${ZITADEL_VERSION}-debug --target builder ./apps/api
docker buildx build --push --platform linux/amd64,linux/arm64 -t ghcr.io/zitadel/zitadel:v${ZITADEL_VERSION} ./apps/api
