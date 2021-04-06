#! /bin/sh

if [ -n $1 ]; then
    ZITADEL_IMPORT=$1/zitadel
else
    echo "need message import"
    exit 3
fi

protoc \
    -I=/proto/include \
    -I=/proto/include/zitadel \
    --go_opt=Mzitadel/app.proto=${ZITADEL_IMPORT}/app \
    --go_opt=Mzitadel/auth_n_key.proto=${ZITADEL_IMPORT}/authn \
    --go_opt=Mzitadel/change.proto=${ZITADEL_IMPORT}/change \
    --go_opt=Mzitadel/features.proto=${ZITADEL_IMPORT}/features \
    --go_opt=Mzitadel/idp.proto=${ZITADEL_IMPORT}/idp \
    --go_opt=Mzitadel/member.proto=${ZITADEL_IMPORT}/member \
    --go_opt=Mzitadel/message.proto=${ZITADEL_IMPORT}/message \
    --go_opt=Mzitadel/object.proto=${ZITADEL_IMPORT}/object \
    --go_opt=Mzitadel/org.proto=${ZITADEL_IMPORT}/org \
    --go_opt=Mzitadel/policy.proto=${ZITADEL_IMPORT}/policy \
    --go_opt=Mzitadel/project.proto=${ZITADEL_IMPORT}/project \
    --go_opt=Mzitadel/user.proto=${ZITADEL_IMPORT}/user \
    --go_out /go/src \
    --go-grpc_out /go/src \
    $(find /proto/include/zitadel -iname *.proto)
