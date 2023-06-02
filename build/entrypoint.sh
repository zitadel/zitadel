#!/bin/sh

if [[ $(echo $ZITADEL_ARGS) == sh* ]]; then
    $ZITADEL_ARGS
fi

/app/zitadel $(echo $ZITADEL_ARGS)