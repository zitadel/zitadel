package core

//go:generate mockgen -source current.go -package coremock -destination mock/current.mock.go github.com/caos/internal/operator/database/kinds/databases/core DatabaseCurrent
