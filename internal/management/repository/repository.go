package repository

type Repository interface {
	Health() error
	ProjectRepository
	PolicyRepository
	UserRepository
	UserGrantRepository
}
