#!/bin/sh

case $@ in 
    sh*)
        ${@:3}
        ;;
    bash*)
        ${@:5}
        ;;
    *)
        /app/zitadel $(echo $ZITADEL_ARGS)
        ;;
esac
