package database

//go:generate mockgen -source client.go -package databasemock -destination mock/client.mock.go github.com/caos/zitadel/operator/kinds/iam/zitadel/database ClientInt
