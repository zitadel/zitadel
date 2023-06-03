#!/bin/sh

case $@ in 
    sh* | bash*)
    $@
    exit $?
esac

/app/zitadel $(echo $ZITADEL_ARGS)