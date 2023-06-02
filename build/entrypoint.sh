#!/bin/sh

case $@ in 
    sh*)
    $@
    exit $?
esac

/app/zitadel $(echo $ZITADEL_ARGS)