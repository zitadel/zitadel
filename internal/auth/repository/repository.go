package repository

type Repository interface {
	Health() error
	UserAgentRepository
	UserRepository
}
