env=$1
bucket=$2
db=$3
folder=$4
safile=$5
certs=$6
name=$7

filenamelocal=zitadel-${db}.sql
filenamebucket=zitadel-${db}-${name}.sql

/cockroach/cockroach.sh dump --dump-mode=data --certs-dir=${certs} --host=cockroachdb-public:26257 ${db} > ${folder}/${filenamelocal}
curl -X POST \
    -H "$(oauth2l header --json ${safile} cloud-platform)" \
    -H "Content-Type: application/json" \
    --data-binary @${folder}/${filenamelocal} \
    "https://storage.googleapis.com/upload/storage/v1/b/${bucket}/o?uploadType=media&name=${env}/${name}/${filenamebucket}"