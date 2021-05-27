#!/bin/sh

# ----------------------------------------------------------------
# generates necessary ZITADEL keys
# ----------------------------------------------------------------

set -e


KEY_PATH=$(echo "/zitadel/$(dirname ${ZITADEL_KEY_PATH})")
KEY_FILE=${KEY_PATH}/local_keys.yaml

mkdir -p ${KEY_PATH}
if [ ! -f ${KEY_FILE} ]; then
    touch ${KEY_FILE}
fi

for key in $(env | grep "ZITADEL_.*_KEY" | cut -d'=' -f2); do
    if [ $(grep -L ${key} ${KEY_FILE}) ]; then
        echo -e "${key}: $(head -c22 /dev/urandom | base64)" >> ${KEY_FILE}
    fi
done
