#!/bin/sh

case $@ in 
    sh*)
        ${@:3}
        ;;
    bash*)
        ${@:5}
        ;;
    *)
        /app/zitadel $([[ ! -z $@ ]] && echo $@ || echo ${ZITADEL_ARGS})
        ;;
esac
