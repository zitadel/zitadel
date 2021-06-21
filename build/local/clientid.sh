#!/bin/bash

# ------------------------------
# sets the client id in environment.json
# ------------------------------

clientid=""
while [ -z $clientid ]; do
    echo "no from zitadel ==> retry"
    sleep 2
    clientid=$(curl -s http://${HOST}:${PORT}/clientID)
    if [[ "$clientid" != *@zitadel* ]]; then
        echo "invalid response from zitadel ==> retry"
        clientid=""
    fi
done

echo "$(jq ".clientid = $clientid" /environment.json)" > environment.json