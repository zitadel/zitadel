#!/bin/sh

case $@ in 
    sh*)
    echo "came in if"
    $@
    exit $?
esac

/app/zitadel $(echo $ZITADEL_ARGS)