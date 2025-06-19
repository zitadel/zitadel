#!/usr/bin/env bash

# Usage: ./run_or_skip.sh <Make target> <images>
# Example: ./run_or_skip.sh lint-force "img1;img2"

set -euo pipefail

if [ -z "$CACHE_DIR" ]; then
    echo "CACHE_DIR is not set. Please set it to a valid directory."
    exit 1
fi

MAKE_TARGET=$1
IMAGES=$2
FORCE=${FORCE:-false}

CACHE_FILE="$CACHE_DIR/$MAKE_TARGET.digests"
mkdir -p "$CACHE_DIR"

inspect_image() {
  local image=$1
  local format=$2
  docker image inspect "$image" --format="$format" 2>/dev/null || true
}

get_digest() {
  local image=$1
  echo "id=$(inspect_image $image '{{ .Id }}'),digest=$(inspect_image $image '{{ index RepoDigests 0 }}')"
}

get_image_digests() {
  local digests=""
	for img in $(echo "$IMAGES"); do
		local digest=$(get_digest $img)
		if [[ -z $digest ]]; then
		  docker pull "$img" >/dev/null 2>&1 || true
		  digest=$(get_digest $img)
    fi
    if [[ -z $digest ]]; then
		  digest=$(get_digest $img)
    fi
		digest="${img}@${digest}"
		digests="${digests}${digest} "
	done
	digests=${digests% }  # Remove trailing space
	echo "$digests"
}

CACHE_CONTENT=$(cat "$CACHE_FILE" 2>/dev/null || echo "")
CACHED_STATUS=$(echo "$CACHE_CONTENT" | cut -d ';' -f1)
CACHED_DIGESTS=$(echo "$CACHE_CONTENT" | cut -d ';' -f2)
CURRENT_DIGESTS="$(get_image_digests)"

echo "Comparing cached vs current image digests..."
echo
echo "$CACHED_DIGESTS"
echo
echo "$CURRENT_DIGESTS"

IMAGE_CHANGED=false
for current_digest in $CURRENT_DIGESTS; do
  current_digest_image_id=$(echo "$current_digest" | cut -d ',' -f1)
  current_digest_repo_digest=$(echo "$current_digest" | cut -d ',' -f2)
  for cached_digest in $CACHED_DIGESTS; do
    cached_digest_image_id=$(echo "$current_digest" | cut -d ',' -f1)
    cached_digest_repo_digest=$(echo "$current_digest" | cut -d ',' -f2)
    if [[ "$current_digest_image_id" != "$cached_digest_image_id" && "$current_digest_repo_digest" != "$cached_digest_repo_digest" ]]; then
      echo "Image digest mismatch:"
      echo "Current: $current_digest"
      echo "Cached:  $cached_digest"
      IMAGE_CHANGED=true
      break 2
    fi
  done
done

if [[ "$IMAGE_CHANGED" == "false" ]]; then
    if [[ "$FORCE" == "true" ]]; then
        echo "\$FORCE=$FORCE - Running $MAKE_TARGET despite unchanged images."
    else
        echo "Skipping $MAKE_TARGET – all images unchanged, returning cached status $CACHED_STATUS"
        exit $CACHED_STATUS
    fi
fi

echo "Images have changed"
echo
docker images
echo
echo "Running $MAKE_TARGET..."
set +e
make -j $MAKE_TARGET
STATUS=$?
set -e
echo "${STATUS};$(get_image_digests)" > $CACHE_FILE
exit $STATUS
