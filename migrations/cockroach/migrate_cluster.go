package migrations

//go:generate flyway -url=jdbc:postgresql://cockroachdb-public:26257/defaultdb -user=root -password= -locations=filesystem:./ migrate
