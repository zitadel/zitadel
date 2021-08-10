certs=$1

/cockroach/cockroach.sh sql --certs-dir=${certs} --host=cockroachdb-public:26257 -e "TRUNCATE defaultdb.flyway_schema_history;"