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

DIGEST_FILE="$CACHE_DIR/$MAKE_TARGET.digests"
mkdir -p "$CACHE_DIR"

get_image_ids() {
  local ids=""
	for img in $(echo "$IMAGES"); do
		local id=$(docker image inspect "$img" --format='{{.Id}}' 2>/dev/null || true)
		if [[ -z $id ]]; then
		  docker pull "$img" >/dev/null 2>&1 || true
		  id=$(docker image inspect "$img" --format='{{.Id}}' 2>/dev/null || true)
    fi
    if [[ -z $id ]]; then
		  id=$(docker image inspect "$img" --format='{{.Id}}' 2>/dev/null || true)
    fi
		id=${id:-new-and-not-pullable-or-failed-to-build}
		id="${img}@${id}"
		ids="${ids}${id};"
	done
	ids=${ids%;}  # Remove trailing semicolon
	echo "$ids"
}

PREVIOUS_DIGEST=$(cat "$DIGEST_FILE" 2>/dev/null || echo "")
PREVIOUS_STATUS=$(echo "$PREVIOUS_DIGEST" | cut -d ';' -f1)
PREVIOUS_IMAGE_IDS=$(echo "$PREVIOUS_DIGEST" | cut -d ';' -f2-99)
CURRENT_IMAGE_IDS="$(get_image_ids)"
  if [[ "$PREVIOUS_IMAGE_IDS" == "$CURRENT_IMAGE_IDS" ]]; then
    if [[ "$FORCE" == "true" ]]; then
        echo "\$FORCE=$FORCE - Running $MAKE_TARGET despite unchanged images."
    else
        echo "Skipping $MAKE_TARGET – all images unchanged, returning cached status $PREVIOUS_STATUS"
        exit $PREVIOUS_STATUS
    fi
fi
echo "Images have changed"
echo
echo "PREVIOUS_IMAGE_IDS does not match CURRENT_IMAGE_IDS"
echo
echo "$PREVIOUS_IMAGE_IDS"
echo
echo "$CURRENT_IMAGE_IDS"
echo
docker images
echo
echo "Running $MAKE_TARGET..."
set +e
make -j $MAKE_TARGET
STATUS=$?
set -e
echo "${STATUS};$(get_image_ids)" > $DIGEST_FILE
exit $STATUS
