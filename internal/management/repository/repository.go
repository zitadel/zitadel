package repository

type Repository interface {
	Health() error
	ProjectRepository
	PolicyRepository
	OrgRepository
	OrgMemberRepository
	UserRepository
	UserGrantRepository
}
