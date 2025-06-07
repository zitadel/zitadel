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

DIGEST_FILE="$CACHE_DIR/$MAKE_TARGET.digests"
mkdir -p "$CACHE_DIR"

get_image_ids() {
	local ids=""
	for img in $(echo "$IMAGES" | tr ';' ' '); do
		local id=$(docker image inspect "$img" --format='{{.Id}}' 2>/dev/null || true)
		id=${id:-new-or-error}
		ids="${ids}${id};"
	done
	ids=${ids%;}  # Remove trailing semicolon
	echo "$ids"
}

OLD_DIGEST=$(cat "$DIGEST_FILE" 2>/dev/null || echo "")
OLD_STATUS=$(echo "$OLD_DIGEST" | cut -d ';' -f1)
OLD_IDS=$(echo "$OLD_DIGEST" | cut -d ';' -f2-9)
if [[ "$OLD_IDS" == "$(get_image_ids)" ]]; then
    echo "Skipping $MAKE_TARGET – all images unchanged, returning cached status $OLD_STATUS"
    exit $OLD_STATUS
else
    echo "Running $MAKE_TARGET..."
    set +e
    make $MAKE_TARGET
    STATUS=$?
    set -e
    echo "${STATUS};$(get_image_ids)" > $DIGEST_FILE
    exit $STATUS
fi
