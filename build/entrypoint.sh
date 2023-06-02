#!/bin/sh

if [[ $@ == sh* ]]; then
    $@
    exit $?
fi

/app/zitadel $(echo $ZITADEL_ARGS)