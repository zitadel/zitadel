//+build ignore

package migrations

//go:generate flyway -url=jdbc:postgresql://localhost:26257/defaultdb -user=root -password= -locations=filesystem:./ migrate
