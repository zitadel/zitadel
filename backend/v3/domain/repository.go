package domain

import "github.com/zitadel/zitadel/backend/v3/storage/database"

// Repository is the base interface for all repositories.
type Repository interface {
	// PrimaryKeyColumns returns the columns for the primary key fields
	PrimaryKeyColumns() []database.Column
}
