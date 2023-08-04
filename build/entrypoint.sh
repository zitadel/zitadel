#!/bin/bash

case $@ in 
    sh*)
        ${@:3}
        ;;
    bash*)
        ${@:5}
        ;;
    *)
        if [[ ! -z "$@" ]]
        then
            ZITADEL_ARGS="$@"
        fi
        /app/zitadel ${ZITADEL_ARGS}
        ;;
esac
