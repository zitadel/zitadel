
## local development

### start cockroach in docker

```bash
docker rm -f zitadel-db &&
rm -rf ${GOPATH}/src/github.com/caos/zitadel/cockroach-data &&
docker run -d \
--name=zitadel-db \
--hostname=zitadel-db \
-p 26257:26257 -p 8080:8080  \
-v "${GOPATH}/src/github.com/caos/zitadel/cockroach-data/zitadel1:/cockroach/cockroach-data"  \
cockroachdb/cockroach:v19.2.2 start --insecure
``` 

### local database migrations

#### local migrate

`go generate $GOPATH/src/github.com/caos/zitadel/migrations/cockroach/migrate_local.go`

#### local cleanup

`go generate $GOPATH/src/github.com/caos/zitadel/migrations/cockroach/clean_local.go`


### Connect to Cockroach

`docker exec -it "zitadel-db" /cockroach/cockroach sql --insecure`

#### Should show eventstore, management, admin, auth
`show databases;`

