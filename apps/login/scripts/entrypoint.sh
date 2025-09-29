#!/bin/sh

if [ -f /.env-file/.env ]; then
    set -o allexport
    . /.env-file/.env
    set +o allexport
fi

exec $@
