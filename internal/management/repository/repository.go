package repository

type Repository interface {
	Health() error
	ProjectRepository
	PolicyRepository
	OrgRepository
	UserRepository
	UserGrantRepository
	IamRepository
}
