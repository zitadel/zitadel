//+build ignore

package migrations

//go:generate flyway -url=jdbc:sqlite:/Users/silvanreusser/go/src/github.com/caos/zitadel/.local/zitadel.db -user=admin -password= -schemas=eventstore, -locations=filesystem:./ migrate
