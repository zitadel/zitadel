#!/bin/bash

set -e

./build/operator/prebuild.sh "./migrations"
./build/operator/build.sh
exec ./zitadelctl "$@"
