package db

//go:generate mockgen -source client.go -package databasemock -destination mock/client.mock.go github.com/caos/zitadel/pkg/databases/db ClientInt
