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
IGNORE_RUN_CACHE=${IGNORE_RUN_CACHE:-false}

CACHE_FILE="$CACHE_DIR/$MAKE_TARGET.digests"
mkdir -p "$CACHE_DIR"

get_image_creation_dates() {
  local values=""
	for img in $(echo "$IMAGES"); do
		local value=$(docker image inspect "$img" --format='{{.Created}}' 2>/dev/null || true)
		if [[ -z $value ]]; then
		  docker pull "$img" >/dev/null 2>&1 || true
		  value=$(docker image inspect "$img" --format='{{.Created}}' 2>/dev/null || true)
    fi
    if [[ -z $value ]]; then
		  value=$(docker image inspect "$img" --format='{{.Created}}' 2>/dev/null || true)
    fi
		value=${value:-new-and-not-pullable-or-failed-to-build}
		value="${img}@${value}"
		values="${values}${value};"
	done
	values=${values%;}  # Remove trailing semicolon
	echo "$values"
}

CACHE_FILE_CONTENT=$(cat "$CACHE_FILE" 2>/dev/null || echo "")
CACHED_STATUS=$(echo "$CACHE_FILE_CONTENT" | cut -d ';' -f1)
CACHED_IMAGE_CREATED_VALUES=$(echo "$CACHE_FILE_CONTENT" | cut -d ';' -f2-99)
CURRENT_IMAGE_CREATED_VALUES="$(get_image_creation_dates)"
  if [[ "$CACHED_IMAGE_CREATED_VALUES" == "$CURRENT_IMAGE_CREATED_VALUES" ]]; then
    if [[ "$IGNORE_RUN_CACHE" == "true" ]]; then
        echo "\$IGNORE_RUN_CACHE=$IGNORE_RUN_CACHE - Running $MAKE_TARGET despite unchanged images."
    else
        echo "Skipping $MAKE_TARGET – all images unchanged, returning cached status $CACHED_STATUS"
        exit $CACHED_STATUS
    fi
fi
echo "Images have changed"
echo
echo "CACHED_IMAGE_CREATED_VALUES does not match CURRENT_IMAGE_CREATED_VALUES"
echo
echo "$CACHED_IMAGE_CREATED_VALUES"
echo
echo "$CURRENT_IMAGE_CREATED_VALUES"
echo
docker images
echo
echo "Running $MAKE_TARGET..."
set +e
make -j $MAKE_TARGET
STATUS=$?
set -e
echo "${STATUS};$(get_image_creation_dates)" > $CACHE_FILE
exit $STATUS
