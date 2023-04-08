#!/bin/sh

set -eux

echo "Init deps"

( cd console ; npm ci )
( cd docs ; yarn install )

echo "Done deps"
exit 0