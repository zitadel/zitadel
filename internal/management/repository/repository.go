package repository

type Repository interface {
	Health() error
	ProjectRepository
	OrgRepository
	OrgMemberRepository
}
