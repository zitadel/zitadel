package repository

type Repository interface {
	ProjectRepository
	OrgRepository
	UserRepository
	IamRepository
}
