package repository

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type user struct {
	database.QueryExecutor
}

func User(client database.QueryExecutor) domain.UserRepository {
	// return &user{QueryExecutor: client}
	return nil
}

// On implements [domain.UserRepository].
func (exec *user) On(clauses ...domain.UserClause) domain.UserOperation {
	return &userOperation{
		QueryExecutor: exec.QueryExecutor,
		clauses:       clauses,
	}
}

// OnHuman implements [domain.UserRepository].
func (exec *user) OnHuman(clauses ...domain.UserClause) domain.HumanOperation {
	return &humanOperation{
		userOperation: *exec.On(clauses...).(*userOperation),
	}
}

// OnMachine implements [domain.UserRepository].
func (exec *user) OnMachine(clauses ...domain.UserClause) domain.MachineOperation {
	return &machineOperation{
		userOperation: *exec.On(clauses...).(*userOperation),
	}
}

// var _ domain.UserRepository = (*user)(nil)
