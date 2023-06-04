#!/bin/sh

case $@ in 
    sh* | bash*)
        $@
        ;;
    *)
        /app/zitadel $(echo $ZITADEL_ARGS)
        ;;
esac
