#!/usr/bin/env bash

ACTION=$1
ENVFILE=$2

shift
shift

export projectRoot=".."

if [ -z ${ENVFILE:x} ]; then
    echo "Not sourcing any env file"
else
    set -a; source $ENVFILE; set +a
fi

NPX=""
if ! command -v cypress &> /dev/null; then
    NPX="npx"
fi

###############################
function zitadel-health-check {
###############################
# this is needed to ensure zitadel has finished its startup
i=1
until [ $i -gt 10 ]; do
    HEALTHSTATUS=$(curl  zitadel:8080/healthz | tail -n1)
    if  [ $HEALTHSTATUS == '{"status":"ok"}' ]; then
        echo "ZITADEL is up and running"
        break
    else echo "ZITADEL is starting" 
        sleep 1
    fi
    i=$[$i+1]; 
done
}


#install missing packages manually
npm install debug jsonwebtoken mochawesome typesript

zitadel-health-check

$NPX cypress $ACTION \
    --port "${E2E_CYPRESSPORT}" \
    --env org="${ZITADEL_E2E_ORG}",org_owner_password="${ZITADEL_E2E_ORGOWNERPW}",org_owner_viewer_password="${ZITADEL_E2E_ORGOWNERVIEWERPW}",org_project_creator_password="${ZITADEL_E2E_ORGPROJECTCREATORPW}",login_policy_user_password="${ZITADEL_E2E_LOGINPOLICYUSERPW}",password_complexity_user_password="${ZITADEL_E2E_PASSWORDCOMPLEXITYUSERPW}",baseUrl="${ZITADEL_E2E_BASEURL}",serviceAccountKeyPath="${ZITADEL_E2E_MACHINEKEYPATH}",otherZitadelIdpInstance="${ZITADEL_E2E_OTHERZITADELIDPINSTANCE}",zitadelProjectResourceId="${ZITADEL_E2E_ZITADELPROJECTRESOURCEID}" \
    "$@"
