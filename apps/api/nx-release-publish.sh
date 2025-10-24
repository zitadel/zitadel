#!/bin/bash

set -ex

for os in linux darwin windows; do
    for arch in amd64 arm64; do
        echo "Releasing for $os-$arch..."
        GOOS=$os GOARCH=$arch pnpm nx run @zitadel/api:pack
  done
done

gh release upload ${ZITADEL_VERSION} ./.artifacts/pack/zitadel-*-*.tar.gz
docker buildx build  --file ./apps/api/Dockerfile --push --platform linux/amd64,linux/arm64 -t ghcr.io/eliobischof/api:${ZITADEL_VERSION}-debug --target builder .
docker buildx build  --file ./apps/api/Dockerfile --push --platform linux/amd64,linux/arm64 -t ghcr.io/eliobischof/api:${ZITADEL_VERSION} .
