package repository

type Repository interface {
	Health() error
	ProjectRepository
	OrgRepository
	UserRepository
	UserGrantRepository
	IamRepository
	FeaturesRepository
}
