package handler

import "database/sql"

type SQLHandler interface {
	Handler
}

type sqlHandler struct {
	client *sql.DB
}
