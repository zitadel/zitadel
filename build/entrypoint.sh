#!/bin/sh

case $@ in 
    sh*)
    $@
    exit $?
esac

./zitadel $(echo $ZITADEL_ARGS)