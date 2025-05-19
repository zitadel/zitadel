package database

//go:generate mockgen -typed -package dbmock -destination ./dbmock/database.mock.go github.com/zitadel/zitadel/backend/v3/storage/database Pool,Client,Row,Rows,Transaction
