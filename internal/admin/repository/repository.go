package repository

type Repository interface {
	Health() error
	OrgRepository
}
