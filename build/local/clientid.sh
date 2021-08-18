#!/bin/bash

# ------------------------------
# sets the client id in environment.json
# ------------------------------

clientid=""
while [ -z $clientid ]; do
    echo "no from zitadel ==> retrying in 5 seconds"
    sleep 5
    clientid=$(curl -s http://${HOST}:${PORT}/clientID)
    if [[ "$clientid" != *@zitadel* ]]; then
        echo "invalid response from zitadel ==> retrying in 5 seconds"
        clientid=""
    fi
done

echo "$(jq ".clientid = $clientid" /environment.json)" > environment.json